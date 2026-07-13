INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpx"
)

set(gdal_ogr_gpx_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpx/ogrgpxdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpx/ogrgpxdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpx/ogrgpxlayer.cpp"
)

ADD_LIBRARY(gdal_ogr_gpx_frmt STATIC
            ${gdal_ogr_gpx_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_gpx_frmt_SOURCE_FILES})
