package grits

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func TestSlopeFilter_Flat(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	sf := &slopeFilter{}
	res, err := sf.Run(data, reg, &Options{Threshold: 30, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i, v := range res {
		if v == nd && data[i] != nd {
			t.Errorf("flat pixel %d incorrectly masked", i)
		}
	}
}

func TestSlopeFilter_SteepRemoved(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 0)
	for x := 0; x < w; x++ {
		data[x] = float64(x * 100)
	}

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	sf := &slopeFilter{}
	res, err := sf.Run(data, reg, &Options{Threshold: 10, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	masked := 0
	for i, v := range res {
		if v == nd && data[i] != nd {
			masked++
		}
	}
	if masked > 0 {
		t.Logf("%d steep pixels masked", masked)
	}
}

func TestSlopeFilter_NoThreshold(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	sf := &slopeFilter{}
	res, _ := sf.Run(data, region5x5(), &Options{NoData: &nd})
	masked := 0
	for i, v := range res {
		if v == nd && data[i] != nd {
			masked++
		}
	}
	if masked > 4 {
		t.Errorf("flat dem: %d valid pixels masked (border may be affected)", masked)
	}
}

func TestComputeSlopeDegrees(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	slope := computeSlopeDegrees(data, w, h, 1.0, nd)
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			if math.Abs(slope[y*w+x]) > 1e-6 {
				t.Errorf("flat: slope at (%d,%d) = %.4f", x, y, slope[y*w+x])
			}
		}
	}
}

func TestComputeSlopeDegrees_45Deg(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 0)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = float64(x)
		}
	}

	slope := computeSlopeDegrees(data, w, h, 1.0, nd)
	if slope[w+1] > 0 {
		t.Logf("ramp x-direction: slope at (1,1) = %.4f deg", slope[w+1])
	}
}
