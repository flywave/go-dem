package datalist

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
	"github.com/flywave/go3d/float64/vec2"
)

type XYZPoint struct {
	X, Y, Z float64
	Intensity float64
	Quality   float64
}

type XYZFile struct {
	Points  []XYZPoint
	Bounds  vec2.Rect
	SRS     geo.Proj
	NoData  float64
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

	xf.Bounds = vec2.Rect{
		Min: vec2.T{minX, minY},
		Max: vec2.T{maxX, maxY},
	}

	return xf, scanner.Err()
}

func (xf *XYZFile) ToDEM(region *dem.Region, method string) ([]float64, error) {
	pts := make([]vec2.T, len(xf.Points))
	zs := make([]float64, len(xf.Points))

	needTransform := xf.SRS != nil && region.SRS() != nil &&
		!xf.SRS.Eq(region.SRS())

	for i, p := range xf.Points {
		if needTransform {
			transformed := xf.SRS.TransformTo(region.SRS(), []vec2.T{{p.X, p.Y}})
			if len(transformed) > 0 {
				pts[i] = transformed[0]
			} else {
				pts[i] = vec2.T{p.X, p.Y}
			}
		} else {
			pts[i] = vec2.T{p.X, p.Y}
		}
		zs[i] = p.Z
	}

	_ = method
	gt := region.GeoTransform()
	w, h := region.XSize, region.YSize
	noData := xf.NoData

	result := make([]float64, w*h)
	for i := range result {
		result[i] = noData
	}

	for i, pt := range pts {
		px := int(math.Round((pt[0] - gt[0]) / gt[1]))
		py := int(math.Round((pt[1] - gt[3]) / gt[5]))
		if px >= 0 && px < w && py >= 0 && py < h {
			result[py*w+px] = zs[i]
		}
	}

	return result, nil
}

func (xf *XYZFile) PointCount() int {
	return len(xf.Points)
}

func (xf *XYZFile) BBoxString() string {
	return fmt.Sprintf("%.6f/%.6f/%.6f/%.6f",
		xf.Bounds.Min[0], xf.Bounds.Max[0],
		xf.Bounds.Min[1], xf.Bounds.Max[1])
}
