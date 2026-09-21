package dem

import (
	"math"
	"testing"

	gdal "github.com/flywave/flywave-gdal"
)

func makeTestRaster(t *testing.T, path string, w, h int, val float64) {
	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		t.Fatal(err)
	}
	ds := driver.Create(path, w, h, 1, gdal.Float32,
		[]string{"COMPRESS=DEFLATE"})
	if ds == (gdal.Dataset{}) {
		t.Fatal("create failed")
	}
	ds.SetGeoTransform([6]float64{0, 1, 0, float64(h), 0, -1})
	ds.SetProjection(`GEOGCS["WGS 84",DATUM["WGS_1984",SPHEROID["WGS 84",6378137,298.257223563]],PRIMEM["Greenwich",0],UNIT["degree",0.0174532925199433]]`)
	band := ds.RasterBand(1)
	band.SetNoDataValue(-9999)
	data := make([]float64, w*h)
	for i := range data {
		data[i] = val
	}
	band.IO(gdal.Write, 0, 0, w, h, data, w, h, 0, 0)
	ds.Close()
}

func TestCutRaster_Basic(t *testing.T) {
	src := t.TempDir() + "/src.tif"
	dst := t.TempDir() + "/cut.tif"
	makeTestRaster(t, src, 10, 10, 100)

	region := NewRegionFromBBox(2, 2, 8, 8, nil, 1, 1)
	err := CutRaster(src, dst, region)
	if err != nil {
		t.Fatalf("cut error: %v", err)
	}

	gdal.WithDatasetReadonly(dst, func(ds gdal.Dataset) error {
		if ds.RasterXSize() != 6 || ds.RasterYSize() != 6 {
			t.Errorf("expected 6x6, got %dx%d", ds.RasterXSize(), ds.RasterYSize())
		}
		return nil
	})
}

func TestCropRaster_NoNoData(t *testing.T) {
	src := t.TempDir() + "/src.tif"
	dst := t.TempDir() + "/crop.tif"
	makeTestRaster(t, src, 10, 10, 100)

	err := CropRaster(src, dst)
	if err != nil {
		t.Fatalf("crop error: %v", err)
	}

	gdal.WithDatasetReadonly(dst, func(ds gdal.Dataset) error {
		if ds.RasterXSize() != 10 || ds.RasterYSize() != 10 {
			t.Errorf("all valid: expected 10x10, got %dx%d", ds.RasterXSize(), ds.RasterYSize())
		}
		return nil
	})
}

func TestCropRaster_WithBorders(t *testing.T) {
	src := t.TempDir() + "/src.tif"
	dst := t.TempDir() + "/crop.tif"

	driver, _ := gdal.GetDriverByName("GTiff")
	ds := driver.Create(src, 10, 10, 1, gdal.Float32, []string{"COMPRESS=DEFLATE"})
	ds.SetGeoTransform([6]float64{0, 1, 0, 10, 0, -1})
	band := ds.RasterBand(1)
	band.SetNoDataValue(-9999)
	data := make([]float64, 100)
	for y := 2; y < 8; y++ {
		for x := 2; x < 8; x++ {
			data[y*10+x] = 100
		}
	}
	for i := range data {
		if data[i] == 0 {
			data[i] = -9999
		}
	}
	band.IO(gdal.Write, 0, 0, 10, 10, data, 10, 10, 0, 0)
	ds.Close()

	err := CropRaster(src, dst)
	if err != nil {
		t.Fatalf("crop error: %v", err)
	}

	gdal.WithDatasetReadonly(dst, func(ds gdal.Dataset) error {
		if ds.RasterXSize() != 6 || ds.RasterYSize() != 6 {
			t.Errorf("expected 6x6 (cropped), got %dx%d", ds.RasterXSize(), ds.RasterYSize())
		}
		return nil
	})
}

