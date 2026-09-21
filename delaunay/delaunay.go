package delaunay

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
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
	for i := range x {
		if math.IsNaN(x[i]) || math.IsInf(x[i], 0) {
			return nil, nil, fmt.Errorf("delaunay: non-finite x coordinate at index %d", i)
		}
		if math.IsNaN(y[i]) || math.IsInf(y[i], 0) {
			return nil, nil, fmt.Errorf("delaunay: non-finite y coordinate at index %d", i)
		}
	}

	tris, err := bowyerWatson(x, y)
	if err != nil {
		return nil, nil, err
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

type bwTri struct {
	v     [3]int
	nb    [3]int
	alive bool
}

type cavityEdge struct {
	a, b int
	ext  int
}

type pointPreds struct {
	px, py []float64
	rx, ry []*big.Rat
}

func (pp *pointPreds) orient(a, b, c int) int {
	if s, ok := orientFilter(pp.px[a], pp.py[a], pp.px[b], pp.py[b], pp.px[c], pp.py[c]); ok {
		return s
	}
	return orientRat(pp.rx[a], pp.ry[a], pp.rx[b], pp.ry[b], pp.rx[c], pp.ry[c])
}

func (pp *pointPreds) inCircle(a, b, c, d int) int {
	if s, ok := incircleFilter(pp.px[a], pp.py[a], pp.px[b], pp.py[b], pp.px[c], pp.py[c], pp.px[d], pp.py[d]); ok {
		return s
	}
	return incircleRat(pp.rx[a], pp.ry[a], pp.rx[b], pp.ry[b], pp.rx[c], pp.ry[c], pp.rx[d], pp.ry[d])
}

func bowyerWatson(x, y []float64) ([][3]int, error) {
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
		return nil, fmt.Errorf("delaunay: all points coincide")
	}
	midX := 0.5*minX + 0.5*maxX
	midY := 0.5*minY + 0.5*maxY

	const k = 100000.0
	px := append(append([]float64(nil), x...), midX-k*delta, midX, midX+k*delta)
	py := append(append([]float64(nil), y...), midY-delta, midY+k*delta, midY-delta)
	for i := n; i < n+3; i++ {
		if math.IsNaN(px[i]) || math.IsInf(px[i], 0) || math.IsNaN(py[i]) || math.IsInf(py[i], 0) {
			return nil, fmt.Errorf("delaunay: coordinate extent too large to enclose with a super-triangle")
		}
	}

	rx := make([]*big.Rat, n+3)
	ry := make([]*big.Rat, n+3)
	for i := range px {
		rx[i] = bigRat(px[i])
		ry[i] = bigRat(py[i])
	}
	pp := &pointPreds{px: px, py: py, rx: rx, ry: ry}

	super := bwTri{v: [3]int{n, n + 1, n + 2}, nb: [3]int{-1, -1, -1}, alive: true}
	if s := pp.orient(n, n+1, n+2); s < 0 {
		super.v[0], super.v[1] = super.v[1], super.v[0]
	}
	tris := []bwTri{super}

	order := rand.New(rand.NewSource(1)).Perm(n)

	mark := make([]int, 2*n+4)
	stamp := 0
	startIdx := make([]int, n+3)
	endIdx := make([]int, n+3)

	bad := make([]int, 0, 64)
	stack := make([]int, 0, 64)
	edges := make([]cavityEdge, 0, 64)

	last := 0
	for _, p := range order {
		last = locatePoint(p, last, tris, pp)
		if pp.inCircle(tris[last].v[0], tris[last].v[1], tris[last].v[2], p) <= 0 {
			continue
		}

		stamp++
		bad = bad[:0]
		stack = append(stack[:0], last)
		mark[last] = stamp
		for len(stack) > 0 {
			t := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			bad = append(bad, t)
			tr := &tris[t]
			for j := 0; j < 3; j++ {
				nb := tr.nb[j]
				if nb < 0 || mark[nb] == stamp {
					continue
				}
				if pp.inCircle(tris[nb].v[0], tris[nb].v[1], tris[nb].v[2], p) > 0 {
					mark[nb] = stamp
					stack = append(stack, nb)
				}
			}
		}

		edges = edges[:0]
		for _, t := range bad {
			tr := &tris[t]
			for j := 0; j < 3; j++ {
				nb := tr.nb[j]
				if nb < 0 || mark[nb] != stamp {
					edges = append(edges, cavityEdge{a: tr.v[j], b: tr.v[(j+1)%3], ext: nb})
				}
			}
		}

		if need := len(tris) + len(edges); need > len(mark) {
			mark = append(mark, make([]int, need+16-len(mark))...)
		}

		for _, t := range bad {
			tris[t].alive = false
		}

		for i, e := range edges {
			startIdx[e.a] = i
			endIdx[e.b] = i
		}

		first := len(tris)
		for _, e := range edges {
			tris = append(tris, bwTri{
				v:     [3]int{e.a, e.b, p},
				nb:    [3]int{e.ext, first + startIdx[e.b], first + endIdx[e.a]},
				alive: true,
			})
		}
		for i, e := range edges {
			if e.ext < 0 {
				continue
			}
			ext := &tris[e.ext]
			for j := 0; j < 3; j++ {
				if ext.v[j] == e.b && ext.v[(j+1)%3] == e.a {
					ext.nb[j] = first + i
					break
				}
			}
		}
		last = first
	}

	out := make([][3]int, 0, len(tris))
	for _, tr := range tris {
		if !tr.alive || tr.v[0] >= n || tr.v[1] >= n || tr.v[2] >= n {
			continue
		}
		out = append(out, tr.v)
	}
	sort.Slice(out, func(i, j int) bool {
		for k := 0; k < 3; k++ {
			if out[i][k] != out[j][k] {
				return out[i][k] < out[j][k]
			}
		}
		return false
	})
	if len(out) == 0 {
		return nil, fmt.Errorf("delaunay: degenerate input (all points collinear or fewer than 3 distinct points)")
	}
	return out, nil
}

func locatePoint(p, start int, tris []bwTri, pp *pointPreds) int {
	cur := start
	limit := 2 * len(tris)
	if limit < 64 {
		limit = 64
	}
	for step := 0; step < limit; step++ {
		tr := &tris[cur]
		next := -1
		for j := 0; j < 3; j++ {
			if pp.orient(tr.v[j], tr.v[(j+1)%3], p) < 0 {
				next = tr.nb[j]
				break
			}
		}
		if next < 0 {
			return cur
		}
		if next >= len(tris) {
			break
		}
		cur = next
	}
	for i := range tris {
		tr := &tris[i]
		if !tr.alive {
			continue
		}
		inside := true
		for j := 0; j < 3; j++ {
			if pp.orient(tr.v[j], tr.v[(j+1)%3], p) < 0 {
				inside = false
				break
			}
		}
		if inside {
			return i
		}
	}
	return 0
}
