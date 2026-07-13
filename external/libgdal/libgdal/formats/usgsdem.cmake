
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/usgsdem"
)

set(gdal_usgsdem_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/usgsdem/usgsdem_create.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/usgsdem/usgsdemdataset.cpp"

)

ADD_LIBRARY(gdal_usgsdem_frmt STATIC
            ${gdal_usgsdem_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_usgsdem_frmt_SOURCE_FILES})
