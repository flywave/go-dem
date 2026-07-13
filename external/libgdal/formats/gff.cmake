INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gff"
)

set(gdal_gff_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gff/gff_dataset.cpp"
)

ADD_LIBRARY(gdal_gff_frmt STATIC
            ${gdal_gff_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_gff_frmt_SOURCE_FILES})
