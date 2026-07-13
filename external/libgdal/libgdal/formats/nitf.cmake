INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff"
)

set(gdal_nitf_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/nitfdes.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/rpftocdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/nitfdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/nitfrasterband.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/nitfwritejpeg.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/ecrgtocdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/mgrs.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/rpftocfile.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/nitfaridpcm.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/nitfimage.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/nitffile.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/nitfwritejpeg_12.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/nitf_gcprpc.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/nitfbilevel.cpp"
)
IF(NOT MINGW)
LIST(APPEND  gdal_nitf_frmt_SOURCE_FILES "${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/nitf/nitfdump.c")
ENDIF()

ADD_LIBRARY(gdal_nitf_frmt STATIC
            ${gdal_nitf_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_nitf_frmt_SOURCE_FILES})
