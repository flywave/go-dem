package dem

import (
	"fmt"
	"math"

	gdal "github.com/flywave/flywave-gdal"
)

func ComputeEuclideanDistance(data []float64, w, h int, nd float64, resX, resY float64) []float64 {
	dist := make([]float64, len(data))
	inf := math.MaxFloat64
	for i := range dist {
		if data[i] == nd || math.IsNaN(data[i]) {
			dist[i] = 0
		} else {
			dist[i] = inf
		}
	}

	diag := math.Hypot(resX, resY)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if dist[idx] == 0 {
				continue
			}
			if x > 0 {
				dist[idx] = math.Min(dist[idx], dist[y*w+(x-1)]+resX)
			}
			if y > 0 {
				dist[idx] = math.Min(dist[idx], dist[(y-1)*w+x]+resY)
			}
			if x > 0 && y > 0 {
				dist[idx] = math.Min(dist[idx], dist[(y-1)*w+(x-1)]+diag)
			}
			if x < w-1 && y > 0 {
				dist[idx] = math.Min(dist[idx], dist[(y-1)*w+(x+1)]+diag)
			}
		}
	}

	for y := h - 1; y >= 0; y-- {
		for x := w - 1; x >= 0; x-- {
			idx := y*w + x
			if dist[idx] == 0 {
				continue
			}
			if x < w-1 {
				dist[idx] = math.Min(dist[idx], dist[y*w+(x+1)]+resX)
			}
			if y < h-1 {
				dist[idx] = math.Min(dist[idx], dist[(y+1)*w+x]+resY)
			}
			if x < w-1 && y < h-1 {
				dist[idx] = math.Min(dist[idx], dist[(y+1)*w+(x+1)]+diag)
			}
			if x > 0 && y < h-1 {
				dist[idx] = math.Min(dist[idx], dist[(y+1)*w+(x-1)]+diag)
			}
		}
	}

	maxDist := math.Hypot(float64(w)*resX, float64(h)*resY)
	for i := range dist {
		if dist[i] >= math.MaxFloat64/2 {
			dist[i] = maxDist
		}
	}

	return dist
}

func EuclideanMergeDEMs(dems [][]float64, region *Region, nd float64) ([]float64, error) {
	if len(dems) == 0 {
		return nil, fmt.Errorf("no DEMs provided")
	}
	w, h := region.XSize, region.YSize
	n := w * h
	if n == 0 {
		return nil, fmt.Errorf("invalid region")
	}
	if region.XRes <= 0 || region.YRes <= 0 {
		return nil, fmt.Errorf("invalid region resolution: %f x %f", region.XRes, region.YRes)
	}
	for i, d := range dems {
		if len(d) != n {
			return nil, fmt.Errorf("DEM %d: length %d doesn't match region %dx%d=%d", i, len(d), w, h, n)
		}
	}

	dists := make([][]float64, len(dems))
	for i, d := range dems {
		dists[i] = ComputeEuclideanDistance(d, w, h, nd, region.XRes, region.YRes)
	}

	sum := make([]float64, n)
	weight := make([]float64, n)
	for i := 0; i < len(dems); i++ {
		for j := 0; j < n; j++ {
			if IsNoData(dems[i][j], nd) {
				continue
			}
			wgt := dists[i][j]
			sum[j] += dems[i][j] * wgt
			weight[j] += wgt
		}
	}

	result := make([]float64, n)
	for j := 0; j < n; j++ {
		if weight[j] > 0 {
			result[j] = sum[j] / weight[j]
		} else {
			result[j] = nd
		}
	}

	return result, nil
}

func CutRaster(srcPath, dstPath string, region *Region) error {
	return gdal.WithDatasetReadonly(srcPath, func(src gdal.Dataset) error {
		xSize := src.RasterXSize()
		ySize := src.RasterYSize()
		gt := src.GeoTransform()
		if gt[1] == 0 || gt[5] == 0 {
			return fmt.Errorf("degenerate geotransform in %s", srcPath)
		}

		pxMin := int(math.Floor((region.Extent.BBox.Min[0] - gt[0]) / gt[1]))
		pxMax := int(math.Ceil((region.Extent.BBox.Max[0] - gt[0]) / gt[1]))
		pyMin := int(math.Floor((region.Extent.BBox.Max[1] - gt[3]) / gt[5]))
		pyMax := int(math.Ceil((region.Extent.BBox.Min[1] - gt[3]) / gt[5]))

		if pxMin < 0 {
			pxMin = 0
		}
		if pyMin < 0 {
			pyMin = 0
		}
		if pxMax > xSize {
			pxMax = xSize
		}
		if pyMax > ySize {
			pyMax = ySize
		}
		winW := pxMax - pxMin
		winH := pyMax - pyMin
		if winW <= 0 || winH <= 0 {
			return fmt.Errorf("cut region outside raster extent")
		}

		dstGT := [6]float64{
			gt[0] + float64(pxMin)*gt[1],
			gt[1], 0,
			gt[3] + float64(pyMin)*gt[5],
			0, gt[5],
		}

		bands := src.RasterCount()
		driver, err := gdal.GetDriverByName("GTiff")
		if err != nil {
			return err
		}
		band1 := src.RasterBand(1)
		dt := band1.RasterDataType()

		dst := driver.Create(dstPath, winW, winH, bands, dt,
			[]string{"COMPRESS=DEFLATE", "TILED=YES", "BIGTIFF=IF_SAFER"})
		if dst == (gdal.Dataset{}) {
			return fmt.Errorf("failed to create %s", dstPath)
		}
		dst.SetGeoTransform(dstGT)
		if proj := src.Projection(); proj != "" {
			dst.SetProjection(proj)
		}

		for b := 1; b <= bands; b++ {
			srcBand := src.RasterBand(b)
			data, err := srcBand.ReadWindow(pxMin, pyMin, winW, winH, winW, winH, gdal.Nearest)
			if err != nil {
				dst.Close()
				return fmt.Errorf("read band %d: %v", b, err)
			}
			dstBand := dst.RasterBand(b)
			if ndv, valid := srcBand.NoDataValue(); valid {
				dstBand.SetNoDataValue(ndv)
			}
			if err := dstBand.IO(gdal.Write, 0, 0, winW, winH, data, winW, winH, 0, 0); err != nil {
				dst.Close()
				return fmt.Errorf("write band %d: %v", b, err)
			}
		}
		dst.Close()
		return nil
	})
}

