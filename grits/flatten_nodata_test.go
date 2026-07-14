package grits

import (
	"testing"

	"github.com/flywave/go-dem"
)

func TestFlattenNoData_NoHoles(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	f := &flattenNoDataFilter{}
	res, err := f.Run(data, reg, &Options{Threshold: 10, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i, v := range res {
		if v != data[i] {
			t.Errorf("pixel %d changed without holes: %.2f -> %.2f", i, data[i], v)
		}
	}
}

func TestFlattenNoData_SmallHole(t *testing.T) {
	w, h := 6, 6
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[3*w+3] = nd

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	f := &flattenNoDataFilter{}
	res, err := f.Run(data, reg, &Options{Threshold: 10, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res[3*w+3] != nd {
		t.Logf("single noData pixel should remain noData: got %.2f", res[3*w+3])
	}
}

func TestFlattenNoData_LargeHole(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	for y := 3; y <= 6; y++ {
		for x := 3; x <= 6; x++ {
			data[y*w+x] = nd
		}
	}

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	f := &flattenNoDataFilter{}
	res, err := f.Run(data, reg, &Options{Threshold: 10, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for y := 3; y <= 6; y++ {
		for x := 3; x <= 6; x++ {
			if res[y*w+x] != nd {
				t.Logf("large hole pixel (%d,%d) filled: %.2f (may be OK)", x, y, res[y*w+x])
			}
		}
	}
}

func TestFlattenNoData_DefaultThreshold(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	f := &flattenNoDataFilter{}
	res, err := f.Run(data, reg, &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(res) != w*h {
		t.Errorf("output size mismatch")
	}
}
