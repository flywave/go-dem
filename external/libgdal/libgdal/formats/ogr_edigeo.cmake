INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/edigeo"
)

set(gdal_ogr_edigeo_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/edigeo/ogredigeodatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/edigeo/ogredigeodriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/edigeo/ogredigeolayer.cpp"
)

ADD_LIBRARY(gdal_ogr_edigeo_frmt STATIC
            ${gdal_ogr_edigeo_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_edigeo_frmt_SOURCE_FILES})
