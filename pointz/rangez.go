package pointz

import "math"

type RangeZOptions struct {
	MinZ   float64
	MaxZ   float64
	Invert bool
}

func RangeZFilter(points []Point3D, opts *RangeZOptions) []bool {
	if len(points) == 0 {
		return nil
	}
	mask := make([]bool, len(points))
	if opts == nil {
		return mask
	}

	for i, p := range points {
		outside := false
		if !math.IsNaN(opts.MinZ) && p.Z < opts.MinZ {
			outside = true
		}
		if !math.IsNaN(opts.MaxZ) && p.Z > opts.MaxZ {
			outside = true
		}
		if opts.Invert {
			mask[i] = !outside
		} else {
			mask[i] = outside
		}
	}
	return mask
}
