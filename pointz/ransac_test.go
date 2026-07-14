package pointz

import (
	"math"
	"testing"
)

func TestFitPlaneRANSAC_Flat(t *testing.T) {
	pts := make([]Point3D, 20)
	for i := 0; i < 20; i++ {
		pts[i] = Point3D{
			X: float64(i % 5),
			Y: float64(i / 5),
			Z: 2*float64(i%5) + 3*float64(i/5) + 10,
		}
	}
	plane := FitPlaneRANSAC(pts, 200, 0.1)
	pred := plane.PredictZ(2, 2)
	expected := 2*2 + 3*2 + 10
	if math.Abs(pred-float64(expected)) > 0.5 {
		t.Errorf("expected ~%d, got %.2f", expected, pred)
	}
}

func TestFitPlaneRANSAC_WithOutliers(t *testing.T) {
	pts := make([]Point3D, 30)
	for i := 0; i < 25; i++ {
		pts[i] = Point3D{
			X: float64(i % 5),
			Y: float64(i / 5),
			Z: float64(i%5) + float64(i/5) + 5,
		}
	}
	for i := 25; i < 30; i++ {
		pts[i] = Point3D{X: float64(i), Y: float64(i), Z: 9999}
	}
	plane := FitPlaneRANSAC(pts, 200, 1.0)
	dist := plane.DistanceToPoint(Point3D{X: 2, Y: 2, Z: 9})
	if dist > 1.0 {
		t.Errorf("RANSAC should ignore outliers, distance=%.2f", dist)
	}
}

func TestFitPlaneLeastSquares(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 1},
		{X: 0, Y: 1, Z: 1},
	}
	plane := FitPlaneLeastSquares(pts)
	pred := plane.PredictZ(1, 1)
	if math.Abs(pred-2) > 0.01 {
		t.Errorf("expected 2, got %.4f", pred)
	}
}

func TestPlaneAngle(t *testing.T) {
	plane := Plane{A: 0, B: 0, C: 10}
	angle := plane.AngleDeg()
	if angle > 1 {
		t.Errorf("flat plane angle: expected ~0, got %.2f", angle)
	}
}
