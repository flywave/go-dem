
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iso8211"
)

set(gdal_sdts_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtsattrreader.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdts2shp.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtsindexedreader.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtsdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtsxref.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtsiref.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtscatd.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtslinereader.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtslib.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtsrasterreader.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtspointreader.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtspolygonreader.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts/sdtstransfer.cpp"
)

ADD_LIBRARY(gdal_sdts_frmt STATIC
            ${gdal_sdts_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_sdts_frmt_SOURCE_FILES})
