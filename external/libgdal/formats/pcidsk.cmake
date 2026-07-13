
INCLUDE_DIRECTORIES(
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/raw"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jpeg/libjpeg"
)

set(gdal_pcidsk_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/gdal_edb.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/pcidskdataset2.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/ogrpcidsklayer.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/blockdir/asciitiledir.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/blockdir/asciitilelayer.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/blockdir/binarytiledir.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/blockdir/binarytilelayer.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/blockdir/blockdir.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/blockdir/blockfile.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/blockdir/blocklayer.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/blockdir/blocktiledir.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/blockdir/blocktilelayer.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/edb_pcidsk.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/pcidsk_pubutils.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/cpcidskblockfile.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/pcidskbuffer.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/clinksegment.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/pcidskexception.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/pcidskinterfaces.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/metadataset_p.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/pcidskcreate.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/libjpeg_io.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/cpcidskfile.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/pcidsk_utils.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/pcidskopen.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/pcidsk_raster.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/core/pcidsk_scanint.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidsk_array.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskbinarysegment.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskbitmap.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskblut.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskbpct.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskephemerissegment.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskgcp2segment.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskgeoref.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidsklut.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskpct.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskpolymodel.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskrpcmodel.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidsksegment.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidsk_tex.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidsktoutinmodel.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskvectorsegment_consistencycheck.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/cpcidskvectorsegment.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/metadatasegment_p.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/systiledir.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/vecsegdataindex.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/segment/vecsegheader.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/channel/cpixelinterleavedchannel.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/channel/cbandinterleavedchannel.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/channel/cexternalchannel.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/channel/ctiledchannel.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/channel/cpcidskchannel.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/port/io_stdio.cpp"
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/vsi_pcidsk_io.cpp"
)
IF (WIN32)
  LIST(APPEND gdal_pcidsk_frmt_SOURCE_FILES
		"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/port/io_win32.cpp"
		"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/port/win32_mutex.cpp"
	)
ELSE()
  LIST(APPEND gdal_pcidsk_frmt_SOURCE_FILES "${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/pcidsk/sdk/port/pthread_mutex.cpp")
ENDIF()


ADD_LIBRARY(gdal_pcidsk_frmt STATIC
            ${gdal_pcidsk_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_pcidsk_frmt_SOURCE_FILES})
