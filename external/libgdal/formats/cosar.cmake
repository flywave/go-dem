INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/cosar"
)

set(gdal_cosar_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/cosar/cosar_dataset.cpp"
)

ADD_LIBRARY(gdal_cosar_frmt STATIC
            ${gdal_cosar_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_cosar_frmt_SOURCE_FILES})
