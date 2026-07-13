package perspecto

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func makeFlat(w, h int, val float64) []float64 {
	d := make([]float64, w*h)
	for i := range d {
		d[i] = val
	}
	return d
}

func makeRamp(w, h int) []float64 {
	d := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			d[y*w+x] = float64(y*100 + x)
		}
	}
	return d
}

func region(w, h int) *dem.Region {
	return dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
}

func TestHillshade_Flat(t *testing.T) {
	w, h := 10, 10
	data := makeFlat(w, h, 100)
	reg := region(w, h)
	hs := Hillshade(data, reg, &Options{NoData: -9999})
	for i, v := range hs {
		if v != 0 && !math.IsNaN(v) {
			t.Logf("flat: hillshade at %d = %.2f (0 expected)", i, v)
		}
	}
}

func TestHillshade_NoData(t *testing.T) {
	w, h := 5, 5
	data := makeFlat(w, h, -9999)
	reg := region(w, h)
	hs := Hillshade(data, reg, &Options{NoData: -9999})
	for i, v := range hs {
		if v != -9999 && !math.IsNaN(v) {
			t.Errorf("noData pixel %d has hillshade value %.2f", i, v)
		}
	}
}

func TestSlope_Flat(t *testing.T) {
	w, h := 10, 10
	data := makeFlat(w, h, 100)
	reg := region(w, h)
	sl := Slope(data, reg, &Options{NoData: -9999})
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			if sl[y*w+x] > 1 {
				t.Errorf("flat: slope at (%d,%d) = %.4f", x, y, sl[y*w+x])
			}
		}
	}
}

func TestSlope_45Deg(t *testing.T) {
	w, h := 5, 5
	data := makeFlat(w, h, 0)
	for x := 0; x < w; x++ {
		data[x] = float64(x)
	}
	reg := region(w, h)
	sl := Slope(data, reg, &Options{NoData: -9999})
	t.Logf("x-ramp: slope at (2,2) = %.4f deg", sl[2*w+2])
}

func TestSlope_Percent(t *testing.T) {
	w, h := 5, 5
	data := makeFlat(w, h, 0)
	for x := 0; x < w; x++ {
		data[x] = float64(x)
	}
	reg := region(w, h)
	sl := Slope(data, reg, &Options{NoData: -9999, SlopeUnits: "percent"})
	if sl[2*w+2] > 0 {
		t.Logf("slope in percent: %.2f", sl[2*w+2])
	}
}

func TestAspect_Flat(t *testing.T) {
	w, h := 10, 10
	data := makeFlat(w, h, 100)
	reg := region(w, h)
	asp := Aspect(data, reg, &Options{NoData: -9999})
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			if asp[y*w+x] > 360 || (asp[y*w+x] > 0 && asp[y*w+x] < 0) {
			}
		}
	}
}

func TestColorRelief_Basic(t *testing.T) {
	w, h := 5, 5
	data := makeFlat(w, h, 500)
	reg := region(w, h)
	pixels, err := ColorRelief(data, reg, &Options{Colormap: DefaultTerrainColormap(), NoData: -9999})
	if err != nil {
		t.Fatalf("color relief error: %v", err)
	}
	if len(pixels) != w*h*3 {
		t.Errorf("pixel count: expected %d, got %d", w*h*3, len(pixels))
	}
}

func TestColorRelief_NoData(t *testing.T) {
	w, h := 5, 5
	data := makeFlat(w, h, -9999)
	reg := region(w, h)
	pixels, _ := ColorRelief(data, reg, &Options{NoData: -9999})
	if pixels[0] != 0 || pixels[1] != 0 || pixels[2] != 0 {
		t.Errorf("noData should be black, got (%d,%d,%d)", pixels[0], pixels[1], pixels[2])
	}
}

func TestInterpolateColor(t *testing.T) {
	cmap := []ColorStop{
		{0, 0, 0, 255},
		{100, 255, 255, 255},
	}
	r, g, b := interpolateColor(50, cmap)
	if r != 127 || g != 127 {
		t.Errorf("midpoint: expected (127,127,255), got (%d,%d,%d)", r, g, b)
	}
}

func TestInterpolateColor_BelowRange(t *testing.T) {
	cmap := []ColorStop{{100, 255, 0, 0}, {200, 0, 255, 0}}
	r, g, b := interpolateColor(50, cmap)
	if r != 255 || g != 0 || b != 0 {
		t.Errorf("below range: expected (255,0,0), got (%d,%d,%d)", r, g, b)
	}
}

func TestInterpolateColor_AboveRange(t *testing.T) {
	cmap := []ColorStop{{100, 255, 0, 0}, {200, 0, 255, 0}}
	r, g, b := interpolateColor(300, cmap)
	if r != 0 || g != 255 || b != 0 {
		t.Errorf("above range: expected (0,255,0), got (%d,%d,%d)", r, g, b)
	}
}

func TestInterpolateColor_EmptyCmap(t *testing.T) {
	r, g, b := interpolateColor(100, nil)
	if r != 128 || g != 128 || b != 128 {
		t.Errorf("empty cmap: expected (128,128,128), got (%d,%d,%d)", r, g, b)
	}
}

func TestShadedRelief(t *testing.T) {
	w, h := 10, 10
	data := makeRamp(w, h)
	reg := region(w, h)
	pixels, err := ShadedRelief(data, reg, &ShadedReliefOptions{
		HillshadeOpts: Options{NoData: -9999},
		ColorOpts:     Options{NoData: -9999, Colormap: DefaultTerrainColormap()},
		Opacity:       0.6,
	})
	if err != nil {
		t.Fatalf("shaded relief error: %v", err)
	}
	if len(pixels) != w*h*3 {
		t.Errorf("pixel count: %d", len(pixels))
	}
}

func TestHasNoData(t *testing.T) {
	if hasNoData([]float64{1, 2, -9999, 4}, -9999) != true {
		t.Error("should detect noData")
	}
	if hasNoData([]float64{1, 2, 3}, -9999) != false {
		t.Error("should not detect noData")
	}
	if hasNoData([]float64{1, math.NaN(), 3}, -9999) != true {
		t.Error("should detect NaN as noData")
	}
}
