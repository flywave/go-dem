add_definitions( -DAVCBIN_ENABLED=1)

INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc"
)

set(gdal_ogr_avc_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/avc_binwr.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/ogravcbinlayer.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/avc_e00read.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/avc_mbyte.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/ogravcbindriver.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/avc_misc.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/ogravce00driver.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/ogravce00datasource.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/avc_rawbin.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/ogravce00layer.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/avc_e00gen.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/avc_e00write.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/ogravclayer.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/avc_bin.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/ogravcbindatasource.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/avc_e00parse.cpp"
				"${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc/ogravcdatasource.cpp"
)

ADD_LIBRARY(gdal_ogr_avc_frmt STATIC
            ${gdal_ogr_avc_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_avc_frmt_SOURCE_FILES})
