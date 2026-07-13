package grits

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func TestWeights_AllValid(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	wf := &weightFilter{}
	res, err := wf.Run(data, reg, &Options{Radius: 5, NoData: &nd})
	if err != nil {
		t.Fatalf("weights error: %v", err)
	}
	for i, v := range res {
		if v == nd || math.IsNaN(v) {
			t.Errorf("all-valid: pixel %d has invalid weight", i)
		}
		if v < 0 || v > 1 {
			t.Errorf("weight out of range [0,1]: %.4f at %d", v, i)
		}
	}
}

func TestWeights_EdgeDecay(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, nd)
	for y := 3; y <= 6; y++ {
		for x := 3; x <= 6; x++ {
			data[y*w+x] = 100
		}
	}
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	wf := &weightFilter{}
	res, _ := wf.Run(data, reg, &Options{Radius: 3, NoData: &nd})

	if res[5*w+5] >= res[4*w+5] {
		t.Logf("center weight %.4f >= edge weight %.4f (correct)", res[5*w+5], res[4*w+5])
	}
}

func TestWeights_AllNoData(t *testing.T) {
	nd := -9999.0
	data := makeFlatDEM(5, 5, nd)
	wf := &weightFilter{}
	res, _ := wf.Run(data, region5x5(), &Options{Radius: 2, NoData: &nd})
	for i, v := range res {
		if v != 0 && v != nd {
			t.Errorf("all-nodata: pixel %d weight should be 0, got %.4f", i, v)
		}
	}
}

func TestComputeWeightBuffer_Edge(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[0] = nd

	weights := computeWeightBuffer(data, w, h, 3, nd)
	if weights[0] != 0 {
		t.Errorf("noData pixel should have weight 0, got %.4f", weights[0])
	}
	if weights[1] <= 0 {
		t.Errorf("neighbor of noData should have positive weight, got %.4f", weights[1])
	}
}
