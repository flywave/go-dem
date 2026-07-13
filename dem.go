package dem

import "math"

const (
	DefaultNoData float64 = -9999
)

type InterpMethod string

const (
	MethodIDW             InterpMethod = "idw"
	MethodKriging         InterpMethod = "kriging"
	MethodLinear          InterpMethod = "linear"
	MethodCubic           InterpMethod = "cubic"
	MethodNearest         InterpMethod = "nearest"
	MethodNaturalNeighbor InterpMethod = "natural_neighbor"
	MethodCUDEM           InterpMethod = "cudem"
	MethodInpaint         InterpMethod = "inpaint"
	MethodCUBE            InterpMethod = "cube"
)

type FilterType string

const (
	FilterGaussian   FilterType = "gaussian"
	FilterMedian     FilterType = "median"
	FilterClip       FilterType = "clip"
	FilterFill       FilterType = "fill"
	FilterBlend      FilterType = "blend"
	FilterErode      FilterType = "erode"
	FilterDilate     FilterType = "dilate"
	FilterOpen       FilterType = "open"
	FilterClose      FilterType = "close"
	FilterHydro      FilterType = "hydro"
	FilterDiff       FilterType = "diff"
	FilterWeights    FilterType = "weights"
	FilterRivers     FilterType = "rivers"
	FilterSlope      FilterType = "slope_filter"
)

func IsNoData(val, noData float64) bool {
	return val == noData || math.IsNaN(val)
}

func IsNoDataValue(val float64) bool {
	return val == DefaultNoData || math.IsNaN(val)
}

func CoalesceNoData(userVal *float64) float64 {
	if userVal != nil {
		return *userVal
	}
	return DefaultNoData
}
