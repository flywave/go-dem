package pipeline

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	gdal "github.com/flywave/flywave-gdal"
	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/datum"
	"github.com/flywave/go-dem/waffle"
	"github.com/flywave/go-geo"
	"github.com/flywave/go-geoid"
	"github.com/flywave/go3d/float64/vec2"
)

const testElevation = 100.0

func testRegion() *dem.Region {
	return dem.NewRegionFromBBox(-122.5, 40.0, -122.4, 40.1, geo.NewProj("EPSG:4326"), 0.05, 0.05)
}

func testPoints() []waffle.Point {
	return []waffle.Point{
		{Position: vec2.T{-122.475, 40.075}, Z: testElevation},
		{Position: vec2.T{-122.425, 40.075}, Z: testElevation},
		{Position: vec2.T{-122.475, 40.025}, Z: testElevation},
		{Position: vec2.T{-122.425, 40.025}, Z: testElevation},
	}
}

func pixelLonLat(region *dem.Region, x, y int) (float64, float64) {
	bbox := region.BBox()
	return bbox.Min[0] + (float64(x)+0.5)*region.XRes,
		bbox.Max[1] - (float64(y)+0.5)*region.YRes
}

func readOutputDEM(t *testing.T, path string, region *dem.Region) []float64 {
	t.Helper()
	data, outRegion, err := dem.ReadDEM(path)
	if err != nil {
		t.Fatalf("read output %s: %v", path, err)
	}
	if outRegion.XSize != region.XSize || outRegion.YSize != region.YSize {
		t.Fatalf("output region is %dx%d, want %dx%d", outRegion.XSize, outRegion.YSize, region.XSize, region.YSize)
	}
	return data
}

func assertNoOutputFile(t *testing.T, outPath string) {
	t.Helper()
	if _, statErr := os.Stat(outPath); !os.IsNotExist(statErr) {
		t.Errorf("failed run must not write an output file, stat err=%v", statErr)
	}
}

func TestRunDEMUnknownSourceEpsgFailsExplicitly(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "out.tif")
	cfg := Config{SourceEpsg: 4326, VerticalDatum: geoid.EGM96, NoData: -9999}
	err := RunDEM(testPoints(), testRegion(), dem.MethodIDW, nil, outPath, cfg)
	if err == nil {
		t.Fatal("SourceEpsg=4326 is not a registered vertical frame and must fail explicitly instead of silently skipping the correction")
	}
	if !strings.Contains(err.Error(), "registered") {
		t.Errorf("error must explain the registered-frame requirement, got: %v", err)
	}
	assertNoOutputFile(t, outPath)
}

func TestRunDEMVerticalDatumCorrectionMatchesGeoid(t *testing.T) {
	region := testRegion()
	outPath := filepath.Join(t.TempDir(), "out.tif")
	cfg := Config{SourceEpsg: 5773, VerticalDatum: geoid.EGM2008, NoData: -9999}
	if err := RunDEM(testPoints(), region, dem.MethodIDW, nil, outPath, cfg); err != nil {
		t.Fatal(err)
	}

	data := readOutputDEM(t, outPath, region)
	maxWantDelta := 0.0
	maxErr := 0.0
	for y := 0; y < region.YSize; y++ {
		for x := 0; x < region.XSize; x++ {
			lon, lat := pixelLonLat(region, x, y)
			n96 := datum.WGS84ToMSL(lon, lat, testElevation, geoid.EGM96)
			n2008 := datum.WGS84ToMSL(lon, lat, testElevation, geoid.EGM2008)
			want := testElevation + n2008 - n96
			got := data[y*region.XSize+x]
			if err := math.Abs(got - want); err > 1e-6 {
				t.Errorf("pixel(%d,%d)=%f, want %f (= h−(N2008−N96) from WGS84ToMSL)", x, y, got, want)
			}
			if d := math.Abs(want - testElevation); d > maxWantDelta {
				maxWantDelta = d
			}
			if d := math.Abs(got - want); d > maxErr {
				maxErr = d
			}
		}
	}
	if maxWantDelta < 0.01 {
		t.Errorf("test region must have a meaningful EGM96/EGM2008 model difference, max is only %f", maxWantDelta)
	}
	if maxErr >= 1e-6 {
		t.Errorf("output must match the geoid model difference per pixel, max deviation %f", maxErr)
	}
}

