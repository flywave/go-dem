INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/msgn"
)

set(gdal_msgn_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/msgn/msg_basic_types.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/msgn/msg_reader_core.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/msgn/msgndataset.cpp"
)

ADD_LIBRARY(gdal_msgn_frmt STATIC
            ${gdal_msgn_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_msgn_frmt_SOURCE_FILES})
