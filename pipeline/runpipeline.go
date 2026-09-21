package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/datalist"
	"github.com/flywave/go-dem/datum"
	"github.com/flywave/go-dem/grits"
	"github.com/flywave/go-dem/perspecto"
	"github.com/flywave/go-dem/uncertainty"
	"github.com/flywave/go-dem/waffle"
	"github.com/flywave/go-geo"
	"github.com/flywave/go-geoid"
)

const (
	stageLoad        = "load"
	stageInterpolate = "interpolate"
	stageFilter      = "filter"
	stageDatum       = "datum"
	stageUncertainty = "uncertainty"
	stageVisualize   = "visualize"
	stageWrite       = "write"
)

var stageOrder = []string{stageLoad, stageInterpolate, stageFilter, stageDatum, stageUncertainty, stageVisualize, stageWrite}

type PipelineConfig struct {
	Config
	Method      dem.InterpMethod
	Filters     []string
	FilterOpts  *grits.Options
	Uncertainty *uncertainty.Options
	Hillshade   bool
	Histogram   bool
	OutputDir   string
}

type PipelineResult struct {
	DEMPath          string
	HillshadePath    string
	HistogramPath    string
	UncertaintyPaths []string
	PointCount       int
	Region           *dem.Region
}

func RunPipeline(inputs []string, region *dem.Region, cfg *PipelineConfig, progress dem.ProgressFunc, ctx context.Context) (*PipelineResult, error) {
	if len(inputs) == 0 {
		return nil, fmt.Errorf("inputs are required")
	}
	if region == nil {
		return nil, fmt.Errorf("region is required")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if cfg.OutputDir == "" {
		return nil, fmt.Errorf("output dir is required")
	}
	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}
	if nilSRS(region.SRS()) {
		fixed := *region
		fixed.Extent = &geo.MapExtent{BBox: region.BBox(), Srs: geo.NewProj("EPSG:4326")}
		region = &fixed
	}

	if cfg.TargetEpsg > 0 && cfg.TargetEpsg != cfg.SourceEpsg && datum.GetFrameByEPSG(cfg.TargetEpsg) == nil {
		return nil, fmt.Errorf("horizontal reprojection to EPSG %d is not supported: TargetEpsg must be a registered vertical datum EPSG, horizontal reprojection requires gdal warp resampling", cfg.TargetEpsg)
	}
	for _, name := range cfg.Filters {
		if _, err := grits.New(grits.FilterType(name)); err != nil {
			return nil, fmt.Errorf("%s %s: %v", stageFilter, name, err)
		}
	}

	noData := cfg.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	windows := stageWindows(&cfg.Config, cfg.Filters, cfg.Uncertainty != nil, cfg.Hillshade || cfg.Histogram)
	tracker := &progressTracker{progress: progress}

	if err := dem.CheckCtx(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", stageLoad, err)
	}

	loadWin := windows[stageLoad]
	points, err := loadPoints(inputs, tracker.stageProgress(loadWin.base, loadWin.span), ctx)
	if err != nil {
		return nil, err
	}
	tracker.seal(loadWin)

	method := cfg.Method
	if method == "" {
		method = dem.MethodIDW
	}
	interpWin := windows[stageInterpolate]
	if err := dem.CheckCtx(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", stageInterpolate, err)
	}
	w, err := waffle.New(method)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", stageInterpolate, err)
	}
	result, err := w.Run(points, &waffle.Options{
		Region:   region,
		NoData:   noData,
		Progress: tracker.stageProgress(interpWin.base, interpWin.span),
		Ctx:      ctx,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %v", stageInterpolate, err)
	}
	data := result.DEM
	tracker.seal(interpWin)

	if len(cfg.Filters) > 0 {
		filterWin := windows[stageFilter]
		if err := dem.CheckCtx(ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stageFilter, err)
		}
		fOpts := &grits.Options{}
		if cfg.FilterOpts != nil {
			*fOpts = *cfg.FilterOpts
		}
		if fOpts.NoData == nil {
			fOpts.NoData = &noData
		}
		fOpts.Ctx = ctx
		per := filterWin.span / len(cfg.Filters)
		rem := filterWin.span % len(cfg.Filters)
		subBase := filterWin.base
		for j, name := range cfg.Filters {
			f, err := grits.New(grits.FilterType(name))
			if err != nil {
				return nil, fmt.Errorf("%s %s: %v", stageFilter, name, err)
			}
			subSpan := per
			if j == len(cfg.Filters)-1 {
				subSpan += rem
			}
			fOpts.Progress = tracker.stageProgress(subBase, subSpan)
			data, err = f.Run(data, region, fOpts)
			if err != nil {
				return nil, fmt.Errorf("%s %s: %v", stageFilter, name, err)
			}
			subBase += subSpan
		}
		tracker.seal(filterWin)
	}

	datumWin := windows[stageDatum]
	if err := dem.CheckCtx(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", stageDatum, err)
	}
	data, currentVerticalEpsg, transformed, err := applyVerticalDatum(data, region, cfg.Config, noData,
		tracker.stageProgress(datumWin.base, datumWin.span), ctx)
	if err != nil {
		return nil, err
	}
	tracker.seal(datumWin)

	var unc *uncertainty.Result
	if cfg.Uncertainty != nil {
		uncWin := windows[stageUncertainty]
		if err := dem.CheckCtx(ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stageUncertainty, err)
		}
		uOpts := *cfg.Uncertainty
		if uOpts.Method == "" {
			uOpts.Method = uncertainty.MethodProximity
		}
		if uOpts.NoData == 0 {
			uOpts.NoData = noData
		}
		uOpts.Progress = tracker.stageProgress(uncWin.base, uncWin.span)
		uOpts.Ctx = ctx
		unc, err = uncertainty.Estimate(data, region, &uOpts)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", stageUncertainty, err)
		}
		tracker.seal(uncWin)
	}

	hillshadePath := ""
	histogramPath := ""
	if cfg.Hillshade || cfg.Histogram {
		visWin := windows[stageVisualize]
		if err := dem.CheckCtx(ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stageVisualize, err)
		}
		tasks := 0
		if cfg.Hillshade {
			tasks++
		}
		if cfg.Histogram {
			tasks++
		}
		visProgress := tracker.stageProgress(visWin.base, visWin.span)
		done := 0
		if cfg.Hillshade {
			hillshadePath = filepath.Join(cfg.OutputDir, "hillshade.tif")
			if err := writeHillshade(data, region, hillshadePath, noData); err != nil {
				return nil, fmt.Errorf("%s: %v", stageVisualize, err)
			}
			done++
			dem.ReportProgress(visProgress, "hillshade", done, tasks)
		}
		if cfg.Histogram {
			histogramPath = filepath.Join(cfg.OutputDir, "histogram.png")
			if err := perspecto.WriteHistogramPNG(data, histogramPath, &perspecto.HistogramOptions{NoData: noData}); err != nil {
				return nil, fmt.Errorf("%s: %v", stageVisualize, err)
			}
			done++
			dem.ReportProgress(visProgress, "histogram", done, tasks)
		}
		tracker.seal(visWin)
	}

	writeWin := windows[stageWrite]
	if err := dem.CheckCtx(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", stageWrite, err)
	}
	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		return nil, fmt.Errorf("%s: %v", stageWrite, err)
	}
	demPath := filepath.Join(cfg.OutputDir, "dem.tif")
	var uncPaths []string
	if unc != nil {
		if err := uncertainty.WriteUncertainty(unc, region, demPath, noData); err != nil {
			return nil, fmt.Errorf("%s: %v", stageWrite, err)
		}
		uncBase := strings.TrimSuffix(demPath, ".tif")
		uncPaths = []string{
			uncBase + "_tvu.tif",
			uncBase + "_interp_u.tif",
			uncBase + "_src_u.tif",
			uncBase + "_prox.tif",
		}
	}
	outDatum := cfg.VerticalDatum
	if transformed {
		if m := datum.EPSGToVerticalDatum(currentVerticalEpsg); m != geoid.HAE && m != geoid.UNKNOWN {
			outDatum = m
		}
	}
	if err := dem.CreateDEMWithConfig(data, region, demPath, dem.OutputConfig{NoData: noData, VerticalDatum: outDatum}); err != nil {
		return nil, fmt.Errorf("%s: %v", stageWrite, err)
	}
	tracker.seal(writeWin)

	return &PipelineResult{
		DEMPath:          demPath,
		HillshadePath:    hillshadePath,
		HistogramPath:    histogramPath,
		UncertaintyPaths: uncPaths,
		PointCount:       len(points),
		Region:           region,
	}, nil
}

