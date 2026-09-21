package grits

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
)

type blendFilter struct{ baseGrits }

func init() {
	Register(FilterBlend, func() Grits { return &blendFilter{baseGrits{name: string(FilterBlend)}} })
}

func (f *blendFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	maskPath := opts.SourceMask
	if maskPath == "" {
		return data, nil
	}

	maskData, maskRegion, err := dem.ReadDEM(maskPath)
	if err != nil {
		return nil, err
	}

	noData := opts.GetNoData()
	blendWidth := opts.MaxDistance
	if blendWidth <= 0 {
		blendWidth = float64(region.XRes * 10)
	}

	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	var result []float64
	if maskRegion.XSize != region.XSize || maskRegion.YSize != region.YSize {
		result, err = linearBlendResampledOpts(data, region, maskData, maskRegion, noData, blendWidth, opts, f.Name())
	} else {
		result, err = linearBlendOpts(data, maskData, region, noData, blendWidth, opts, f.Name())
	}
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return result, nil
}

func linearBlend(demData, mask []float64, region *dem.Region, noData, blendWidth float64) []float64 {
	result, _ := linearBlendOpts(demData, mask, region, noData, blendWidth, nil, "")
	return result
}

func linearBlendOpts(demData, mask []float64, region *dem.Region, noData, blendWidth float64, opts *Options, stage string) ([]float64, error) {
	prog, ctx := progressOf(opts)
	w, h := region.XSize, region.YSize
	result := make([]float64, w*h)
	copy(result, demData)

	distMap := dem.ComputeEuclideanDistance(mask, w, h, noData, region.XRes, region.YRes)

	for y := 0; y < h; y++ {
		if err := dem.CheckCtx(ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := 0; x < w; x++ {
			idx := y*w + x
			demVal := demData[idx]
			maskVal := mask[idx]

			if demVal == noData || math.IsNaN(demVal) {
				result[idx] = maskVal
				continue
			}
			if maskVal == noData || math.IsNaN(maskVal) {
				continue
			}
			if math.Abs(demVal-maskVal) < 1e-10 {
				continue
			}

			dist := distMap[idx]
			if dist >= blendWidth {
				continue
			}

			t := dist / blendWidth
			result[idx] = demVal*(1-t) + maskVal*t
		}
		dem.ReportProgress(prog, stage, y+1, h)
	}

	return result, nil
}

func linearBlendResampled(demData []float64, region *dem.Region, mask []float64, maskRegion *dem.Region, noData, blendWidth float64) []float64 {
	result, _ := linearBlendResampledOpts(demData, region, mask, maskRegion, noData, blendWidth, nil, "")
	return result
}

func linearBlendResampledOpts(demData []float64, region *dem.Region, mask []float64, maskRegion *dem.Region, noData, blendWidth float64, opts *Options, stage string) ([]float64, error) {
	prog, ctx := progressOf(opts)
	w, h := region.XSize, region.YSize
	result := make([]float64, w*h)
	copy(result, demData)

	mgt := maskRegion.GeoTransform()
	mw, mh := maskRegion.XSize, maskRegion.YSize

	distMap := dem.ComputeEuclideanDistance(mask, mw, mh, noData, maskRegion.XRes, maskRegion.YRes)

	for y := 0; y < h; y++ {
		if err := dem.CheckCtx(ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := 0; x < w; x++ {
			idx := y*w + x
			demVal := demData[idx]
			if demVal == noData || math.IsNaN(demVal) {
				continue
			}

			geoX, geoY := region.PixelCenterGeo(x, y)

			mx := int((geoX - mgt[0]) / mgt[1])
			my := int((geoY - mgt[3]) / mgt[5])
			if mx < 0 || mx >= mw || my < 0 || my >= mh {
				continue
			}
			maskVal := mask[my*mw+mx]
			if maskVal == noData || math.IsNaN(maskVal) {
				continue
			}

			if math.Abs(demVal-maskVal) < 1e-10 {
				continue
			}

			dist := distMap[my*mw+mx]
			if dist >= blendWidth {
				continue
			}
			t := dist / blendWidth
			result[idx] = demVal*(1-t) + maskVal*t
		}
		dem.ReportProgress(prog, stage, y+1, h)
	}

	return result, nil
}
