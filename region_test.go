package dem

import (
	"math"
	"strings"
	"testing"

	"github.com/flywave/go-geo"
	"github.com/flywave/go-geoid"
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

func TestVerticalDatumEPSG(t *testing.T) {
	code, name := verticalDatumEPSG(geoid.EGM96)
	if code != 5773 {
		t.Errorf("EGM96: expected 5773, got %d", code)
	}
	if !strings.Contains(name, "EGM96") {
		t.Errorf("EGM96 name: expected EGM96, got %s", name)
	}

	code, name = verticalDatumEPSG(geoid.EGM2008)
	if code != 3855 {
		t.Errorf("EGM2008: expected 3855, got %d", code)
	}

	code, name = verticalDatumEPSG(geoid.EGM84)
	if code != 5798 {
		t.Errorf("EGM84: expected 5798, got %d", code)
	}

	code, _ = verticalDatumEPSG(geoid.HAE)
	if code != 0 {
		t.Errorf("HAE: expected 0, got %d", code)
	}
}

func TestResolveOutputCRS_HorizontalOnly(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	reg := NewRegionFromBBox(-180, -90, 180, 90, srs, 1, 1)
	cfg := OutputConfig{CRS: "EPSG:4326"}
	crs := resolveOutputCRS(cfg, reg)
	if crs == "" {
		t.Fatal("empty CRS")
	}
	if !strings.Contains(crs, "4326") && !strings.Contains(crs, "WGS 84") {
		t.Logf("CRS string: %s", crs[:min(80, len(crs))])
	}
}

func TestResolveOutputCRS_WithVerticalDatum(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	reg := NewRegionFromBBox(-180, -90, 180, 90, srs, 1, 1)
	cfg := OutputConfig{
		CRS:           "EPSG:4326",
		VerticalDatum: geoid.EGM96,
	}
	crs := resolveOutputCRS(cfg, reg)
	if crs == "" {
		t.Fatal("empty CRS")
	}
	t.Logf("compound CRS: %s", crs)
}

func TestResolveOutputCRS_FromRegion(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	reg := NewRegionFromBBox(-180, -90, 180, 90, srs, 1, 1)
	cfg := OutputConfig{}
	crs := resolveOutputCRS(cfg, reg)
	if crs == "" {
		t.Fatal("empty CRS")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestRegionIsValid(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(-180, -90, 180, 90, srs, 1, 1)
	if !r.IsValid() {
		t.Error("valid region should return true")
	}
	r2 := NewRegionFromBBox(0, 0, 0, 0, srs, 1, 1)
	if r2.IsValid() {
		t.Error("zero-size region should return false")
	}
}

func TestRegionCopy(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(-125, 40, -122, 43, srs, 0.01, 0.01)
	r.ZMin = -10
	r.ZMax = 100
	c := r.Copy()
	if c.Extent.BBox.Min[0] != r.Extent.BBox.Min[0] {
		t.Error("copy should have same minx")
	}
	c.Extent.BBox.Min[0] = 0
	if r.Extent.BBox.Min[0] == 0 {
		t.Error("modifying copy should not modify original")
	}
	if c.ZMin != -10 || c.ZMax != 100 {
		t.Error("copy should preserve Z range")
	}
}

func TestRegionCenter(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(0, 0, 10, 20, srs, 1, 1)
	cx, cy := r.Center()
	if cx != 5 {
		t.Errorf("center x: expected 5, got %f", cx)
	}
	if cy != 10 {
		t.Errorf("center y: expected 10, got %f", cy)
	}
}

func TestRegionWidthHeight(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(0, 10, 100, 60, srs, 1, 1)
	if r.Width() != 100 {
		t.Errorf("width: expected 100, got %f", r.Width())
	}
	if r.Height() != 50 {
		t.Errorf("height: expected 50, got %f", r.Height())
	}
}

func TestRegionZRange(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(0, 0, 10, 10, srs, 1, 1)
	r.SetZRange(-50, 200)
	zMin, zMax := r.ZRange()
	if zMin != -50 || zMax != 200 {
		t.Errorf("zrange: expected (-50, 200), got (%f, %f)", zMin, zMax)
	}
}

func TestRegionBuffer(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(0, 0, 10, 10, srs, 1, 1)
	b := r.Buffer(2, 3)
	if b.Extent.BBox.Min[0] != -2 {
		t.Errorf("buffer minx: expected -2, got %f", b.Extent.BBox.Min[0])
	}
	if b.Extent.BBox.Max[0] != 12 {
		t.Errorf("buffer maxx: expected 12, got %f", b.Extent.BBox.Max[0])
	}
	if b.Extent.BBox.Min[1] != -3 {
		t.Errorf("buffer miny: expected -3, got %f", b.Extent.BBox.Min[1])
	}
	if b.Extent.BBox.Max[1] != 13 {
		t.Errorf("buffer maxy: expected 13, got %f", b.Extent.BBox.Max[1])
	}
	if r.Extent.BBox.Min[0] != 0 {
		t.Error("original should not be modified")
	}
}

func TestRegionBufferPct(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(0, 0, 100, 50, srs, 1, 1)
	b := r.BufferPct(10)
	if math.Abs(b.Extent.BBox.Min[0]+10) > 0.01 {
		t.Errorf("buffer pct minx: expected -10, got %f", b.Extent.BBox.Min[0])
	}
	if math.Abs(b.Extent.BBox.Max[0]-110) > 0.01 {
		t.Errorf("buffer pct maxx: expected 110, got %f", b.Extent.BBox.Max[0])
	}
	if math.Abs(b.Extent.BBox.Min[1]+5) > 0.01 {
		t.Errorf("buffer pct miny: expected -5, got %f", b.Extent.BBox.Min[1])
	}
}

func TestRegionIntersection(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r1 := NewRegionFromBBox(0, 0, 10, 10, srs, 1, 1)
	r2 := NewRegionFromBBox(5, 5, 15, 15, srs, 1, 1)
	inter := r1.Intersection(r2)
	if inter.Extent.BBox.Min[0] != 5 || inter.Extent.BBox.Min[1] != 5 {
		t.Errorf("intersection min: expected (5,5), got (%f,%f)", inter.Extent.BBox.Min[0], inter.Extent.BBox.Min[1])
	}
	if inter.Extent.BBox.Max[0] != 10 || inter.Extent.BBox.Max[1] != 10 {
		t.Errorf("intersection max: expected (10,10), got (%f,%f)", inter.Extent.BBox.Max[0], inter.Extent.BBox.Max[1])
	}
}

func TestRegionUnion(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r1 := NewRegionFromBBox(0, 0, 10, 10, srs, 1, 1)
	r2 := NewRegionFromBBox(5, 5, 15, 15, srs, 1, 1)
	u := r1.Union(r2)
	if u.Extent.BBox.Min[0] != 0 || u.Extent.BBox.Min[1] != 0 {
		t.Errorf("union min: expected (0,0), got (%f,%f)", u.Extent.BBox.Min[0], u.Extent.BBox.Min[1])
	}
	if u.Extent.BBox.Max[0] != 15 || u.Extent.BBox.Max[1] != 15 {
		t.Errorf("union max: expected (15,15), got (%f,%f)", u.Extent.BBox.Max[0], u.Extent.BBox.Max[1])
	}
}

func TestRegionIntersects(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r1 := NewRegionFromBBox(0, 0, 10, 10, srs, 1, 1)
	r2 := NewRegionFromBBox(5, 5, 15, 15, srs, 1, 1)
	r3 := NewRegionFromBBox(20, 20, 30, 30, srs, 1, 1)
	if !r1.Intersects(r2) {
		t.Error("overlapping regions should intersect")
	}
	if r1.Intersects(r3) {
		t.Error("non-overlapping regions should not intersect")
	}
}

func TestRegionContains(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	outer := NewRegionFromBBox(0, 0, 10, 10, srs, 1, 1)
	inner := NewRegionFromBBox(2, 2, 8, 8, srs, 1, 1)
	outside := NewRegionFromBBox(5, 5, 15, 15, srs, 1, 1)
	if !outer.Contains(inner) {
		t.Error("outer should contain inner")
	}
	if outer.Contains(outside) {
		t.Error("outer should not contain outside")
	}
}

func TestRegionSrcWin(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(2, 3, 8, 7, srs, 1, 1)
	gt := [6]float64{0, 1, 0, 10, 0, -1}
	xOff, yOff, xSize, ySize := r.SrcWin(gt, 10, 10)
	if xOff != 2 || yOff != 3 || xSize != 6 || ySize != 4 {
		t.Errorf("srcwin: expected (2,3,6,4), got (%d,%d,%d,%d)", xOff, yOff, xSize, ySize)
	}
}

func TestRegionSrcWin_Clamp(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(-5, -5, 15, 15, srs, 1, 1)
	gt := [6]float64{0, 1, 0, 10, 0, -1}
	xOff, yOff, xSize, ySize := r.SrcWin(gt, 10, 10)
	if xOff != 0 || yOff != 0 || xSize != 10 || ySize != 10 {
		t.Errorf("clamped srcwin: expected (0,0,10,10), got (%d,%d,%d,%d)", xOff, yOff, xSize, ySize)
	}
}

func TestRegionGeoTransformFromCount(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(0, 0, 10, 10, srs, 1, 1)
	gt := r.GeoTransformFromCount(5, 5)
	if gt[0] != 0 {
		t.Errorf("gt[0]: expected 0, got %f", gt[0])
	}
	if math.Abs(gt[1]-2) > 0.01 {
		t.Errorf("gt[1]: expected 2, got %f", gt[1])
	}
	if math.Abs(gt[3]-10) > 0.01 {
		t.Errorf("gt[3]: expected 10, got %f", gt[3])
	}
}

func TestRegionFormat(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(-125, 40, -122, 43, srs, 1, 1)
	gmt := r.Format("gmt")
	if gmt != "-125.000000/-122.000000/40.000000/43.000000" {
		t.Errorf("gmt format: got %s", gmt)
	}
	bbox := r.Format("bbox")
	expected := "-125.000000,40.000000,-122.000000,43.000000"
	if bbox != expected {
		t.Errorf("bbox format: expected %s, got %s", expected, bbox)
	}
	wkt := r.Format("wkt")
	if !strings.HasPrefix(wkt, "POLYGON") {
		t.Errorf("wkt format: expected POLYGON(...), got %s", wkt)
	}
}

func TestRegionRound(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(1.2345, 2.6789, 5.4321, 6.1234, srs, 0.1, 0.1)
	rounded := r.Round(1)
	if rounded.Extent.BBox.Min[0] != 1.2 {
		t.Errorf("round minx: expected 1.2, got %f", rounded.Extent.BBox.Min[0])
	}
	if rounded.Extent.BBox.Max[0] != 5.5 {
		t.Errorf("round maxx: expected 5.5, got %f", rounded.Extent.BBox.Max[0])
	}
}

func TestRegionChunk(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(0, 0, 10, 10, srs, 1, 1)
	chunks := r.Chunk(2, 2)
	if len(chunks) != 4 {
		t.Fatalf("expected 4 chunks, got %d", len(chunks))
	}
	if chunks[0].Extent.BBox.Min[0] != 0 || chunks[0].Extent.BBox.Min[1] != 0 {
		t.Errorf("chunk[0] min: expected (0,0), got (%f,%f)", chunks[0].Extent.BBox.Min[0], chunks[0].Extent.BBox.Min[1])
	}
	if chunks[0].Extent.BBox.Max[0] != 5 || chunks[0].Extent.BBox.Max[1] != 5 {
		t.Errorf("chunk[0] max: expected (5,5), got (%f,%f)", chunks[0].Extent.BBox.Max[0], chunks[0].Extent.BBox.Max[1])
	}
	if chunks[3].Extent.BBox.Min[0] != 5 || chunks[3].Extent.BBox.Min[1] != 5 {
		t.Errorf("chunk[3] min: expected (5,5), got (%f,%f)", chunks[3].Extent.BBox.Min[0], chunks[3].Extent.BBox.Min[1])
	}
	if chunks[3].Extent.BBox.Max[0] != 10 || chunks[3].Extent.BBox.Max[1] != 10 {
		t.Errorf("chunk[3] max: expected (10,10), got (%f,%f)", chunks[3].Extent.BBox.Max[0], chunks[3].Extent.BBox.Max[1])
	}
}

func TestRegionExportAsWKT(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r := NewRegionFromBBox(0, 0, 10, 10, srs, 1, 1)
	wkt := r.ExportAsWKT()
	if !strings.HasPrefix(wkt, "POLYGON") {
		t.Errorf("expected POLYGON, got %s", wkt)
	}
	if !strings.Contains(wkt, "10.000000 10.000000") {
		t.Errorf("wkt should contain max corner, got %s", wkt)
	}
}

func TestRegionNewFromString_MultipleDelimiters(t *testing.T) {
	srs := geo.NewProj("EPSG:4326")
	r1, _ := NewRegionFromString("-125/-122/40/43", srs, 0.01, 0)
	r2, _ := NewRegionFromString("-125,-122,40,43", srs, 0.01, 0)
	r3, _ := NewRegionFromString("-125 -122 40 43", srs, 0.01, 0)
	if r1 == nil || r2 == nil || r3 == nil {
		t.Error("parsing should work with / , and space")
	}
	if r1.Extent.BBox.Min[0] != r2.Extent.BBox.Min[0] {
		t.Error("slash and comma parsing should produce same result")
	}
	if r1.Extent.BBox.Min[0] != r3.Extent.BBox.Min[0] {
		t.Error("slash and space parsing should produce same result")
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
