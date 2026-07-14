package grits

import (
	"math"

	"github.com/flywave/go-dem"
)

type zscoreFilter struct {
	baseGrits
}

func init() {
	Register(FilterZScore, func() Grits { return &zscoreFilter{baseGrits{name: string(FilterZScore)}} })
}

func (f *zscoreFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	threshold := opts.Threshold
	if threshold <= 0 {
		threshold = 3.0
	}
	kSize := opts.KernelSize
	if kSize < 3 {
		kSize = 5
	}
	if kSize%2 == 0 {
		kSize++
	}

	w, h := region.XSize, region.YSize
	nd := opts.GetNoData()
	result := make([]float64, len(data))
	copy(result, data)

	half := kSize / 2

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if result[idx] == nd || math.IsNaN(result[idx]) {
				continue
			}

			var sum, sumSq, count float64
			for ky := -half; ky <= half; ky++ {
				for kx := -half; kx <= half; kx++ {
					ix, iy := x+kx, y+ky
					if ix < 0 || ix >= w || iy < 0 || iy >= h {
						continue
					}
					val := data[iy*w+ix]
					if val == nd || math.IsNaN(val) {
						continue
					}
					sum += val
					sumSq += val * val
					count++
				}
			}
			if count < 2 {
				continue
			}
			mean := sum / count
			variance := (sumSq / count) - (mean * mean)
			if variance < 1e-12 {
				continue
			}
			std := math.Sqrt(variance)
			z := math.Abs((result[idx] - mean) / std)
			if z > threshold {
				result[idx] = nd
			}
		}
	}

	return result, nil
}
