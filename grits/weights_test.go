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
		if math.Abs(v-1.0) > 1e-9 {
			t.Errorf("all-valid: pixel %d weight should be 1.0, got %.4f", i, v)
		}
	}
}

func TestWeights_EdgeDecay(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	for y := 0; y < h; y++ {
		data[y*w] = nd
	}
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	wf := &weightFilter{}
	res, err := wf.Run(data, reg, &Options{Radius: 5, NoData: &nd})
	if err != nil {
		t.Fatalf("weights error: %v", err)
	}

	if res[4*w] != 0 {
		t.Errorf("noData pixel weight should be 0, got %.4f", res[4*w])
	}
	w1 := res[4*w+1]
	w3 := res[4*w+3]
	w5 := res[4*w+5]
	if math.Abs(w1-(1-1.0/5)) > 1e-9 {
		t.Errorf("weight at dist 1: expected %.4f, got %.4f", 1-1.0/5, w1)
	}
	if math.Abs(w3-(1-3.0/5)) > 1e-9 {
		t.Errorf("weight at dist 3: expected %.4f, got %.4f", 1-3.0/5, w3)
	}
	if math.Abs(w5-0.1) > 1e-9 {
		t.Errorf("weight at dist >= radius should clamp to 0.1, got %.4f", w5)
	}
	if !(w1 > w3 && w3 >= w5) {
		t.Errorf("weight should decay with distance from noData: %.4f, %.4f, %.4f", w1, w3, w5)
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
	expected := 1 - 1.0/3
	if math.Abs(weights[1]-expected) > 1e-9 {
		t.Errorf("neighbor of noData (dist 1, radius 3): expected %.4f, got %.4f", expected, weights[1])
	}
	farExpected := 1 - 5.0/3
	if weights[4*w+5] < farExpected-1e-9 {
		t.Errorf("interior weight should not be below floor computation, got %.4f", weights[4*w+5])
	}
}
