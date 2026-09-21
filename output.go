package dem

import (
	"fmt"
	"math"
	"strings"

	"github.com/flywave/flywave-gdal"
	"github.com/flywave/go-geo"
	"github.com/flywave/go-geoid"
)

type OutputConfig struct {
	Format        string
	DataType      gdal.DataType
	Compress      string
	Tiled         bool
	BlockXSize    int
	BlockYSize    int
	NoData        float64
	CRS           string
	VerticalDatum geoid.VerticalDatum
}

var DefaultGTiffConfig = OutputConfig{
	Format:     "GTiff",
	DataType:   gdal.Float64,
	Compress:   "DEFLATE",
	Tiled:      true,
	BlockXSize: 256,
	BlockYSize: 256,
	NoData:     DefaultNoData,
}

func buildProfileConfig(region *Region, bands int, cfg OutputConfig) gdal.Profile {
	if cfg.DataType == 0 {
		cfg.DataType = gdal.Float64
	}
	if cfg.Format == "" {
		cfg.Format = "GTiff"
	}
	profile := gdal.DefaultProfile()
	profile = profile.Update(
		gdal.WithDimensions(region.XSize, region.YSize, bands),
		gdal.WithDataType(cfg.DataType),
		gdal.WithTransform(region.GeoTransform()),
		gdal.WithNodata(cfg.NoData),
	)

	if cfg.Format == "GTiff" || cfg.Format == "COG" {
		if cfg.Compress != "" {
			profile = profile.Update(gdal.WithCompress(cfg.Compress))
		}
		profile = profile.Update(gdal.WithTiled(cfg.Tiled))
		if cfg.BlockXSize > 0 {
			profile = profile.Update(gdal.WithBlockSize(cfg.BlockXSize, cfg.BlockYSize))
		}
	}

	crsStr := resolveOutputCRS(cfg, region)
	if crsStr != "" {
		profile = profile.Update(gdal.WithCRS(crsStr))
	}

	return profile
}

func verticalDatumEPSG(vd geoid.VerticalDatum) (code int, name string) {
	switch vd {
	case geoid.EGM84:
		return 5798, "EGM84 height"
	case geoid.EGM96:
		return 5773, "EGM96 geoid height"
	case geoid.EGM2008:
		return 3855, "EGM2008 geoid height"
	default:
		return 0, ""
	}
}

func resolveOutputCRS(cfg OutputConfig, region *Region) string {
	hCRS := strings.TrimSpace(cfg.CRS)
	if hCRS == "" && region != nil {
		if srs := region.SRS(); srs != nil {
			hCRS = strings.TrimSpace(srs.GetDef())
		}
	}
	if hCRS == "" {
		return ""
	}

	var hWkt string
	if strings.HasPrefix(hCRS, "EPSG:") || strings.HasPrefix(hCRS, "+proj=") {
		if crs, err := gdal.NewCRS(hCRS); err == nil {
			defer crs.SpatialReference().Destroy()
			if wkt, err := crs.ToWKT(); err == nil && wkt != "" {
				hWkt = wkt
			}
		}
	}
	if hWkt == "" {
		hWkt = hCRS
	}

	if cfg.VerticalDatum != geoid.HAE && cfg.VerticalDatum != geoid.UNKNOWN {
		vEPSG, vName := verticalDatumEPSG(cfg.VerticalDatum)
		if vEPSG > 0 {
			vCrs, err := gdal.NewCRS(fmt.Sprintf("EPSG:%d", vEPSG))
			if err == nil {
				defer vCrs.SpatialReference().Destroy()
				vWkt, err := vCrs.ToWKT()
				if err == nil && vWkt != "" {
					compound := fmt.Sprintf(
						`COMPD_CS["%s + %s",%s,%s]`,
						hCRS, vName, hWkt, vWkt,
					)
					return compound
				}
			}
		}
	}

	return hWkt
}

func buildProfile(region *Region, bands int, noData float64) gdal.Profile {
	cfg := DefaultGTiffConfig
	cfg.NoData = noData
	cfg.CRS = region.SRS().GetDef()
	return buildProfileConfig(region, bands, cfg)
}

func writeBand(ds gdal.Dataset, bandIdx int, data []float64, xSize, ySize int) error {
	band := ds.RasterBand(bandIdx)
	return band.IO(gdal.Write, 0, 0, xSize, ySize, data, xSize, ySize, 0, 0)
}

func CreateDEM(data []float64, region *Region, outputPath string, noData float64) error {
	return CreateDEMWithConfig(data, region, outputPath, OutputConfig{NoData: noData})

}

func CreateDEMWithConfig(data []float64, region *Region, outputPath string, cfg OutputConfig) error {
	if cfg.DataType == 0 {
		cfg.DataType = gdal.Float64
	}
	profile := buildProfileConfig(region, 1, cfg)
	return gdal.WithOutput(outputPath, profile, func(ds gdal.Dataset) error {
		return writeBand(ds, 1, data, region.XSize, region.YSize)
	})
}

