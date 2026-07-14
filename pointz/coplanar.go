package pointz

import (
	"math"
)

type CoplanarOptions struct {
	Radius       float64
	Threshold    float64
	MinNeighbors int
	Invert       bool
}

func CoplanarFilter(points []Point3D, opts *CoplanarOptions) []bool {
	if len(points) == 0 {
		return nil
	}
	radius := opts.Radius
	if radius <= 0 {
		radius = 10
	}
	threshold := opts.Threshold
	if threshold <= 0 {
		threshold = 0.5
	}
	minN := opts.MinNeighbors
	if minN <= 0 {
		minN = 3
	}

	pts := make([]vec2, len(points))
	for i, p := range points {
		pts[i] = vec2{p.X, p.Y}
	}
	tree := newKDTree2D(pts)

	mask := make([]bool, len(points))
	for i, p := range points {
		neighbors, _ := tree.radiusSearch(p.X, p.Y, radius)
		if len(neighbors) < minN+1 {
			mask[i] = true
			continue
		}

		var sumX, sumY, sumZ, sumXX, sumYY, sumXY, sumXZ, sumYZ float64
		count := 0
		for _, ni := range neighbors {
			np := points[ni]
			dx := np.X - p.X
			dy := np.Y - p.Y
			sumX += dx
			sumY += dy
			sumZ += np.Z
			sumXX += dx * dx
			sumYY += dy * dy
			sumXY += dx * dy
			sumXZ += dx * np.Z
			sumYZ += dy * np.Z
			count++
		}

		if count < 3 {
			mask[i] = true
			continue
		}
		n := float64(count)
		sxx := sumXX - sumX*sumX/n
		syy := sumYY - sumY*sumY/n
		sxy := sumXY - sumX*sumY/n
		sxz := sumXZ - sumX*sumZ/n
		syz := sumYZ - sumY*sumZ/n

		det := sxx*syy - sxy*sxy
		if math.Abs(det) < 1e-12 {
			mask[i] = true
			continue
		}
		a := (syy*sxz - sxy*syz) / det
		b := (sxx*syz - sxy*sxz) / det
		c := (sumZ - a*sumX - b*sumY) / n
		fittedZ := a*0 + b*0 + c
		dev := math.Abs(p.Z - fittedZ)
		if dev > threshold {
			mask[i] = true
		}
	}

	if opts.Invert {
		for i := range mask {
			mask[i] = !mask[i]
		}
	}
	return mask
}

type vec2 [2]float64

type kdNode2D struct {
	pt    vec2
	idx   int
	left  *kdNode2D
	right *kdNode2D
	axis  int
}

type kdTree2D struct {
	root *kdNode2D
}

func newKDTree2D(pts []vec2) *kdTree2D {
	if len(pts) == 0 {
		return &kdTree2D{}
	}
	idxs := make([]int, len(pts))
	for i := range idxs {
		idxs[i] = i
	}
	return &kdTree2D{root: buildKDTree2D(pts, idxs, 0)}
}

func buildKDTree2D(pts []vec2, idxs []int, depth int) *kdNode2D {
	if len(idxs) == 0 {
		return nil
	}
	axis := depth % 2
	sortIdxsByAxis(pts, idxs, axis)
	mid := len(idxs) / 2
	return &kdNode2D{
		pt:    pts[idxs[mid]],
		idx:   idxs[mid],
		left:  buildKDTree2D(pts, idxs[:mid], depth+1),
		right: buildKDTree2D(pts, idxs[mid+1:], depth+1),
		axis:  axis,
	}
}

func sortIdxsByAxis(pts []vec2, idxs []int, axis int) {
	for i := 0; i < len(idxs); i++ {
		for j := i + 1; j < len(idxs); j++ {
			if pts[idxs[i]][axis] > pts[idxs[j]][axis] {
				idxs[i], idxs[j] = idxs[j], idxs[i]
			}
		}
	}
}

func (t *kdTree2D) radiusSearch(qx, qy, radius float64) ([]int, []float64) {
	if t.root == nil || radius <= 0 {
		return nil, nil
	}
	r2 := radius * radius
	var results []struct {
		idx  int
		dist float64
	}
	t.radiusSearchNode(t.root, qx, qy, r2, 0, &results)
	idxs := make([]int, len(results))
	dists := make([]float64, len(results))
	for i, r := range results {
		idxs[i] = r.idx
		dists[i] = r.dist
	}
	return idxs, dists
}

func (t *kdTree2D) radiusSearchNode(node *kdNode2D, qx, qy, r2 float64, depth int, results *[]struct{ idx int; dist float64 }) {
	if node == nil {
		return
	}
	dx := qx - node.pt[0]
	dy := qy - node.pt[1]
	d2 := dx*dx + dy*dy
	if d2 <= r2 {
		*results = append(*results, struct {
			idx  int
			dist float64
		}{node.idx, math.Sqrt(d2)})
	}
	axis := depth % 2
	diff := qx - node.pt[0]
	if axis == 1 {
		diff = qy - node.pt[1]
	}
	var first, second *kdNode2D
	if diff < 0 {
		first, second = node.left, node.right
	} else {
		first, second = node.right, node.left
	}
	t.radiusSearchNode(first, qx, qy, r2, depth+1, results)
	if diff*diff <= r2 {
		t.radiusSearchNode(second, qx, qy, r2, depth+1, results)
	}
}

func (t *kdTree2D) knnSearch(qx, qy float64, k int) []int {
	return nil
}
