package datum

import (
	"context"
	"errors"
	"testing"
)

type datumProgressEvent struct {
	stage string
	done  int
	total int
}

type datumProgressRecorder struct {
	events []datumProgressEvent
}

func (r *datumProgressRecorder) record(stage string, done, total int) {
	r.events = append(r.events, datumProgressEvent{stage, done, total})
}

func TestDatumProgressRunStepReports(t *testing.T) {
	rec := &datumProgressRecorder{}
	vt := NewVerticalTransform(TransformOptions{
		EpsgIn:   5773,
		EpsgOut:  3855,
		Region:   testRegion(),
		Progress: rec.record,
	})
	if _, err := vt.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(rec.events) == 0 {
		t.Fatal("no progress callbacks")
	}
	for _, e := range rec.events {
		if e.stage != "cdn2cdn" {
			t.Fatalf("stage %q, want cdn2cdn", e.stage)
		}
	}
	final := rec.events[len(rec.events)-1]
	if final.done != final.total || final.total != 1 {
		t.Fatalf("final callback %+v, want (cdn2cdn,1,1)", final)
	}
	last := -1
	for _, e := range rec.events {
		if e.total != 1 {
			continue
		}
		if e.done < last {
			t.Fatalf("step reports not monotonic: %d after %d", e.done, last)
		}
		last = e.done
	}
	hasRows := false
	for _, e := range rec.events {
		if e.total == testRegion().YSize {
			hasRows = true
		}
	}
	if !hasRows {
		t.Error("expected per-row geoid grid callbacks")
	}
}

func TestDatumProgressGeoidRowsMonotonic(t *testing.T) {
	rec := &datumProgressRecorder{}
	vt := NewVerticalTransform(TransformOptions{
		Region:   testRegion(),
		Progress: rec.record,
	})
	grid, _, err := vt.executeStep(transformStep{from: mslEPSG, to: ellipsoidEPSG, via: "msl2ellipsoid"})
	if err != nil {
		t.Fatalf("executeStep: %v", err)
	}
	if grid == nil {
		t.Fatal("nil grid")
	}
	h := testRegion().YSize
	last := -1
	for i, e := range rec.events {
		if e.stage != "msl2ellipsoid" {
			t.Fatalf("event %d stage %q, want msl2ellipsoid", i, e.stage)
		}
		if e.total != h {
			t.Fatalf("event %d total %d, want %d", i, e.total, h)
		}
		if e.done < last {
			t.Fatalf("event %d done %d after %d, not monotonic", i, e.done, last)
		}
		last = e.done
	}
	if last != h {
		t.Fatalf("final done %d, want %d", last, h)
	}
}

func TestDatumCancelRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	vt := NewVerticalTransform(TransformOptions{
		EpsgIn:  5773,
		EpsgOut: 3855,
		Region:  testRegion(),
		Ctx:     ctx,
	})
	res, err := vt.Run()
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	if res != nil {
		t.Errorf("expected nil result on cancel")
	}
}

func TestDatumCancelGeoidRow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	vt := NewVerticalTransform(TransformOptions{
		Region: testRegion(),
		Ctx:    ctx,
	})
	_, _, err := vt.executeStep(transformStep{from: mslEPSG, to: ellipsoidEPSG, via: "msl2ellipsoid"})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