func CreateStack(stackData [][]float64, region *Region, outputPath string, noData float64) error {
	count := len(stackData)
	if count == 0 {
		return fmt.Errorf("no stack bands provided")
	}
	profile := buildProfile(region, count, noData)
	return gdal.WithOutput(outputPath, profile, func(ds gdal.Dataset) error {
		for i := 0; i < count; i++ {
			if err := writeBand(ds, i+1, stackData[i], region.XSize, region.YSize); err != nil {
				return fmt.Errorf("band %d: %v", i+1, err)
			}
		}
		return nil
	})
}

func CreateRGB(pixels []uint8, region *Region, outputPath string) error {
	w, h := region.XSize, region.YSize
	if len(pixels) < 3*w*h {
		return fmt.Errorf("insufficient pixel data: need at least %d values for %dx%d RGB, got %d",
			3*w*h, w, h, len(pixels))
	}
	cfg := DefaultGTiffConfig
	cfg.DataType = gdal.Byte
	cfg.CRS = region.SRS().GetDef()

	n := w * h
	planes := [3][]float64{
		make([]float64, n),
		make([]float64, n),
		make([]float64, n),
	}
	for i := 0; i < n; i++ {
		planes[0][i] = float64(pixels[i*3])
		planes[1][i] = float64(pixels[i*3+1])
		planes[2][i] = float64(pixels[i*3+2])
	}

	profile := buildProfileConfig(region, 3, cfg)
	return gdal.WithOutput(outputPath, profile, func(ds gdal.Dataset) error {
		for band := 0; band < 3; band++ {
			bandData := ds.RasterBand(band + 1)
			if err := bandData.IO(gdal.Write, 0, 0, w, h, planes[band], w, h, 0, 0); err != nil {
				return fmt.Errorf("band %d: %v", band+1, err)
			}
		}
		return nil
	})
}

func parseCRSFromWKT(wkt string) geo.Proj {
	if wkt == "" {
		return nil
	}
	crs, err := gdal.NewCRS(wkt)
	if err != nil {
		return nil
	}
	defer crs.SpatialReference().Destroy()
	if proj4, err := crs.ToProj4(); err == nil && proj4 != "" {
		return geo.NewProj(proj4)
	}
	if epsg, ok := crs.ToEPSG(); ok {
		return geo.NewProj(fmt.Sprintf("EPSG:%d", epsg))
	}
	return nil
}

func ReadDEM(path string) ([]float64, *Region, error) {
	var data []float64
	var region *Region

	err := gdal.WithDatasetReadonly(path, func(ds gdal.Dataset) error {
		xSize := ds.RasterXSize()
		ySize := ds.RasterYSize()
		gt := ds.GeoTransform()

		if gt[2] != 0 || gt[4] != 0 {
			return fmt.Errorf("rotated rasters not supported in ReadDEM")
		}

		band := ds.RasterBand(1)
		var err error
		data, err = band.ReadWindow(0, 0, xSize, ySize, xSize, ySize, gdal.Nearest)
		if err != nil {
			return fmt.Errorf("read data: %v", err)
		}

		wkt := ds.Projection()
		srs := parseCRSFromWKT(wkt)
		if srs == nil {
			srs = geo.NewProj("EPSG:4326")
		}

		west := gt[0]
		north := gt[3]
		east := gt[0] + float64(xSize)*gt[1]
		south := gt[3] + float64(ySize)*gt[5]

		region = NewRegionFromBBox(
			west, south, east, north,
			srs, math.Abs(gt[1]), math.Abs(gt[5]),
		)
		region.XSize = xSize
		region.YSize = ySize
		return nil
	})

	return data, region, err
}

func ReadDEMBand(path string, bandIdx int) ([]float64, *Region, error) {
	var data []float64
	var region *Region

	err := gdal.WithDatasetReadonly(path, func(ds gdal.Dataset) error {
		if bandIdx < 1 || bandIdx > ds.RasterCount() {
			return fmt.Errorf("band %d out of range (1-%d)", bandIdx, ds.RasterCount())
		}
		xSize := ds.RasterXSize()
		ySize := ds.RasterYSize()
		gt := ds.GeoTransform()

		if gt[2] != 0 || gt[4] != 0 {
			return fmt.Errorf("rotated rasters not supported in ReadDEMBand")
		}

		band := ds.RasterBand(bandIdx)
		var err error
		data, err = band.ReadWindow(0, 0, xSize, ySize, xSize, ySize, gdal.Nearest)
		if err != nil {
			return fmt.Errorf("read band %d: %v", bandIdx, err)
		}

		wkt := ds.Projection()
		srs := parseCRSFromWKT(wkt)
		if srs == nil {
			srs = geo.NewProj("EPSG:4326")
		}

		region = NewRegionFromBBox(
			gt[0], gt[3]+float64(ySize)*gt[5],
			gt[0]+float64(xSize)*gt[1], gt[3],
			srs, math.Abs(gt[1]), math.Abs(gt[5]),
		)
		region.XSize = xSize
		region.YSize = ySize
		return nil
	})

	return data, region, err
}
