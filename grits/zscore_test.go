package grits

import (
	"testing"

	"github.com/flywave/go-dem"
)

func TestZScore_Basic(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[5*w+5] = 500

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	z := &zscoreFilter{}
	res, err := z.Run(data, reg, &Options{Threshold: 2.0, KernelSize: 5, NoData: &nd})
	if err != nil {
		t.Fatalf("zscore error: %v", err)
	}
	if res[5*w+5] == nd {
		t.Log("zscore: spike correctly detected")
	} else {
		t.Errorf("spike not masked: %.2f", res[5*w+5])
	}
	if res[0] == nd {
		t.Error("normal pixel should not be masked")
	}
}

func TestZScore_DefaultThreshold(t *testing.T) {
	w, h := 8, 8
	nd := -9999.0
	data := makeFlatDEM(w, h, 50)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	z := &zscoreFilter{}
	res, err := z.Run(data, reg, &Options{KernelSize: 3, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(res) != len(data) {
		t.Errorf("output size mismatch")
	}
}

func TestZScore_AllFlat(t *testing.T) {
	w, h := 6, 6
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	z := &zscoreFilter{}
	res, err := z.Run(data, reg, &Options{Threshold: 3.0, KernelSize: 3, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i := range res {
		if res[i] == nd {
			t.Errorf("flat pixel %d should not be masked", i)
		}
	}
}

func TestZScore_PitDetection(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[5*w+5] = -500

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	z := &zscoreFilter{}
	res, err := z.Run(data, reg, &Options{Threshold: 2.0, KernelSize: 5, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res[5*w+5] == nd {
		t.Log("zscore: pit correctly detected")
	} else {
		t.Errorf("pit not masked: %.2f", res[5*w+5])
	}
}
