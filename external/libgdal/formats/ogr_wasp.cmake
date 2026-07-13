
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/wasp"
)

set(gdal_ogr_wasp_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/wasp/ogrwaspdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/wasp/ogrwaspdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/wasp/ogrwasplayer.cpp"
)

ADD_LIBRARY(gdal_ogr_wasp_frmt STATIC
            ${gdal_ogr_wasp_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_wasp_frmt_SOURCE_FILES})
