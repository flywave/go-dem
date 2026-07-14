package pointz

import (
	"math"
	"sort"
)

type BlockThinMode string

const (
	BlockThinMin    BlockThinMode = "min"
	BlockThinMax    BlockThinMode = "max"
	BlockThinMean   BlockThinMode = "mean"
	BlockThinMedian BlockThinMode = "median"
)

type BlockThinOptions struct {
	Resolution float64
	Mode       BlockThinMode
}

func BlockThinFilter(points []Point3D, opts *BlockThinOptions) []bool {
	if len(points) == 0 {
		return nil
	}

	res := opts.Resolution
	if res <= 0 {
		res = 10
	}
	mode := opts.Mode
	if mode == "" {
		mode = BlockThinMin
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

	type cellPoint struct {
		idx   int
		z     float64
	}
	cells := make(map[int64][]cellPoint)

	for i, p := range points {
		gx := int(math.Floor((p.X - minX) / res))
		gy := int(math.Floor((p.Y - minY) / res))
		key := int64(gx)<<32 | int64(gy)&0xffffffff
		cells[key] = append(cells[key], cellPoint{idx: i, z: p.Z})
	}

	keep := make(map[int]bool)

	for _, pts := range cells {
		if len(pts) == 0 {
			continue
		}
		var winner int
		switch mode {
		case BlockThinMax:
			maxZ := pts[0].z
			winner = pts[0].idx
			for _, cp := range pts[1:] {
				if cp.z > maxZ {
					maxZ = cp.z
					winner = cp.idx
				}
			}
		case BlockThinMean:
			var sum float64
			for _, cp := range pts {
				sum += cp.z
			}
			mean := sum / float64(len(pts))
			bestDist := math.Abs(pts[0].z - mean)
			winner = pts[0].idx
			for _, cp := range pts[1:] {
				d := math.Abs(cp.z - mean)
				if d < bestDist {
					bestDist = d
					winner = cp.idx
				}
			}
		case BlockThinMedian:
			zs := make([]float64, len(pts))
			for i, cp := range pts {
				zs[i] = cp.z
			}
			sort.Float64s(zs)
			med := zs[len(zs)/2]
			bestDist := math.Abs(pts[0].z - med)
			winner = pts[0].idx
			for _, cp := range pts[1:] {
				d := math.Abs(cp.z - med)
				if d < bestDist {
					bestDist = d
					winner = cp.idx
				}
			}
		default:
			minZ := pts[0].z
			winner = pts[0].idx
			for _, cp := range pts[1:] {
				if cp.z < minZ {
					minZ = cp.z
					winner = cp.idx
				}
			}
		}
		keep[winner] = true
	}

	mask := make([]bool, len(points))
	for i := range mask {
		mask[i] = !keep[i]
	}
	return mask
}
