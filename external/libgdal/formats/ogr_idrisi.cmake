INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/idrisi"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/idrisi"
)

set(gdal_ogr_idrisi_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/idrisi/ogridrisidatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/idrisi/ogridrisidriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/idrisi/ogridrisilayer.cpp"
)

ADD_LIBRARY(gdal_ogr_idrisi_frmt STATIC
            ${gdal_ogr_idrisi_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_idrisi_frmt_SOURCE_FILES})
