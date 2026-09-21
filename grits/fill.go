package grits

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/flywave/flywave-gdal"
	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/waffle"
	"github.com/flywave/go3d/float64/vec2"
)

type fillFilter struct{ baseGrits }

func init() {
	Register(FilterFill, func() Grits { return &fillFilter{baseGrits{name: string(FilterFill)}} })
}

func (f *fillFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	maxDist := opts.MaxDistance
	if maxDist <= 0 {
		maxDist = 100
	}
	noData := opts.GetNoData()
	width := region.XSize
	height := region.YSize

	result := make([]float64, len(data))
	copy(result, data)

	if err := startFilter(opts, f.Name(), height); err != nil {
		return nil, err
	}
	fillNoDataPixels(result, width, height, noData)
	if err := fillWithInverseDistanceOpts(result, width, height, noData, opts, f.Name()); err != nil {
		return nil, err
	}
	if err := dem.CheckCtx(opts.Ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", f.Name(), err)
	}

	if hasRemainingNoData(result, noData) {
		if filled, err := gdalFillNoData(result, region, maxDist, noData); err == nil {
			result = filled
		}
	}

	finishFilter(opts, f.Name(), height)
	return result, nil
}

func hasRemainingNoData(data []float64, noData float64) bool {
	for _, v := range data {
		if v == noData || math.IsNaN(v) {
			return true
		}
	}
	return false
}

func gdalFillNoData(data []float64, region *dem.Region, maxDist float64, noData float64) ([]float64, error) {
	tmpDir, err := os.MkdirTemp("", "gdal_fill_*")
	if err != nil {
		return data, fmt.Errorf("fill temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.tif")
	tmp := make([]float64, len(data))
	for i, v := range data {
		if math.IsNaN(v) {
			tmp[i] = noData
		} else {
			tmp[i] = v
		}
	}
	if err := dem.CreateDEM(tmp, region, inputPath, noData); err != nil {
		return data, fmt.Errorf("fill temp dem: %v", err)
	}

	err = gdal.WithDatasetUpdate(inputPath, func(ds gdal.Dataset) error {
		band := ds.RasterBand(1)
		return band.FillWithAutoMask(maxDist, 0)
	})
	if err != nil {
		return data, fmt.Errorf("gdal fillnodata: %v", err)
	}

	outputData, _, err := dem.ReadDEM(inputPath)
	if err != nil {
		return data, fmt.Errorf("fill read back: %v", err)
	}

	return outputData, nil
}

func fillNoDataPixels(data []float64, w, h int, noData float64) {
	type point struct{ x, y int }

	isNoData := func(v float64) bool { return v == noData || math.IsNaN(v) }

	queue := make([]point, 0)
	seeded := make([]bool, w*h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if isNoData(data[idx]) || seeded[idx] {
				continue
			}
			for dy := -1; dy <= 1 && !seeded[idx]; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx, ny := x+dx, y+dy
					if nx < 0 || nx >= w || ny < 0 || ny >= h {
						continue
					}
					if isNoData(data[ny*w+nx]) {
						seeded[idx] = true
						queue = append(queue, point{x, y})
						break
					}
				}
			}
		}
	}

	for head := 0; head < len(queue); head++ {
		ep := queue[head]
		z := data[ep.y*w+ep.x]
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				nx, ny := ep.x+dx, ep.y+dy
				if nx < 0 || nx >= w || ny < 0 || ny >= h {
					continue
				}
				nidx := ny*w + nx
				if isNoData(data[nidx]) {
					data[nidx] = z
					queue = append(queue, point{nx, ny})
				}
			}
		}
	}
}

func fillWithInverseDistance(data []float64, w, h int, noData float64) {
	fillWithInverseDistanceOpts(data, w, h, noData, nil, "")
}

func fillWithInverseDistanceOpts(data []float64, w, h int, noData float64, opts *Options, stage string) error {
	prog, ctx := progressOf(opts)
	type hole struct{ x, y int }

	var holes []hole
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] == noData || math.IsNaN(data[idx]) {
				holes = append(holes, hole{x, y})
			}
		}
	}
	if len(holes) == 0 {
		return nil
	}

	pts := make([]vec2.T, 0, w*h)
	vals := make([]float64, 0, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] != noData && !math.IsNaN(data[idx]) {
				pts = append(pts, vec2.T{float64(x), float64(y)})
				vals = append(vals, data[idx])
			}
		}
	}
	if len(pts) == 0 {
		return nil
	}

	tree := waffle.NewKDTree(pts)
	k := 8
	if k > len(pts) {
		k = len(pts)
	}

	stride := len(holes) / 100
	if stride < 1 {
		stride = 1
	}
	for hi, hp := range holes {
		if hi%stride == 0 {
			if err := dem.CheckCtx(ctx); err != nil {
				return fmt.Errorf("%s: %w", stage, err)
			}
			dem.ReportProgress(prog, stage, (hi+1)*h/len(holes), h)
		}
		idxs, dists := tree.KNN(vec2.T{float64(hp.x), float64(hp.y)}, k)
		var sumWeight, sumVal float64
		for i, pi := range idxs {
			d := dists[i]
			if d < 1e-10 {
				d = 1e-10
			}
			wt := 1.0 / (d * d)
			sumWeight += wt
			sumVal += wt * vals[pi]
		}
		if sumWeight > 0 {
			data[hp.y*w+hp.x] = sumVal / sumWeight
		}
	}
	return nil
}
