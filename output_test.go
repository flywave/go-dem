package dem

import (
	"testing"

	gdal "github.com/flywave/flywave-gdal"
	"github.com/flywave/go-geo"
)

func TestCreateRGB_ShortInput(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	region := NewRegionFromBBox(0, 0, 2, 2, srs, 1, 1)
	err := CreateRGB(make([]uint8, 3*2*2-1), region, t.TempDir()+"/rgb.tif")
	if err == nil {
		t.Error("expected error for pixel input shorter than 3*w*h")
	}
}

func TestReadDEM_Rotated(t *testing.T) {
	src := t.TempDir() + "/rot.tif"
	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		t.Fatal(err)
	}
	ds := driver.Create(src, 4, 4, 1, gdal.Float32, nil)
	if ds == (gdal.Dataset{}) {
		t.Fatal("create failed")
	}
	ds.SetGeoTransform([6]float64{0, 1, 0.1, 10, 0.1, -1})
	band := ds.RasterBand(1)
	data := make([]float64, 16)
	for i := range data {
		data[i] = 1
	}
	band.IO(gdal.Write, 0, 0, 4, 4, data, 4, 4, 0, 0)
	ds.Close()

	gdal.WithDatasetReadonly(src, func(ds gdal.Dataset) error {
		gt := ds.GeoTransform()
		if gt[2] == 0 && gt[4] == 0 {
			t.Skip("driver did not preserve rotation, cannot exercise the check")
		}
		return nil
	})

	_, _, err = ReadDEM(src)
	if err == nil {
		t.Error("expected error for rotated raster in ReadDEM")
	}

	_, _, err = ReadDEMBand(src, 1)
	if err == nil {
		t.Error("expected error for rotated raster in ReadDEMBand")
	}
}
