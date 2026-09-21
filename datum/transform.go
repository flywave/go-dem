package datum

import (
	"context"
	"fmt"
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
	"github.com/flywave/go-geoid"
	"github.com/flywave/go3d/float64/vec2"
)

const mslEPSG = 5714
const ellipsoidEPSG = 7912

var mslModel = geoid.EGM96

type TransformOptions struct {
	EpsgIn   int
	EpsgOut  int
	GeoidIn  string
	GeoidOut string
	Region   *dem.Region
	NoData   float64
	Progress dem.ProgressFunc
	Ctx      context.Context
}

type TransformResult struct {
	Grid        []float64
	Uncertainty []float64
	EpsgOut     int
	NoData      float64
}

type VerticalTransform struct {
	opts       TransformOptions
	xCount     int
	yCount     int
	geoTrans   [6]float64
	geoidGrids map[geoid.VerticalDatum][]float64
}

func NewVerticalTransform(opts TransformOptions) *VerticalTransform {
	region := opts.Region
	gt := [6]float64{}
	xCount, yCount := 0, 0
	if region != nil {
		gt = region.GeoTransform()
		xCount, yCount = region.XSize, region.YSize
	}
	return &VerticalTransform{
		opts:       opts,
		xCount:     xCount,
		yCount:     yCount,
		geoTrans:   gt,
		geoidGrids: make(map[geoid.VerticalDatum][]float64),
	}
}

func (vt *VerticalTransform) Run() (*TransformResult, error) {
	grid, unc, outEpsg, err := vt.verticalTransform(vt.opts.EpsgIn, vt.opts.EpsgOut)
	if err != nil {
		return nil, err
	}
	return &TransformResult{
		Grid:        grid,
		Uncertainty: unc,
		EpsgOut:     outEpsg,
		NoData:      vt.opts.NoData,
	}, nil
}

type transformStep struct {
	from int
	to   int
	via  string
}

func (vt *VerticalTransform) verticalTransform(epsgIn, epsgOut int) ([]float64, []float64, int, error) {
	if vt.opts.Region == nil {
		return nil, nil, epsgOut, fmt.Errorf("no region specified for vertical transform")
	}

	n := vt.xCount * vt.yCount
	transArray := make([]float64, n)
	uncArray := make([]float64, n)

	if epsgIn == epsgOut {
		return transArray, uncArray, epsgOut, nil
	}

	frameIn := GetFrameByEPSG(epsgIn)
	frameOut := GetFrameByEPSG(epsgOut)
	if frameIn == nil || frameOut == nil {
		return nil, nil, epsgOut, fmt.Errorf("unknown vertical frame: EPSG %d or %d is not a registered vertical datum", epsgIn, epsgOut)
	}

	for i := range uncArray {
		uncArray[i] = frameIn.Uncertainty
	}

	steps, err := planSteps(epsgIn, epsgOut, frameIn, frameOut)
	if err != nil {
		return nil, nil, epsgOut, err
	}
	if len(steps) == 0 {
		return nil, nil, epsgOut, fmt.Errorf("no transform path from EPSG %d to %d", epsgIn, epsgOut)
	}

	currentEpsg := epsgIn
	for i, step := range steps {
		if step.from != currentEpsg {
			return nil, nil, currentEpsg, fmt.Errorf("broken transform chain: step %d→%d but current frame is %d", step.from, step.to, currentEpsg)
		}
		if err := dem.CheckCtx(vt.opts.Ctx); err != nil {
			return nil, nil, currentEpsg, fmt.Errorf("%s: %w", step.via, err)
		}
		dem.ReportProgress(vt.opts.Progress, step.via, i, len(steps))
		stepGrid, geoidUnc, err := vt.executeStep(step)
		if err != nil {
			return nil, nil, currentEpsg, err
		}
		dem.ReportProgress(vt.opts.Progress, step.via, i+1, len(steps))
		for i := range transArray {
			transArray[i] += stepGrid[i]
		}
		currentEpsg = step.to

		stepUnc := math.Hypot(FrameUncertainty(step.to), geoidUnc)
		if stepUnc > 0 {
			for i := range uncArray {
				uncArray[i] = math.Sqrt(uncArray[i]*uncArray[i] + stepUnc*stepUnc)
			}
		}
	}

	return transArray, uncArray, currentEpsg, nil
}

