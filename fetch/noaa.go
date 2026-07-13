package fetch

import (
	"fmt"
	"path/filepath"
)

type noaaMultibeamFetcher struct{ baseFetcher }

func init() {
	Register(&noaaMultibeamFetcher{baseFetcher: newBaseFetcher(SourceNOAAMultibeam)})
}

func (f *noaaMultibeamFetcher) Name() DataSource { return SourceNOAAMultibeam }

func (f *noaaMultibeamFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	outputPath := filepath.Join(opts.OutputDir, "noaa_multibeam.datalist")

	url := fmt.Sprintf(
		"https://www.ngdc.noaa.gov/thredds/ncss/grid/MultibeamSurvey/aggregate?"+
			"west=%.6f&south=%.6f&east=%.6f&north=%.6f&"+
			"disableProjSubset=on&horizStride=1&"+
			"var=topo&disableLLSubset=on&addLatLon=true",
		opts.MinX, opts.MinY, opts.MaxX, opts.MaxY,
	)

	if err := f.download(url, outputPath, opts); err != nil {
		return nil, fmt.Errorf("noaa multibeam: %v", err)
	}

	return []string{outputPath}, nil
}

type usgsTnmFetcher struct{ baseFetcher }

func init() {
	Register(&usgsTnmFetcher{baseFetcher: newBaseFetcher(SourceUSGSTNM)})
}

func (f *usgsTnmFetcher) Name() DataSource { return SourceUSGSTNM }

func (f *usgsTnmFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	outputPath := filepath.Join(opts.OutputDir, "usgs_tnm_dem.tif")

	url := fmt.Sprintf(
		"https://tnmaccess.nationalmap.gov/api/v1/products?"+
			"bbox=%.6f,%.6f,%.6f,%.6f&"+
			"prodFormats=GeoTIFF&"+
			"datasets=National+Elevation+Dataset+(NED)+1+arc-second&"+
			"outputFormat=JSON",
		opts.MinY, opts.MinX, opts.MaxY, opts.MaxX,
	)

	if err := f.download(url, outputPath, opts); err != nil {
		return nil, fmt.Errorf("usgs tnm: %v", err)
	}

	return []string{outputPath}, nil
}

type emodnetFetcher struct{ baseFetcher }

func init() {
	Register(&emodnetFetcher{baseFetcher: newBaseFetcher(SourceEMODNet)})
}

func (f *emodnetFetcher) Name() DataSource { return SourceEMODNet }

func (f *emodnetFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	outputPath := filepath.Join(opts.OutputDir, "emodnet_dem.tif")

	url := fmt.Sprintf(
		"https://ws.emodnet-bathymetry.eu/wcs?"+
			"service=WCS&version=2.0.1&request=GetCoverage&"+
			"CoverageId=emodnet_bathymetry&"+
			"subset=Lat(%.6f,%.6f)&subset=Long(%.6f,%.6f)&"+
			"format=image/tiff",
		opts.MinY, opts.MaxY, opts.MinX, opts.MaxX,
	)

	if err := f.download(url, outputPath, opts); err != nil {
		return nil, fmt.Errorf("emodnet: %v", err)
	}

	return []string{outputPath}, nil
}

type arcticDemFetcher struct{ baseFetcher }

func init() {
	Register(&arcticDemFetcher{baseFetcher: newBaseFetcher(SourceArcticDEM)})
}

func (f *arcticDemFetcher) Name() DataSource { return SourceArcticDEM }

func (f *arcticDemFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	outputPath := filepath.Join(opts.OutputDir, "arctic_dem.tif")

	url := fmt.Sprintf(
		"https://api.pgc.umn.edu/api/1/datasets/arcticdem/mosaic/tiles?"+
			"bbox=%.6f,%.6f,%.6f,%.6f&output_format=GeoTIFF",
		opts.MinX, opts.MinY, opts.MaxX, opts.MaxY,
	)

	if err := f.download(url, outputPath, opts); err != nil {
		return nil, fmt.Errorf("arctic dem: %v", err)
	}

	return []string{outputPath}, nil
}
