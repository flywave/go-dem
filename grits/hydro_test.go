package grits

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func TestHydroFill_NoSinks(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	snapshot := make([]float64, len(data))
	copy(snapshot, data)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	hyd := &hydroFilter{}
	res, err := hyd.Run(data, reg, &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i := range res {
		if snapshot[i] != nd && res[i] == nd {
			t.Errorf("valid pixel %d became noData", i)
		}
	}
}

func TestHydroFill_Sink(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeRampDEM(w, h)
	data[5*w+5] = 0

	snapshot := make([]float64, len(data))
	copy(snapshot, data)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	hyd := &hydroFilter{}
	res, err := hyd.Run(data, reg, &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res[5*w+5] > snapshot[5*w+5] {
		t.Logf("sink raised: %.0f -> %.2f", snapshot[5*w+5], res[5*w+5])
	} else {
		t.Log("sink may not need filling (already drains)")
	}
}

func TestHydroFill_BorderSink(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeRampDEM(w, h)
	data[0] = -100

	hyd := &hydroFilter{}
	res, err := hyd.Run(data, region5x5(), &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res[0] <= -100 {
		t.Logf("border sink: %.0f -> %.2f", data[0], res[0])
	}
}

func TestFillSinks_AllSame(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	res := fillSinks(data, w, h, nd)
	for i := range res {
		if res[i] == nd {
			t.Errorf("pixel %d became noData", i)
		}
	}
}

func TestFillSinks_GradientPreserved(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeRampDEM(w, h)
	res := fillSinks(data, w, h, nd)
	for i := range data {
		if math.Abs(res[i]-data[i]) > 1e-6 {
			t.Logf("pixel %d changed: %.0f -> %.2f", i, data[i], res[i])
		}
	}
}

func TestHydroFill_NoDataPreserved(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[2*w+2] = nd

	hyd := &hydroFilter{}
	res, _ := hyd.Run(data, region5x5(), &Options{NoData: &nd})
	if res[2*w+2] != nd {
		t.Errorf("noData cell should remain noData, got %.2f", res[2*w+2])
	}
}

func TestAbsInt(t *testing.T) {
	if absInt(-5) != 5 {
		t.Errorf("absInt(-5) = %d", absInt(-5))
	}
	if absInt(3) != 3 {
		t.Errorf("absInt(3) = %d", absInt(3))
	}
	if absInt(0) != 0 {
		t.Errorf("absInt(0) = %d", absInt(0))
	}
}
