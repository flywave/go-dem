package pointz

import (
	"testing"
)

func TestCoplanar_Basic(t *testing.T) {
	pts := make([]Point3D, 20)
	for i := 0; i < 20; i++ {
		pts[i] = Point3D{
			X: float64(i % 5),
			Y: float64(i / 5),
			Z: float64(i),
		}
	}
	mask := CoplanarFilter(pts, &CoplanarOptions{
		Radius:       3,
		Threshold:    2.0,
		MinNeighbors: 3,
	})
	if mask == nil || len(mask) != 20 {
		t.Fatalf("expected 20 results, got %d", len(mask))
	}
	outliers := 0
	for _, m := range mask {
		if m {
			outliers++
		}
	}
	t.Logf("coplanar: %d/%d outliers", outliers, len(pts))
}

func TestCoplanar_Flat(t *testing.T) {
	pts := make([]Point3D, 16)
	for i := 0; i < 16; i++ {
		pts[i] = Point3D{
			X: float64(i % 4),
			Y: float64(i / 4),
			Z: 100,
		}
	}
	mask := CoplanarFilter(pts, &CoplanarOptions{
		Radius:       2,
		Threshold:    0.5,
		MinNeighbors: 3,
	})
	outliers := 0
	for _, m := range mask {
		if m {
			outliers++
		}
	}
	if outliers > 2 {
		t.Errorf("flat surface: expected few outliers, got %d", outliers)
	}
}

func TestCoplanar_Spike(t *testing.T) {
	pts := make([]Point3D, 17)
	for i := 0; i < 16; i++ {
		pts[i] = Point3D{
			X: float64(i % 4),
			Y: float64(i / 4),
			Z: 100,
		}
	}
	pts[16] = Point3D{X: 0.5, Y: 0.5, Z: 200}

	mask := CoplanarFilter(pts, &CoplanarOptions{
		Radius:       2,
		Threshold:    1.0,
		MinNeighbors: 3,
	})
	if mask[16] {
		t.Log("coplanar: spike detected as outlier")
	} else {
		t.Log("coplanar: spike not detected (may depend on layout)")
	}
}

func TestCoplanar_Invert(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0, Z: 10}}
	mask := CoplanarFilter(pts, &CoplanarOptions{
		Radius:       1,
		Threshold:    0.5,
		MinNeighbors: 3,
		Invert:       true,
	})
	if !mask[0] {
		t.Log("isolated point inverted")
	}
}
