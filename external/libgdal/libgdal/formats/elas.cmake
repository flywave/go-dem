INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/elas"
)

set(gdal_elas_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/elas/elasdataset.cpp"
)

ADD_LIBRARY(gdal_elas_frmt STATIC
            ${gdal_elas_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_elas_frmt_SOURCE_FILES})
