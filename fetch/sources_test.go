package fetch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func overrideURL(t *testing.T, target *string, value string) {
	t.Helper()
	old := *target
	*target = value
	t.Cleanup(func() { *target = old })
}

func TestSRTMFetch(t *testing.T) {
	payload := []byte("srtm geotiff payload")
	var mu sync.Mutex
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		query = r.URL.Query()
		mu.Unlock()
		w.Write(payload)
	}))
	defer srv.Close()
	overrideURL(t, &opentopographyGlobalDEMURL, srv.URL)

	dir := t.TempDir()
	opts := testOptions(dir)
	opts.APIKey = "ot-key-123"

	files, err := Fetch(SourceSRTM, opts)
	if err != nil {
		t.Fatalf("srtm fetch: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("got %d files, want 1", len(files))
	}
	if files[0] != filepath.Join(dir, "srtm_dem.tif") {
		t.Errorf("output path = %q, want %q", files[0], filepath.Join(dir, "srtm_dem.tif"))
	}
	got, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if string(got) != string(payload) {
		t.Errorf("output content = %q, want %q", got, payload)
	}
	if _, err := os.Stat(files[0] + ".part"); !os.IsNotExist(err) {
		t.Errorf("temp file left behind")
	}

	mu.Lock()
	q := query
	mu.Unlock()
	for k, want := range map[string]string{
		"demtype":      "SRTMGL1",
		"west":         "-122.500000",
		"south":        "37.500000",
		"east":         "-122.000000",
		"north":        "38.000000",
		"outputFormat": "GTiff",
		"API_Key":      "ot-key-123",
	} {
		if q.Get(k) != want {
			t.Errorf("query %s = %q, want %q", k, q.Get(k), want)
		}
	}
}

func TestSRTMFetchWithoutAPIKey(t *testing.T) {
	var mu sync.Mutex
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		query = r.URL.Query()
		mu.Unlock()
		w.Write([]byte("x"))
	}))
	defer srv.Close()
	overrideURL(t, &opentopographyGlobalDEMURL, srv.URL)

	if _, err := Fetch(SourceSRTM, testOptions(t.TempDir())); err != nil {
		t.Fatalf("srtm fetch: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if _, ok := query["API_Key"]; ok {
		t.Errorf("API_Key should be omitted when opts.APIKey is empty, got %q", query.Get("API_Key"))
	}
}

