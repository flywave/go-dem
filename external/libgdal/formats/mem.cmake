INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/mem"
)

set(gdal_mem_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/mem/memdataset.cpp"
)

ADD_LIBRARY(gdal_mem_frmt STATIC
            ${gdal_mem_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_mem_frmt_SOURCE_FILES})