func planSteps(epsgIn, epsgOut int, frameIn, frameOut *Frame) ([]transformStep, error) {
	if frameIn.Type == FrameTidal && epsgIn != mslEPSG {
		return nil, fmt.Errorf("no offset grid available for tidal datum EPSG %d (%s): only MSL (%d) is supported", epsgIn, frameIn.Name, mslEPSG)
	}
	if frameOut.Type == FrameTidal && epsgOut != mslEPSG {
		return nil, fmt.Errorf("no offset grid available for tidal datum EPSG %d (%s): only MSL (%d) is supported", epsgOut, frameOut.Name, mslEPSG)
	}

	switch {
	case frameIn.Type == FrameTidal && frameOut.Type == FrameCDN:
		return []transformStep{{from: epsgIn, to: epsgOut, via: "msl2cdn"}}, nil

	case frameIn.Type == FrameCDN && frameOut.Type == FrameTidal:
		return []transformStep{{from: epsgIn, to: epsgOut, via: "cdn2msl"}}, nil

	case frameIn.Type == FrameCDN && frameOut.Type == FrameCDN:
		return []transformStep{{from: epsgIn, to: epsgOut, via: "cdn2cdn"}}, nil

	case frameIn.Type == FrameHTDP && frameOut.Type == FrameHTDP:
		return []transformStep{{from: epsgIn, to: epsgOut, via: "htdp2htdp"}}, nil

	case frameIn.Type == FrameHTDP && frameOut.Type == FrameCDN:
		return []transformStep{
			{from: epsgIn, to: ellipsoidEPSG, via: "htdp2ellipsoid"},
			{from: ellipsoidEPSG, to: epsgOut, via: "ellipsoid2cdn"},
		}, nil

	case frameIn.Type == FrameCDN && frameOut.Type == FrameHTDP:
		return []transformStep{
			{from: epsgIn, to: ellipsoidEPSG, via: "cdn2ellipsoid"},
			{from: ellipsoidEPSG, to: epsgOut, via: "ellipsoid2htdp"},
		}, nil

	case frameIn.Type == FrameHTDP && frameOut.Type == FrameTidal:
		return []transformStep{
			{from: epsgIn, to: ellipsoidEPSG, via: "htdp2ellipsoid"},
			{from: ellipsoidEPSG, to: epsgOut, via: "ellipsoid2msl"},
		}, nil

	case frameIn.Type == FrameTidal && frameOut.Type == FrameHTDP:
		return []transformStep{
			{from: epsgIn, to: ellipsoidEPSG, via: "msl2ellipsoid"},
			{from: ellipsoidEPSG, to: epsgOut, via: "ellipsoid2htdp"},
		}, nil
	}

	return nil, fmt.Errorf("no transform path from EPSG %d to %d", epsgIn, epsgOut)
}

func (vt *VerticalTransform) executeStep(step transformStep) ([]float64, float64, error) {
	switch step.via {
	case "msl2ellipsoid":
		g, err := vt.computeGeoidGrid(mslModel, step.via)
		if err != nil {
			return nil, 0, err
		}
		return invertGrid(g), geoidUncertainty(mslModel), nil
	case "ellipsoid2msl":
		g, err := vt.computeGeoidGrid(mslModel, step.via)
		if err != nil {
			return nil, 0, err
		}
		return g, geoidUncertainty(mslModel), nil
	case "msl2cdn":
		toModel, err := vt.endpointModel(step.to, vt.opts.GeoidOut)
		if err != nil {
			return nil, 0, err
		}
		g, err := vt.differenceGrid(mslModel, toModel, step.via)
		if err != nil {
			return nil, 0, err
		}
		return g, geoidUncertainty(mslModel, toModel), nil
	case "cdn2msl":
		fromModel, err := vt.endpointModel(step.from, vt.opts.GeoidIn)
		if err != nil {
			return nil, 0, err
		}
		g, err := vt.differenceGrid(fromModel, mslModel, step.via)
		if err != nil {
			return nil, 0, err
		}
		return g, geoidUncertainty(fromModel, mslModel), nil
	case "cdn2cdn":
		return vt.cdnTransform(step.from, step.to, step.via)
	case "cdn2ellipsoid":
		fromModel, err := vt.endpointModel(step.from, vt.opts.GeoidIn)
		if err != nil {
			return nil, 0, err
		}
		g, err := vt.computeGeoidGrid(fromModel, step.via)
		if err != nil {
			return nil, 0, err
		}
		return invertGrid(g), geoidUncertainty(fromModel), nil
	case "ellipsoid2cdn":
		toModel, err := vt.endpointModel(step.to, vt.opts.GeoidOut)
		if err != nil {
			return nil, 0, err
		}
		g, err := vt.computeGeoidGrid(toModel, step.via)
		if err != nil {
			return nil, 0, err
		}
		return g, geoidUncertainty(toModel), nil
	case "htdp2htdp":
		return vt.htdpStepGrid(step.from, step.to)
	case "htdp2ellipsoid":
		return vt.htdpStepGrid(step.from, ellipsoidEPSG)
	case "ellipsoid2htdp":
		return vt.htdpStepGrid(ellipsoidEPSG, step.to)
	default:
		return nil, 0, fmt.Errorf("unknown transform step %q", step.via)
	}
}

