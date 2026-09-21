package delaunay

import (
	"math"
	"math/rand"
	"strings"
	"testing"
)

func TestTriangulate_UnitSquare(t *testing.T) {
	x := []float64{0, 1, 0, 1}
	y := []float64{0, 0, 1, 1}
	tris, nbrs, err := Triangulate(x, y)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tris) != 2 {
		t.Errorf("expected 2 triangles, got %d", len(tris))
	} else {
		t.Logf("triangles: %v", tris)
	}
	if len(nbrs) != 2 {
		t.Errorf("expected 2 neighbor entries, got %d", len(nbrs))
	}
	checkTriangulation(t, x, y, tris, nbrs)
}

func TestTriangulate_Triangle(t *testing.T) {
	x := []float64{0, 1, 0.5}
	y := []float64{0, 0, 1}
	tris, nbrs, err := Triangulate(x, y)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tris) != 1 {
		t.Errorf("expected 1 triangle, got %d", len(tris))
	}
	checkTriangulation(t, x, y, tris, nbrs)
}

func TestTriangulate_Collinear(t *testing.T) {
	x := []float64{0, 1, 2}
	y := []float64{0, 0, 0}
	_, _, err := Triangulate(x, y)
	if err == nil {
		t.Error("expected error for collinear input")
	} else if !strings.Contains(err.Error(), "collinear") {
		t.Errorf("expected collinear message, got %q", err.Error())
	}
}

func TestTriangulate_Coincident(t *testing.T) {
	_, _, err := Triangulate([]float64{1, 1, 1, 1}, []float64{2, 2, 2, 2})
	if err == nil {
		t.Fatal("expected error for coincident input")
	}
	if !strings.Contains(err.Error(), "coincide") {
		t.Errorf("expected coincident message, got %q", err.Error())
	}
}

func TestTriangulate_NonFinite(t *testing.T) {
	if _, _, err := Triangulate([]float64{0, 1, math.NaN()}, []float64{0, 0, 1}); err == nil {
		t.Error("expected error for NaN x coordinate")
	}
	if _, _, err := Triangulate([]float64{0, 1, 2}, []float64{0, math.Inf(-1), 1}); err == nil {
		t.Error("expected error for -Inf y coordinate")
	}
	if _, _, err := Triangulate([]float64{0, math.Inf(1), 2}, []float64{0, 0, 1}); err == nil {
		t.Error("expected error for +Inf x coordinate")
	}
}

func TestTriangulate_TooFew(t *testing.T) {
	_, _, err := Triangulate([]float64{0, 1}, []float64{0, 0})
	if err == nil {
		t.Error("expected error for <3 points")
	}
}

func TestTriangulate_MismatchedLength(t *testing.T) {
	_, _, err := Triangulate([]float64{0, 1, 2}, []float64{0, 0})
	if err == nil {
		t.Error("expected error for mismatched lengths")
	}
}

func checkTriangulation(t *testing.T, x, y []float64, tris, nbrs [][3]int) {
	t.Helper()
	n := len(x)
	if len(nbrs) != len(tris) {
		t.Fatalf("neighbor count %d != triangle count %d", len(nbrs), len(tris))
	}
	used := make([]bool, n)
	boundary := 0
	for i, tr := range tris {
		if o := Orient2D(x[tr[0]], y[tr[0]], x[tr[1]], y[tr[1]], x[tr[2]], y[tr[2]]); o <= 0 {
			t.Errorf("triangle %d (%v) is not anticlockwise (orient=%d)", i, tr, o)
		}
		for j := 0; j < 3; j++ {
			if tr[j] < 0 || tr[j] >= n {
				t.Fatalf("triangle %d (%v) references invalid vertex", i, tr)
			}
			used[tr[j]] = true
			nb := nbrs[i][j]
			if nb < 0 {
				boundary++
				continue
			}
			if nb >= len(tris) {
				t.Fatalf("triangle %d has out-of-range neighbor %d", i, nb)
			}
			opp := -1
			for _, v := range tris[nb] {
				if v != tr[0] && v != tr[1] && v != tr[2] {
					opp = v
					break
				}
			}
			if opp < 0 {
				t.Fatalf("triangle %d and neighbor %d share no opposite vertex", i, nb)
			}
			if ic := InCircle(x[tr[0]], y[tr[0]], x[tr[1]], y[tr[1]], x[tr[2]], y[tr[2]], x[opp], y[opp]); ic > 0 {
				t.Errorf("triangle %d edge %d violates empty-circumcircle (opposite vertex %d, incircle=%d)", i, j, opp, ic)
			}
		}
	}
	for v, u := range used {
		if !u {
			t.Errorf("point %d is not a vertex of any triangle", v)
		}
	}
	if expect := 2*n - 2 - boundary; len(tris) != expect {
		t.Errorf("triangle count %d != 2*%d-2-%d = %d", len(tris), n, boundary, expect)
	}
}

func TestTriangulate_Grid5x4(t *testing.T) {
	n := 20
	x := make([]float64, n)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i % 5)
		y[i] = float64(i / 5)
	}
	tris, nbrs, err := Triangulate(x, y)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tris) != 24 {
		t.Errorf("grid 5x4: expected 24 triangles, got %d", len(tris))
	}
	checkTriangulation(t, x, y, tris, nbrs)
}

func TestTriangulate_Random5000(t *testing.T) {
	n := 5000
	x := make([]float64, n)
	y := make([]float64, n)
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < n; i++ {
		x[i] = rng.Float64()
		y[i] = rng.Float64()
	}
	tris, nbrs, err := Triangulate(x, y)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkTriangulation(t, x, y, tris, nbrs)
}

func BenchmarkTriangulate_5000(b *testing.B) {
	n := 5000
	x := make([]float64, n)
	y := make([]float64, n)
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < n; i++ {
		x[i] = rng.Float64()
		y[i] = rng.Float64()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := Triangulate(x, y); err != nil {
			b.Fatal(err)
		}
	}
}

func TestOrient2D(t *testing.T) {
	if o := Orient2D(0, 0, 1, 0, 0, 1); o <= 0 {
		t.Errorf("anticlockwise triangle: expected >0, got %d", o)
	}
	if o := Orient2D(0, 0, 0, 1, 1, 0); o >= 0 {
		t.Errorf("clockwise triangle: expected <0, got %d", o)
	}
	if o := Orient2D(0, 0, 1, 0, 2, 0); o != 0 {
		t.Errorf("collinear: expected 0, got %d", o)
	}
}

func TestInCircle(t *testing.T) {
	// Unit square: (0,0), (1,0), (1,1), (0,1) — all cocircular
	ic := InCircle(0, 0, 1, 0, 1, 1, 0, 1)
	if ic != 0 {
		t.Errorf("cocircular: expected 0, got %d", ic)
	}
	// (0.5, 0.5) is inside the circumcircle of (0,0), (1,0), (0,1)
	ic2 := InCircle(0, 0, 1, 0, 0, 1, 0.5, 0.5)
	if ic2 <= 0 {
		t.Errorf("inside circle: expected >0, got %d", ic2)
	}
}
