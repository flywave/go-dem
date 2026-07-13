INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gtm"
)

set(gdal_ogr_gtm_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gtm/gtm.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gtm/gtmtracklayer.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gtm/gtmwaypointlayer.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gtm/ogrgtmdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gtm/ogrgtmdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gtm/ogrgtmlayer.cpp"
)

ADD_LIBRARY(gdal_ogr_gtm_frmt STATIC
            ${gdal_ogr_gtm_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_gtm_frmt_SOURCE_FILES})