func loadPoints(inputs []string, progress dem.ProgressFunc, ctx context.Context) ([]waffle.Point, error) {
	dl, err := datalist.BuildDataList(inputs)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", stageLoad, err)
	}
	if len(dl.Entries) == 0 {
		return nil, fmt.Errorf("%s: no readable sources found in inputs", stageLoad)
	}
	var points []waffle.Point
	for i := range dl.Entries {
		if err := dem.CheckCtx(ctx); err != nil {
			return nil, fmt.Errorf("%s: %w", stageLoad, err)
		}
		pts, err := datalist.PointsFromDataList(&datalist.DataList{Entries: dl.Entries[i : i+1]})
		if err != nil {
			return nil, fmt.Errorf("%s %s: %v", stageLoad, dl.Entries[i].Path, err)
		}
		points = append(points, pts...)
		dem.ReportProgress(progress, stageLoad, i+1, len(dl.Entries))
	}
	if len(points) == 0 {
		return nil, fmt.Errorf("%s: no valid points found in inputs", stageLoad)
	}
	return points, nil
}

func writeHillshade(data []float64, region *dem.Region, path string, noData float64) error {
	shade := perspecto.Hillshade(data, region, &perspecto.Options{NoData: noData})
	pixels := make([]uint8, len(shade)*3)
	for i, v := range shade {
		b := uint8(0)
		if v > 0 {
			b = uint8(v)
		}
		pixels[i*3] = b
		pixels[i*3+1] = b
		pixels[i*3+2] = b
	}
	return perspecto.WriteRGB(pixels, region, path)
}

