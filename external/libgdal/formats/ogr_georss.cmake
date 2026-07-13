
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/georss"
)

set(gdal_ogr_georss_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/georss/ogrgeorssdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/georss/ogrgeorssdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/georss/ogrgeorsslayer.cpp"
)

ADD_LIBRARY(gdal_ogr_georss_frmt STATIC
            ${gdal_ogr_georss_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_georss_frmt_SOURCE_FILES})
