INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/openfilegdb"
)

set(gdal_ogr_openfilegdb_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/openfilegdb/filegdbindex.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/openfilegdb/filegdbtable.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/openfilegdb/ogropenfilegdbdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/openfilegdb/ogropenfilegdbdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/openfilegdb/ogropenfilegdblayer.cpp"
)

ADD_LIBRARY(gdal_ogr_openfilegdb_frmt STATIC
            ${gdal_ogr_openfilegdb_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_openfilegdb_frmt_SOURCE_FILES})
