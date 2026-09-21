package pipeline

import (
	"context"
	"fmt"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/datum"
	"github.com/flywave/go-dem/waffle"
	"github.com/flywave/go-geoid"
)

type Config struct {
	SourceEpsg    int
	TargetEpsg    int
	VerticalDatum geoid.VerticalDatum
	NoData        float64
}

func RunDEM(points []waffle.Point, region *dem.Region, method dem.InterpMethod, opts *waffle.Options, outPath string, cfg Config) error {
	if region == nil {
		return fmt.Errorf("region is required")
	}

	noData := cfg.NoData
	if noData == 0 && opts != nil {
		noData = opts.NoData
	}
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	if opts == nil {
		opts = &waffle.Options{}
	} else {
		cp := *opts
		opts = &cp
	}
	opts.Region = region
	opts.NoData = noData

	if cfg.TargetEpsg > 0 && cfg.TargetEpsg != cfg.SourceEpsg && datum.GetFrameByEPSG(cfg.TargetEpsg) == nil {
		return fmt.Errorf("horizontal reprojection to EPSG %d is not supported: TargetEpsg must be a registered vertical datum EPSG, horizontal reprojection requires gdal warp resampling", cfg.TargetEpsg)
	}

	w, err := waffle.New(method)
	if err != nil {
		return fmt.Errorf("waffle: %v", err)
	}

	result, err := w.Run(points, opts)
	if err != nil {
		return fmt.Errorf("interpolation: %v", err)
	}

	out, currentVerticalEpsg, transformed, err := applyVerticalDatum(result.DEM, region, cfg, noData, nil, nil)
	if err != nil {
		return err
	}
	result.DEM = out

	outDatum := cfg.VerticalDatum
	if transformed {
		if m := datum.EPSGToVerticalDatum(currentVerticalEpsg); m != geoid.HAE && m != geoid.UNKNOWN {
			outDatum = m
		}
	}

	outCfg := dem.OutputConfig{
		NoData:        noData,
		VerticalDatum: outDatum,
	}
	return dem.CreateDEMWithConfig(result.DEM, region, outPath, outCfg)
}

func applyVerticalDatum(data []float64, region *dem.Region, cfg Config, noData float64, progress dem.ProgressFunc, ctx context.Context) ([]float64, int, bool, error) {
	currentVerticalEpsg := cfg.SourceEpsg
	transformed := false

	if cfg.VerticalDatum != geoid.HAE && cfg.VerticalDatum != geoid.UNKNOWN {
		dstEpsg := verticalEPSG(cfg.VerticalDatum)
		if dstEpsg == 0 {
			return nil, 0, false, fmt.Errorf("unsupported vertical datum %s", cfg.VerticalDatum.ToString())
		}
		out, err := transformVertical(data, region, currentVerticalEpsg, dstEpsg, noData, progress, ctx)
		if err != nil {
			return nil, 0, false, verticalTransformError("vertical datum", currentVerticalEpsg, dstEpsg, err)
		}
		data = out
		currentVerticalEpsg = dstEpsg
		transformed = true
	}

	if cfg.TargetEpsg > 0 && cfg.TargetEpsg != currentVerticalEpsg && datum.GetFrameByEPSG(cfg.TargetEpsg) != nil {
		out, err := transformVertical(data, region, currentVerticalEpsg, cfg.TargetEpsg, noData, progress, ctx)
		if err != nil {
			return nil, 0, false, verticalTransformError("target vertical datum", currentVerticalEpsg, cfg.TargetEpsg, err)
		}
		data = out
		currentVerticalEpsg = cfg.TargetEpsg
		transformed = true
	}

	return data, currentVerticalEpsg, transformed, nil
}

func transformVertical(data []float64, region *dem.Region, fromEpsg, toEpsg int, noData float64, progress dem.ProgressFunc, ctx context.Context) ([]float64, error) {
	if fromEpsg == toEpsg {
		return data, nil
	}
	vt := datum.NewVerticalTransform(datum.TransformOptions{
		EpsgIn:   fromEpsg,
		EpsgOut:  toEpsg,
		Region:   region,
		NoData:   noData,
		Progress: progress,
		Ctx:      ctx,
	})
	res, err := vt.Run()
	if err != nil {
		return nil, err
	}
	vg := datum.VDatumGrid{
		Data:    res.Grid,
		Region:  region,
		SrcEpsg: fromEpsg,
		DstEpsg: toEpsg,
		NoData:  noData,
	}
	return vg.ApplyToDEM(data, false), nil
}

func verticalTransformError(phase string, fromEpsg, toEpsg int, err error) error {
	return fmt.Errorf("%s: vertical datum transform %d→%d: %v (both EPSG codes must be registered vertical or ellipsoidal frames such as 4979, 7912, 5714, 5773 or 5703; see datum.SupportedFrames())", phase, fromEpsg, toEpsg, err)
}

func verticalEPSG(vd geoid.VerticalDatum) int {
	switch vd {
	case geoid.EGM84:
		return 5798
	case geoid.EGM96:
		return 5773
	case geoid.EGM2008:
		return 3855
	default:
		return 0
	}
}