func TestUSGSTNMTwoStageFetch(t *testing.T) {
	tifA := []byte("tif-a-payload")
	tifB := []byte("tif-b-payload")

	var mu sync.Mutex
	var productsQuery url.Values
	var downloads []string

	var baseURL string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/products", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		productsQuery = r.URL.Query()
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w,
			`{"total":3,"items":[`+
				`{"title":"Product A","format":"GeoTIFF","downloadURL":"%s/data/a.tif"},`+
				`{"title":"Product B","format":"GeoTIFF","downloadURL":"%s/data/b.tif"},`+
				`{"title":"No download URL"}]}`,
			baseURL, baseURL)
	})
	mux.HandleFunc("/data/a.tif", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		downloads = append(downloads, "a")
		mu.Unlock()
		w.Write(tifA)
	})
	mux.HandleFunc("/data/b.tif", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		downloads = append(downloads, "b")
		mu.Unlock()
		w.Write(tifB)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	baseURL = srv.URL
	overrideURL(t, &tnmProductsURL, srv.URL+"/api/v1/products")

	dir := t.TempDir()
	files, err := Fetch(SourceUSGSTNM, testOptions(dir))
	if err != nil {
		t.Fatalf("usgs tnm fetch: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("got %d files, want 2", len(files))
	}
	wantFiles := []string{
		filepath.Join(dir, "usgs_tnm_0.tif"),
		filepath.Join(dir, "usgs_tnm_1.tif"),
	}
	for i, want := range wantFiles {
		if files[i] != want {
			t.Errorf("files[%d] = %q, want %q", i, files[i], want)
		}
	}
	wantContent := [][]byte{tifA, tifB}
	for i, want := range wantContent {
		got, err := os.ReadFile(files[i])
		if err != nil {
			t.Fatalf("read %s: %v", files[i], err)
		}
		if string(got) != string(want) {
			t.Errorf("content of %s = %q, want %q", files[i], got, want)
		}
	}

	mu.Lock()
	q := productsQuery
	mu.Unlock()
	if q.Get("bbox") != "37.500000,-122.500000,38.000000,-122.000000" {
		t.Errorf("bbox = %q, want south,west,north,east ordering", q.Get("bbox"))
	}
	if q.Get("prodFormats") != "GeoTIFF" {
		t.Errorf("prodFormats = %q", q.Get("prodFormats"))
	}
	if q.Get("datasets") != tnmNED1ArcSecond {
		t.Errorf("datasets = %q, want %q", q.Get("datasets"), tnmNED1ArcSecond)
	}
	if q.Get("outputFormat") != "JSON" {
		t.Errorf("outputFormat = %q", q.Get("outputFormat"))
	}

	mu.Lock()
	defer mu.Unlock()
	if len(downloads) != 2 {
		t.Errorf("downloaded %d products, want 2", len(downloads))
	}
}

func TestUSGSTNMNoProducts(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"total":0,"items":[]}`))
	}))
	defer srv.Close()
	overrideURL(t, &tnmProductsURL, srv.URL)

	dir := t.TempDir()
	_, err := Fetch(SourceUSGSTNM, testOptions(dir))
	if err == nil {
		t.Fatal("expected error when no products match")
	}
	if !strings.Contains(err.Error(), "no GeoTIFF products") {
		t.Errorf("error = %v, want a no-products error", err)
	}
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		t.Fatalf("read dir: %v", readErr)
	}
	if len(entries) != 0 {
		t.Errorf("output dir not empty after failure: %v", entries)
	}
}

func TestUSGSTNMInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html>not json</html>"))
	}))
	defer srv.Close()
	overrideURL(t, &tnmProductsURL, srv.URL)

	_, err := Fetch(SourceUSGSTNM, testOptions(t.TempDir()))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "decode products") {
		t.Errorf("error = %v, want a JSON decode error", err)
	}
}

func TestUSGSTNMDownloadFailurePropagates(t *testing.T) {
	var baseURL string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/products", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"total":1,"items":[{"downloadURL":"%s/data/broken.tif"}]}`, baseURL)
	})
	mux.HandleFunc("/data/broken.tif", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	baseURL = srv.URL
	overrideURL(t, &tnmProductsURL, srv.URL+"/api/v1/products")

	_, err := Fetch(SourceUSGSTNM, testOptions(t.TempDir()))
	if err == nil {
		t.Fatal("expected error when product download fails")
	}
	if !strings.Contains(err.Error(), "usgs tnm") {
		t.Errorf("error = %v, want it to be wrapped with the source name", err)
	}
}

func TestArcticDEMFetch(t *testing.T) {
	payload := []byte("arcticdem geotiff payload")
	var mu sync.Mutex
	var query url.Values
	var authorization string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		query = r.URL.Query()
		authorization = r.Header.Get("Authorization")
		mu.Unlock()
		w.Write(payload)
	}))
	defer srv.Close()
	overrideURL(t, &pgcArcticDEMURL, srv.URL)

	dir := t.TempDir()
	opts := testOptions(dir)
	opts.Token = "Bearer pgc-token"

	files, err := Fetch(SourceArcticDEM, opts)
	if err != nil {
		t.Fatalf("arctic dem fetch: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("got %d files, want 1", len(files))
	}
	if files[0] != filepath.Join(dir, "arctic_dem.tif") {
		t.Errorf("output path = %q, want %q", files[0], filepath.Join(dir, "arctic_dem.tif"))
	}
	got, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if string(got) != string(payload) {
		t.Errorf("output content = %q, want %q", got, payload)
	}

	mu.Lock()
	defer mu.Unlock()
	if query.Get("bbox") != "-122.500000,37.500000,-122.000000,38.000000" {
		t.Errorf("bbox = %q, want west,south,east,north ordering", query.Get("bbox"))
	}
	if query.Get("output_format") != "GeoTIFF" {
		t.Errorf("output_format = %q", query.Get("output_format"))
	}
	if authorization != "Bearer pgc-token" {
		t.Errorf("Authorization = %q, want %q", authorization, "Bearer pgc-token")
	}
}

func TestArcticDEMFetchWithoutToken(t *testing.T) {
	var mu sync.Mutex
	var authorization string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		authorization = r.Header.Get("Authorization")
		mu.Unlock()
		w.Write([]byte("x"))
	}))
	defer srv.Close()
	overrideURL(t, &pgcArcticDEMURL, srv.URL)

	if _, err := Fetch(SourceArcticDEM, testOptions(t.TempDir())); err != nil {
		t.Fatalf("arctic dem fetch: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if authorization != "" {
		t.Errorf("Authorization = %q, want empty when Token is not set", authorization)
	}
}
