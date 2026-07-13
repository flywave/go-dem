INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/coasp"
)

set(gdal_coasp_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/coasp/coasp_dataset.cpp"
)

ADD_LIBRARY(gdal_coasp_frmt STATIC
            ${gdal_coasp_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_coasp_frmt_SOURCE_FILES})
