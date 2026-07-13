IF(UNIX)
  SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wno-float-conversion -Wno-shadow -Wno-unused-function -Wno-format-extra-args -Wno-unknown-pragmas -Wno-switch -Wno-tautological-compare -Wno-attributes -Wno-unused-const-variable -Wno-sign-compare -Wno-deprecated-declarations")
ENDIF()

INCLUDE_DIRECTORIES(
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/third_party/LercLib"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jpeg/libjpeg"
)

INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/../zlib" )
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_BINARY_DIR}/../zlib")

set(gdal_gtiff_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/cogdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/gt_jpeg_copy.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_ojpeg.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_dirwrite.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_tile.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_vsi.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_packbits.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_warning.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_webp.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_lzma.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_thunder.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_swab.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_compress.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_codec.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_luv.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_open.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_dumpmode.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_strip.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_print.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_flush.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_dirinfo.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_fax3.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_pixarlog.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_getimage.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_dir.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_write.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_jpeg_12.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_fax3sm.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_predict.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_aux.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_close.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_extension.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_color.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_jpeg.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_dirread.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_version.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_zip.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_read.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_error.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_zstd.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_lzw.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff/tif_next.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_names.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/xtiff.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_new.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_print.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_extra.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_free.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geotiff_proj4.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_trans.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_get.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_tiffp.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_set.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_simpletags.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_normalize.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libgeotiff/geo_write.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/geotiff.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/tifvsi.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/gt_wkt_srs.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/tif_float.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/gt_citation.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/gt_overview.cpp"
)

ADD_LIBRARY(gdal_gtiff_frmt STATIC
            ${gdal_gtiff_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_gtiff_frmt_SOURCE_FILES})
