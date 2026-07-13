INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/dimap"
)

set(gdal_dimap_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/dimap/dimapdataset.cpp"
)

ADD_LIBRARY(gdal_dimap_frmt STATIC
            ${gdal_dimap_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_dimap_frmt_SOURCE_FILES})
