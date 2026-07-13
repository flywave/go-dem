INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/aigrid")

set(gdal_aigrid_frmt_SOURCE_FILES "${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/aigrid/aigccitt.c"
"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/aigrid/aigdataset.cpp"
"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/aigrid/aigopen.c"
#"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/aigrid/aitest.c"
"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/aigrid/gridlib.c"
)

ADD_LIBRARY(gdal_aigrid_frmt STATIC
            ${gdal_aigrid_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_aigrid_frmt_SOURCE_FILES})
