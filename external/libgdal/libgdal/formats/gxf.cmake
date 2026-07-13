INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gxf"
)

set(gdal_gxf_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gxf/gxf_ogcwkt.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gxf/gxf_proj4.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gxf/gxfdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gxf/gxfopen.c"
)

ADD_LIBRARY(gdal_gxf_frmt STATIC
            ${gdal_gxf_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_gxf_frmt_SOURCE_FILES})
