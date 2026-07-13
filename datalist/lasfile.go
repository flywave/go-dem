package datalist

import (
	"fmt"

	"github.com/flywave/flywave-pointcloud"
	"github.com/flywave/go-geo"
	"github.com/flywave/go3d/float64/vec2"
)

type LASFile struct {
	Path        string
	PointCount  int64
	Bounds      vec2.Rect
	SRS         string
}

func OpenLAS(path string) (*LASFile, error) {
	meta, err := pointcloud.GetPointCloudMetadata(path)
	if err != nil {
		return nil, fmt.Errorf("read LAS metadata: %v", err)
	}

	lf := &LASFile{
		Path: path,
	}

	if meta != nil {
		lf.Bounds = vec2.Rect{
			Min: vec2.T{meta.Bounds[0], meta.Bounds[1]},
			Max: vec2.T{meta.Bounds[3], meta.Bounds[4]},
		}
		lf.SRS = meta.SpatialReference
	}

	return lf, nil
}

func (lf *LASFile) Reproject(dstPath, dstSRS string) error {
	return pointcloud.Reproject(pointcloud.ReprojectOptions{
		InputPath:  lf.Path,
		OutputPath: dstPath,
		InSRS:      lf.SRS,
		OutSRS:     dstSRS,
	})
}

func (lf *LASFile) MergeWith(outputPath string, others ...string) error {
	paths := append([]string{lf.Path}, others...)
	return pointcloud.MergePointCloudsWithPaths(paths, outputPath)
}

func (lf *LASFile) BBoxString() string {
	return fmt.Sprintf("%.6f/%.6f/%.6f/%.6f",
		lf.Bounds.Min[0], lf.Bounds.Max[0],
		lf.Bounds.Min[1], lf.Bounds.Max[1])
}

func (lf *LASFile) GetSRS() geo.Proj {
	if lf.SRS != "" {
		return geo.NewProj(lf.SRS)
	}
	return geo.NewProj("EPSG:4326")
}
