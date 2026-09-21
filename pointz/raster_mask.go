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
	ndv, ndvValid := band.NoDataValue()

	det := gt[1]*gt[5] - gt[2]*gt[4]
	if det == 0 {
		return nil, fmt.Errorf("degenerate geo transform in raster mask: %s", opts.MaskPath)
	}

	bw, bh := band.BlockSize()
	if bw <= 0 {
		bw = xSize
	}
	if bh <= 0 {
		bh = ySize
	}
	bcols := (xSize + bw - 1) / bw

	type maskBlock struct {
		data      []float64
		w, x0, y0 int
	}
	cache := make(map[int]maskBlock)

	mask := make([]bool, len(points))
	for i, p := range points {
		dx := p.X - gt[0]
		dy := p.Y - gt[3]
		px := int(math.Floor((gt[5]*dx - gt[2]*dy) / det))
		py := int(math.Floor((gt[1]*dy - gt[4]*dx) / det))
		if px < 0 || px >= xSize || py < 0 || py >= ySize {
			if !opts.Invert {
				mask[i] = true
			}
			continue
		}
		bk := (py/bh)*bcols + px/bw
		blk, ok := cache[bk]
		if !ok {
			x0 := (px / bw) * bw
			y0 := (py / bh) * bh
			ww := bw
			if x0+ww > xSize {
				ww = xSize - x0
			}
			hh := bh
			if y0+hh > ySize {
				hh = ySize - y0
			}
			blkData, err := band.ReadWindow(x0, y0, ww, hh, ww, hh, gdal.Nearest)
			if err != nil {
				return nil, fmt.Errorf("read raster mask: %v", err)
			}
			blk = maskBlock{data: blkData, w: ww, x0: x0, y0: y0}
			cache[bk] = blk
		}
		val := blk.data[(py-blk.y0)*blk.w+(px-blk.x0)]
		isInside := !math.IsNaN(val) && val != 0
		if ndvValid && val == ndv {
			isInside = false
		}
		if opts.Invert {
			mask[i] = isInside
		} else {
			mask[i] = !isInside
		}
	}

	return mask, nil
}
