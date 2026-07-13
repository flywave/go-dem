
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/osm"
  "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite"
)

INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/../libsqlite3" )

set(gdal_ogr_osm_frmt_SOURCE_FILES
  "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/osm/ogrosmlayer.cpp"
  "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/osm/ogrosmdriver.cpp"
  "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/osm/ogrosmdatasource.cpp"
  "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/osm/osm2osm.cpp"
  "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/osm/osm_parser.cpp"
)

ADD_LIBRARY(gdal_ogr_osm_frmt STATIC
            ${gdal_ogr_osm_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_osm_frmt_SOURCE_FILES})
