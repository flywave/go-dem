INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jdem"
)

set(gdal_jdem_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jdem/jdemdataset.cpp"
)

ADD_LIBRARY(gdal_jdem_frmt STATIC
            ${gdal_jdem_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_jdem_frmt_SOURCE_FILES})
