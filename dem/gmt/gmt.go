package gmt

/*
#cgo CFLAGS: -I${SRCDIR}/../../external/gmt/capi -I${SRCDIR}/../../external/gmt/src
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../../libs/darwin_arm -lgmt -lpsl -lm -L/opt/homebrew/lib -lnetcdf -L${SRCDIR}/../../libs/darwin_arm -lhdf5 -lhdf5_hl -lgdal_c -lgdal_aaigrid_frmt -lgdal_adrg_frmt -lgdal_aigrid_frmt -lgdal_airsar_frmt -lgdal_blx_frmt -lgdal_bmp_frmt -lgdal_ceos_frmt -lgdal_ceos2_frmt -lgdal_coasp_frmt -lgdal_cosar_frmt -lgdal_ctg_frmt -lgdal_dimap_frmt -lgdal_dted_frmt -lgdal_elas_frmt -lgdal_envisat_frmt -lgdal_ers_frmt -lgdal_fit_frmt -lgdal_gff_frmt -lgdal_gsg_frmt -lgdal_gtiff_frmt -lgdal_hf2_frmt -lgdal_hfa_frmt -lgdal_idrisi_frmt -lgdal_ilwis_frmt -lgdal_ingr_frmt -lgdal_iris_frmt -lgdal_iso8211_frmt -lgdal_jaxapalsar_frmt -lgdal_jdem_frmt -lgdal_jpeg_frmt -lgdal_kmlsuperoverlay_frmt -lgdal_l1b_frmt -lgdal_leveller_frmt -lgdal_map_frmt -lgdal_mbtiles_frmt -lgdal_mem_frmt -lgdal_ngsgeoid_frmt -lgdal_nitf_frmt -lgdal_northwood_frmt -lgdal_ogr_avc_frmt -lgdal_ogr_csv_frmt -lgdal_ogr_dgn_frmt -lgdal_ogr_dxf_frmt -lgdal_ogr_edigeo_frmt -lgdal_ogr_geoconcept_frmt -lgdal_ogr_geojson_frmt -lgdal_ogr_georss_frmt -lgdal_ogr_gml_frmt -lgdal_ogr_gmt_frmt -lgdal_ogr_gpkg_frmt -lgdal_ogr_gpsbabel_frmt -lgdal_ogr_gpx_frmt -lgdal_ogr_gtm_frmt -lgdal_ogr_idrisi_frmt -lgdal_ogr_kml_frmt -lgdal_ogr_mem_frmt -lgdal_ogr_mitab_frmt -lgdal_ogr_mvt_frmt -lgdal_ogr_ntf_frmt -lgdal_ogr_openfilegdb_frmt -lgdal_ogr_osm_frmt -lgdal_ogr_pds_frmt -lgdal_ogr_pgdump_frmt -lgdal_ogr_rec_frmt -lgdal_ogr_s57_frmt -lgdal_ogr_sdts_frmt -lgdal_ogr_shape_frmt -lgdal_ogr_sqlite_frmt -lgdal_ogr_svg_frmt -lgdal_ogr_sxf_frmt -lgdal_ogr_vrt_frmt -lgdal_ogr_wasp_frmt -lgdal_pcidsk_frmt -lgdal_pds_frmt -lgdal_png_frmt -lgdal_r_frmt -lgdal_raw_frmt -lgdal_rmf_frmt -lgdal_rs2_frmt -lgdal_saga_frmt -lgdal_sdts_frmt -lgdal_sgi_frmt -lgdal_srtmhgt_frmt -lgdal_terragen_frmt -lgdal_til_frmt -lgdal_tsx_frmt -lgdal_usgsdem_frmt -lgdal_vrt_frmt -lgdal_webp_frmt -lgdal_xpm_frmt -lgdal_xyz_frmt -lgdal_zmap_frmt -lproj -lgeos -lsqlite3 -lzlib -lpng -ljpeg -lexpat -liconv -lwebp -lc++ -framework Accelerate
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../../libs/linux -lgmt -lpsl -lm -lnetcdf -lhdf5 -lhdf5_hl -lgdal_c -lproj -lgeos -lsqlite3 -lz -lpng -ljpeg -lexpat -liconv
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../../libs/linux_arm -lgmt -lpsl -lm -lnetcdf -lhdf5 -lhdf5_hl -lgdal_c -lproj -lgeos -lsqlite3 -lz -lpng -ljpeg -lexpat -liconv

#include <stdlib.h>
#include "gmt_capi.h"
#include "../../external/gmt/capi/gmt_capi.c"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type GridConfig struct {
	XInc, YInc   float64
	XMin, XMax   float64
	YMin, YMax   float64
	Tension      float64
	SearchRadius float64
	EmptyValue   int
}

func init() {
	C.gdemo_gmt_begin()
}

func Surface(inputPath, outputPath string, cfg *GridConfig) error {
	if cfg.Tension <= 0 {
		cfg.Tension = 0.25
	}
	cIn := C.CString(inputPath)
	cOut := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cIn))
	defer C.free(unsafe.Pointer(cOut))
	ret := C.gmt_surface(cIn, cOut,
		C.double(cfg.Tension),
		C.double(cfg.XInc), C.double(cfg.YInc),
		C.double(cfg.XMin), C.double(cfg.XMax),
		C.double(cfg.YMin), C.double(cfg.YMax))
	if ret != 0 {
		return fmt.Errorf("gmt surface failed with code %d", int(ret))
	}
	return nil
}

func Grdfilter(inputPath, outputPath, filterType, distFlag string) error {
	cIn := C.CString(inputPath)
	cOut := C.CString(outputPath)
	cFilt := C.CString(filterType)
	cDist := C.CString(distFlag)
	defer C.free(unsafe.Pointer(cIn))
	defer C.free(unsafe.Pointer(cOut))
	defer C.free(unsafe.Pointer(cFilt))
	defer C.free(unsafe.Pointer(cDist))
	ret := C.gmt_grdfilter(cIn, cOut, cFilt, cDist)
	if ret != 0 {
		return fmt.Errorf("gmt grdfilter failed with code %d", int(ret))
	}
	return nil
}

func Triangulate(inputPath, outputPath string, cfg *GridConfig) error {
	cIn := C.CString(inputPath)
	cOut := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cIn))
	defer C.free(unsafe.Pointer(cOut))
	ret := C.gmt_triangulate(cIn, cOut,
		C.double(cfg.XInc), C.double(cfg.YInc),
		C.double(cfg.XMin), C.double(cfg.XMax),
		C.double(cfg.YMin), C.double(cfg.YMax))
	if ret != 0 {
		return fmt.Errorf("gmt triangulate failed with code %d", int(ret))
	}
	return nil
}

func Blockmean(inputPath, outputPath string, cfg *GridConfig) error {
	cIn := C.CString(inputPath)
	cOut := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cIn))
	defer C.free(unsafe.Pointer(cOut))
	ret := C.gmt_blockmean(cIn, cOut,
		C.double(cfg.XInc), C.double(cfg.YInc),
		C.double(cfg.XMin), C.double(cfg.XMax),
		C.double(cfg.YMin), C.double(cfg.YMax))
	if ret != 0 {
		return fmt.Errorf("gmt blockmean failed with code %d", int(ret))
	}
	return nil
}

func Nearneighbor(inputPath, outputPath string, cfg *GridConfig) error {
	cIn := C.CString(inputPath)
	cOut := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cIn))
	defer C.free(unsafe.Pointer(cOut))
	sr := cfg.SearchRadius
	if sr <= 0 {
		sr = cfg.XInc * 5
	}
	if cfg.EmptyValue == 0 {
		cfg.EmptyValue = -9999
	}
	ret := C.gmt_nearneighbor(cIn, cOut,
		C.double(cfg.XInc), C.double(cfg.YInc),
		C.double(cfg.XMin), C.double(cfg.XMax),
		C.double(cfg.YMin), C.double(cfg.YMax),
		C.double(sr), C.int(cfg.EmptyValue))
	if ret != 0 {
		return fmt.Errorf("gmt nearneighbor failed with code %d", int(ret))
	}
	return nil
}
