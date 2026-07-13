INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ingr"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff/libtiff"
)

set(gdal_ingr_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ingr/IngrTypes.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ingr/IntergraphBand.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ingr/IntergraphDataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ingr/JpegHelper.cpp"
)

ADD_LIBRARY(gdal_ingr_frmt STATIC
            ${gdal_ingr_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ingr_frmt_SOURCE_FILES})
