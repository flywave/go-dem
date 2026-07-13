package datum

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
	"github.com/flywave/go-geoid"
	vec2d "github.com/flywave/go3d/float64/vec2"
)

type VDatumGrid struct {
	Data   []float64
	Region *dem.Region
	Model  geoid.VerticalDatum
}

func GenerateVDatumGrid(region *dem.Region, model geoid.VerticalDatum, cubic bool) (*VDatumGrid, error) {
	if model == geoid.HAE || model == geoid.UNKNOWN {
		return nil, fmt.Errorf("invalid vertical datum model: %s", model.ToString())
	}

	g := geoid.NewGeoid(model, cubic)
	if g == nil {
		return nil, fmt.Errorf("failed to initialize geoid: %s", model.ToString())
	}

	size := region.XSize * region.YSize
	data := make([]float64, size)
	gt := region.GeoTransform()
	w := region.XSize
	h := region.YSize

	var srs4326 geo.Proj
	srs4326 = geo.NewProj("EPSG:4326")
	needTransform := region.SRS() != nil && !region.SRS().Eq(srs4326) && !region.SRS().IsLatLong()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
			geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]

			lon, lat := geoX, geoY
			if needTransform {
				pts := region.SRS().TransformTo(srs4326, []vec2d.T{{geoX, geoY}})
				if len(pts) > 0 {
					lon, lat = pts[0][0], pts[0][1]
				}
			}

			undulation := g.GetHeight(lat, lon)
			if math.IsNaN(undulation) || math.IsInf(undulation, 0) {
				data[y*w+x] = 0
			} else {
				data[y*w+x] = undulation
			}
		}
	}

	return &VDatumGrid{
		Data:   data,
		Region: region,
		Model:  model,
	}, nil
}

func (vg *VDatumGrid) ApplyToDEM(demData []float64, fromEllipsoidal bool) []float64 {
	result := make([]float64, len(demData))
	noData := dem.DefaultNoData

	for i := range demData {
		if demData[i] == noData || math.IsNaN(demData[i]) {
			result[i] = noData
			continue
		}
		if fromEllipsoidal {
			result[i] = demData[i] - vg.Data[i]
		} else {
			result[i] = demData[i] + vg.Data[i]
		}
	}

	return result
}

func (vg *VDatumGrid) Write(path string) error {
	return dem.CreateDEM(vg.Data, vg.Region, path, -9999)
}
