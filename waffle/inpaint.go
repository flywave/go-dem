package waffle

import (
	"math"
	"sort"

	"github.com/flywave/go-dem"
)

type inpaintWaffle struct {
	baseWaffle
}

func init() {
	Register("inpaint", func() Waffle {
		return &inpaintWaffle{baseWaffle: baseWaffle{name: "inpaint"}}
	})
}

func (iw *inpaintWaffle) Run(points []Point, opts *Options) (*Result, error) {
	if len(points) == 0 {
		return nil, nil
	}

	region := opts.Region
	if region.XSize <= 0 || region.YSize <= 0 {
		region.XSize = int(math.Round((region.BBox().Max[0] - region.BBox().Min[0]) / region.XRes))
		region.YSize = int(math.Round((region.BBox().Max[1] - region.BBox().Min[1]) / region.YRes))
	}

	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	width := region.XSize
	height := region.YSize
	demData := make([]float64, width*height)
	for i := range demData {
		demData[i] = noData
	}

	gt := region.GeoTransform()
	for _, p := range points {
		px := int(math.Round((p.Position[0] - gt[0]) / gt[1]))
		py := int(math.Round((p.Position[1] - gt[3]) / gt[5]))
		if px >= 0 && px < width && py >= 0 && py < height {
			demData[py*width+px] = p.Z
		}
	}

	demData = inpaintFMM(demData, width, height, noData)
	return &Result{DEM: demData, Region: region}, nil
}

const (
	BAND_KNOWN  = 0
	BAND_BAND   = 1
	BAND_INSIDE = 2
)

type fmmPixel struct {
	x, y int
	dist float64
}

type fmmPriorityQueue []fmmPixel

func (pq fmmPriorityQueue) Len() int           { return len(pq) }
func (pq fmmPriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq fmmPriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *fmmPriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(fmmPixel)) }
func (pq *fmmPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func inpaintFMM(data []float64, w, h int, noData float64) []float64 {
	result := make([]float64, len(data))
	copy(result, data)

	flag := make([]int, w*h)
	dist := make([]float64, w*h)
	for i := range dist {
		dist[i] = 1e18
	}

	unknownCount := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if result[idx] == noData || math.IsNaN(result[idx]) {
				flag[idx] = BAND_INSIDE
				unknownCount++
			} else {
				flag[idx] = BAND_KNOWN
			}
		}
	}

	if unknownCount == 0 {
		return result
	}

	pq := make(fmmPriorityQueue, 0)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if flag[y*w+x] != BAND_INSIDE {
				continue
			}
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx, ny := x+dx, y+dy
					if nx < 0 || nx >= w || ny < 0 || ny >= h {
						continue
					}
					if flag[ny*w+nx] == BAND_KNOWN {
						flag[y*w+x] = BAND_BAND
						dist[y*w+x] = 0
						pq = append(pq, fmmPixel{x, y, 0})
						goto nextBandInit
					}
				}
			}
		nextBandInit:
		}
	}

	if len(pq) == 0 {
		return result
	}

	for len(pq) > 0 {
		sort.Sort(pq)
		p := pq[0]
		pq = pq[1:]

		px, py := p.x, p.y
		pIdx := py*w + px
		if flag[pIdx] != BAND_BAND {
			continue
		}

		flag[pIdx] = BAND_KNOWN

		if result[pIdx] == noData || math.IsNaN(result[pIdx]) {
			result[pIdx] = fmmInpaintValue(result, flag, dist, px, py, w, h, noData)
		}

		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				nx, ny := px+dx, py+dy
				if nx < 0 || nx >= w || ny < 0 || ny >= h {
					continue
				}
				nIdx := ny*w + nx
				if flag[nIdx] == BAND_INSIDE {
					flag[nIdx] = BAND_BAND
					newDist := math.Sqrt(float64(dx*dx + dy*dy))
					if newDist < dist[nIdx] {
						dist[nIdx] = newDist
					}
					pq = append(pq, fmmPixel{nx, ny, dist[nIdx]})
				} else if flag[nIdx] == BAND_BAND {
					newDist := math.Sqrt(float64(dx*dx + dy*dy))
					if newDist < dist[nIdx] {
						dist[nIdx] = newDist
					}
				}
			}
		}
	}

	for i := range result {
		if result[i] == noData || math.IsNaN(result[i]) {
			result[i] = finalFallbackValue(result, i, w, h, noData)
		}
	}

	return result
}

