
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/xyz"
)

set(gdal_xyz_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/xyz/xyzdataset.cpp"
)

ADD_LIBRARY(gdal_xyz_frmt STATIC
            ${gdal_xyz_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_xyz_frmt_SOURCE_FILES})
