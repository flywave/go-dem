package grits

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/flywave/go-dem"
)

type hydroFilter struct{ baseGrits }

func init() {
	Register("hydro", func() Grits { return &hydroFilter{baseGrits{name: "hydro"}} })
}

func (f *hydroFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	noData := opts.GetNoData()

	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	result, err := fillSinksOpts(data, region.XSize, region.YSize, noData, opts, f.Name())
	if err != nil {
		return nil, err
	}
	finishFilter(opts, f.Name(), region.YSize)
	return result, nil
}

type hydrologyEdge struct {
	x, y int
	z    float64
}

func fillSinks(data []float64, w, h int, noData float64) []float64 {
	result, _ := fillSinksOpts(data, w, h, noData, nil, "")
	return result
}

func fillSinksOpts(data []float64, w, h int, noData float64, opts *Options, stage string) ([]float64, error) {
	prog, ctx := progressOf(opts)
	result := make([]float64, len(data))
	copy(result, data)

	rowsPerIter := h - 2
	if rowsPerIter < 0 {
		rowsPerIter = 0
	}
	totalRows := 10 * rowsPerIter
	rowsDone := 0

	for iteration := 0; iteration < 10; iteration++ {
		changed := 0
		for y := 1; y < h-1; y++ {
			if err := dem.CheckCtx(ctx); err != nil {
				return nil, fmt.Errorf("%s: %w", stage, err)
			}
			for x := 1; x < w-1; x++ {
				idx := y*w + x
				z := result[idx]
				if z == noData || math.IsNaN(z) {
					continue
				}

				minNeighbor := math.MaxFloat64
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						if dx == 0 && dy == 0 {
							continue
						}
						nx, ny := x+dx, y+dy
						nidx := ny*w + nx
						nval := result[nidx]
						if nval != noData && !math.IsNaN(nval) && nval < minNeighbor {
							minNeighbor = nval
						}
					}
				}

				if minNeighbor < math.MaxFloat64 && z < minNeighbor-1e-10 {

					flowTo := -1
					minSlope := math.Inf(-1)
					for dy := -1; dy <= 1; dy++ {
						for dx := -1; dx <= 1; dx++ {
							if dx == 0 && dy == 0 {
								continue
							}
							nx, ny := x+dx, y+dy
							nidx := ny*w + nx
							nval := result[nidx]
							if nval == noData || math.IsNaN(nval) {
								continue
							}
							dist := math.Sqrt(float64(dx*dx + dy*dy))
							slope := (z - nval) / dist
							if slope > minSlope {
								minSlope = slope
								flowTo = nidx
							}
						}
					}

					if flowTo >= 0 {
						result[idx] = result[flowTo] + 0.0001
						changed++
					}
				}
			}
			rowsDone++
			if totalRows > 0 {
				dem.ReportProgress(prog, stage, rowsDone*h/totalRows, h)
			}
		}
		if changed == 0 {
			break
		}
	}

	result = fillFlatAreas(result, w, h, noData)

	return result, nil
}

type floodItem struct {
	z   float64
	idx int
}

type floodQueue []floodItem

func (q floodQueue) Len() int { return len(q) }
func (q floodQueue) Less(i, j int) bool {
	return q[i].z < q[j].z || (q[i].z == q[j].z && q[i].idx < q[j].idx)
}
func (q floodQueue) Swap(i, j int)       { q[i], q[j] = q[j], q[i] }
func (q *floodQueue) Push(x interface{}) { *q = append(*q, x.(floodItem)) }
func (q *floodQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	*q = old[:n-1]
	return item
}

func fillFlatAreas(data []float64, w, h int, noData float64) []float64 {
	const eps = 0.0001

	result := make([]float64, len(data))
	copy(result, data)

	isNoData := func(v float64) bool { return v == noData || math.IsNaN(v) }
	valid := func(v float64) bool { return !isNoData(v) }

	visited := make([]bool, w*h)
	queue := make(floodQueue, 0, w*h)

	seed := func(idx int) {
		if !visited[idx] && valid(result[idx]) {
			visited[idx] = true
			queue = append(queue, floodItem{result[idx], idx})
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if !valid(result[idx]) {
				continue
			}
			if x == 0 || x == w-1 || y == 0 || y == h-1 {
				seed(idx)
				continue
			}
			if isNoData(result[idx-1]) || isNoData(result[idx+1]) ||
				isNoData(result[idx-w]) || isNoData(result[idx+w]) {
				seed(idx)
			}
		}
	}

	if len(queue) == 0 {
		return result
	}
	heap.Init(&queue)

	for queue.Len() > 0 {
		cur := heap.Pop(&queue).(floodItem)
		cx, cy := cur.idx%w, cur.idx/w

		for _, d := range [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			nx, ny := cx+d[0], cy+d[1]
			if nx < 0 || nx >= w || ny < 0 || ny >= h {
				continue
			}
			nidx := ny*w + nx
			if visited[nidx] || !valid(result[nidx]) {
				continue
			}
			visited[nidx] = true
			nz := result[nidx]
			if nz <= cur.z {
				nz = cur.z + eps
				result[nidx] = nz
			}
			heap.Push(&queue, floodItem{nz, nidx})
		}
	}

	return result
}
