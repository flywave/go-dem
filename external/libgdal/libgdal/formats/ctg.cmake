INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ctg"
)

set(gdal_ctg_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ctg/ctgdataset.cpp"
)

ADD_LIBRARY(gdal_ctg_frmt STATIC
            ${gdal_ctg_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ctg_frmt_SOURCE_FILES})
