package fetch

import (
	"fmt"
	"path/filepath"
)

type srtmFetcher struct{ baseFetcher }

func init() {
	Register(&srtmFetcher{baseFetcher: newBaseFetcher(SourceSRTM)})
}

func (f *srtmFetcher) Name() DataSource { return SourceSRTM }

func (f *srtmFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	baseURL := "https://opentopography.org/otr/getdem"
	url := fmt.Sprintf("%s?demtype=SRTM&west=%.6f&south=%.6f&east=%.6f&north=%.6f&output=GTiff",
		baseURL, opts.MinX, opts.MinY, opts.MaxX, opts.MaxY)

	outputPath := filepath.Join(opts.OutputDir, "srtm_dem.tif")
	if err := f.download(url, outputPath, opts); err != nil {
		return nil, fmt.Errorf("srtm: %v", err)
	}

	return []string{outputPath}, nil
}

type gebcoFetcher struct{ baseFetcher }

func init() {
	Register(&gebcoFetcher{baseFetcher: newBaseFetcher(SourceGEBCO)})
}

func (f *gebcoFetcher) Name() DataSource { return SourceGEBCO }

func (f *gebcoFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	baseURL := "https://www.gebco.net/data_and_products/gebco_web_service/web_service"
	url := fmt.Sprintf("%s?action=grid&bbox=%s&format=geotiff",
		baseURL, f.bboxParam(opts.MinX, opts.MinY, opts.MaxX, opts.MaxY))

	outputPath := filepath.Join(opts.OutputDir, "gebco_dem.tif")
	if err := f.download(url, outputPath, opts); err != nil {
		return nil, fmt.Errorf("gebco: %v", err)
	}

	return []string{outputPath}, nil
}

type copernicusFetcher struct{ baseFetcher }

func init() {
	Register(&copernicusFetcher{baseFetcher: newBaseFetcher(SourceCopernicus)})
}

func (f *copernicusFetcher) Name() DataSource { return SourceCopernicus }

func (f *copernicusFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	baseURL := "https://prism-dem-open.copernicus.eu/pd-desk-open-service/open-service"
	url := fmt.Sprintf("%s/dem/v1?bbox=%s&product=COPDEM&format=GeoTIFF",
		baseURL, f.bboxParam(opts.MinX, opts.MinY, opts.MaxX, opts.MaxY))

	outputPath := filepath.Join(opts.OutputDir, "copernicus_dem.tif")
	if err := f.download(url, outputPath, opts); err != nil {
		return nil, fmt.Errorf("copernicus: %v", err)
	}

	return []string{outputPath}, nil
}
