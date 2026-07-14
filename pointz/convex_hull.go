package pointz

import (
	"math"

	"github.com/flywave/go-dem/delaunay"
)

type ConvexHull struct {
	Points   []Point3D
	delaunay [][3]int
}

func computeConvexHull(points []Point3D) ConvexHull {
	if len(points) < 3 {
		return ConvexHull{Points: points}
	}

	start := 0
	for i := 1; i < len(points); i++ {
		if points[i].X < points[start].X || (points[i].X == points[start].X && points[i].Y < points[start].Y) {
			start = i
		}
	}

	result := []Point3D{}
	p := start
	for {
		result = append(result, points[p])
		q := (p + 1) % len(points)
		for i := 0; i < len(points); i++ {
			if orientation(points[p], points[i], points[q]) == 2 {
				q = i
			}
		}
		p = q
		if p == start {
			break
		}
	}

	h := ConvexHull{Points: result}
	h.initDelaunay()
	return h
}

func (h *ConvexHull) initDelaunay() {
	if len(h.Points) < 3 {
		return
	}
	x := make([]float64, len(h.Points))
	y := make([]float64, len(h.Points))
	for i, p := range h.Points {
		x[i] = p.X
		y[i] = p.Y
	}
	tris, _, err := delaunay.Triangulate(x, y)
	if err == nil {
		h.delaunay = tris
	}
}

func orientation(p, q, r Point3D) int {
	val := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)
	if math.Abs(val) < 1e-12 {
		return 0
	}
	if val > 0 {
		return 1
	}
	return 2
}

func (h *ConvexHull) KeepPointsInside(pts []Point3D) []Point3D {
	if len(h.Points) < 3 {
		return pts
	}
	xMin, xMax := h.Points[0].X, h.Points[0].X
	yMin, yMax := h.Points[0].Y, h.Points[0].Y
	for _, p := range h.Points[1:] {
		if p.X < xMin {
			xMin = p.X
		}
		if p.X > xMax {
			xMax = p.X
		}
		if p.Y < yMin {
			yMin = p.Y
		}
		if p.Y > yMax {
			yMax = p.Y
		}
	}

	var result []Point3D
	for _, p := range pts {
		if p.X < xMin || p.X > xMax || p.Y < yMin || p.Y > yMax {
			continue
		}
		if h.pointInConvexPolygon(p) {
			result = append(result, p)
		}
	}
	return result
}

func (h *ConvexHull) CalculateMask(pts []Point3D) []bool {
	mask := make([]bool, len(pts))
	if len(h.Points) < 3 {
		return mask
	}
	xMin, xMax := h.Points[0].X, h.Points[0].X
	yMin, yMax := h.Points[0].Y, h.Points[0].Y
	for _, p := range h.Points[1:] {
		if p.X < xMin {
			xMin = p.X
		}
		if p.X > xMax {
			xMax = p.X
		}
		if p.Y < yMin {
			yMin = p.Y
		}
		if p.Y > yMax {
			yMax = p.Y
		}
	}
	for i, p := range pts {
		if p.X < xMin || p.X > xMax || p.Y < yMin || p.Y > yMax {
			continue
		}
		mask[i] = h.pointInConvexPolygon(p)
	}
	return mask
}

func (h *ConvexHull) findSimplex(pt Point3D) int {
	if len(h.delaunay) == 0 {
		return -1
	}
	for i, tr := range h.delaunay {
		a, b, c := h.Points[tr[0]], h.Points[tr[1]], h.Points[tr[2]]
		o1 := orientation(a, b, pt)
		o2 := orientation(b, c, pt)
		o3 := orientation(c, a, pt)
		if o1 >= 0 && o2 >= 0 && o3 >= 0 {
			return i
		}
		if o1 <= 0 && o2 <= 0 && o3 <= 0 {
			return i
		}
	}
	return -1
}

func (h *ConvexHull) pointInConvexPolygon(pt Point3D) bool {
	if len(h.delaunay) > 0 {
		return h.findSimplex(pt) >= 0
	}
	n := len(h.Points)
	prev := 0
	for i := 0; i < n; i++ {
		cur := orientation(h.Points[i], h.Points[(i+1)%n], pt)
		if cur != 0 {
			if prev == 0 {
				prev = cur
			} else if cur != prev {
				return false
			}
		}
	}
	return true
}

func (h *ConvexHull) Bounds() (xMin, xMax, yMin, yMax float64) {
	if len(h.Points) == 0 {
		return 0, 0, 0, 0
	}
	xMin, xMax = h.Points[0].X, h.Points[0].X
	yMin, yMax = h.Points[0].Y, h.Points[0].Y
	for _, p := range h.Points[1:] {
		if p.X < xMin {
			xMin = p.X
		}
		if p.X > xMax {
			xMax = p.X
		}
		if p.Y < yMin {
			yMin = p.Y
		}
		if p.Y > yMax {
			yMax = p.Y
		}
	}
	return
}
