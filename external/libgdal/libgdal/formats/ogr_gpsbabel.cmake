INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpsbabel"
)

set(gdal_ogr_gpsbabel_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpsbabel/ogrgpsbabeldatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpsbabel/ogrgpsbabeldriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpsbabel/ogrgpsbabelwritedatasource.cpp"
)

ADD_LIBRARY(gdal_ogr_gpsbabel_frmt STATIC
            ${gdal_ogr_gpsbabel_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_gpsbabel_frmt_SOURCE_FILES})
