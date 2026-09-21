package grits

import (
	"context"
	"errors"
	"testing"

	"github.com/flywave/go-dem"
)

type gritsProgressEvent struct {
	stage string
	done  int
	total int
}

type gritsProgressRecorder struct {
	events []gritsProgressEvent
}

func (r *gritsProgressRecorder) record(stage string, done, total int) {
	r.events = append(r.events, gritsProgressEvent{stage, done, total})
}

func assertGritsProgress(t *testing.T, rec *gritsProgressRecorder, stage string, h int) {
	t.Helper()
	if len(rec.events) == 0 {
		t.Fatal("no progress callbacks")
	}
	last := make(map[int]int)
	for _, e := range rec.events {
		if e.stage != stage {
			t.Fatalf("stage %q, want %q", e.stage, stage)
		}
		if e.total != h {
			t.Fatalf("total %d, want %d", e.total, h)
		}
		if e.done < last[e.total] {
			t.Fatalf("done not monotonic: %d after %d", e.done, last[e.total])
		}
		last[e.total] = e.done
	}
	final := rec.events[len(rec.events)-1]
	if final.done != final.total {
		t.Fatalf("final callback done=%d total=%d, want done==total", final.done, final.total)
	}
}

func gritsTestRegion() *dem.Region {
	return dem.NewRegionFromBBox(0, 0, 8, 8, nil, 1, 1)
}

func gritsTestData() []float64 {
	w, h := 8, 8
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = 100 + float64(x) + float64(y)
		}
	}
	data[3*w+3] = dem.DefaultNoData
	data[4*w+4] = dem.DefaultNoData
	return data
}

func TestGritsProgressGaussian(t *testing.T) {
	rec := &gritsProgressRecorder{}
	f, err := New(FilterGaussian)
	if err != nil {
		t.Fatal(err)
	}
	opts := &Options{Progress: rec.record}
	out, err := f.Run(gritsTestData(), gritsTestRegion(), opts)
	if err != nil {
		t.Fatalf("gaussian: %v", err)
	}
	if len(out) != 64 {
		t.Fatalf("output size %d", len(out))
	}
	assertGritsProgress(t, rec, "gaussian", 8)
}

func TestGritsProgressOpenTwoPhases(t *testing.T) {
	rec := &gritsProgressRecorder{}
	f, err := New(FilterOpen)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Run(gritsTestData(), gritsTestRegion(), &Options{Progress: rec.record}); err != nil {
		t.Fatalf("open: %v", err)
	}
	assertGritsProgress(t, rec, "open", 8)
}

func TestGritsProgressOutliersMultipass(t *testing.T) {
	rec := &gritsProgressRecorder{}
	f, err := New(FilterOutliers)
	if err != nil {
		t.Fatal(err)
	}
	opts := &Options{Progress: rec.record, Iterations: 2}
	if _, err := f.Run(gritsTestData(), gritsTestRegion(), opts); err != nil {
		t.Fatalf("outliers: %v", err)
	}
	assertGritsProgress(t, rec, "outliers", 8)
}

func TestGritsProgressZScore(t *testing.T) {
	rec := &gritsProgressRecorder{}
	f, err := New(FilterZScore)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Run(gritsTestData(), gritsTestRegion(), &Options{Progress: rec.record}); err != nil {
		t.Fatalf("zscore: %v", err)
	}
	assertGritsProgress(t, rec, "zscore", 8)
}

func TestGritsProgressMedian(t *testing.T) {
	rec := &gritsProgressRecorder{}
	f, err := New(FilterMedian)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Run(gritsTestData(), gritsTestRegion(), &Options{Progress: rec.record}); err != nil {
		t.Fatalf("median: %v", err)
	}
	assertGritsProgress(t, rec, "median", 8)
}

func TestGritsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	reg := gritsTestRegion()
	data := gritsTestData()
	for _, ft := range ListFilters() {
		if ft == FilterBlend || ft == "diff" {
			continue
		}
		f, err := New(ft)
		if err != nil {
			t.Fatalf("new %s: %v", ft, err)
		}
		opts := &Options{
			Ctx:        ctx,
			PolygonWKT: "POLYGON((1 1,7 1,7 7,1 7,1 1))",
			CutBounds:  []float64{1, 1, 7, 7},
		}
		res, err := f.Run(data, reg, opts)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("%s: expected context.Canceled, got %v", ft, err)
		}
		if res != nil {
			t.Errorf("%s: expected nil result on cancel", ft)
		}
	}
}
