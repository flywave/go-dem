
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/til"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt"

)

set(gdal_til_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/til/tildataset.cpp"
)

ADD_LIBRARY(gdal_til_frmt STATIC
            ${gdal_til_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_til_frmt_SOURCE_FILES})
