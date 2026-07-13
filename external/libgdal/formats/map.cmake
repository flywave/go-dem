INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/map"
)

set(gdal_map_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/map/mapdataset.cpp"
)

ADD_LIBRARY(gdal_map_frmt STATIC
            ${gdal_map_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_map_frmt_SOURCE_FILES})
