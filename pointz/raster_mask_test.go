package pointz

import (
	"os"
	"testing"

	gdal "github.com/flywave/flywave-gdal"
)

func makeTestRasterMask(t *testing.T) string {
	tmp, err := os.CreateTemp("", "mask_*.tif")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		t.Fatal(err)
	}

	ds := driver.Create(tmp.Name(), 10, 10, 1, gdal.Float32, nil)
	if ds == (gdal.Dataset{}) {
		t.Fatal("failed to create raster")
	}

	ds.SetGeoTransform([6]float64{0, 1, 0, 10, 0, -1})
	ds.SetProjection(`GEOGCS["WGS 84",DATUM["WGS_1984",SPHEROID["WGS 84",6378137,298.257223563]],PRIMEM["Greenwich",0],UNIT["degree",0.0174532925199433]]`)

	band := ds.RasterBand(1)
	band.SetNoDataValue(-9999)

	data := make([]float64, 100)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			if x >= 2 && x <= 7 && y >= 2 && y <= 7 {
				data[y*10+x] = 255
			} else {
				data[y*10+x] = -9999
			}
		}
	}

	err = band.IO(gdal.Write, 0, 0, 10, 10, data, 10, 10, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	ds.Close()
	return tmp.Name()
}

func TestRasterMaskFilter_Basic(t *testing.T) {
	path := makeTestRasterMask(t)
	defer os.Remove(path)

	pts := []Point3D{
		{X: 5, Y: 5, Z: 100},
		{X: 0, Y: 0, Z: 100},
		{X: 9, Y: 9, Z: 100},
	}

	mask, err := RasterMaskFilter(pts, &RasterMaskOptions{
		MaskPath: path,
		Invert:   false,
	})
	if err != nil {
		t.Fatalf("raster mask error: %v", err)
	}
	if mask[0] {
		t.Error("point inside mask should not be masked (invert=false)")
	}
	if !mask[1] || !mask[2] {
		t.Error("points outside mask should be masked (invert=false)")
	}
}

func TestRasterMaskFilter_Invert(t *testing.T) {
	path := makeTestRasterMask(t)
	defer os.Remove(path)

	pts := []Point3D{
		{X: 5, Y: 5, Z: 100},
		{X: 0, Y: 0, Z: 100},
	}

	mask, err := RasterMaskFilter(pts, &RasterMaskOptions{
		MaskPath: path,
		Invert:   true,
	})
	if err != nil {
		t.Fatalf("raster mask error: %v", err)
	}
	if !mask[0] {
		t.Error("invert: point inside mask should be masked")
	}
	if mask[1] {
		t.Error("invert: point outside mask should not be masked")
	}
}

func TestRasterMaskFilter_NilOpts(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0}}
	mask, err := RasterMaskFilter(pts, nil)
	if err != nil {
		t.Fatal(err)
	}
	if mask[0] {
		t.Error("nil opts should not mask")
	}
}

func TestRasterMaskFilter_EmptyPoints(t *testing.T) {
	mask, err := RasterMaskFilter(nil, &RasterMaskOptions{MaskPath: "test.tif"})
	if err != nil {
		t.Fatal(err)
	}
	if mask != nil {
		t.Error("nil points should return nil")
	}
}
