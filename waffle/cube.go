package waffle

import (
	"fmt"
	"math"
	"sort"

	"github.com/flywave/go-dem"
)

type cubeWaffle struct {
	baseWaffle
}

func init() {
	Register("cube", func() Waffle {
		return &cubeWaffle{baseWaffle: baseWaffle{name: "cube"}}
	})
}

type cubeParams struct {
	Resolution      float64
	SearchRadius    float64
	MinPoints       int
	MaxPoints       int
	IQRMultiplier   float64
	VerticalUnc     float64
}

type cubeHypothesis struct {
	mean   float64
	stdDev float64
	count  int
}

func (cw *cubeWaffle) Run(sources []string, opts *Options) (*Result, error) {
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

	params := cubeParams{
		Resolution:    region.XRes,
		SearchRadius:  region.XRes * 3,
		MinPoints:     5,
		MaxPoints:     30,
		IQRMultiplier: 1.5,
		VerticalUnc:   0.2,
	}

	gt := region.GeoTransform()
	width := region.XSize
	height := region.YSize

	demData := make([]float64, width*height)
	uncData := make([]float64, width*height)

	for i := range demData {
		demData[i] = noData
		uncData[i] = noData
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
			geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]

			var depthsInCell []float64
			for i, pt := range pts {
				dx := geoX - pt[0]
				dy := geoY - pt[1]
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist <= params.SearchRadius {
					depthsInCell = append(depthsInCell, zs[i])
				}
			}

			if len(depthsInCell) < params.MinPoints {
				continue
			}

			if len(depthsInCell) > params.MaxPoints {
				sort.Float64s(depthsInCell)
				depthsInCell = depthsInCell[:params.MaxPoints]
			}

			hypotheses := buildCUBEHypotheses(depthsInCell, params)
			if len(hypotheses) == 0 {
				continue
			}

			bestHyp := selectBestHypothesis(hypotheses, depthsInCell)
			if bestHyp == nil {
				continue
			}

			idx := y*width + x
			demData[idx] = bestHyp.mean
			uncData[idx] = bestHyp.stdDev
		}
	}

	result := &Result{
		DEM:    demData,
		Region: region,
		Stack:  [][]float64{demData, uncData},
	}

	return result, nil
}

func buildCUBEHypotheses(depths []float64, params cubeParams) []cubeHypothesis {
	if len(depths) < params.MinPoints {
		return nil
	}

	sorted := make([]float64, len(depths))
	copy(sorted, depths)
	sort.Float64s(sorted)

	var hypotheses []cubeHypothesis
	used := make([]bool, len(sorted))

	for i := 0; i < len(sorted); i++ {
		if used[i] {
			continue
		}

		var cluster []float64
		cluster = append(cluster, sorted[i])
		used[i] = true

		for j := i + 1; j < len(sorted); j++ {
			if used[j] {
				continue
			}
			if math.Abs(sorted[j]-sorted[i]) <= params.IQRMultiplier*params.VerticalUnc {
				cluster = append(cluster, sorted[j])
				used[j] = true
			}
		}

		if len(cluster) >= params.MinPoints {
			mean, stdDev := computeStats(cluster)
			hypotheses = append(hypotheses, cubeHypothesis{
				mean:   mean,
				stdDev: stdDev,
				count:  len(cluster),
			})
		}
	}

	return hypotheses
}

func selectBestHypothesis(hypotheses []cubeHypothesis, depths []float64) *cubeHypothesis {
	if len(hypotheses) == 0 {
		return nil
	}
	if len(hypotheses) == 1 {
		return &hypotheses[0]
	}

	bestIdx := 0
	bestScore := math.MaxFloat64

	for i, h := range hypotheses {
		var score float64
		for _, d := range depths {
			diff := math.Abs(d - h.mean)
			if h.stdDev > 0 {
				score += diff / h.stdDev
			} else {
				score += diff
			}
		}
		score /= float64(len(depths))
		score -= float64(h.count) * 0.1

		if score < bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	return &hypotheses[bestIdx]
}

func computeStats(vals []float64) (mean, stdDev float64) {
	n := float64(len(vals))
	if n == 0 {
		return 0, 0
	}
	for _, v := range vals {
		mean += v
	}
	mean /= n
	for _, v := range vals {
		diff := v - mean
		stdDev += diff * diff
	}
	stdDev = math.Sqrt(stdDev / n)
	return
}
