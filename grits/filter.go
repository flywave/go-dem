package grits

import (
	"fmt"
	"math"
	"sort"

	"github.com/flywave/go-dem"
)

type gaussianFilter struct{ baseGrits }
type medianFilter struct{ baseGrits }
type blurFilter struct{ baseGrits }

func init() {
	Register(FilterGaussian, func() Grits { return &gaussianFilter{baseGrits{name: string(FilterGaussian)}} })
	Register(FilterMedian, func() Grits { return &medianFilter{baseGrits{name: string(FilterMedian)}} })
	Register(FilterBlur, func() Grits { return &blurFilter{baseGrits{name: string(FilterBlur)}} })
}

func (f *blurFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	sigma := opts.Sigma
	if sigma <= 0 {
		sigma = 1.0
	}
	radius := opts.Radius
	if radius <= 0 {
		radius = int(math.Ceil(sigma * 2))
	}
	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	result, err := gaussianBlur(data, region.XSize, region.YSize, sigma, radius, opts.GetNoData(), opts, f.Name())
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return result, nil
}

func (f *gaussianFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	sigma := opts.Sigma
	if sigma <= 0 {
		sigma = 1.0
	}
	radius := opts.Radius
	if radius <= 0 {
		radius = int(math.Ceil(sigma * 2))
	}
	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	result, err := gaussianBlur(data, region.XSize, region.YSize, sigma, radius, opts.GetNoData(), opts, f.Name())
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return result, nil
}

func (f *medianFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	kSize := normalizeKernelSize(opts.KernelSize)
	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	result, err := medianFilter2D(data, region.XSize, region.YSize, kSize, opts.GetNoData(), opts, f.Name())
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return result, nil
}

func normalizeKernelSize(kSize int) int {
	if kSize < 3 {
		kSize = 3
	}
	if kSize%2 == 0 {
		kSize++
	}
	return kSize
}

func gaussianBlur(data []float64, w, h int, sigma float64, radius int, noData float64, opts *Options, stage string) ([]float64, error) {
	kernel := make1DGaussianKernel(sigma, radius)
	result := make([]float64, len(data))
	copy(result, data)

	scratch := make([]float64, len(data))

	for y := 0; y < h; y++ {
		if err := dem.CheckCtx(opts.Ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] == noData || math.IsNaN(data[idx]) {
				scratch[idx] = data[idx]
				continue
			}
			var sum, weightSum float64
			for k := -radius; k <= radius; k++ {
				ix := x + k
				if ix < 0 || ix >= w {
					continue
				}
				val := data[y*w+ix]
				if val == noData || math.IsNaN(val) {
					continue
				}
				kw := kernel[k+radius]
				sum += val * kw
				weightSum += kw
			}
			if weightSum > 0 {
				scratch[y*w+x] = sum / weightSum
			} else {
				scratch[y*w+x] = data[y*w+x]
			}
		}
		dem.ReportProgress(opts.Progress, stage, y+1, h)
	}

	for x := 0; x < w; x++ {
		if err := dem.CheckCtx(opts.Ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for y := 0; y < h; y++ {
			if scratch[y*w+x] == noData || math.IsNaN(scratch[y*w+x]) {
				continue
			}
			var sum, weightSum float64
			for k := -radius; k <= radius; k++ {
				iy := y + k
				if iy < 0 || iy >= h {
					continue
				}
				val := scratch[iy*w+x]
				if val == noData || math.IsNaN(val) {
					continue
				}
				kw := kernel[k+radius]
				sum += val * kw
				weightSum += kw
			}
			if weightSum > 0 {
				result[y*w+x] = sum / weightSum
			} else {
				result[y*w+x] = scratch[y*w+x]
			}
		}
	}

	return result, nil
}

func make1DGaussianKernel(sigma float64, radius int) []float64 {
	kernel := make([]float64, 2*radius+1)
	var sum float64
	for i := -radius; i <= radius; i++ {
		val := math.Exp(-float64(i*i) / (2 * sigma * sigma))
		kernel[i+radius] = val
		sum += val
	}
	for i := range kernel {
		kernel[i] /= sum
	}
	return kernel
}

func medianFilter2D(data []float64, w, h, kSize int, noData float64, opts *Options, stage string) ([]float64, error) {
	result := make([]float64, len(data))
	copy(result, data)

	half := kSize / 2
	neighbors := make([]float64, 0, kSize*kSize)

	for y := 0; y < h; y++ {
		if err := dem.CheckCtx(opts.Ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] == noData || math.IsNaN(data[idx]) {
				continue
			}

			neighbors = neighbors[:0]
			for ky := -half; ky <= half; ky++ {
				for kx := -half; kx <= half; kx++ {
					ix, iy := x+kx, y+ky
					if ix < 0 || ix >= w || iy < 0 || iy >= h {
						continue
					}
					val := data[iy*w+ix]
					if val == noData || math.IsNaN(val) {
						continue
					}
					neighbors = append(neighbors, val)
				}
			}

			if len(neighbors) > 0 {
				result[idx] = median(neighbors)
			}
		}
		dem.ReportProgress(opts.Progress, stage, y+1, h)
	}

	return result, nil
}

func median(vals []float64) float64 {
	n := len(vals)
	if n == 0 {
		return 0
	}
	sorted := make([]float64, n)
	copy(sorted, vals)
	sort.Float64s(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}
