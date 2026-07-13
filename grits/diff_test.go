package grits

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func TestDiff_Identical(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeRampDEM(w, h)
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)

	df := ComputeDiff(data, data, reg, false, nd)
	for i, v := range df {
		if v == nd || math.IsNaN(v) {
			continue
		}
		if math.Abs(v) > 1e-10 {
			t.Errorf("identical dems: diff at %d = %.6f", i, v)
		}
	}
}

func TestDiff_ConstantOffset(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeRampDEM(w, h)
	offset := make([]float64, w*h)
	copy(offset, data)
	for i := range offset {
		offset[i] += 50
	}

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	df := ComputeDiff(offset, data, reg, false, nd)
	for i, v := range df {
		if v == nd || math.IsNaN(v) {
			continue
		}
		if math.Abs(v-50) > 1e-10 {
			t.Errorf("offset diff: expected 50 at %d, got %.2f", i, v)
		}
	}
}

func TestDiff_AbsDiff(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	a := makeFlatDEM(w, h, 100)
	b := makeFlatDEM(w, h, 50)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	df := ComputeDiff(a, b, reg, true, nd)
	for i, v := range df {
		if v == nd || math.IsNaN(v) {
			continue
		}
		if math.Abs(v-50) > 1e-10 {
			t.Errorf("abs diff: expected 50 at %d, got %.2f", i, v)
		}
	}
}

func TestDiff_NoDataHandling(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	a := makeFlatDEM(w, h, 100)
	b := makeFlatDEM(w, h, 50)
	a[2*w+2] = nd

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	df := ComputeDiff(a, b, reg, false, nd)
	if df[2*w+2] != nd {
		t.Error("noData should remain noData in diff")
	}
}

func TestDiff_Threshold(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	a := makeFlatDEM(w, h, 100)
	b := makeFlatDEM(w, h, 99.5)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	df := ComputeDiff(a, b, reg, false, nd)
	for i, v := range df {
		if v == nd || math.IsNaN(v) {
			continue
		}
		if math.Abs(v-0.5) > 1e-10 {
			t.Errorf("expected 0.5 at %d, got %.4f", i, v)
		}
	}
}

func TestDiff_ResampleDiff(t *testing.T) {
	nd := -9999.0
	reg1 := dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
	reg2 := dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
	data := makeRampDEM(10, 10)
	ref := makeRampDEM(10, 10)

	result, err := resampleDiff(data, reg1, ref, reg2, nd)
	if err != nil {
		t.Fatalf("resampleDiff error: %v", err)
	}
	if len(result) != 100 {
		t.Errorf("output size: expected 100, got %d", len(result))
	}
}
