package datalist

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func TestNewStack(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
	s := NewStack(reg, -9999)
	if s == nil {
		t.Fatal("nil stack")
	}
	if len(s.Elevation) != 100 {
		t.Errorf("elevation size: expected 100, got %d", len(s.Elevation))
	}
	if len(s.Count) != 100 {
		t.Errorf("count size: expected 100, got %d", len(s.Count))
	}
}

func TestStack_InitialNoData(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 5, 5, nil, 1, 1)
	s := NewStack(reg, -9999)
	for i, v := range s.Elevation {
		if v != -9999 {
			t.Errorf("elevation[%d] should be noData, got %.2f", i, v)
		}
	}
	for i, v := range s.Weight {
		if v != 0 {
			t.Errorf("weight[%d] should be 0, got %.2f", i, v)
		}
	}
}

func TestMerge_Mean(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	a.Elevation[0] = 100
	a.Count[0] = 1
	b.Elevation[0] = 200
	b.Count[0] = 1

	err := a.Merge(b, StackModeMean)
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if math.Abs(a.Elevation[0]-150) > 1e-10 {
		t.Errorf("merge mean: expected 150, got %.2f", a.Elevation[0])
	}
	if math.Abs(a.Count[0]-2) > 1e-10 {
		t.Errorf("merge count: expected 2, got %.2f", a.Count[0])
	}
}

func TestMerge_Min(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	a.Elevation[0] = 100
	b.Elevation[0] = 50

	a.Merge(b, StackModeMin)
	if math.Abs(a.Elevation[0]-50) > 1e-10 {
		t.Errorf("merge min: expected 50, got %.2f", a.Elevation[0])
	}
}

func TestMerge_Max(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	a.Elevation[0] = 100
	b.Elevation[0] = 150

	a.Merge(b, StackModeMax)
	if math.Abs(a.Elevation[0]-150) > 1e-10 {
		t.Errorf("merge max: expected 150, got %.2f", a.Elevation[0])
	}
}

func TestMerge_NoDataHandling(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	a.Elevation[0] = -9999
	b.Elevation[0] = 100

	a.Merge(b, StackModeMean)
	if math.Abs(a.Elevation[0]-100) > 1e-10 {
		t.Errorf("merge with noData: expected 100, got %.2f", a.Elevation[0])
	}
}

func TestMerge_SizeMismatch(t *testing.T) {
	reg1 := dem.NewRegionFromBBox(0, 0, 5, 5, nil, 1, 1)
	reg2 := dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
	a := NewStack(reg1, -9999)
	b := NewStack(reg2, -9999)

	err := a.Merge(b, StackModeMean)
	if err == nil {
		t.Error("expected error for size mismatch")
	}
}

func TestStackToBands(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 3, 3, nil, 1, 1)
	s := NewStack(reg, -9999)
	bands := s.ToBands()
	if len(bands) != 5 {
		t.Errorf("expected 5 bands, got %d", len(bands))
	}
	for i, b := range bands {
		if len(b) != 9 {
			t.Errorf("band %d: expected 9, got %d", i, len(b))
		}
	}
}

func TestBuildDataList_Empty(t *testing.T) {
	dl, err := BuildDataList(nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if dl == nil {
		t.Fatal("nil datalist")
	}
	if len(dl.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(dl.Entries))
	}
}

func TestDetectType(t *testing.T) {
	tests := []struct {
		path string
		want DataSourceType
	}{
		{"dem.tif", SourceRaster},
		{"dem.tiff", SourceRaster},
		{"points.las", SourcePoint},
		{"points.laz", SourcePoint},
		{"data.csv", SourcePoint},
		{"area.shp", SourceVector},
		{"area.geojson", SourceVector},
		{"unknown.xyz", SourcePoint},
	}
	for _, tt := range tests {
		got := detectType(tt.path)
		if got != tt.want {
			t.Errorf("detectType(%s) = %s, want %s", tt.path, got, tt.want)
		}
	}
}
