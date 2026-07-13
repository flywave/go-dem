
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/rmf"
)

set(gdal_rmf_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/rmf/rmfdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/rmf/rmfdem.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/rmf/rmfjpeg.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/rmf/rmflzw.cpp"
)

ADD_LIBRARY(gdal_rmf_frmt STATIC
            ${gdal_rmf_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_rmf_frmt_SOURCE_FILES})
