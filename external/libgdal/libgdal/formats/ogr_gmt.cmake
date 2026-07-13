
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gmt"
)

set(gdal_ogr_gmt_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gmt/ogrgmtdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gmt/ogrgmtdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gmt/ogrgmtlayer.cpp"
)

ADD_LIBRARY(gdal_ogr_gmt_frmt STATIC
            ${gdal_ogr_gmt_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_gmt_frmt_SOURCE_FILES})
