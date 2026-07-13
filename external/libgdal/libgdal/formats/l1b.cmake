INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/l1b"
)

set(gdal_l1b_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/l1b/l1bdataset.cpp"
)

ADD_LIBRARY(gdal_l1b_frmt STATIC
            ${gdal_l1b_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_l1b_frmt_SOURCE_FILES})
