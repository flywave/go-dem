package grits

import (
	"testing"

	"github.com/flywave/go-dem"
)

func TestDenoise_Median(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[5*w+5] = 9999

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	d := &denoiseFilter{}
	res, err := d.Run(data, reg, &Options{
		Method:     "median",
		KernelSize: 3,
		NoData:     &nd,
	})
	if err != nil {
		t.Fatalf("denoise median error: %v", err)
	}
	if res[5*w+5] < data[5*w+5] {
		t.Logf("denoise median: spike reduced %.0f -> %.2f", data[5*w+5], res[5*w+5])
	} else {
		t.Errorf("spike not reduced: %.0f -> %.2f", data[5*w+5], res[5*w+5])
	}
}

func TestDenoise_Bilateral(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] += float64(y * 2)
		}
	}
	data[5*w+5] = 500

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	d := &denoiseFilter{}
	res, err := d.Run(data, reg, &Options{
		Method:     "bilateral",
		Sigma:      2.0,
		SigmaColor: 200.0,
		Radius:     3,
		NoData:     &nd,
	})
	if err != nil {
		t.Fatalf("denoise bilateral error: %v", err)
	}
	if res[5*w+5] < data[5*w+5] {
		t.Logf("denoise bilateral: spike reduced %.0f -> %.2f", data[5*w+5], res[5*w+5])
	} else {
		t.Errorf("spike not reduced: %.0f -> %.2f", data[5*w+5], res[5*w+5])
	}
}

func TestBilateral_AllFlat(t *testing.T) {
	w, h := 6, 6
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	b := &bilateralFilter{}
	res, err := b.Run(data, reg, &Options{
		Sigma:      1.0,
		SigmaColor: 20.0,
		Radius:     2,
		NoData:     &nd,
	})
	if err != nil {
		t.Fatalf("bilateral error: %v", err)
	}
	eps := 1e-8
	for i, v := range res {
		if v == nd {
			t.Errorf("flat pixel %d became noData", i)
		}
		if v < data[i]-eps || v > data[i]+eps {
			t.Errorf("flat pixel %d changed: %.10f -> %.10f", i, data[i], v)
		}
	}
}

func TestDenoise_NodataHandling(t *testing.T) {
	w, h := 8, 8
	nd := -9999.0
	data := makeCheckNoData(w, h, nd)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	d := &denoiseFilter{}
	res, err := d.Run(data, reg, &Options{
		Method:     "median",
		KernelSize: 3,
		NoData:     &nd,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i := range data {
		if data[i] == nd && res[i] != nd {
			t.Errorf("nodata pixel %d became valid", i)
		}
	}
}
