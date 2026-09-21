package grits

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
)

type erodeFilter struct{ baseGrits }
type dilateFilter struct{ baseGrits }
type openFilter struct{ baseGrits }
type closeFilter struct{ baseGrits }

func init() {
	Register(FilterErode, func() Grits { return &erodeFilter{baseGrits{name: string(FilterErode)}} })
	Register(FilterDilate, func() Grits { return &dilateFilter{baseGrits{name: string(FilterDilate)}} })
	Register(FilterOpen, func() Grits { return &openFilter{baseGrits{name: string(FilterOpen)}} })
	Register(FilterClose, func() Grits { return &closeFilter{baseGrits{name: string(FilterClose)}} })
}

func (f *erodeFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	radius := opts.Radius
	if radius <= 0 {
		radius = 1
	}
	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	result, err := erode(data, region.XSize, region.YSize, radius, opts.GetNoData(), opts, f.Name(), 0, region.YSize)
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return result, nil
}

func (f *dilateFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	radius := opts.Radius
	if radius <= 0 {
		radius = 1
	}
	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	result, err := dilate(data, region.XSize, region.YSize, radius, opts.GetNoData(), opts, f.Name(), 0, region.YSize)
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return result, nil
}

func (f *openFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	radius := opts.Radius
	if radius <= 0 {
		radius = 1
	}
	noData := opts.GetNoData()
	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	result, err := erode(data, region.XSize, region.YSize, radius, noData, opts, f.Name(), 0, region.YSize/2)
	if err != nil {
		return nil, err
	}
	if err := dem.CheckCtx(opts.Ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", f.Name(), err)
	}
	dilatedResult, err := dilate(result, region.XSize, region.YSize, radius, noData, opts, f.Name(), region.YSize/2, region.YSize-region.YSize/2)
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return dilatedResult, nil
}

func (f *closeFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	radius := opts.Radius
	if radius <= 0 {
		radius = 1
	}
	noData := opts.GetNoData()
	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	result, err := dilate(data, region.XSize, region.YSize, radius, noData, opts, f.Name(), 0, region.YSize/2)
	if err != nil {
		return nil, err
	}
	if err := dem.CheckCtx(opts.Ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", f.Name(), err)
	}
	erodedResult, err := erode(result, region.XSize, region.YSize, radius, noData, opts, f.Name(), region.YSize/2, region.YSize-region.YSize/2)
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return erodedResult, nil
}

func erode(data []float64, w, h, radius int, noData float64, opts *Options, stage string, doneBase, doneSpan int) ([]float64, error) {
	result := make([]float64, w*h)
	copy(result, data)

	for y := 0; y < h; y++ {
		if err := dem.CheckCtx(opts.Ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] == noData || math.IsNaN(data[idx]) {
				result[idx] = noData
				continue
			}

			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					nx, ny := x+dx, y+dy
					if nx < 0 || nx >= w || ny < 0 || ny >= h {
						result[idx] = noData
						goto nextPixel
					}
					nidx := ny*w + nx
					if data[nidx] == noData || math.IsNaN(data[nidx]) {
						result[idx] = noData
						goto nextPixel
					}
				}
			}
		nextPixel:
		}
		dem.ReportProgress(opts.Progress, stage, doneBase+(y+1)*doneSpan/h, h)
	}

	return result, nil
}

func dilate(data []float64, w, h, radius int, noData float64, opts *Options, stage string, doneBase, doneSpan int) ([]float64, error) {
	result := make([]float64, w*h)

	for y := 0; y < h; y++ {
		if err := dem.CheckCtx(opts.Ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := 0; x < w; x++ {
			idx := y*w + x

			maxVal := noData
			hasValid := false

			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					nx, ny := x+dx, y+dy
					if nx < 0 || nx >= w || ny < 0 || ny >= h {
						continue
					}
					nidx := ny*w + nx
					nval := data[nidx]
					if nval == noData || math.IsNaN(nval) {
						continue
					}
					if !hasValid || nval > maxVal {
						maxVal = nval
					}
					hasValid = true
				}
			}

			if hasValid {
				result[idx] = maxVal
			} else {
				result[idx] = noData
			}
		}
		dem.ReportProgress(opts.Progress, stage, doneBase+(y+1)*doneSpan/h, h)
	}

	return result, nil
}
