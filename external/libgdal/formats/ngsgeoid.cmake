INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ngsgeoid"
)

set(gdal_ngsgeoid_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/ngsgeoid/ngsgeoiddataset.cpp"
)

ADD_LIBRARY(gdal_ngsgeoid_frmt STATIC
            ${gdal_ngsgeoid_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ngsgeoid_frmt_SOURCE_FILES})
