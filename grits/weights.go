package grits

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
)

type weightFilter struct{ baseGrits }

func init() {
	Register("weights", func() Grits { return &weightFilter{baseGrits{name: "weights"}} })
}

func (f *weightFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	radius := opts.Radius
	if radius <= 0 {
		radius = 5
	}

	noData := opts.GetNoData()

	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	weights, err := computeWeightBufferOpts(data, region.XSize, region.YSize, radius, noData, opts, f.Name())
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return weights, nil
}

func computeWeightBuffer(data []float64, w, h, radius int, noData float64) []float64 {
	weights, _ := computeWeightBufferOpts(data, w, h, radius, noData, nil, "")
	return weights
}

func computeWeightBufferOpts(data []float64, w, h, radius int, noData float64, opts *Options, stage string) ([]float64, error) {
	prog, ctx := progressOf(opts)
	weights := make([]float64, w*h)

	const far = 1 << 30
	distMap := make([]float64, w*h)
	hasNoData := false
	for i := range distMap {
		if data[i] == noData || math.IsNaN(data[i]) {
			distMap[i] = 0
			hasNoData = true
		} else {
			distMap[i] = far
		}
	}

	for y := 0; y < h; y++ {
		if err := dem.CheckCtx(ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := 0; x < w; x++ {
			idx := y*w + x
			if distMap[idx] == 0 {
				continue
			}
			if x > 0 && distMap[idx-1]+1 < distMap[idx] {
				distMap[idx] = distMap[idx-1] + 1
			}
			if y > 0 && distMap[idx-w]+1 < distMap[idx] {
				distMap[idx] = distMap[idx-w] + 1
			}
		}
		dem.ReportProgress(prog, stage, (y+1)/2, h)
	}

	for y := h - 1; y >= 0; y-- {
		if err := dem.CheckCtx(ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := w - 1; x >= 0; x-- {
			idx := y*w + x
			if distMap[idx] == 0 {
				continue
			}
			if x < w-1 && distMap[idx+1]+1 < distMap[idx] {
				distMap[idx] = distMap[idx+1] + 1
			}
			if y < h-1 && distMap[idx+w]+1 < distMap[idx] {
				distMap[idx] = distMap[idx+w] + 1
			}
		}
		dem.ReportProgress(prog, stage, (h+h-1-y+1)/2, h)
	}

	if !hasNoData {
		for i := range weights {
			weights[i] = 1.0
		}
		return weights, nil
	}

	for i := range weights {
		if distMap[i] == 0 {
			weights[i] = 0
			continue
		}
		wt := 1.0 - distMap[i]/float64(radius)
		if wt < 0.1 {
			wt = 0.1
		}
		weights[i] = wt
	}

	return weights, nil
}
