package pointz

import (
	"math"
	"sort"
)

type Point3D struct {
	X, Y, Z float64
}

type OutlierZOptions struct {
	Percentile    float64
	MaxPercentile float64
	Multipass     int
	Resolution    float64
	MaxResolution float64
	Invert        bool
}

func OutlierZFilter(points []Point3D, opts *OutlierZOptions) []bool {
	if len(points) == 0 {
		return nil
	}

	percentile := opts.Percentile
	if percentile <= 0 {
		percentile = 98
	}
	maxPercentile := opts.MaxPercentile
	if maxPercentile <= 0 {
		maxPercentile = 99.9
	}
	multipass := opts.Multipass
	if multipass <= 0 {
		multipass = 4
	}
	res := opts.Resolution
	if res <= 0 {
		res = 50
	}
	maxRes := opts.MaxResolution
	if maxRes <= 0 {
		maxRes = 5000
	}

	mask := make([]bool, len(points))

	for pass := 0; pass < multipass; pass++ {
		f := float64(pass) / float64(multipass-1)
		passPerc := percentile + f*(maxPercentile-percentile)
		passRes := maxRes - f*(maxRes-res)

		count := 0
		for _, m := range mask {
			if !m {
				count++
			}
		}
		if count == 0 {
			break
		}

		validIdxs := make([]int, 0, count)
		for i, m := range mask {
			if !m {
				validIdxs = append(validIdxs, i)
			}
		}

		var minX, minY, maxX, maxY float64
		for _, idx := range validIdxs {
			p := points[idx]
			if idx == validIdxs[0] {
				minX, maxX = p.X, p.X
				minY, maxY = p.Y, p.Y
			} else {
				if p.X < minX {
					minX = p.X
				}
				if p.X > maxX {
					maxX = p.X
				}
				if p.Y < minY {
					minY = p.Y
				}
				if p.Y > maxY {
					maxY = p.Y
				}
			}
		}

		type cell struct {
			sumZ  float64
			count int
		}
		xSize := int((maxX-minX)/passRes) + 1
		ySize := int((maxY-minY)/passRes) + 1
		if xSize <= 0 || ySize <= 0 {
			continue
		}
		grid := make([]cell, xSize*ySize)

		for _, idx := range validIdxs {
			p := points[idx]
			gx := int((p.X - minX) / passRes)
			gy := int((p.Y - minY) / passRes)
			if gx < 0 {
				gx = 0
			}
			if gx >= xSize {
				gx = xSize - 1
			}
			if gy < 0 {
				gy = 0
			}
			if gy >= ySize {
				gy = ySize - 1
			}
			c := &grid[gy*xSize+gx]
			c.sumZ += p.Z
			c.count++
		}

		residuals := make([]float64, len(validIdxs))
		for j, idx := range validIdxs {
			p := points[idx]
			gx := int((p.X - minX) / passRes)
			gy := int((p.Y - minY) / passRes)
			if gx < 0 {
				gx = 0
			}
			if gx >= xSize {
				gx = xSize - 1
			}
			if gy < 0 {
				gy = 0
			}
			if gy >= ySize {
				gy = ySize - 1
			}
			c := grid[gy*xSize+gx]
			meanZ := c.sumZ / float64(c.count)
			diff := p.Z - meanZ
			if diff < 0 {
				diff = -diff
			}
			residuals[j] = diff
		}

		var validResiduals []float64
		for _, r := range residuals {
			if !math.IsNaN(r) && !math.IsInf(r, 0) {
				validResiduals = append(validResiduals, r)
			}
		}
		if len(validResiduals) == 0 {
			continue
		}
		sort.Float64s(validResiduals)

		threshold := percentileValueFloat64(validResiduals, passPerc)

		for j, idx := range validIdxs {
			if residuals[j] > threshold {
				mask[idx] = true
			}
		}
	}

	if opts.Invert {
		for i := range mask {
			mask[i] = !mask[i]
		}
	}

	return mask
}

func percentileValueFloat64(sorted []float64, pct float64) float64 {
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
