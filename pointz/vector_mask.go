package pointz

import (
	"fmt"

	gdal "github.com/flywave/flywave-gdal"
)

type VectorMaskOptions struct {
	MaskPath string
	Invert   bool
}

func VectorMaskFilter(points []Point3D, opts *VectorMaskOptions) ([]bool, error) {
	if len(points) == 0 {
		return nil, nil
	}
	if opts == nil || opts.MaskPath == "" {
		return make([]bool, len(points)), nil
	}

	ds := gdal.OpenDataSource(opts.MaskPath, 0)
	if ds.LayerCount() == 0 {
		return nil, fmt.Errorf("no layers in vector mask: %s", opts.MaskPath)
	}
	defer ds.Destroy()

	inside := make([]bool, len(points))

	for li := 0; li < ds.LayerCount(); li++ {
		layer := ds.LayerByIndex(li)
		layer.ResetReading()

		for {
			feat := layer.NextFeature()
			if feat == nil {
				break
			}
			geom := feat.Geometry()
			if geom.IsEmpty() {
				continue
			}

			for i, p := range points {
				if inside[i] {
					continue
				}
				srs := layer.SpatialReference()
				ptGeom, err := gdal.CreateFromWKT(
					fmt.Sprintf("POINT (%.10f %.10f)", p.X, p.Y),
					srs,
				)
				if err != nil {
					continue
				}
				if geom.Contains(ptGeom) {
					inside[i] = true
				}
			}
		}
	}

	mask := make([]bool, len(points))
	for i := range points {
		if opts.Invert {
			mask[i] = inside[i]
		} else {
			mask[i] = !inside[i]
		}
	}
	return mask, nil
}
