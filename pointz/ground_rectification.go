package pointz

import (
	"fmt"
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

// RectifyGround reclassifies and/or extends ground points of a classified
// point cloud. NOT IMPLEMENTED YET: readClassifiedCloud and
// writeClassifiedCloud are placeholders that always return errors, so this
// entry point currently fails at the read step. The internal
// reclassifyCloud/extendCloud logic is exercised through tests only.
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

	plan := opts.ReclassifyPlan
	if plan == "" {
		plan = PartitionMedian
	}
	partitioner := SelectPartitionPlan(plan, ground)
	partitions := partitioner.Execute(ground, opts.MinPoints, opts.MinArea)

	dists := make([]float64, len(ground))
	for j := range dists {
		dists[j] = -1
	}

	groundIndex := make(map[Point3D]int, len(ground))
	for k, gp := range ground {
		groundIndex[gp] = k
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
			k, ok := groundIndex[gp]
			if !ok {
				continue
			}
			dist := plane.AbsDistance(gp)
			if dists[k] < 0 || dist < dists[k] {
				dists[k] = dist
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

	plan := opts.ExtendPlan
	if plan == "" {
		plan = PartitionMedian
	}
	partitioner := SelectPartitionPlan(plan, ground)
	partitions := partitioner.Execute(ground, opts.MinPoints, opts.MinArea)

	groundColored := make([]ClassifiedPoint, 0)
	for _, p := range points {
		if p.Classification == 2 {
			groundColored = append(groundColored, p)
		}
	}

	result := make([]ClassifiedPoint, len(points))
	copy(result, points)
	assigned := make([]bool, len(grid3D))

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
		for _, p := range groundColored {
			if part.Bounds.Contains(p.X, p.Y) {
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
			if assigned[gi] {
				continue
			}
			assigned[gi] = true

			grid3D[gi].Z = plane.ProjectZ(gp.X, gp.Y)
			grid3D[gi].Classification = 2
			grid3D[gi].R = avgR
			grid3D[gi].G = avgG
			grid3D[gi].B = avgB
		}
	}

	for gi, gp := range grid3D {
		if !assigned[gi] {
			continue
		}
		if !bbox.Contains(gp.X, gp.Y) {
			continue
		}
		gp.Classification = 2
		result = append(result, gp)
	}

	return result
}

const maxGridPoints = 5000000

func buildGridForBounds(bounds BoxBounds, hull ConvexHull, cloud []Point3D, distance float64) []Point3D {
	if distance <= 0 {
		distance = 5
	}
	nx := int((bounds.XMax-bounds.XMin)/distance) + 2
	ny := int((bounds.YMax-bounds.YMin)/distance) + 2
	if nx <= 0 || ny <= 0 || nx > maxGridPoints || ny > maxGridPoints || nx*ny > maxGridPoints {
		return nil
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
	if len(inside) == 0 || len(cloud) == 0 {
		return inside
	}

	pts := make([]vec2, len(cloud))
	for i, cp := range cloud {
		pts[i] = vec2{cp.X, cp.Y}
	}
	tree := newKDTree2D(pts)

	var lonely []Point3D
	for _, gp := range inside {
		if tree.radiusCount(gp.X, gp.Y, distance) == 0 {
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
