INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/fit"
)

set(gdal_fit_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/fit/fit.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/fit/fitdataset.cpp"
)

ADD_LIBRARY(gdal_fit_frmt STATIC
            ${gdal_fit_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ers_frmt_SOURCE_FILES})
