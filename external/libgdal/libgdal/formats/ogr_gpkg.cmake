
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpkg"
"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite"
)

INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/../libsqlite3" )

set(gdal_ogr_gpkg_frmt_SOURCE_FILES
    "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpkg/gdalgeopackagerasterband.cpp"
    "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpkg/ogrgeopackagelayer.cpp"
    "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpkg/ogrgeopackagetablelayer.cpp"
    "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpkg/ogrgeopackagedriver.cpp"
    "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpkg/ogrgeopackageselectlayer.cpp"
    "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpkg/ogrgeopackagedatasource.cpp"
    "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpkg/ogrgeopackageutility.cpp"
)

ADD_LIBRARY(gdal_ogr_gpkg_frmt STATIC
            ${gdal_ogr_gpkg_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_gpkg_frmt_SOURCE_FILES})
