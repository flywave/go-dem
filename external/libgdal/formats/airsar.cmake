INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/airsar")

set(gdal_airsar_frmt_SOURCE_FILES "${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/airsar/airsardataset.cpp"
)

ADD_LIBRARY(gdal_airsar_frmt STATIC
            ${gdal_airsar_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_aigrid_frmt_SOURCE_FILES})
