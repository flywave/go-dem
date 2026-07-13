
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/terragen"
)

set(gdal_terragen_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/terragen/terragendataset.cpp"
)

ADD_LIBRARY(gdal_terragen_frmt STATIC
            ${gdal_terragen_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_terragen_frmt_SOURCE_FILES})
