package grits

import (
	"fmt"

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
	if len(bounds) != 4 {
		return data, nil
	}

	noData := opts.GetNoData()
	w, h := region.XSize, region.YSize
	result := make([]float64, len(data))
	copy(result, data)

	xMin, yMin, xMax, yMax := bounds[0], bounds[1], bounds[2], bounds[3]

	if err := startFilter(opts, f.Name(), h); err != nil {
		return nil, err
	}
	for y := 0; y < h; y++ {
		if err := dem.CheckCtx(opts.Ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", f.Name(), err)
		}
		for x := 0; x < w; x++ {
			geoX, geoY := region.PixelCenterGeo(x, y)

			inside := geoX >= xMin && geoX <= xMax && geoY >= yMin && geoY <= yMax

			mask := (!opts.CutInvert && !inside) || (opts.CutInvert && inside)
			if mask {
				result[y*w+x] = noData
			}
		}
		dem.ReportProgress(opts.Progress, f.Name(), y+1, h)
	}

	finishFilter(opts, f.Name(), h)
	return result, nil
}
