
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sgi"
)

set(gdal_sgi_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sgi/sgidataset.cpp"
)

ADD_LIBRARY(gdal_sgi_frmt STATIC
            ${gdal_sgi_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_sgi_frmt_SOURCE_FILES})
