package datalist

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
)

func writeTempXYZ(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "points.xyz")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseXYZFile_Valid(t *testing.T) {
	path := writeTempXYZ(t, "# comment\n0 0 10 100\n1 2 20\n\nbad line\n3 4 30 200\n// note\n5 6\n")

	xf, err := ParseXYZFile(path, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(xf.Points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(xf.Points))
	}
	if xf.Points[0].Intensity != 100 || xf.Points[2].Intensity != 200 {
		t.Errorf("intensity: expected 100/200, got %.2f/%.2f", xf.Points[0].Intensity, xf.Points[2].Intensity)
	}
	if xf.Bounds.Min[0] != 0 || xf.Bounds.Min[1] != 0 || xf.Bounds.Max[0] != 3 || xf.Bounds.Max[1] != 4 {
		t.Errorf("bounds: expected 0/0/3/4, got %v/%v", xf.Bounds.Min, xf.Bounds.Max)
	}
	if xf.BBoxString() != "0.000000/3.000000/0.000000/4.000000" {
		t.Errorf("bbox string: %s", xf.BBoxString())
	}
}

func TestParseXYZFile_Empty(t *testing.T) {
	if _, err := ParseXYZFile(writeTempXYZ(t, ""), nil); err == nil {
		t.Error("expected error for empty xyz file")
	}
	if _, err := ParseXYZFile(writeTempXYZ(t, "\n# only a comment\n"), nil); err == nil {
		t.Error("expected error for comment-only xyz file")
	}
}

func TestParseXYZFile_AllInvalidLines(t *testing.T) {
	path := writeTempXYZ(t, "abc def ghi\n1 2\n3 x 5\n")
	if _, err := ParseXYZFile(path, nil); err == nil {
		t.Error("expected error when every line is invalid")
	}
}

func TestParseXYZFile_NaNInfZFiltered(t *testing.T) {
	path := writeTempXYZ(t, "0 0 NaN\n1 1 5\n2 2 Inf\n3 3 -Inf\n4 4 nan\n")
	xf, err := ParseXYZFile(path, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(xf.Points) != 1 {
		t.Fatalf("expected 1 valid point after NaN/Inf filter, got %d", len(xf.Points))
	}
	if xf.Points[0].Z != 5 {
		t.Errorf("kept wrong point: %v", xf.Points[0])
	}
	if xf.Bounds.Min[0] != 1 || xf.Bounds.Max[0] != 1 || xf.Bounds.Min[1] != 1 || xf.Bounds.Max[1] != 1 {
		t.Errorf("bounds must come from valid points only, got %v/%v", xf.Bounds.Min, xf.Bounds.Max)
	}
}

func TestParseXYZFile_MissingFile(t *testing.T) {
	if _, err := ParseXYZFile(filepath.Join(t.TempDir(), "nope.xyz"), nil); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestToDEM_IDWAggregatesSameCell(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 4, 4, nil, 1, 1)
	xf := &XYZFile{NoData: -9999}
	xf.Points = []XYZPoint{
		{X: 1.2, Y: 3.2, Z: 100},
		{X: 1.4, Y: 3.4, Z: 200},
		{X: 2.5, Y: 1.5, Z: 150},
	}

	got, err := xf.ToDEM(reg, "idw")
	if err != nil {
		t.Fatalf("toDEM idw: %v", err)
	}
	if len(got) != 16 {
		t.Fatalf("expected 16 cells, got %d", len(got))
	}
	for i, v := range got {
		if v == -9999 || math.IsNaN(v) {
			t.Fatalf("idw must fill every cell, cell %d = %v", i, v)
		}
		if v <= 100 || v >= 200 {
			t.Errorf("idw cell %d = %.4f, want strictly between min/max point z", i, v)
		}
	}
	if got[5] == 100 || got[5] == 200 {
		t.Errorf("same-cell multi-point aggregation failed: cell 5 = %.4f (last-write-wins would be 200)", got[5])
	}
}

func TestToDEM_NearestFillsGrid(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 4, 4, nil, 1, 1)
	xf := &XYZFile{NoData: -9999}
	xf.Points = []XYZPoint{
		{X: 1.0, Y: 3.0, Z: 100},
		{X: 3.0, Y: 1.0, Z: 200},
		{X: 2.0, Y: 2.0, Z: 150},
	}

	got, err := xf.ToDEM(reg, "nearest")
	if err != nil {
		t.Fatalf("toDEM nearest: %v", err)
	}
	filled := 0
	for i, v := range got {
		if v != -9999 && !math.IsNaN(v) {
			filled++
		} else {
			t.Errorf("nearest must fill every cell, cell %d = %v", i, v)
		}
	}
	if filled != 16 {
		t.Errorf("expected 16 filled cells, got %d", filled)
	}
	if got[0] != 100 {
		t.Errorf("nearest at corner (0,0): expected 100, got %.2f", got[0])
	}
	if got[15] != 200 {
		t.Errorf("nearest at corner (3,3): expected 200, got %.2f", got[15])
	}

	again, err := xf.ToDEM(reg, "nearest")
	if err != nil {
		t.Fatalf("second toDEM: %v", err)
	}
	for i := range got {
		if got[i] != again[i] {
			t.Errorf("cached points changed result at %d: %.4f vs %.4f", i, got[i], again[i])
		}
	}
}

func TestToDEM_TransformsCoordinates(t *testing.T) {
	src := geo.NewProj("EPSG:4326")
	dst := geo.NewProj("EPSG:3857")

	reg := dem.NewRegionFromBBox(1110000, 1110000, 1120000, 1125000, dst, 1000, 1000)
	xf := &XYZFile{SRS: src, NoData: -9999}
	xf.Points = []XYZPoint{{X: 10, Y: 10, Z: 50}}

	got, err := xf.ToDEM(reg, "nearest")
	if err != nil {
		t.Fatalf("toDEM with reproject: %v", err)
	}
	if len(got) != 150 {
		t.Fatalf("expected 150 cells, got %d", len(got))
	}
	if got[0] != 50 {
		t.Errorf("point must be rasterized, corner = %.2f", got[0])
	}
	if len(xf.transformed) != 1 {
		t.Fatalf("expected cached transformed points, got %d", len(xf.transformed))
	}
	tx, ty := xf.transformed[0][0], xf.transformed[0][1]
	if math.Abs(tx-1113194.91) > 1 || math.Abs(ty-1118890.48) > 1 {
		t.Errorf("batch transform 4326->3857: expected ~(1113194.91, 1118890.48), got (%.2f, %.2f)", tx, ty)
	}
	if tx < reg.BBox().Min[0] || tx > reg.BBox().Max[0] || ty < reg.BBox().Min[1] || ty > reg.BBox().Max[1] {
		t.Errorf("transformed point (%.2f, %.2f) outside region bbox", tx, ty)
	}
}

func TestToDEM_UnknownMethod(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 4, 4, nil, 1, 1)
	xf := &XYZFile{NoData: -9999}
	xf.Points = []XYZPoint{{X: 1, Y: 1, Z: 10}}

	_, err := xf.ToDEM(reg, "bogus_method")
	if err == nil {
		t.Error("expected error for unknown interpolation method")
	}
	if !strings.Contains(err.Error(), "bogus_method") {
		t.Errorf("error should name the bad method, got: %v", err)
	}
}

func TestToDEM_EmptyPoints(t *testing.T) {
	reg := dem.NewRegionFromBBox(0, 0, 4, 4, nil, 1, 1)
	if _, err := (&XYZFile{NoData: -9999}).ToDEM(reg, "idw"); err == nil {
		t.Error("expected error for xyz file without points")
	}
}
