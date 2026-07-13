INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/kmlsuperoverlay"
)

set(gdal_kmlsuperoverlay_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/kmlsuperoverlay/kmlsuperoverlaydataset.cpp"
)

ADD_LIBRARY(gdal_kmlsuperoverlay_frmt STATIC
            ${gdal_kmlsuperoverlay_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_kmlsuperoverlay_frmt_SOURCE_FILES})
