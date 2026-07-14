package pointz

import "math"

type DiffZOptions struct {
	MinDiff float64
	MaxDiff float64
	Invert  bool
}

func DiffZFilter(points []Point3D, opts *DiffZOptions) []bool {
	if len(points) == 0 {
		return nil
	}
	mask := make([]bool, len(points))
	if opts == nil {
		return mask
	}

	for i, p := range points {
		inside := true
		if !math.IsNaN(opts.MinDiff) && p.Z < opts.MinDiff {
			inside = false
		}
		if !math.IsNaN(opts.MaxDiff) && p.Z > opts.MaxDiff {
			inside = false
		}
		if opts.Invert {
			mask[i] = inside
		} else {
			mask[i] = !inside
		}
	}
	return mask
}