type stageWindow struct {
	name string
	base int
	span int
}

func stageWindows(cfg *Config, filters []string, withUncertainty, withVisualize bool) map[string]stageWindow {
	active := map[string]bool{
		stageLoad:        true,
		stageInterpolate: true,
		stageFilter:      len(filters) > 0,
		stageDatum:       hasVerticalDatum(cfg),
		stageUncertainty: withUncertainty,
		stageVisualize:   withVisualize,
		stageWrite:       true,
	}
	weights := map[string]int{
		stageLoad:        10,
		stageInterpolate: 40,
		stageFilter:      15,
		stageDatum:       10,
		stageUncertainty: 15,
		stageVisualize:   5,
		stageWrite:       5,
	}
	for _, name := range stageOrder {
		if !active[name] {
			weights[stageInterpolate] += weights[name]
			weights[name] = 0
		}
	}
	windows := make(map[string]stageWindow, len(stageOrder))
	base := 0
	for _, name := range stageOrder {
		windows[name] = stageWindow{name: name, base: base, span: weights[name]}
		base += weights[name]
	}
	return windows
}

func hasVerticalDatum(cfg *Config) bool {
	if cfg.VerticalDatum != geoid.HAE && cfg.VerticalDatum != geoid.UNKNOWN {
		return true
	}
	return cfg.TargetEpsg > 0
}

type progressTracker struct {
	progress dem.ProgressFunc
	last     int
}

func (t *progressTracker) stageProgress(base, span int) dem.ProgressFunc {
	if t.progress == nil {
		return nil
	}
	return func(stage string, done, total int) {
		if total <= 0 || done < 0 {
			return
		}
		frac := float64(done) / float64(total)
		if frac > 1 {
			frac = 1
		}
		global := base + int(frac*float64(span))
		if global > base+span {
			global = base + span
		}
		if global < t.last {
			global = t.last
		}
		t.last = global
		t.progress(stage, global, 100)
	}
}

func (t *progressTracker) seal(w stageWindow) {
	if w.span <= 0 {
		return
	}
	end := w.base + w.span
	if end > t.last {
		t.last = end
	}
	dem.ReportProgress(t.progress, w.name, end, 100)
}

func nilSRS(p geo.Proj) bool {
	if p == nil {
		return true
	}
	sp, ok := p.(*geo.SRSProj4)
	return ok && sp == nil
}
