package pointz

import (
	"math"
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
		t.Errorf("reclassify: outlier ground point (Z=50) should be reclassified to 1, got %d", result[25].Classification)
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
	var pts []ClassifiedPoint
	for x := 0; x <= 4; x++ {
		pts = append(pts,
			ClassifiedPoint{Point3D: Point3D{X: float64(x), Y: 0, Z: 0.1*float64(x) + 5}, Classification: 2, R: 100, G: 150, B: 200},
			ClassifiedPoint{Point3D: Point3D{X: float64(x), Y: 4, Z: 0.1*float64(x) + 0.8 + 5}, Classification: 2, R: 100, G: 150, B: 200})
	}
	for y := 1; y <= 3; y++ {
		pts = append(pts,
			ClassifiedPoint{Point3D: Point3D{X: 0, Y: float64(y), Z: 0.2*float64(y) + 5}, Classification: 2, R: 100, G: 150, B: 200},
			ClassifiedPoint{Point3D: Point3D{X: 4, Y: float64(y), Z: 0.4 + 0.2*float64(y) + 5}, Classification: 2, R: 100, G: 150, B: 200})
	}

	opts := &GroundRectificationOptions{
		ExtendPlan:         PartitionOne,
		ExtendGridDistance: 1,
		MinPoints:          3,
		MinArea:            1,
	}
	result := extendCloud(pts, opts)
	if result == nil {
		t.Fatal("nil result")
	}
	added := result[len(pts):]
	if len(added) == 0 {
		t.Fatal("extend: interior grid points should be added for the boundary-ring ground")
	}
	for _, p := range added {
		if p.Classification != 2 {
			t.Errorf("extended point should be ground, got class %d", p.Classification)
		}
		expected := 0.1*p.X + 0.2*p.Y + 5
		if math.Abs(p.Z-expected) > 1e-6 {
			t.Errorf("extended point (%.1f,%.1f) should lie on plane z=0.1x+0.2y+5, got %.2f (expected %.2f)", p.X, p.Y, p.Z, expected)
		}
	}
}

func TestExtendCloud_ZeroPlane(t *testing.T) {
	var pts []ClassifiedPoint
	for x := 0; x <= 4; x++ {
		pts = append(pts,
			ClassifiedPoint{Point3D: Point3D{X: float64(x), Y: 0, Z: 0}, Classification: 2},
			ClassifiedPoint{Point3D: Point3D{X: float64(x), Y: 4, Z: 0}, Classification: 2})
	}
	for y := 1; y <= 3; y++ {
		pts = append(pts,
			ClassifiedPoint{Point3D: Point3D{X: 0, Y: float64(y), Z: 0}, Classification: 2},
			ClassifiedPoint{Point3D: Point3D{X: 4, Y: float64(y), Z: 0}, Classification: 2})
	}

	opts := &GroundRectificationOptions{
		ExtendPlan:         PartitionOne,
		ExtendGridDistance: 1,
		MinPoints:          3,
		MinArea:            1,
	}
	result := extendCloud(pts, opts)
	added := result[len(pts):]
	if len(added) == 0 {
		t.Fatal("extend: a flat z=0 ground must still produce grid points (Z==0 is not a sentinel)")
	}
	for _, p := range added {
		if p.Classification != 2 {
			t.Errorf("extended point should be ground, got class %d", p.Classification)
		}
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
