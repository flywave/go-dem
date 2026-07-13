package grits

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/flywave/flywave-gdal"
	"github.com/flywave/go-dem"
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

	fillNoDataPixels(result, width, height, noData)
	fillWithInverseDistance(result, width, height, noData)

	if hasRemainingNoData(result, noData) {
		if filled, err := gdalFillNoData(result, region, maxDist, noData); err == nil {
			result = filled
		}
	}

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
		return data, nil
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.tif")
	if err := dem.CreateDEM(data, region, inputPath, noData); err != nil {
		return data, fmt.Errorf("fill temp: %v", err)
	}

	err = gdal.WithDatasetUpdate(inputPath, func(ds gdal.Dataset) error {
		band := ds.RasterBand(1)
		var maskBand gdal.RasterBand
		return band.FillNodata(maskBand, maxDist, 0)
	})
	if err != nil {
		return data, nil
	}

	outputData, _, err := dem.ReadDEM(inputPath)
	if err != nil {
		return data, nil
	}

	return outputData, nil
}

func fillNoDataPixels(data []float64, w, h int, noData float64) {
	type edgePoint struct{ x, y int }
	var edges []edgePoint

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] != noData && !math.IsNaN(data[idx]) {
				continue
			}
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx, ny := x+dx, y+dy
					if nx < 0 || nx >= w || ny < 0 || ny >= h {
						continue
					}
					if data[ny*w+nx] != noData && !math.IsNaN(data[ny*w+nx]) {
						edges = append(edges, edgePoint{x, y})
						goto nextPixel
					}
				}
			}
		nextPixel:
		}
	}

	for len(edges) > 0 {
		queue := edges[:1]
		edges = edges[1:]
		idx := 0

		for idx < len(queue) {
			ep := queue[idx]
			idx++

			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx, ny := ep.x+dx, ep.y+dy
					if nx < 0 || nx >= w || ny < 0 || ny >= h {
						continue
					}
					if data[ny*w+nx] == noData || math.IsNaN(data[ny*w+nx]) {
						if data[ep.y*w+ep.x] != noData && !math.IsNaN(data[ep.y*w+ep.x]) {
							data[ny*w+nx] = data[ep.y*w+ep.x]
						}
						queue = append(queue, edgePoint{nx, ny})
					}
				}
			}
		}
	}
}

func fillWithInverseDistance(data []float64, w, h int, noData float64) {
	var validPts []struct{ x, y int; val float64 }
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] != noData && !math.IsNaN(data[idx]) {
				validPts = append(validPts, struct {
					x, y int
					val  float64
				}{x, y, data[idx]})
			}
		}
	}

	if len(validPts) == 0 {
		return
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] != noData && !math.IsNaN(data[idx]) {
				continue
			}
			var sumWeight, sumVal float64
			for _, vp := range validPts {
				dx := float64(x - vp.x)
				dy := float64(y - vp.y)
				distSq := dx*dx + dy*dy
				if distSq < 1e-10 {
					sumVal = vp.val
					sumWeight = 1
					break
				}
				w := 1.0 / distSq
				sumWeight += w
				sumVal += w * vp.val
			}
			if sumWeight > 0 {
				data[idx] = sumVal / sumWeight
			}
		}
	}
}

