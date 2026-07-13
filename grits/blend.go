package grits

import (
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

	if maskRegion.XSize != region.XSize || maskRegion.YSize != region.YSize {
		return linearBlendResampled(data, region, maskData, maskRegion, noData, blendWidth), nil
	}

	return linearBlend(data, maskData, region, noData, blendWidth), nil
}

func linearBlend(dem, mask []float64, region *dem.Region, noData, blendWidth float64) []float64 {
	w, h := region.XSize, region.YSize
	result := make([]float64, w*h)
	copy(result, dem)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			demVal := dem[idx]
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

			dist := edgeDistance(idx, mask, w, h, noData)
			if dist < 0 {
				continue
			}
			if dist >= blendWidth {
				continue
			}

			t := dist / blendWidth
			result[idx] = demVal*(1-t) + maskVal*t
		}
	}

	return result
}

func linearBlendResampled(dem []float64, region *dem.Region, mask []float64, maskRegion *dem.Region, noData, blendWidth float64) []float64 {
	w, h := region.XSize, region.YSize
	result := make([]float64, w*h)
	copy(result, dem)

	gt := region.GeoTransform()
	mgt := maskRegion.GeoTransform()
	mw, mh := maskRegion.XSize, maskRegion.YSize

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			demVal := dem[idx]
			if demVal == noData || math.IsNaN(demVal) {
				continue
			}

			geoX := gt[0] + float64(x)*gt[1]
			geoY := gt[3] + float64(y)*gt[5]

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

			dist := edgeDistance(my*mw+mx, mask, mw, mh, mask[my*mw+mx])
			if dist < 0 || dist >= blendWidth {
				continue
			}
			t := dist / blendWidth
			result[idx] = demVal*(1-t) + maskVal*t
		}
	}

	return result
}

func edgeDistance(idx int, data []float64, w, h int, noData float64) float64 {
	x, y := idx%w, idx/w
	if data[idx] == noData || math.IsNaN(data[idx]) {
		return -1
	}

	minDist := -1.0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx < 0 || nx >= w || ny < 0 || ny >= h {
				continue
			}
			nidx := ny*w + nx
			if nidx < 0 || nidx >= len(data) {
				continue
			}
			if data[nidx] == noData || math.IsNaN(data[nidx]) {
				dist := math.Sqrt(float64(dx*dx + dy*dy))
				if minDist < 0 || dist < minDist {
					minDist = dist
				}
			}
		}
	}
	return minDist
}
