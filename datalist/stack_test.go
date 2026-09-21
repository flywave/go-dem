package datalist

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
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

func fillCell(s *Stack, i int, elev, count, weight, unc, src float64) {
	s.Elevation[i] = elev
	s.Count[i] = count
	s.Weight[i] = weight
	s.Uncertainty[i] = unc
	s.SourceID[i] = src
}

func TestMerge_First(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	fillCell(a, 0, 100, 1, 0.5, 0.3, 7)
	fillCell(b, 0, 200, 2, 1.5, 0.4, 9)

	if err := a.Merge(b, StackModeFirst); err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if a.Elevation[0] != 100 {
		t.Errorf("first elevation: expected 100, got %.2f", a.Elevation[0])
	}
	if a.Count[0] != 1 || a.Weight[0] != 0.5 || a.SourceID[0] != 7 {
		t.Errorf("first cell fields changed: count=%.2f weight=%.2f src=%.2f", a.Count[0], a.Weight[0], a.SourceID[0])
	}
	if a.Uncertainty[0] != 0.3 {
		t.Errorf("first uncertainty must stay untouched: expected 0.3, got %.4f", a.Uncertainty[0])
	}
}

func TestMerge_Last(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	fillCell(a, 0, 100, 1, 0.5, 0.3, 7)
	fillCell(b, 0, 200, 2, 1.5, 0.4, 9)

	if err := a.Merge(b, StackModeLast); err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if a.Elevation[0] != 200 {
		t.Errorf("last elevation: expected 200, got %.2f", a.Elevation[0])
	}
	if a.Count[0] != 2 || a.Weight[0] != 1.5 || a.SourceID[0] != 9 {
		t.Errorf("last cell fields should come from other: count=%.2f weight=%.2f src=%.2f", a.Count[0], a.Weight[0], a.SourceID[0])
	}
	if a.Uncertainty[0] != 0.4 {
		t.Errorf("last uncertainty: expected 0.4, got %.4f", a.Uncertainty[0])
	}
}

func TestMerge_LastOtherNoData(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	fillCell(a, 0, 100, 1, 0.5, 0.3, 7)
	b.Elevation[0] = -9999

	if err := a.Merge(b, StackModeLast); err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if a.Elevation[0] != 100 {
		t.Errorf("last with nodata other: expected 100, got %.2f", a.Elevation[0])
	}
}

func TestMerge_Uncertainty_Mean(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	a.Elevation[0], a.Count[0], a.Uncertainty[0] = 100, 1, 0.3
	b.Elevation[0], b.Count[0], b.Uncertainty[0] = 200, 1, 0.4

	if err := a.Merge(b, StackModeMean); err != nil {
		t.Fatalf("merge error: %v", err)
	}
	want := math.Sqrt(0.3*0.3 + 0.4*0.4)
	if math.Abs(a.Uncertainty[0]-want) > 1e-9 {
		t.Errorf("mean uncertainty RSS: expected %.6f, got %.6f", want, a.Uncertainty[0])
	}
}

func TestMerge_Uncertainty_Weight(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	a.Elevation[0], a.Weight[0], a.Uncertainty[0] = 100, 1, 0.3
	b.Elevation[0], b.Weight[0], b.Uncertainty[0] = 200, 1, 0.4

	if err := a.Merge(b, StackModeWeight); err != nil {
		t.Fatalf("merge error: %v", err)
	}
	want := math.Sqrt(0.3*0.3 + 0.4*0.4)
	if math.Abs(a.Uncertainty[0]-want) > 1e-9 {
		t.Errorf("weight uncertainty RSS: expected %.6f, got %.6f", want, a.Uncertainty[0])
	}
	if a.Elevation[0] != 150 {
		t.Errorf("weight elevation: expected 150, got %.2f", a.Elevation[0])
	}
}

