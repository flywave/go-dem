package perspecto

type Options struct {
	NoData   float64
	ZFactor  float64
	Azimuth  float64
	Altitude float64
	Colormap []ColorStop
	SlopeUnits string
}

type HistogramOptions struct {
	Bins     int
	Type     string
	Width    int
	Height   int
	Title    string
	ShowStats bool
	NoData   float64
}

type ColorbarOptions struct {
	Width  int
	Height int
	Label  string
}

type ColorStop struct {
	Value    float64
	R, G, B  uint8
}

func DefaultTerrainColormap() []ColorStop {
	return []ColorStop{
		{0, 0, 100, 0},
		{200, 34, 139, 34},
		{500, 107, 142, 35},
		{1000, 139, 119, 26},
		{1500, 160, 120, 60},
		{2000, 139, 90, 43},
		{3000, 160, 130, 100},
		{4000, 200, 180, 150},
		{5000, 240, 230, 210},
	}
}

func DefaultBathymetryColormap() []ColorStop {
	return []ColorStop{
		{-8000, 10, 10, 80},
		{-4000, 20, 40, 120},
		{-2000, 30, 80, 150},
		{-500, 50, 130, 180},
		{-100, 80, 170, 200},
		{-20, 140, 210, 230},
		{0, 200, 230, 240},
	}
}

func RescaleColormap(cmap []ColorStop, minVal, maxVal float64) []ColorStop {
	if len(cmap) == 0 || maxVal-minVal == 0 {
		return cmap
	}
	oldMin, oldMax := cmap[0].Value, cmap[len(cmap)-1].Value
	if oldMax-oldMin == 0 {
		return cmap
	}
	out := make([]ColorStop, len(cmap))
	for i, c := range cmap {
		f := (c.Value - oldMin) / (oldMax - oldMin)
		out[i] = ColorStop{
			Value: minVal + f*(maxVal-minVal),
			R:     c.R, G: c.G, B: c.B,
		}
	}
	return out
}

func InterpolateColor(cmap []ColorStop, val float64) (uint8, uint8, uint8) {
	if len(cmap) == 0 {
		return 0, 0, 0
	}
	if val <= cmap[0].Value {
		return cmap[0].R, cmap[0].G, cmap[0].B
	}
	if val >= cmap[len(cmap)-1].Value {
		return cmap[len(cmap)-1].R, cmap[len(cmap)-1].G, cmap[len(cmap)-1].B
	}
	for i := 0; i < len(cmap)-1; i++ {
		if val >= cmap[i].Value && val <= cmap[i+1].Value {
			f := (val - cmap[i].Value) / (cmap[i+1].Value - cmap[i].Value)
			r := uint8(float64(cmap[i].R) + f*float64(int(cmap[i+1].R)-int(cmap[i].R)))
			g := uint8(float64(cmap[i].G) + f*float64(int(cmap[i+1].G)-int(cmap[i].G)))
			b := uint8(float64(cmap[i].B) + f*float64(int(cmap[i+1].B)-int(cmap[i].B)))
			return r, g, b
		}
	}
	return cmap[len(cmap)-1].R, cmap[len(cmap)-1].G, cmap[len(cmap)-1].B
}
