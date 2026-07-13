package dem

import (
	"math"
	"testing"

	"github.com/flywave/go-geo"
)

func TestNewRegionFromBBox(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(-180, -90, 180, 90, srs, 1.0, 1.0)
	if r == nil {
		t.Fatal("region is nil")
	}
	if r.XSize != 360 {
		t.Errorf("xsize: expected 360, got %d", r.XSize)
	}
	if r.YSize != 180 {
		t.Errorf("ysize: expected 180, got %d", r.YSize)
	}
	if r.XRes != 1.0 {
		t.Errorf("xres: expected 1.0, got %f", r.XRes)
	}
	if r.YRes != 1.0 {
		t.Errorf("yres: expected 1.0, got %f", r.YRes)
	}
}

func TestNewRegionFromString(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r, err := NewRegionFromString("-125/-122/40/43", srs, 0.001, 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if r == nil {
		t.Fatal("region is nil")
	}
	eps := 1e-6
	if math.Abs(r.BBox().Min[0]+125) > eps {
		t.Errorf("minx: expected -125, got %f", r.BBox().Min[0])
	}
	if math.Abs(r.BBox().Max[0]-(-122)) > eps {
		t.Errorf("maxx: expected -122, got %f", r.BBox().Max[0])
	}
	if math.Abs(r.XRes-0.001) > eps {
		t.Errorf("xres: expected 0.001, got %f", r.XRes)
	}
}

func TestRegionGeoTransform(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(-125, 40, -122, 43, srs, 0.01, 0.01)
	gt := r.GeoTransform()
	if math.Abs(gt[0]+125) > 1e-10 {
		t.Errorf("gt[0]: expected -125, got %f", gt[0])
	}
	if math.Abs(gt[1]-0.01) > 1e-10 {
		t.Errorf("gt[1]: expected 0.01, got %f", gt[1])
	}
	if gt[2] != 0 {
		t.Errorf("gt[2]: expected 0, got %f", gt[2])
	}
	if math.Abs(gt[3]-43) > 1e-10 {
		t.Errorf("gt[3]: expected 43, got %f", gt[3])
	}
	if gt[4] != 0 {
		t.Errorf("gt[4]: expected 0, got %f", gt[4])
	}
	if math.Abs(gt[5]+0.01) > 1e-10 {
		t.Errorf("gt[5]: expected -0.01, got %f", gt[5])
	}
}

func TestRegionTransformTo(t *testing.T) {
	srcSRS := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(-125, 40, -122, 43, srcSRS, 0.01, 0.01)
	if r == nil {
		t.Fatal("source region is nil")
	}

	r2 := r.TransformTo(srcSRS)
	if r2 == nil {
		t.Log("transform to same SRS returned nil (skipping)")
	} else {
		if math.Abs(r2.BBox().Min[0]+125) > 1e-6 {
			t.Errorf("transformed minx: expected -125, got %f", r2.BBox().Min[0])
		}
	}
}

func TestRegionZeroResolution(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(0, 0, 1, 1, srs, 0, 0)
	if r.XRes != 0 {
		t.Errorf("expected 0, got %f", r.XRes)
	}
	if r.YRes != 0 {
		t.Errorf("yres should be 0 when xres is 0, got %f", r.YRes)
	}
}

func TestRegionSizeCalculation(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(0, 0, 10, 10, srs, 2.0, 5.0)
	if r.XSize != 5 {
		t.Errorf("xsize: expected 5, got %d", r.XSize)
	}
	if r.YSize != 2 {
		t.Errorf("ysize: expected 2, got %d", r.YSize)
	}
}

func TestRegionString(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(-125, 40, -122, 43, srs, 0.5, 0.5)
	s := r.String()
	if s == "" {
		t.Error("empty string")
	}
}

func TestNewRegionFromStringInvalid(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	_, err := NewRegionFromString("invalid", srs, 0.1, 0.1)
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestIsNoData(t *testing.T) {
	if !IsNoData(DefaultNoData, DefaultNoData) {
		t.Error("should detect noData")
	}
	if IsNoData(0, DefaultNoData) {
		t.Error("0 should not be noData")
	}
	if IsNoData(100.5, DefaultNoData) {
		t.Error("100.5 should not be noData")
	}
	if !IsNoData(math.NaN(), DefaultNoData) {
		t.Error("NaN should be noData")
	}
}
