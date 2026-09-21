package dem

import (
	"context"
	"errors"
	"testing"
)

func TestReportProgress(t *testing.T) {
	type event struct {
		stage string
		done  int
		total int
	}
	var events []event
	fn := ProgressFunc(func(stage string, done, total int) {
		events = append(events, event{stage, done, total})
	})

	ReportProgress(fn, "stage_a", 3, 10)
	ReportProgress(fn, "stage_a", 7, 10)
	ReportProgress(fn, "stage_b", 1, 1)

	if len(events) != 3 {
		t.Fatalf("expected 3 callbacks, got %d", len(events))
	}
	if events[0].stage != "stage_a" || events[0].done != 3 || events[0].total != 10 {
		t.Errorf("event 0 = %+v", events[0])
	}
	if events[1].done != 7 {
		t.Errorf("event 1 = %+v", events[1])
	}
	if events[2].stage != "stage_b" || events[2].done != 1 || events[2].total != 1 {
		t.Errorf("event 2 = %+v", events[2])
	}
}

func TestReportProgressNilFunc(t *testing.T) {
	ReportProgress(nil, "noop", 1, 1)
}

func TestCheckCtxNil(t *testing.T) {
	if err := CheckCtx(nil); err != nil {
		t.Errorf("nil ctx must return nil, got %v", err)
	}
}

func TestCheckCtxLive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := CheckCtx(ctx); err != nil {
		t.Errorf("live ctx must return nil, got %v", err)
	}
}

func TestCheckCtxCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := CheckCtx(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("canceled ctx must return context.Canceled, got %v", err)
	}
}
