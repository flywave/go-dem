package fetch

import (
	"fmt"
	"net/url"
	"path/filepath"
)

var opentopographyGlobalDEMURL = "https://portal.opentopography.org/API/globaldem"

type srtmFetcher struct{ baseFetcher }

func init() {
	Register(&srtmFetcher{baseFetcher: newBaseFetcher(SourceSRTM)})
}

func (f *srtmFetcher) Name() DataSource { return SourceSRTM }

func (f *srtmFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	q := url.Values{}
	q.Set("demtype", "SRTMGL1")
	q.Set("west", fmt.Sprintf("%.6f", opts.MinX))
	q.Set("south", fmt.Sprintf("%.6f", opts.MinY))
	q.Set("east", fmt.Sprintf("%.6f", opts.MaxX))
	q.Set("north", fmt.Sprintf("%.6f", opts.MaxY))
	q.Set("outputFormat", "GTiff")
	if opts.APIKey != "" {
		q.Set("API_Key", opts.APIKey)
	}

	outputPath := filepath.Join(opts.OutputDir, "srtm_dem.tif")
	if err := f.download(opentopographyGlobalDEMURL+"?"+q.Encode(), outputPath, opts); err != nil {
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
	// 端点未经官方文档确认：GEBCO 官方仅公开 WMS 服务，此直链地址与参数未能在公开文档中考证。
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
	// 端点未经官方文档确认：此 open-service 直链地址与参数未能在 Copernicus 公开文档中考证。
	baseURL := "https://prism-dem-open.copernicus.eu/pd-desk-open-service/open-service"
	url := fmt.Sprintf("%s/dem/v1?bbox=%s&product=COPDEM&format=GeoTIFF",
		baseURL, f.bboxParam(opts.MinX, opts.MinY, opts.MaxX, opts.MaxY))

	outputPath := filepath.Join(opts.OutputDir, "copernicus_dem.tif")
	if err := f.download(url, outputPath, opts); err != nil {
		return nil, fmt.Errorf("copernicus: %v", err)
	}

	return []string{outputPath}, nil
}
