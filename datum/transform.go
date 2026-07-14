package datum

import (
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
	"github.com/flywave/go-geoid"
	"github.com/flywave/go3d/float64/vec2"
)

type TransformOptions struct {
	EpsgIn   int
	EpsgOut  int
	GeoidIn  string
	GeoidOut string
	Region   *dem.Region
	NoData   float64
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

type transformStep struct {
	from int
	to   int
	via  string
}

func (vt *VerticalTransform) verticalTransform(epsgIn, epsgOut int) ([]float64, []float64, int) {
	n := vt.xCount * vt.yCount
	transArray := make([]float64, n)
	uncArray := make([]float64, n)

	if epsgIn == epsgOut {
		return transArray, uncArray, epsgOut
	}

	frameIn := GetFrameByEPSG(epsgIn)
	frameOut := GetFrameByEPSG(epsgOut)
	if frameIn == nil || frameOut == nil {
		return transArray, uncArray, epsgOut
	}

	baseUnc := frameIn.Uncertainty
	if frameOut.Uncertainty > baseUnc {
		baseUnc = frameOut.Uncertainty
	}
	for i := range uncArray {
		uncArray[i] = baseUnc
	}

	steps := planSteps(epsgIn, epsgOut, frameIn, frameOut)

	currentEpsg := epsgIn
	for _, step := range steps {
		stepGrid := vt.executeStep(step, currentEpsg)
		for i := range transArray {
			transArray[i] += stepGrid[i]
		}
		currentEpsg = step.to

		stepUnc := FrameUncertainty(step.to)
		if stepUnc > 0 {
			for i := range uncArray {
				uncArray[i] = math.Sqrt(uncArray[i]*uncArray[i] + stepUnc*stepUnc)
			}
		}
	}

	return transArray, uncArray, currentEpsg
}

func planSteps(epsgIn, epsgOut int, frameIn, frameOut *Frame) []transformStep {
	var steps []transformStep

	switch {
	case frameIn.Type == FrameTidal && frameOut.Type == FrameTidal:
		steps = append(steps, transformStep{from: epsgIn, to: 5714, via: "tidal2msl"})
		steps = append(steps, transformStep{from: 5714, to: epsgOut, via: "msl2tidal"})

	case frameIn.Type == FrameTidal && frameOut.Type == FrameCDN:
		steps = append(steps, transformStep{from: epsgIn, to: 5714, via: "tidal2msl"})
		steps = append(steps, transformStep{from: 5714, to: epsgOut, via: "msl2geoid"})

	case frameIn.Type == FrameCDN && frameOut.Type == FrameTidal:
		steps = append(steps, transformStep{from: epsgIn, to: 5714, via: "geoid2msl"})
		steps = append(steps, transformStep{from: 5714, to: epsgOut, via: "msl2tidal"})

	case frameIn.Type == FrameCDN && frameOut.Type == FrameCDN:
		steps = append(steps, transformStep{from: epsgIn, to: epsgOut, via: "cdn2cdn"})

	case frameIn.Type == FrameHTDP || frameOut.Type == FrameHTDP:
		if frameIn.Type == FrameCDN {
			steps = append(steps, transformStep{from: epsgIn, to: 7912, via: "cdn2ellipsoid"})
		}
		if frameIn.Type == FrameTidal {
			steps = append(steps, transformStep{from: epsgIn, to: 5714, via: "tidal2msl"})
			steps = append(steps, transformStep{from: 5714, to: 7912, via: "msl2ellipsoid"})
		}
		if frameIn.Type == FrameHTDP && frameOut.Type == FrameHTDP {
			steps = append(steps, transformStep{from: epsgIn, to: epsgOut, via: "htdp2htdp"})
		} else if frameIn.Type == FrameHTDP && frameOut.Type == FrameCDN {
			steps = append(steps, transformStep{from: epsgIn, to: 7912, via: "htdp2ellipsoid"})
			steps = append(steps, transformStep{from: 7912, to: epsgOut, via: "ellipsoid2cdn"})
		}
		if frameOut.Type == FrameTidal {
			steps = append(steps, transformStep{from: 7912, to: 5714, via: "ellipsoid2msl"})
			steps = append(steps, transformStep{from: 5714, to: epsgOut, via: "msl2tidal"})
		}

	case frameIn.Type == FrameCDN && frameOut.Type == FrameHTDP:
		steps = append(steps, transformStep{from: epsgIn, to: 7912, via: "cdn2ellipsoid"})
		steps = append(steps, transformStep{from: 7912, to: epsgOut, via: "ellipsoid2htdp"})
	}

	return steps
}

func (vt *VerticalTransform) executeStep(step transformStep, currentEpsg int) []float64 {
	switch step.via {
	case "tidal2msl":
		return vt.computeGeoidGrid(geoid.EGM96)
	case "msl2tidal":
		return invertGrid(vt.computeGeoidGrid(geoid.EGM96))
	case "msl2geoid":
		return vt.computeGeoidGrid(geoid.EGM96)
	case "geoid2msl":
		return invertGrid(vt.computeGeoidGrid(geoid.EGM96))
	case "cdn2cdn":
		return vt.cdnTransform(step.from, step.to)
	case "cdn2ellipsoid":
		return invertGrid(vt.computeGeoidGrid(EPSGToVerticalDatum(step.from)))
	case "ellipsoid2cdn":
		return vt.computeGeoidGrid(EPSGToVerticalDatum(step.to))
	case "msl2ellipsoid":
		return vt.computeGeoidGrid(geoid.EGM96)
	case "ellipsoid2msl":
		return invertGrid(vt.computeGeoidGrid(geoid.EGM96))
	case "htdp2htdp":
		return vt.htdpGrid(step.from, step.to)
	case "htdp2ellipsoid":
		return invertGrid(vt.htdpGrid(step.from, 7912))
	case "ellipsoid2htdp":
		return vt.htdpGrid(7912, step.to)
	default:
		return make([]float64, vt.xCount*vt.yCount)
	}
}

func invertGrid(grid []float64) []float64 {
	out := make([]float64, len(grid))
	for i, v := range grid {
		out[i] = -v
	}
	return out
}

func (vt *VerticalTransform) cdnTransform(fromEPSG, toEPSG int) []float64 {
	n := vt.xCount * vt.yCount
	grid := make([]float64, n)
	model := geoid.EGM96

	if fromEPSG == 3855 || fromEPSG == 5773 || fromEPSG == 5798 {
		model = EPSGToVerticalDatum(fromEPSG)
	}
	if toEPSG == 3855 || toEPSG == 5773 || toEPSG == 5798 {
		toModel := EPSGToVerticalDatum(toEPSG)
		if toModel != geoid.HAE {
			g := geoid.NewGeoid(toModel, true)
			if g != nil {
				geoGrid := vt.computeGeoidGrid(toModel)
				fromGrid := vt.computeGeoidGrid(model)
				for i := range grid {
					grid[i] = geoGrid[i] - fromGrid[i]
				}
				return grid
			}
		}
	}

	return vt.computeGeoidGrid(model)
}

func (vt *VerticalTransform) htdpGrid(fromEPSG, toEPSG int) []float64 {
	frameIn := GetFrameByEPSG(fromEPSG)
	frameOut := GetFrameByEPSG(toEPSG)
	if frameIn == nil || frameOut == nil {
		return make([]float64, vt.xCount*vt.yCount)
	}

	gridDef := [6]float64{
		vt.geoTrans[0],
		vt.geoTrans[3],
		vt.geoTrans[0] + float64(vt.xCount)*vt.geoTrans[1],
		vt.geoTrans[3] + float64(vt.yCount)*vt.geoTrans[5],
		float64(vt.xCount),
		float64(vt.yCount),
	}

	srcEpoch := frameIn.Epoch
	if srcEpoch == 0 {
		srcEpoch = 1997.0
	}
	dstEpoch := frameOut.Epoch
	if dstEpoch == 0 {
		dstEpoch = 2000.0
	}

	return cHTDPGrid(gridDef, frameIn.HTDPID, frameOut.HTDPID, srcEpoch, dstEpoch)
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

func (vt *VerticalTransform) computeGeoidGrid(model geoid.VerticalDatum) []float64 {
	return computeGeoidGrid(vt.opts.Region, model)
}
