package delaunay

import (
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
}

func TestTriangulate_Triangle(t *testing.T) {
	x := []float64{0, 1, 0.5}
	y := []float64{0, 0, 1}
	tris, _, err := Triangulate(x, y)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tris) != 1 {
		t.Errorf("expected 1 triangle, got %d", len(tris))
	}
}

func TestTriangulate_Collinear(t *testing.T) {
	x := []float64{0, 1, 2}
	y := []float64{0, 0, 0}
	_, _, err := Triangulate(x, y)
	if err == nil {
		t.Error("expected error for collinear input")
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
	if len(tris) != 2*n-2-5-4+1 {
		t.Logf("grid 5x4: %d triangles", len(tris))
	}
	if len(nbrs) != len(tris) {
		t.Errorf("neighbor count %d != triangle count %d", len(nbrs), len(tris))
	}
	for i, tr := range tris {
		o := Orient2D(x[tr[0]], y[tr[0]], x[tr[1]], y[tr[1]], x[tr[2]], y[tr[2]])
		if o <= 0 {
			t.Errorf("triangle %d (%v) is not anticlockwise (orient=%d)", i, tr, o)
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