func TestRunDEM7912EllipsoidSourceNoSilentPass(t *testing.T) {
	region := testRegion()
	outPath := filepath.Join(t.TempDir(), "out.tif")
	cfg := Config{SourceEpsg: 7912, VerticalDatum: geoid.EGM96, NoData: -9999}
	err := RunDEM(testPoints(), region, dem.MethodIDW, nil, outPath, cfg)
	if err != nil {
		if !strings.Contains(err.Error(), "registered") {
			t.Errorf("explicit failure must explain the registered-frame requirement, got: %v", err)
		}
		assertNoOutputFile(t, outPath)
		return
	}

	data := readOutputDEM(t, outPath, region)
	maxErr := 0.0
	for y := 0; y < region.YSize; y++ {
		for x := 0; x < region.XSize; x++ {
			lon, lat := pixelLonLat(region, x, y)
			n96 := datum.WGS84ToMSL(lon, lat, testElevation, geoid.EGM96)
			want := testElevation - n96
			if math.Abs(want-testElevation) < 1 {
				t.Fatalf("test region must have a meaningful geoid undulation, N=%f", n96)
			}
			if d := math.Abs(data[y*region.XSize+x] - want); d > maxErr {
				maxErr = d
			}
		}
	}
	if maxErr > 1e-3 {
		t.Errorf("successful 7912→EGM96 run must apply the geoid undulation per pixel, max deviation %f", maxErr)
	}
}

func TestRunDEMHorizontalTargetEpsgRejected(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "out.tif")
	cfg := Config{SourceEpsg: 4326, TargetEpsg: 32610}
	err := RunDEM(testPoints(), testRegion(), dem.MethodIDW, nil, outPath, cfg)
	if err == nil {
		t.Fatal("horizontal TargetEpsg must fail explicitly: datum cannot horizontally reproject")
	}
	if !strings.Contains(err.Error(), "horizontal reprojection") {
		t.Errorf("error must mention horizontal reprojection, got: %v", err)
	}
	assertNoOutputFile(t, outPath)
}

func TestRunDEMPlainRunAndSameTargetEpsg(t *testing.T) {
	for _, cfg := range []Config{
		{SourceEpsg: 4326},
		{SourceEpsg: 4326, TargetEpsg: 4326},
	} {
		outPath := filepath.Join(t.TempDir(), "out.tif")
		if err := RunDEM(testPoints(), testRegion(), dem.MethodIDW, nil, outPath, cfg); err != nil {
			t.Fatalf("cfg %+v: %v", cfg, err)
		}
		data := readOutputDEM(t, outPath, testRegion())
		for i, v := range data {
			if math.Abs(v-testElevation) > 1e-9 {
				t.Fatalf("cfg %+v: constant input must interpolate to %f, got %f at %d", cfg, testElevation, v, i)
			}
		}
	}
}

func TestRunDEMTargetEpsgMergedIntoSingleCorrection(t *testing.T) {
	region := testRegion()
	outPath := filepath.Join(t.TempDir(), "out.tif")
	cfg := Config{SourceEpsg: 5773, VerticalDatum: geoid.EGM2008, TargetEpsg: 3855}
	if err := RunDEM(testPoints(), region, dem.MethodIDW, nil, outPath, cfg); err != nil {
		t.Fatal(err)
	}

	data := readOutputDEM(t, outPath, region)
	for y := 0; y < region.YSize; y++ {
		for x := 0; x < region.XSize; x++ {
			lon, lat := pixelLonLat(region, x, y)
			n96 := datum.WGS84ToMSL(lon, lat, testElevation, geoid.EGM96)
			n2008 := datum.WGS84ToMSL(lon, lat, testElevation, geoid.EGM2008)
			single := testElevation + n2008 - n96
			double := testElevation + 2*(n2008-n96)
			got := data[y*region.XSize+x]
			if math.Abs(got-single) > 1e-6 {
				t.Errorf("pixel(%d,%d)=%f, want single correction %f (double correction would be %f)", x, y, got, single, double)
			}
		}
	}
}

func TestRunDEMNoDataLinkage(t *testing.T) {
	cases := []struct {
		name   string
		cfg    Config
		opts   *waffle.Options
		expect float64
	}{
		{"cfg nodata", Config{SourceEpsg: 4326, NoData: -32767}, nil, -32767},
		{"opts nodata", Config{SourceEpsg: 4326}, &waffle.Options{NoData: -32767}, -32767},
		{"default", Config{SourceEpsg: 4326}, nil, dem.DefaultNoData},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outPath := filepath.Join(t.TempDir(), "out.tif")
			if err := RunDEM(testPoints(), testRegion(), dem.MethodIDW, tc.opts, outPath, tc.cfg); err != nil {
				t.Fatal(err)
			}
			var nd float64
			var valid bool
			if err := gdal.WithDatasetReadonly(outPath, func(ds gdal.Dataset) error {
				nd, valid = ds.RasterBand(1).NoDataValue()
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			if !valid || nd != tc.expect {
				t.Errorf("output nodata label = %f (valid=%v), want %f", nd, valid, tc.expect)
			}
		})
	}
}
