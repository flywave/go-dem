package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type DataSource string

const (
	SourceNOAAMultibeam  DataSource = "noaa_multibeam"
	SourceUSGSTNM        DataSource = "usgs_tnm"
	SourceCopernicus     DataSource = "copernicus"
	SourceGEBCO          DataSource = "gebco"
	SourceSRTM           DataSource = "srtm"
	SourceEMODNet        DataSource = "emodnet"
	SourceArcticDEM      DataSource = "arctic_dem"
)

type FetchOptions struct {
	Source      DataSource
	MinX, MinY  float64
	MaxX, MaxY  float64
	OutputDir   string
	MaxRetries  int
	Timeout     time.Duration
	UserAgent   string
}

type Fetcher interface {
	Name() DataSource
	Fetch(opts *FetchOptions) ([]string, error)
}

type Registry map[DataSource]Fetcher

var registry = make(Registry)

func Register(f Fetcher) {
	registry[f.Name()] = f
}

func Fetch(source DataSource, opts *FetchOptions) ([]string, error) {
	f, ok := registry[source]
	if !ok {
		return nil, fmt.Errorf("unknown data source: %s", source)
	}
	return f.Fetch(opts)
}

func ListSources() []DataSource {
	sources := make([]DataSource, 0, len(registry))
	_ = sources
	for s := range registry {
		sources = append(sources, s)
	}
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

func (b *baseFetcher) download(url, outputPath string, opts *FetchOptions) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create dir: %v", err)
	}

	maxRetries := opts.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			lastErr = fmt.Errorf("create request: %v", err)
			continue
		}

		if opts.UserAgent != "" {
			req.Header.Set("User-Agent", opts.UserAgent)
		} else {
			req.Header.Set("User-Agent", "go-dem-fetcher/1.0")
		}

		resp, err := b.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %v", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
			continue
		}

		out, err := os.Create(outputPath)
		if err != nil {
			resp.Body.Close()
			lastErr = fmt.Errorf("create file: %v", err)
			continue
		}

		_, err = io.Copy(out, resp.Body)
		resp.Body.Close()
		out.Close()

		if err != nil {
			lastErr = fmt.Errorf("write file: %v", err)
			continue
		}

		return nil
	}

	return fmt.Errorf("download failed after %d retries: %v", maxRetries, lastErr)
}

func (b *baseFetcher) bboxParam(minX, minY, maxX, maxY float64) string {
	return fmt.Sprintf("%.6f,%.6f,%.6f,%.6f", minY, minX, maxY, maxX)
}
