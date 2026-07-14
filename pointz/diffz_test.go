package pointz

import (
	"testing"
)

func TestDiffZ_Basic(t *testing.T) {
	pts := []Point3D{{Z: -5}, {Z: 0}, {Z: 5}, {Z: 10}, {Z: 15}}
	mask := DiffZFilter(pts, &DiffZOptions{MinDiff: 0, MaxDiff: 10})
	if !mask[0] || !mask[4] {
		t.Error("points outside diff range should be masked")
	}
	if mask[1] || mask[2] || mask[3] {
		t.Error("points inside diff range should not be masked")
	}
}

func TestDiffZ_Invert(t *testing.T) {
	pts := []Point3D{{Z: 5}, {Z: 15}}
	mask := DiffZFilter(pts, &DiffZOptions{MinDiff: 10, MaxDiff: 20, Invert: true})
	if mask[0] {
		t.Error("invert: point below range should not be masked")
	}
	if !mask[1] {
		t.Error("invert: point inside range should be masked")
	}
}

func TestDiffZ_NilOptions(t *testing.T) {
	mask := DiffZFilter([]Point3D{{Z: 10}}, nil)
	if mask[0] {
		t.Error("nil options should not mask")
	}
}

func TestDiffZ_Empty(t *testing.T) {
	mask := DiffZFilter(nil, &DiffZOptions{})
	if mask != nil {
		t.Error("nil input should return nil")
	}
}
