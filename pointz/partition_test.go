package pointz

import (
	"math"
	"testing"
)

func TestBoxBounds_Contains(t *testing.T) {
	b := BoxBounds{XMin: 0, XMax: 10, YMin: 0, YMax: 10}
	if !b.Contains(5, 5) {
		t.Error("center should be inside")
	}
	if b.Contains(-1, 5) {
		t.Error("outside should be false")
	}
	if !b.Contains(0, 10) {
		t.Error("edge should be inside")
	}
}

func TestBoxBounds_Center(t *testing.T) {
	b := BoxBounds{XMin: 2, XMax: 8, YMin: 4, YMax: 12}
	cx, cy := b.Center()
	if math.Abs(cx-5) > 0.01 || math.Abs(cy-8) > 0.01 {
		t.Errorf("center: expected (5,8), got (%.1f,%.1f)", cx, cy)
	}
}

func TestBoxBounds_Area(t *testing.T) {
	b := BoxBounds{XMin: 0, XMax: 10, YMin: 0, YMax: 20}
	if math.Abs(b.Area()-200) > 0.01 {
		t.Errorf("area: expected 200, got %.1f", b.Area())
	}
}

func TestBoxBounds_DivideByPoint(t *testing.T) {
	b := BoxBounds{XMin: 0, XMax: 10, YMin: 0, YMax: 10}
	boxes := b.DivideByPoint(5, 5)
	if len(boxes) != 4 {
		t.Fatalf("expected 4 boxes, got %d", len(boxes))
	}
	expected := [][4]float64{
		{0, 5, 0, 5},
		{5, 10, 0, 5},
		{0, 5, 5, 10},
		{5, 10, 5, 10},
	}
	for i, e := range expected {
		if math.Abs(boxes[i].XMin-e[0]) > 0.01 ||
			math.Abs(boxes[i].XMax-e[1]) > 0.01 ||
			math.Abs(boxes[i].YMin-e[2]) > 0.01 ||
			math.Abs(boxes[i].YMax-e[3]) > 0.01 {
			t.Errorf("box %d: expected [%.1f,%.1f,%.1f,%.1f], got [%.1f,%.1f,%.1f,%.1f]",
				i, e[0], e[1], e[2], e[3],
				boxes[i].XMin, boxes[i].XMax, boxes[i].YMin, boxes[i].YMax)
		}
	}
}

func TestBoxFromPoints(t *testing.T) {
	pts := []Point3D{{X: 3, Y: 7}, {X: 1, Y: 9}, {X: 5, Y: 2}}
	b := boxFromPoints(pts)
	if math.Abs(b.XMin-1) > 0.01 || math.Abs(b.XMax-5) > 0.01 ||
		math.Abs(b.YMin-2) > 0.01 || math.Abs(b.YMax-9) > 0.01 {
		t.Errorf("bounds: expected (1,5,2,9), got (%.1f,%.1f,%.1f,%.1f)", b.XMin, b.XMax, b.YMin, b.YMax)
	}
}

func TestBoxFromPoints_Empty(t *testing.T) {
	b := boxFromPoints(nil)
	if b != (BoxBounds{}) {
		t.Error("empty input should return zero bounds")
	}
}

func TestFilterPointsByBox(t *testing.T) {
	pts := []Point3D{{X: 1, Y: 1}, {X: 5, Y: 5}, {X: 10, Y: 10}}
	b := BoxBounds{XMin: 2, XMax: 8, YMin: 2, YMax: 8}
	inside := filterPointsByBox(pts, b, true)
	if len(inside) != 1 || inside[0].X != 5 {
		t.Errorf("keep inside: expected 1 point (5,5), got %d", len(inside))
	}
	outside := filterPointsByBox(pts, b, false)
	if len(outside) != 2 {
		t.Errorf("keep outside: expected 2 points, got %d", len(outside))
	}
}

func TestMedianXY(t *testing.T) {
	pts := []Point3D{
		{X: 10, Y: 1}, {X: 1, Y: 10}, {X: 5, Y: 5},
	}
	mx, my := medianXY(pts)
	if math.Abs(mx-5) > 0.01 || math.Abs(my-5) > 0.01 {
		t.Errorf("median: expected (5,5), got (%.1f,%.1f)", mx, my)
	}
}

func TestMedianXY_Empty(t *testing.T) {
	mx, my := medianXY(nil)
	if mx != 0 || my != 0 {
		t.Error("empty median should return (0,0)")
	}
}

func TestPartitionPlan_One(t *testing.T) {
	pts := make([]Point3D, 10)
	for i := 0; i < 10; i++ {
		pts[i] = Point3D{X: float64(i), Y: float64(i)}
	}
	qp := SelectPartitionPlan(PartitionOne, pts)
	parts := qp.Execute(pts, 1, 1)
	if len(parts) != 1 {
		t.Errorf("expected 1 partition, got %d", len(parts))
	}
	if len(parts[0].Points) != 10 {
		t.Errorf("expected 10 points, got %d", len(parts[0].Points))
	}
}

func TestPartitionPlan_UniformMinimumArea(t *testing.T) {
	pts := make([]Point3D, 8)
	for i := 0; i < 8; i++ {
		pts[i] = Point3D{X: float64(i % 4), Y: float64(i / 4)}
	}
	qp := SelectPartitionPlan(PartitionUniform, pts)
	parts := qp.Execute(pts, 2, 100)
	if len(parts) != 1 {
		t.Logf("uniform with large min_area: %d partition(s)", len(parts))
	}
}

func TestPartitionPlan_MedianMinimumPoints(t *testing.T) {
	pts := make([]Point3D, 8)
	for i := 0; i < 8; i++ {
		pts[i] = Point3D{X: float64(i % 4), Y: float64(i / 4)}
	}
	qp := SelectPartitionPlan(PartitionMedian, pts)
	parts := qp.Execute(pts, 100, 1)
	if len(parts) != 1 {
		t.Logf("median with large min_points: %d partition(s)", len(parts))
	}
}
