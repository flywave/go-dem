package grits

import (
	"math"
	"sort"

	"github.com/flywave/go-dem"
)

type outliersFilter struct {
	baseGrits
}

func init() {
	Register(FilterOutliers, func() Grits {
		return &outliersFilter{baseGrits{name: string(FilterOutliers)}}
	})
}

func (f *outliersFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	nd := opts.GetNoData()
	w, h := region.XSize, region.YSize

	multipass := opts.Iterations
	if multipass <= 0 {
		multipass = 1
	}
	percentile := opts.Percentile
	if percentile <= 0 {
		percentile = 75
	}
	k := opts.Threshold
	if k <= 0 {
		k = 1.5
	}
	mode := opts.Method

	result := make([]float64, len(data))
	copy(result, data)

	slopeData := computeSlopeLocal(data, w, h, nd)

	minWin := 5
	maxWin := 25
	if multipass > 1 {
		maxWin = 50
	}

	score := make([]float64, len(data))

	for pass := 0; pass < multipass; pass++ {
		winSize := minWin
		if multipass > 1 {
			f := float64(pass) / float64(multipass-1)
			winSize = minWin + int(f*float64(maxWin-minWin))
		}
		if winSize < 3 {
			winSize = 3
		}
		if winSize%2 == 0 {
			winSize++
		}
		half := winSize / 2

		passScore := make([]float64, len(data))
		attrs := []struct {
			name  string
			vals  []float64
			weight float64
		}{
			{"elevation", data, 1.0},
			{"slope", slopeData, 0.25},
		}

		if mode == "aggressive" || multipass > 1 {
			roughness := computeRoughness(data, w, h, nd)
			attrs = append(attrs, struct {
				name   string
				vals   []float64
				weight float64
			}{"roughness", roughness, 0.25})
		}

		for _, attr := range attrs {
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					idx := y*w + x
					if result[idx] == nd || math.IsNaN(result[idx]) {
						continue
					}

					var vals []float64
					for ky := -half; ky <= half; ky++ {
						for kx := -half; kx <= half; kx++ {
							ix, iy := x+kx, y+ky
							if ix < 0 || ix >= w || iy < 0 || iy >= h {
								continue
							}
							v := attr.vals[iy*w+ix]
							if math.IsNaN(v) {
								continue
							}
							vals = append(vals, v)
						}
					}
					if len(vals) < 4 {
						continue
					}

					sorted := make([]float64, len(vals))
					copy(sorted, vals)
					sort.Float64s(sorted)

					q1 := percentileValue(sorted, 25)
					q3 := percentileValue(sorted, percentile)
					iqr := q3 - q1
					if iqr < 1e-12 {
						continue
					}
					upper := q3 + k*iqr
					lower := q1 - k*iqr

					val := attr.vals[idx]
					if val > upper || val < lower {
						diff := 0.0
						if val > upper {
							diff = (val - upper) / (q3 - q1 + 1e-12)
						} else {
							diff = (lower - val) / (q3 - q1 + 1e-12)
						}
						passScore[idx] += attr.weight * math.Min(diff, 3.0)
					}
				}
			}
		}

		for i := range score {
			score[i] = math.Sqrt(score[i]*score[i] + passScore[i]*passScore[i])
		}
	}

	scoreThreshold := float64(multipass) * 0.5
	for i := range result {
		if result[i] == nd || math.IsNaN(result[i]) {
			continue
		}
		if score[i] > scoreThreshold {
			result[i] = nd
		}
	}

	return result, nil
}

func computeSlopeLocal(data []float64, w, h int, nd float64) []float64 {
	slope := make([]float64, len(data))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] == nd || math.IsNaN(data[idx]) {
				slope[idx] = math.NaN()
				continue
			}
			dzdx, dzdy := 0.0, 0.0
			c := 0
			if x > 0 && data[idx-1] != nd && !math.IsNaN(data[idx-1]) {
				dzdx += data[idx] - data[idx-1]
				c++
			}
			if x < w-1 && data[idx+1] != nd && !math.IsNaN(data[idx+1]) {
				dzdx += data[idx+1] - data[idx]
				c++
			}
			if c > 0 {
				dzdx /= float64(c)
			}
			c = 0
			if y > 0 && data[idx-w] != nd && !math.IsNaN(data[idx-w]) {
				dzdy += data[idx] - data[idx-w]
				c++
			}
			if y < h-1 && data[idx+w] != nd && !math.IsNaN(data[idx+w]) {
				dzdy += data[idx+w] - data[idx]
				c++
			}
			if c > 0 {
				dzdy /= float64(c)
			}
			slope[idx] = math.Atan(math.Sqrt(dzdx*dzdx+dzdy*dzdy)) * 180 / math.Pi
		}
	}
	return slope
}

func computeRoughness(data []float64, w, h int, nd float64) []float64 {
	rough := make([]float64, len(data))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] == nd || math.IsNaN(data[idx]) {
				rough[idx] = math.NaN()
				continue
			}
			minZ, maxZ := data[idx], data[idx]
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					ix, iy := x+kx, y+ky
					if ix < 0 || ix >= w || iy < 0 || iy >= h {
						continue
					}
					v := data[iy*w+ix]
					if v == nd || math.IsNaN(v) {
						continue
					}
					if v < minZ {
						minZ = v
					}
					if v > maxZ {
						maxZ = v
					}
				}
			}
			rough[idx] = maxZ - minZ
		}
	}
	return rough
}

func percentileValue(sorted []float64, pct float64) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	idx := pct / 100.0 * float64(n-1)
	lower := int(math.Floor(idx))
	upper := int(math.Ceil(idx))
	if lower == upper || upper >= n {
		return sorted[lower]
	}
	frac := idx - float64(lower)
	return sorted[lower]*(1-frac) + sorted[upper]*frac
}
