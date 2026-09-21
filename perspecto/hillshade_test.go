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
	const flat = 0.5 * 255
	for i, v := range hs {
		if v < 0 || v > 255 {
			t.Errorf("flat hillshade at %d out of range: %.6f", i, v)
		}
	}
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			if math.Abs(hs[y*w+x]-flat) > 1e-6 {
				t.Errorf("flat hillshade interior at (%d,%d): expected %.4f, got %.6f", x, y, flat, hs[y*w+x])
			}
		}
	}
}

func TestHillshade_GDALFormula(t *testing.T) {
	w, h := 5, 5
	reg := region(w, h)

	east := makeGrid(w, h, func(x, y int) float64 { return float64(x) })
	hs := Hillshade(east, reg, &Options{NoData: -9999})
	if math.Abs(hs[2*w+2]-217.6561146) > 1e-4 {
		t.Errorf("east ramp hillshade at (2,2): expected ~217.6561, got %.6f", hs[2*w+2])
	}

	south := makeGrid(w, h, func(x, y int) float64 { return float64(y) })
	hs = Hillshade(south, reg, &Options{NoData: -9999})
	if math.Abs(hs[2*w+2]-217.6561146) > 1e-4 {
		t.Errorf("south ramp hillshade at (2,2): expected ~217.6561, got %.6f", hs[2*w+2])
	}

	diag := makeGrid(w, h, func(x, y int) float64 { return float64(x + 2*y) })
	hs = Hillshade(diag, reg, &Options{NoData: -9999})
	if math.Abs(hs[2*w+2]-234.4364166) > 1e-4 {
		t.Errorf("diagonal ramp hillshade at (2,2): expected ~234.4364, got %.6f", hs[2*w+2])
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
	for i, v := range sl {
		if v != 0 {
			t.Errorf("flat slope at %d: expected 0, got %.6f", i, v)
		}
	}
}

func TestSlope_InteriorExact(t *testing.T) {
	w, h := 5, 5
	data := makeGrid(w, h, func(x, y int) float64 { return float64(x) })
	reg := region(w, h)
	sl := Slope(data, reg, &Options{NoData: -9999})
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			if math.Abs(sl[y*w+x]-45) > 1e-9 {
				t.Errorf("x-ramp slope at (%d,%d): expected 45, got %.6f", x, y, sl[y*w+x])
			}
		}
	}
}

func TestSlope_45Deg(t *testing.T) {
	w, h := 5, 5
	data := makeGrid(w, h, func(x, y int) float64 { return float64(x) })
	reg := region(w, h)
	sl := Slope(data, reg, &Options{NoData: -9999})
	if math.Abs(sl[2*w+2]-45) > 1e-9 {
		t.Errorf("x-ramp slope at (2,2): expected 45 deg, got %.4f", sl[2*w+2])
	}
}

func TestSlope_Percent(t *testing.T) {
	w, h := 5, 5
	data := makeGrid(w, h, func(x, y int) float64 { return float64(x) })
	reg := region(w, h)
	sl := Slope(data, reg, &Options{NoData: -9999, SlopeUnits: "percent"})
	if math.Abs(sl[2*w+2]-100) > 1e-9 {
		t.Errorf("x-ramp slope percent at (2,2): expected 100, got %.4f", sl[2*w+2])
	}
}

func TestSlope_BordersFilled(t *testing.T) {
	w, h := 5, 5
	data := makeGrid(w, h, func(x, y int) float64 { return float64(x) })
	reg := region(w, h)
	sl := Slope(data, reg, &Options{NoData: -9999})
	for x := 0; x < w; x++ {
		if sl[x] <= 0 {
			t.Errorf("top border slope at x=%d: expected >0, got %.4f", x, sl[x])
		}
		if sl[(h-1)*w+x] <= 0 {
			t.Errorf("bottom border slope at x=%d: expected >0, got %.4f", x, sl[(h-1)*w+x])
		}
	}
	for y := 0; y < h; y++ {
		if sl[y*w] <= 0 {
			t.Errorf("left border slope at y=%d: expected >0, got %.4f", y, sl[y*w])
		}
		if sl[y*w+w-1] <= 0 {
			t.Errorf("right border slope at y=%d: expected >0, got %.4f", y, sl[y*w+w-1])
		}
	}
}

func TestSlope_BorderNoNoDataLeak(t *testing.T) {
	w, h := 5, 5
	data := makeGrid(w, h, func(x, y int) float64 { return float64(x) })
	data[2*w+2] = -9999
	reg := region(w, h)
	sl := Slope(data, reg, &Options{NoData: -9999})
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if y == 0 || y == h-1 || x == 0 || x == w-1 {
				if sl[y*w+x] == -9999 {
					t.Errorf("border pixel (%d,%d) leaked noData", x, y)
				}
			}
		}
	}
}

