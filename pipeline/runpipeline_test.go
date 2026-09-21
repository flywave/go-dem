package pipeline

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/uncertainty"
	"github.com/flywave/go-geo"
)

type progressCall struct {
	stage       string
	done, total int
}

func pipelineRegion() *dem.Region {
	return dem.NewRegionFromBBox(-122.5, 40.0, -122.4, 40.1, geo.NewProj("EPSG:4326"), 0.02, 0.02)
}

func writeTestXYZ(t *testing.T, path string) {
	t.Helper()
	var sb strings.Builder
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			lon := -122.5 + (float64(x)+0.5)*0.02
			lat := 40.1 - (float64(y)+0.5)*0.02
			z := 100.0 + float64(x) + 2*float64(y)
			sb.WriteString(fmt.Sprintf("%f %f %f\n", lon, lat, z))
		}
	}
	if err := os.WriteFile(path, []byte(sb.String()), 0o644); err != nil {
		t.Fatalf("write xyz: %v", err)
	}
}

func pipelineTestInput(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "points.xyz")
	writeTestXYZ(t, path)
	return path
}

func fullPipelineConfig(outDir string) *PipelineConfig {
	return &PipelineConfig{
		Config:    Config{SourceEpsg: 4326},
		Method:    dem.MethodIDW,
		Filters:   []string{"fill"},
		Hillshade: true,
		Histogram: true,
		OutputDir: outDir,
	}
}

func TestRunPipelineFullRun(t *testing.T) {
	input := pipelineTestInput(t)
	outDir := t.TempDir()
	region := pipelineRegion()

	res, err := RunPipeline([]string{input}, region, fullPipelineConfig(outDir), nil, context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if res.PointCount != 25 {
		t.Errorf("PointCount = %d, want 25", res.PointCount)
	}
	if res.Region == nil {
		t.Error("result region must not be nil")
	}

	wantDEM := filepath.Join(outDir, "dem.tif")
	wantShade := filepath.Join(outDir, "hillshade.tif")
	wantHist := filepath.Join(outDir, "histogram.png")
	if res.DEMPath != wantDEM {
		t.Errorf("DEMPath = %s, want %s", res.DEMPath, wantDEM)
	}
	if res.HillshadePath != wantShade {
		t.Errorf("HillshadePath = %s, want %s", res.HillshadePath, wantShade)
	}
	if res.HistogramPath != wantHist {
		t.Errorf("HistogramPath = %s, want %s", res.HistogramPath, wantHist)
	}
	for _, p := range []string{res.DEMPath, res.HillshadePath, res.HistogramPath} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("output %s: %v", p, err)
		}
	}

	data, outRegion, err := dem.ReadDEM(res.DEMPath)
	if err != nil {
		t.Fatal(err)
	}
	if outRegion.XSize != region.XSize || outRegion.YSize != region.YSize {
		t.Fatalf("output region is %dx%d, want %dx%d", outRegion.XSize, outRegion.YSize, region.XSize, region.YSize)
	}
	for y := 0; y < region.YSize; y++ {
		for x := 0; x < region.XSize; x++ {
			want := 100.0 + float64(x) + 2*float64(y)
			got := data[y*region.XSize+x]
			if math.Abs(got-want) > 1e-9 {
				t.Errorf("pixel(%d,%d) = %f, want %f", x, y, got, want)
			}
		}
	}
}

