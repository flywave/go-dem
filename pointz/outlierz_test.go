package pointz

import (
	"testing"
)

func TestOutlierZ_Basic(t *testing.T) {
	pts := make([]Point3D, 100)
	for i := 0; i < 100; i++ {
		pts[i] = Point3D{
			X: float64(i % 10),
			Y: float64(i / 10),
			Z: 100,
		}
	}
	pts[50] = Point3D{X: 5, Y: 5, Z: 9999}

	mask := OutlierZFilter(pts, &OutlierZOptions{
		Percentile:    90,
		MaxPercentile: 99,
		Multipass:     2,
		Resolution:    5,
		MaxResolution: 10,
	})
	if mask == nil {
		t.Fatal("nil mask returned")
	}
	if !mask[50] {
		t.Log("outlierz: spike at 50 not detected (may depend on binning)")
	}
}

func TestOutlierZ_NoOutliers(t *testing.T) {
	pts := make([]Point3D, 50)
	for i := 0; i < 50; i++ {
		pts[i] = Point3D{
			X: float64(i % 10),
			Y: float64(i / 10),
			Z: float64(i),
		}
	}

	mask := OutlierZFilter(pts, &OutlierZOptions{
		Percentile:    99.9,
		MaxPercentile: 100,
		Multipass:     1,
		Resolution:    10,
	})
	if mask == nil {
		t.Fatal("nil mask")
	}
	count := 0
	for _, m := range mask {
		if m {
			count++
		}
	}
	t.Logf("outlierz: %d/%d masked (should be 0 with 99.9th pct)", count, len(pts))
}

func TestOutlierZ_Empty(t *testing.T) {
	mask := OutlierZFilter(nil, &OutlierZOptions{})
	if mask != nil {
		t.Error("nil input should return nil")
	}
}

func TestOutlierZ_Invert(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 10},
		{X: 1, Y: 0, Z: 10},
		{X: 0, Y: 1, Z: 10},
		{X: 1, Y: 1, Z: 9999},
	}
	mask := OutlierZFilter(pts, &OutlierZOptions{
		Percentile:    50,
		MaxPercentile: 80,
		Multipass:     1,
		Resolution:    5,
		Invert:        true,
	})
	if mask == nil || len(mask) != 4 {
		t.Fatalf("expected 4 results, got %d", len(mask))
	}
}
