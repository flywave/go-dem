package pipeline

import (
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
	w, err := waffle.New(method)
	if err != nil {
		return fmt.Errorf("waffle: %v", err)
	}

	result, err := w.Run(points, opts)
	if err != nil {
		return fmt.Errorf("interpolation: %v", err)
	}

	if cfg.VerticalDatum != geoid.HAE && cfg.VerticalDatum != geoid.UNKNOWN {
		dstEpsg := verticalEPSG(cfg.VerticalDatum)
		if dstEpsg > 0 && dstEpsg != cfg.SourceEpsg {
			transformed, err := datum.TransformDEM(result.DEM, region, cfg.SourceEpsg, dstEpsg)
			if err != nil {
				return fmt.Errorf("datum: %v", err)
			}
			result.DEM = transformed
		}
	}

	if cfg.TargetEpsg > 0 && cfg.TargetEpsg != cfg.SourceEpsg {
		transformed, err := datum.TransformDEM(result.DEM, region, cfg.SourceEpsg, cfg.TargetEpsg)
		if err != nil {
			return fmt.Errorf("datum target: %v", err)
		}
		result.DEM = transformed
	}

	outCfg := dem.OutputConfig{
		NoData:        cfg.NoData,
		VerticalDatum: cfg.VerticalDatum,
	}
	return dem.CreateDEMWithConfig(result.DEM, region, outPath, outCfg)
}

func RunGMTBlockmean(points []waffle.Point, region *dem.Region, method dem.InterpMethod, opts *waffle.Options, outPath string, cfg Config) error {
	w, err := waffle.New(method)
	if err != nil {
		return fmt.Errorf("waffle: %v", err)
	}

	result, err := w.Run(points, opts)
	if err != nil {
		return fmt.Errorf("interpolation: %v", err)
	}

	if cfg.TargetEpsg > 0 && cfg.TargetEpsg != cfg.SourceEpsg {
		transformed, err := datum.TransformDEM(result.DEM, region, cfg.SourceEpsg, cfg.TargetEpsg)
		if err != nil {
			return fmt.Errorf("datum target: %v", err)
		}
		result.DEM = transformed
	}

	outCfg := dem.OutputConfig{
		NoData:        cfg.NoData,
		VerticalDatum: cfg.VerticalDatum,
	}
	return dem.CreateDEMWithConfig(result.DEM, region, outPath, outCfg)
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
