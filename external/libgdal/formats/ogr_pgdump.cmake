
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/pgdump"
)

set(gdal_ogr_pgdump_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/pgdump/ogrpgdumpdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/pgdump/ogrpgdumpdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/pgdump/ogrpgdumplayer.cpp"
)

ADD_LIBRARY(gdal_ogr_pgdump_frmt STATIC
            ${gdal_ogr_pgdump_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_pgdump_frmt_SOURCE_FILES})
