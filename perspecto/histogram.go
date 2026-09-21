package perspecto

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
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

	vals := make([]float64, 0, len(data))
	minV, maxV := math.Inf(1), math.Inf(-1)
	for _, v := range data {
		if v == nd || math.IsNaN(v) {
			continue
		}
		vals = append(vals, v)
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	if len(vals) == 0 {
		return nil, fmt.Errorf("no valid data")
	}
	if opts.ShowStats {
		sort.Float64s(vals)
	}
	if maxV-minV == 0 {
		maxV = minV + 1
	}

	hist := make([]int, bins)
	for _, v := range vals {
		idx := int((v - minV) / (maxV - minV) * float64(bins))
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

	draw.Draw(img, img.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)

	stride := img.Stride

	margin := 60
	plotW := width - 2*margin
	plotH := height - 2*margin
	plotTop := margin
	plotBottom := height - margin
	plotLeft := margin
	plotRight := width - margin

	if plotBottom >= 0 {
		off := plotBottom*stride + plotLeft*4
		for x := plotLeft; x <= plotRight; x++ {
			img.Pix[off] = fg.R
			img.Pix[off+1] = fg.G
			img.Pix[off+2] = fg.B
			img.Pix[off+3] = fg.A
			off += 4
		}
	}
	if plotLeft < width {
		for y := plotTop; y <= plotBottom; y++ {
			off := y*stride + plotLeft*4
			img.Pix[off] = fg.R
			img.Pix[off+1] = fg.G
			img.Pix[off+2] = fg.B
			img.Pix[off+3] = fg.A
		}
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
				off := barY*stride + x*4
				for y := barY; y <= plotBottom; y++ {
					img.Pix[off] = bar.R
					img.Pix[off+1] = bar.G
					img.Pix[off+2] = bar.B
					img.Pix[off+3] = bar.A
					off += stride
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
				off := (plotBottom-barH+1)*stride + (x+dx)*4
				for dy := 0; dy < barH; dy++ {
					img.Pix[off] = bar.R
					img.Pix[off+1] = bar.G
					img.Pix[off+2] = bar.B
					img.Pix[off+3] = bar.A
					off += stride
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
					off := y*stride + meanX*4
					img.Pix[off] = line.R
					img.Pix[off+1] = line.G
					img.Pix[off+2] = line.B
					img.Pix[off+3] = line.A
				}
			}
		}

		median := vals[len(vals)/2]
		medX := plotLeft + int((median-minV)/(maxV-minV)*float64(plotW))
		if medX >= plotLeft && medX <= plotRight {
			for dy := -2; dy <= 2; dy++ {
				y := plotBottom*3/4 + dy
				if y >= 0 && y < height {
					off := y*stride + medX*4
					img.Pix[off] = 50
					img.Pix[off+1] = 180
					img.Pix[off+2] = 50
					img.Pix[off+3] = 255
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
