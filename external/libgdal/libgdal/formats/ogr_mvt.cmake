INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/mvt"
)

set(gdal_ogr_mvt_frmt_SOURCE_FILES
    "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/mvt/mvtutils.cpp"
    "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/mvt/ogrmvtdataset.cpp"
    "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/mvt/mvt_tile.cpp"
)

ADD_LIBRARY(gdal_ogr_mvt_frmt STATIC
            ${gdal_ogr_mvt_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_mvt_frmt_SOURCE_FILES})
