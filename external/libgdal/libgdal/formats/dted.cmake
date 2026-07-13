INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/dted"
)

set(gdal_dted_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/dted/dted_api.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/dted/dted_create.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/dted/dted_ptstream.c"
	#"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/dted/dted_test.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/dted/dteddataset.cpp"
)

ADD_LIBRARY(gdal_dted_frmt STATIC
            ${gdal_dted_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_dted_frmt_SOURCE_FILES})
