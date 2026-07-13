package waffle

import (
	"math"

	"github.com/flywave/go3d/float64/vec2"
)

type triangleIndex struct {
	cx, cy float64
	pts    [3]int
}

func nearestInterp(x, y float64, pts []vec2.T, zs []float64) float64 {
	if len(pts) == 0 {
		return math.NaN()
	}
	minDist := math.MaxFloat64
	bestZ := 0.0
	for i, pt := range pts {
		dx := x - pt[0]
		dy := y - pt[1]
		dist := dx*dx + dy*dy
		if dist < minDist {
			minDist = dist
			bestZ = zs[i]
		}
	}
	return bestZ
}

func linearInterpGrid(x, y float64, triList [][3]int, pts []vec2.T, zs []float64,
	gridIdx *triangleGridIndex) float64 {

	candidateTris := gridIdx.findTriangles(x, y)
	if candidateTris == nil {
		candidateTris = make([]int, len(triList))
		for i := range triList {
			candidateTris[i] = i
		}
	}
	for _, ti := range candidateTris {
		if ti >= len(triList) {
			continue
		}
		t := triList[ti]
		p0, p1, p2 := pts[t[0]], pts[t[1]], pts[t[2]]
		z0, z1, z2 := zs[t[0]], zs[t[1]], zs[t[2]]
		val, found := barycentricInterp(x, y, p0, p1, p2, z0, z1, z2)
		if found {
			return val
		}
	}
	return math.NaN()
}

func barycentricInterp(x, y float64, p0, p1, p2 vec2.T, z0, z1, z2 float64) (float64, bool) {
	denom := (p1[1]-p2[1])*(p0[0]-p2[0]) + (p2[0]-p1[0])*(p0[1]-p2[1])
	if math.Abs(denom) < 1e-15 {
		return 0, false
	}
	w0 := ((p1[1]-p2[1])*(x-p2[0]) + (p2[0]-p1[0])*(y-p2[1])) / denom
	w1 := ((p2[1]-p0[1])*(x-p2[0]) + (p0[0]-p2[0])*(y-p2[1])) / denom
	w2 := 1 - w0 - w1
	if w0 < -1e-10 || w1 < -1e-10 || w2 < -1e-10 {
		return 0, false
	}
	return w0*z0 + w1*z1 + w2*z2, true
}

func boundingBox(tri [3]int, pts []vec2.T) [4]float64 {
	xMin := math.Min(pts[tri[0]][0], math.Min(pts[tri[1]][0], pts[tri[2]][0]))
	xMax := math.Max(pts[tri[0]][0], math.Max(pts[tri[1]][0], pts[tri[2]][0]))
	yMin := math.Min(pts[tri[0]][1], math.Min(pts[tri[1]][1], pts[tri[2]][1]))
	yMax := math.Max(pts[tri[0]][1], math.Max(pts[tri[1]][1], pts[tri[2]][1]))
	return [4]float64{xMin, xMax, yMin, yMax}
}

func distSq(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}

func interpolateLaplace(x, y float64, pts []vec2.T, zs []float64, tris []triangleIndex) float64 {
	found, p0, p1, p2, z0, z1, z2 := findContainingTriangle(x, y, pts, zs, tris)
	if !found {
		return math.NaN()
	}
	return laplaceWeightedInterp(x, y, p0, p1, p2, z0, z1, z2)
}

func findContainingTriangle(x, y float64, pts []vec2.T, zs []float64, tris []triangleIndex) (bool, vec2.T, vec2.T, vec2.T, float64, float64, float64) {
	bestTri := -1
	var bestDist float64
	for i, tri := range tris {
		dx := x - tri.cx
		dy := y - tri.cy
		dist := dx*dx + dy*dy
		if bestTri < 0 || dist < bestDist {
			t0, t1, t2 := tri.pts[0], tri.pts[1], tri.pts[2]
			_, inside := barycentricInterp(x, y, pts[t0], pts[t1], pts[t2], zs[t0], zs[t1], zs[t2])
			if inside {
				bestTri = i
				bestDist = dist
			}
		}
	}
	if bestTri >= 0 {
		tri := tris[bestTri]
		t0, t1, t2 := tri.pts[0], tri.pts[1], tri.pts[2]
		return true, pts[t0], pts[t1], pts[t2], zs[t0], zs[t1], zs[t2]
	}
	return false, vec2.T{}, vec2.T{}, vec2.T{}, 0, 0, 0
}

func laplaceWeightedInterp(x, y float64, p0, p1, p2 vec2.T, z0, z1, z2 float64) float64 {
	d0 := math.Sqrt(distSq(x, y, p0[0], p0[1]))
	d1 := math.Sqrt(distSq(x, y, p1[0], p1[1]))
	d2 := math.Sqrt(distSq(x, y, p2[0], p2[1]))
	if d0 < 1e-12 {
		return z0
	}
	if d1 < 1e-12 {
		return z1
	}
	if d2 < 1e-12 {
		return z2
	}

	alpha1 := angleBetween(p0, x, y, p2)
	beta1 := angleBetween(p0, x, y, p1)
	w0 := (math.Tan(alpha1/2) + math.Tan(beta1/2)) / d0

	alpha2 := angleBetween(p1, x, y, p0)
	beta2 := angleBetween(p1, x, y, p2)
	w1 := (math.Tan(alpha2/2) + math.Tan(beta2/2)) / d1

	alpha3 := angleBetween(p2, x, y, p1)
	beta3 := angleBetween(p2, x, y, p0)
	w2 := (math.Tan(alpha3/2) + math.Tan(beta3/2)) / d2

	totalWeight := w0 + w1 + w2
	if totalWeight <= 0 {
		return math.NaN()
	}
	return (w0*z0 + w1*z1 + w2*z2) / totalWeight
}

func angleBetween(p1 vec2.T, x, y float64, p2 vec2.T) float64 {
	ba := math.Atan2(p1[1]-y, p1[0]-x)
	bc := math.Atan2(p2[1]-y, p2[0]-x)
	angle := bc - ba
	for angle > math.Pi {
		angle -= 2 * math.Pi
	}
	for angle < -math.Pi {
		angle += 2 * math.Pi
	}
	return math.Abs(angle)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
