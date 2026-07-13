INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/mem"
)

set(gdal_ogr_mem_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/mem/ogrmemdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/mem/ogrmemdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/mem/ogrmemlayer.cpp"
)

ADD_LIBRARY(gdal_ogr_mem_frmt STATIC
            ${gdal_ogr_mem_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_idrisi_frmt_SOURCE_FILES})
