package pointz

import (
	"testing"
)

func TestClassifyGroundSMRF_Basic(t *testing.T) {
	pts := make([]Point3D, 30)
	for i := 0; i < 25; i++ {
		x := float64(i % 5)
		y := float64(i / 5)
		pts[i] = Point3D{X: x, Y: y, Z: 10 + x*0.1 + y*0.1}
	}
	for i := 25; i < 30; i++ {
		x := float64(i % 5)
		y := float64(i / 5)
		pts[i] = Point3D{X: x, Y: y, Z: 10 + x*0.1 + y*0.1 + 50}
	}

	classification := ClassifyGroundSMRF(pts, &SMRFGroundClassificationOptions{
		CellSize:  1.0,
		Slope:     0.3,
		Window:    3.0,
		Scalar:    1.0,
		Threshold: 0.5,
	})
	if classification == nil {
		t.Fatal("nil classification")
	}
	if len(classification) != len(pts) {
		t.Fatalf("classification len %d != points %d", len(classification), len(pts))
	}
	groundCount := 0
	highGroundCount := 0
	for i, c := range classification {
		if c == 2 {
			groundCount++
			if i >= 25 {
				highGroundCount++
			}
		}
	}
	t.Logf("smrf: %d/%d ground (low), %d/%d misclassified as ground (high)",
		groundCount-(5-highGroundCount), 25, highGroundCount, 5)
	if groundCount < 20 {
		t.Errorf("expected most low points to be ground, got %d low points ground", groundCount)
	}
}

func TestClassifyGroundSMRF_AllGround(t *testing.T) {
	pts := make([]Point3D, 20)
	for i := 0; i < 20; i++ {
		pts[i] = Point3D{
			X: float64(i % 5),
			Y: float64(i / 5),
			Z: 100,
		}
	}
	classification := ClassifyGroundSMRF(pts, &SMRFGroundClassificationOptions{
		CellSize:  1.0,
		Slope:     0.15,
		Window:    5.0,
		Scalar:    0.5,
		Threshold: 0.3,
	})
	groundCount := 0
	for _, c := range classification {
		if c == 2 {
			groundCount++
		}
	}
	if groundCount < 18 {
		t.Errorf("all flat: expected most points ground, got %d/20", groundCount)
	}
}

func TestClassifyGroundSMRF_DefaultOptions(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0, Z: 10}}
	classification := ClassifyGroundSMRF(pts, nil)
	if classification == nil || len(classification) != 1 {
		t.Error("should return classification for single point")
	}
}

func TestClassifyGroundSMRF_Empty(t *testing.T) {
	classification := ClassifyGroundSMRF(nil, nil)
	if classification != nil {
		t.Error("nil input should return nil")
	}
}