func TestMerge_Uncertainty_Min(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	a.Elevation[0], a.Uncertainty[0] = 100, 0.3
	b.Elevation[0], b.Uncertainty[0] = 50, 0.4

	if err := a.Merge(b, StackModeMin); err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if a.Elevation[0] != 50 {
		t.Errorf("min elevation: expected 50, got %.2f", a.Elevation[0])
	}
	if a.Uncertainty[0] != 0.4 {
		t.Errorf("min uncertainty should follow winning side: expected 0.4, got %.4f", a.Uncertainty[0])
	}

	c := NewStack(reg, -9999)
	c.Elevation[0], c.Uncertainty[0] = 100, 0.3
	d := NewStack(reg, -9999)
	d.Elevation[0], d.Uncertainty[0] = 200, 0.4

	if err := c.Merge(d, StackModeMin); err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if c.Uncertainty[0] != 0.3 {
		t.Errorf("min uncertainty must stay when other loses: expected 0.3, got %.4f", c.Uncertainty[0])
	}
}

func TestMerge_Uncertainty_Max(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	a.Elevation[0], a.Uncertainty[0] = 100, 0.3
	b.Elevation[0], b.Uncertainty[0] = 150, 0.4

	if err := a.Merge(b, StackModeMax); err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if a.Elevation[0] != 150 {
		t.Errorf("max elevation: expected 150, got %.2f", a.Elevation[0])
	}
	if a.Uncertainty[0] != 0.4 {
		t.Errorf("max uncertainty should follow winning side: expected 0.4, got %.4f", a.Uncertainty[0])
	}

	c := NewStack(reg, -9999)
	c.Elevation[0], c.Uncertainty[0] = 100, 0.3
	d := NewStack(reg, -9999)
	d.Elevation[0], d.Uncertainty[0] = 50, 0.4

	if err := c.Merge(d, StackModeMax); err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if c.Uncertainty[0] != 0.3 {
		t.Errorf("max uncertainty must stay when other loses: expected 0.3, got %.4f", c.Uncertainty[0])
	}
}

func TestMerge_MinMaxUpdateCountWeight(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	fillCell(a, 0, 100, 1, 0.5, 0.3, 7)
	fillCell(b, 0, 50, 3, 1.5, 0.4, 9)

	if err := a.Merge(b, StackModeMin); err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if a.Count[0] != 3 || a.Weight[0] != 1.5 || a.SourceID[0] != 9 {
		t.Errorf("min should take losing-side count/weight/source: count=%.2f weight=%.2f src=%.2f",
			a.Count[0], a.Weight[0], a.SourceID[0])
	}

	c := NewStack(reg, -9999)
	d := NewStack(reg, -9999)
	fillCell(c, 0, 100, 1, 0.5, 0.3, 7)
	fillCell(d, 0, 150, 3, 1.5, 0.4, 9)

	if err := c.Merge(d, StackModeMax); err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if c.Count[0] != 3 || c.Weight[0] != 1.5 || c.SourceID[0] != 9 {
		t.Errorf("max should take winning-side count/weight/source: count=%.2f weight=%.2f src=%.2f",
			c.Count[0], c.Weight[0], c.SourceID[0])
	}
}

func TestMerge_BBoxMismatch(t *testing.T) {
	reg1 := dem.NewRegionFromBBox(0, 0, 5, 5, nil, 1, 1)
	reg2 := dem.NewRegionFromBBox(100, 100, 105, 105, nil, 1, 1)
	a := NewStack(reg1, -9999)
	b := NewStack(reg2, -9999)

	if a.Region.XSize != b.Region.XSize || a.Region.YSize != b.Region.YSize {
		t.Fatalf("test setup: regions must share size to isolate the bbox check")
	}
	err := a.Merge(b, StackModeMean)
	if err == nil {
		t.Error("expected error for bounds mismatch with identical sizes")
	}
}

func TestMerge_SRSMismatch(t *testing.T) {
	reg1 := dem.NewRegionFromBBox(0, 0, 5, 5, nil, 1, 1)
	reg2 := dem.NewRegionFromBBox(0, 0, 5, 5, geo.NewProj("EPSG:4326"), 1, 1)
	a := NewStack(reg1, -9999)
	b := NewStack(reg2, -9999)

	if err := a.Merge(b, StackModeMean); err == nil {
		t.Error("expected error for SRS mismatch (nil vs EPSG:4326)")
	}
}

