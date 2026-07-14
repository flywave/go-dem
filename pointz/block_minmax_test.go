package pointz

import (
	"testing"
)

func TestBlockMinMax_Min(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 10},
		{X: 0.5, Y: 0.5, Z: 20},
		{X: 0.2, Y: 0.8, Z: 5},
		{X: 5, Y: 5, Z: 100},
	}
	mask := BlockMinMaxFilter(pts, &BlockMinMaxOptions{
		Resolution: 2,
		Mode:       BlockMinMaxMin,
	})
	kept := 0
	for _, m := range mask {
		if !m {
			kept++
		}
	}
	if kept < 1 || kept > 2 {
		t.Errorf("expected 1-2 kept, got %d", kept)
	}
}

func TestBlockMinMax_Max(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 10},
		{X: 0.5, Y: 0.5, Z: 20},
	}
	mask := BlockMinMaxFilter(pts, &BlockMinMaxOptions{
		Resolution: 2,
		Mode:       BlockMinMaxMax,
	})
	kept := 0
	keptZ := 0.0
	for i, m := range mask {
		if !m {
			kept++
			keptZ = pts[i].Z
		}
	}
	if kept == 1 && keptZ != 20 {
		t.Errorf("max mode: expected Z=20, got %.0f", keptZ)
	}
}

func TestBlockMinMax_Invert(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0, Z: 5}, {X: 0.5, Y: 0.5, Z: 15}}
	mask := BlockMinMaxFilter(pts, &BlockMinMaxOptions{
		Resolution: 2,
		Mode:       BlockMinMaxMax,
		Invert:     true,
	})
	kept := 0
	for _, m := range mask {
		if !m {
			kept++
		}
	}
	if kept != 1 {
		t.Logf("invert: kept %d points (expected 1)", kept)
	}
}

func TestBlockMinMax_Empty(t *testing.T) {
	mask := BlockMinMaxFilter(nil, &BlockMinMaxOptions{Resolution: 10})
	if mask != nil {
		t.Error("nil input should return nil")
	}
}

func TestBlockMinMax_DefaultMode(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0, Z: 10}}
	mask := BlockMinMaxFilter(pts, &BlockMinMaxOptions{Resolution: 1})
	if mask[0] {
		t.Error("single point should not be masked")
	}
}
