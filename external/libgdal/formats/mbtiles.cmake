INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/mbtiles"
"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sqlite"
"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gpkg"
"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/mvt"
)

INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/../libsqlite3" )

set(gdal_mbtiles_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/mbtiles/mbtilesdataset.cpp"
)

ADD_LIBRARY(gdal_mbtiles_frmt STATIC
            ${gdal_mbtiles_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_mbtiles_frmt_SOURCE_FILES})
