
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/../libpng" )
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/../zlib" )
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_BINARY_DIR}/../zlib" )

INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/webp"
)

set(gdal_webp_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/webp/webpdataset.cpp"
)

ADD_LIBRARY(gdal_webp_frmt STATIC
            ${gdal_webp_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_webp_frmt_SOURCE_FILES})
