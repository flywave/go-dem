
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/vrt"
)

set(gdal_ogr_vrt_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/vrt/ogrvrtdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/vrt/ogrvrtdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/vrt/ogrvrtlayer.cpp"
)

ADD_LIBRARY(gdal_ogr_vrt_frmt STATIC
            ${gdal_ogr_vrt_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_vrt_frmt_SOURCE_FILES})
