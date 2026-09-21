package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type DataSource string

const (
	SourceNOAAMultibeam DataSource = "noaa_multibeam"
	SourceUSGSTNM       DataSource = "usgs_tnm"
	SourceCopernicus    DataSource = "copernicus"
	SourceGEBCO         DataSource = "gebco"
	SourceSRTM          DataSource = "srtm"
	SourceEMODNet       DataSource = "emodnet"
	SourceArcticDEM     DataSource = "arctic_dem"
)

type FetchOptions struct {
	Source     DataSource
	MinX, MinY float64
	MaxX, MaxY float64
	OutputDir  string
	// MaxRetries is the total number of HTTP attempts per request; <= 0 means 3.
	MaxRetries int
	// Timeout bounds each HTTP request (connect plus body read) when positive.
	Timeout time.Duration
	// APIKey is the OpenTopography API key, sent as the API_Key query parameter.
	APIKey    string
	UserAgent string
	// Token is sent verbatim as the Authorization header of every request
	// (e.g. "Bearer <token>" or "Basic <token>" for the PGC ArcticDEM API).
	Token string
}

type Fetcher interface {
	Name() DataSource
	Fetch(opts *FetchOptions) ([]string, error)
}

type Registry map[DataSource]Fetcher

var (
	registryMu sync.RWMutex
	registry   = make(Registry)
)

func Register(f Fetcher) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[f.Name()] = f
}

func Fetch(source DataSource, opts *FetchOptions) ([]string, error) {
	registryMu.RLock()
	f, ok := registry[source]
	registryMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown data source: %s", source)
	}
	return f.Fetch(opts)
}

func ListSources() []DataSource {
	registryMu.RLock()
	sources := make([]DataSource, 0, len(registry))
	for s := range registry {
		sources = append(sources, s)
	}
	registryMu.RUnlock()

	sort.Slice(sources, func(i, j int) bool { return sources[i] < sources[j] })
	return sources
}

type baseFetcher struct {
	name   DataSource
	client *http.Client
}

func newBaseFetcher(name DataSource) baseFetcher {
	return baseFetcher{
		name: name,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (b *baseFetcher) clientFor(opts *FetchOptions) *http.Client {
	if opts != nil && opts.Timeout > 0 {
		return &http.Client{Timeout: opts.Timeout}
	}
	return b.client
}

func (b *baseFetcher) request(url string, opts *FetchOptions) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %v", err)
	}

	userAgent := opts.UserAgent
	if userAgent == "" {
		userAgent = "go-dem-fetcher/1.0"
	}
	req.Header.Set("User-Agent", userAgent)
	if opts.Token != "" {
		req.Header.Set("Authorization", opts.Token)
	}

	resp, err := b.clientFor(opts).Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}
	return resp, nil
}

func withRetries(opts *FetchOptions, fn func() error) error {
	attempts := opts.MaxRetries
	if attempts <= 0 {
		attempts = 3
	}

	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return fmt.Errorf("failed after %d attempts: %v", attempts, lastErr)
}

func (b *baseFetcher) get(url string, opts *FetchOptions) (*http.Response, error) {
	var resp *http.Response
	err := withRetries(opts, func() error {
		r, err := b.request(url, opts)
		if err != nil {
			return err
		}
		resp = r
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *baseFetcher) download(url, outputPath string, opts *FetchOptions) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create dir: %v", err)
	}

	partPath := outputPath + ".part"
	err := withRetries(opts, func() error {
		return b.saveResponse(url, partPath, opts)
	})
	if err != nil {
		os.Remove(partPath)
		return fmt.Errorf("download %s: %v", url, err)
	}

	if err := os.Rename(partPath, outputPath); err != nil {
		os.Remove(partPath)
		return fmt.Errorf("rename %s: %v", partPath, err)
	}
	return nil
}

func (b *baseFetcher) saveResponse(url, path string, opts *FetchOptions) error {
	resp, err := b.request(url, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %v", err)
	}

	_, copyErr := io.Copy(out, resp.Body)
	closeErr := out.Close()
	if copyErr != nil {
		return fmt.Errorf("write file: %v", copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close file: %v", closeErr)
	}
	return nil
}

// bboxParam formats a bbox as "west,south,east,north"
// (minLon,minLat,maxLon,maxLat), the ordering used by most REST/WMS services.
func (b *baseFetcher) bboxParam(west, south, east, north float64) string {
	return fmt.Sprintf("%.6f,%.6f,%.6f,%.6f", west, south, east, north)
}
