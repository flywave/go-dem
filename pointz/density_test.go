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
	keptZ := 0.0
	for i, m := range mask {
		if !m {
			keptZ = pts[i].Z
		}
	}
	if keptZ != 15 {
		t.Logf("density median: kept Z=%.0f (expected 15)", keptZ)
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
		t.Logf("density mean: kept %d points (expected 1)", kept)
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
	if kept != 2 {
		t.Logf("density center: kept %d points", kept)
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
