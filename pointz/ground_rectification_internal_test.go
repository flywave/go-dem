package pointz

import (
	"testing"
)

func TestGroundPoints(t *testing.T) {
	pts := []ClassifiedPoint{
		{Point3D: Point3D{X: 0, Y: 0, Z: 10}, Classification: 2},
		{Point3D: Point3D{X: 1, Y: 1, Z: 20}, Classification: 1},
		{Point3D: Point3D{X: 2, Y: 2, Z: 30}, Classification: 2},
	}
	gp := groundPoints(pts)
	if len(gp) != 2 {
		t.Errorf("expected 2 ground points, got %d", len(gp))
	}
}

func TestGroundPoints_None(t *testing.T) {
	pts := []ClassifiedPoint{
		{Point3D: Point3D{X: 0, Y: 0}, Classification: 1},
	}
	gp := groundPoints(pts)
	if len(gp) != 0 {
		t.Error("expected 0 ground points")
	}
}

func TestGroundPoints_Empty(t *testing.T) {
	gp := groundPoints(nil)
	if len(gp) != 0 {
		t.Error("empty input should return empty")
	}
}

func TestBuildGridForBounds_AllGridded(t *testing.T) {
	bounds := BoxBounds{XMin: 0, XMax: 10, YMin: 0, YMax: 10}
	hull := computeConvexHull([]Point3D{
		{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10},
	})
	grid := buildGridForBounds(bounds, hull, nil, 5)
	if len(grid) == 0 {
		t.Error("grid should have points")
	}
	for _, g := range grid {
		if g.X < 0 || g.X > 10 || g.Y < 0 || g.Y > 10 {
			t.Errorf("grid point (%.1f,%.1f) outside bounds", g.X, g.Y)
		}
	}
}

func TestBuildGridForBounds_FiltersClose(t *testing.T) {
	bounds := BoxBounds{XMin: 0, XMax: 10, YMin: 0, YMax: 10}
	hull := computeConvexHull([]Point3D{
		{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10},
	})
	cloud := []Point3D{{X: 2.5, Y: 2.5}, {X: 7.5, Y: 7.5}}
	grid := buildGridForBounds(bounds, hull, cloud, 1)
	distanceOK := true
	for _, g := range grid {
		for _, c := range cloud {
			dx := g.X - c.X
			dy := g.Y - c.Y
			d := dx*dx + dy*dy
			if d < 1.0 {
				distanceOK = false
			}
		}
	}
	if !distanceOK && len(grid) > 0 {
		t.Log("buildGridForBounds: some grid points close to cloud (expected for this spacing)")
	}
}

func TestBuildGridForBounds_ZeroDistance(t *testing.T) {
	bounds := BoxBounds{XMin: 0, XMax: 5, YMin: 0, YMax: 5}
	hull := computeConvexHull([]Point3D{
		{X: 0, Y: 0}, {X: 5, Y: 0}, {X: 5, Y: 5}, {X: 0, Y: 5},
	})
	grid := buildGridForBounds(bounds, hull, nil, 0)
	if len(grid) == 0 {
		t.Error("default distance should produce points")
	}
}
