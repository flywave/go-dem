INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ceos"
)

set(gdal_ceos_frmt_SOURCE_FILES "${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ceos/ceosdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ceos/ceosopen.c"
	#"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ceos/ceostest.c"
)

ADD_LIBRARY(gdal_ceos_frmt STATIC
            ${gdal_ceos_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ceos_frmt_SOURCE_FILES})
