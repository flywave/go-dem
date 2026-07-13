INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/csv"
)

set(gdal_ogr_csv_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/csv/ogrcsvdatasource.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/csv/ogrcsvdriver.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/csv/ogrcsvlayer.cpp"
)

ADD_LIBRARY(gdal_ogr_csv_frmt STATIC
            ${gdal_ogr_csv_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_csv_frmt_SOURCE_FILES})
