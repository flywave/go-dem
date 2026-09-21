package pointz

import (
	"testing"
)

func TestBlockThin_Min(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 10},
		{X: 0.5, Y: 0.5, Z: 20},
		{X: 0.2, Y: 0.8, Z: 5},
		{X: 5, Y: 5, Z: 100},
	}
	mask := BlockThinFilter(pts, &BlockThinOptions{
		Resolution: 2,
		Mode:       BlockThinMin,
	})
	if mask == nil || len(mask) != 4 {
		t.Fatalf("expected 4 results, got %d", len(mask))
	}
	kept := 0
	for _, m := range mask {
		if !m {
			kept++
		}
	}
	t.Logf("blockthin min: kept %d/%d points (2x2 grid, 2 cells)", kept, len(pts))
	if kept < 1 || kept > 2 {
		t.Errorf("expected 1-2 kept points, got %d", kept)
	}
}

func TestBlockThin_Max(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 10},
		{X: 0.5, Y: 0.5, Z: 20},
		{X: 0.2, Y: 0.8, Z: 5},
	}
	mask := BlockThinFilter(pts, &BlockThinOptions{
		Resolution: 2,
		Mode:       BlockThinMax,
	})
	kept := 0
	keptZ := 0.0
	for i, m := range mask {
		if !m {
			kept++
			keptZ = pts[i].Z
		}
	}
	if kept != 1 {
		t.Errorf("blockthin max: kept %d points (expected 1)", kept)
	}
	if keptZ != 20 {
		t.Errorf("blockthin max: expected Z=20 kept, got %.0f", keptZ)
	}
}

func TestBlockThin_Mean(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 10},
		{X: 0.5, Y: 0, Z: 20},
		{X: 0, Y: 0.5, Z: 30},
	}
	mask := BlockThinFilter(pts, &BlockThinOptions{
		Resolution: 2,
		Mode:       BlockThinMean,
	})
	kept := 0
	keptZ := 0.0
	for i, m := range mask {
		if !m {
			kept++
			keptZ = pts[i].Z
		}
	}
	if kept != 1 {
		t.Errorf("blockthin mean: kept %d points (expected 1)", kept)
	}
	if keptZ != 20 {
		t.Errorf("blockthin mean: mean of {10,20,30} is 20, kept Z=%.0f", keptZ)
	}
}

func TestBlockThin_Median(t *testing.T) {
	pts := []Point3D{
		{X: 0, Y: 0, Z: 10},
		{X: 0.5, Y: 0.5, Z: 20},
		{X: 0.2, Y: 0.2, Z: 30},
	}
	mask := BlockThinFilter(pts, &BlockThinOptions{
		Resolution: 2,
		Mode:       BlockThinMedian,
	})
	kept := 0
	keptZ := 0.0
	for i, m := range mask {
		if !m {
			kept++
			keptZ = pts[i].Z
		}
	}
	if kept != 1 {
		t.Errorf("blockthin median: kept %d points (expected 1)", kept)
	}
	if keptZ != 20 {
		t.Errorf("blockthin median: median of {10,20,30} is 20, kept Z=%.0f", keptZ)
	}
}

func TestBlockThin_NilOptions(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0, Z: 10}}
	mask := BlockThinFilter(pts, nil)
	if mask == nil || len(mask) != 1 || mask[0] {
		t.Error("nil opts should keep the single point")
	}
}

func TestBlockThin_Empty(t *testing.T) {
	mask := BlockThinFilter(nil, &BlockThinOptions{Resolution: 10})
	if mask != nil {
		t.Error("nil input should return nil")
	}
}

func TestBlockThin_DefaultMode(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0, Z: 10}}
	mask := BlockThinFilter(pts, &BlockThinOptions{Resolution: 1})
	if mask == nil || len(mask) != 1 {
		t.Fatalf("expected 1 result")
	}
	if mask[0] {
		t.Error("single point should be kept (not masked)")
	}
}
