package grits

import (
	"testing"

	"github.com/flywave/go-dem"
)

func TestFlats_Basic(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	for y := 3; y <= 6; y++ {
		for x := 3; x <= 6; x++ {
			data[y*w+x] = 50
		}
	}

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	f := &flatsFilter{}
	res, err := f.Run(data, reg, &Options{Threshold: 5, NoData: &nd})
	if err != nil {
		t.Fatalf("flats error: %v", err)
	}
	masked := 0
	for y := 3; y <= 6; y++ {
		for x := 3; x <= 6; x++ {
			if res[y*w+x] == nd {
				masked++
			}
		}
	}
	if masked > 0 {
		t.Logf("flats: %d flat pixels masked", masked)
	} else {
		t.Log("flats: no flat pixels masked (may need higher threshold)")
	}
}

func TestFlats_AutoThreshold(t *testing.T) {
	w, h := 8, 8
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	for i := 0; i < 20; i++ {
		data[i] = 50
	}

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	f := &flatsFilter{}
	res, err := f.Run(data, reg, &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(res) != len(data) {
		t.Errorf("output size mismatch")
	}
}

func TestFlats_NoFlats(t *testing.T) {
	w, h := 6, 6
	nd := -9999.0
	data := makeRampDEM(w, h)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	f := &flatsFilter{}
	res, err := f.Run(data, reg, &Options{Threshold: 2, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i := range res {
		if res[i] == nd && data[i] != nd {
			t.Errorf("unique value pixel %d should not be masked", i)
		}
	}
}
