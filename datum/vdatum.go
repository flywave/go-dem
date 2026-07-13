package datum

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
	"github.com/flywave/go-geoid"
	"github.com/flywave/go3d/float64/vec2"
)

type VDatumGrid struct {
	Data      []float64
	Uncertainty []float64
	Region    *dem.Region
	Model     geoid.VerticalDatum
	SrcEpsg   int
	DstEpsg   int
}

func GenerateTransformGrid(region *dem.Region, epsgIn, epsgOut int, model geoid.VerticalDatum) (*VDatumGrid, error) {
	vt := NewVerticalTransform(TransformOptions{
		EpsgIn:    epsgIn,
		EpsgOut:   epsgOut,
		Region:    region,
	})
	result, err := vt.Run()
	if err != nil {
		return nil, err
	}

	return &VDatumGrid{
		Data:    result.Grid,
		Uncertainty: result.Uncertainty,
		Region:  region,
		SrcEpsg: epsgIn,
		DstEpsg: epsgOut,
	}, nil
}

func GenerateGeoidGrid(region *dem.Region, model geoid.VerticalDatum) (*VDatumGrid, error) {
	data := computeGeoidGrid(region, model)
	size := region.XSize * region.YSize

	return &VDatumGrid{
		Data:    data,
		Uncertainty: make([]float64, size),
		Region:  region,
		Model:   model,
	}, nil
}

func MultiStepTransform(region *dem.Region, epsgIn, epsgOut int) (*VDatumGrid, error) {
	vt := NewVerticalTransform(TransformOptions{
		EpsgIn:  epsgIn,
		EpsgOut: epsgOut,
		Region:  region,
	})
	result, err := vt.Run()
	if err != nil {
		return nil, err
	}

	return &VDatumGrid{
		Data:    result.Grid,
		Uncertainty: result.Uncertainty,
		Region:  region,
		SrcEpsg: epsgIn,
		DstEpsg: epsgOut,
	}, nil
}

func (vg *VDatumGrid) ApplyToDEM(demData []float64, inverse bool) []float64 {
	result := make([]float64, len(demData))
	noData := dem.DefaultNoData

	for i := range demData {
		if demData[i] == noData || math.IsNaN(demData[i]) {
			result[i] = noData
			continue
		}
		if i >= len(vg.Data) {
			result[i] = demData[i]
			continue
		}
		if inverse {
			result[i] = demData[i] + vg.Data[i]
		} else {
			result[i] = demData[i] - vg.Data[i]
		}
	}

	return result
}

func (vg *VDatumGrid) Write(path string) error {
	return dem.CreateDEM(vg.Data, vg.Region, path, -9999)
}

func (vg *VDatumGrid) WriteUncertainty(path string) error {
	return dem.CreateDEM(vg.Uncertainty, vg.Region, path, -9999)
}

func EPSGToVerticalDatum(epsg int) geoid.VerticalDatum {
	switch epsg {
	case 3855:
		return geoid.EGM2008
	case 5773:
		return geoid.EGM96
	case 5798:
		return geoid.EGM84
	case 5703, 6360, 8228:
		return geoid.EGM96
	default:
		return geoid.HAE
	}
}

func ResolveTransform(fromEPSG, toEPSG int, region *dem.Region) (*VDatumGrid, error) {
	frameIn := GetFrameByEPSG(fromEPSG)
	frameOut := GetFrameByEPSG(toEPSG)

	if frameIn == nil || frameOut == nil {
		return nil, fmt.Errorf("unsupported EPSG: %d -> %d", fromEPSG, toEPSG)
	}

	if frameIn.Type == frameOut.Type && frameIn.Type == FrameCDN {
		model := EPSGToVerticalDatum(fromEPSG)
		g := geoid.NewGeoid(model, true)
		if g == nil {
			return nil, fmt.Errorf("failed to initialize geoid for EPSG:%d", fromEPSG)
		}

		size := region.XSize * region.YSize
		data := make([]float64, size)
		noData := dem.DefaultNoData

		var srs4326 geo.Proj = geo.NewProj("EPSG:4326")
		needTransform := region.SRS() != nil && !region.SRS().Eq(srs4326)

		for y := 0; y < region.YSize; y++ {
			for x := 0; x < region.XSize; x++ {
				geoX := region.BBox().Min[0] + float64(x)*region.XRes
				geoY := region.BBox().Min[1] + float64(y)*region.YRes

				lon, lat := geoX, geoY
				if needTransform {
					pts := region.SRS().TransformTo(srs4326, []vec2.T{{geoX, geoY}})
					if len(pts) > 0 {
						lon, lat = pts[0][0], pts[0][1]
					}
				}

				und := g.GetHeight(lat, lon)
				if math.IsNaN(und) || math.IsInf(und, 0) {
					data[y*region.XSize+x] = noData
				} else {
					data[y*region.XSize+x] = und
				}
			}
		}

		return &VDatumGrid{
			Data:    data,
			Uncertainty: make([]float64, size),
			Region:  region,
			SrcEpsg: fromEPSG,
			DstEpsg: toEPSG,
		}, nil
	}

	return MultiStepTransform(region, fromEPSG, toEPSG)
}

func SupportedFrames() string {
	return ListFrames()
}
