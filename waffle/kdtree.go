package waffle

import (
	"container/heap"
	"math"
	"sort"

	"github.com/flywave/go3d/float64/vec2"
)

type KDNode struct {
	Point vec2.T
	Index int
	Left  *KDNode
	Right *KDNode
	Axis  int
}

type KDTree struct {
	Root *KDNode
}

func NewKDTree(pts []vec2.T) *KDTree {
	if len(pts) == 0 {
		return &KDTree{}
	}
	indices := make([]int, len(pts))
	for i := range indices {
		indices[i] = i
	}
	return &KDTree{Root: buildKDTree(pts, indices, 0)}
}

type kdSortable struct {
	pts  []vec2.T
	idx  []int
	axis int
}

func (s kdSortable) Len() int           { return len(s.idx) }
func (s kdSortable) Less(i, j int) bool { return s.pts[s.idx[i]][s.axis] < s.pts[s.idx[j]][s.axis] }
func (s kdSortable) Swap(i, j int)      { s.idx[i], s.idx[j] = s.idx[j], s.idx[i] }

func buildKDTree(pts []vec2.T, indices []int, depth int) *KDNode {
	if len(indices) == 0 {
		return nil
	}
	axis := depth % 2
	sort.Sort(kdSortable{pts: pts, idx: indices, axis: axis})
	mid := len(indices) / 2

	return &KDNode{
		Point: pts[indices[mid]],
		Index: indices[mid],
		Left:  buildKDTree(pts, indices[:mid], depth+1),
		Right: buildKDTree(pts, indices[mid+1:], depth+1),
		Axis:  axis,
	}
}

type maxHeapEntry struct {
	dist2 float64
	index int
}

type maxHeap []maxHeapEntry

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].dist2 > h[j].dist2 }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(maxHeapEntry)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func (t *KDTree) KNN(q vec2.T, k int) ([]int, []float64) {
	if t.Root == nil || k <= 0 {
		return nil, nil
	}
	h := &maxHeap{}
	heap.Init(h)
	t.knnSearch(t.Root, q, k, 0, h)

	n := h.Len()
	idxs := make([]int, n)
	dists := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		entry := heap.Pop(h).(maxHeapEntry)
		idxs[i] = entry.index
		dists[i] = math.Sqrt(entry.dist2)
	}
	return idxs, dists
}

func (t *KDTree) knnSearch(node *KDNode, q vec2.T, k int, depth int, h *maxHeap) {
	if node == nil {
		return
	}

	dx := q[0] - node.Point[0]
	dy := q[1] - node.Point[1]
	dist2 := dx*dx + dy*dy

	if h.Len() < k {
		heap.Push(h, maxHeapEntry{dist2, node.Index})
	} else if dist2 < (*h)[0].dist2 {
		heap.Pop(h)
		heap.Push(h, maxHeapEntry{dist2, node.Index})
	}

	axis := depth % 2
	diff := q[axis] - node.Point[axis]

	var first, second *KDNode
	if diff < 0 {
		first, second = node.Left, node.Right
	} else {
		first, second = node.Right, node.Left
	}

	t.knnSearch(first, q, k, depth+1, h)

	if h.Len() < k || diff*diff < (*h)[0].dist2 {
		t.knnSearch(second, q, k, depth+1, h)
	}
}

func (t *KDTree) RadiusSearch(q vec2.T, radius float64) ([]int, []float64) {
	if t.Root == nil || radius <= 0 {
		return nil, nil
	}
	radius2 := radius * radius
	var results []neighborEntry
	t.radiusSearch(t.Root, q, radius2, 0, &results)

	idxs := make([]int, len(results))
	dists := make([]float64, len(results))
	for i, r := range results {
		idxs[i] = r.index
		dists[i] = math.Sqrt(r.dist2)
	}
	return idxs, dists
}

type neighborEntry struct {
	index int
	dist2 float64
}

func (t *KDTree) radiusSearch(node *KDNode, q vec2.T, radius2 float64, depth int, results *[]neighborEntry) {
	if node == nil {
		return
	}

	dx := q[0] - node.Point[0]
	dy := q[1] - node.Point[1]
	dist2 := dx*dx + dy*dy

	if dist2 <= radius2 {
		*results = append(*results, neighborEntry{node.Index, dist2})
	}

	axis := depth % 2
	diff := q[axis] - node.Point[axis]

	var first, second *KDNode
	if diff < 0 {
		first, second = node.Left, node.Right
	} else {
		first, second = node.Right, node.Left
	}

	t.radiusSearch(first, q, radius2, depth+1, results)
	if diff*diff <= radius2 {
		t.radiusSearch(second, q, radius2, depth+1, results)
	}
}

func (t *KDTree) Points() []vec2.T {
	var pts []vec2.T
	var idxs []int
	t.collectPoints(t.Root, &pts, &idxs)
	return pts
}

func (t *KDTree) collectPoints(node *KDNode, pts *[]vec2.T, idxs *[]int) {
	if node == nil {
		return
	}
	*pts = append(*pts, node.Point)
	*idxs = append(*idxs, node.Index)
	t.collectPoints(node.Left, pts, idxs)
	t.collectPoints(node.Right, pts, idxs)
}
