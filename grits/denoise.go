package grits

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
)

type denoiseFilter struct {
	baseGrits
}

type bilateralFilter struct {
	baseGrits
}

func init() {
	Register(FilterDenoise, func() Grits { return &denoiseFilter{baseGrits{name: string(FilterDenoise)}} })
	Register(FilterBilateral, func() Grits { return &bilateralFilter{baseGrits{name: string(FilterBilateral)}} })
}

func (f *bilateralFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	result, err := applyBilateralFilter(data, region.XSize, region.YSize, opts, f.Name())
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return result, nil
}

func (f *denoiseFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	method := opts.Method
	if method == "" {
		method = "median"
	}

	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	var (
		result []float64
		err    error
	)
	switch method {
	case "bilateral":
		result, err = applyBilateralFilter(data, region.XSize, region.YSize, opts, f.Name())
	default:
		kSize := normalizeKernelSize(opts.KernelSize)
		result, err = medianFilter2D(data, region.XSize, region.YSize, kSize, opts.GetNoData(), opts, f.Name())
	}
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return result, nil
}

func applyBilateralFilter(data []float64, w, h int, opts *Options, stage string) ([]float64, error) {
	sigmaSpatial := opts.Sigma
	if sigmaSpatial <= 0 {
		sigmaSpatial = 1.0
	}
	sigmaColor := opts.SigmaColor
	if sigmaColor <= 0 {
		sigmaColor = 20.0
	}
	radius := opts.Radius
	if radius <= 0 {
		radius = int(math.Ceil(sigmaSpatial * 2))
	}
	nd := opts.GetNoData()

	result := make([]float64, len(data))
	copy(result, data)

	spatialKernel := make1DGaussianKernel(sigmaSpatial, radius)

	for y := 0; y < h; y++ {
		if err := dem.CheckCtx(opts.Ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := 0; x < w; x++ {
			idx := y*w + x
			centerVal := data[idx]
			if centerVal == nd || math.IsNaN(centerVal) {
				continue
			}

			var sumWeight, sumValue float64
			for ky := -radius; ky <= radius; ky++ {
				for kx := -radius; kx <= radius; kx++ {
					ix, iy := x+kx, y+ky
					if ix < 0 || ix >= w || iy < 0 || iy >= h {
						continue
					}
					val := data[iy*w+ix]
					if val == nd || math.IsNaN(val) {
						continue
					}

					spatialW := spatialKernel[kx+radius] * spatialKernel[ky+radius]
					diff := centerVal - val
					rangeW := math.Exp(-(diff * diff) / (2 * sigmaColor * sigmaColor))
					weight := spatialW * rangeW

					sumWeight += weight
					sumValue += val * weight
				}
			}
			if sumWeight > 0 {
				result[idx] = sumValue / sumWeight
			}
		}
		dem.ReportProgress(opts.Progress, stage, y+1, h)
	}

	return result, nil
}
