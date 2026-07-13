package pointz

import (
	"fmt"

	"github.com/flywave/flywave-pointcloud"
	_ "github.com/flywave/flywave-pointcloud/pdal"
)

type SubsampleMethod string

const (
	MethodRandom  SubsampleMethod = "random"
	MethodSpatial SubsampleMethod = "spatial"
	MethodRadial  SubsampleMethod = "radial"
)

type SubsampleOptions struct {
	Method     SubsampleMethod
	SampleSize uint64
	VoxelSize  float64
	Radius     float64
	InputPath  string
	OutputPath string
}

func SubsamplePointCloud(opts *SubsampleOptions) error {
	method := pointcloud.RandomMethod
	switch opts.Method {
	case MethodRandom:
		method = pointcloud.RandomMethod
	case MethodSpatial:
		method = pointcloud.SpatialMethod
	case MethodRadial:
		method = pointcloud.RadialDensityMethod
	default:
		return fmt.Errorf("unknown subsample method: %s", opts.Method)
	}

	cfg := pointcloud.Metadata{
		Input:  opts.InputPath,
		Output: opts.OutputPath,
	}
	ctx := pointcloud.NewReaderContext(&cfg)
	if ctx == nil {
		return fmt.Errorf("failed to create reader context for %s", opts.InputPath)
	}

	subOpts := pointcloud.SubsampleOptions{
		Method:     method,
		SampleSize: opts.SampleSize,
		VoxelSize:  opts.VoxelSize,
		Radius:     opts.Radius,
		OutputPath: opts.OutputPath,
	}

	return pointcloud.Subsample(&ctx.ReaderContext, subOpts)
}

func VoxelDownsample(inputPath, outputPath string, voxelSize float64) error {
	return SubsamplePointCloud(&SubsampleOptions{
		Method:     MethodSpatial,
		VoxelSize:  voxelSize,
		InputPath:  inputPath,
		OutputPath: outputPath,
	})
}

func RandomDownsample(inputPath, outputPath string, sampleSize uint64) error {
	return SubsamplePointCloud(&SubsampleOptions{
		Method:     MethodRandom,
		SampleSize: sampleSize,
		InputPath:  inputPath,
		OutputPath: outputPath,
	})
}
