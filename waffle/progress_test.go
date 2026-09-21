package waffle

import (
	"context"
	"errors"
	"testing"

	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

type progressEvent struct {
	stage string
	done  int
	total int
}

type progressRecorder struct {
	events []progressEvent
}

func (r *progressRecorder) record(stage string, done, total int) {
	r.events = append(r.events, progressEvent{stage, done, total})
}

func assertProgressMonotonic(t *testing.T, rec *progressRecorder, stage string) {
	t.Helper()
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
}

func assertProgressFinal(t *testing.T, rec *progressRecorder) {
	t.Helper()
	if len(rec.events) == 0 {
		t.Fatal("no progress callbacks")
	}
	last := rec.events[len(rec.events)-1]
	if last.done != last.total {
		t.Fatalf("final callback done=%d total=%d, want done==total", last.done, last.total)
	}
}

func testWafflePoints() []Point {
	pts := make([]Point, 0, 25)
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			pts = append(pts, Point{
				Position: vec2.T{float64(x)*2 + 1, float64(y)*2 + 1},
				Z:        float64(x + y),
			})
		}
	}
	return pts
}

func TestWaffleProgress(t *testing.T) {
	pts := testWafflePoints()
	methods := []dem.InterpMethod{
		dem.MethodIDW,
		dem.MethodKriging,
		dem.MethodLinear,
		dem.MethodNearest,
		dem.MethodCUDEM,
		dem.MethodInpaint,
		dem.MethodCUBE,
		dem.MethodMovingAverage,
		dem.MethodNaturalNeighbor,
	}
	for _, m := range methods {
		rec := &progressRecorder{}
		w, err := New(m)
		if err != nil {
			t.Fatalf("new %s: %v", m, err)
		}
		res, err := w.Run(pts, &Options{Region: testRegion(), Progress: rec.record})
		if err != nil {
			t.Fatalf("%s: %v", m, err)
		}
		if res == nil {
			t.Fatalf("%s: nil result", m)
		}
		if len(rec.events) == 0 {
			t.Errorf("%s: no progress callbacks", m)
			continue
		}
		for _, e := range rec.events {
			if e.stage != string(m) {
				t.Errorf("%s: stage %q, want %q", m, e.stage, string(m))
				break
			}
		}
		assertProgressMonotonic(t, rec, string(m))
		assertProgressFinal(t, rec)
	}
}

func TestWaffleProgressNilCallbacksSafe(t *testing.T) {
	pts := testWafflePoints()
	w, err := New(dem.MethodIDW)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Run(pts, &Options{Region: testRegion()}); err != nil {
		t.Fatalf("run without progress: %v", err)
	}
}

func TestWaffleCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	pts := testWafflePoints()
	for _, m := range ListMethods() {
		w, err := New(m)
		if err != nil {
			t.Fatalf("new %s: %v", m, err)
		}
		res, err := w.Run(pts, &Options{Region: testRegion(), Ctx: ctx})
		if !errors.Is(err, context.Canceled) {
			t.Errorf("%s: expected context.Canceled, got %v", m, err)
		}
		if res != nil {
			t.Errorf("%s: expected nil result on cancel", m)
		}
	}
}
