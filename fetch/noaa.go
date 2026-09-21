package fetch

import (
	"encoding/json"
	"fmt"
	"net/url"
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

var tnmProductsURL = "https://tnmaccess.nationalmap.gov/api/v1/products"

const tnmNED1ArcSecond = "National Elevation Dataset (NED) 1 arc-second"

type tnmProduct struct {
	Title       string `json:"title"`
	Format      string `json:"format"`
	DownloadURL string `json:"downloadURL"`
}

type tnmProductsResponse struct {
	Items []tnmProduct `json:"items"`
	Total int          `json:"total"`
}

type usgsTnmFetcher struct{ baseFetcher }

func init() {
	Register(&usgsTnmFetcher{baseFetcher: newBaseFetcher(SourceUSGSTNM)})
}

func (f *usgsTnmFetcher) Name() DataSource { return SourceUSGSTNM }

func (f *usgsTnmFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	q := url.Values{}
	// TNM products API 的 bbox 顺序为 south,west,north,east。
	q.Set("bbox", fmt.Sprintf("%.6f,%.6f,%.6f,%.6f", opts.MinY, opts.MinX, opts.MaxY, opts.MaxX))
	q.Set("prodFormats", "GeoTIFF")
	q.Set("datasets", tnmNED1ArcSecond)
	q.Set("outputFormat", "JSON")

	resp, err := f.get(tnmProductsURL+"?"+q.Encode(), opts)
	if err != nil {
		return nil, fmt.Errorf("usgs tnm: %v", err)
	}
	defer resp.Body.Close()

	var products tnmProductsResponse
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("usgs tnm: decode products: %v", err)
	}

	outputs := make([]string, 0, len(products.Items))
	for _, product := range products.Items {
		if product.DownloadURL == "" {
			continue
		}
		outputPath := filepath.Join(opts.OutputDir, fmt.Sprintf("usgs_tnm_%d.tif", len(outputs)))
		if err := f.download(product.DownloadURL, outputPath, opts); err != nil {
			return nil, fmt.Errorf("usgs tnm: %v", err)
		}
		outputs = append(outputs, outputPath)
	}

	if len(outputs) == 0 {
		return nil, fmt.Errorf("usgs tnm: no GeoTIFF products found for bbox (%.6f,%.6f)-(%.6f,%.6f)",
			opts.MinX, opts.MinY, opts.MaxX, opts.MaxY)
	}

	return outputs, nil
}

type emodnetFetcher struct{ baseFetcher }

func init() {
	Register(&emodnetFetcher{baseFetcher: newBaseFetcher(SourceEMODNet)})
}

func (f *emodnetFetcher) Name() DataSource { return SourceEMODNet }

func (f *emodnetFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	// 端点未经官方文档确认：CoverageId=emodnet_bathymetry 未能在 EMODnet WCS GetCapabilities
	// 中考证（实际发布的 CoverageId 形如 emodnet:mean），使用前应先执行 GetCapabilities 确认。
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

var pgcArcticDEMURL = "https://api.pgc.umn.edu/api/1/datasets/arcticdem/mosaic/tiles"

type arcticDemFetcher struct{ baseFetcher }

func init() {
	Register(&arcticDemFetcher{baseFetcher: newBaseFetcher(SourceArcticDEM)})
}

func (f *arcticDemFetcher) Name() DataSource { return SourceArcticDEM }

func (f *arcticDemFetcher) Fetch(opts *FetchOptions) ([]string, error) {
	outputPath := filepath.Join(opts.OutputDir, "arctic_dem.tif")

	url := fmt.Sprintf("%s?bbox=%.6f,%.6f,%.6f,%.6f&output_format=GeoTIFF",
		pgcArcticDEMURL, opts.MinX, opts.MinY, opts.MaxX, opts.MaxY)

	if err := f.download(url, outputPath, opts); err != nil {
		return nil, fmt.Errorf("arctic dem: %v", err)
	}

	return []string{outputPath}, nil
}
