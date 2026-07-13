package grits

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/flywave/go-dem"
)

type clipFilter struct{ baseGrits }

func init() {
	Register(FilterClip, func() Grits { return &clipFilter{baseGrits{name: string(FilterClip)}} })
}

func (f *clipFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	if opts.PolygonWKT == "" {
		return data, nil
	}

	noData := opts.GetNoData()
	w, h := region.XSize, region.YSize
	result := make([]float64, len(data))
	copy(result, data)

	polygon, xMin, xMax, yMin, yMax, err := parsePolygonWKT(opts.PolygonWKT)
	if err != nil {
		return nil, err
	}

	gt := region.GeoTransform()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			geoX := gt[0] + float64(x)*gt[1]
			geoY := gt[3] + float64(y)*gt[5]

			if geoX < xMin || geoX > xMax || geoY < yMin || geoY > yMax {
				result[y*w+x] = noData
				continue
			}

			if !pointInPolygon(geoX, geoY, polygon) {
				result[y*w+x] = noData
			}
		}
	}

	return result, nil
}

type ring []struct{ x, y float64 }

func parsePolygonWKT(wkt string) (ring, float64, float64, float64, float64, error) {
	var poly ring

	xMin, xMax := math.MaxFloat64, -math.MaxFloat64
	yMin, yMax := math.MaxFloat64, -math.MaxFloat64

	if len(wkt) < 10 {
		return nil, 0, 0, 0, 0, fmt.Errorf("invalid WKT: too short")
	}

	parenStart := -1
	for i := 0; i < len(wkt); i++ {
		if wkt[i] == '(' {
			parenStart = i + 1
			break
		}
	}
	if parenStart < 0 {
		return nil, 0, 0, 0, 0, fmt.Errorf("invalid WKT: no opening paren")
	}

	innerStart := -1
	for i := parenStart; i < len(wkt); i++ {
		if wkt[i] == '(' {
			innerStart = i + 1
			break
		}
	}
	if innerStart < 0 {
		return nil, 0, 0, 0, 0, fmt.Errorf("invalid WKT: no inner opening paren")
	}

	closeParen := -1
	depth := 0
	for i := innerStart; i < len(wkt); i++ {
		if wkt[i] == '(' {
			depth++
		} else if wkt[i] == ')' {
			depth--
			if depth < 0 {
				closeParen = i
				break
			}
		}
	}
	if closeParen < 0 {
		return nil, 0, 0, 0, 0, fmt.Errorf("invalid WKT: no closing paren")
	}

	coordStr := wkt[innerStart:closeParen]
	parts := strings.FieldsFunc(coordStr, func(r rune) bool {
		return r == ',' || r == ' '
	})

	if len(parts) < 6 {
		return nil, 0, 0, 0, 0, fmt.Errorf("need at least 3 coordinate pairs, got %d values", len(parts))
	}

	for i := 0; i+1 < len(parts); i += 2 {
		xi, errX := strconv.ParseFloat(parts[i], 64)
		yi, errY := strconv.ParseFloat(parts[i+1], 64)
		if errX != nil || errY != nil {
			continue
		}
		poly = append(poly, struct{ x, y float64 }{xi, yi})
		if xi < xMin {
			xMin = xi
		}
		if xi > xMax {
			xMax = xi
		}
		if yi < yMin {
			yMin = yi
		}
		if yi > yMax {
			yMax = yi
		}
	}

	if len(poly) < 3 {
		return nil, 0, 0, 0, 0, fmt.Errorf("polygon needs at least 3 points, got %d", len(poly))
	}

	return poly, xMin, xMax, yMin, yMax, nil
}

func pointInPolygon(x, y float64, poly ring) bool {
	inside := false
	n := len(poly)
	j := n - 1
	for i := 0; i < n; i++ {
		if ((poly[i].y > y) != (poly[j].y > y)) &&
			(x < (poly[j].x-poly[i].x)*(y-poly[i].y)/(poly[j].y-poly[i].y)+poly[i].x) {
			inside = !inside
		}
		j = i
	}
	return inside
}
