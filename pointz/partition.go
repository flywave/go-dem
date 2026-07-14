package pointz

import (
	"sort"
)

type PartitionPlan string

const (
	PartitionOne       PartitionPlan = "one"
	PartitionUniform   PartitionPlan = "uniform"
	PartitionMedian    PartitionPlan = "median"
)

type Partition struct {
	Points []Point3D
	Bounds BoxBounds
}

type BoxBounds struct {
	XMin, XMax, YMin, YMax float64
}

func (b BoxBounds) Contains(x, y float64) bool {
	return x >= b.XMin && x <= b.XMax && y >= b.YMin && y <= b.YMax
}

func (b BoxBounds) Center() (float64, float64) {
	return (b.XMin + b.XMax) / 2, (b.YMin + b.YMax) / 2
}

func (b BoxBounds) Area() float64 {
	return (b.XMax - b.XMin) * (b.YMax - b.YMin)
}

func (b BoxBounds) DivideByPoint(x, y float64) [4]BoxBounds {
	eps := 1e-8
	return [4]BoxBounds{
		{XMin: b.XMin, XMax: x, YMin: b.YMin, YMax: y},
		{XMin: x + eps, XMax: b.XMax, YMin: b.YMin, YMax: y},
		{XMin: b.XMin, XMax: x, YMin: y + eps, YMax: b.YMax},
		{XMin: x + eps, XMax: b.XMax, YMin: y + eps, YMax: b.YMax},
	}
}

func boxFromPoints(points []Point3D) BoxBounds {
	if len(points) == 0 {
		return BoxBounds{}
	}
	xMin, xMax := points[0].X, points[0].X
	yMin, yMax := points[0].Y, points[0].Y
	for _, p := range points[1:] {
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
	return BoxBounds{XMin: xMin, XMax: xMax, YMin: yMin, YMax: yMax}
}

func filterPointsByBox(points []Point3D, bounds BoxBounds, keepInside bool) []Point3D {
	var result []Point3D
	for _, p := range points {
		inside := bounds.Contains(p.X, p.Y)
		if (keepInside && inside) || (!keepInside && !inside) {
			result = append(result, p)
		}
	}
	return result
}

type quadPartitioner struct {
	plan PartitionPlan
}

func SelectPartitionPlan(plan PartitionPlan, points []Point3D) *quadPartitioner {
	return &quadPartitioner{plan: plan}
}

func (qp *quadPartitioner) Execute(points []Point3D, minPoints, minArea float64) []Partition {
	bounds := boxFromPoints(points)
	return qp.divideUntil(points, bounds, minPoints, minArea)
}

func (qp *quadPartitioner) divideUntil(points []Point3D, bounds BoxBounds, minPoints, minArea float64) []Partition {
	if len(points) == 0 {
		return nil
	}

	cx, cy := bounds.Center()
	if qp.plan == PartitionMedian {
		cx, cy = medianXY(points)
	}

	newBoxes := bounds.DivideByPoint(cx, cy)
	for _, nb := range newBoxes {
		if nb.Area() < minArea {
			return []Partition{{Points: points, Bounds: bounds}}
		}
	}

	var result []Partition
	for _, nb := range newBoxes {
		sub := filterPointsByBox(points, nb, true)
		if len(sub) < int(minPoints) {
			return []Partition{{Points: points, Bounds: bounds}}
		}
		result = append(result, qp.divideUntil(sub, nb, minPoints, minArea)...)
	}
	return result
}

func medianXY(points []Point3D) (float64, float64) {
	n := len(points)
	if n == 0 {
		return 0, 0
	}
	xs := make([]float64, n)
	ys := make([]float64, n)
	for i, p := range points {
		xs[i] = p.X
		ys[i] = p.Y
	}
	sort.Float64s(xs)
	sort.Float64s(ys)
	return xs[n/2], ys[n/2]
}
