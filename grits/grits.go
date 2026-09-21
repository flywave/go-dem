package grits

import (
	"context"
	"fmt"
	"math"

	"github.com/flywave/go-dem"
)

type FilterType string

const (
	FilterGaussian          FilterType = "gaussian"
	FilterMedian            FilterType = "median"
	FilterBilateral         FilterType = "bilateral"
	FilterClip              FilterType = "clip"
	FilterFill              FilterType = "fill"
	FilterBlend             FilterType = "blend"
	FilterErode             FilterType = "erode"
	FilterDilate            FilterType = "dilate"
	FilterOpen              FilterType = "open"
	FilterClose             FilterType = "close"
	FilterZScore            FilterType = "zscore"
	FilterBlur              FilterType = "blur"
	FilterCut               FilterType = "cut"
	FilterDenoise           FilterType = "denoise"
	FilterFlats             FilterType = "flats"
	FilterOutliers          FilterType = "outliers"
	FilterEuclideanDistance FilterType = "euclidean_distance"
	FilterFlattenNoData     FilterType = "flatten_nodata"
)

type Options struct {
	Sigma       float64
	SigmaColor  float64
	Radius      int
	KernelSize  int
	Threshold   float64
	MaxDistance float64
	Iterations  int
	NoData      *float64
	PolygonWKT  string
	SourceMask  string
	CutBounds   []float64
	CutInvert   bool
	Method      string
	Percentile  float64
	Progress    dem.ProgressFunc
	Ctx         context.Context
}

func (o *Options) GetNoData() float64 {
	if o.NoData != nil {
		return *o.NoData
	}
	return dem.DefaultNoData
}

func (o *Options) IsNoData(val float64) bool {
	nd := o.GetNoData()
	return val == nd || math.IsNaN(val)
}

type Grits interface {
	Name() string
	Run(data []float64, region *dem.Region, opts *Options) ([]float64, error)
}

type Constructor func() Grits

var registry = make(map[FilterType]Constructor)

func Register(ft FilterType, c Constructor) {
	registry[ft] = c
}

func New(ft FilterType) (Grits, error) {
	c, ok := registry[ft]
	if !ok {
		return nil, fmt.Errorf("unknown filter: %s", ft)
	}
	return c(), nil
}

func ListFilters() []FilterType {
	fts := make([]FilterType, 0, len(registry))
	for ft := range registry {
		fts = append(fts, ft)
	}
	return fts
}

type baseGrits struct {
	name string
}

func (b *baseGrits) Name() string { return b.name }

func progressOf(opts *Options) (dem.ProgressFunc, context.Context) {
	if opts == nil {
		return nil, nil
	}
	return opts.Progress, opts.Ctx
}

func startFilter(opts *Options, name string, total int) error {
	dem.ReportProgress(opts.Progress, name, 0, total)
	if err := dem.CheckCtx(opts.Ctx); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

func finishFilter(opts *Options, name string, total int) {
	dem.ReportProgress(opts.Progress, name, total, total)
}
