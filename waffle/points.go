package waffle

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

func PointsFromRaster(path string) ([]Point, error) {
	data, region, err := dem.ReadDEM(path)
	if err != nil {
		return nil, fmt.Errorf("read raster %s: %v", path, err)
	}

	gt := region.GeoTransform()
	noData := dem.DefaultNoData
	pts := make([]Point, 0, len(data)/4)

	for y := 0; y < region.YSize; y++ {
		for x := 0; x < region.XSize; x++ {
			z := data[y*region.XSize+x]
			if z == noData || math.IsNaN(z) {
				continue
			}
			geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
			geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]
			pts = append(pts, Point{
				Position: vec2.T{geoX, geoY},
				Z:        z,
			})
		}
	}
	return pts, nil
}

func PointsFromMultiple(sources []string) ([]Point, error) {
	var all []Point
	for _, src := range sources {
		pts, err := PointsFromRaster(src)
		if err != nil {
			return nil, err
		}
		all = append(all, pts...)
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("no valid points found in sources")
	}
	return all, nil
}
