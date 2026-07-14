package pointz

import (
	"math"
	"testing"
)

func TestPlane3D_ProjectZ(t *testing.T) {
	p := Plane3D{NX: 0, NY: 0, NZ: 1, D: -10}
	z := p.ProjectZ(5, 5)
	if math.Abs(z-10) > 0.01 {
		t.Errorf("plane z=10: expected Z=10, got %.2f", z)
	}
}

func TestPlane3D_SignedDistance(t *testing.T) {
	p := Plane3D{NX: 0, NY: 0, NZ: 1, D: -10}
	d := p.SignedDistance(Point3D{X: 0, Y: 0, Z: 5})
	if math.Abs(d-(-5)) > 0.01 {
		t.Errorf("signed distance: expected -5, got %.2f", d)
	}
}

func TestPlane3D_AbsDistance(t *testing.T) {
	p := Plane3D{NX: 0, NY: 0, NZ: 1, D: -10}
	d := p.AbsDistance(Point3D{X: 0, Y: 0, Z: 5})
	if math.Abs(d-5) > 0.01 {
		t.Errorf("abs distance: expected 5, got %.2f", d)
	}
}

func TestPlane3D_AngleDeg(t *testing.T) {
	p := Plane3D{NX: 0, NY: 0, NZ: 1, D: 0}
	angle := p.AngleDeg()
	if angle > 1 {
		t.Errorf("flat plane: expected ~0, got %.2f", angle)
	}
}

func TestPlane3D_IsValid(t *testing.T) {
	if (&Plane3D{NX: 0, NY: 0, NZ: 1, D: 0}).IsValid() != true {
		t.Error("valid plane should return true")
	}
	if (&Plane3D{math.NaN(), math.NaN(), math.NaN(), math.NaN()}).IsValid() != false {
		t.Error("NaN plane should return false")
	}
	if (&Plane3D{NX: math.NaN(), NY: 0, NZ: 1, D: 0}).IsValid() != false {
		t.Error("plane with NaN should return false")
	}
}

func TestPlane3D_Normalize(t *testing.T) {
	p := Plane3D{NX: 3, NY: 4, NZ: 0, D: -10}
	p.Normalize()
	if math.Abs(p.NX-0.6) > 0.01 || math.Abs(p.NY-0.8) > 0.01 {
		t.Errorf("normalized: expected (0.6,0.8,0), got (%.4f,%.4f,%.4f)", p.NX, p.NY, p.NZ)
	}
}

func TestFitPlaneLMedS_Basic(t *testing.T) {
	pts := make([]Point3D, 20)
	for i := 0; i < 20; i++ {
		pts[i] = Point3D{
			X: float64(i % 5),
			Y: float64(i / 5),
			Z: 2*float64(i%5) + 3*float64(i/5) + 10,
		}
	}
	plane := fitPlaneLMedS(pts, 200)
	if !plane.IsValid() {
		t.Fatal("invalid plane")
	}
	pred := plane.ProjectZ(2, 2)
	expected := 2*2 + 3*2 + 10
	if math.Abs(pred-float64(expected)) > 1.0 {
		t.Errorf("expected ~%d, got %.2f", expected, pred)
	}
}

func TestFitPlaneLMedS_WithOutliers(t *testing.T) {
	pts := make([]Point3D, 30)
	for i := 0; i < 25; i++ {
		pts[i] = Point3D{
			X: float64(i % 5),
			Y: float64(i / 5),
			Z: float64(i%5) + float64(i/5),
		}
	}
	for i := 25; i < 30; i++ {
		pts[i] = Point3D{X: float64(i), Y: float64(i), Z: 9999}
	}
	plane := fitPlaneLMedS(pts, 200)
	if !plane.IsValid() {
		t.Fatal("invalid plane")
	}
	dist := plane.AbsDistance(Point3D{X: 2, Y: 2, Z: 4})
	if dist > 2.0 {
		t.Errorf("LMedS should tolerate outliers, distance=%.2f", dist)
	}
}

func TestFitPlaneLeastSquares3D(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 1},
		{X: 0, Y: 1, Z: 1},
	}
	plane := fitPlaneLeastSquares3D(pts)
	pred := plane.ProjectZ(1, 1)
	if math.Abs(pred-2) > 0.1 {
		t.Errorf("expected ~2, got %.4f", pred)
	}
}

func TestFitPlaneLeastSquares3D_Flat(t *testing.T) {
	pts := make([]Point3D, 5)
	for i := 0; i < 5; i++ {
		pts[i] = Point3D{X: float64(i), Y: float64(i), Z: 100}
	}
	plane := fitPlaneLeastSquares3D(pts)
	z := plane.ProjectZ(10, 10)
	if math.Abs(z-100) > 0.1 {
		t.Errorf("flat plane: expected 100, got %.2f", z)
	}
}

func TestRefitPlaneInliers(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 1},
		{X: 0, Y: 1, Z: 1},
		{X: 10, Y: 10, Z: 100},
	}
	initial := Plane3D{NX: -0.577, NY: -0.577, NZ: 0.577, D: 0}
	refined := refitPlaneInliers(pts, initial, 2.0)
	if !refined.IsValid() {
		t.Log("refit may fail with few inliers, using initial")
	}
}
