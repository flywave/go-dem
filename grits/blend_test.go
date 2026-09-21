package grits

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func TestBlend_NoMask(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeRampDEM(w, h)
	bf := &blendFilter{}
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	res, err := bf.Run(data, reg, &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i, v := range res {
		if v != data[i] {
			t.Errorf("no-mask: pixel %d changed %.0f->%.2f", i, data[i], v)
		}
	}
}

func TestLinearBlend_NoOverlap(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	dem1 := makeFlatDEM(w, h, 100)
	dem2 := makeFlatDEM(w, h, 200)
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)

	result := linearBlend(dem1, dem2, reg, nd, 2)
	for i, v := range result {
		if math.Abs(v-dem1[i]) > 1e-10 {
			t.Errorf("no mask holes: pixel %d should be unchanged, got %.2f", i, v)
		}
	}
}

func TestLinearBlend_NearHole(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	demData := makeFlatDEM(w, h, 100)
	mask := makeFlatDEM(w, h, 200)
	for y := 4; y <= 5; y++ {
		for x := 4; x <= 5; x++ {
			mask[y*w+x] = nd
		}
	}
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)

	result := linearBlend(demData, mask, reg, nd, 4)

	if math.Abs(result[4*w+3]-125) > 1e-6 {
		t.Errorf("pixel adjacent to mask hole: expected 125, got %.4f", result[4*w+3])
	}
	if result[0] != 100 {
		t.Errorf("pixel far from mask hole: expected 100, got %.2f", result[0])
	}
	if result[4*w+4] != 100 {
		t.Errorf("mask hole pixel: should keep dem value 100, got %.2f", result[4*w+4])
	}
}

func TestLinearBlend_PartialBlend(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	demData := makeFlatDEM(w, h, 100)
	mask := makeFlatDEM(w, h, 200)
	mask[4*w+5] = nd
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)

	result := linearBlend(demData, mask, reg, nd, 4)

	if result[4*w+5] != 100 {
		t.Errorf("hole pixel should keep dem value, got %.2f", result[4*w+5])
	}
	v := result[4*w+4]
	if v <= 100 || v >= 200 {
		t.Errorf("pixel adjacent to hole should blend between 100 and 200, got %.2f", v)
	}
	if result[0] != 100 {
		t.Errorf("pixel far from hole should be unchanged, got %.2f", result[0])
	}
}

func TestLinearBlendResampled_Blends(t *testing.T) {
	nd := -9999.0
	w, h := 10, 10
	demData := makeFlatDEM(w, h, 100)
	mask := makeFlatDEM(5, 5, 200)
	mask[2*5+2] = nd

	reg1 := dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
	reg2 := dem.NewRegionFromBBox(0, 0, 10, 10, nil, 2, 2)

	result := linearBlendResampled(demData, reg1, mask, reg2, nd, 4)

	if result[4*w+4] != 100 {
		t.Errorf("resampled: pixel over mask hole should keep dem value, got %.2f", result[4*w+4])
	}
	v := result[4*w+3]
	if v <= 100 || v >= 200 {
		t.Errorf("resampled: pixel near mask hole should blend, got %.2f", v)
	}
	if result[0] != 100 {
		t.Errorf("resampled: far pixel should be unchanged, got %.2f", result[0])
	}
}

func TestLinearBlendResampled_SameSize(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	dem1 := makeFlatDEM(w, h, 100)
	dem2 := makeFlatDEM(w, h, 200)
	reg1 := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	reg2 := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)

	result := linearBlendResampled(dem1, reg1, dem2, reg2, nd, 2)
	if len(result) != w*h {
		t.Errorf("output size mismatch: %d vs %d", len(result), w*h)
	}
	for i, v := range result {
		if math.Abs(v-dem1[i]) > 1e-10 {
			t.Errorf("same-size no holes: pixel %d should be unchanged, got %.2f", i, v)
		}
	}
}
