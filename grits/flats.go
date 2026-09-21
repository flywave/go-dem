package grits

import (
	"fmt"
	"math"
	"sort"

	"github.com/flywave/go-dem"
)

type flatsFilter struct {
	baseGrits
}

func init() {
	Register(FilterFlats, func() Grits { return &flatsFilter{baseGrits{name: string(FilterFlats)}} })
}

const flatsQuantStep = 0.001

func (f *flatsFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	nd := opts.GetNoData()
	result := make([]float64, len(data))
	copy(result, data)

	quantKey := func(v float64) int64 { return int64(math.Round(v / flatsQuantStep)) }

	if err := startFilter(opts, f.Name(), region.YSize); err != nil {
		return nil, err
	}

	n := len(data)
	stride := n / 100
	if stride < 1 {
		stride = 1
	}
	h := region.YSize
	totalRows := 2 * n

	counts := make(map[int64]int)
	for i, v := range data {
		if i%stride == 0 {
			if err := dem.CheckCtx(opts.Ctx); err != nil {
				return nil, fmt.Errorf("%s: %w", f.Name(), err)
			}
			dem.ReportProgress(opts.Progress, f.Name(), i*h/totalRows, h)
		}
		if v == nd || math.IsNaN(v) {
			continue
		}
		counts[quantKey(v)]++
	}

	threshold := int(opts.Threshold)
	if threshold <= 0 {
		threshold = autoThreshold(counts)
	}

	for i, v := range result {
		if i%stride == 0 {
			if err := dem.CheckCtx(opts.Ctx); err != nil {
				return nil, fmt.Errorf("%s: %w", f.Name(), err)
			}
			dem.ReportProgress(opts.Progress, f.Name(), (n+i)*h/totalRows, h)
		}
		if v == nd || math.IsNaN(v) {
			continue
		}
		if counts[quantKey(v)] > threshold {
			result[i] = nd
		}
	}

	finishFilter(opts, f.Name(), h)
	return result, nil
}

func autoThreshold(counts map[int64]int) int {
	if len(counts) < 2 {
		return 100
	}

	freqs := make([]int, 0, len(counts))
	for _, c := range counts {
		freqs = append(freqs, c)
	}
	sort.Ints(freqs)

	idx := int(math.Round(0.99 * float64(len(freqs)-1)))
	if idx >= len(freqs) {
		idx = len(freqs) - 1
	}
	return freqs[idx]
}
