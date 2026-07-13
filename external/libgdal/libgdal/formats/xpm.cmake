
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/mem"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/xpm"
)

set(gdal_xpm_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/xpm/xpmdataset.cpp"
)

ADD_LIBRARY(gdal_xpm_frmt STATIC
            ${gdal_xpm_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_xpm_frmt_SOURCE_FILES})
