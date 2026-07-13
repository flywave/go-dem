

INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/rs2"
)

set(gdal_rs2_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/rs2/rs2dataset.cpp"
)

ADD_LIBRARY(gdal_rs2_frmt STATIC
            ${gdal_rs2_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_rs2_frmt_SOURCE_FILES})
