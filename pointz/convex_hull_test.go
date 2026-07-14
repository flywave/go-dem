package pointz

import (
	"math"
	"testing"
)

func TestConvexHull_Square(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10},
		{X: 5, Y: 5},
	}
	hull := computeConvexHull(pts)
	if len(hull.Points) < 3 {
		t.Fatalf("convex hull should have >= 3 points, got %d", len(hull.Points))
	}
}

func TestConvexHull_Triangle(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0}, {X: 5, Y: 10}, {X: 10, Y: 0}}
	hull := computeConvexHull(pts)
	if len(hull.Points) != 3 {
		t.Errorf("triangle hull: expected 3, got %d", len(hull.Points))
	}
}

func TestConvexHull_KeepPointsInside(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10}}
	hull := computeConvexHull(pts)

	inside := []Point3D{{X: 5, Y: 5}, {X: 1, Y: 1}}
	outside := []Point3D{{X: -1, Y: 5}, {X: 15, Y: 15}, {X: 5, Y: -1}}

	all := append(inside, outside...)
	filtered := hull.KeepPointsInside(all)
	if len(filtered) != 2 {
		t.Errorf("expected 2 inside, got %d", len(filtered))
	}
}

func TestConvexHull_CalculateMask(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}}
	hull := computeConvexHull(pts)

	test := []Point3D{{X: 2, Y: 2}, {X: -1, Y: 2}, {X: 6, Y: 2}}
	mask := hull.CalculateMask(test)
	if !mask[0] {
		t.Error("expected (2,2) to be inside")
	}
	if mask[1] || mask[2] {
		t.Error("expected (-1,2) and (6,2) to be outside")
	}
}

func TestConvexHull_Bounds(t *testing.T) {
	pts := []Point3D{{X: 2, Y: 3}, {X: 8, Y: 5}, {X: 6, Y: 9}, {X: 1, Y: 7}}
	hull := computeConvexHull(pts)
	xMin, xMax, yMin, yMax := hull.Bounds()
	if math.Abs(xMin-1) > 0.01 || math.Abs(xMax-8) > 0.01 ||
		math.Abs(yMin-3) > 0.01 || math.Abs(yMax-9) > 0.01 {
		t.Errorf("bounds: expected (1,8,3,9), got (%.1f,%.1f,%.1f,%.1f)", xMin, xMax, yMin, yMax)
	}
}

func TestOrientation(t *testing.T) {
	a, b := Point3D{X: 0, Y: 0}, Point3D{X: 5, Y: 0}
	o1 := orientation(a, b, Point3D{X: 3, Y: 1})
	o2 := orientation(a, b, Point3D{X: 3, Y: -1})
	o3 := orientation(a, b, Point3D{X: 3, Y: 0})
	if o1 == o2 {
		t.Errorf("opposite sides should give different orientations: %d vs %d", o1, o2)
	}
	if o3 != 0 {
		t.Errorf("collinear expected 0, got %d", o3)
	}
	t.Logf("orientation: up=%d, down=%d, collinear=%d", o1, o2, o3)
}

func TestConvexHull_SinglePoint(t *testing.T) {
	hull := computeConvexHull([]Point3D{{X: 5, Y: 5}})
	if len(hull.Points) != 1 {
		t.Errorf("single point: expected 1, got %d", len(hull.Points))
	}
	mask := hull.CalculateMask([]Point3D{{X: 5, Y: 5}, {X: 10, Y: 10}})
	if len(mask) != 2 {
		t.Error("mask length mismatch")
	}
}

func TestConvexHull_Empty(t *testing.T) {
	hull := computeConvexHull(nil)
	if len(hull.Points) != 0 {
		t.Error("empty hull should have 0 points")
	}
}
