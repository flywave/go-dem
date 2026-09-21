package grits

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
)

type zscoreFilter struct {
	baseGrits
}

func init() {
	Register(FilterZScore, func() Grits { return &zscoreFilter{baseGrits{name: string(FilterZScore)}} })
}

func (f *zscoreFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	threshold := opts.Threshold
	if threshold <= 0 {
		threshold = 3.0
	}
	kSize := opts.KernelSize
	if kSize < 3 {
		kSize = 5
	}
	if kSize%2 == 0 {
		kSize++
	}

	w, h := region.XSize, region.YSize
	nd := opts.GetNoData()
	result := make([]float64, len(data))
	copy(result, data)

	if err := startFilter(opts, f.Name(), h); err != nil {
		return nil, err
	}

	half := kSize / 2

	var sum, count float64
	for _, v := range data {
		if v == nd || math.IsNaN(v) {
			continue
		}
		sum += v
		count++
	}
	if count == 0 {
		finishFilter(opts, f.Name(), h)
		return result, nil
	}
	offset := sum / count

	spw := w + 1
	sat := make([]float64, spw*(h+1))
	sqs := make([]float64, spw*(h+1))
	cnt := make([]int32, spw*(h+1))
	for y := 0; y < h; y++ {
		if err := dem.CheckCtx(opts.Ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", f.Name(), err)
		}
		rowSat, rowSqs := 0.0, 0.0
		var rowCnt int32
		for x := 0; x < w; x++ {
			v := data[y*w+x]
			if v != nd && !math.IsNaN(v) {
				d := v - offset
				rowSat += d
				rowSqs += d * d
				rowCnt++
			}
			sat[(y+1)*spw+x+1] = sat[y*spw+x+1] + rowSat
			sqs[(y+1)*spw+x+1] = sqs[y*spw+x+1] + rowSqs
			cnt[(y+1)*spw+x+1] = cnt[y*spw+x+1] + rowCnt
		}
		dem.ReportProgress(opts.Progress, f.Name(), (y+1)/2, h)
	}

	for y := 0; y < h; y++ {
		if err := dem.CheckCtx(opts.Ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", f.Name(), err)
		}
		wy0 := y - half
		if wy0 < 0 {
			wy0 = 0
		}
		wy1 := y + half + 1
		if wy1 > h {
			wy1 = h
		}
		a0 := wy0 * spw
		a1 := wy1 * spw
		for x := 0; x < w; x++ {
			idx := y*w + x
			if result[idx] == nd || math.IsNaN(result[idx]) {
				continue
			}
			wx0 := x - half
			if wx0 < 0 {
				wx0 = 0
			}
			wx1 := x + half + 1
			if wx1 > w {
				wx1 = w
			}
			c := int(cnt[a1+wx1]) - int(cnt[a0+wx1]) - int(cnt[a1+wx0]) + int(cnt[a0+wx0])
			if c < 2 {
				continue
			}
			cf := float64(c)
			s := sat[a1+wx1] - sat[a0+wx1] - sat[a1+wx0] + sat[a0+wx0]
			ms := s / cf
			q := sqs[a1+wx1] - sqs[a0+wx1] - sqs[a1+wx0] + sqs[a0+wx0]
			variance := q/cf - ms*ms
			if variance < 1e-12 {
				continue
			}
			std := math.Sqrt(variance)
			mean := offset + ms
			z := math.Abs((result[idx] - mean) / std)
			if z > threshold {
				result[idx] = nd
			}
		}
		dem.ReportProgress(opts.Progress, f.Name(), (h+y+1)/2, h)
	}

	finishFilter(opts, f.Name(), h)
	return result, nil
}
