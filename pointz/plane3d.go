package pointz

import (
	"math"
	"math/rand"
	"sort"
)

type Plane3D struct {
	NX, NY, NZ, D float64
}

func (p *Plane3D) IsValid() bool {
	return !math.IsNaN(p.NX) && !math.IsNaN(p.NY) && !math.IsNaN(p.NZ) && !math.IsNaN(p.D)
}

func (p *Plane3D) Normalize() {
	norm := math.Sqrt(p.NX*p.NX + p.NY*p.NY + p.NZ*p.NZ)
	if norm < 1e-12 {
		return
	}
	p.NX /= norm
	p.NY /= norm
	p.NZ /= norm
	p.D /= norm
}

func (p *Plane3D) ProjectZ(x, y float64) float64 {
	if math.Abs(p.NZ) < 1e-12 {
		return 0
	}
	return -(p.D + p.NX*x + p.NY*y) / p.NZ
}

func (p *Plane3D) SignedDistance(pt Point3D) float64 {
	return p.NX*pt.X + p.NY*pt.Y + p.NZ*pt.Z + p.D
}

func (p *Plane3D) AbsDistance(pt Point3D) float64 {
	d := p.SignedDistance(pt)
	if d < 0 {
		return -d
	}
	return d
}

func (p *Plane3D) AngleDeg() float64 {
	cos := math.Abs(p.NZ)
	if cos > 1 {
		cos = 1
	}
	return math.Acos(cos) * 180.0 / math.Pi
}

func fitPlaneLMedS(points []Point3D, iterations int) Plane3D {
	n := len(points)
	if n < 3 {
		return Plane3D{}
	}
	if iterations <= 0 {
		iterations = 200
	}

	best := Plane3D{}
	bestMedian := math.MaxFloat64

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
		nx := (p1.Y-p0.Y)*(p2.Z-p0.Z) - (p1.Z-p0.Z)*(p2.Y-p0.Y)
		ny := (p1.Z-p0.Z)*(p2.X-p0.X) - (p1.X-p0.X)*(p2.Z-p0.Z)
		nz := (p1.X-p0.X)*(p2.Y-p0.Y) - (p1.Y-p0.Y)*(p2.X-p0.X)
		norm := math.Sqrt(nx*nx + ny*ny + nz*nz)
		if norm < 1e-12 {
			continue
		}
		nx /= norm
		ny /= norm
		nz /= norm
		d := -(nx*p0.X + ny*p0.Y + nz*p0.Z)

		dists := make([]float64, n)
		for j, p := range points {
			dists[j] = math.Abs(nx*p.X + ny*p.Y + nz*p.Z + d)
		}
		sort.Float64s(dists)
		median := dists[n/2]

		if median < bestMedian {
			bestMedian = median
			best = Plane3D{NX: nx, NY: ny, NZ: nz, D: d}
		}
	}

	if !best.IsValid() || bestMedian > 1e10 {
		return fitPlaneLeastSquares3D(points)
	}

	best = refitPlaneInliers(points, best, bestMedian*2)
	return best
}

func fitPlaneLeastSquares3D(points []Point3D) Plane3D {
	n := len(points)
	if n < 3 {
		return Plane3D{}
	}
	var cx, cy, cz float64
	for _, p := range points {
		cx += p.X
		cy += p.Y
		cz += p.Z
	}
	cx /= float64(n)
	cy /= float64(n)
	cz /= float64(n)

	var xx, xy, xz, yy, yz, zz float64
	for _, p := range points {
		dx := p.X - cx
		dy := p.Y - cy
		dz := p.Z - cz
		xx += dx * dx
		xy += dx * dy
		xz += dx * dz
		yy += dy * dy
		yz += dy * dz
		zz += dz * dz
	}

	det := xx*yy - xy*xy
	if math.Abs(det) < 1e-12 {
		norm := math.Sqrt(xz*xz + yz*yz + zz*zz)
		if norm < 1e-12 {
			return Plane3D{NX: 0, NY: 0, NZ: 1, D: -cz}
		}
		return Plane3D{NX: 0, NY: 0, NZ: 1, D: -cz}
	}
	nx := (yy*xz - xy*yz) / det
	ny := (xx*yz - xy*xz) / det
	nz := -1.0
	norm := math.Sqrt(nx*nx + ny*ny + nz*nz)
	nx /= norm
	ny /= norm
	nz /= norm
	d := -(nx*cx + ny*cy + nz*cz)
	return Plane3D{NX: nx, NY: ny, NZ: nz, D: d}
}

func refitPlaneInliers(points []Point3D, initial Plane3D, threshold float64) Plane3D {
	var inliers []Point3D
	for _, p := range points {
		if initial.AbsDistance(p) < threshold {
			inliers = append(inliers, p)
		}
	}
	if len(inliers) < 3 {
		return initial
	}
	return fitPlaneLeastSquares3D(inliers)
}
