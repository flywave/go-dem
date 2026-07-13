package waffle

import (
	"fmt"
	"math"
	"sort"

	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

type cudemWaffle struct {
	baseWaffle
}

func init() {
	Register(dem.MethodCUDEM, func() Waffle {
		return &cudemWaffle{baseWaffle: baseWaffle{name: string(dem.MethodCUDEM)}}
	})
}

type cudemCell struct {
	mean   float64
	stdDev float64
	count  int
	weight float64
}

func (cw *cudemWaffle) Run(sources []string, opts *Options) (*Result, error) {
	pts, zs, err := collectPoints(sources)
	if err != nil {
		return nil, err
	}
	if len(pts) == 0 {
		return nil, fmt.Errorf("no data points")
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
	gt := region.GeoTransform()

	kdtree := NewKDTree(pts)

	levels := computeLevelsByDensity(pts, kdtree, region.XRes, width, height, gt)
	if len(levels) == 0 {
		levels = []float64{region.XRes}
	}

	type weightedCell struct {
		mean   float64
		weight float64
		count  int
	}

	result := make([]weightedCell, width*height)
	for i := range result {
		result[i] = weightedCell{weight: 0, count: 0}
	}

	sort.Float64s(levels)
	for li := len(levels) - 1; li >= 0; li-- {
		res := levels[li]
		scale := int(math.Round(res / region.XRes))
		if scale < 1 {
			scale = 1
		}

		gw := (width + scale - 1) / scale
		gh := (height + scale - 1) / scale

		grid := make([]*cudemCell, gw*gh)
		searchR := res * 3

		for i, pt := range pts {
			gx := int(math.Floor((pt[0] - gt[0]) / res))
			gy := int(math.Floor((pt[1] - gt[3]) / res))
			if gx < 0 || gx >= gw || gy < 0 || gy >= gh {
				continue
			}
			idx := gy*gw + gx
			if grid[idx] == nil {
				grid[idx] = &cudemCell{mean: zs[i], count: 1, weight: 1.0}
			} else {
				c := grid[idx]
				c.mean = (c.mean*float64(c.count) + zs[i]) / float64(c.count+1)
				c.count++
			}
		}

		for gy := 0; gy < gh; gy++ {
			for gx := 0; gx < gw; gx++ {
				idx := gy*gw + gx
				if grid[idx] != nil {
					continue
				}
				geoX := gt[0] + (float64(gx)+0.5)*res
				geoY := gt[3] + (float64(gy)+0.5)*res
				q := vec2.T{geoX, geoY}
				idxs, dists := kdtree.RadiusSearch(q, searchR)
				if len(idxs) < 3 {
					idxs2, dists2 := kdtree.KNN(q, 5)
					idxs, dists = idxs2, dists2
				}
				if len(idxs) < 3 {
					continue
				}
				var sumW, sumV, minDist float64
				for i, idx := range idxs {
					d := dists[i]
					if d < 1e-10 {
						sumV, sumW = zs[idx], 1
						minDist = 0
						break
					}
					w := 1.0 / (d*d + 1e-15)
					sumW += w
					sumV += w * zs[idx]
					if i == 0 || d < minDist {
						minDist = d
					}
				}
				if sumW > 0 {
					grid[idx] = &cudemCell{
						mean:   sumV / sumW,
						count:  len(idxs),
						weight: 1.0 / (minDist + res*0.1),
					}
				}
			}
		}

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := y*width + x
				gx, gy := x/scale, y/scale
				if gx >= gw || gy >= gh {
					continue
				}
				cell := grid[gy*gw+gx]
				if cell == nil {
					continue
				}

				geoX := gt[0] + float64(x)*gt[1]
				geoY := gt[3] + float64(y)*gt[5]
				q := vec2.T{geoX, geoY}
				localIdxs, _ := kdtree.RadiusSearch(q, res)
				localCount := len(localIdxs)

				if localCount >= 3 && res <= region.XRes*1.5 {
					continue
				}

				cellWeight := cell.weight
				if localCount > 0 {
					cellWeight *= float64(cell.count) / float64(cell.count+localCount)
				}

				current := &result[idx]
				if cellWeight > current.weight {
					current.mean = cell.mean
					current.weight = cellWeight
					current.count = cell.count
				} else if cellWeight > 0 && current.count > 0 {
					totalW := current.weight + cellWeight
					current.mean = (current.mean*current.weight + cell.mean*cellWeight) / totalW
					current.weight = totalW
					current.count += cell.count
				}
			}
		}
	}

	demData := make([]float64, width*height)
	for i := range demData {
		if result[i].count > 0 {
			demData[i] = result[i].mean
		} else {
			demData[i] = noData
		}
	}

	return &Result{DEM: demData, Region: region}, nil
}

func computeLevelsByDensity(pts []vec2.T, kdtree *KDTree, baseRes float64, width, height int, gt [6]float64) []float64 {
	type cellDensity struct {
		cx, cy int
		count  int
	}

	densityGrid := make(map[int]map[int]int)
	step := 10
	if step > width/4 {
		step = width / 4
	}
	if step < 1 {
		step = 1
	}

	for _, pt := range pts {
		cx := int((pt[0]-gt[0])/gt[1]) / step
		if cx < 0 {
			cx = 0
		}
		cy := int((gt[3]-pt[1])/(-gt[5])) / step
		if densityGrid[cy] == nil {
			densityGrid[cy] = make(map[int]int)
		}
		densityGrid[cy][cx]++
	}

	var densities []int
	for _, row := range densityGrid {
		for _, c := range row {
			if c > 0 {
				densities = append(densities, c)
			}
		}
	}

	if len(densities) == 0 {
		return []float64{baseRes}
	}

	sort.Ints(densities)
	median := densities[len(densities)/2]

	nLevels := 2
	if median < 5 {
		nLevels = 4
	} else if median < 20 {
		nLevels = 3
	}

	levels := make([]float64, nLevels)
	for i := 0; i < nLevels; i++ {
		levels[i] = baseRes * math.Pow(2, float64(nLevels-1-i))
	}

	var unique []float64
	seen := make(map[float64]bool)
	for _, l := range levels {
		if !seen[l] {
			seen[l] = true
			unique = append(unique, l)
		}
	}

	return unique
}
