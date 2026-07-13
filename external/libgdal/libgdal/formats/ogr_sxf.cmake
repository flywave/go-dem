
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sxf"
)

set(gdal_ogr_sxf_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sxf/ogrsxfdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sxf/ogrsxfdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sxf/ogrsxflayer.cpp"
)

ADD_LIBRARY(gdal_ogr_sxf_frmt STATIC
            ${gdal_ogr_sxf_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_sxf_frmt_SOURCE_FILES})
