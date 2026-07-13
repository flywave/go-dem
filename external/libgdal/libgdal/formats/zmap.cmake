
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/zmap"
)

set(gdal_zmap_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/zmap/zmapdataset.cpp"
)

ADD_LIBRARY(gdal_zmap_frmt STATIC
            ${gdal_zmap_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_zmap_frmt_SOURCE_FILES})
