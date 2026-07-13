IF(UNIX)
  SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wno-incompatible-pointer-types -Wno-implicit-function-declaration -Wno-unused-result  -Wno-shadow -Wno-unused-function -Wno-format-extra-args -Wno-unknown-pragmas -Wno-tautological-compare -Wno-attributes -Wno-unused-const-variable -Wno-sign-compare -Wno-deprecated-declarations")
IF("${CMAKE_CXX_COMPILER_ID}" STREQUAL "GNU")
 SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wno-unused-but-set-variable")
 SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall")
ENDIF()
ENDIF()

SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS}  -fpermissive")

 
IF(APPLE)
 SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=gnu++14")
ELSE()
SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fPIC  -std=c++14")
ENDIF()
 

IF(MINGW)
  SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wno-stringop-overflow -Wno-fpermissive -Wno-gnu-designator")
  SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wno-stringop-overflow -Wno-fpermissive -Wno-gnu-designator")
ENDIF()

IF (WIN32)
  INCLUDE_DIRECTORIES(BEFORE ${CMAKE_CURRENT_SOURCE_DIR}/arch/win/)
ELSEIF(APPLE)
  INCLUDE_DIRECTORIES(BEFORE ${CMAKE_CURRENT_SOURCE_DIR}/arch/osx/)
ELSE()
  INCLUDE_DIRECTORIES(BEFORE ${CMAKE_CURRENT_SOURCE_DIR}/arch/unix/)
ENDIF()

INCLUDE_DIRECTORIES(BEFORE
  ${CMAKE_CURRENT_SOURCE_DIR}/../libproj/proj/include
  ${CMAKE_CURRENT_SOURCE_DIR}/../libproj/proj/src
)

INCLUDE_DIRECTORIES(BEFORE
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/alg/marching_squares
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/alg
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/gcore
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/apps
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/port
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/gtiff
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/vrt
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/mem
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/jpeg
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/adrg
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/sdts
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/mem
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/generic
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/shape
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/avc
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/geojson
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/gpkg
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/geojson/libjson
  ${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/s57

)

LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../zlib" )
LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_BINARY_DIR}/../zlib" )
LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libiconv/include" )
LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libpng" )
LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../webp/include" )
LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libsqlite3" )
LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_BINARY_DIR}/../libsqlite3" )
IF (WIN32)

  LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libproj/arch/win/src")
  LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libgeos/arch/win")
  LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libexpat/arch/win")
ELSE()
  LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libproj/arch/unix/src")
  LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libgeos/arch/unix")
  LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libexpat/arch/unix")
ENDIF()

LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libgeos/arch/common")
LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libgeos/geos/include")
LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libgeos/geos/capi")
LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libproj/proj/src")
LIST(APPEND GDAL_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/../libexpat/expat/lib")

INCLUDE_DIRECTORIES(BEFORE ${GDAL_INCLUDE_DIRS})

add_definitions(
  -D_LARGEFILE_SOURCE
  -D_FILE_OFFSET_BITS=64
  -DINTERNAL_LIBTIFF=1
  -DINTERNAL_LIBGEOTIFF=1
  -DGDAL_COMPILATION=1
  -DPAM_ENABLED=1
  -DOGR_ENABLED=1
  -DHAVE_EXPAT=1
  -DHAVE_LIBPROJ=1
  -DHAVE_GEOS=1
  -DPROJ_STATIC=1
  -DCPU_LITTLE_ENDIAN=1
  -DBIGTIFF_SUPPORT=1
  -DZIP_SUPPORT=1
  -DJPEG_SUPPORT=1
)

IF(MINGW)
add_definitions(
  -DDONT_DEPRECATE_SPRINTF
  -DHAVE_LONG_LONG
)
ENDIF()

set(gdal_SOURCE_FILES
  gdal/apps/ogr2ogr_lib.cpp
  gdal/apps/gdal_rasterize_lib.cpp
  gdal/apps/gdal_translate_lib.cpp
  gdal/apps/gdalwarp_lib.cpp
  gdal/apps/gdalbuildvrt_lib.cpp
  gdal/apps/commonutils.cpp

  gdal/frmts/gdalallregister.cpp
  gdal/frmts/derived/deriveddataset.cpp
  gdal/frmts/derived/derivedlist.c

  gdal/ogr/swq_select.cpp
  gdal/ogr/swq_parser.cpp
  gdal/ogr/swq_op_registrar.cpp
  gdal/ogr/swq_op_general.cpp
  gdal/ogr/swq_expr_node.cpp
  gdal/ogr/swq.cpp
  gdal/ogr/ogrutils.cpp
  gdal/ogr/ogrtriangulatedsurface.cpp
  gdal/ogr/ogrtriangle.cpp
  gdal/ogr/ogrsurface.cpp
  gdal/ogr/ogrspatialreference.cpp
  gdal/ogr/ogrpolyhedralsurface.cpp
  gdal/ogr/ogrpolygon.cpp
  gdal/ogr/ogrpoint.cpp
  gdal/ogr/ogrpgeogeometry.cpp
  gdal/ogr/ogrmultisurface.cpp
  gdal/ogr/ogrmultipolygon.cpp
  gdal/ogr/ogrmultipoint.cpp
  gdal/ogr/ogrmultilinestring.cpp
  gdal/ogr/ogrmulticurve.cpp
  gdal/ogr/ogrlinestring.cpp
  gdal/ogr/ogrlinearring.cpp
  gdal/ogr/ogrgeomfielddefn.cpp
  gdal/ogr/ogrgeometryfactory.cpp
  gdal/ogr/ogrgeometrycollection.cpp
  gdal/ogr/ogrgeometry.cpp
  gdal/ogr/ogrgeomediageometry.cpp
  gdal/ogr/ogrfielddefn.cpp
  gdal/ogr/ogrfeaturestyle.cpp
  gdal/ogr/ogrfeaturequery.cpp
  gdal/ogr/ogrfeaturedefn.cpp
  gdal/ogr/ogrfeature.cpp
  gdal/ogr/ogrcurvepolygon.cpp
  gdal/ogr/ogrcurvecollection.cpp
  gdal/ogr/ogrcurve.cpp
  gdal/ogr/ogrct.cpp
  gdal/ogr/ogrcompoundcurve.cpp
  gdal/ogr/ogrcircularstring.cpp
  gdal/ogr/ograssemblepolygon.cpp
  gdal/ogr/ograpispy.cpp
  gdal/ogr/ogr_xerces.cpp
  gdal/ogr/ogr_srsnode.cpp
  gdal/ogr/ogr_srs_xml.cpp
  gdal/ogr/ogr_srs_usgs.cpp
  gdal/ogr/ogr_srs_pci.cpp
  gdal/ogr/ogr_srs_panorama.cpp
  gdal/ogr/ogr_srs_ozi.cpp
  gdal/ogr/ogr_srs_esri.cpp
  gdal/ogr/ogr_srs_erm.cpp
  gdal/ogr/ogr_srs_dict.cpp
  gdal/ogr/ogr_proj_p.cpp
  gdal/ogr/ogr_geocoding.cpp
  gdal/ogr/ogr_geo_utils.cpp
  gdal/ogr/ogr_fromepsg.cpp
  gdal/ogr/ogr_expat.cpp
  gdal/ogr/ogr_api.cpp
  gdal/ogr/ogr2gmlgeometry.cpp
  gdal/ogr/gml2ogrgeometry.cpp
  # gdal/ogr/ogr_opt.cpp
  # gdal/ogr/generate_encoding_table.c
  gdal/ogr/ogrsf_frmts/generic/ogr_gensql.cpp
  gdal/ogr/ogrsf_frmts/generic/ogremulatedtransaction.cpp
  gdal/ogr/ogrsf_frmts/generic/ogr_miattrind.cpp
  gdal/ogr/ogrsf_frmts/generic/ogrunionlayer.cpp
  gdal/ogr/ogrsf_frmts/generic/ogrsfdriver.cpp
  gdal/ogr/ogrsf_frmts/generic/ogrregisterall.cpp
  gdal/ogr/ogrsf_frmts/generic/ogreditablelayer.cpp
  gdal/ogr/ogrsf_frmts/generic/ogr_attrind.cpp
  gdal/ogr/ogrsf_frmts/generic/ogrdatasource.cpp
  gdal/ogr/ogrsf_frmts/generic/ogrmutexedlayer.cpp
  gdal/ogr/ogrsf_frmts/generic/ogrwarpedlayer.cpp
  gdal/ogr/ogrsf_frmts/generic/ogrlayerpool.cpp
  gdal/ogr/ogrsf_frmts/generic/ogrsfdriverregistrar.cpp
  gdal/ogr/ogrsf_frmts/generic/ogrmutexeddatasource.cpp
  gdal/ogr/ogrsf_frmts/generic/ogrlayer.cpp
  gdal/ogr/ogrsf_frmts/generic/ogrlayerdecorator.cpp

  gdal/alg/contour.cpp         
  gdal/alg/gdalrasterpolygonenumerator.cpp
  gdal/alg/gdalapplyverticalshiftgrid.cpp  
  gdal/alg/gdal_rpc.cpp
  gdal/alg/gdalchecksum.cpp    
  gdal/alg/gdalsievefilter.cpp
  gdal/alg/gdal_crs.cpp        
  gdal/alg/gdal_simplesurf.cpp
  gdal/alg/gdalcutline.cpp     
  gdal/alg/gdalsimplewarp.cpp
  gdal/alg/gdaldither.cpp      
  gdal/alg/gdal_tps.cpp
  gdal/alg/gdalgeoloc.cpp      
  gdal/alg/gdaltransformer.cpp
  gdal/alg/gdalgridavx.cpp     
  gdal/alg/gdaltransformgeolocs.cpp
  gdal/alg/gdalgrid.cpp        
  gdal/alg/gdalwarper.cpp
  gdal/alg/gdalgridsse.cpp     
  gdal/alg/gdalwarpkernel.cpp
  gdal/alg/gdallinearsystem.cpp      
  gdal/alg/gdalwarpkernel_opencl.cpp
  gdal/alg/gdalmatching.cpp    
  gdal/alg/gdalwarpoperation.cpp
  gdal/alg/gdalmediancut.cpp         
  gdal/alg/llrasterize.cpp
  gdal/alg/gdal_octave.cpp     
  gdal/alg/polygonize.cpp
  gdal/alg/gdalpansharpen.cpp        
  gdal/alg/rasterfill.cpp
  gdal/alg/gdalproximity.cpp         
  gdal/alg/thinplatespline.cpp
  gdal/alg/gdalrasterize.cpp         
  gdal/alg/viewshed.cpp
  gdal/alg/delaunay.c

  gdal/gcore/gdalabstractbandblockcache.cpp  
  gdal/gcore/gdalnodatamaskband.cpp
  gdal/gcore/gdalallvalidmaskband.cpp  
  gdal/gcore/gdalnodatavaluesmaskband.cpp
  gdal/gcore/gdalarraybandblockcache.cpp     
  gdal/gcore/gdalopeninfo.cpp
  gdal/gcore/gdalcolortable.cpp        
  gdal/gcore/gdaloverviewdataset.cpp
  gdal/gcore/gdaldataset.cpp     
  gdal/gcore/gdalpamdataset.cpp
  gdal/gcore/gdaldefaultasync.cpp      
  gdal/gcore/gdalpamproxydb.cpp
  gdal/gcore/gdaldefaultoverviews.cpp  
  gdal/gcore/gdalpamrasterband.cpp
  gdal/gcore/gdaldllmain.cpp     
  gdal/gcore/gdalproxydataset.cpp
  gdal/gcore/gdaldriver.cpp      
  gdal/gcore/gdalproxypool.cpp
  gdal/gcore/gdaldrivermanager.cpp     
  gdal/gcore/gdalpython.cpp
  gdal/gcore/gdalexif.cpp        
  gdal/gcore/gdalpythondriverloader.cpp
  gdal/gcore/gdalgeorefpamdataset.cpp  
  gdal/gcore/gdalrasterband.cpp
  gdal/gcore/gdalhashsetbandblockcache.cpp   
  gdal/gcore/gdalrasterblock.cpp
  gdal/gcore/gdaljp2abstractdataset.cpp      
  gdal/gcore/gdal_rat.cpp
  gdal/gcore/gdaljp2box.cpp      
  gdal/gcore/gdalrescaledalphaband.cpp
  gdal/gcore/gdaljp2metadata.cpp       
  gdal/gcore/gdal_thread_pool.cpp
  gdal/gcore/gdaljp2metadatagenerator.cpp    
  gdal/gcore/gdalvirtualmem.cpp
  gdal/gcore/gdaljp2structure.cpp      
  gdal/gcore/nasakeywordhandler.cpp
  gdal/gcore/gdalmajorobject.cpp       
  gdal/gcore/overview.cpp
  gdal/gcore/gdal_mdreader.cpp         
  gdal/gcore/rasterio.cpp
  gdal/gcore/gdal_misc.cpp       
  gdal/gcore/rasterio_ssse3.cpp
  gdal/gcore/gdalmultidim.cpp    
  gdal/gcore/rawdataset.cpp
  gdal/gcore/gdalmultidomainmetadata.cpp     
  gdal/gcore/tilematrixset.cpp
  gdal/gcore/mdreader/reader_alos.cpp
  gdal/gcore/mdreader/reader_digital_globe.cpp
  gdal/gcore/mdreader/reader_eros.cpp
  gdal/gcore/mdreader/reader_geo_eye.cpp
  gdal/gcore/mdreader/reader_kompsat.cpp
  gdal/gcore/mdreader/reader_landsat.cpp
  gdal/gcore/mdreader/reader_orb_view.cpp
  gdal/gcore/mdreader/reader_pleiades.cpp
  gdal/gcore/mdreader/reader_rapid_eye.cpp
  gdal/gcore/mdreader/reader_rdk1.cpp
  gdal/gcore/mdreader/reader_spot.cpp

  # gdal/port/vsipreload.cpp
  # gdal/port/cpl_vsil_simple.cpp
  gdal/port/cpl_compressor.cpp
  gdal/port/xmlreformat.cpp
  gdal/port/cplstringlist.cpp
  gdal/port/cplstring.cpp
  gdal/port/cplkeywordparser.cpp
  gdal/port/cplgetsymbol.cpp
  gdal/port/cpl_xml_validate.cpp
  gdal/port/cpl_worker_thread_pool.cpp
  gdal/port/cpl_vsisimple.cpp
  gdal/port/cpl_vsil_win32.cpp
  gdal/port/cpl_vsil_webhdfs.cpp
  gdal/port/cpl_vsil_unix_stdio_64.cpp
  gdal/port/cpl_vsil_tar.cpp
  gdal/port/cpl_vsil_swift.cpp
  gdal/port/cpl_vsil_subfile.cpp
  gdal/port/cpl_vsil_stdout.cpp
  gdal/port/cpl_vsil_stdin.cpp
  gdal/port/cpl_vsil_sparsefile.cpp
  gdal/port/cpl_vsil_s3.cpp
  gdal/port/cpl_vsil_plugin.cpp
  gdal/port/cpl_vsil_oss.cpp
  gdal/port/cpl_vsil_hdfs.cpp
  gdal/port/cpl_vsil_gzip.cpp
  gdal/port/cpl_vsil_gs.cpp
  gdal/port/cpl_vsil_curl_streaming.cpp
  gdal/port/cpl_vsil_curl.cpp
  gdal/port/cpl_vsil_crypt.cpp
  gdal/port/cpl_vsil_cache.cpp
  gdal/port/cpl_vsil_buffered_reader.cpp
  gdal/port/cpl_vsil_az.cpp
  gdal/port/cpl_vsil_abstract_archive.cpp
  gdal/port/cpl_vsil.cpp
  gdal/port/cpl_vsi_mem.cpp
  gdal/port/cpl_vsi_error.cpp
  gdal/port/cpl_virtualmem.cpp
  gdal/port/cpl_userfaultfd.cpp
  gdal/port/cpl_time.cpp
  gdal/port/cpl_swift.cpp
  gdal/port/cpl_strtod.cpp
  gdal/port/cpl_string.cpp
  gdal/port/cpl_spawn.cpp
  gdal/port/cpl_sha256.cpp
  gdal/port/cpl_sha1.cpp
  gdal/port/cpl_recode_stub.cpp
  gdal/port/cpl_recode_iconv.cpp
  gdal/port/cpl_recode.cpp
  gdal/port/cpl_quad_tree.cpp
  gdal/port/cpl_progress.cpp
  gdal/port/cpl_path.cpp
  #gdal/port/cpl_odbc.cpp
  gdal/port/cpl_multiproc.cpp
  gdal/port/cpl_minizip_zip.cpp
  gdal/port/cpl_minizip_unzip.cpp
  gdal/port/cpl_minizip_ioapi.cpp
  gdal/port/cpl_minixml.cpp
  gdal/port/cpl_md5.cpp
  gdal/port/cpl_list.cpp
  gdal/port/cpl_json_streaming_parser.cpp
  gdal/port/cpl_json.cpp
  gdal/port/cpl_http.cpp
  gdal/port/cpl_hash_set.cpp
  gdal/port/cpl_google_oauth2.cpp
  gdal/port/cpl_google_cloud.cpp
  gdal/port/cpl_getexecpath.cpp
  gdal/port/cpl_findfile.cpp
  gdal/port/cpl_error.cpp
  gdal/port/cpl_csv.cpp
  gdal/port/cpl_cpu_features.cpp
  gdal/port/cpl_conv.cpp
  gdal/port/cpl_base64.cpp
  gdal/port/cpl_azure.cpp
  gdal/port/cpl_aws.cpp
  gdal/port/cpl_atomic_ops.cpp
  gdal/port/cpl_alibaba_oss.cpp
  gdal/port/cpl_vax.cpp
)

IF(WIN32)
LIST(APPEND gdal_SOURCE_FILES gdal/port/cpl_odbc.cpp)
ENDIF()
  
ADD_LIBRARY(gdal_c STATIC
      ${gdal_SOURCE_FILES}
      ${gdal_frmts_gpkg}
    )
SET_TARGET_PROPERTIES(gdal_c PROPERTIES
    ARCHIVE_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
    ARCHIVE_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
SET_TARGET_PROPERTIES(gdal_c PROPERTIES 
    LIBRARY_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
    LIBRARY_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
SET_TARGET_PROPERTIES(gdal_c PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
    RUNTIME_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})

LIST(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/formats")

SET(FMTS_DEPS
  gtiff
  hfa
  aigrid
  aaigrid
  ceos
  ceos2
  iso8211
  xpm
  sdts
  raw
  dted
  mem
  jdem
  envisat
  elas
  fit
  vrt
  usgsdem
  l1b
  nitf
  bmp
  pcidsk
  airsar
  rs2
  ilwis
  rmf
  leveller
  sgi
  srtmhgt
  idrisi
  gsg
  ingr
  ers
  jaxapalsar
  dimap
  gff
  cosar
  pds
  adrg
  coasp
  tsx
  terragen
  blx
  til
  r
  northwood
  saga
  xyz
  hf2
  kmlsuperoverlay
  ctg
  zmap
  ngsgeoid
  iris
  map
  jpeg
  png
  webp
  mbtiles
  ogr_shape
  ogr_vrt
  ogr_avc
  ogr_geojson
  ogr_mem
  ogr_mitab
  ogr_kml
  ogr_gpx
  ogr_dxf
  ogr_csv
  ogr_edigeo
  ogr_geoconcept
  ogr_georss
  ogr_gml
  ogr_gmt
  ogr_gpsbabel
  ogr_gtm
  ogr_idrisi
  ogr_dgn
  ogr_openfilegdb
  ogr_osm
  ogr_pds
  ogr_pgdump
  ogr_rec
  ogr_s57
  ogr_sdts
  ogr_svg
  ogr_sxf
  ogr_ntf
  ogr_wasp
  ogr_sqlite
  ogr_gpkg
  ogr_mvt
)

add_definitions(
  -DFRMT_gtiff=1
  -DFRMT_hfa=1
  -DFRMT_aigrid=1
  -DFRMT_aaigrid=1
  -DFRMT_ceos=1
  -DFRMT_ceos2=1
  -DFRMT_iso8211=1
  -DFRMT_xpm=1
  -DFRMT_sdts=1
  -DFRMT_raw=1
  -DFRMT_dted=1
  -DFRMT_mem=1
  -DFRMT_jdem=1
  -DFRMT_envisat=1
  -DFRMT_elas=1
  -DFRMT_fit=1
  -DFRMT_vrt=1
  -DFRMT_usgsdem=1
  -DFRMT_l1b=1
  -DFRMT_nitf=1
  -DFRMT_bmp=1
  -DFRMT_pcidsk=1
  -DFRMT_airsar=1
  -DFRMT_rs2=1
  -DFRMT_ilwis=1
  -DFRMT_rmf=1
  -DFRMT_leveller=1
  -DFRMT_sgi=1
  -DFRMT_srtmhgt=1
  -DFRMT_idrisi=1
  -DFRMT_gsg=1
  -DFRMT_ingr=1
  -DFRMT_ers=1
  -DFRMT_jaxapalsar=1
  -DFRMT_dimap=1
  -DFRMT_gff=1
  -DFRMT_cosar=1
  -DFRMT_pds=1
  -DFRMT_adrg=1
  -DFRMT_coasp=1
  -DFRMT_tsx=1
  -DFRMT_terragen=1
  -DFRMT_blx=1
  -DFRMT_til=1
  -DFRMT_r=1
  -DFRMT_northwood=1
  -DFRMT_saga=1
  -DFRMT_xyz=1
  -DFRMT_hf2=1
  -DFRMT_kmlsuperoverlay=1
  -DFRMT_ctg=1
  -DFRMT_zmap=1
  -DFRMT_ngsgeoid=1
  -DFRMT_iris=1
  -DFRMT_map=1
  -DFRMT_zlib=1
  -DFRMT_jpeg=1
  -DFRMT_webp=1
  -DFRMT_png=1
  -DFRMT_mbtiles=1
  -DSHAPE_ENABLED=1
  -DVRT_ENABLED=1
  -DAVC_ENABLED=1
  -DGEOJSON_ENABLED=1
  -DMEM_ENABLED=1
  -DTAB_ENABLED=1
  -DMITAB_ENABLED=1
  -DKML_ENABLED=1
  -DGPX_ENABLED=1
  -DDXF_ENABLED=1
  -DCSV_ENABLED=1
  -DEDIGEO_ENABLED=1
  -DGEOCONCEPT_ENABLED=1
  -DGEORSS_ENABLED=1
  -DGML_ENABLED=1
  -DGMT_ENABLED=1
  -DGPSBABEL_ENABLED=1
  -DGTM_ENABLED=1
  -DIDRISI_ENABLED=1
  -DDGN_ENABLED=1
  -DOPENFILEGDB_ENABLED=1
  -DPDS_ENABLED=1
  -DPGDUMP_ENABLED=1
  -DREC_ENABLED=1
  -DS57_ENABLED=1
  -DSDTS_ENABLED=1
  -DOSM_ENABLED=1
  -DSVG_ENABLED=1
  -DSXF_ENABLED=1
  -DNTF_ENABLED=1
  -DSQLITE_ENABLE=1
  -DWASP_ENABLED=1
  -DSQLITE_ENABLED=1
  -DGPKG_ENABLED=1
	-DMVT_ENABLED=1
  -DLIBICONV_PLUG
)

FOREACH(format ${FMTS_DEPS})
  INCLUDE(${format})
ENDFOREACH()

IF(FLYWAVE_ENABLE_SOLUTION_FOLDERS)
  SET_TARGET_PROPERTIES(gdal_c PROPERTIES FOLDER external)

  FOREACH(format ${FMTS_DEPS})
    SET_TARGET_PROPERTIES(gdal_${format}_frmt PROPERTIES FOLDER external/gdal_format)
    SET_TARGET_PROPERTIES(gdal_${format}_frmt PROPERTIES
    ARCHIVE_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
    ARCHIVE_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
    SET_TARGET_PROPERTIES(gdal_${format}_frmt PROPERTIES 
  LIBRARY_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
  LIBRARY_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
    SET_TARGET_PROPERTIES(gdal_${format}_frmt PROPERTIES
  RUNTIME_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
  RUNTIME_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
  ENDFOREACH()
ENDIF(FLYWAVE_ENABLE_SOLUTION_FOLDERS)
