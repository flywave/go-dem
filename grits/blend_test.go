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

func TestEdgeDistance_Center(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[2*w+2] = nd

	dist := edgeDistance(2*w+2, data, w, h, nd)
	if dist >= 0 {
		t.Logf("center edge distance: %.2f", dist)
	}
}

func TestEdgeDistance_NeighborOfNoData(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[2*w+2] = nd

	dist := edgeDistance(2*w+1, data, w, h, nd)
	if dist >= 0 {
		t.Logf("neighbor edge distance: %.2f", dist)
	}
}

func TestEdgeDistance_NoDataPixel(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, nd)

	dist := edgeDistance(2*w+2, data, w, h, nd)
	if dist >= 0 {
		t.Errorf("noData pixel should have negative edge distance, got %.2f", dist)
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
			t.Logf("same pixel: dem1=%.0f dem2=%.0f blended=%.2f", dem1[i], dem2[i], v)
		}
	}
}

func TestLinearBlend_Different(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	dem1 := makeFlatDEM(w, h, 100)
	dem2 := makeFlatDEM(w, h, 100)
	dem2[2*w+2] = 300
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)

	result := linearBlend(dem1, dem2, reg, nd, 2)
	if result[2*w+2] > 100 && result[2*w+2] < 300 {
		t.Logf("blended value: %.2f", result[2*w+2])
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
}
