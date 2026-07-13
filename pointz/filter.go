package pointz

import (
	"fmt"

	"github.com/flywave/flywave-pointcloud"
)

type OutlierMethod string

const (
	MethodStatistical OutlierMethod = "statistical"
	MethodRadius     OutlierMethod = "radius"
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

func RemoveOutliers(opts *FilterOptions) error {
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

func RemoveOutliersRemove(opts *FilterOptions) error {
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
