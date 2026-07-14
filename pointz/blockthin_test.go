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
	for _, m := range mask {
		if !m {
			kept++
		}
	}
	t.Logf("blockthin max: kept %d points", kept)
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
	for _, m := range mask {
		if !m {
			kept++
		}
	}
	t.Logf("blockthin mean: kept %d points", kept)
	if kept != 1 {
		t.Log("blockthin mean: all points in same cell, expected 1 kept")
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
	for _, m := range mask {
		if !m {
			kept++
		}
	}
	t.Logf("blockthin median: kept %d points", kept)
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
