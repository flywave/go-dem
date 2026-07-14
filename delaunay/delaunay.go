package delaunay

import (
	"fmt"
	"sort"
)

func Triangulate(x, y []float64) (triangles, neighbors [][3]int, err error) {
	return DelaunayFast(x, y)
}

func DelaunayFast(x, y []float64) (triangles, neighbors [][3]int, err error) {
	if len(x) != len(y) {
		return nil, nil, fmt.Errorf("delaunay: x and y length mismatch (%d vs %d)", len(x), len(y))
	}
	if len(x) < 3 {
		return nil, nil, fmt.Errorf("delaunay: need at least 3 points, got %d", len(x))
	}

	tris, ok := bowyerWatson(x, y)
	if !ok {
		return nil, nil, fmt.Errorf("delaunay: degenerate input (all points collinear?)")
	}
	neighbors = computeNeighbors(tris)
	return tris, neighbors, nil
}

type triEdge struct{ u, v int }

func computeNeighbors(tris [][3]int) [][3]int {
	loc := make(map[triEdge][2]int, len(tris)*3)
	for i, tr := range tris {
		for j := 0; j < 3; j++ {
			loc[triEdge{tr[j], tr[(j+1)%3]}] = [2]int{i, j}
		}
	}
	nbrs := make([][3]int, len(tris))
	for i, tr := range tris {
		nbrs[i] = [3]int{-1, -1, -1}
		for j := 0; j < 3; j++ {
			if l, ok := loc[triEdge{tr[(j+1)%3], tr[j]}]; ok {
				nbrs[i][j] = l[0]
			}
		}
	}
	return nbrs
}

type bwTri struct{ a, b, c int }

func bowyerWatson(x, y []float64) ([][3]int, bool) {
	n := len(x)
	minX, maxX, minY, maxY := x[0], x[0], y[0], y[0]
	for i := 1; i < n; i++ {
		if x[i] < minX {
			minX = x[i]
		}
		if x[i] > maxX {
			maxX = x[i]
		}
		if y[i] < minY {
			minY = y[i]
		}
		if y[i] > maxY {
			maxY = y[i]
		}
	}
	dx, dy := maxX-minX, maxY-minY
	delta := dx
	if dy > delta {
		delta = dy
	}
	if delta == 0 {
		return nil, false
	}
	midX, midY := (minX+maxX)/2, (minY+maxY)/2

	const k = 100000.0
	px := append(append([]float64(nil), x...), midX-k*delta, midX, midX+k*delta)
	py := append(append([]float64(nil), y...), midY-delta, midY+k*delta, midY-delta)

	super, ok := orientTri(n, n+1, n+2, px, py)
	if !ok {
		return nil, false
	}
	tris := []bwTri{super}

	for p := 0; p < n; p++ {
		bad := make([]bool, len(tris))
		boundary := make(map[[2]int]int)
		for i, tr := range tris {
			if inCircle(px[tr.a], py[tr.a], px[tr.b], py[tr.b], px[tr.c], py[tr.c], px[p], py[p]) > 0 {
				bad[i] = true
				boundary[sortEdge(tr.a, tr.b)]++
				boundary[sortEdge(tr.b, tr.c)]++
				boundary[sortEdge(tr.c, tr.a)]++
			}
		}
		kept := tris[:0]
		for i, tr := range tris {
			if !bad[i] {
				kept = append(kept, tr)
			}
		}
		tris = kept
		for edge, count := range boundary {
			if count != 1 {
				continue
			}
			if tr, ok := orientTri(edge[0], edge[1], p, px, py); ok {
				tris = append(tris, tr)
			}
		}
	}

	out := make([][3]int, 0, len(tris))
	for _, tr := range tris {
		if tr.a >= n || tr.b >= n || tr.c >= n {
			continue
		}
		out = append(out, [3]int{tr.a, tr.b, tr.c})
	}
	sort.Slice(out, func(i, j int) bool {
		for k := 0; k < 3; k++ {
			if out[i][k] != out[j][k] {
				return out[i][k] < out[j][k]
			}
		}
		return false
	})
	return out, len(out) > 0
}

func orientTri(a, b, c int, x, y []float64) (bwTri, bool) {
	s := Orient2D(x[a], y[a], x[b], y[b], x[c], y[c])
	if s == 0 {
		return bwTri{}, false
	}
	if s < 0 {
		a, b = b, a
	}
	return bwTri{a, b, c}, true
}

func inCircle(ax, ay, bx, by, cx, cy, dx, dy float64) int {
	return InCircle(ax, ay, bx, by, cx, cy, dx, dy)
}

func sortEdge(a, b int) [2]int {
	if a < b {
		return [2]int{a, b}
	}
	return [2]int{b, a}
}
