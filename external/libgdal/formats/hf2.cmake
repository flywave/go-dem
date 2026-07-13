INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hf2"
)

set(gdal_hf2_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hf2/hf2dataset.cpp"
)

ADD_LIBRARY(gdal_hf2_frmt STATIC
            ${gdal_hf2_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_hf2_frmt_SOURCE_FILES})
