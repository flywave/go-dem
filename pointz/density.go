package pointz

import (
	"math"
)

type DensityMode string

const (
	DensityRandom DensityMode = "random"
	DensityMedian DensityMode = "median"
	DensityMean   DensityMode = "mean"
	DensityCenter DensityMode = "center"
)

type DensityOptions struct {
	Resolution float64
	Mode       DensityMode
}

func DensityFilter(points []Point3D, opts *DensityOptions) []bool {
	if len(points) == 0 {
		return nil
	}
	res := opts.Resolution
	if res <= 0 {
		res = 10
	}
	mode := opts.Mode
	if mode == "" {
		mode = DensityRandom
	}

	var minX, minY float64
	for i, p := range points {
		if i == 0 {
			minX, minY = p.X, p.Y
		} else {
			if p.X < minX {
				minX = p.X
			}
			if p.Y < minY {
				minY = p.Y
			}
		}
	}

	type cpoint struct {
		orig int
		x, y, z float64
	}
	cells := make(map[int64][]cpoint)

	for i, p := range points {
		gx := int(math.Floor((p.X - minX) / res))
		gy := int(math.Floor((p.Y - minY) / res))
		key := int64(gx)<<32 | int64(gy)&0xffffffff
		cells[key] = append(cells[key], cpoint{orig: i, x: p.X, y: p.Y, z: p.Z})
	}

	keep := make(map[int]bool)
	for _, pts := range cells {
		if len(pts) == 0 {
			continue
		}
		var winner int

		switch mode {
		case DensityMedian, DensityMean:
			var sumZ float64
			for _, cp := range pts {
				sumZ += cp.z
			}
			meanZ := sumZ / float64(len(pts))
			if mode == DensityMean {
				bestDist := math.Abs(pts[0].z - meanZ)
				winner = pts[0].orig
				for _, cp := range pts[1:] {
					d := math.Abs(cp.z - meanZ)
					if d < bestDist {
						bestDist = d
						winner = cp.orig
					}
				}
			} else {
				winner = pts[0].orig
				bestDist := math.Abs(pts[0].z - meanZ)
				for _, cp := range pts[1:] {
					d := math.Abs(cp.z - meanZ)
					if d < bestDist {
						bestDist = d
						winner = cp.orig
					}
				}
			}
		case DensityCenter:
			cx := (pts[0].x + pts[len(pts)-1].x) / 2
			cy := (pts[0].y + pts[len(pts)-1].y) / 2
			bestDist := (pts[0].x-cx)*(pts[0].x-cx) + (pts[0].y-cy)*(pts[0].y-cy)
			winner = pts[0].orig
			for _, cp := range pts[1:] {
				d := (cp.x-cx)*(cp.x-cx) + (cp.y-cy)*(cp.y-cy)
				if d < bestDist {
					bestDist = d
					winner = cp.orig
				}
			}
		default:
			winner = pts[0].orig
		}
		keep[winner] = true
	}

	mask := make([]bool, len(points))
	for i := range mask {
		mask[i] = !keep[i]
	}
	return mask
}
