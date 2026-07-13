package waffle

import (
	"math"
	"testing"

	"github.com/flywave/go3d/float64/vec2"
)

func TestKDTree_Empty(t *testing.T) {
	tree := NewKDTree(nil)
	if tree.Root != nil {
		t.Error("empty tree should have nil root")
	}
	pts := tree.Points()
	if len(pts) != 0 {
		t.Error("empty tree should return no points")
	}
}

func TestKDTree_SinglePoint(t *testing.T) {
	pts := []vec2.T{{5, 5}}
	tree := NewKDTree(pts)
	if tree.Root == nil {
		t.Fatal("tree root is nil")
	}
	if tree.Root.Point[0] != 5 || tree.Root.Point[1] != 5 {
		t.Errorf("root: expected (5,5), got (%v)", tree.Root.Point)
	}
}

func TestKDTree_KNNAll(t *testing.T) {
	pts := []vec2.T{{0, 0}, {1, 1}, {2, 2}, {10, 10}}
	tree := NewKDTree(pts)
	idxs, dists := tree.KNN(vec2.T{0, 0}, 4)
	if len(idxs) != 4 {
		t.Errorf("expected 4 neighbors, got %d", len(idxs))
	}
	if idxs[0] != 0 {
		t.Errorf("nearest should be index 0, got %d", idxs[0])
	}
	if math.Abs(dists[0]) > 1e-10 {
		t.Errorf("distance to self should be 0, got %.4f", dists[0])
	}
}

func TestKDTree_KNNOrdered(t *testing.T) {
	pts := []vec2.T{{0, 0}, {10, 0}, {5, 0}}
	tree := NewKDTree(pts)
	idxs, dists := tree.KNN(vec2.T{3, 0}, 3)
	if len(idxs) != 3 {
		t.Fatalf("expected 3, got %d", len(idxs))
	}
	if math.Abs(dists[0]-2) > 1e-10 {
		t.Errorf("nearest distance should be 2 (to (5,0)), got %.4f", dists[0])
	}
	if dists[0] > dists[1] || dists[1] > dists[2] {
		t.Error("distances not sorted ascending")
	}
}

func TestKDTree_RadiusSearch(t *testing.T) {
	pts := []vec2.T{{0, 0}, {1, 0}, {2, 0}, {10, 10}}
	tree := NewKDTree(pts)
	idxs, dists := tree.RadiusSearch(vec2.T{0, 0}, 1.5)
	if len(idxs) != 2 {
		t.Errorf("expected 2 within radius 1.5, got %d", len(idxs))
	}

	foundSelf := false
	for i, idx := range idxs {
		if idx == 0 {
			foundSelf = true
			if math.Abs(dists[i]) > 1e-10 {
				t.Errorf("distance to (0,0) should be 0, got %.4f", dists[i])
			}
		}
	}
	if !foundSelf {
		t.Error("radius search did not find point (0,0)")
	}
}

func TestKDTree_RadiusSearchNone(t *testing.T) {
	pts := []vec2.T{{0, 0}, {1, 1}}
	tree := NewKDTree(pts)
	idxs, _ := tree.RadiusSearch(vec2.T{10, 10}, 1)
	if len(idxs) != 0 {
		t.Errorf("expected 0, got %d", len(idxs))
	}
}

func TestKDTree_Points(t *testing.T) {
	pts := []vec2.T{{3, 1}, {1, 4}, {2, 2}, {5, 3}}
	tree := NewKDTree(pts)
	got := tree.Points()
	if len(got) != len(pts) {
		t.Errorf("expected %d points, got %d", len(pts), len(got))
	}
}

func TestKDTree_KNNWithDupes(t *testing.T) {
	pts := []vec2.T{{0, 0}, {0, 0}, {0, 0}}
	tree := NewKDTree(pts)
	idxs, dists := tree.KNN(vec2.T{0, 0}, 3)
	if len(idxs) != 3 {
		t.Errorf("expected 3, got %d", len(idxs))
	}
	for i, d := range dists {
		if d != 0 {
			t.Errorf("dist[%d] should be 0, got %.4f", i, d)
		}
	}
}

func BenchmarkKDTree_KNN(b *testing.B) {
	n := 10000
	pts := make([]vec2.T, n)
	for i := range pts {
		pts[i] = vec2.T{float64(i % 100), float64(i / 100)}
	}
	tree := NewKDTree(pts)
	q := vec2.T{50, 50}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.KNN(q, 10)
	}
}
