
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sdts"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211"
)

set(gdal_ogr_sdts_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sdts/ogrsdtsdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sdts/ogrsdtsdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sdts/ogrsdtslayer.cpp"
)

ADD_LIBRARY(gdal_ogr_sdts_frmt STATIC
            ${gdal_ogr_sdts_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_sdts_frmt_SOURCE_FILES})
