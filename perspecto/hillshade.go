package perspecto

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
)

func Hillshade(data []float64, region *dem.Region, opts *Options) []float64 {
	azimuth := opts.Azimuth
	if azimuth == 0 {
		azimuth = 315
	}
	altitude := opts.Altitude
	if altitude == 0 {
		altitude = 45
	}
	zFactor := opts.ZFactor
	if zFactor <= 0 {
		zFactor = 1.0
	}
	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	azimuthRad := azimuth * math.Pi / 180
	altitudeRad := altitude * math.Pi / 180

	w, h := region.XSize, region.YSize
	resX := region.XRes
	resY := region.YRes
	if resY <= 0 {
		resY = resX
	}

	result := make([]float64, w*h)
	for i := range result {
		if data[i] == noData || math.IsNaN(data[i]) {
			result[i] = noData
		}
	}

	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			idx := y*w + x
			z := data[idx]
			if z == noData || math.IsNaN(z) {
				continue
			}

			zNW := data[(y-1)*w+(x-1)]
			zN := data[(y-1)*w+x]
			zNE := data[(y-1)*w+(x+1)]
			zW := data[y*w+(x-1)]
			zE := data[y*w+(x+1)]
			zSW := data[(y+1)*w+(x-1)]
			zS := data[(y+1)*w+x]
			zSE := data[(y+1)*w+(x+1)]

			if hasNoData([]float64{zNW, zN, zNE, zW, zE, zSW, zS, zSE}, noData) {
				result[idx] = noData
				continue
			}

			dzdx := ((zNE + 2*zE + zSE) - (zNW + 2*zW + zSW)) / (8 * resX)
			dzdy := ((zSW + 2*zS + zSE) - (zNW + 2*zN + zNE)) / (8 * resY)

			dzdx *= zFactor
			dzdy *= zFactor

			slope := math.Atan(math.Sqrt(dzdx*dzdx + dzdy*dzdy))
			aspect := math.Atan2(dzdy, -dzdx)

			shade := math.Cos(altitudeRad)*math.Sin(slope) +
				math.Sin(altitudeRad)*math.Cos(slope)*math.Cos(azimuthRad-aspect)
			if shade < 0 {
				shade = 0
			}
			if shade > 1 {
				shade = 1
			}

			result[idx] = shade * 255
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if y == 0 || y == h-1 || x == 0 || x == w-1 {
				idx := y*w + x
				if data[idx] != noData && !math.IsNaN(data[idx]) {
					var sum, count float64
					for dy := -1; dy <= 1; dy++ {
						for dx := -1; dx <= 1; dx++ {
							nx, ny := x+dx, y+dy
							if nx < 0 || nx >= w || ny < 0 || ny >= h {
								continue
							}
							nidx := ny*w + nx
							if result[nidx] != noData && !math.IsNaN(result[nidx]) {
								sum += result[nidx]
								count++
							}
						}
					}
					if count > 0 {
						result[idx] = sum / count
					}
				}
			}
		}
	}

	return result
}

func hasNoData(vals []float64, noData float64) bool {
	for _, v := range vals {
		if v == noData || math.IsNaN(v) {
			return true
		}
	}
	return false
}

func Slope(data []float64, region *dem.Region, opts *Options) []float64 {
	zFactor := opts.ZFactor
	if zFactor <= 0 {
		zFactor = 1.0
	}
	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	w, h := region.XSize, region.YSize
	resX := region.XRes
	resY := region.YRes
	if resY <= 0 {
		resY = resX
	}

	result := make([]float64, w*h)

	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			idx := y*w + x
			z := data[idx]
			if z == noData || math.IsNaN(z) {
				result[idx] = noData
				continue
			}

			zW := data[y*w+(x-1)]
			zE := data[y*w+(x+1)]
			zN := data[(y-1)*w+x]
			zS := data[(y+1)*w+x]

			if hasNoData([]float64{zW, zE, zN, zS}, noData) {
				result[idx] = noData
				continue
			}

			dzdx := (zE - zW) / (2 * resX) * zFactor
			dzdy := (zS - zN) / (2 * resY) * zFactor

			slopeRad := math.Atan(math.Sqrt(dzdx*dzdx + dzdy*dzdy))

			if opts.SlopeUnits == "percent" {
				result[idx] = math.Sqrt(dzdx*dzdx+dzdy*dzdy) * 100
			} else {
				result[idx] = slopeRad * 180 / math.Pi
			}
		}
	}

	for y := 0; y < h; y++ {
		if data[y*w] != noData && result[y*w] == 0 {
			result[y*w] = result[y*w+1]
		}
		if data[y*w+w-1] != noData && result[y*w+w-1] == 0 {
			result[y*w+w-1] = result[y*w+w-2]
		}
	}

	return result
}

