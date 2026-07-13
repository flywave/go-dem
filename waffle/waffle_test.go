package waffle

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

func testRegion() *dem.Region {
	return dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
}

func TestWaffleFactory_Registry(t *testing.T) {
	methods := ListMethods()
	if len(methods) == 0 {
		t.Error("no methods registered")
	}
	for _, m := range methods {
		w, err := New(m)
		if err != nil {
			t.Errorf("factory: method %s: %v", m, err)
		}
		if w == nil {
			t.Errorf("factory: method %s returned nil", m)
		}
	}
}

func TestCollectPoints_Empty(t *testing.T) {
	pts, zs, err := collectPoints(nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(pts) != 0 || len(zs) != 0 {
		t.Errorf("expected empty, got %d pts %d zs", len(pts), len(zs))
	}
}

func TestBarycentricInterp_Inside(t *testing.T) {
	val, found := barycentricInterp(5, 5,
		vec2.T{0, 0}, vec2.T{10, 0}, vec2.T{0, 10},
		0, 10, 10)
	if !found {
		t.Error("point should be inside triangle")
	}
	if math.Abs(val-10) > 1e-10 {
		t.Errorf("expected 10, got %.4f", val)
	}
}

func TestBarycentricInterp_Outside(t *testing.T) {
	_, found := barycentricInterp(15, 15,
		vec2.T{0, 0}, vec2.T{10, 0}, vec2.T{0, 10},
		0, 10, 10)
	if found {
		t.Error("point should be outside triangle")
	}
}

func TestBarycentricInterp_Vertex(t *testing.T) {
	val, found := barycentricInterp(0, 0,
		vec2.T{0, 0}, vec2.T{10, 0}, vec2.T{0, 10},
		5, 10, 10)
	if !found {
		t.Error("vertex should be inside")
	}
	if math.Abs(val-5) > 1e-10 {
		t.Errorf("vertex: expected 5, got %.4f", val)
	}
}

func TestBarycentricInterp_Degenerate(t *testing.T) {
	_, found := barycentricInterp(0, 0,
		vec2.T{0, 0}, vec2.T{0, 0}, vec2.T{0, 0},
		0, 0, 0)
	if found {
		t.Error("degenerate triangle should not contain point")
	}
}

func TestDistSq(t *testing.T) {
	if math.Abs(distSq(0, 0, 3, 4)-25) > 1e-10 {
		t.Errorf("distSq(0,0,3,4) = %.4f, expected 25", distSq(0, 0, 3, 4))
	}
	if distSq(0, 0, 0, 0) != 0 {
		t.Errorf("distSq same point should be 0, got %.4f", distSq(0, 0, 0, 0))
	}
}

func TestNearestInterp_Basic(t *testing.T) {
	pts := []vec2.T{{0, 0}, {10, 0}, {0, 10}}
	zs := []float64{0, 20, 30}
	val := nearestInterp(1, 1, pts, zs)
	if math.Abs(val-0) > 1e-10 {
		t.Errorf("nearest to (0,0): expected 0, got %.2f", val)
	}
}

func TestNearestInterp_Duplicate(t *testing.T) {
	pts := []vec2.T{{5, 5}, {10, 10}}
	zs := []float64{100, 200}
	val := nearestInterp(5, 5, pts, zs)
	if math.Abs(val-100) > 1e-10 {
		t.Errorf("exact point: expected 100, got %.2f", val)
	}
}

func TestNearestInterp_Empty(t *testing.T) {
	val := nearestInterp(0, 0, nil, nil)
	if !math.IsNaN(val) {
		t.Errorf("empty: expected NaN, got %.2f", val)
	}
}

func TestAngleBetween(t *testing.T) {
	angle := angleBetween(vec2.T{0, 1}, 0, 0, vec2.T{1, 0})
	expected := math.Pi / 2
	if math.Abs(angle-expected) > 1e-10 {
		t.Errorf("expected %.4f, got %.4f", expected, angle)
	}
}

func TestAngleBetween_StraightLine(t *testing.T) {
	angle := angleBetween(vec2.T{0, 1}, 0, 0, vec2.T{0, 2})
	if math.Abs(angle) > 1e-10 {
		t.Errorf("collinear points: expected 0, got %.4f", angle)
	}
}

func TestTVU(t *testing.T) {
	v := tvu(100, 0.2, 0.01)
	expected := math.Sqrt(0.04 + 1.0)
	if math.Abs(v-expected) > 1e-10 {
		t.Errorf("tvu(100, 0.2, 0.01) = %.4f, expected %.4f", v, expected)
	}
}

func TestTVU_Shallow(t *testing.T) {
	v := tvu(0, 0.2, 0.01)
	if math.Abs(v-0.2) > 1e-10 {
		t.Errorf("tvu(0, 0.2, 0.01) = %.4f, expected 0.2", v)
	}
}

func TestDensityClusterHypotheses(t *testing.T) {
	sp := soundingParams{TVUa: 0.2, TVUb: 0.01, THU: 2.0}
	depths := []float64{10, 10.1, 10.2, 10.3, 20, 20.1, 20.2}
	weights := []float64{1, 1, 1, 1, 1, 1, 1}
	h := densityClusterHypotheses(depths, weights, sp)
	if len(h) == 0 {
		t.Error("should find at least one hypothesis")
	}
	if len(h) >= 1 {
		t.Logf("found %d clusters", len(h))
	}
}

func TestSelectHypothesisByIC(t *testing.T) {
	h := []cubeHypothesis{
		{mean: 10, stdDev: 0.1, count: 10},
		{mean: 20, stdDev: 0.5, count: 3},
	}
	depths := []float64{9.9, 10.1, 10.0, 10.2}
	best := selectHypothesisByIC(h, depths, nil)
	if best == nil {
		t.Fatal("no hypothesis selected")
	}
	if math.Abs(best.mean-10) > 1.0 {
		t.Errorf("expected hypothesis 0 (mean=10), got mean=%.4f", best.mean)
	}
}

func TestBoundingBox(t *testing.T) {
	tri := [3]int{0, 1, 2}
	pts := []vec2.T{{-5, -5}, {10, 3}, {2, 15}}
	bbox := boundingBox(tri, pts)
	if math.Abs(bbox[0]+5) > 1e-10 || math.Abs(bbox[2]+5) > 1e-10 {
		t.Errorf("bbox y-min: expected -5, got %.2f", bbox[2])
	}
}

func TestBuildTriangleGridIndex(t *testing.T) {
	pts := make([]vec2.T, 10)
	for i := range pts {
		pts[i] = vec2.T{float64(i % 5), float64(i / 5)}
	}
	triList := [][3]int{
		{0, 1, 5},
		{1, 2, 6},
		{2, 3, 7},
	}
	idx := buildTriangleGridIndex(triList, pts, 10)
	if idx.gridW != 10 || idx.gridH != 10 {
		t.Errorf("grid size: expected 10x10, got %dx%d", idx.gridW, idx.gridH)
	}
}

func TestGridIndex_FindTriangles(t *testing.T) {
	pts := []vec2.T{{0, 0}, {10, 0}, {0, 10}}
	triList := [][3]int{{0, 1, 2}}
	idx := buildTriangleGridIndex(triList, pts, 10)

	tris := idx.findTriangles(2, 2)
	_ = tris
}

func TestIsLASFile(t *testing.T) {
	if !isLASFile("data.las") {
		t.Error(".las should be detected")
	}
	if !isLASFile("data.laz") {
		t.Error(".laz should be detected")
	}
	if isLASFile("data.tif") {
		t.Error(".tif should not be las")
	}
	if isLASFile("short") {
		t.Error("short string should not be las")
	}
}

func TestCoalesceNoData(t *testing.T) {
	v := 0.0
	if dem.CoalesceNoData(&v) != 0 {
		t.Error("should return user value")
	}
	if dem.CoalesceNoData(nil) != dem.DefaultNoData {
		t.Error("should return default")
	}
}

func TestIsNoDataValue(t *testing.T) {
	if !dem.IsNoDataValue(dem.DefaultNoData) {
		t.Error("default noData should be detected")
	}
	if dem.IsNoDataValue(0) {
		t.Error("0 should not be noData")
	}
}

func TestLaplaceWeightedInterp_ExactVertex(t *testing.T) {
	val := laplaceWeightedInterp(0, 0,
		vec2.T{0, 0}, vec2.T{10, 0}, vec2.T{0, 10},
		100, 200, 300)
	if math.Abs(val-100) > 1e-6 {
		t.Errorf("vertex: expected 100, got %.2f", val)
	}
}

func TestMax(t *testing.T) {
	if max(5, 3) != 5 {
		t.Errorf("max(5,3) = %d", max(5, 3))
	}
	if max(-1, 0) != 0 {
		t.Errorf("max(-1,0) = %d", max(-1, 0))
	}
}
