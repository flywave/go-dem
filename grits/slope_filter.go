package grits

import (
	"math"

	"github.com/flywave/go-dem"
)

type slopeFilter struct{ baseGrits }

func init() {
	Register("slope_filter", func() Grits { return &slopeFilter{baseGrits{name: "slope_filter"}} })
}

func (f *slopeFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	maxSlope := opts.Threshold
	if maxSlope <= 0 {
		maxSlope = 30.0
	}
	noData := opts.GetNoData()

	slope := computeSlopeDegrees(data, region.XSize, region.YSize, region.XRes, noData)

	result := make([]float64, len(data))
	copy(result, data)

	w, h := region.XSize, region.YSize
	for i := 0; i < w*h; i++ {
		if data[i] == noData || math.IsNaN(data[i]) {
			result[i] = noData
			continue
		}
		if slope[i] > maxSlope {
			result[i] = noData
		}
	}

	return result, nil
}

func computeSlopeDegrees(data []float64, w, h int, res float64, noData float64) []float64 {
	slope := make([]float64, w*h)

	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			idx := y*w + x
			z := data[idx]
			if z == noData || math.IsNaN(z) {
				slope[idx] = noData
				continue
			}

			zW := data[y*w+(x-1)]
			zE := data[y*w+(x+1)]
			zN := data[(y-1)*w+x]
			zS := data[(y+1)*w+x]

			if dem.IsNoData(zW, noData) || dem.IsNoData(zE, noData) || dem.IsNoData(zN, noData) || dem.IsNoData(zS, noData) {
				slope[idx] = 0
				continue
			}

			dzdx := (zE - zW) / (2 * res)
			dzdy := (zS - zN) / (2 * res)

			slopeRad := math.Atan(math.Sqrt(dzdx*dzdx + dzdy*dzdy))
			slope[idx] = slopeRad * 180 / math.Pi
		}
	}

	return slope
}
