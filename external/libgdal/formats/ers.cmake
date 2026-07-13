INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ers"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw"
)

set(gdal_ers_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ers/ersdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ers/ershdrnode.cpp"
)

ADD_LIBRARY(gdal_ers_frmt STATIC
            ${gdal_ers_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ers_frmt_SOURCE_FILES})
