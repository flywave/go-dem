INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hfa"
)

set(gdal_hfa_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hfa/hfatest.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hfa/hfafield.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hfa/hfaopen.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hfa/hfacompress.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hfa/hfadataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hfa/hfa_overviews.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hfa/hfadictionary.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hfa/hfaentry.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hfa/hfatype.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/hfa/hfaband.cpp"
)

ADD_LIBRARY(gdal_hfa_frmt STATIC
            ${gdal_hfa_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_hfa_frmt_SOURCE_FILES})
