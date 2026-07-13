package grits

import (
	"math"

	"github.com/flywave/go-dem"
)

type riverFilter struct{ baseGrits }

func init() {
	Register("rivers", func() Grits { return &riverFilter{baseGrits{name: "rivers"}} })
}

func (f *riverFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	threshold := opts.Threshold
	if threshold <= 0 {
		threshold = 0.01
	}
	noData := opts.GetNoData()

	flowAcc := computeFlowAccumulationIterative(data, region.XSize, region.YSize, noData)
	rivers := thresholdRiverNetwork(flowAcc, region.XSize, region.YSize, threshold, noData)

	return rivers, nil
}

func computeFlowDirection(data []float64, w, h int, noData float64) []int {
	dir := make([]int, w*h)
	for i := range dir {
		dir[i] = -1
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] == noData || math.IsNaN(data[idx]) {
				continue
			}

			minElev := data[idx]
			minDir := -1

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
					nval := data[nidx]
					if nval == noData || math.IsNaN(nval) {
						continue
					}
					if nval < minElev {
						minElev = nval
						minDir = ny*w + nx
					}
				}
			}

			if minDir >= 0 && minDir != idx {
				dir[idx] = minDir
			}
		}
	}

	return dir
}

func computeFlowAccumulationIterative(data []float64, w, h int, noData float64) []float64 {
	flowDir := computeFlowDirection(data, w, h, noData)
	acc := make([]float64, w*h)
	indeg := make([]int, w*h)

	for i := 0; i < w*h; i++ {
		if flowDir[i] >= 0 {
			indeg[flowDir[i]]++
		}
	}

	queue := make([]int, 0, w*h)
	for i := 0; i < w*h; i++ {
		if data[i] != noData && !math.IsNaN(data[i]) && indeg[i] == 0 {
			acc[i] = 1
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		idx := queue[0]
		queue = queue[1:]

		down := flowDir[idx]
		if down >= 0 {
			acc[down] += acc[idx]
			indeg[down]--
			if indeg[down] == 0 {
				queue = append(queue, down)
			}
		}
	}

	return acc
}

func thresholdRiverNetwork(flowAcc []float64, w, h int, threshold float64, noData float64) []float64 {
	rivers := make([]float64, w*h)
	maxAcc := 0.0

	for i := range flowAcc {
		if flowAcc[i] > maxAcc {
			maxAcc = flowAcc[i]
		}
	}

	if maxAcc == 0 {
		return rivers
	}

	for i := range flowAcc {
		normalized := flowAcc[i] / maxAcc
		if normalized >= threshold {
			rivers[i] = normalized
		} else {
			rivers[i] = noData
		}
	}

	return rivers
}
