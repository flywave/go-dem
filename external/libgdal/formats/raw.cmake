
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw"
)

set(gdal_raw_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/gtxdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/roipacdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/landataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/ntv2dataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/loslasdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/ace2dataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/snodasdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/fastdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/doq2dataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/envidataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/eirdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/byndataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/hkvdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/atlsci_spheroid.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/btdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/pnmdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/fujibasdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/genbindataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/krodataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/idadataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/doq1dataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/dipxdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/gscdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/ehdrdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/iscedataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/ctable2dataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/lcpdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/cpgdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/ndfdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/mffdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/pauxdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/rrasterdataset.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw/ntv2dataset.cpp"
)

ADD_LIBRARY(gdal_raw_frmt STATIC
            ${gdal_raw_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_raw_frmt_SOURCE_FILES})