func TestRunPipelineProgressMonotonic(t *testing.T) {
	input := pipelineTestInput(t)
	outDir := t.TempDir()

	cfg := fullPipelineConfig(outDir)
	cfg.Uncertainty = &uncertainty.Options{Method: uncertainty.MethodProximity}

	var calls []progressCall
	progress := func(stage string, done, total int) {
		calls = append(calls, progressCall{stage: stage, done: done, total: total})
	}

	_, err := RunPipeline([]string{input}, pipelineRegion(), cfg, progress, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(calls) == 0 {
		t.Fatal("no progress callbacks collected")
	}

	last := 0
	for i, c := range calls {
		if c.total != 100 {
			t.Errorf("callback %d (%s): total = %d, want 100", i, c.stage, c.total)
		}
		if c.done < last {
			t.Errorf("callback %d (%s): done = %d, decreased from %d", i, c.stage, c.done, last)
		}
		last = c.done
	}
	if last != 100 {
		t.Errorf("final done = %d, want 100", last)
	}
	found := false
	for _, c := range calls {
		if c.done == 100 && c.total == 100 {
			found = true
		}
	}
	if !found {
		t.Error("no callback reported done == 100")
	}
}

func TestRunPipelineCanceledContext(t *testing.T) {
	input := pipelineTestInput(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := RunPipeline([]string{input}, pipelineRegion(), fullPipelineConfig(t.TempDir()), nil, ctx)
	if err == nil {
		t.Fatal("canceled context must fail the run")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}
}

func TestRunPipelineUncertaintyOutputs(t *testing.T) {
	input := pipelineTestInput(t)
	outDir := t.TempDir()

	cfg := fullPipelineConfig(outDir)
	cfg.Filters = nil
	cfg.Hillshade = false
	cfg.Histogram = false
	cfg.Uncertainty = &uncertainty.Options{Method: uncertainty.MethodProximity}

	res, err := RunPipeline([]string{input}, pipelineRegion(), cfg, nil, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.UncertaintyPaths) == 0 {
		t.Fatal("UncertaintyPaths is empty")
	}
	if filepath.Base(res.UncertaintyPaths[0]) != "dem_tvu.tif" {
		t.Errorf("first uncertainty path = %s, want dem_tvu.tif", res.UncertaintyPaths[0])
	}
	for _, p := range res.UncertaintyPaths {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("uncertainty output %s: %v", p, err)
		}
	}
}

func TestRunPipelineValidationErrors(t *testing.T) {
	region := pipelineRegion()
	input := pipelineTestInput(t)

	if _, err := RunPipeline(nil, region, fullPipelineConfig(t.TempDir()), nil, context.Background()); err == nil {
		t.Error("empty inputs must fail")
	}
	if _, err := RunPipeline([]string{filepath.Join(t.TempDir(), "missing.xyz")}, region, fullPipelineConfig(t.TempDir()), nil, context.Background()); err == nil {
		t.Error("missing input file must fail")
	}
	if _, err := RunPipeline([]string{input}, nil, fullPipelineConfig(t.TempDir()), nil, context.Background()); err == nil {
		t.Error("nil region must fail")
	}
	if _, err := RunPipeline([]string{input}, region, nil, nil, context.Background()); err == nil {
		t.Error("nil config must fail")
	}
	if _, err := RunPipeline([]string{input}, region, &PipelineConfig{Config: Config{SourceEpsg: 4326}}, nil, context.Background()); err == nil {
		t.Error("empty OutputDir must fail")
	}
}

func TestRunPipelineUnknownFilter(t *testing.T) {
	input := pipelineTestInput(t)
	cfg := fullPipelineConfig(t.TempDir())
	cfg.Filters = []string{"nope"}

	_, err := RunPipeline([]string{input}, pipelineRegion(), cfg, nil, context.Background())
	if err == nil {
		t.Fatal("unknown filter must fail")
	}
	if !strings.Contains(err.Error(), "unknown filter") {
		t.Errorf("error = %v, want unknown filter message", err)
	}
}

func TestRunPipelineNilRegionSRSDefaultsTo4326(t *testing.T) {
	input := pipelineTestInput(t)
	outDir := t.TempDir()
	region := dem.NewRegionFromBBox(-122.5, 40.0, -122.4, 40.1, nil, 0.02, 0.02)

	res, err := RunPipeline([]string{input}, region, fullPipelineConfig(outDir), nil, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res.Region.SRS() == nil || res.Region.SRS().GetSrsCode() != "EPSG:4326" {
		t.Errorf("resolved SRS = %v, want EPSG:4326", res.Region.SRS())
	}
	if !nilSRS(region.SRS()) {
		t.Error("input region must not be mutated by the SRS default")
	}
	data, outRegion, err := dem.ReadDEM(res.DEMPath)
	if err != nil {
		t.Fatal(err)
	}
	if outRegion.XSize != 5 || outRegion.YSize != 5 {
		t.Fatalf("output region is %dx%d, want 5x5", outRegion.XSize, outRegion.YSize)
	}
	for i, v := range data {
		if dem.IsNoData(v, dem.DefaultNoData) {
			t.Fatalf("pixel %d is nodata, want valid elevation", i)
		}
	}
}
