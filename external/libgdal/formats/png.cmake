
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/../libpng" )
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/../zlib" )
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_BINARY_DIR}/../zlib" )

INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/png"
)

set(gdal_png_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/png/pngdataset.cpp"
)

ADD_LIBRARY(gdal_png_frmt STATIC
            ${gdal_png_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_png_frmt_SOURCE_FILES})
