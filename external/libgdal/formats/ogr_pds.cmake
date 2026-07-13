
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/pds"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pds"
)

set(gdal_ogr_pds_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/pds/ogrpdsdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/pds/ogrpdsdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/pds/ogrpdslayer.cpp"
)

ADD_LIBRARY(gdal_ogr_pds_frmt STATIC
            ${gdal_ogr_pds_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_pds_frmt_SOURCE_FILES})