func Aspect(data []float64, region *dem.Region, opts *Options) []float64 {
	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	w, h := region.XSize, region.YSize
	resX := region.XRes

	result := make([]float64, w*h)

	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			idx := y*w + x
			z := data[idx]
			if z == noData || math.IsNaN(z) {
				result[idx] = noData
				continue
			}

			zNE := data[(y-1)*w+(x+1)]
			zSE := data[(y+1)*w+(x+1)]
			zSW := data[(y+1)*w+(x-1)]
			zNW := data[(y-1)*w+(x-1)]

			if hasNoData([]float64{zNE, zSE, zSW, zNW}, noData) {
				result[idx] = noData
				continue
			}

			dzdx := (zNE + 2*data[(y)*w+(x+1)] + zSE) -
				(zNW + 2*data[(y)*w+(x-1)] + zSW)
			dzdx /= (8 * resX)

			dzdy := (zSW + 2*data[(y+1)*w+x] + zSE) -
				(zNW + 2*data[(y-1)*w+x] + zNE)
			dzdy /= (8 * resX)

			aspect := math.Atan2(dzdx, dzdy) * 180 / math.Pi
			aspect = 90 - aspect
			if aspect < 0 {
				aspect += 360
			}

			result[idx] = aspect
		}
	}

	return result
}

func ColorRelief(data []float64, region *dem.Region, opts *Options) ([]uint8, error) {
	cmap := opts.Colormap
	if len(cmap) == 0 {
		cmap = DefaultTerrainColormap()
	}

	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	w, h := region.XSize, region.YSize
	pixels := make([]uint8, w*h*3)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			z := data[idx]

			bIdx := idx * 3
			if z == noData || math.IsNaN(z) {
				pixels[bIdx] = 0
				pixels[bIdx+1] = 0
				pixels[bIdx+2] = 0
				continue
			}

			r, g, b := interpolateColor(z, cmap)
			pixels[bIdx] = r
			pixels[bIdx+1] = g
			pixels[bIdx+2] = b
		}
	}

	return pixels, nil
}

func interpolateColor(z float64, cmap []ColorStop) (uint8, uint8, uint8) {
	if len(cmap) == 0 {
		return 128, 128, 128
	}
	if z <= cmap[0].Value {
		return cmap[0].R, cmap[0].G, cmap[0].B
	}
	if z >= cmap[len(cmap)-1].Value {
		return cmap[len(cmap)-1].R, cmap[len(cmap)-1].G, cmap[len(cmap)-1].B
	}

	for i := 0; i < len(cmap)-1; i++ {
		if z >= cmap[i].Value && z <= cmap[i+1].Value {
			t := (z - cmap[i].Value) / (cmap[i+1].Value - cmap[i].Value)
			r := uint8(float64(cmap[i].R) + t*float64(int(cmap[i+1].R)-int(cmap[i].R)))
			g := uint8(float64(cmap[i].G) + t*float64(int(cmap[i+1].G)-int(cmap[i].G)))
			b := uint8(float64(cmap[i].B) + t*float64(int(cmap[i+1].B)-int(cmap[i].B)))
			return r, g, b
		}
	}

	return cmap[len(cmap)-1].R, cmap[len(cmap)-1].G, cmap[len(cmap)-1].B
}

type ShadedReliefOptions struct {
	HillshadeOpts Options
	ColorOpts     Options
	Opacity       float64
}

func ShadedRelief(data []float64, region *dem.Region, opts *ShadedReliefOptions) ([]uint8, error) {
	hillshade := Hillshade(data, region, &opts.HillshadeOpts)
	colors, err := ColorRelief(data, region, &opts.ColorOpts)
	if err != nil {
		return nil, fmt.Errorf("color relief: %v", err)
	}

	opacity := opts.Opacity
	if opacity <= 0 {
		opacity = 0.6
	}

	w, h := region.XSize, region.YSize
	pixels := make([]uint8, w*h*3)

	noData := opts.HillshadeOpts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	for i := 0; i < w*h; i++ {
		bIdx := i * 3
		if data[i] == noData || math.IsNaN(data[i]) {
			pixels[bIdx] = 0
			pixels[bIdx+1] = 0
			pixels[bIdx+2] = 0
			continue
		}

		h := hillshade[i] / 255.0
		if h > 1 {
			h = 1
		}
		if h < 0 {
			h = 0
		}

		blend := h*opacity + (1 - opacity)
		pixels[bIdx] = uint8(float64(colors[bIdx]) * blend)
		pixels[bIdx+1] = uint8(float64(colors[bIdx+1]) * blend)
		pixels[bIdx+2] = uint8(float64(colors[bIdx+2]) * blend)
	}

	return pixels, nil
}

func WriteRGB(pixels []uint8, region *dem.Region, outputPath string) error {
	return dem.CreateRGB(pixels, region, outputPath)
}