func TestCutRaster_InvalidRegion(t *testing.T) {
	src := t.TempDir() + "/src.tif"
	dst := t.TempDir() + "/cut.tif"
	makeTestRaster(t, src, 5, 5, 100)

	region := NewRegionFromBBox(100, 100, 200, 200, nil, 1, 1)
	err := CutRaster(src, dst, region)
	if err == nil {
		t.Error("expected error for region outside raster")
	}
}

func TestCropRaster_AllNoData(t *testing.T) {
	src := t.TempDir() + "/src.tif"
	dst := t.TempDir() + "/crop.tif"

	driver, _ := gdal.GetDriverByName("GTiff")
	ds := driver.Create(src, 5, 5, 1, gdal.Float32, nil)
	ds.SetGeoTransform([6]float64{0, 1, 0, 5, 0, -1})
	band := ds.RasterBand(1)
	band.SetNoDataValue(-9999)
	data := make([]float64, 25)
	for i := range data {
		data[i] = -9999
	}
	band.IO(gdal.Write, 0, 0, 5, 5, data, 5, 5, 0, 0)
	ds.Close()

	err := CropRaster(src, dst)
	if err == nil {
		t.Error("expected error for all-nodata raster")
	}
}

func TestCropRaster_NoNoDataMetadata(t *testing.T) {
	src := t.TempDir() + "/src.tif"
	dst := t.TempDir() + "/crop.tif"

	driver, _ := gdal.GetDriverByName("GTiff")
	ds := driver.Create(src, 8, 8, 1, gdal.Float32, nil)
	ds.SetGeoTransform([6]float64{0, 1, 0, 8, 0, -1})
	band := ds.RasterBand(1)
	data := make([]float64, 64)
	for i := range data {
		data[i] = 100
	}
	data[0] = -1e10
	band.IO(gdal.Write, 0, 0, 8, 8, data, 8, 8, 0, 0)
	ds.Close()

	if err := CropRaster(src, dst); err != nil {
		t.Fatalf("crop error: %v", err)
	}

	gdal.WithDatasetReadonly(dst, func(ds gdal.Dataset) error {
		band := ds.RasterBand(1)
		ndv, valid := band.NoDataValue()
		if valid {
			t.Errorf("source has no nodata metadata: output must not get a fake one, got %v", ndv)
		}
		if ds.RasterXSize() != 8 || ds.RasterYSize() != 8 {
			t.Errorf("-1e10 is valid data here: expected 8x8 untouched, got %dx%d",
				ds.RasterXSize(), ds.RasterYSize())
		}
		return nil
	})
}

func TestComputeEuclideanDistance(t *testing.T) {
	nd := -9999.0
	data := []float64{
		100, 100, 100,
		100, nd, 100,
		100, 100, 100,
	}
	dist := ComputeEuclideanDistance(data, 3, 3, nd, 1.0, 1.0)
	if math.Abs(dist[0]-math.Sqrt2) > 0.01 {
		t.Errorf("corner distance: expected ~1.41, got %.4f", dist[0])
	}
	if dist[4] != 0 {
		t.Errorf("nodata cell should have distance 0, got %.2f", dist[4])
	}
}

func TestComputeEuclideanDistance_Anisotropic(t *testing.T) {
	nd := -9999.0
	data := []float64{
		nd, 100, 100,
		100, 100, 100,
		100, 100, 100,
	}
	dist := ComputeEuclideanDistance(data, 3, 3, nd, 2.0, 5.0)
	if math.Abs(dist[1]-2.0) > 1e-9 {
		t.Errorf("horizontal step: expected 2.0, got %.4f", dist[1])
	}
	if math.Abs(dist[3]-5.0) > 1e-9 {
		t.Errorf("vertical step: expected 5.0, got %.4f", dist[3])
	}
	if math.Abs(dist[4]-math.Hypot(2.0, 5.0)) > 1e-9 {
		t.Errorf("diagonal step: expected %.4f, got %.4f", math.Hypot(2.0, 5.0), dist[4])
	}
}

