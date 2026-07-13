INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211"
)

set(gdal_iso8211_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211/ddffield.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211/ddffielddefn.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211/ddfmodule.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211/ddfrecord.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211/ddfsubfielddefn.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211/ddfutils.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211/8211createfromxml.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211/8211dump.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211/8211view.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211/mkcatalog.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211/timetest.cpp"
)

ADD_LIBRARY(gdal_iso8211_frmt STATIC
            ${gdal_iso8211_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_iso8211_frmt_SOURCE_FILES})
