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

func (nw *naturalNeighborWaffle) Run(points []Point, opts *Options) (*Result, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("need at least 3 points, got %d", len(points))
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
	triList := make([]triangleIndex, 0, len(triMap))
	for _, ti := range triMap {
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
