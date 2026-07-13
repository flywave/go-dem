package uncertainty

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func region(w, h int) *dem.Region {
	return dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
}

func TestProximityUncertainty_AllValid(t *testing.T) {
	w, h := 10, 10
	data := make([]float64, w*h)
	for i := range data {
		data[i] = 100
	}
	reg := region(w, h)
	res, err := Estimate(data, reg, &Options{Method: MethodProximity, NoData: -9999})
	if err != nil {
		t.Fatalf("proximity error: %v", err)
	}
	if res == nil {
		t.Fatal("nil result")
	}
	if len(res.TotalUncertainty) != w*h {
		t.Errorf("output size: %d", len(res.TotalUncertainty))
	}
}

func TestProximityUncertainty_WithHoles(t *testing.T) {
	w, h := 10, 10
	data := make([]float64, w*h)
	for i := range data {
		data[i] = 100
	}
	data[5*w+5] = -9999
	reg := region(w, h)
	res, err := Estimate(data, reg, &Options{Method: MethodProximity, NoData: -9999})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res.Proximity[5*w+5] >= 0 {
		t.Logf("hole proximity: %.2f", res.Proximity[5*w+5])
	}
}

func TestUncertainty_InvalidMethod(t *testing.T) {
	reg := region(5, 5)
	_, err := Estimate(nil, reg, &Options{Method: "invalid"})
	if err == nil {
		t.Error("expected error for invalid method")
	}
}

func TestFillUncertaintyGaps_NoGaps(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	filled := fillUncertaintyGaps(data, 5, 1, -9999)
	for i, v := range filled {
		if v != data[i] {
			t.Errorf("no-gap: pixel %d changed %.0f->%.2f", i, data[i], v)
		}
	}
}

func TestFillUncertaintyGaps_WithGaps(t *testing.T) {
	data := []float64{1, -9999, -9999, 4, 5}
	filled := fillUncertaintyGaps(data, 5, 1, -9999)
	if filled[1] == -9999 || math.IsNaN(filled[1]) {
		t.Error("gap not filled")
	}
}

func TestFillUncertaintyGaps_AllGaps(t *testing.T) {
	data := []float64{-9999, -9999, -9999}
	filled := fillUncertaintyGaps(data, 3, 1, -9999)
	for i, v := range filled {
		if v != data[i] {
			t.Errorf("all-gap: pixel %d changed", i)
		}
	}
}
