package datum

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geoid"
)

type VDatumGrid struct {
	Data        []float64
	Uncertainty []float64
	Region      *dem.Region
	Model       geoid.VerticalDatum
	SrcEpsg     int
	DstEpsg     int
}

func GenerateTransformGrid(region *dem.Region, epsgIn, epsgOut int, model geoid.VerticalDatum) (*VDatumGrid, error) {
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
		if vg.Data[i] == noData {
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
	if fromEPSG == toEPSG {
		return &VDatumGrid{
			Data:    make([]float64, region.XSize*region.YSize),
			Region:  region,
			SrcEpsg: fromEPSG,
			DstEpsg: toEPSG,
		}, nil
	}
	return MultiStepTransform(region, fromEPSG, toEPSG)
}

func TransformDEM(demData []float64, region *dem.Region, fromEPSG, toEPSG int) ([]float64, error) {
	grid, err := ResolveTransform(fromEPSG, toEPSG, region)
	if err != nil {
		return nil, fmt.Errorf("resolve transform %d→%d: %v", fromEPSG, toEPSG, err)
	}
	return grid.ApplyToDEM(demData, false), nil
}

func SupportedFrames() string {
	return ListFrames()
}
