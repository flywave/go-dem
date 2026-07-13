INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/envisat"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw"
)

set(gdal_envisat_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/envisat/adsrange.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/envisat/envisatdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/envisat/EnvisatFile.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/envisat/records.c"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/envisat/unwrapgcps.cpp"
)

ADD_LIBRARY(gdal_envisat_frmt STATIC
            ${gdal_envisat_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_envisat_frmt_SOURCE_FILES})
