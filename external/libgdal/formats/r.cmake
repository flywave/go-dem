
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/r"
)

set(gdal_r_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/r/rcreatecopy.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/r/rdataset.cpp"
)

ADD_LIBRARY(gdal_r_frmt STATIC
            ${gdal_r_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_r_frmt_SOURCE_FILES})
