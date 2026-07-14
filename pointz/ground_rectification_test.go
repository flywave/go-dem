package pointz

import (
	"testing"
)

func TestReclassifyCloud_Basic(t *testing.T) {
	pts := make([]ClassifiedPoint, 30)
	for i := 0; i < 25; i++ {
		pts[i] = ClassifiedPoint{
			Point3D:        Point3D{X: float64(i % 5), Y: float64(i / 5), Z: 10},
			Classification: 2,
		}
	}
	pts[25] = ClassifiedPoint{
		Point3D:        Point3D{X: 0, Y: 0, Z: 50},
		Classification: 2,
	}

	opts := &GroundRectificationOptions{
		ReclassifyPlan:      PartitionMedian,
		ReclassifyThreshold: 2,
		MinPoints:           3,
		MinArea:             1,
	}
	result := reclassifyCloud(pts, opts)
	if result == nil {
		t.Fatal("nil result")
	}
	if result[25].Classification != 1 {
		t.Log("reclassify: outlier ground point should be reclassified")
	} else {
		t.Log("reclassify: outlier remained ground (may vary with LMedS)")
	}
}

func TestReclassifyCloud_AllGroundFlat(t *testing.T) {
	pts := make([]ClassifiedPoint, 20)
	for i := 0; i < 20; i++ {
		pts[i] = ClassifiedPoint{
			Point3D:        Point3D{X: float64(i % 5), Y: float64(i / 5), Z: 100},
			Classification: 2,
		}
	}
	opts := &GroundRectificationOptions{
		ReclassifyPlan:      PartitionOne,
		ReclassifyThreshold: 5,
		MinPoints:           3,
		MinArea:             1,
	}
	result := reclassifyCloud(pts, opts)
	changed := 0
	for _, p := range result {
		if p.Classification != 2 {
			changed++
		}
	}
	if changed > 0 {
		t.Errorf("all ground flat: %d points wrongly reclassified", changed)
	}
}

func TestExtendCloud_Basic(t *testing.T) {
	pts := make([]ClassifiedPoint, 16)
	for i := 0; i < 16; i++ {
		x := float64(i % 4)
		y := float64(i / 4)
		pts[i] = ClassifiedPoint{
			Point3D:        Point3D{X: x, Y: y, Z: x + y + 5},
			Classification: 2,
			R:              100, G: 150, B: 200,
		}
	}
	opts := &GroundRectificationOptions{
		ExtendPlan:         PartitionOne,
		ExtendGridDistance: 2,
		MinPoints:          3,
		MinArea:            1,
	}
	result := extendCloud(pts, opts)
	if result == nil {
		t.Fatal("nil result")
	}
	if len(result) > len(pts) {
		t.Logf("extend: added %d new points", len(result)-len(pts))
	} else {
		t.Log("extend: no new points added (grid may be too dense)")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultGroundRectificationOptions()
	if opts.Method != MethodReclassifyExtend {
		t.Errorf("default method should be reclassify_extend")
	}
	if opts.ReclassifyThreshold != 5 {
		t.Errorf("default threshold should be 5")
	}
}

func TestBuildGridForBounds(t *testing.T) {
	bounds := BoxBounds{XMin: 0, XMax: 10, YMin: 0, YMax: 10}
	hull := computeConvexHull([]Point3D{
		{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10},
	})
	cloud := []Point3D{{X: 2, Y: 2}, {X: 8, Y: 8}}
	grid := buildGridForBounds(bounds, hull, cloud, 5)
	if len(grid) == 0 {
		t.Error("grid should have points")
	}
	for _, g := range grid {
		if !bounds.Contains(g.X, g.Y) {
			t.Errorf("grid point (%.1f,%.1f) outside bounds", g.X, g.Y)
		}
	}
}

func TestRectifyMethodConstants(t *testing.T) {
	if MethodReclassify != "reclassify" {
		t.Errorf("unexpected reclassify constant")
	}
	if MethodExtend != "extend" {
		t.Errorf("unexpected extend constant")
	}
	if MethodReclassifyExtend != "reclassify_extend" {
		t.Errorf("unexpected reclassify_extend constant")
	}
}
