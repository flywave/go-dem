
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/saga"
)

set(gdal_saga_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/saga/sagadataset.cpp"
)

ADD_LIBRARY(gdal_saga_frmt STATIC
            ${gdal_saga_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_saga_frmt_SOURCE_FILES})
