INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ceos2"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw"
)

set(gdal_ceos2_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ceos2/ceos.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ceos2/ceosrecipe.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ceos2/ceossar.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ceos2/link.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ceos2/sar_ceosdataset.cpp"
)

ADD_LIBRARY(gdal_ceos2_frmt STATIC
            ${gdal_ceos2_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ceos2_frmt_SOURCE_FILES})
