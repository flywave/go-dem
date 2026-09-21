package grits

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
)

type diffFilter struct{ baseGrits }

func init() {
	Register("diff", func() Grits { return &diffFilter{baseGrits{name: "diff"}} })
}

func (f *diffFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	if opts.SourceMask == "" {
		return nil, fmt.Errorf("reference DEM path required via SourceMask")
	}

	refData, refRegion, err := dem.ReadDEM(opts.SourceMask)
	if err != nil {
		return nil, fmt.Errorf("read reference DEM: %v", err)
	}

	noData := opts.GetNoData()

	if refRegion.SRS() != nil && region.SRS() != nil &&
		!refRegion.SRS().Eq(region.SRS()) {
		return nil, fmt.Errorf("SRS mismatch: reference is %s, target is %s",
			refRegion.SRS().GetSrsCode(), region.SRS().GetSrsCode())
	}

	width := region.XSize
	height := region.YSize

	if err := startFilter(opts, f.Name(), height); err != nil {
		return nil, err
	}

	if refRegion.XSize != width || refRegion.YSize != height {
		result, err := resampleDiffOpts(data, region, refData, refRegion, noData, opts.Threshold, opts, f.Name())
		if err != nil {
			return nil, err
		}
		finishFilter(opts, f.Name(), height)
		return result, nil
	}

	result := make([]float64, width*height)
	for y := 0; y < height; y++ {
		if err := dem.CheckCtx(opts.Ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", f.Name(), err)
		}
		for x := 0; x < width; x++ {
			idx := y*width + x
			if opts.IsNoData(data[idx]) || opts.IsNoData(refData[idx]) {
				result[idx] = noData
				continue
			}
			diff := data[idx] - refData[idx]
			if opts.Threshold > 0 && math.Abs(diff) < opts.Threshold {
				result[idx] = 0
			} else {
				result[idx] = diff
			}
		}
		dem.ReportProgress(opts.Progress, f.Name(), y+1, height)
	}

	finishFilter(opts, f.Name(), height)
	return result, nil
}

func ComputeDiff(dem1, dem2 []float64, region *dem.Region, absDiff bool, noData float64) []float64 {
	w, h := region.XSize, region.YSize
	result := make([]float64, w*h)

	for i := range result {
		if dem1[i] == noData || math.IsNaN(dem1[i]) || dem2[i] == noData || math.IsNaN(dem2[i]) {
			result[i] = noData
			continue
		}
		if absDiff {
			result[i] = math.Abs(dem1[i] - dem2[i])
		} else {
			result[i] = dem1[i] - dem2[i]
		}
	}

	return result
}

func resampleDiff(data []float64, region *dem.Region, refData []float64, refRegion *dem.Region, noData, threshold float64) ([]float64, error) {
	return resampleDiffOpts(data, region, refData, refRegion, noData, threshold, nil, "")
}

func resampleDiffOpts(data []float64, region *dem.Region, refData []float64, refRegion *dem.Region, noData, threshold float64, opts *Options, stage string) ([]float64, error) {
	prog, ctx := progressOf(opts)
	isNoData := func(v float64) bool { return v == noData || math.IsNaN(v) }
	width, height := region.XSize, region.YSize
	rw, rh := refRegion.XSize, refRegion.YSize

	result := make([]float64, width*height)
	rgt := refRegion.GeoTransform()

	for y := 0; y < height; y++ {
		if err := dem.CheckCtx(ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := 0; x < width; x++ {
			idx := y*width + x
			if isNoData(data[idx]) {
				result[idx] = noData
				continue
			}

			geoX, geoY := region.PixelCenterGeo(x, y)

			rx := int(math.Round((geoX-rgt[0])/rgt[1] - 0.5))
			ry := int(math.Round((geoY-rgt[3])/rgt[5] - 0.5))

			if rx < 0 || rx >= rw || ry < 0 || ry >= rh {
				result[idx] = noData
				continue
			}

			v2 := refData[ry*rw+rx]
			if isNoData(v2) {
				result[idx] = noData
				continue
			}

			diff := data[idx] - v2
			if threshold > 0 && math.Abs(diff) < threshold {
				result[idx] = 0
			} else {
				result[idx] = diff
			}
		}
		dem.ReportProgress(prog, stage, y+1, height)
	}

	return result, nil
}
