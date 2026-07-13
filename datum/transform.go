package datum

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"

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

	steps := planSteps(epsgIn, epsgOut, frameIn, frameOut)

	currentEpsg := epsgIn
	for _, step := range steps {
		stepGrid := vt.executeStep(step, currentEpsg)
		for i := range transArray {
			transArray[i] += stepGrid[i]
		}
		currentEpsg = step.to
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
	n := vt.xCount * vt.yCount
	grid := make([]float64, n)

	frameIn := GetFrameByEPSG(fromEPSG)
	frameOut := GetFrameByEPSG(toEPSG)
	if frameIn == nil || frameOut == nil {
		return grid
	}

	if htdpExec, err := findHTDP(); err == nil {
		return vt.htdpExecGrid(htdpExec, frameIn, frameOut)
	}

	grid = vt.computeGeoidGrid(geoid.EGM96)

	return grid
}

func findHTDP() (string, error) {
	if p, err := exec.LookPath("htdp"); err == nil {
		return p, nil
	}
	candidates := []string{
		filepath.Join("libs", "darwin_arm", "htdp"),
		filepath.Join("libs", "darwin", "htdp"),
		filepath.Join("libs", "linux", "htdp"),
		filepath.Join("libs", "linux_arm", "htdp"),
		"/usr/local/bin/htdp",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			if abs, err := filepath.Abs(c); err == nil {
				return abs, nil
			}
			return c, nil
		}
	}
	return "", fmt.Errorf("htdp not found")
}

func (vt *VerticalTransform) htdpExecGrid(htdpPath string, frameIn, frameOut *Frame) []float64 {
	n := vt.xCount * vt.yCount
	grid := make([]float64, n)

	tmpDir, err := os.MkdirTemp("", "htdp_grid_*")
	if err != nil {
		return grid
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.xyz")
	ctrlPath := filepath.Join(tmpDir, "control.txt")
	outputPath := filepath.Join(tmpDir, "output.xyz")

	srcEpoch := frameIn.Epoch
	if srcEpoch == 0 {
		srcEpoch = 1997.0
	}
	dstEpoch := frameOut.Epoch
	if dstEpoch == 0 {
		dstEpoch = 2000.0
	}

	writeHTDPInput(vt.geoTrans, vt.xCount, vt.yCount, inputPath)
	writeHTDPControl(ctrlPath, outputPath, inputPath, frameIn.HTDPID, srcEpoch, frameOut.HTDPID, dstEpoch)

	cmd := exec.Command(htdpPath, ctrlPath)
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return grid
	} else {
		_ = out
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		return grid
	}

	line := 0
	count := 0
	for y := 0; y < vt.yCount && count < len(data); y++ {
		for x := 0; x < vt.xCount && line < len(data); x++ {
			for line < len(data) && (data[line] == '\n' || data[line] == '\r') {
				line++
			}
			if line >= len(data) {
				break
			}
			var lat, lon, h float64
			var fid string
			consumed := 0
			for line+consumed < len(data) && data[line+consumed] != '\n' {
				consumed++
			}
			lineStr := string(data[line : line+consumed])
			fmt.Sscanf(lineStr, "%s %f %f %f", &fid, &lat, &lon, &h)
			if y*vt.xCount+x < len(grid) {
				grid[y*vt.xCount+x] = h
			}
			line += consumed
		}
	}

	return grid
}

func writeHTDPInput(gt [6]float64, xCount, yCount int, path string) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	xMin := gt[0]
	yMax := gt[3]
	xInc := gt[1]
	yInc := -gt[5]

	for y := 0; y < yCount; y++ {
		lat := yMax + float64(y)*yInc
		for x := 0; x < xCount; x++ {
			lon := xMin + float64(x)*xInc
			fmt.Fprintf(f, "P %s %.8f %.8f 0.0\n", fmt.Sprintf("%d", y*xCount+x+1), lat, lon)
		}
	}
}

func writeHTDPControl(ctrlPath, outPath, inPath string, srcID int, srcEpoch float64, dstID int, dstEpoch float64) {
	content := fmt.Sprintf(
		"I\n%s\n%s\n1\n3\n%.4f %d %.4f %d\n5\n0\n0\n",
		inPath, outPath,
		srcEpoch, srcID, dstEpoch, dstID,
	)
	os.WriteFile(ctrlPath, []byte(content), 0644)
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
