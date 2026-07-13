package perspecto

type Options struct {
	NoData         float64
	ZFactor        float64
	Azimuth        float64
	Altitude       float64
	Colormap       []ColorStop
	SlopeUnits     string
}

type ColorStop struct {
	Value  float64
	R, G, B uint8
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