func CropRaster(srcPath, dstPath string) error {
	return gdal.WithDatasetReadonly(srcPath, func(src gdal.Dataset) error {
		xSize := src.RasterXSize()
		ySize := src.RasterYSize()
		band := src.RasterBand(1)
		ndv, ndvValid := band.NoDataValue()

		data, err := band.ReadWindow(0, 0, xSize, ySize, xSize, ySize, gdal.Nearest)
		if err != nil {
			return err
		}

		firstRow, lastRow := -1, -1
		firstCol, lastCol := -1, -1

		for y := 0; y < ySize; y++ {
			for x := 0; x < xSize; x++ {
				v := data[y*xSize+x]
				isNoData := math.IsNaN(v) || (ndvValid && v == ndv)
				if !isNoData {
					if firstRow < 0 {
						firstRow = y
					}
					lastRow = y
					break
				}
			}
		}
		for x := 0; x < xSize; x++ {
			for y := 0; y < ySize; y++ {
				v := data[y*xSize+x]
				isNoData := math.IsNaN(v) || (ndvValid && v == ndv)
				if !isNoData {
					if firstCol < 0 {
						firstCol = x
					}
					break
				}
			}
		}
		for x := xSize - 1; x >= 0; x-- {
			for y := 0; y < ySize; y++ {
				v := data[y*xSize+x]
				isNoData := math.IsNaN(v) || (ndvValid && v == ndv)
				if !isNoData {
					if lastCol < 0 {
						lastCol = x
					}
					break
				}
			}
			if lastCol >= 0 {
				break
			}
		}
		for y := ySize - 1; y >= 0; y-- {
			for x := 0; x < xSize; x++ {
				v := data[y*xSize+x]
				isNoData := math.IsNaN(v) || (ndvValid && v == ndv)
				if !isNoData {
					if lastRow < 0 {
						lastRow = y
					}
					break
				}
			}
			if lastRow >= 0 {
				break
			}
		}

		if firstRow < 0 || firstCol < 0 {
			return fmt.Errorf("no valid data found in raster")
		}

		winW := lastCol - firstCol + 1
		winH := lastRow - firstRow + 1
		gt := src.GeoTransform()

		cropped := make([]float64, winW*winH)
		for y := 0; y < winH; y++ {
			copy(cropped[y*winW:(y+1)*winW],
				data[(firstRow+y)*xSize+firstCol:(firstRow+y)*xSize+firstCol+winW])
		}

		dstGT := [6]float64{
			gt[0] + float64(firstCol)*gt[1],
			gt[1], 0,
			gt[3] + float64(firstRow)*gt[5],
			0, gt[5],
		}

		driver, err := gdal.GetDriverByName("GTiff")
		if err != nil {
			return err
		}
		dt := band.RasterDataType()
		dst := driver.Create(dstPath, winW, winH, 1, dt,
			[]string{"COMPRESS=DEFLATE", "TILED=YES"})
		if dst == (gdal.Dataset{}) {
			return fmt.Errorf("failed to create %s", dstPath)
		}
		dst.SetGeoTransform(dstGT)
		if proj := src.Projection(); proj != "" {
			dst.SetProjection(proj)
		}
		dstBand := dst.RasterBand(1)
		if ndvValid && !math.IsNaN(ndv) {
			dstBand.SetNoDataValue(ndv)
		}
		if err := dstBand.IO(gdal.Write, 0, 0, winW, winH, cropped, winW, winH, 0, 0); err != nil {
			dst.Close()
			return fmt.Errorf("write cropped band: %v", err)
		}
		dst.Close()
		return nil
	})
}