func makeGrid(w, h int, f func(x, y int) float64) []float64 {
	d := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			d[y*w+x] = f(x, y)
		}
	}
	return d
}

func TestAspect_Flat(t *testing.T) {
	w, h := 10, 10
	data := makeFlat(w, h, 100)
	reg := region(w, h)
	asp := Aspect(data, reg, &Options{NoData: -9999})
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			if math.Abs(asp[y*w+x]-270) > 1e-6 {
				t.Errorf("flat aspect at (%d,%d): expected 270, got %.4f", x, y, asp[y*w+x])
			}
		}
	}
}

func TestAspect_Directions(t *testing.T) {
	w, h := 5, 5
	reg := region(w, h)

	cases := []struct {
		name string
		data []float64
		want float64
	}{
		{"downhill west", makeGrid(w, h, func(x, y int) float64 { return float64(x) }), 270},
		{"downhill north", makeGrid(w, h, func(x, y int) float64 { return float64(y) }), 0},
		{"downhill south", makeGrid(w, h, func(x, y int) float64 { return float64(4 - y) }), 180},
		{"downhill east", makeGrid(w, h, func(x, y int) float64 { return float64(4 - x) }), 90},
	}
	for _, c := range cases {
		asp := Aspect(c.data, reg, &Options{NoData: -9999})
		for y := 1; y < h-1; y++ {
			for x := 1; x < w-1; x++ {
				if math.Abs(asp[y*w+x]-c.want) > 1e-6 {
					t.Errorf("%s: aspect at (%d,%d): expected %.1f, got %.4f", c.name, x, y, c.want, asp[y*w+x])
				}
			}
		}
	}
}

func TestAspect_AnisotropicRes(t *testing.T) {
	w, h := 5, 5
	reg := dem.NewRegionFromBBox(0, 0, 5, 10, nil, 1, 2)
	data := makeGrid(w, h, func(x, y int) float64 { return float64(x + y) })
	asp := Aspect(data, reg, &Options{NoData: -9999})
	if math.Abs(asp[2*w+2]-296.5650512) > 1e-4 {
		t.Errorf("anisotropic aspect at (2,2): expected ~296.5651, got %.4f", asp[2*w+2])
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

func TestShadedRelief_OpacityClamped(t *testing.T) {
	w, h := 10, 10
	data := makeRamp(w, h)
	reg := region(w, h)
	mk := func(op float64) []uint8 {
		pixels, err := ShadedRelief(data, reg, &ShadedReliefOptions{
			HillshadeOpts: Options{NoData: -9999},
			ColorOpts:     Options{NoData: -9999, Colormap: DefaultTerrainColormap()},
			Opacity:       op,
		})
		if err != nil {
			t.Fatalf("shaded relief error: %v", err)
		}
		return pixels
	}
	p1 := mk(1)
	p2 := mk(2)
	for i := range p1 {
		if p1[i] != p2[i] {
			t.Errorf("opacity 2 differs from opacity 1 at %d: %d vs %d", i, p1[i], p2[i])
		}
	}
}

func TestInterpolateColor_EqualAdjacentStops(t *testing.T) {
	cmap := []ColorStop{{100, 10, 20, 30}, {100, 40, 50, 60}, {200, 70, 80, 90}}
	r, g, b := interpolateColor(100, cmap)
	if r != 10 || g != 20 || b != 30 {
		t.Errorf("equal adjacent stops: expected first stop (10,20,30), got (%d,%d,%d)", r, g, b)
	}
	r, g, b = InterpolateColor(cmap, 100)
	if r != 10 || g != 20 || b != 30 {
		t.Errorf("equal adjacent stops (exported): expected first stop (10,20,30), got (%d,%d,%d)", r, g, b)
	}
}

func TestHistogramPNG_FullRangeBins(t *testing.T) {
	data := make([]float64, 0, 12)
	data = append(data, 0)
	for i := 0; i < 10; i++ {
		data = append(data, 0.75)
	}
	data = append(data, 1)
	img, err := HistogramPNG(data, &HistogramOptions{Bins: 2})
	if err != nil {
		t.Fatalf("histogram error: %v", err)
	}
	r0, g0, b0, a0 := img.At(200, 300).RGBA()
	if r0 != 65535 || g0 != 65535 || b0 != 65535 || a0 != 65535 {
		t.Errorf("bin 0 (1 value) at (200,300) should be background, got (%d,%d,%d,%d)", r0, g0, b0, a0)
	}
	r1, g1, b1, a1 := img.At(500, 300).RGBA()
	if r1 == 65535 && g1 == 65535 && b1 == 65535 && a1 == 65535 {
		t.Errorf("bin 1 (11 values) at (500,300) should be bar, got background")
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
