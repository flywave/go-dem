INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/geoconcept"
)

set(gdal_ogr_geoconcept_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/geoconcept/geoconcept.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/geoconcept/geoconcept_syscoord.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/geoconcept/ogrgeoconceptdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/geoconcept/ogrgeoconceptdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/geoconcept/ogrgeoconceptlayer.cpp"
)

ADD_LIBRARY(gdal_ogr_geoconcept_frmt STATIC
            ${gdal_ogr_geoconcept_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_geoconcept_frmt_SOURCE_FILES})
