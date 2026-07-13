
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/tsx"
)

set(gdal_tsx_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/tsx/tsxdataset.cpp"
)

ADD_LIBRARY(gdal_tsx_frmt STATIC
            ${gdal_tsx_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_tsx_frmt_SOURCE_FILES})
