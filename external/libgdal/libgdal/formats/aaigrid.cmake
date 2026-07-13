INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/aaigrid")

set(gdal_aaigrid_frmt_SOURCE_FILES "${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/aaigrid/aaigriddataset.cpp")

ADD_LIBRARY(gdal_aaigrid_frmt STATIC
            ${gdal_aaigrid_frmt_SOURCE_FILES}
          )

TARGET_INCLUDE_DIRECTORIES(gdal_aaigrid_frmt PUBLIC $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/..>)

SOURCE_GROUP("src" FILES ${gdal_aaigrid_frmt_SOURCE_FILES})
