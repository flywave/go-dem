package perspecto

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func ColorbarPNG(cmap []ColorStop, opts *ColorbarOptions) *image.RGBA {
	width := opts.Width
	if width <= 0 {
		width = 600
	}
	height := opts.Height
	if height <= 0 {
		height = 60
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	barTop := 10
	barBottom := height - 20
	barLeft := 40
	barRight := width - 10
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	if len(cmap) == 0 {
		return img
	}

	minV := cmap[0].Value
	maxV := cmap[len(cmap)-1].Value
	drange := maxV - minV
	if drange == 0 {
		drange = 1
	}

	for x := barLeft; x < barRight; x++ {
		f := float64(x-barLeft) / float64(barRight-barLeft)
		val := minV + f*drange
		r, g, b := InterpolateColor(cmap, val)
		for y := barTop; y < barBottom; y++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	for x := barLeft; x <= barRight; x++ {
		img.Set(x, barTop, color.RGBA{0, 0, 0, 255})
		img.Set(x, barBottom-1, color.RGBA{0, 0, 0, 255})
	}
	for y := barTop; y < barBottom; y++ {
		img.Set(barLeft, y, color.RGBA{0, 0, 0, 255})
		img.Set(barRight-1, y, color.RGBA{0, 0, 0, 255})
	}

	return img
}

func WriteColorbarPNG(cmap []ColorStop, path string, opts *ColorbarOptions) error {
	img := ColorbarPNG(cmap, opts)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
