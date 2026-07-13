INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/blx"
)

set(gdal_blx_frmt_SOURCE_FILES "${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/blx/blx.c"
"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/blx/blxdataset.cpp"
)

ADD_LIBRARY(gdal_blx_frmt STATIC
            ${gdal_blx_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_blx_frmt_SOURCE_FILES})
