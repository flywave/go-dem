package waffle

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-delaunay"
	"github.com/flywave/go3d/float64/vec2"
)

type griddataWaffle struct {
	baseWaffle
	method string
}

type gridCell struct {
	triIndices []int
}

type triangleGridIndex struct {
	gridW, gridH int
	xMin, yMin   float64
	xMax, yMax   float64
	cells        []gridCell
}

func buildTriangleGridIndex(triList [][3]int, pts []vec2.T, gridSize int) triangleGridIndex {
	xMin, xMax := pts[0][0], pts[0][0]
	yMin, yMax := pts[0][1], pts[0][1]
	for _, pt := range pts {
		if pt[0] < xMin {
			xMin = pt[0]
		}
		if pt[0] > xMax {
			xMax = pt[0]
		}
		if pt[1] < yMin {
			yMin = pt[1]
		}
		if pt[1] > yMax {
			yMax = pt[1]
		}
	}

	gw := gridSize
	gh := gridSize

	idx := triangleGridIndex{
		gridW: gw, gridH: gh,
		xMin: xMin, yMin: yMin,
		xMax: xMax, yMax: yMax,
		cells: make([]gridCell, gw*gh),
	}
	for i := range idx.cells {
		idx.cells[i].triIndices = make([]int, 0)
	}

	for ti, tri := range triList {
		bbox := boundingBox(tri, pts)
		cxMin := int(math.Floor((bbox[0] - xMin) / (xMax - xMin) * float64(gw)))
		cxMax := int(math.Floor((bbox[1] - xMin) / (xMax - xMin) * float64(gw)))
		cyMin := int(math.Floor((bbox[2] - yMin) / (yMax - yMin) * float64(gh)))
		cyMax := int(math.Floor((bbox[3] - yMin) / (yMax - yMin) * float64(gh)))
		if cxMin < 0 {
			cxMin = 0
		}
		if cxMax >= gw {
			cxMax = gw - 1
		}
		if cyMin < 0 {
			cyMin = 0
		}
		if cyMax >= gh {
			cyMax = gh - 1
		}
		for cy := cyMin; cy <= cyMax; cy++ {
			for cx := cxMin; cx <= cxMax; cx++ {
				idx.cells[cy*gw+cx].triIndices = append(idx.cells[cy*gw+cx].triIndices, ti)
			}
		}
	}
	return idx
}

func (idx *triangleGridIndex) findTriangles(x, y float64) []int {
	cx := int(math.Floor((x - idx.xMin) / (idx.xMax - idx.xMin) * float64(idx.gridW)))
	cy := int(math.Floor((y - idx.yMin) / (idx.yMax - idx.yMin) * float64(idx.gridH)))
	if cx < 0 || cx >= idx.gridW || cy < 0 || cy >= idx.gridH {
		return nil
	}
	return idx.cells[cy*idx.gridW+cx].triIndices
}

func init() {
	Register(dem.MethodLinear, func() Waffle {
		return &griddataWaffle{baseWaffle: baseWaffle{name: string(dem.MethodLinear)}, method: "linear"}
	})
	Register(dem.MethodCubic, func() Waffle {
		return &griddataWaffle{baseWaffle: baseWaffle{name: string(dem.MethodCubic)}, method: "cubic"}
	})
	Register(dem.MethodNearest, func() Waffle {
		return &griddataWaffle{baseWaffle: baseWaffle{name: string(dem.MethodNearest)}, method: "nearest"}
	})
}

func (w *griddataWaffle) Run(points []Point, opts *Options) (*Result, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("need at least 3 points for triangulation, got %d", len(points))
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

	demData := interpolateGrid(region, w.method, triList, pts, zs, noData, &gridIdx)

	return &Result{DEM: demData, Region: region}, nil
}

func interpolateGrid(region *dem.Region, method string, triList [][3]int,
	pts []vec2.T, zs []float64, noData float64, gridIdx *triangleGridIndex) []float64 {

	demData := make([]float64, region.XSize*region.YSize)
	for i := range demData {
		demData[i] = noData
	}

	gt := region.GeoTransform()
	width := region.XSize
	height := region.YSize

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
			geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]

			val := interpAtPoint(geoX, geoY, method, triList, pts, zs, gridIdx)
			if !math.IsNaN(val) {
				demData[y*width+x] = val
			}
		}
	}
	return demData
}

func interpAtPoint(geoX, geoY float64, method string, triList [][3]int,
	pts []vec2.T, zs []float64, gridIdx *triangleGridIndex) float64 {

	switch method {
	case "nearest":
		return nearestInterp(geoX, geoY, pts, zs)
	default:
		return linearInterpGrid(geoX, geoY, triList, pts, zs, gridIdx)
	}
}
