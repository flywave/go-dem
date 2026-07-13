INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jaxapalsar"
)

set(gdal_jaxapalsar_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jaxapalsar/jaxapalsardataset.cpp"
)

ADD_LIBRARY(gdal_jaxapalsar_frmt STATIC
            ${gdal_jaxapalsar_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_jaxapalsar_frmt_SOURCE_FILES})
