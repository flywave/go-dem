INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gsg"
)

set(gdal_gsg_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gsg/gs7bgdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gsg/gsagdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gsg/gsbgdataset.cpp"
)

ADD_LIBRARY(gdal_gsg_frmt STATIC
            ${gdal_gsg_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_gsg_frmt_SOURCE_FILES})
