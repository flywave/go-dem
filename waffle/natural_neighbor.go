package waffle

import (
	"fmt"
	"math"

	"github.com/flywave/go-delaunay"
	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

// naturalNeighborWaffle evaluates the mean-value ("Laplace") weights over the
// Delaunay triangulation of the input points. A mean-value barycentric basis
// reproduces linear precision, so for a query point inside a single triangle
// the normalized weights are exactly the triangle's barycentric coordinates:
// this method is mathematically equivalent to the "linear" TIN interpolation
// (it is not a Sibson natural-neighbour interpolation). The registration name
// is kept for backward compatibility.
type naturalNeighborWaffle struct {
	baseWaffle
}

func init() {
	Register(dem.MethodNaturalNeighbor, func() Waffle {
		return &naturalNeighborWaffle{baseWaffle: baseWaffle{name: "laplace_interp"}}
	})
}

func (nw *naturalNeighborWaffle) Run(points []Point, opts *Options) (*Result, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("need at least 3 points, got %d", len(points))
	}
	if opts == nil || opts.Region == nil {
		return nil, fmt.Errorf("region is required")
	}

	region := opts.Region
	if region.XSize <= 0 || region.YSize <= 0 {
		region.XSize = int(math.Round((region.BBox().Max[0] - region.BBox().Min[0]) / region.XRes))
		region.YSize = int(math.Round((region.BBox().Max[1] - region.BBox().Min[1]) / region.YRes))
	}

	pts := make([]vec2.T, len(points))
	zs := make([]float64, len(points))
	for i, p := range points {
		pts[i] = p.Position
		zs[i] = p.Z
	}

	delaunayPts := make([]delaunay.Point, len(pts))
	for i, pt := range pts {
		delaunayPts[i] = delaunay.Point{pt[0], pt[1]}
	}

	tri, err := delaunay.Triangulate(delaunayPts)
	if err != nil {
		return nil, fmt.Errorf("delaunay triangulation failed: %v", err)
	}

	triMap := tri.GetTrianglesPointsMap()
	triList := make([][3]int, 0, len(triMap))
	for _, ti := range triMap {
		if len(ti) == 3 {
			triList = append(triList, [3]int{ti[0], ti[1], ti[2]})
		}
	}

	gridSize := int(math.Sqrt(float64(len(triList))))
	if gridSize < 10 {
		gridSize = 10
	}
	if gridSize > 100 {
		gridSize = 100
	}
	gridIdx := buildTriangleGridIndex(triList, pts, gridSize)

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

	stage := string(dem.MethodNaturalNeighbor)
	if err := startRun(opts, stage, height); err != nil {
		return nil, err
	}
	for y := 0; y < height; y++ {
		if err := dem.CheckCtx(opts.Ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := 0; x < width; x++ {
			geoX, geoY := region.PixelCenterGeo(x, y)

			val := math.NaN()
			for _, ti := range gridIdx.findTriangles(geoX, geoY) {
				if ti >= len(triList) {
					continue
				}
				t := triList[ti]
				p0, p1, p2 := pts[t[0]], pts[t[1]], pts[t[2]]
				z0, z1, z2 := zs[t[0]], zs[t[1]], zs[t[2]]
				if _, inside := barycentricInterp(geoX, geoY, p0, p1, p2, z0, z1, z2); inside {
					val = laplaceWeightedInterp(geoX, geoY, p0, p1, p2, z0, z1, z2)
					break
				}
			}
			if !math.IsNaN(val) {
				demData[y*width+x] = val
			}
		}
		dem.ReportProgress(opts.Progress, stage, y+1, height)
	}
	finishRun(opts, stage, height)

	return &Result{DEM: demData, Region: region}, nil
}
