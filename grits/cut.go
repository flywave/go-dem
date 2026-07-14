package grits

import (
	"github.com/flywave/go-dem"
)

type cutFilter struct {
	baseGrits
}

func init() {
	Register(FilterCut, func() Grits { return &cutFilter{baseGrits{name: string(FilterCut)}} })
}

func (f *cutFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	bounds := opts.CutBounds
	if bounds[0] == 0 && bounds[1] == 0 && bounds[2] == 0 && bounds[3] == 0 {
		return data, nil
	}

	noData := opts.GetNoData()
	w, h := region.XSize, region.YSize
	result := make([]float64, len(data))
	copy(result, data)

	xMin, yMin, xMax, yMax := bounds[0], bounds[1], bounds[2], bounds[3]
	gt := region.GeoTransform()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			geoX := gt[0] + float64(x)*gt[1]
			geoY := gt[3] + float64(y)*gt[5]

			inside := geoX >= xMin && geoX <= xMax && geoY >= yMin && geoY <= yMax

			mask := (!opts.CutInvert && !inside) || (opts.CutInvert && inside)
			if mask {
				result[y*w+x] = noData
			}
		}
	}

	return result, nil
}
