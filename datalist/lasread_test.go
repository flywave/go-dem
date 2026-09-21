package datalist

import (
	"context"
	"errors"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flywave/flywave-pointcloud/pdal"
	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/waffle"
)

var testLASPoints = [][3]float64{
	{0.5, 0.25, 10.5},
	{1.5, 1.25, 20.5},
	{2.5, 2.75, -15.25},
	{1000, 2000, 0.5},
	{3.5, 4.25, 30.5},
}

func writeTestLASFile(path string, pts [][3]float64, srs string) error {
	pdal.PdalInitStage()

	stg := pdal.NewStage("writers.las")
	if stg == nil {
		return errors.New("failed to create writer stage")
	}
	defer stg.Free()
	opt := pdal.NewOptions()
	defer opt.Free()
	opt.Add("filename", path)
	if srs != "" {
		opt.Add("a_srs", srs)
	}
	stg.AddOptions(opt)

	tab := pdal.NewPointTable()
	defer tab.Free()

	view := pdal.NewPointViewTable(tab)
	defer view.Free()
	l := view.Layout()
	defer l.Free()
	l.RegisterDimAndType(pdal.X, pdal.Double)
	l.RegisterDimAndType(pdal.Y, pdal.Double)
	l.RegisterDimAndType(pdal.Z, pdal.Double)

	if err := stg.Prepare(tab); err != nil {
		return err
	}
	if v := stg.Execute(tab); v == nil {
		return errors.New("first writer execute failed")
	}

	for i, p := range pts {
		view.Set(uint64(i), pdal.X, &p[0])
		view.Set(uint64(i), pdal.Y, &p[1])
		view.Set(uint64(i), pdal.Z, &p[2])
	}

	buff := pdal.NewBufferReader(view)
	defer buff.Free()
	stg.SetBufferReader(buff)
	if err := stg.Prepare(tab); err != nil {
		return err
	}
	if v := stg.Execute(tab); v == nil {
		return errors.New("writer execute failed: " + pdal.LastError())
	}

	_, err := os.Stat(path)
	return err
}

func samePoint(p waffle.Point, x, y, z float64) bool {
	return math.Abs(p.Position[0]-x) < 1e-6 &&
		math.Abs(p.Position[1]-y) < 1e-6 &&
		math.Abs(p.Z-z) < 1e-6
}

func containsPoint(pts []waffle.Point, x, y, z float64) bool {
	for _, p := range pts {
		if samePoint(p, x, y, z) {
			return true
		}
	}
	return false
}

func TestReadLASPoints_LASRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rt.las")
	if err := writeTestLASFile(path, testLASPoints, ""); err != nil {
		t.Fatalf("write test las: %v", err)
	}

	pts, err := ReadLASPoints(path)
	if err != nil {
		t.Fatalf("read las: %v", err)
	}
	if len(pts) != len(testLASPoints) {
		t.Fatalf("expected %d points, got %d", len(testLASPoints), len(pts))
	}
	for _, want := range testLASPoints {
		if !containsPoint(pts, want[0], want[1], want[2]) {
			t.Errorf("point (%v,%v,%v) not found in read-back result", want[0], want[1], want[2])
		}
	}
}

func TestReadLASPointsWithSRS(t *testing.T) {
	noSRS := filepath.Join(t.TempDir(), "nosrs.las")
	if err := writeTestLASFile(noSRS, testLASPoints, ""); err != nil {
		t.Fatalf("write test las: %v", err)
	}
	pts, wkt, err := ReadLASPointsWithSRS(noSRS)
	if err != nil {
		t.Fatalf("read las without srs: %v", err)
	}
	if len(pts) != len(testLASPoints) {
		t.Fatalf("expected %d points, got %d", len(testLASPoints), len(pts))
	}
	if wkt != "" {
		t.Errorf("expected empty srs wkt, got: %.120s", wkt)
	}

	withSRS := filepath.Join("..", "..", "flywave-pointcloud", "tests", "las", "test_epsg_4326.las")
	if _, err := os.Stat(withSRS); err != nil {
		t.Skip("sibling test data with embedded SRS not present")
	}
	srsPts, srsWKT, err := ReadLASPointsWithSRS(withSRS)
	if err != nil {
		t.Fatalf("read las with srs: %v", err)
	}
	if len(srsPts) == 0 {
		t.Fatal("srs las file yielded zero points")
	}
	if srsWKT == "" {
		t.Fatal("expected non-empty srs wkt")
	}
	if !strings.Contains(srsWKT, "WGS 84") {
		t.Errorf("wkt should reference WGS 84, got: %.120s", srsWKT)
	}
}

func TestReadLASPoints_LAZRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rt.laz")
	if err := writeTestLASFile(path, testLASPoints, ""); err != nil {
		t.Skipf("laz compression unavailable in this environment: %v", err)
	}

	pts, err := ReadLASPoints(path)
	if err != nil {
		t.Fatalf("read laz: %v", err)
	}
	if len(pts) != len(testLASPoints) {
		t.Fatalf("expected %d points, got %d", len(testLASPoints), len(pts))
	}
	for _, want := range testLASPoints {
		if !containsPoint(pts, want[0], want[1], want[2]) {
			t.Errorf("point (%v,%v,%v) not found in laz read-back result", want[0], want[1], want[2])
		}
	}
}

