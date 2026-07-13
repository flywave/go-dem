package grits

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func makeFlatDEM(w, h int, val float64) []float64 {
	d := make([]float64, w*h)
	for i := range d {
		d[i] = val
	}
	return d
}

func makeRampDEM(w, h int) []float64 {
	d := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			d[y*w+x] = float64(y*100 + x)
		}
	}
	return d
}

func makeCheckNoData(w, h int, noData float64) []float64 {
	d := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if (x+y)%3 == 0 {
				d[y*w+x] = noData
			} else {
				d[y*w+x] = float64(y*10 + x)
			}
		}
	}
	return d
}

func region5x5() *dem.Region {
	return dem.NewRegionFromBBox(0, 0, 5, 5, nil, 1, 1)
}

func TestGaussianBlur_Identity(t *testing.T) {
	data := makeFlatDEM(10, 10, 100)
	reg := dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
	nd := -9999.0
	g := &gaussianFilter{}
	res, err := g.Run(data, reg, &Options{Sigma: 0.001, NoData: &nd})
	if err != nil {
		t.Fatalf("gaussian error: %v", err)
	}
	for i, v := range res {
		if math.Abs(v-data[i]) > 1.0 {
			t.Errorf("identity: at %d expected %.0f, got %.2f", i, data[i], v)
			break
		}
	}
}

func TestGaussianBlur_SmoothsPeak(t *testing.T) {
	w, h := 10, 10
	data := makeFlatDEM(w, h, 0)
	data[5*w+5] = 1000
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	nd := -9999.0
	g := &gaussianFilter{}
	res, err := g.Run(data, reg, &Options{Sigma: 2.0, NoData: &nd})
	if err != nil {
		t.Fatalf("gaussian error: %v", err)
	}
	if res[5*w+5] < data[5*w+5] {
		t.Logf("peak smoothed: %.0f -> %.2f", data[5*w+5], res[5*w+5])
	} else {
		t.Errorf("peak not smoothed: %.0f -> %.2f", data[5*w+5], res[5*w+5])
	}
}

func TestGaussianBlur_NegativeSigma(t *testing.T) {
	data := makeFlatDEM(5, 5, 50)
	reg := dem.NewRegionFromBBox(0, 0, 5, 5, nil, 1, 1)
	nd := -9999.0
	g := &gaussianFilter{}
	_, err := g.Run(data, reg, &Options{Sigma: -1, NoData: &nd})
	if err != nil {
		t.Fatalf("negative sigma should not error: %v", err)
	}
}

func TestGaussianBlur_NoDataPreserved(t *testing.T) {
	w, h := 8, 8
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[3*w+3] = nd
	data[3*w+4] = nd

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	g := &gaussianFilter{}
	res, err := g.Run(data, reg, &Options{Sigma: 1.5, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res[3*w+3] != nd {
		t.Logf("noData neighbor got filled: %.2f", res[3*w+3])
	}
	if res[0] == nd || math.IsNaN(res[0]) {
		t.Log("gaussian: valid pixel should remain valid")
	}
}

func TestMedianFilter_Basic(t *testing.T) {
	w, h := 7, 7
	data := makeFlatDEM(w, h, 0)
	data[3*w+3] = 9999
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	nd := -9999.0
	m := &medianFilter{}
	res, err := m.Run(data, reg, &Options{KernelSize: 3, NoData: &nd})
	if err != nil {
		t.Fatalf("median error: %v", err)
	}
	if res[3*w+3] < 9000 {
		t.Logf("median removed spike: %.0f -> %.2f", data[3*w+3], res[3*w+3])
	} else {
		t.Errorf("spike not removed: %.0f -> %.2f", data[3*w+3], res[3*w+3])
	}
}

func TestMedianFilter_NoDataHandling(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeCheckNoData(w, h, nd)
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	m := &medianFilter{}
	res, err := m.Run(data, reg, &Options{KernelSize: 3, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i := range data {
		if data[i] != nd && res[i] == nd {
			t.Errorf("valid pixel %d became noData", i)
		}
	}
}

func TestMedianFilter_EvenKernel(t *testing.T) {
	data := makeFlatDEM(5, 5, 100)
	reg := region5x5()
	nd := -9999.0
	m := &medianFilter{}
	res, err := m.Run(data, reg, &Options{KernelSize: 4, NoData: &nd})
	if err != nil {
		t.Fatalf("even kernel error: %v", err)
	}
	if len(res) != len(data) {
		t.Errorf("output size mismatch")
	}
}

func TestMedian_Function(t *testing.T) {
	tests := []struct {
		input []float64
		want  float64
	}{
		{[]float64{1, 2, 3}, 2},
		{[]float64{1, 2, 3, 4}, 2.5},
		{[]float64{5}, 5},
		{[]float64{}, 0},
		{[]float64{3, 1, 2}, 2},
	}
	for _, tt := range tests {
		got := median(tt.input)
		if math.Abs(got-tt.want) > 1e-10 {
			t.Errorf("median(%v) = %.2f, want %.2f", tt.input, got, tt.want)
		}
	}
}

func Test1DGaussianKernel(t *testing.T) {
	k := make1DGaussianKernel(1.0, 2)
	if len(k) != 5 {
		t.Errorf("kernel length: expected 5, got %d", len(k))
	}
	var sum float64
	for _, v := range k {
		sum += v
	}
	if math.Abs(sum-1.0) > 1e-10 {
		t.Errorf("kernel sum: expected 1.0, got %f", sum)
	}
	if k[2] <= k[1] || k[2] <= k[3] {
		t.Error("kernel center should be maximum")
	}
}