func TestMerge_SameSRS(t *testing.T) {
	reg1 := dem.NewRegionFromBBox(0, 0, 5, 5, geo.NewProj("EPSG:4326"), 1, 1)
	reg2 := dem.NewRegionFromBBox(0, 0, 5, 5, geo.NewProj("EPSG:4326"), 1, 1)
	a := NewStack(reg1, -9999)
	b := NewStack(reg2, -9999)

	a.Elevation[0] = 100
	a.Count[0] = 1
	b.Elevation[0] = 200
	b.Count[0] = 1

	if err := a.Merge(b, StackModeMean); err != nil {
		t.Fatalf("merge with equal SRS should succeed: %v", err)
	}
	if math.Abs(a.Elevation[0]-150) > 1e-10 {
		t.Errorf("merge mean: expected 150, got %.2f", a.Elevation[0])
	}
}

func TestMerge_UnknownMode(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, nil, 1, 1)
	a := NewStack(reg, -9999)
	b := NewStack(reg, -9999)

	a.Elevation[0] = 100
	b.Elevation[0] = 200

	if err := a.Merge(b, StackMode("bogus")); err == nil {
		t.Error("expected error for unknown stack mode")
	}
	if a.Elevation[0] != 100 {
		t.Errorf("unknown mode must not modify data: expected 100, got %.2f", a.Elevation[0])
	}
}

func TestBuildDataList_DirFiltersNonSourceFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"dem.tif", "points.las", ".DS_Store", "notes.bin", "README.md"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}

	dl, err := BuildDataList([]string{dir})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(dl.Entries) != 2 {
		t.Fatalf("expected 2 entries (tif+las), got %d: %v", len(dl.Entries), dl.Entries)
	}
	for _, e := range dl.Entries {
		if filepath.Ext(e.Path) == ".DS_Store" || filepath.Ext(e.Path) == ".bin" || filepath.Ext(e.Path) == ".md" {
			t.Errorf("non-source file slipped in: %s", e.Path)
		}
	}
}

func TestReadStack_RoundTrip(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, geo.NewProj("EPSG:4326"), 1, 1)
	s := NewStack(reg, -32767)
	s.Elevation[0] = 100
	s.Count[0] = 2
	s.Weight[0] = 1.5
	s.Uncertainty[0] = 0.25
	s.SourceID[0] = 3

	path := filepath.Join(t.TempDir(), "stack.tif")
	if err := s.Write(path); err != nil {
		t.Fatalf("write stack: %v", err)
	}

	got, err := ReadStack(path)
	if err != nil {
		t.Fatalf("read stack: %v", err)
	}
	if got.NoData != -32767 {
		t.Errorf("nodata from file: expected -32767, got %v", got.NoData)
	}
	if got.Elevation[0] != 100 {
		t.Errorf("elevation: expected 100, got %.2f", got.Elevation[0])
	}
	if got.Count[0] != 2 {
		t.Errorf("count: expected 2, got %.2f", got.Count[0])
	}
	if got.Weight[0] != 1.5 {
		t.Errorf("weight: expected 1.5, got %.2f", got.Weight[0])
	}
	if got.Uncertainty[0] != 0.25 {
		t.Errorf("uncertainty: expected 0.25, got %.4f", got.Uncertainty[0])
	}
	if got.SourceID[0] != 3 {
		t.Errorf("source: expected 3, got %.2f", got.SourceID[0])
	}
}

func TestReadStack_SingleBandFails(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 2, 2, geo.NewProj("EPSG:4326"), 1, 1)
	path := filepath.Join(t.TempDir(), "dem.tif")
	data := make([]float64, 4)
	for i := range data {
		data[i] = 5
	}
	if err := dem.CreateDEM(data, reg, path, -9999); err != nil {
		t.Fatalf("create dem: %v", err)
	}

	if _, err := ReadStack(path); err == nil {
		t.Error("expected error when stack bands are missing from a single-band DEM")
	}
}
