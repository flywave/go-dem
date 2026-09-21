package waffle

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

func TestCUDEM_MeasuredMeanAtDataPoints(t *testing.T) {
	region := dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
	pts := []Point{
		{Position: vec2.T{2.1, 6.2}, Z: 0},
		{Position: vec2.T{2.4, 6.5}, Z: 50},
		{Position: vec2.T{2.7, 6.9}, Z: 100},
		{Position: vec2.T{7.1, 7.2}, Z: 200},
		{Position: vec2.T{7.4, 7.5}, Z: 200},
		{Position: vec2.T{7.7, 7.8}, Z: 200},
	}

	w, err := New(dem.MethodCUDEM)
	if err != nil {
		t.Fatalf("new cudem: %v", err)
	}
	res, err := w.Run(pts, &Options{Region: region})
	if err != nil {
		t.Fatalf("run cudem: %v", err)
	}
	if len(res.DEM) != 100 {
		t.Fatalf("expected 100 cells, got %d", len(res.DEM))
	}

	if got := res.DEM[3*10+2]; math.Abs(got-50) > 1e-6 {
		t.Errorf("pixel (2,3): expected measured mean 50, got %.4f", got)
	}
	if got := res.DEM[2*10+7]; math.Abs(got-200) > 1e-6 {
		t.Errorf("pixel (7,2): expected measured mean 200, got %.4f", got)
	}
}

func TestCUDEM_MeasuredMeanNorthRow(t *testing.T) {
	region := dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
	pts := []Point{
		{Position: vec2.T{4.2, 9.3}, Z: 10},
		{Position: vec2.T{4.5, 9.5}, Z: 20},
		{Position: vec2.T{4.8, 9.7}, Z: 60},
	}

	w, _ := New(dem.MethodCUDEM)
	res, err := w.Run(pts, &Options{Region: region})
	if err != nil {
		t.Fatalf("run cudem: %v", err)
	}
	if got := res.DEM[0*10+4]; math.Abs(got-30) > 1e-6 {
		t.Errorf("north row pixel (4,0): expected measured mean 30, got %.4f", got)
	}
}

func TestCUDEM_GradientHoleInterpolation(t *testing.T) {
	region := dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
	var pts []Point
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			if x >= 3 && x <= 6 && y >= 3 && y <= 6 {
				continue
			}
			geoX := float64(x) + 0.5
			geoY := 10 - (float64(y) + 0.5)
			pts = append(pts, Point{Position: vec2.T{geoX, geoY}, Z: 10 * geoY})
		}
	}

	w, _ := New(dem.MethodCUDEM)
	res, err := w.Run(pts, &Options{Region: region})
	if err != nil {
		t.Fatalf("run cudem: %v", err)
	}

	for y := 3; y <= 6; y++ {
		expected := 10 * (10 - (float64(y) + 0.5))
		for x := 3; x <= 6; x++ {
			got := res.DEM[y*10+x]
			if got == dem.DefaultNoData || math.IsNaN(got) {
				t.Errorf("hole pixel (%d,%d) not filled", x, y)
				continue
			}
			if math.Abs(got-expected) > 25 {
				t.Errorf("hole pixel (%d,%d): expected ~%.1f by north-up row, got %.2f", x, y, expected, got)
			}
		}
	}
}

func TestCUDEM_EmptyPoints(t *testing.T) {
	w, _ := New(dem.MethodCUDEM)
	res, err := w.Run(nil, &Options{Region: dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)})
	if err == nil {
		t.Fatal("expected error for empty points")
	}
	if res != nil {
		t.Fatal("expected nil result with error")
	}
}
