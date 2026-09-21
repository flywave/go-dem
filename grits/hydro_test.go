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
	if res[5*w+5] <= snapshot[5*w+5]+1e-9 {
		t.Errorf("sink not raised: %.0f -> %.2f", snapshot[5*w+5], res[5*w+5])
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
	if res[0] != -100 {
		t.Errorf("border pit is an outlet and should stay, got %.2f", res[0])
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
		if res[i] < 100-1e-9 {
			t.Errorf("flat pixel %d lowered: %.4f", i, res[i])
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
			t.Errorf("pixel %d changed: %.0f -> %.2f", i, data[i], res[i])
		}
	}
}

func TestFillFlatAreas_FlatGetsGradient(t *testing.T) {
	w, h := 9, 9
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	for y := 3; y <= 5; y++ {
		for x := 3; x <= 5; x++ {
			data[y*w+x] = 100
		}
	}
	data[0] = 90

	res := fillFlatAreas(data, w, h, nd)
	center := res[4*w+4]
	if math.Abs(center-100) < 1e-9 {
		t.Errorf("flat interior should get epsilon gradient, stayed %.4f", center)
	}
	if center < 100-1e-9 {
		t.Errorf("flat interior should not be lowered, got %.4f", center)
	}
	if math.Abs(res[0]-90) > 1e-9 {
		t.Errorf("outlet pixel should stay, got %.4f", res[0])
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
