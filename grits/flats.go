package grits

import (
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

func (f *flatsFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	nd := opts.GetNoData()
	result := make([]float64, len(data))
	copy(result, data)

	counts := make(map[float64]int)
	for _, v := range data {
		if v == nd || math.IsNaN(v) {
			continue
		}
		counts[v]++
	}

	threshold := int(opts.Threshold)
	if threshold <= 0 {
		threshold = autoThreshold(counts)
	}

	for i, v := range result {
		if v == nd || math.IsNaN(v) {
			continue
		}
		if counts[v] > threshold {
			result[i] = nd
		}
	}

	return result, nil
}

func autoThreshold(counts map[float64]int) int {
	if len(counts) < 2 {
		return 100
	}

	freqs := make([]int, 0, len(counts))
	for _, c := range counts {
		freqs = append(freqs, c)
	}
	sort.Ints(freqs)

	idx := int(float64(len(freqs)-1) * 0.99)
	return freqs[idx]
}
