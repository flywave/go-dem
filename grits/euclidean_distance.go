package grits

import (
	"math"

	"github.com/flywave/go-dem"
)

type euclideanDistanceFilter struct {
	baseGrits
}

func init() {
	Register(FilterEuclideanDistance, func() Grits {
		return &euclideanDistanceFilter{baseGrits{name: string(FilterEuclideanDistance)}}
	})
}

func (f *euclideanDistanceFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	w, h := region.XSize, region.YSize
	nd := opts.GetNoData()
	res := region.XRes

	dist := make([]float64, len(data))
	inf := math.MaxFloat64
	for i := range dist {
		if data[i] == nd || math.IsNaN(data[i]) {
			dist[i] = 0
		} else {
			dist[i] = inf
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if dist[idx] == 0 {
				continue
			}
			if x > 0 {
				dist[idx] = math.Min(dist[idx], dist[y*w+(x-1)]+res)
			}
			if y > 0 {
				dist[idx] = math.Min(dist[idx], dist[(y-1)*w+x]+res)
			}
			if x > 0 && y > 0 {
				dist[idx] = math.Min(dist[idx], dist[(y-1)*w+(x-1)]+res*math.Sqrt2)
			}
			if x < w-1 && y > 0 {
				dist[idx] = math.Min(dist[idx], dist[(y-1)*w+(x+1)]+res*math.Sqrt2)
			}
		}
	}

	for y := h - 1; y >= 0; y-- {
		for x := w - 1; x >= 0; x-- {
			idx := y*w + x
			if dist[idx] == 0 {
				continue
			}
			if x < w-1 {
				dist[idx] = math.Min(dist[idx], dist[y*w+(x+1)]+res)
			}
			if y < h-1 {
				dist[idx] = math.Min(dist[idx], dist[(y+1)*w+x]+res)
			}
			if x < w-1 && y < h-1 {
				dist[idx] = math.Min(dist[idx], dist[(y+1)*w+(x+1)]+res*math.Sqrt2)
			}
			if x > 0 && y < h-1 {
				dist[idx] = math.Min(dist[idx], dist[(y+1)*w+(x-1)]+res*math.Sqrt2)
			}
		}
	}

	return dist, nil
}
