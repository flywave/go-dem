package datalist

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/waffle"
	"github.com/flywave/go-geo"
	"github.com/flywave/go3d/float64/vec2"
)

type XYZPoint struct {
	X, Y, Z   float64
	Intensity float64
	Quality   float64
}

type XYZFile struct {
	Points []XYZPoint
	Bounds vec2.Rect
	SRS    geo.Proj
	NoData float64

	mu            sync.Mutex
	transformed   []vec2.T
	transformedTo geo.Proj
}

func isFinite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func ParseXYZFile(path string, srs geo.Proj) (*XYZFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open xyz: %v", err)
	}
	defer f.Close()

	xf := &XYZFile{
		SRS:    srs,
		NoData: dem.DefaultNoData,
	}

	scanner := bufio.NewScanner(f)
	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		x, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			continue
		}
		y, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}
		z, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			continue
		}
		if !isFinite(x) || !isFinite(y) || !isFinite(z) {
			continue
		}

		pt := XYZPoint{X: x, Y: y, Z: z}
		if len(parts) >= 4 {
			if v, err := strconv.ParseFloat(parts[3], 64); err == nil {
				pt.Intensity = v
			}
		}

		xf.Points = append(xf.Points, pt)

		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read xyz: %v", err)
	}

	if len(xf.Points) == 0 {
		return nil, fmt.Errorf("no valid points in xyz file: %s", path)
	}

	xf.Bounds = vec2.Rect{
		Min: vec2.T{minX, minY},
		Max: vec2.T{maxX, maxY},
	}

	return xf, nil
}

func (xf *XYZFile) ToDEM(region *dem.Region, method string) ([]float64, error) {
	if len(xf.Points) == 0 {
		return nil, fmt.Errorf("no points in xyz file")
	}

	pts, err := xf.projectedPoints(region)
	if err != nil {
		return nil, err
	}

	w, err := waffle.New(dem.InterpMethod(strings.TrimSpace(method)))
	if err != nil {
		return nil, fmt.Errorf("toDEM: %v", err)
	}

	points := make([]waffle.Point, len(pts))
	for i, pt := range pts {
		points[i] = waffle.Point{Position: pt, Z: xf.Points[i].Z}
	}

	result, err := w.Run(points, &waffle.Options{
		Region: region,
		NoData: xf.NoData,
	})
	if err != nil {
		return nil, fmt.Errorf("toDEM: %v", err)
	}

	return result.DEM, nil
}

func (xf *XYZFile) projectedPoints(region *dem.Region) ([]vec2.T, error) {
	xf.mu.Lock()
	defer xf.mu.Unlock()

	target := region.SRS()
	if xf.transformed != nil && len(xf.transformed) == len(xf.Points) && sameSRS(xf.transformedTo, target) {
		return xf.transformed, nil
	}

	pts := make([]vec2.T, len(xf.Points))
	for i, p := range xf.Points {
		pts[i] = vec2.T{p.X, p.Y}
	}

	if !isNilProj(xf.SRS) && !isNilProj(target) && !xf.SRS.Eq(target) {
		pts = xf.SRS.TransformTo(target, pts)
		if len(pts) != len(xf.Points) {
			return nil, fmt.Errorf("coordinate transform returned %d of %d points", len(pts), len(xf.Points))
		}
	}

	xf.transformed = pts
	xf.transformedTo = target
	return pts, nil
}

func (xf *XYZFile) PointCount() int {
	return len(xf.Points)
}

func (xf *XYZFile) BBoxString() string {
	return fmt.Sprintf("%.6f/%.6f/%.6f/%.6f",
		xf.Bounds.Min[0], xf.Bounds.Max[0],
		xf.Bounds.Min[1], xf.Bounds.Max[1])
}
