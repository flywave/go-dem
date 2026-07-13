package datum

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
	"github.com/flywave/go-geoid"
	vec2d "github.com/flywave/go3d/float64/vec2"
)

type HeightType int

const (
	HeightEllipsoidal HeightType = iota
	HeightOrthometric
	HeightGeoid
)

type ConvertOptions struct {
	From HeightType
	To   HeightType
	Model geoid.VerticalDatum
	Cubic bool
}

func ConvertHeight(data []float64, region *dem.Region, opts *ConvertOptions) ([]float64, error) {
	model := opts.Model
	if model == geoid.UNKNOWN || model == geoid.HAE {
		return data, nil
	}

	g := geoid.NewGeoid(model, opts.Cubic)
	if g == nil {
		return nil, fmt.Errorf("failed to create geoid model: %s", model.ToString())
	}

	result := make([]float64, len(data))
	noData := dem.DefaultNoData

	gt := region.GeoTransform()
	w := region.XSize
	h := region.YSize

	var srs4326 geo.Proj
	srs4326 = geo.NewProj("EPSG:4326")
	needTransform := region.SRS() != nil && !region.SRS().Eq(srs4326) && !region.SRS().IsLatLong()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			z := data[idx]
			if z == noData || math.IsNaN(z) {
				result[idx] = noData
				continue
			}

			geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
			geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]

			lon, lat := geoX, geoY
			if needTransform {
				pts := region.SRS().TransformTo(srs4326, []vec2d.T{{geoX, geoY}})
				if len(pts) > 0 {
					lon, lat = pts[0][0], pts[0][1]
				}
			}

			var resultZ float64
			switch {
			case opts.From == HeightEllipsoidal && opts.To == HeightOrthometric:
				resultZ = g.ConvertHeight(lat, lon, z, geoid.ELLIPSOIDTOGEOID)
			case opts.From == HeightOrthometric && opts.To == HeightEllipsoidal:
				resultZ = g.ConvertHeight(lat, lon, z, geoid.GEOIDTOELLIPSOID)
			default:
				resultZ = z
			}

			if math.IsNaN(resultZ) || math.IsInf(resultZ, 0) {
				result[idx] = noData
			} else {
				result[idx] = resultZ
			}
		}
	}

	return result, nil
}

func OrthometricToEllipsoidal(data []float64, region *dem.Region, model geoid.VerticalDatum) ([]float64, error) {
	return ConvertHeight(data, region, &ConvertOptions{
		From:  HeightOrthometric,
		To:    HeightEllipsoidal,
		Model: model,
	})
}

func EllipsoidalToOrthometric(data []float64, region *dem.Region, model geoid.VerticalDatum) ([]float64, error) {
	return ConvertHeight(data, region, &ConvertOptions{
		From:  HeightEllipsoidal,
		To:    HeightOrthometric,
		Model: model,
	})
}

func WGS84ToMSL(lon, lat, h float64, model geoid.VerticalDatum) float64 {
	g := geoid.NewGeoid(model, true)
	if g == nil {
		return h
	}
	return g.ConvertHeight(lat, lon, h, geoid.ELLIPSOIDTOGEOID)
}

func MSLToWGS84(lon, lat, h float64, model geoid.VerticalDatum) float64 {
	g := geoid.NewGeoid(model, true)
	if g == nil {
		return h
	}
	return g.ConvertHeight(lat, lon, h, geoid.GEOIDTOELLIPSOID)
}
