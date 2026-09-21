package datum

import (
	"fmt"
	"math"
	"sync"

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
	From  HeightType
	To    HeightType
	Model geoid.VerticalDatum
	Cubic bool
}

var geoidCache sync.Map

func getGeoid(model geoid.VerticalDatum, cubic bool) *geoid.Geoid {
	if !cubic {
		return geoid.NewGeoid(model, false)
	}
	if v, ok := geoidCache.Load(model); ok {
		return v.(*geoid.Geoid)
	}
	v, _ := geoidCache.LoadOrStore(model, geoid.NewGeoid(model, true))
	return v.(*geoid.Geoid)
}

func ConvertHeight(data []float64, region *dem.Region, opts *ConvertOptions) ([]float64, error) {
	if opts == nil {
		return nil, fmt.Errorf("nil convert options")
	}
	if opts.From == opts.To {
		return data, nil
	}
	if opts.Model == geoid.UNKNOWN || opts.Model == geoid.HAE {
		return nil, fmt.Errorf("no geoid model specified for height conversion %d→%d", opts.From, opts.To)
	}

	var flag geoid.ConvertFlag
	switch {
	case opts.From == HeightEllipsoidal && opts.To == HeightOrthometric:
		flag = geoid.ELLIPSOIDTOGEOID
	case opts.From == HeightOrthometric && opts.To == HeightEllipsoidal:
		flag = geoid.GEOIDTOELLIPSOID
	default:
		return nil, fmt.Errorf("unsupported height conversion %d→%d", opts.From, opts.To)
	}

	g := getGeoid(opts.Model, opts.Cubic)

	w := region.XSize
	h := region.YSize
	n := w * h
	if len(data) < n {
		return nil, fmt.Errorf("data length %d smaller than region size %d", len(data), n)
	}

	result := make([]float64, len(data))
	noData := dem.DefaultNoData

	gt := region.GeoTransform()

	var srs4326 geo.Proj = geo.NewProj("EPSG:4326")
	needTransform := region.SRS() != nil && !region.SRS().Eq(srs4326) && !region.SRS().IsLatLong()

	pts := make([]vec2d.T, n)
	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			pts[i] = vec2d.T{
				gt[0] + (float64(x)+0.5)*gt[1] + (float64(y)+0.5)*gt[2],
				gt[3] + (float64(x)+0.5)*gt[4] + (float64(y)+0.5)*gt[5],
			}
			i++
		}
	}
	if needTransform {
		pts = region.SRS().TransformTo(srs4326, pts)
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			z := data[idx]
			if z == noData || math.IsNaN(z) {
				result[idx] = noData
				continue
			}
			p := pts[idx]
			resultZ := g.ConvertHeight(p[1], p[0], z, flag)
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
	g := getGeoid(model, true)
	return g.ConvertHeight(lat, lon, h, geoid.ELLIPSOIDTOGEOID)
}

func MSLToWGS84(lon, lat, h float64, model geoid.VerticalDatum) float64 {
	g := getGeoid(model, true)
	return g.ConvertHeight(lat, lon, h, geoid.GEOIDTOELLIPSOID)
}
