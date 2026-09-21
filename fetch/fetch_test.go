package fetch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

func testOptions(dir string) *FetchOptions {
	return &FetchOptions{
		MinX: -122.5, MinY: 37.5, MaxX: -122.0, MaxY: 38.0,
		OutputDir:  dir,
		MaxRetries: 1,
	}
}

type stubFetcher struct {
	baseFetcher
	name DataSource
}

func (f *stubFetcher) Name() DataSource { return f.name }

func (f *stubFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	return []string{"stub"}, nil
}

func TestBBoxParam(t *testing.T) {
	b := newBaseFetcher("test")
	got := b.bboxParam(-180, -45, 90, 80)
	want := "-180.000000,-45.000000,90.000000,80.000000"
	if got != want {
		t.Errorf("bboxParam = %q, want %q", got, want)
	}
}

func TestListSources(t *testing.T) {
	sources := ListSources()
	if !sort.SliceIsSorted(sources, func(i, j int) bool { return sources[i] < sources[j] }) {
		t.Errorf("ListSources not sorted: %v", sources)
	}

	seen := make(map[DataSource]bool)
	for _, s := range sources {
		if seen[s] {
			t.Errorf("duplicate source %q", s)
		}
		seen[s] = true
	}

	for _, want := range []DataSource{
		SourceNOAAMultibeam, SourceUSGSTNM, SourceCopernicus, SourceGEBCO,
		SourceSRTM, SourceEMODNet, SourceArcticDEM,
	} {
		if !seen[want] {
			t.Errorf("ListSources missing %q, got %v", want, sources)
		}
	}
}

func TestFetchUnknownSource(t *testing.T) {
	if _, err := Fetch("not_a_source", nil); err == nil {
		t.Fatal("expected error for unknown source")
	}
}

func TestRegisterConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		name := DataSource(fmt.Sprintf("test_source_%d", i))
		wg.Add(1)
		go func() {
			defer wg.Done()
			Register(&stubFetcher{name: name, baseFetcher: newBaseFetcher(name)})
		}()
	}
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ListSources()
		}()
	}
	wg.Wait()

	files, err := Fetch("test_source_7", nil)
	if err != nil {
		t.Fatalf("Fetch registered stub: %v", err)
	}
	if len(files) != 1 || files[0] != "stub" {
		t.Errorf("stub Fetch returned %v", files)
	}
}

func TestClientFor(t *testing.T) {
	b := newBaseFetcher("test")

	opts := &FetchOptions{}
	if got := b.clientFor(opts); got != b.client {
		t.Errorf("clientFor with zero Timeout should return the default client")
	}

	opts.Timeout = 5 * time.Second
	got := b.clientFor(opts)
	if got == b.client {
		t.Errorf("clientFor with positive Timeout should not return the default client")
	}
	if got.Timeout != 5*time.Second {
		t.Errorf("clientFor Timeout = %v, want %v", got.Timeout, 5*time.Second)
	}
}

func TestDownloadSuccess(t *testing.T) {
	payload := []byte("fake geotiff payload")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(payload)
	}))
	defer srv.Close()

	dir := t.TempDir()
	outputPath := filepath.Join(dir, "data.tif")

	b := newBaseFetcher("test")
	if err := b.download(srv.URL+"/data.tif", outputPath, testOptions(dir)); err != nil {
		t.Fatalf("download: %v", err)
	}

	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read downloaded file: %v", err)
	}
	if string(got) != string(payload) {
		t.Errorf("file content = %q, want %q", got, payload)
	}
	if _, err := os.Stat(outputPath + ".part"); !os.IsNotExist(err) {
		t.Errorf("temp file %s still exists", outputPath+".part")
	}
}

