package datum

import (
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
	"github.com/flywave/go-geoid"
	"github.com/flywave/go3d/float64/vec2"
)

type TransformOptions struct {
	EpsgIn    int
	EpsgOut   int
	GeoidIn   string
	GeoidOut  string
	Region    *dem.Region
	NoData    float64
}

type TransformResult struct {
	Grid        []float64
	Uncertainty []float64
	EpsgOut     int
}

type VerticalTransform struct {
	opts     TransformOptions
	xCount   int
	yCount   int
	geoTrans [6]float64
}

func NewVerticalTransform(opts TransformOptions) *VerticalTransform {
	region := opts.Region
	gt := region.GeoTransform()
	return &VerticalTransform{
		opts:     opts,
		xCount:   region.XSize,
		yCount:   region.YSize,
		geoTrans: gt,
	}
}

func (vt *VerticalTransform) Run() (*TransformResult, error) {
	grid, unc, outEpsg := vt.verticalTransform(vt.opts.EpsgIn, vt.opts.EpsgOut)
	return &TransformResult{
		Grid:        grid,
		Uncertainty: unc,
		EpsgOut:     outEpsg,
	}, nil
}

func (vt *VerticalTransform) verticalTransform(epsgIn, epsgOut int) ([]float64, []float64, int) {
	n := vt.xCount * vt.yCount
	transArray := make([]float64, n)
	uncArray := make([]float64, n)

	noData := vt.opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}
	_ = noData
	_ = uncArray

	if epsgIn == epsgOut {
		return transArray, uncArray, epsgOut
	}

	frameIn := GetFrameByEPSG(epsgIn)
	frameOut := GetFrameByEPSG(epsgOut)

	if frameIn == nil || frameOut == nil {
		return transArray, uncArray, epsgOut
	}

	currentEpsg := epsgIn
	_ = currentEpsg

	var steps []transformStep

	if frameIn.Type == FrameTidal && frameOut.Type == FrameTidal {
		steps = append(steps, transformStep{from: epsgIn, to: 5714, via: "tidal2msl"})
		steps = append(steps, transformStep{from: 5714, to: epsgOut, via: "msl2tidal"})
	} else if frameIn.Type == FrameTidal {
		steps = append(steps, transformStep{from: epsgIn, to: 5714, via: "tidal2msl"})
		if frameOut.Type == FrameCDN {
			steps = append(steps, transformStep{from: 5714, to: epsgOut, via: "geoid"})
		} else {
			steps = append(steps, transformStep{from: 5714, to: epsgOut, via: "geoid"})
		}
	} else if frameOut.Type == FrameTidal {
		if frameIn.Type == FrameCDN {
			steps = append(steps, transformStep{from: epsgIn, to: 5714, via: "geoid_invert"})
		}
		steps = append(steps, transformStep{from: 5714, to: epsgOut, via: "msl2tidal"})
	} else if frameIn.Type == FrameCDN && frameOut.Type == FrameCDN {
		steps = append(steps, transformStep{from: epsgIn, to: epsgOut, via: "cdn2cdn"})
	} else if frameIn.Type == FrameHTDP || frameOut.Type == FrameHTDP {
		if frameIn.Type == FrameCDN {
			steps = append(steps, transformStep{from: epsgIn, to: 7912, via: "cdn2ellipsoid"})
		}
		if frameIn.Type == FrameTidal {
			steps = append(steps, transformStep{from: epsgIn, to: 5714, via: "tidal2msl"})
			steps = append(steps, transformStep{from: 5714, to: 7912, via: "geoid_invert"})
		}
		if frameIn.Type == FrameHTDP && frameOut.Type == FrameHTDP {
			steps = append(steps, transformStep{from: epsgIn, to: epsgOut, via: "htdp2htdp"})
		} else if frameIn.Type == FrameHTDP && frameOut.Type == FrameCDN {
			steps = append(steps, transformStep{from: epsgIn, to: 7912, via: "htdp2ellipsoid"})
			steps = append(steps, transformStep{from: 7912, to: epsgOut, via: "ellipsoid2cdn"})
		}
		if frameOut.Type == FrameTidal {
			steps = append(steps, transformStep{from: 7912, to: 5714, via: "geoid"})
			steps = append(steps, transformStep{from: 5714, to: epsgOut, via: "msl2tidal"})
		}
	} else if frameIn.Type == FrameCDN && frameOut.Type == FrameHTDP {
		steps = append(steps, transformStep{from: epsgIn, to: 7912, via: "cdn2ellipsoid"})
		steps = append(steps, transformStep{from: 7912, to: epsgOut, via: "ellipsoid2htdp"})
	} else {
		steps = append(steps, transformStep{from: epsgIn, to: epsgOut, via: "direct"})
	}

	_ = steps

	transArray = computeGeoidGrid(vt.opts.Region, geoid.EGM96)

	return transArray, uncArray, epsgOut
}

type transformStep struct {
	from int
	to   int
	via  string
}

func computeGeoidGrid(region *dem.Region, model geoid.VerticalDatum) []float64 {
	g := geoid.NewGeoid(model, true)
	n := region.XSize * region.YSize
	grid := make([]float64, n)
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
				grid[y*region.XSize+x] = noData
			} else {
				grid[y*region.XSize+x] = und
			}
		}
	}
	return grid
}
