package pointz

import (
	"testing"
)

func TestRangeZ_Basic(t *testing.T) {
	pts := []Point3D{
		{Z: -10}, {Z: 0}, {Z: 10}, {Z: 50}, {Z: 100},
	}
	mask := RangeZFilter(pts, &RangeZOptions{MinZ: 0, MaxZ: 50})
	if mask == nil {
		t.Fatal("nil mask")
	}
	if !mask[0] || !mask[4] {
		t.Error("points outside range should be masked")
	}
	if mask[1] || mask[2] {
		t.Error("points inside range should not be masked")
	}
}

func TestRangeZ_Invert(t *testing.T) {
	pts := []Point3D{{Z: 5}, {Z: 15}, {Z: 25}}
	mask := RangeZFilter(pts, &RangeZOptions{MinZ: 10, MaxZ: 20, Invert: true})
	if mask[0] || mask[2] {
		t.Error("invert: outside range should not be masked")
	}
	if !mask[1] {
		t.Error("invert: inside range should be masked")
	}
}

func TestRangeZ_NilOptions(t *testing.T) {
	pts := []Point3D{{Z: 10}}
	mask := RangeZFilter(pts, nil)
	if mask == nil || len(mask) != 1 {
		t.Error("nil options should return all-false mask")
	}
	if mask[0] {
		t.Error("nil options should not mask any points")
	}
}

func TestRangeZ_Empty(t *testing.T) {
	mask := RangeZFilter(nil, &RangeZOptions{})
	if mask != nil {
		t.Error("nil input should return nil")
	}
}
