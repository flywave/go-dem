INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/leveller"
)

set(gdal_leveller_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/leveller/levellerdataset.cpp"
)

ADD_LIBRARY(gdal_leveller_frmt STATIC
            ${gdal_leveller_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_leveller_frmt_SOURCE_FILES})
