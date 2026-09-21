package pointz

import (
	"testing"
)

func TestDensity_Random(t *testing.T) {
	pts := makeTestGrid(10, 10)
	mask := DensityFilter(pts, &DensityOptions{
		Resolution: 5,
		Mode:       DensityRandom,
	})
	kept := 0
	for _, m := range mask {
		if !m {
			kept++
		}
	}
	if kept < 1 || kept > 10 {
		t.Errorf("expected ~4 kept with res=5 on 10x10 grid, got %d", kept)
	}
}

func TestDensity_Median(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 10},
		{X: 0.5, Y: 0.5, Z: 20},
		{X: 0.2, Y: 0.2, Z: 15},
	}
	mask := DensityFilter(pts, &DensityOptions{
		Resolution: 2,
		Mode:       DensityMedian,
	})
	for i, m := range mask {
		if !m && i != 2 {
			t.Errorf("median mode should keep the Z=15 point, kept index %d (Z=%.0f)", i, pts[i].Z)
		}
		if m && i == 2 {
			t.Error("median mode should keep the Z=15 point")
		}
	}
}

func TestDensity_MedianAsymmetric(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 1},
		{X: 0.1, Y: 0.1, Z: 2},
		{X: 0.2, Y: 0.2, Z: 3},
		{X: 0.3, Y: 0.3, Z: 4},
		{X: 0.4, Y: 0.4, Z: 10},
	}
	mask := DensityFilter(pts, &DensityOptions{
		Resolution: 2,
		Mode:       DensityMedian,
	})
	for i, m := range mask {
		if !m && pts[i].Z != 3 {
			t.Errorf("median of {1,2,3,4,10} is 3, kept index %d (Z=%.0f)", i, pts[i].Z)
		}
	}
}

func TestDensity_Mean(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 10},
		{X: 0.5, Y: 0.5, Z: 20},
		{X: 0.2, Y: 0.2, Z: 30},
	}
	mask := DensityFilter(pts, &DensityOptions{
		Resolution: 2,
		Mode:       DensityMean,
	})
	kept := 0
	for _, m := range mask {
		if !m {
			kept++
		}
	}
	if kept != 1 {
		t.Errorf("density mean: kept %d points (expected 1)", kept)
	}
}

func TestDensity_Center(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 10},
		{X: 3, Y: 3, Z: 20},
	}
	mask := DensityFilter(pts, &DensityOptions{
		Resolution: 5,
		Mode:       DensityCenter,
	})
	kept := 0
	for _, m := range mask {
		if !m {
			kept++
		}
	}
	if kept != 1 {
		t.Errorf("density center: both points share one cell, kept %d points (expected 1)", kept)
	}
	for i, m := range mask {
		if !m && i != 1 {
			t.Errorf("cell center is (2.5,2.5): point (3,3) is closer and should be kept, kept index %d", i)
		}
	}
}

func TestDensity_NilOptions(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0, Z: 1}}
	mask := DensityFilter(pts, nil)
	if mask == nil || len(mask) != 1 || mask[0] {
		t.Error("nil opts should keep the single point")
	}
}

func TestDensity_Empty(t *testing.T) {
	mask := DensityFilter(nil, &DensityOptions{Resolution: 10})
	if mask != nil {
		t.Error("nil input should return nil")
	}
}

func makeTestGrid(w, h int) []Point3D {
	pts := make([]Point3D, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			pts[y*w+x] = Point3D{
				X: float64(x),
				Y: float64(y),
				Z: float64(x + y),
			}
		}
	}
	return pts
}