func fmmInpaintValue(data []float64, flag []int, dist []float64, x, y, w, h int, noData float64) float64 {
	const searchRadius = 6

	type neighbor struct {
		val    float64
		weight float64
	}
	var neighbors []neighbor

	for dy := -searchRadius; dy <= searchRadius; dy++ {
		for dx := -searchRadius; dx <= searchRadius; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx < 0 || nx >= w || ny < 0 || ny >= h {
				continue
			}
			nIdx := ny*w + nx
			if flag[nIdx] != BAND_KNOWN {
				continue
			}
			if data[nIdx] == noData || math.IsNaN(data[nIdx]) {
				continue
			}
			d := math.Sqrt(float64(dx*dx + dy*dy))
			weight := 1.0 / (d*d + 1e-15)
			neighbors = append(neighbors, neighbor{val: data[nIdx], weight: weight})
		}
	}

	if len(neighbors) < 2 {
		if len(neighbors) == 1 {
			return neighbors[0].val
		}
		return fallbackIDW(data, flag, x, y, w, h, noData)
	}

	var sumWeight, sumVal float64
	for _, n := range neighbors {
		if !math.IsNaN(n.val) && !math.IsInf(n.val, 0) {
			sumWeight += n.weight
			sumVal += n.val * n.weight
		}
	}

	if sumWeight > 0 {
		return sumVal / sumWeight
	}
	return noData
}

func estimateGradientX(data []float64, flag []int, x, y, w, h int, noData float64) float64 {
	left := 0
	if x > 0 && flag[y*w+(x-1)] == BAND_KNOWN {
		left = 1
	}
	right := 0
	if x < w-1 && flag[y*w+(x+1)] == BAND_KNOWN {
		right = 1
	}
	if left == 0 && right == 0 {
		return 0
	}
	if left == 0 {
		return data[y*w+(x+1)] - data[y*w+x]
	}
	if right == 0 {
		return data[y*w+x] - data[y*w+(x-1)]
	}
	return (data[y*w+(x+1)] - data[y*w+(x-1)]) / 2
}

func estimateGradientY(data []float64, flag []int, x, y, w, h int, noData float64) float64 {
	up := 0
	if y > 0 && flag[(y-1)*w+x] == BAND_KNOWN {
		up = 1
	}
	down := 0
	if y < h-1 && flag[(y+1)*w+x] == BAND_KNOWN {
		down = 1
	}
	if up == 0 && down == 0 {
		return 0
	}
	if up == 0 {
		return data[(y+1)*w+x] - data[y*w+x]
	}
	if down == 0 {
		return data[y*w+x] - data[(y-1)*w+x]
	}
	return (data[(y+1)*w+x] - data[(y-1)*w+x]) / 2
}

func fallbackIDW(data []float64, flag []int, x, y, w, h int, noData float64) float64 {
	const maxSearch = 20
	var sumVal, sumWeight float64
	for r := 1; r <= maxSearch; r++ {
		for dy := -r; dy <= r; dy++ {
			for dx := -r; dx <= r; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				if math.Abs(float64(dx)) != float64(r) && math.Abs(float64(dy)) != float64(r) {
					continue
				}
				nx, ny := x+dx, y+dy
				if nx < 0 || nx >= w || ny < 0 || ny >= h {
					continue
				}
				nIdx := ny*w + nx
				if flag[nIdx] != BAND_KNOWN {
					continue
				}
				if data[nIdx] == noData || math.IsNaN(data[nIdx]) {
					continue
				}
				d := math.Sqrt(float64(dx*dx + dy*dy))
				w := 1.0 / (d + 1e-15)
				sumVal += data[nIdx] * w
				sumWeight += w
			}
		}
		if sumWeight > 0 {
			return sumVal / sumWeight
		}
	}
	return noData
}

func finalFallbackValue(data []float64, idx, w, h int, noData float64) float64 {
	x, y := idx%w, idx/w
	searchDirs := [][2]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		{-1, -1}, {-1, 1}, {1, -1}, {1, 1},
		{-2, 0}, {2, 0}, {0, -2}, {0, 2},
		{-3, 0}, {3, 0}, {0, -3}, {0, 3},
	}
	var sum, count float64
	for _, d := range searchDirs {
		nx, ny := x+d[0], y+d[1]
		if nx < 0 || nx >= w || ny < 0 || ny >= h {
			continue
		}
		val := data[ny*w+nx]
		if val != noData && !math.IsNaN(val) {
			dist := math.Sqrt(float64(d[0]*d[0] + d[1]*d[1]))
			sum += val / (dist + 1e-15)
			count += 1.0 / (dist + 1e-15)
		}
	}
	if count > 0 {
		return sum / count
	}
	return noData
}
