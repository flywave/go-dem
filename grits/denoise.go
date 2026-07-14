package grits

import (
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
	return applyBilateralFilter(data, region.XSize, region.YSize, opts)
}

func (f *denoiseFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	method := opts.Method
	if method == "" {
		method = "median"
	}

	switch method {
	case "bilateral":
		return applyBilateralFilter(data, region.XSize, region.YSize, opts)
	default:
		return medianFilter2D(data, region.XSize, region.YSize, opts.KernelSize, opts.GetNoData()), nil
	}
}

func applyBilateralFilter(data []float64, w, h int, opts *Options) ([]float64, error) {
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
	}

	return result, nil
}
