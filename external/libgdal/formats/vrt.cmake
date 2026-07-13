
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw"
)

set(gdal_vrt_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/vrtsourcedrasterband.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/vrtrasterband.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/pixelfunctions.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/vrtsources.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/vrtdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/vrtfilters.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/vrtpansharpened.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/vrtwarped.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/vrtdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/vrtderivedrasterband.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/vrtrawrasterband.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt/vrtmultidim.cpp"
)

ADD_LIBRARY(gdal_vrt_frmt STATIC
            ${gdal_vrt_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_vrt_frmt_SOURCE_FILES})
