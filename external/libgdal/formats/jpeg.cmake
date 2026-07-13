
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jpeg"
"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jpeg/libjpeg"
)

set(gdal_jpeg_frmt_SOURCE_FILES
  	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jpeg/jpgdataset_12.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jpeg/jpgdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jpeg/vsidataio_12.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jpeg/vsidataio.cpp"
)

ADD_LIBRARY(gdal_jpeg_frmt STATIC
            ${gdal_jpeg_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_jpeg_frmt_SOURCE_FILES})