func TestReadLASPoints_SiblingTestData(t *testing.T) {
	lasPath := filepath.Join("..", "..", "flywave-pointcloud", "tests", "las", "simple.las")
	if _, err := os.Stat(lasPath); err != nil {
		t.Skip("sibling flywave-pointcloud test data not present")
	}
	pts, err := ReadLASPoints(lasPath)
	if err != nil {
		t.Fatalf("read real las: %v", err)
	}
	if len(pts) == 0 {
		t.Fatal("real las file yielded zero points")
	}

	lazCandidates := []string{
		filepath.Join("..", "..", "flywave-pointcloud", "tests", "laz", "autzen_trim.laz"),
		filepath.Join("..", "..", "flywave-pointcloud", "tests", "mytest.laz"),
		filepath.Join("..", "..", "flywave-pointcloud", "tests", "laz", "simple.laz"),
	}
	var lazPath string
	for _, c := range lazCandidates {
		if _, err := os.Stat(c); err == nil {
			lazPath = c
			break
		}
	}
	if lazPath == "" {
		t.Skip("no real laz sample present")
	}
	lazPts, err := ReadLASPoints(lazPath)
	if err != nil {
		t.Skipf("real laz sample %s not readable by bundled lazperf: %v", lazPath, err)
	}
	if len(lazPts) == 0 {
		t.Fatalf("real laz file %s yielded zero points", lazPath)
	}
}

func TestReadLASPoints_MissingFile(t *testing.T) {
	if _, err := ReadLASPoints(filepath.Join(t.TempDir(), "nope.las")); err == nil {
		t.Error("expected error for missing las file")
	}
}

func TestReadLASPoints_EmptyFile(t *testing.T) {
	dir := t.TempDir()

	empty := filepath.Join(dir, "empty.las")
	if err := os.WriteFile(empty, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadLASPoints(empty); err == nil {
		t.Error("expected error for empty las file")
	}

	garbage := filepath.Join(dir, "garbage.las")
	if err := os.WriteFile(garbage, []byte(strings.Repeat("not a las file", 20)), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadLASPoints(garbage); err == nil {
		t.Error("expected error for non-las content")
	}
}

func TestReadLASPointsWithProgress(t *testing.T) {
	path := filepath.Join(t.TempDir(), "prog.las")
	if err := writeTestLASFile(path, testLASPoints, ""); err != nil {
		t.Fatalf("write test las: %v", err)
	}

	var calls int
	var lastStage string
	var lastDone, lastTotal int
	pts, err := ReadLASPointsWithProgress(path, func(stage string, done, total int) {
		calls++
		lastStage = stage
		lastDone, lastTotal = done, total
	}, context.Background())
	if err != nil {
		t.Fatalf("read las with progress: %v", err)
	}
	if calls == 0 {
		t.Fatal("progress callback never invoked")
	}
	if lastStage != "las_read" {
		t.Errorf("stage name = %q, want las_read", lastStage)
	}
	if lastDone != lastTotal || lastTotal != len(testLASPoints) {
		t.Errorf("final progress = %d/%d, want %d/%d", lastDone, lastTotal, len(testLASPoints), len(testLASPoints))
	}
	if len(pts) != len(testLASPoints) {
		t.Fatalf("expected %d points, got %d", len(testLASPoints), len(pts))
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := ReadLASPointsWithProgress(path, nil, ctx); err == nil {
		t.Error("expected error with canceled context")
	} else if !strings.Contains(err.Error(), "cancel") {
		t.Errorf("error should mention cancellation, got: %v", err)
	}
}

func TestPointsFromDataList_Mixed(t *testing.T) {
	dir := t.TempDir()

	rasterPath := filepath.Join(dir, "grid.tif")
	reg := dem.NewRegionFromBBox(0, 0, 4, 4, nil, 1, 1)
	data := make([]float64, reg.XSize*reg.YSize)
	for i := range data {
		data[i] = float64(i)
	}
	data[0] = dem.DefaultNoData
	data[15] = dem.DefaultNoData
	if err := dem.CreateDEM(data, reg, rasterPath, dem.DefaultNoData); err != nil {
		t.Fatalf("create raster fixture: %v", err)
	}

	xyzPath := filepath.Join(dir, "pts.xyz")
	if err := os.WriteFile(xyzPath, []byte("1 1 11\n2 2 22\n3 3 33\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	lasPath := filepath.Join(dir, "pts.las")
	if err := writeTestLASFile(lasPath, testLASPoints, ""); err != nil {
		t.Fatalf("write test las: %v", err)
	}

	dl := &DataList{Entries: []DataEntry{
		{Path: rasterPath, Type: SourceRaster},
		{Path: xyzPath, Type: SourcePoint},
		{Path: lasPath, Type: SourcePoint},
	}}
	pts, err := PointsFromDataList(dl)
	if err != nil {
		t.Fatalf("points from datalist: %v", err)
	}
	wantRaster := int(reg.XSize*reg.YSize) - 2
	if len(pts) != wantRaster+3+len(testLASPoints) {
		t.Fatalf("expected %d points, got %d", wantRaster+3+len(testLASPoints), len(pts))
	}
	if !containsPoint(pts, 1, 1, 11) {
		t.Error("xyz source point (1,1,11) missing")
	}
	if !containsPoint(pts, 1000, 2000, 0.5) {
		t.Error("las source point (1000,2000,0.5) missing")
	}
	rasterFound := false
	for _, p := range pts {
		if p.Z == 1 {
			rasterFound = true
			break
		}
	}
	if !rasterFound {
		t.Error("no raster source points found")
	}

	vectorList := &DataList{Entries: []DataEntry{{Path: filepath.Join(dir, "b.shp"), Type: SourceVector}}}
	if _, err := PointsFromDataList(vectorList); err == nil {
		t.Error("expected error for vector entry")
	}
	if _, err := PointsFromDataList(nil); err == nil {
		t.Error("expected error for nil datalist")
	}
	if _, err := PointsFromDataList(&DataList{}); err == nil {
		t.Error("expected error for empty datalist")
	}
}
