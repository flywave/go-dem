package pointz

import "math"

type DiffZOptions struct {
	MinDiff   float64
	MaxDiff   float64
	MinMaxSet bool
	Invert    bool
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
		inside := !math.IsNaN(p.Z)
		if inside && opts.MinMaxSet {
			inside = p.Z >= opts.MinDiff && p.Z <= opts.MaxDiff
		}
		if opts.Invert {
			mask[i] = inside
		} else {
			mask[i] = !inside
		}
	}
	return mask
}
