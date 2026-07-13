package grits

import (
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

	return computeWeightBuffer(data, region.XSize, region.YSize, radius, noData), nil
}

func computeWeightBuffer(data []float64, w, h, radius int, noData float64) []float64 {
	weights := make([]float64, w*h)

	distMap := make([][]float64, h)
	for y := 0; y < h; y++ {
		distMap[y] = make([]float64, w)
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] != noData && !math.IsNaN(data[idx]) {
				distMap[y][x] = 0
			} else {
				distMap[y][x] = float64(w + h)
			}
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x > 0 && distMap[y][x-1]+1 < distMap[y][x] {
				distMap[y][x] = distMap[y][x-1] + 1
			}
			if y > 0 && distMap[y-1][x]+1 < distMap[y][x] {
				distMap[y][x] = distMap[y-1][x] + 1
			}
		}
	}

	for y := h - 1; y >= 0; y-- {
		for x := w - 1; x >= 0; x-- {
			if x < w-1 && distMap[y][x+1]+1 < distMap[y][x] {
				distMap[y][x] = distMap[y][x+1] + 1
			}
			if y < h-1 && distMap[y+1][x]+1 < distMap[y][x] {
				distMap[y][x] = distMap[y+1][x] + 1
			}
		}
	}

	maxDist := 0.0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if distMap[y][x] > maxDist && distMap[y][x] < float64(w+h) {
				maxDist = distMap[y][x]
			}
		}
	}

	if maxDist <= 0 {
		hasValid := false
		for _, v := range data {
			if v != noData && !math.IsNaN(v) {
				hasValid = true
				break
			}
		}
		if !hasValid {
			return weights
		}
		for i := range weights {
			weights[i] = 1.0
		}
		return weights
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] != noData && !math.IsNaN(data[idx]) {
				dist := distMap[y][x]
				if float64(radius) > 0 && dist <= float64(radius) {
					weights[idx] = 1.0 - dist/float64(radius)
					if weights[idx] < 0.1 {
						weights[idx] = 0.1
					}
				} else {
					weights[idx] = 1.0
				}
			} else {
				weights[idx] = 0
			}
		}
	}

	return weights
}
