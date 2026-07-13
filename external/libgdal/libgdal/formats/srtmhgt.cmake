
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/srtmhgt"
)

set(gdal_srtmhgt_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/srtmhgt/srtmhgtdataset.cpp"
)

ADD_LIBRARY(gdal_srtmhgt_frmt STATIC
            ${gdal_srtmhgt_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_srtmhgt_frmt_SOURCE_FILES})
