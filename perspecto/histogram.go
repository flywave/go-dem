package perspecto

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sort"
)

func HistogramPNG(data []float64, opts *HistogramOptions) (*image.RGBA, error) {
	bins := opts.Bins
	if bins <= 0 {
		bins = 256
	}
	width := opts.Width
	if width <= 0 {
		width = 800
	}
	height := opts.Height
	if height <= 0 {
		height = 400
	}
	nd := opts.NoData
	if nd == 0 {
		nd = -9999
	}

	var vals []float64
	for _, v := range data {
		if v == nd || math.IsNaN(v) {
			continue
		}
		vals = append(vals, v)
	}
	if len(vals) == 0 {
		return nil, fmt.Errorf("no valid data")
	}

	sort.Float64s(vals)
	minV, maxV := vals[0], vals[len(vals)-1]
	if maxV-minV == 0 {
		maxV = minV + 1
	}

	hist := make([]int, bins)
	for _, v := range vals {
		idx := int((v - minV) / (maxV - minV) * float64(bins-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= bins {
			idx = bins - 1
		}
		hist[idx]++
	}

	maxCount := 0
	for _, c := range hist {
		if c > maxCount {
			maxCount = c
		}
	}
	if maxCount == 0 {
		maxCount = 1
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	bg := color.RGBA{255, 255, 255, 255}
	fg := color.RGBA{50, 50, 50, 255}
	bar := color.RGBA{70, 130, 180, 200}
	line := color.RGBA{200, 50, 50, 255}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bg)
		}
	}

	margin := 60
	plotW := width - 2*margin
	plotH := height - 2*margin
	plotTop := margin
	plotBottom := height - margin
	plotLeft := margin
	plotRight := width - margin

	for x := plotLeft; x <= plotRight; x++ {
		img.Set(x, plotBottom, fg)
	}
	for y := plotTop; y <= plotBottom; y++ {
		img.Set(plotLeft, y, fg)
	}

	isCDF := opts.Type == "cdf"
	if isCDF {
		total := float64(len(vals))
		cumulative := 0
		for i, c := range hist {
			x0 := plotLeft + i*plotW/bins
			x1 := plotLeft + (i+1)*plotW/bins
			cumulative += c
			frac := float64(cumulative) / total
			barH := int(frac * float64(plotH))
			barY := plotBottom - barH
			for x := x0; x < x1 && x <= plotRight; x++ {
				for y := barY; y <= plotBottom; y++ {
					img.Set(x, y, bar)
				}
			}
		}
	} else {
		barW := plotW / bins
		if barW < 1 {
			barW = 1
		}
		for i, c := range hist {
			x := plotLeft + i*barW
			barH := int(float64(c) / float64(maxCount) * float64(plotH))
			if barH == 0 && c > 0 {
				barH = 1
			}
			for dx := 0; dx < barW && x+dx <= plotRight; dx++ {
				for dy := 0; dy < barH; dy++ {
					img.Set(x+dx, plotBottom-dy, bar)
				}
			}
		}
	}

	if opts.ShowStats {
		mean := 0.0
		for _, v := range vals {
			mean += v
		}
		mean /= float64(len(vals))
		meanX := plotLeft + int((mean-minV)/(maxV-minV)*float64(plotW))
		if meanX >= plotLeft && meanX <= plotRight {
			for dy := -2; dy <= 2; dy++ {
				y := plotBottom/2 + dy
				if y >= 0 && y < height {
					img.Set(meanX, y, line)
				}
			}
		}

		median := vals[len(vals)/2]
		medX := plotLeft + int((median-minV)/(maxV-minV)*float64(plotW))
		if medX >= plotLeft && medX <= plotRight {
			for dy := -2; dy <= 2; dy++ {
				y := plotBottom*3/4 + dy
				if y >= 0 && y < height {
					img.Set(medX, y, color.RGBA{50, 180, 50, 255})
				}
			}
		}
	}

	return img, nil
}

func WriteHistogramPNG(data []float64, path string, opts *HistogramOptions) error {
	img, err := HistogramPNG(data, opts)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
