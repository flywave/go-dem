package waffle

import (
	"fmt"
	"math"
	"sort"

	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

type cubeWaffle struct {
	baseWaffle
}

type soundingParams struct {
	TVUa float64
	TVUb float64
	THU  float64
}

func init() {
	Register("cube", func() Waffle {
		return &cubeWaffle{baseWaffle: baseWaffle{name: "cube"}}
	})
}

type cubeHypothesis struct {
	mean   float64
	stdDev float64
	count  int
	score  float64
}

func tvu(depth, a, b float64) float64 {
	return math.Sqrt(a*a + (b*depth)*(b*depth))
}

func (cw *cubeWaffle) Run(points []Point, opts *Options) (*Result, error) {
	if len(points) == 0 {
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

	pts := make([]vec2.T, len(points))
	zs := make([]float64, len(points))
	for i, p := range points {
		pts[i] = p.Position
		zs[i] = p.Z
	}

	sParams := soundingParams{
		TVUa: 0.2,
		TVUb: 0.01,
		THU:  2.0,
	}

	kdtree := NewKDTree(pts)
	width := region.XSize
	height := region.YSize

	demData := make([]float64, width*height)
	uncData := make([]float64, width*height)
	hypCount := make([]int, width*height)

	for i := range demData {
		demData[i] = noData
		uncData[i] = noData
	}

	cellHypotheses := make([][]cubeHypothesis, width*height)
	for i := range cellHypotheses {
		cellHypotheses[i] = nil
	}

	gt := region.GeoTransform()
	searchRadius := region.XRes * 3
	thuCells := int(math.Ceil(sParams.THU / region.XRes))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
			geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]

			q := vec2.T{geoX, geoY}
			searchR := searchRadius + sParams.THU
			idxs, _ := kdtree.RadiusSearch(q, searchR)
			if len(idxs) < 5 {
				idxs2, _ := kdtree.KNN(q, 5)
				idxs = idxs2
			}
			if len(idxs) < 3 {
				continue
			}

			depths := make([]float64, len(idxs))
			weights := make([]float64, len(idxs))
			for i, idx := range idxs {
				depths[i] = zs[idx]
				d := math.Sqrt(distSq(geoX, geoY, pts[idx][0], pts[idx][1]))
				weights[i] = 1.0 / (tvu(zs[idx], sParams.TVUa, sParams.TVUb) + 0.01)
				if d > 0 {
					weights[i] /= d
				}
			}

			h := densityClusterHypotheses(depths, weights, sParams)

			for _, neighborHyp := range cellHypotheses[max(0, y-thuCells)*width+max(0, x-thuCells)] {
				if y+thuCells < height && x+thuCells < width {
					_ = neighborHyp
				}
			}

			if len(h) == 0 {
				continue
			}

			best := selectHypothesisByIC(h, depths, weights)
			if best == nil {
				continue
			}

			idx := y*width + x
			demData[idx] = best.mean
			uncData[idx] = best.stdDev
			hypCount[idx] = best.count

			cellHypotheses[idx] = h
		}
	}

	for iter := 0; iter < 2; iter++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := y*width + x
				if demData[idx] != noData {
					continue
				}

				var merged []cubeHypothesis
				for dy := -2; dy <= 2; dy++ {
					for dx := -2; dx <= 2; dx++ {
						nx, ny := x+dx, y+dy
						if nx < 0 || nx >= width || ny < 0 || ny >= height {
							continue
						}
						nidx := ny*width + nx
						merged = append(merged, cellHypotheses[nidx]...)
					}
				}
				if len(merged) > 0 {
					best := selectHypothesisByIC(merged, nil, nil)
					if best != nil {
						demData[idx] = best.mean
						uncData[idx] = best.stdDev * 1.5
						hypCount[idx] = best.count
					}
				}
			}
		}
	}

	return &Result{
		DEM:    demData,
		Region: region,
		Stack:  [][]float64{demData, uncData},
	}, nil
}

func densityClusterHypotheses(depths, weights []float64, sp soundingParams) []cubeHypothesis {
	if len(depths) < 3 {
		return nil
	}

	type dp struct {
		depth  float64
		weight float64
	}
	sorted := make([]dp, len(depths))
	for i := range depths {
		sorted[i] = dp{depths[i], weights[i]}
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].depth < sorted[j].depth })

	q1 := sorted[len(sorted)/4].depth
	q3 := sorted[len(sorted)*3/4].depth
	iqr := q3 - q1
	bandwidth := iqr * 1.5
	if bandwidth < tvu(sorted[len(sorted)/2].depth, sp.TVUa, sp.TVUb) {
		bandwidth = tvu(sorted[len(sorted)/2].depth, sp.TVUa, sp.TVUb)
	}

	var hypotheses []cubeHypothesis
	used := make([]bool, len(sorted))

	for i := 0; i < len(sorted); i++ {
		if used[i] {
			continue
		}
		var cluster []dp
		cluster = append(cluster, sorted[i])
		used[i] = true

		clusterMean := sorted[i].depth
		for j := i + 1; j < len(sorted); j++ {
			if used[j] {
				continue
			}
			if math.Abs(sorted[j].depth-clusterMean) <= bandwidth {
				cluster = append(cluster, sorted[j])
				used[j] = true
				n := float64(len(cluster))
				clusterMean = clusterMean*(n-1)/n + sorted[j].depth/n
			}
		}

		if len(cluster) >= 3 {
			var mean, sumW float64
			for _, p := range cluster {
				mean += p.depth * p.weight
				sumW += p.weight
			}
			if sumW > 0 {
				mean /= sumW
			}
			var variance float64
			for _, p := range cluster {
				d := p.depth - mean
				variance += d * d * p.weight
			}
			stdDev := math.Sqrt(variance / sumW)

			hypotheses = append(hypotheses, cubeHypothesis{
				mean:   mean,
				stdDev: stdDev,
				count:  len(cluster),
			})
		}
	}

	return hypotheses
}

func selectHypothesisByIC(hypotheses []cubeHypothesis, depths, weights []float64) *cubeHypothesis {
	if len(hypotheses) == 0 {
		return nil
	}
	if len(hypotheses) == 1 {
		return &hypotheses[0]
	}

	n := float64(len(depths))
	if n == 0 {
		n = 1
	}

	bestScore := math.MaxFloat64
	bestIdx := 0

	for i, h := range hypotheses {
		logLikelihood := 0.0
		if depths != nil {
			for _, d := range depths {
				if h.stdDev > 1e-15 {
					dev := (d - h.mean) / h.stdDev
					logLikelihood -= 0.5 * dev * dev
				}
			}
		} else {
			logLikelihood = -float64(h.count)
		}

		k := float64(len(hypotheses))
		score := -2*logLikelihood + k*math.Log(n)

		bonus := float64(h.count) * 0.2
		score -= bonus

		if score < bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	return &hypotheses[bestIdx]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (cw *cubeWaffle) Name() string { return cw.name }
