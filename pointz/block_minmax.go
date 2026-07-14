package pointz

import (
	"math"
)

type BlockMinMaxMode string

const (
	BlockMinMaxMin BlockMinMaxMode = "min"
	BlockMinMaxMax BlockMinMaxMode = "max"
)

type BlockMinMaxOptions struct {
	Resolution float64
	Mode       BlockMinMaxMode
	Invert     bool
}

func BlockMinMaxFilter(points []Point3D, opts *BlockMinMaxOptions) []bool {
	if len(points) == 0 {
		return nil
	}
	res := opts.Resolution
	if res <= 0 {
		res = 10
	}
	mode := opts.Mode
	if mode == "" {
		mode = BlockMinMaxMin
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

	type cellIdx struct {
		orig int
		z    float64
	}
	cells := make(map[int64][]cellIdx)

	for i, p := range points {
		gx := int(math.Floor((p.X - minX) / res))
		gy := int(math.Floor((p.Y - minY) / res))
		key := int64(gx)<<32 | int64(gy)&0xffffffff
		cells[key] = append(cells[key], cellIdx{orig: i, z: p.Z})
	}

	keep := make(map[int]bool)
	for _, pts := range cells {
		if len(pts) == 0 {
			continue
		}
		winner := pts[0].orig
		if mode == BlockMinMaxMax {
			maxZ := pts[0].z
			for _, cp := range pts[1:] {
				if cp.z > maxZ {
					maxZ = cp.z
					winner = cp.orig
				}
			}
		} else {
			minZ := pts[0].z
			for _, cp := range pts[1:] {
				if cp.z < minZ {
					minZ = cp.z
					winner = cp.orig
				}
			}
		}
		keep[winner] = true
	}

	mask := make([]bool, len(points))
	for i := range mask {
		mask[i] = !keep[i]
	}
	if opts.Invert {
		for i := range mask {
			mask[i] = !mask[i]
		}
	}
	return mask
}
