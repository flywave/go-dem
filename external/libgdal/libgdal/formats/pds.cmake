
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pds"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw"
)

set(gdal_pds_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pds/isis2dataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pds/pdsdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pds/pds4vector.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pds/isis3dataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pds/vicardataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pds/vicarkeywordhandler.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pds/pds4dataset.cpp"
)

ADD_LIBRARY(gdal_pds_frmt STATIC
            ${gdal_pds_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_pds_frmt_SOURCE_FILES})
