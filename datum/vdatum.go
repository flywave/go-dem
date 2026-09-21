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
	NoData      float64
}

func GenerateTransformGrid(region *dem.Region, epsgIn, epsgOut int, model geoid.VerticalDatum) (*VDatumGrid, error) {
	opts := TransformOptions{
		EpsgIn:  epsgIn,
		EpsgOut: epsgOut,
		Region:  region,
	}
	if model != geoid.HAE && model != geoid.UNKNOWN {
		if EPSGToVerticalDatum(epsgIn) == geoid.HAE {
			opts.GeoidIn = model.ToString()
		}
		if EPSGToVerticalDatum(epsgOut) == geoid.HAE {
			opts.GeoidOut = model.ToString()
		}
	}

	vt := NewVerticalTransform(opts)
	result, err := vt.Run()
	if err != nil {
		return nil, err
	}
	return newVDatumGrid(region, epsgIn, epsgOut, result), nil
}

func GenerateGeoidGrid(region *dem.Region, model geoid.VerticalDatum) (*VDatumGrid, error) {
	data := computeGeoidGrid(region, model)
	size := region.XSize * region.YSize
	unc := make([]float64, size)
	if u := geoidUncertainty(model); u > 0 {
		for i := range unc {
			unc[i] = u
		}
	}

	return &VDatumGrid{
		Data:        data,
		Uncertainty: unc,
		Region:      region,
		Model:       model,
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
	return newVDatumGrid(region, epsgIn, epsgOut, result), nil
}

func newVDatumGrid(region *dem.Region, epsgIn, epsgOut int, result *TransformResult) *VDatumGrid {
	return &VDatumGrid{
		Data:        result.Grid,
		Uncertainty: result.Uncertainty,
		Region:      region,
		SrcEpsg:     epsgIn,
		DstEpsg:     epsgOut,
		NoData:      result.NoData,
	}
}

func (vg *VDatumGrid) noDataValue() float64 {
	if vg.NoData != 0 {
		return vg.NoData
	}
	return dem.DefaultNoData
}

func (vg *VDatumGrid) ApplyToDEM(demData []float64, inverse bool) []float64 {
	result := make([]float64, len(demData))
	noData := vg.noDataValue()

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
	return dem.CreateDEM(vg.Data, vg.Region, path, vg.noDataValue())
}

func (vg *VDatumGrid) WriteUncertainty(path string) error {
	return dem.CreateDEM(vg.Uncertainty, vg.Region, path, vg.noDataValue())
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
