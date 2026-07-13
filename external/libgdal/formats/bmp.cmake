INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/bmp"
)

set(gdal_bmp_frmt_SOURCE_FILES "${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/bmp/bmpdataset.cpp"
)

ADD_LIBRARY(gdal_bmp_frmt STATIC
            ${gdal_bmp_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_bmp_frmt_SOURCE_FILES})
