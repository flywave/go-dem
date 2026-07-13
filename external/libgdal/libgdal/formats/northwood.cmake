INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/northwood"
)

set(gdal_northwood_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/northwood/grcdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/northwood/grddataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/northwood/northwood.cpp"
)

ADD_LIBRARY(gdal_northwood_frmt STATIC
            ${gdal_northwood_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_northwood_frmt_SOURCE_FILES})
