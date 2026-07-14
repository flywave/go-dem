package waffle

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

func region() *dem.Region {
	return dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
}

func TestMovingAverage_Basic(t *testing.T) {
	pts := []Point{
		{Position: vec2.T{2, 2}, Z: 10},
		{Position: vec2.T{2, 3}, Z: 20},
		{Position: vec2.T{3, 2}, Z: 30},
		{Position: vec2.T{3, 3}, Z: 40},
	}
	w, err := New(dem.MethodMovingAverage)
	if err != nil {
		t.Fatalf("new error: %v", err)
	}
	res, err := w.Run(pts, &Options{
		Region:       region(),
		SearchRadius: 5,
		MinPoints:    1,
	})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if len(res.DEM) != 100 {
		t.Errorf("expected 100 cells, got %d", len(res.DEM))
	}
	centerIdx := 5*10 + 5
	if math.Abs(res.DEM[centerIdx]-25) > 5 {
		t.Logf("near point average: expected ~25, got %.2f", res.DEM[centerIdx])
	}
}

func TestMovingAverage_SinglePoint(t *testing.T) {
	pts := []Point{
		{Position: vec2.T{5, 5}, Z: 42},
	}
	w, err := New(dem.MethodMovingAverage)
	if err != nil {
		t.Fatalf("new error: %v", err)
	}
	res, err := w.Run(pts, &Options{
		Region:       region(),
		SearchRadius: 2,
		MinPoints:    1,
	})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	centerIdx := 5*10 + 5
	if res.DEM[centerIdx] != 42 {
		t.Errorf("single point: expected 42, got %.2f", res.DEM[centerIdx])
	}
}

func TestMovingAverage_NoPoints(t *testing.T) {
	_, err := New(dem.MethodMovingAverage)
	if err != nil {
		t.Fatalf("new error: %v", err)
	}
	_, err = (&movingAverageWaffle{}).Run(nil, &Options{Region: region()})
	if err == nil {
		t.Error("expected error for no points")
	}
}

func TestMovingAverage_OutOfRange(t *testing.T) {
	pts := []Point{
		{Position: vec2.T{100, 100}, Z: 50},
	}
	w, _ := New(dem.MethodMovingAverage)
	res, err := w.Run(pts, &Options{
		Region:       region(),
		SearchRadius: 1,
		MinPoints:    2,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	allNoData := true
	for _, v := range res.DEM {
		if v != dem.DefaultNoData {
			allNoData = false
			break
		}
	}
	if !allNoData {
		t.Error("distant point: all cells should be noData when not enough neighbors")
	}
}
