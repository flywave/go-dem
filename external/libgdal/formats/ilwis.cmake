
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ilwis"
)

set(gdal_ilwis_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ilwis/ilwiscoordinatesystem.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ilwis/ilwisdataset.cpp"
)

ADD_LIBRARY(gdal_ilwis_frmt STATIC
            ${gdal_ilwis_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ilwis_frmt_SOURCE_FILES})
