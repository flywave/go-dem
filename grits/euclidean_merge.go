package grits

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
)

type euclideanMergeFilter struct {
	baseGrits
}

func init() {
	Register(FilterEuclideanMerge, func() Grits {
		return &euclideanMergeFilter{baseGrits{name: string(FilterEuclideanMerge)}}
	})
}

func (f *euclideanMergeFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	if opts == nil || opts.SourceMask == "" {
		return data, nil
	}

	otherData, otherRegion, err := dem.ReadDEM(opts.SourceMask)
	if err != nil {
		return nil, fmt.Errorf("read source mask: %v", err)
	}

	if otherRegion.XSize != region.XSize || otherRegion.YSize != region.YSize {
		return nil, fmt.Errorf("region size mismatch: base %dx%d, other %dx%d",
			region.XSize, region.YSize, otherRegion.XSize, otherRegion.YSize)
	}

	nd := opts.GetNoData()
	w, h := region.XSize, region.YSize
	res := region.XRes

	smallDist := 0.001953125

	baseDist := computeEuclideanDistance(data, w, h, nd, res)
	otherDist := computeEuclideanDistance(otherData, w, h, nd, res)

	for i := range otherDist {
		if otherDist[i] == 0 {
			otherDist[i] = smallDist
		}
	}
	for i := range baseDist {
		if baseDist[i] == 0 {
			baseDist[i] = smallDist
		}
	}

	result := make([]float64, len(data))
	totalWeight := make([]float64, len(data))
	distanceSum := make([]float64, len(data))

	for i := range result {
		if data[i] != nd && !math.IsNaN(data[i]) {
			w := baseDist[i]
			result[i] += data[i] * w
			totalWeight[i] += w
			distanceSum[i] += w
		}
		if otherData[i] != nd && !math.IsNaN(otherData[i]) {
			w := otherDist[i]
			result[i] += otherData[i] * w
			totalWeight[i] += w
			distanceSum[i] += w
		}
	}

	for i := range result {
		if totalWeight[i] > 0 {
			if distanceSum[i] > smallDist*2 {
				nearestIdx := i
				if baseDist[i] > otherDist[i] {
					nearestIdx = i
				}
				_ = nearestIdx
			}
			result[i] /= totalWeight[i]
		} else {
			result[i] = nd
		}
	}

	return result, nil
}

func computeEuclideanDistance(data []float64, w, h int, nd float64, res float64) []float64 {
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

	return dist
}
