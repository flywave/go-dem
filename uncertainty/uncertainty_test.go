package uncertainty

import (
	"math"
	"os"
	"path/filepath"
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
	if got := res.Proximity[5*w+5]; got != -9999 {
		t.Errorf("hole proximity: expected -9999, got %v", got)
	}
	if got := res.Proximity[4*w+5]; math.Abs(got-1) > 1e-9 {
		t.Errorf("neighbor of hole proximity: expected 1, got %v", got)
	}
}

func TestProximityUncertainty_IsolatedPixel(t *testing.T) {
	w, h := 5, 5
	data := make([]float64, w*h)
	for i := range data {
		data[i] = -9999
	}
	data[2*w+2] = 100
	reg := region(w, h)
	res, err := Estimate(data, reg, &Options{Method: MethodProximity})
	if err != nil {
		t.Fatalf("proximity error: %v", err)
	}
	if got := res.Proximity[2*w+2]; got != dem.DefaultNoData {
		t.Errorf("isolated pixel proximity: expected %v, got %v", dem.DefaultNoData, got)
	}
	if got := res.TotalUncertainty[2*w+2]; got != dem.DefaultNoData {
		t.Errorf("isolated pixel uncertainty: expected %v, got %v", dem.DefaultNoData, got)
	}
}

func TestProximityUncertainty_NaNNeighborExcluded(t *testing.T) {
	w, h := 5, 5
	data := make([]float64, w*h)
	for i := range data {
		data[i] = 100
	}
	data[1*w+2] = math.NaN()
	data[3*w+2] = math.NaN()
	data[2*w+1] = math.NaN()
	data[2*w+3] = math.NaN()
	reg := region(w, h)
	res, err := Estimate(data, reg, &Options{Method: MethodProximity, NoData: -9999})
	if err != nil {
		t.Fatalf("proximity error: %v", err)
	}
	if got := res.Proximity[2*w+2]; math.Abs(got-math.Sqrt2) > 1e-9 {
		t.Errorf("proximity with NaN orthogonal neighbors: expected %v, got %v", math.Sqrt2, got)
	}
}

func TestProximityUncertainty_MapUnits(t *testing.T) {
	w, h := 5, 5
	data := make([]float64, w*h)
	for i := range data {
		data[i] = 100
	}
	data[2*w+2] = -9999
	reg := dem.NewRegionFromBBox(0, 0, 10, 10, nil, 2, 2)
	res, err := Estimate(data, reg, &Options{Method: MethodProximity, NoData: -9999})
	if err != nil {
		t.Fatalf("proximity error: %v", err)
	}
	if got := res.Proximity[2*w+1]; math.Abs(got-2) > 1e-9 {
		t.Errorf("proximity in map units (XRes=2): expected 2, got %v", got)
	}
}

func TestCombinedUncertainty_DefaultNoData(t *testing.T) {
	w, h := 12, 12
	data := make([]float64, w*h)
	for i := range data {
		data[i] = 100
	}
	for y := 5; y <= 6; y++ {
		for x := 5; x <= 6; x++ {
			data[y*w+x] = -9999
		}
	}
	reg := region(w, h)
	res, err := Estimate(data, reg, &Options{Method: MethodCombined})
	if err != nil {
		t.Fatalf("combined error: %v", err)
	}
	for y := 5; y <= 6; y++ {
		for x := 5; x <= 6; x++ {
			if got := res.TotalUncertainty[y*w+x]; got != dem.DefaultNoData {
				t.Errorf("nodata cell (%d,%d): expected %v, got %v", x, y, dem.DefaultNoData, got)
			}
		}
	}
	for i, v := range res.TotalUncertainty {
		if v != dem.DefaultNoData && v > 1000 {
			t.Errorf("garbage uncertainty at %d: %v", i, v)
		}
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

func TestFillUncertaintyGaps_KNN(t *testing.T) {
	data := []float64{1, -9999, 4, 5, 6}
	filled := fillUncertaintyGaps(data, 5, 1, -9999)
	if math.IsNaN(filled[1]) {
		t.Fatal("gap filled with NaN")
	}
	if filled[1] <= 1 || filled[1] >= 4 {
		t.Errorf("gap filled outside neighbor value range: %v", filled[1])
	}
}

func TestWriteUncertainty_Extension(t *testing.T) {
	dir := t.TempDir()
	w, h := 4, 4
	srs, err := dem.ParseSRS("EPSG:4326")
	if err != nil {
		t.Fatalf("parse srs: %v", err)
	}
	reg := dem.NewRegionFromBBox(0, 0, 4, 4, srs, 1, 1)
	unc := &Result{
		TotalUncertainty:         make([]float64, w*h),
		SourceUncertainty:        make([]float64, w*h),
		InterpolationUncertainty: make([]float64, w*h),
		Proximity:                make([]float64, w*h),
	}

	cases := []struct{ in, base string }{
		{"dem.tif", "dem"},
		{"dem.TIF", "dem"},
		{"dem.tiff", "dem.tiff"},
		{"dem", "dem"},
	}
	for _, c := range cases {
		path := filepath.Join(dir, c.in)
		if err := WriteUncertainty(unc, reg, path, -9999); err != nil {
			t.Fatalf("write %s: %v", c.in, err)
		}
		for _, suffix := range []string{"_tvu.tif", "_interp_u.tif", "_src_u.tif", "_prox.tif"} {
			if _, err := os.Stat(filepath.Join(dir, c.base+suffix)); err != nil {
				t.Errorf("write %s: missing %s: %v", c.in, c.base+suffix, err)
			}
		}
	}
}
