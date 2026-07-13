INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/idrisi"
)

set(gdal_idrisi_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/idrisi/IdrisiDataset.cpp"
)

ADD_LIBRARY(gdal_idrisi_frmt STATIC
            ${gdal_idrisi_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_idrisi_frmt_SOURCE_FILES})
