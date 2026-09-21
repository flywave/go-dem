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

func TestPointsFromRaster_NoFile(t *testing.T) {
	_, err := PointsFromRaster("nonexistent.tif")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestPointsFromMultiple_Empty(t *testing.T) {
	_, err := PointsFromMultiple(nil)
	if err == nil {
		t.Error("expected error for nil input")
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
	if len(h) != 1 {
		t.Fatalf("expected single cluster within bandwidth, got %d", len(h))
	}
	if h[0].count != 7 {
		t.Errorf("expected cluster of 7, got %d", h[0].count)
	}
	expectedMean := (10 + 10.1 + 10.2 + 10.3 + 20 + 20.1 + 20.2) / 7
	if math.Abs(h[0].mean-expectedMean) > 1e-9 {
		t.Errorf("expected cluster mean %.4f, got %.4f", expectedMean, h[0].mean)
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
	if len(tris) == 0 {
		t.Fatal("expected candidate triangles at (2,2)")
	}
	found := false
	for _, ti := range tris {
		if ti == 0 {
			found = true
		}
	}
	if !found {
		t.Error("triangle 0 should be a candidate at (2,2)")
	}
	if idx.findTriangles(-5, -5) != nil {
		t.Error("expected nil candidates outside bbox (min side)")
	}
	if idx.findTriangles(11, 11) != nil {
		t.Error("expected nil candidates outside bbox (max side)")
	}
}

func TestLinearInterpGrid_OutsideBBox(t *testing.T) {
	pts := []vec2.T{{0, 0}, {10, 0}, {0, 10}}
	zs := []float64{1, 2, 3}
	triList := [][3]int{{0, 1, 2}}
	idx := buildTriangleGridIndex(triList, pts, 10)

	val := linearInterpGrid(50, 50, triList, pts, zs, &idx)
	if !math.IsNaN(val) {
		t.Errorf("expected NaN outside point bbox, got %.2f", val)
	}
	val = linearInterpGrid(2, 2, triList, pts, zs, &idx)
	if math.IsNaN(val) {
		t.Error("expected valid interpolation inside triangle")
	}
}

func TestWaffle_NilOptions(t *testing.T) {
	pts := []Point{
		{Position: vec2.T{1, 1}, Z: 1},
		{Position: vec2.T{2, 2}, Z: 2},
		{Position: vec2.T{3, 3}, Z: 3},
		{Position: vec2.T{4, 4}, Z: 4},
	}
	for _, m := range ListMethods() {
		w, err := New(m)
		if err != nil {
			t.Fatalf("new %s: %v", m, err)
		}
		res, err := w.Run(pts, nil)
		if err == nil {
			t.Errorf("%s: expected error for nil options", m)
		}
		if res != nil {
			t.Errorf("%s: expected nil result for nil options", m)
		}
		res, err = w.Run(pts, &Options{})
		if err == nil {
			t.Errorf("%s: expected error for nil region", m)
		}
		if res != nil {
			t.Errorf("%s: expected nil result for nil region", m)
		}
	}
}

func TestWaffle_EmptyPointsError(t *testing.T) {
	for _, m := range ListMethods() {
		w, err := New(m)
		if err != nil {
			t.Fatalf("new %s: %v", m, err)
		}
		res, err := w.Run(nil, &Options{Region: testRegion()})
		if err == nil {
			t.Errorf("%s: expected error for empty points", m)
		}
		if res != nil {
			t.Errorf("%s: expected nil result for empty points", m)
		}
	}
}

func TestPointStructure(t *testing.T) {
	p := Point{
		Position:    vec2.T{1, 2},
		Z:           100,
		Uncertainty: 0.5,
	}
	if p.Position[0] != 1 || p.Position[1] != 2 {
		t.Errorf("position: expected (1,2), got %v", p.Position)
	}
	if p.Z != 100 {
		t.Errorf("z: expected 100, got %.2f", p.Z)
	}
	if p.Uncertainty != 0.5 {
		t.Errorf("uncertainty: expected 0.5, got %.2f", p.Uncertainty)
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

func TestMaxInt(t *testing.T) {
	if maxInt(5, 3) != 5 {
		t.Errorf("maxInt(5,3) = %d", maxInt(5, 3))
	}
	if maxInt(-1, 0) != 0 {
		t.Errorf("maxInt(-1,0) = %d", maxInt(-1, 0))
	}
}
