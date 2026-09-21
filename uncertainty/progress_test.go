package uncertainty

import (
	"context"
	"errors"
	"testing"
)

type uncProgressEvent struct {
	stage string
	done  int
	total int
}

type uncProgressRecorder struct {
	events []uncProgressEvent
}

func (r *uncProgressRecorder) record(stage string, done, total int) {
	r.events = append(r.events, uncProgressEvent{stage, done, total})
}

func assertUncProgress(t *testing.T, rec *uncProgressRecorder, stage string) {
	t.Helper()
	if len(rec.events) == 0 {
		t.Fatal("no progress callbacks")
	}
	last := make(map[int]int)
	for _, e := range rec.events {
		if e.stage != stage {
			t.Fatalf("stage %q, want %q", e.stage, stage)
		}
		if e.done < last[e.total] {
			t.Fatalf("done not monotonic: %d after %d (total %d)", e.done, last[e.total], e.total)
		}
		last[e.total] = e.done
	}
	final := rec.events[len(rec.events)-1]
	if final.done != final.total {
		t.Fatalf("final callback done=%d total=%d, want done==total", final.done, final.total)
	}
}

func uncTestData(w, h int) []float64 {
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = 100 + float64(x%3) + float64(y%2)
		}
	}
	return data
}

func TestUncertaintyProgressSplitSample(t *testing.T) {
	w, h := 20, 20
	rec := &uncProgressRecorder{}
	opts := &Options{Method: MethodSplitSample, NoData: -9999, SampleFraction: 0.3, Progress: rec.record}
	res, err := Estimate(uncTestData(w, h), region(w, h), opts)
	if err != nil {
		t.Fatalf("split_sample: %v", err)
	}
	if res == nil {
		t.Fatal("nil result")
	}
	assertUncProgress(t, rec, "split_sample")
}

func TestUncertaintyProgressProximity(t *testing.T) {
	w, h := 12, 12
	rec := &uncProgressRecorder{}
	opts := &Options{Method: MethodProximity, NoData: -9999, Progress: rec.record}
	if _, err := Estimate(uncTestData(w, h), region(w, h), opts); err != nil {
		t.Fatalf("proximity: %v", err)
	}
	assertUncProgress(t, rec, "proximity")
}

func TestUncertaintyCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	w, h := 12, 12
	data := uncTestData(w, h)
	reg := region(w, h)
	for _, m := range []Method{MethodSplitSample, MethodProximity, MethodCombined} {
		res, err := Estimate(data, reg, &Options{Method: m, NoData: -9999, Ctx: ctx})
		if !errors.Is(err, context.Canceled) {
			t.Errorf("%s: expected context.Canceled, got %v", m, err)
		}
		if res != nil {
			t.Errorf("%s: expected nil result on cancel", m)
		}
	}
}
