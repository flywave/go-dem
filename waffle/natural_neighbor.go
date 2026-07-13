package waffle

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-delaunay"
	"github.com/flywave/go3d/float64/vec2"
)

type naturalNeighborWaffle struct {
	baseWaffle
}

func init() {
	Register(dem.MethodNaturalNeighbor, func() Waffle {
		return &naturalNeighborWaffle{baseWaffle: baseWaffle{name: string(dem.MethodNaturalNeighbor)}}
	})
}

type triangleIndex struct {
	cx, cy float64
	pts    [3]int
}

func (nw *naturalNeighborWaffle) Run(sources []string, opts *Options) (*Result, error) {
	pts, zs, err := collectPoints(sources)
	if err != nil {
		return nil, err
	}
	if len(pts) < 3 {
		return nil, fmt.Errorf("need at least 3 points, got %d", len(pts))
	}

	region := opts.Region
	if region.XSize <= 0 || region.YSize <= 0 {
		region.XSize = int(math.Round((region.BBox().Max[0] - region.BBox().Min[0]) / region.XRes))
		region.YSize = int(math.Round((region.BBox().Max[1] - region.BBox().Min[1]) / region.YRes))
	}

	delaunayPts := make([]delaunay.Point, len(pts))
	for i, pt := range pts {
		delaunayPts[i] = delaunay.Point{pt[0], pt[1]}
	}

	tri, err := delaunay.Triangulate(delaunayPts)
	if err != nil {
		return nil, fmt.Errorf("delaunay triangulation failed: %v", err)
	}

	tris := tri.GetTrianglesPointsMap()
	triList := make([]triangleIndex, 0, len(tris))
	for _, ti := range tris {
		if len(ti) != 3 {
			continue
		}
		cx := (pts[ti[0]][0] + pts[ti[1]][0] + pts[ti[2]][0]) / 3
		cy := (pts[ti[0]][1] + pts[ti[1]][1] + pts[ti[2]][1]) / 3
		triList = append(triList, triangleIndex{cx: cx, cy: cy, pts: [3]int{ti[0], ti[1], ti[2]}})
	}

	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	width := region.XSize
	height := region.YSize
	demData := make([]float64, width*height)
	for i := range demData {
		demData[i] = noData
	}

	gt := region.GeoTransform()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
			geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]

			val := interpolateLaplace(geoX, geoY, pts, zs, triList)
			if !math.IsNaN(val) {
				demData[y*width+x] = val
			}
		}
	}

	return &Result{DEM: demData, Region: region}, nil
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
			_, inside := barycentricInterp(x, y,
				pts[t0], pts[t1], pts[t2],
				zs[t0], zs[t1], zs[t2])
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

func distSq(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}
