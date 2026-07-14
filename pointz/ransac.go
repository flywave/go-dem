package pointz

import (
	"math"
	"math/rand"
)

type Plane struct {
	A, B, C float64
}

func FitPlaneRANSAC(points []Point3D, iterations int, threshold float64) Plane {
	n := len(points)
	if n < 3 {
		return Plane{}
	}
	if iterations <= 0 {
		iterations = 100
	}
	if threshold <= 0 {
		threshold = 0.5
	}

	var best Plane
	bestInliers := 0

	for i := 0; i < iterations; i++ {
		idx0 := rand.Intn(n)
		idx1 := rand.Intn(n)
		idx2 := rand.Intn(n)
		for idx1 == idx0 {
			idx1 = rand.Intn(n)
		}
		for idx2 == idx0 || idx2 == idx1 {
			idx2 = rand.Intn(n)
		}

		p0, p1, p2 := points[idx0], points[idx1], points[idx2]
		v1x, v1y, v1z := p1.X-p0.X, p1.Y-p0.Y, p1.Z-p0.Z
		v2x, v2y, v2z := p2.X-p0.X, p2.Y-p0.Y, p2.Z-p0.Z
		nx := v1y*v2z - v1z*v2y
		ny := v1z*v2x - v1x*v2z
		nz := v1x*v2y - v1y*v2x
		norm := math.Sqrt(nx*nx + ny*ny + nz*nz)
		if norm < 1e-12 {
			continue
		}
		if math.Abs(nz) < 1e-12 {
			continue
		}

		a := -nx / nz
		b := -ny / nz
		c := (nx*p0.X + ny*p0.Y + nz*p0.Z) / nz

		inliers := 0
		for _, p := range points {
			pred := a*p.X + b*p.Y + c
			diff := math.Abs(p.Z - pred)
			if diff < threshold {
				inliers++
			}
		}
		if inliers > bestInliers {
			bestInliers = inliers
			best = Plane{A: a, B: b, C: c}
		}
	}

	if bestInliers < 3 {
		return leastSquaresPlane(points)
	}

	best = leastSquaresPlaneInliers(points, best, threshold, bestInliers)
	return best
}

func FitPlaneLeastSquares(points []Point3D) Plane {
	return leastSquaresPlane(points)
}

func leastSquaresPlane(points []Point3D) Plane {
	n := len(points)
	if n < 3 {
		return Plane{}
	}
	var sumX, sumY, sumZ, sumX2, sumY2, sumXY, sumXZ, sumYZ float64
	for _, p := range points {
		sumX += p.X
		sumY += p.Y
		sumZ += p.Z
		sumX2 += p.X * p.X
		sumY2 += p.Y * p.Y
		sumXY += p.X * p.Y
		sumXZ += p.X * p.Z
		sumYZ += p.Y * p.Z
	}
	sxx := sumX2 - sumX*sumX/float64(n)
	syy := sumY2 - sumY*sumY/float64(n)
	sxy := sumXY - sumX*sumY/float64(n)
	sxz := sumXZ - sumX*sumZ/float64(n)
	syz := sumYZ - sumY*sumZ/float64(n)

	det := sxx*syy - sxy*sxy
	if math.Abs(det) < 1e-12 {
		return Plane{}
	}
	a := (syy*sxz - sxy*syz) / det
	b := (sxx*syz - sxy*sxz) / det
	c := (sumZ - a*sumX - b*sumY) / float64(n)
	return Plane{A: a, B: b, C: c}
}

func leastSquaresPlaneInliers(points []Point3D, initial Plane, threshold float64, minInliers int) Plane {
	var inliers []Point3D
	for _, p := range points {
		pred := initial.A*p.X + initial.B*p.Y + initial.C
		if math.Abs(p.Z-pred) < threshold {
			inliers = append(inliers, p)
		}
	}
	if len(inliers) < minInliers {
		return initial
	}
	return leastSquaresPlane(inliers)
}

func (p *Plane) PredictZ(x, y float64) float64 {
	return p.A*x + p.B*y + p.C
}

func (p *Plane) DistanceToPoint(pt Point3D) float64 {
	pred := p.PredictZ(pt.X, pt.Y)
	diff := pt.Z - pred
	if diff < 0 {
		diff = -diff
	}
	return diff
}

func (p *Plane) AngleDeg() float64 {
	denom := math.Sqrt(p.A*p.A + p.B*p.B + 1)
	if denom < 1e-12 {
		return 0
	}
	cos := 1.0 / denom
	if cos > 1 {
		cos = 1
	}
	if cos < -1 {
		cos = -1
	}
	return math.Acos(cos) * 180.0 / math.Pi
}
