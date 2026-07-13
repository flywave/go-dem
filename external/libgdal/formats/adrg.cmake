INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/adrg"
"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211"
)

set(gdal_adrg_frmt_SOURCE_FILES "${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/adrg/adrgdataset.cpp"
"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/adrg/srpdataset.cpp"
)

ADD_LIBRARY(gdal_adrg_frmt STATIC
            ${gdal_adrg_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_adrg_frmt_SOURCE_FILES})
