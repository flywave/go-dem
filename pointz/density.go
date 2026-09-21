package pointz

import (
	"math"
	"math/rand"
	"sort"
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
	if opts == nil {
		opts = &DensityOptions{}
	}
	res := opts.Resolution
	if res <= 0 {
		res = 10
	}
	mode := opts.Mode
	if mode == "" {
		mode = DensityRandom
	}

	minX, minY := points[0].X, points[0].Y
	for _, p := range points[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
	}

	type cpoint struct {
		orig    int
		x, y, z float64
		gx, gy  int
	}
	cells := make(map[int64][]cpoint)

	for i, p := range points {
		gx := int(math.Floor((p.X - minX) / res))
		gy := int(math.Floor((p.Y - minY) / res))
		key := int64(gx)<<32 | int64(gy)&0xffffffff
		cells[key] = append(cells[key], cpoint{orig: i, x: p.X, y: p.Y, z: p.Z, gx: gx, gy: gy})
	}

	keep := make(map[int]bool)
	for _, pts := range cells {
		if len(pts) == 0 {
			continue
		}
		var winner int

		switch mode {
		case DensityMedian, DensityMean:
			var center float64
			if mode == DensityMean {
				var sumZ float64
				for _, cp := range pts {
					sumZ += cp.z
				}
				center = sumZ / float64(len(pts))
			} else {
				zs := make([]float64, len(pts))
				for i, cp := range pts {
					zs[i] = cp.z
				}
				sort.Float64s(zs)
				center = zs[len(zs)/2]
			}
			winner = pts[0].orig
			bestDist := math.Abs(pts[0].z - center)
			for _, cp := range pts[1:] {
				d := math.Abs(cp.z - center)
				if d < bestDist {
					bestDist = d
					winner = cp.orig
				}
			}
		case DensityCenter:
			cx := minX + (float64(pts[0].gx)+0.5)*res
			cy := minY + (float64(pts[0].gy)+0.5)*res
			winner = pts[0].orig
			bestDist := (pts[0].x-cx)*(pts[0].x-cx) + (pts[0].y-cy)*(pts[0].y-cy)
			for _, cp := range pts[1:] {
				d := (cp.x-cx)*(cp.x-cx) + (cp.y-cy)*(cp.y-cy)
				if d < bestDist {
					bestDist = d
					winner = cp.orig
				}
			}
		default:
			winner = pts[rand.Intn(len(pts))].orig
		}
		keep[winner] = true
	}

	mask := make([]bool, len(points))
	for i := range mask {
		mask[i] = !keep[i]
	}
	return mask
}
