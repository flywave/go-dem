package pointz

import (
	"fmt"
	"math"

	gdal "github.com/flywave/flywave-gdal"
)

type RasterMaskOptions struct {
	MaskPath string
	Invert   bool
}

func RasterMaskFilter(points []Point3D, opts *RasterMaskOptions) ([]bool, error) {
	if len(points) == 0 {
		return nil, nil
	}
	if opts == nil || opts.MaskPath == "" {
		return make([]bool, len(points)), nil
	}

	ds, err := gdal.Open(opts.MaskPath, gdal.ReadOnly)
	if err != nil {
		return nil, fmt.Errorf("open raster mask: %v", err)
	}
	defer ds.Close()

	xSize := ds.RasterXSize()
	ySize := ds.RasterYSize()
	gt := ds.GeoTransform()
	band := ds.RasterBand(1)
	ndv, _ := band.NoDataValue()

	data, err := band.ReadWindow(0, 0, xSize, ySize, xSize, ySize, gdal.Nearest)
	if err != nil {
		return nil, fmt.Errorf("read raster mask: %v", err)
	}

	mask := make([]bool, len(points))
	for i, p := range points {
		px := int(math.Floor((p.X - gt[0]) / gt[1]))
		py := int(math.Floor((p.Y - gt[3]) / gt[5]))
		if px < 0 || px >= xSize || py < 0 || py >= ySize {
			if !opts.Invert {
				mask[i] = true
			}
			continue
		}
		val := data[py*xSize+px]
		isInside := val != ndv && !math.IsNaN(val) && val != 0
		if opts.Invert {
			mask[i] = isInside
		} else {
			mask[i] = !isInside
		}
	}

	return mask, nil
}
