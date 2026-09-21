package pointz

import (
	"github.com/flywave/go-dem/delaunay"
)

type ConvexHull struct {
	Points      []Point3D
	delaunay    [][3]int
	trisBox     [][4]float64
	delaunayErr error
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
			if i == p {
				continue
			}
			o := orientation(points[p], points[i], points[q])
			if o == 2 || (o == 0 && dist2(points[p], points[i]) > dist2(points[p], points[q])) {
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

func dist2(a, b Point3D) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
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
	if err != nil {
		h.delaunayErr = err
		return
	}
	h.delaunay = tris
	boxes := make([][4]float64, len(tris))
	for i, tr := range tris {
		a, b, c := h.Points[tr[0]], h.Points[tr[1]], h.Points[tr[2]]
		boxes[i] = [4]float64{
			min(a.X, min(b.X, c.X)),
			max(a.X, max(b.X, c.X)),
			min(a.Y, min(b.Y, c.Y)),
			max(a.Y, max(b.Y, c.Y)),
		}
	}
	h.trisBox = boxes
}

func orientation(p, q, r Point3D) int {
	cross := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)
	l1 := (q.X-p.X)*(q.X-p.X) + (q.Y-p.Y)*(q.Y-p.Y)
	l2 := (r.X-q.X)*(r.X-q.X) + (r.Y-q.Y)*(r.Y-q.Y)
	if l1 == 0 || l2 == 0 {
		return 0
	}
	if cross*cross < 1e-24*l1*l2 {
		return 0
	}
	if cross > 0 {
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

// CalculateMask reports for each point whether it lies inside the hull
// (true = inside). Note this is the opposite of the package-wide filter
// convention where mask true means "remove".
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
		bb := h.trisBox[i]
		if pt.X < bb[0] || pt.X > bb[1] || pt.Y < bb[2] || pt.Y > bb[3] {
			continue
		}
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