func TestDownloadFailureLeavesNoFiles(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	dir := t.TempDir()
	outputPath := filepath.Join(dir, "data.tif")

	b := newBaseFetcher("test")
	err := b.download(srv.URL+"/data.tif", outputPath, testOptions(dir))
	if err == nil {
		t.Fatal("expected error when server returns 500")
	}
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		t.Errorf("output file exists after failed download")
	}
	if _, err := os.Stat(outputPath + ".part"); !os.IsNotExist(err) {
		t.Errorf("temp file %s still exists after failed download", outputPath+".part")
	}
}

func TestDownloadRetryUntilSuccess(t *testing.T) {
	var mu sync.Mutex
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		hits++
		n := hits
		mu.Unlock()
		if n < 2 {
			http.Error(w, "flaky", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	dir := t.TempDir()
	outputPath := filepath.Join(dir, "data.tif")

	opts := testOptions(dir)
	opts.MaxRetries = 2
	b := newBaseFetcher("test")
	if err := b.download(srv.URL+"/data.tif", outputPath, opts); err != nil {
		t.Fatalf("download: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if hits != 2 {
		t.Errorf("attempts = %d, want 2 (MaxRetries is the total number of attempts)", hits)
	}
}

func TestDownloadMaxRetriesDefault(t *testing.T) {
	var mu sync.Mutex
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		hits++
		mu.Unlock()
		http.Error(w, "always failing", http.StatusInternalServerError)
	}))
	defer srv.Close()

	dir := t.TempDir()
	opts := testOptions(dir)
	opts.MaxRetries = 0

	b := newBaseFetcher("test")
	err := b.download(srv.URL+"/data.tif", filepath.Join(dir, "data.tif"), opts)
	if err == nil {
		t.Fatal("expected error when all attempts fail")
	}

	mu.Lock()
	defer mu.Unlock()
	if hits != 3 {
		t.Errorf("attempts = %d, want 3 (default MaxRetries)", hits)
	}
	if !strings.Contains(err.Error(), "after 3 attempts") {
		t.Errorf("error %q should state the total number of attempts", err)
	}
}

func TestDownloadTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte("late"))
	}))
	defer srv.Close()

	dir := t.TempDir()
	outputPath := filepath.Join(dir, "slow.tif")

	opts := testOptions(dir)
	opts.Timeout = 50 * time.Millisecond

	b := newBaseFetcher("test")
	start := time.Now()
	err := b.download(srv.URL+"/slow.tif", outputPath, opts)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("download took %v; Timeout option was not applied", elapsed)
	}
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		t.Errorf("output file exists after timeout")
	}
	if _, err := os.Stat(outputPath + ".part"); !os.IsNotExist(err) {
		t.Errorf("temp file %s still exists after timeout", outputPath+".part")
	}
}

func TestDownloadHeaders(t *testing.T) {
	var mu sync.Mutex
	var userAgent, authorization string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		userAgent = r.Header.Get("User-Agent")
		authorization = r.Header.Get("Authorization")
		mu.Unlock()
		w.Write([]byte("x"))
	}))
	defer srv.Close()

	dir := t.TempDir()
	b := newBaseFetcher("test")

	opts := testOptions(dir)
	opts.Token = "Bearer test-token"
	if err := b.download(srv.URL+"/a.tif", filepath.Join(dir, "a.tif"), opts); err != nil {
		t.Fatalf("download: %v", err)
	}
	mu.Lock()
	if userAgent != "go-dem-fetcher/1.0" {
		t.Errorf("User-Agent = %q, want default", userAgent)
	}
	if authorization != "Bearer test-token" {
		t.Errorf("Authorization = %q, want %q", authorization, "Bearer test-token")
	}
	mu.Unlock()

	opts = testOptions(dir)
	opts.UserAgent = "custom-agent/2.0"
	if err := b.download(srv.URL+"/b.tif", filepath.Join(dir, "b.tif"), opts); err != nil {
		t.Fatalf("download: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if userAgent != "custom-agent/2.0" {
		t.Errorf("User-Agent = %q, want custom-agent/2.0", userAgent)
	}
	if authorization != "" {
		t.Errorf("Authorization = %q, want empty when Token is not set", authorization)
	}
}
