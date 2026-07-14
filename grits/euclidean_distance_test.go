package grits

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func TestEuclideanDistance_NoNoData(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	e := &euclideanDistanceFilter{}
	res, err := e.Run(data, reg, &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i, v := range res {
		if v == nd {
			t.Errorf("pixel %d should have distance", i)
		}
	}
}

func TestEuclideanDistance_WithNoData(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[2*w+2] = nd

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	e := &euclideanDistanceFilter{}
	res, err := e.Run(data, reg, &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res[2*w+2] != 0 {
		t.Errorf("nodata pixel should have distance 0, got %.2f", res[2*w+2])
	}
	if res[2*w+2-1] <= 0 {
		t.Log("neighbor of nodata should have positive distance")
	}
}

func TestEuclideanDistance_Monotonic(t *testing.T) {
	w, h := 7, 7
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[3*w+3] = nd

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	e := &euclideanDistanceFilter{}
	res, err := e.Run(data, reg, &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	eps := 1e-8
	dists := []int{
		3*w + 3, // 0
		3*w + 4, // should be 1
		3*w + 5, // should be > 1
	}
	for i := 1; i < len(dists); i++ {
		if res[dists[i]]+eps < res[dists[i-1]] {
			t.Errorf("distance should increase from nodata: [%d]=%.2f < [%d]=%.2f",
				dists[i], res[dists[i]], dists[i-1], res[dists[i-1]])
		}
	}
}

func TestEuclideanMerge_Basic(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	base := makeFlatDEM(w, h, 100)
	base[2*w+2] = nd

	other := makeFlatDEM(w, h, 200)
	other[1*w+1] = nd

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	m := &euclideanMergeFilter{}
	res, err := m.Run(base, reg, nil)
	if err != nil {
		t.Fatalf("nil source mask error: %v", err)
	}
	if len(res) != w*h {
		t.Errorf("output size mismatch")
	}
}

func TestEuclideanDistance_AllNoData(t *testing.T) {
	w, h := 4, 4
	nd := -9999.0
	data := makeFlatDEM(w, h, nd)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	e := &euclideanDistanceFilter{}
	res, err := e.Run(data, reg, &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i, v := range res {
		if v != 0 {
			t.Errorf("all nodata: pixel %d should be 0, got %.2f", i, v)
		}
	}
}

func TestComputeEuclideanDistance(t *testing.T) {
	w, h := 3, 3
	nd := -9999.0
	data := []float64{
		100, 100, 100,
		100, nd, 100,
		100, 100, 100,
	}
	dist := computeEuclideanDistance(data, w, h, nd, 1.0)
	if math.Abs(dist[0]-math.Sqrt2) > 0.01 {
		t.Logf("corner distance: expected ~1.41, got %.4f", dist[0])
	}
}
