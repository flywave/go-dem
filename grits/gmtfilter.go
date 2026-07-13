package grits

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/dem/gmt"
)

type GMTFilterType string

const (
	GMTBoxcar    GMTFilterType = "b"
	GMTCosine    GMTFilterType = "c"
	GMTGaussian  GMTFilterType = "g"
	GTMMedian    GMTFilterType = "m"
	GMTMaxLike   GMTFilterType = "p"
	GTMLower     GMTFilterType = "l"
	GMTUpper     GMTFilterType = "u"
)

type gmtFilter struct {
	baseGrits
}

func init() {
	Register("gmt_filter", func() Grits {
		return &gmtFilter{baseGrits{name: "gmt_filter"}}
	})
}

func (f *gmtFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	tmpDir, err := os.MkdirTemp("", "gmt_filter_*")
	if err != nil {
		return nil, fmt.Errorf("temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	inPath := filepath.Join(tmpDir, "input.tif")
	outPath := filepath.Join(tmpDir, "output.tif")

	if err := dem.CreateDEM(data, region, inPath, opts.GetNoData()); err != nil {
		return nil, fmt.Errorf("write input: %v", err)
	}

	filterType := fmt.Sprintf("c%.0f", opts.Threshold)
	if opts.Threshold <= 0 {
		filterType = "c100"
	}

	distFlag := fmt.Sprintf("%.10f", region.XRes)

	if err := gmt.Grdfilter(inPath, outPath, filterType, distFlag); err != nil {
		return nil, fmt.Errorf("gmt grdfilter: %v", err)
	}

	result, _, err := dem.ReadDEM(outPath)
	if err != nil {
		return nil, fmt.Errorf("read result: %v", err)
	}

	return result, nil
}
