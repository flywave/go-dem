package grits

import (
	"github.com/flywave/go-dem"
)

type euclideanDistanceFilter struct {
	baseGrits
}

func init() {
	Register(FilterEuclideanDistance, func() Grits {
		return &euclideanDistanceFilter{baseGrits{name: string(FilterEuclideanDistance)}}
	})
}

func (f *euclideanDistanceFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	nd := opts.GetNoData()
	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}
	dist := dem.ComputeEuclideanDistance(data, region.XSize, region.YSize, nd, region.XRes, region.YRes)
	finishFilter(opts, f.Name(), region.YSize)
	return dist, nil
}
