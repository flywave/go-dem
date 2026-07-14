package grits

import (
	"math"
	"sort"

	"github.com/flywave/go-dem"
)

type flattenNoDataFilter struct {
	baseGrits
}

func init() {
	Register(FilterFlattenNoData, func() Grits {
		return &flattenNoDataFilter{baseGrits{name: string(FilterFlattenNoData)}}
	})
}

func (f *flattenNoDataFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	nd := opts.GetNoData()
	w, h := region.XSize, region.YSize
	threshold := opts.Threshold
	if threshold <= 0 {
		threshold = 100
	}

	labeled := make([]int, len(data))
	for i := range labeled {
		labeled[i] = -1
	}

	type cell struct{ x, y int }
	nextLabel := 0
	sizes := map[int]int{}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] == nd || math.IsNaN(data[idx]) {
				continue
			}
			if labeled[idx] >= 0 {
				continue
			}

			queue := []cell{{x, y}}
			labeled[idx] = nextLabel
			sizes[nextLabel] = 1
			head := 0
			for head < len(queue) {
				cx, cy := queue[head].x, queue[head].y
				head++
				dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
				for _, d := range dirs {
					nx, ny := cx+d[0], cy+d[1]
					if nx < 0 || nx >= w || ny < 0 || ny >= h {
						continue
					}
					nidx := ny*w + nx
					if data[nidx] == nd || math.IsNaN(data[nidx]) {
						continue
					}
					if labeled[nidx] >= 0 {
						continue
					}
					labeled[nidx] = nextLabel
					sizes[nextLabel]++
					queue = append(queue, cell{nx, ny})
				}
			}
			nextLabel++
		}
	}

	smallLabels := map[int]bool{}
	for l, sz := range sizes {
		if sz < int(threshold) {
			smallLabels[l] = true
		}
	}
	if len(smallLabels) == 0 {
		return data, nil
	}

	result := make([]float64, len(data))
	copy(result, data)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] == nd || math.IsNaN(data[idx]) {
				continue
			}
			if !smallLabels[labeled[idx]] {
				continue
			}

			var vals []float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					nx, ny := x+kx, y+ky
					if nx < 0 || nx >= w || ny < 0 || ny >= h {
						continue
					}
					nidx := ny*w + nx
					if data[nidx] == nd || math.IsNaN(data[nidx]) {
						continue
					}
					if smallLabels[labeled[nidx]] {
						continue
					}
					vals = append(vals, data[nidx])
				}
			}
			if len(vals) > 0 {
				sort.Float64s(vals)
				pct5 := vals[int(float64(len(vals)-1)*0.05)]
				result[idx] = pct5
			}
		}
	}

	return result, nil
}
