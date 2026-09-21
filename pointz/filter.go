package pointz

import (
	"fmt"

	"github.com/flywave/flywave-pointcloud"
)

type OutlierMethod string

const (
	MethodStatistical OutlierMethod = "statistical"
	MethodRadius      OutlierMethod = "radius"
)

type FilterOptions struct {
	Method     OutlierMethod
	MeanK      int
	Multiplier float64
	Radius     float64
	MinK       int
	InputPath  string
	OutputPath string
}

// MarkOutliers flags outliers in the output cloud (classification code 7)
// without removing the points.
func MarkOutliers(opts *FilterOptions) error {
	if opts == nil {
		return fmt.Errorf("filter options are required")
	}
	switch opts.Method {
	case MethodStatistical:
		meanK := opts.MeanK
		if meanK <= 0 {
			meanK = 8
		}
		multiplier := opts.Multiplier
		if multiplier <= 0 {
			multiplier = 2.0
		}
		return pointcloud.StatisticalOutlier(opts.InputPath, opts.OutputPath, meanK, multiplier)

	case MethodRadius:
		radius := opts.Radius
		if radius <= 0 {
			radius = 1.0
		}
		minK := opts.MinK
		if minK <= 0 {
			minK = 3
		}
		return pointcloud.RadiusOutlier(opts.InputPath, opts.OutputPath, radius, minK)

	default:
		return fmt.Errorf("unknown outlier method: %s", opts.Method)
	}
}

// RemoveOutliers removes outlier points from the output cloud.
func RemoveOutliers(opts *FilterOptions) error {
	if opts == nil {
		return fmt.Errorf("filter options are required")
	}
	switch opts.Method {
	case MethodStatistical:
		meanK := opts.MeanK
		if meanK <= 0 {
			meanK = 8
		}
		multiplier := opts.Multiplier
		if multiplier <= 0 {
			multiplier = 2.0
		}
		return pointcloud.RemoveStatisticalOutlier(opts.InputPath, opts.OutputPath, meanK, multiplier)

	case MethodRadius:
		radius := opts.Radius
		if radius <= 0 {
			radius = 1.0
		}
		minK := opts.MinK
		if minK <= 0 {
			minK = 3
		}
		return pointcloud.RemoveRadiusOutlier(opts.InputPath, opts.OutputPath, radius, minK)

	default:
		return fmt.Errorf("unknown outlier method: %s", opts.Method)
	}
}

// RemoveOutliersRemove is a deprecated alias for RemoveOutliers.
func RemoveOutliersRemove(opts *FilterOptions) error {
	return RemoveOutliers(opts)
}