func TestEuclideanMergeDEMs_Basic(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	base := make([]float64, w*h)
	other := make([]float64, w*h)
	for i := 0; i < w*h; i++ {
		base[i] = 100
		other[i] = 200
	}
	base[2*w+2] = nd

	region := NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	result, err := EuclideanMergeDEMs([][]float64{base, other}, region, nd)
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if len(result) != w*h {
		t.Errorf("expected %d results, got %d", w*h, len(result))
	}
	if math.Abs(result[2*w+2]-200) > 1e-6 {
		t.Errorf("nodata in base should be filled from other dem: expected 200, got %.2f", result[2*w+2])
	}
	if result[0] <= 150 || result[0] >= 200 {
		t.Errorf("merged at (0,0): expected blend in (150,200), got %.2f", result[0])
	}
	for i, v := range result {
		if v < 100-1e-6 || v > 200+1e-6 {
			t.Errorf("cell %d out of input value range: %.2f", i, v)
			break
		}
	}
}

func TestEuclideanMergeDEMs_Single(t *testing.T) {
	region := NewRegionFromBBox(0, 0, 3, 3, nil, 1, 1)
	data := make([]float64, 9)
	for i := range data {
		data[i] = 42
	}
	result, err := EuclideanMergeDEMs([][]float64{data}, region, -9999)
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	for _, v := range result {
		if v != 42 {
			t.Errorf("single dem: expected 42, got %.2f", v)
		}
	}
}

func TestEuclideanMergeDEMs_Empty(t *testing.T) {
	region := NewRegionFromBBox(0, 0, 5, 5, nil, 1, 1)
	_, err := EuclideanMergeDEMs([][]float64{}, region, -9999)
	if err == nil {
		t.Error("expected error for empty DEMs")
	}
}

func TestEuclideanMergeDEMs_GeographicCentralHole(t *testing.T) {
	w, h := 60, 100
	nd := -9999.0
	res := 1.0 / 3600.0
	region := NewRegionFromBBox(100, 40, 100+float64(w)*res, 40+float64(h)*res, nil, res, res)

	base := make([]float64, w*h)
	other := make([]float64, w*h)
	for i := range base {
		base[i] = 100
		other[i] = 200
	}
	hx0, hx1 := w/2-10, w/2+10
	hy0, hy1 := h/2-15, h/2+15
	hole := make([]bool, w*h)
	for y := hy0; y < hy1; y++ {
		for x := hx0; x < hx1; x++ {
			base[y*w+x] = nd
			hole[y*w+x] = true
		}
	}

	result, err := EuclideanMergeDEMs([][]float64{base, other}, region, nd)
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}

	for i := range result {
		if hole[i] {
			if math.Abs(result[i]-200) > 1e-6 {
				t.Fatalf("hole cell %d should be filled from other dem: expected 200, got %.4f", i, result[i])
			}
		} else if result[i] <= 100 || result[i] >= 200 {
			t.Fatalf("valid cell %d should stay a blend of inputs, got %.4f", i, result[i])
		}
	}
}

func TestEuclideanMergeDEMs_GeographicSingleValidUnchanged(t *testing.T) {
	w, h := 60, 100
	nd := -9999.0
	res := 1.0 / 3600.0
	region := NewRegionFromBBox(0, 0, float64(w)*res, float64(h)*res, nil, res, res)

	base := make([]float64, w*h)
	for i := range base {
		base[i] = 10 + float64(i%50)
	}
	hx0, hx1 := w/2-10, w/2+10
	hy0, hy1 := h/2-15, h/2+15
	for y := hy0; y < hy1; y++ {
		for x := hx0; x < hx1; x++ {
			base[y*w+x] = nd
		}
	}

	result, err := EuclideanMergeDEMs([][]float64{base}, region, nd)
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}

	for i := range result {
		if base[i] == nd {
			if result[i] != nd {
				t.Fatalf("hole cell %d should stay nodata, got %.4f", i, result[i])
			}
		} else if math.Abs(result[i]-base[i]) > 1e-6 {
			t.Fatalf("valid cell %d must not be rewritten: input=%.2f merged=%.4f", i, base[i], result[i])
		}
	}
}
