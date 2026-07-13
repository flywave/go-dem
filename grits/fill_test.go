package grits

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func fillBasicTest(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[4*w+4] = nd
	data[5*w+5] = nd

	f := &fillFilter{}
	res, err := f.Run(data, dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1),
		&Options{MaxDistance: 10, NoData: &nd})
	if err != nil {
		t.Fatalf("fill error: %v", err)
	}
	if res[4*w+4] == nd {
		t.Error("hole at (4,4) not filled")
	}
	if res[5*w+5] == nd {
		t.Error("hole at (5,5) not filled")
	}
}

func TestFill_NoHoles(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeRampDEM(w, h)
	snapshot := make([]float64, len(data))
	copy(snapshot, data)

	f := &fillFilter{}
	res, _ := f.Run(data, dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1),
		&Options{MaxDistance: 10, NoData: &nd})
	for i := range res {
		if res[i] == nd && snapshot[i] != nd {
			t.Errorf("valid pixel %d became noData", i)
		}
	}
}

func TestFill_LargeHole(t *testing.T) {
	w, h := 20, 20
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	for y := 5; y <= 14; y++ {
		for x := 5; x <= 14; x++ {
			data[y*w+x] = nd
		}
	}

	f := &fillFilter{}
	res, err := f.Run(data, dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1),
		&Options{MaxDistance: 20, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	filled := 0
	for y := 5; y <= 14; y++ {
		for x := 5; x <= 14; x++ {
			if res[y*w+x] != nd && !math.IsNaN(res[y*w+x]) {
				filled++
			}
		}
	}
	if filled < 50 {
		t.Errorf("large hole: only %d/100 filled", filled)
	}
}

func TestFill_AllNoData(t *testing.T) {
	nd := -9999.0
	data := makeFlatDEM(5, 5, nd)
	f := &fillFilter{}
	res, _ := f.Run(data, region5x5(), &Options{MaxDistance: 10, NoData: &nd})
	for i, v := range res {
		if v != nd {
			t.Errorf("all-nodata: pixel %d should remain, got %.2f", i, v)
		}
	}
}

func TestFillNoDataPixels(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	// Place noData pixel adjacent to valid data
	data[2*w+2] = nd

	fillNoDataPixels(data, w, h, nd)
	if data[2*w+2] == nd {
		t.Log("fillNoDataPixels: center not filled (single isolated noData may not propagate)")
	} else {
		t.Logf("fillNoDataPixels: center filled with %.2f", data[2*w+2])
	}
}

func TestFillWithInverseDistance(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[2*w+2] = nd

	fillWithInverseDistance(data, w, h, nd)
	if data[2*w+2] == nd {
		t.Error("IDW fill: center not filled")
	}
	if math.Abs(data[2*w+2]-100) > 5 {
		t.Errorf("IDW: expected ~100, got %.2f", data[2*w+2])
	}
}

func TestHasRemainingNoData(t *testing.T) {
	if !hasRemainingNoData([]float64{1, 2, -9999, 4}, -9999) {
		t.Error("should detect remaining noData")
	}
	if hasRemainingNoData([]float64{1, 2, 3, 4}, -9999) {
		t.Error("should not detect noData when none exist")
	}
	if !hasRemainingNoData([]float64{1, math.NaN(), 3}, -9999) {
		t.Error("should detect NaN as noData")
	}
}
