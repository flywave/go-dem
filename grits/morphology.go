package grits

import (
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
	return erode(data, region.XSize, region.YSize, radius, opts.GetNoData())
}

func (f *dilateFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	radius := opts.Radius
	if radius <= 0 {
		radius = 1
	}
	return dilate(data, region.XSize, region.YSize, radius, opts.GetNoData())
}

func (f *openFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	radius := opts.Radius
	if radius <= 0 {
		radius = 1
	}
	noData := opts.GetNoData()
	result, _ := erode(data, region.XSize, region.YSize, radius, noData)
	dilatedResult, _ := dilate(result, region.XSize, region.YSize, radius, noData)
	return dilatedResult, nil
}

func (f *closeFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	radius := opts.Radius
	if radius <= 0 {
		radius = 1
	}
	noData := opts.GetNoData()
	result, _ := dilate(data, region.XSize, region.YSize, radius, noData)
	erodedResult, _ := erode(result, region.XSize, region.YSize, radius, noData)
	return erodedResult, nil
}

func erode(data []float64, w, h, radius int, noData float64) ([]float64, error) {
	result := make([]float64, w*h)
	copy(result, data)

	for y := 0; y < h; y++ {
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
	}

	return result, nil
}

func dilate(data []float64, w, h, radius int, noData float64) ([]float64, error) {
	result := make([]float64, w*h)
	copy(result, data)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x

			hasValid := false
			var sum, count float64

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
					hasValid = true
					sum += nval
					count++
				}
			}

			if hasValid {
				result[idx] = sum / count
			} else {
				result[idx] = noData
			}
		}
	}

	return result, nil
}