func invertGrid(grid []float64) []float64 {
	out := make([]float64, len(grid))
	for i, v := range grid {
		out[i] = -v
	}
	return out
}

func (vt *VerticalTransform) differenceGrid(fromModel, toModel geoid.VerticalDatum, stage string) ([]float64, error) {
	fromGrid, err := vt.computeGeoidGrid(fromModel, stage)
	if err != nil {
		return nil, err
	}
	toGrid, err := vt.computeGeoidGrid(toModel, stage)
	if err != nil {
		return nil, err
	}
	grid := make([]float64, vt.xCount*vt.yCount)
	for i := range grid {
		grid[i] = toGrid[i] - fromGrid[i]
	}
	return grid, nil
}

func (vt *VerticalTransform) endpointModel(epsg int, override string) (geoid.VerticalDatum, error) {
	if override != "" {
		if m := geoid.VerticalDatumFromString(override); m != geoid.HAE && m != geoid.UNKNOWN {
			return m, nil
		}
	}
	m := EPSGToVerticalDatum(epsg)
	if m == geoid.HAE || m == geoid.UNKNOWN {
		return m, fmt.Errorf("cannot resolve geoid model for vertical datum EPSG %d", epsg)
	}
	return m, nil
}

func (vt *VerticalTransform) cdnTransform(fromEPSG, toEPSG int, stage string) ([]float64, float64, error) {
	fromModel, err := vt.endpointModel(fromEPSG, vt.opts.GeoidIn)
	if err != nil {
		return nil, 0, err
	}
	toModel, err := vt.endpointModel(toEPSG, vt.opts.GeoidOut)
	if err != nil {
		return nil, 0, err
	}
	grid, err := vt.differenceGrid(fromModel, toModel, stage)
	if err != nil {
		return nil, 0, err
	}
	return grid, geoidUncertainty(fromModel, toModel), nil
}

func (vt *VerticalTransform) htdpStepGrid(fromEPSG, toEPSG int) ([]float64, float64, error) {
	grid, err := vt.htdpGrid(fromEPSG, toEPSG)
	if err != nil {
		return nil, 0, err
	}
	return invertGrid(grid), 0, nil
}

func (vt *VerticalTransform) htdpGrid(fromEPSG, toEPSG int) ([]float64, error) {
	frameIn := GetFrameByEPSG(fromEPSG)
	frameOut := GetFrameByEPSG(toEPSG)
	if frameIn == nil || frameOut == nil {
		return nil, fmt.Errorf("unknown htdp frame: EPSG %d or %d", fromEPSG, toEPSG)
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
	g, _ := computeGeoidGridProgress(region, model, model.ToString(), nil, nil)
	return g
}

func computeGeoidGridProgress(region *dem.Region, model geoid.VerticalDatum, stage string, prog dem.ProgressFunc, ctx context.Context) ([]float64, error) {
	g := getGeoid(model, true)
	n := region.XSize * region.YSize
	grid := make([]float64, n)
	noData := dem.DefaultNoData

	gt := region.GeoTransform()
	pts := make([]vec2.T, n)
	i := 0
	for y := 0; y < region.YSize; y++ {
		for x := 0; x < region.XSize; x++ {
			pts[i] = vec2.T{
				gt[0] + (float64(x)+0.5)*gt[1] + (float64(y)+0.5)*gt[2],
				gt[3] + (float64(x)+0.5)*gt[4] + (float64(y)+0.5)*gt[5],
			}
			i++
		}
	}

	var srs4326 geo.Proj = geo.NewProj("EPSG:4326")
	needTransform := region.SRS() != nil && !region.SRS().Eq(srs4326)
	if needTransform {
		pts = region.SRS().TransformTo(srs4326, pts)
	}

	dem.ReportProgress(prog, stage, 0, region.YSize)
	for y := 0; y < region.YSize; y++ {
		if err := dem.CheckCtx(ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stage, err)
		}
		for x := 0; x < region.XSize; x++ {
			p := pts[y*region.XSize+x]
			und := g.GetHeight(p[1], p[0])
			if math.IsNaN(und) || math.IsInf(und, 0) {
				grid[y*region.XSize+x] = noData
			} else {
				grid[y*region.XSize+x] = und
			}
		}
		dem.ReportProgress(prog, stage, y+1, region.YSize)
	}
	dem.ReportProgress(prog, stage, region.YSize, region.YSize)
	return grid, nil
}

func (vt *VerticalTransform) computeGeoidGrid(model geoid.VerticalDatum, stage string) ([]float64, error) {
	if g, ok := vt.geoidGrids[model]; ok {
		return g, nil
	}
	g, err := computeGeoidGridProgress(vt.opts.Region, model, stage, vt.opts.Progress, vt.opts.Ctx)
	if err != nil {
		return nil, err
	}
	vt.geoidGrids[model] = g
	return g, nil
}
