package pointz

import (
	"fmt"
	"math"
)

type RectifyMethod string

const (
	MethodReclassify       RectifyMethod = "reclassify"
	MethodExtend           RectifyMethod = "extend"
	MethodReclassifyExtend RectifyMethod = "reclassify_extend"
)

type GroundRectificationOptions struct {
	InputPath           string
	OutputPath          string
	Method              RectifyMethod
	ReclassifyPlan      PartitionPlan
	ReclassifyThreshold float64
	ExtendPlan          PartitionPlan
	ExtendGridDistance  float64
	MinPoints           float64
	MinArea             float64
}

func DefaultGroundRectificationOptions() *GroundRectificationOptions {
	return &GroundRectificationOptions{
		Method:              MethodReclassifyExtend,
		ReclassifyPlan:      PartitionMedian,
		ReclassifyThreshold: 5,
		ExtendPlan:          PartitionMedian,
		ExtendGridDistance:  5,
		MinPoints:           500,
		MinArea:             750,
	}
}

func RectifyGround(opts *GroundRectificationOptions) error {
	if opts == nil {
		opts = DefaultGroundRectificationOptions()
	}
	if opts.InputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if opts.OutputPath == "" {
		opts.OutputPath = opts.InputPath
	}

	points, err := readClassifiedCloud(opts.InputPath)
	if err != nil {
		return fmt.Errorf("read cloud: %v", err)
	}

	switch opts.Method {
	case MethodReclassify:
		points = reclassifyCloud(points, opts)
	case MethodExtend:
		points = extendCloud(points, opts)
	case MethodReclassifyExtend:
		points = reclassifyCloud(points, opts)
		points = extendCloud(points, opts)
	default:
		return fmt.Errorf("unknown rectify method: %s", opts.Method)
	}

	if err := writeClassifiedCloud(points, opts.OutputPath); err != nil {
		return fmt.Errorf("write cloud: %v", err)
	}
	return nil
}

type ClassifiedPoint struct {
	Point3D
	Classification uint8
	R, G, B        float64
}

func reclassifyCloud(points []ClassifiedPoint, opts *GroundRectificationOptions) []ClassifiedPoint {
	result := make([]ClassifiedPoint, len(points))
	copy(result, points)

	ground := make([]Point3D, 0)
	groundIdx := make([]int, 0)
	for i, p := range points {
		if p.Classification == 2 {
			ground = append(ground, p.Point3D)
			groundIdx = append(groundIdx, i)
		}
	}
	if len(ground) < 3 {
		return result
	}

	partitioner := SelectPartitionPlan(opts.ReclassifyPlan, ground)
	partitions := partitioner.Execute(ground, opts.MinPoints, opts.MinArea)

	dists := make([]float64, len(ground))
	for j := range dists {
		dists[j] = -1
	}

	for _, part := range partitions {
		if len(part.Points) < 3 {
			continue
		}
		plane := fitPlaneLMedS(part.Points, 200)
		if !plane.IsValid() {
			continue
		}
		if plane.AngleDeg() >= 45 {
			continue
		}

		for _, gp := range part.Points {
			for k, gidx := range groundIdx {
				if gp.X == points[gidx].X && gp.Y == points[gidx].Y && gp.Z == points[gidx].Z {
					dist := plane.AbsDistance(gp)
					if dists[k] < 0 || dist < dists[k] {
						dists[k] = dist
					}
					break
				}
			}
		}
	}

	for j, gidx := range groundIdx {
		if dists[j] >= 0 && dists[j] > opts.ReclassifyThreshold {
			result[gidx].Classification = 1
		}
	}

	return result
}

func extendCloud(points []ClassifiedPoint, opts *GroundRectificationOptions) []ClassifiedPoint {
	ground := make([]Point3D, 0)
	for _, p := range points {
		if p.Classification == 2 {
			ground = append(ground, p.Point3D)
		}
	}
	if len(ground) < 3 {
		return points
	}

	bbox := boxFromPoints(ground)
	hull := computeConvexHull(ground)

	grid2D := buildGridForBounds(bbox, hull, ground, opts.ExtendGridDistance)
	if len(grid2D) == 0 {
		return points
	}

	grid3D := make([]ClassifiedPoint, len(grid2D))
	for i, gp := range grid2D {
		grid3D[i] = ClassifiedPoint{
			Point3D: Point3D{X: gp.X, Y: gp.Y},
		}
	}

	partitioner := SelectPartitionPlan(opts.ExtendPlan, ground)
	if opts.ExtendPlan == "" {
		partitioner = SelectPartitionPlan(PartitionMedian, ground)
	}
	partitions := partitioner.Execute(ground, opts.MinPoints, opts.MinArea)

	result := make([]ClassifiedPoint, len(points))
	copy(result, points)
	usedGrid := make(map[int]bool)

	for _, part := range partitions {
		if len(part.Points) < 3 {
			continue
		}
		plane := fitPlaneLMedS(part.Points, 200)
		if !plane.IsValid() {
			continue
		}
		if plane.AngleDeg() >= 45 {
			continue
		}

		var avgR, avgG, avgB, count float64
		for _, p := range points {
			if p.Classification == 2 && part.Bounds.Contains(p.X, p.Y) {
				avgR += p.R
				avgG += p.G
				avgB += p.B
				count++
			}
		}
		if count > 0 {
			avgR /= count
			avgG /= count
			avgB /= count
		}

		for gi, gp := range grid3D {
			if !part.Bounds.Contains(gp.X, gp.Y) {
				continue
			}
			key := int(gp.X*10000 + gp.Y*10000)
			if usedGrid[key] {
				continue
			}
			usedGrid[key] = true

			grid3D[gi].Z = plane.ProjectZ(gp.X, gp.Y)
			grid3D[gi].Classification = 2
			grid3D[gi].R = avgR
			grid3D[gi].G = avgG
			grid3D[gi].B = avgB
		}
	}

	finalBox := boxFromPoints(groundPoints(points))
	for _, gp := range grid3D {
		if gp.Z == 0 {
			continue
		}
		if !finalBox.Contains(gp.X, gp.Y) {
			continue
		}
		gp.Classification = 2
		result = append(result, gp)
	}

	return result
}

func buildGridForBounds(bounds BoxBounds, hull ConvexHull, cloud []Point3D, distance float64) []Point3D {
	if distance <= 0 {
		distance = 5
	}
	var raw []Point3D
	for x := bounds.XMin; x <= bounds.XMax; x += distance {
		for y := bounds.YMin; y <= bounds.YMax; y += distance {
			if bounds.Contains(x, y) {
				raw = append(raw, Point3D{X: x, Y: y})
			}
		}
	}

	inside := hull.KeepPointsInside(raw)

	var lonely []Point3D
	for _, gp := range inside {
		close := false
		for _, cp := range cloud {
			dx := gp.X - cp.X
			dy := gp.Y - cp.Y
			if math.Sqrt(dx*dx+dy*dy) < distance {
				close = true
				break
			}
		}
		if !close {
			lonely = append(lonely, gp)
		}
	}
	return lonely
}

func groundPoints(points []ClassifiedPoint) []Point3D {
	var result []Point3D
	for _, p := range points {
		if p.Classification == 2 {
			result = append(result, p.Point3D)
		}
	}
	return result
}

func readClassifiedCloud(path string) ([]ClassifiedPoint, error) {
	return nil, fmt.Errorf("not implemented: read cloud from %s", path)
}

func writeClassifiedCloud(points []ClassifiedPoint, path string) error {
	return fmt.Errorf("not implemented: write cloud to %s", path)
}
