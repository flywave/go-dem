package dem

import "math"

const (
	DefaultNoData float64 = -9999
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


type InterpMethod string

const (
	MethodIDW             InterpMethod = "idw"
	MethodKriging         InterpMethod = "kriging"
	MethodLinear          InterpMethod = "linear"
	MethodCubic           InterpMethod = "cubic"
	MethodNearest         InterpMethod = "nearest"
	MethodNaturalNeighbor InterpMethod = "natural_neighbor"
	MethodCUDEM           InterpMethod = "cudem"
)

type OutputFormat string

const (
	FormatGeoTIFF OutputFormat = "GTiff"
	FormatCOG     OutputFormat = "COG"
	FormatNetCDF  OutputFormat = "NetCDF"
)

type ElevationLimit struct {
	Upper *float64
	Lower *float64
}

type WaffleOptions struct {
	Method         InterpMethod
	Power          float64
	MinPoints      int
	SearchRadius   float64
	NoData         float64
	Limits         *ElevationLimit
	WantUncertainty bool
	WantMask       bool
	WantStack      bool
	ChunkSize      *[2]int
}

type DEMResult struct {
	Path        string
	StackPath   string
	MaskPath    string
	UncPath     string
}

func DefaultWaffleOptions() WaffleOptions {
	return WaffleOptions{
		Method:         MethodIDW,
		Power:          2.0,
		MinPoints:      3,
		SearchRadius:   0,
		NoData:         DefaultNoData,
		WantUncertainty: false,
		WantMask:       false,
		WantStack:      true,
	}
}
