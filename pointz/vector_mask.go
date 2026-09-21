package pointz

import (
	"fmt"
	"os"

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

	if _, err := os.Stat(opts.MaskPath); err != nil {
		return nil, fmt.Errorf("vector mask not accessible: %s: %v", opts.MaskPath, err)
	}

	ds := gdal.OpenDataSource(opts.MaskPath, 0)
	if ds == (gdal.DataSource{}) {
		return nil, fmt.Errorf("open vector mask failed: %s", opts.MaskPath)
	}
	defer ds.Destroy()

	inside := make([]bool, len(points))

	for li := 0; li < ds.LayerCount(); li++ {
		if err := applyVectorLayer(points, inside, ds.LayerByIndex(li)); err != nil {
			return nil, err
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

func applyVectorLayer(points []Point3D, inside []bool, layer gdal.Layer) error {
	layer.ResetReading()
	srs := layer.SpatialReference()

	ptGeoms := make([]gdal.Geometry, len(points))
	for i, p := range points {
		g, err := gdal.CreateFromWKT(fmt.Sprintf("POINT (%.10f %.10f)", p.X, p.Y), srs)
		if err != nil {
			for j := 0; j < i; j++ {
				ptGeoms[j].Destroy()
			}
			return fmt.Errorf("create point geometry: %v", err)
		}
		ptGeoms[i] = g
	}
	defer func() {
		for i := range ptGeoms {
			ptGeoms[i].Destroy()
		}
	}()

	for {
		feat := layer.NextFeature()
		if feat == nil {
			break
		}
		geom := feat.Geometry()
		if geom.IsNull() || geom.IsEmpty() {
			feat.Destroy()
			continue
		}

		env := geom.Envelope()
		hasEnv := env.IsInit()
		for i, p := range points {
			if inside[i] {
				continue
			}
			if hasEnv && (p.X < env.MinX() || p.X > env.MaxX() || p.Y < env.MinY() || p.Y > env.MaxY()) {
				continue
			}
			if geom.Contains(ptGeoms[i]) {
				inside[i] = true
			}
		}
		feat.Destroy()
	}
	return nil
}
