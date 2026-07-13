# go-dem GMT build target
# Uses GMT's own cmake externally for reliable build
cmake_minimum_required(VERSION 3.10)
project(gmt C)

get_filename_component(GODEM "${CMAKE_CURRENT_SOURCE_DIR}/../.." ABSOLUTE)
if(APPLE)
  if(CMAKE_SYSTEM_PROCESSOR MATCHES "^arm")
    set(PLAT "darwin_arm")
  else()
    set(PLAT "darwin")
  endif()
elseif(WIN32)
  set(PLAT "windows")
else()
  if(CMAKE_SYSTEM_PROCESSOR MATCHES "aarch64")
    set(PLAT "linux_arm")
  else()
    set(PLAT "linux")
  endif()
endif()
set(LIBS "${GODEM}/libs/${PLAT}")

set(GMT_SRC "${CMAKE_CURRENT_SOURCE_DIR}")
set(GMT_BLD "${CMAKE_CURRENT_BINARY_DIR}/gmt_build")

add_custom_target(gmt_build ALL
  COMMAND ${CMAKE_COMMAND} -E make_directory "${GMT_BLD}"
  COMMAND ${CMAKE_COMMAND} -E env
    CC=${CMAKE_C_COMPILER}
    CXX=${CMAKE_CXX_COMPILER}
    ${CMAKE_COMMAND}
    -DCMAKE_BUILD_TYPE=${CMAKE_BUILD_TYPE}
    -DCMAKE_INSTALL_PREFIX=${LIBS}
    -DBUILD_SHARED_LIBS=OFF
    -DBUILD_TESTING=OFF
    -DNETCDF_INCLUDE_DIR=/opt/homebrew/include
    -DNETCDF_LIBRARY=/opt/homebrew/lib/libnetcdf.dylib
    -DHDF5_INCLUDE_DIR=${LIBS}/include
    -DHDF5_LIBRARY=${LIBS}/libhdf5.a
    -DHDF5_HL_LIBRARY=${LIBS}/libhdf5_hl.a
    -DGDAL_INCLUDE_DIR="${GODEM}/external/libgdal/gdal;${GODEM}/external/libgdal/gdal/gcore;${GODEM}/external/libgdal/gdal/port;${GODEM}/external/libgdal/gdal/alg;${GODEM}/../../flywave-gdal/libs"
    -DGDAL_LIBRARY=${LIBS}/libgdal_c.a
    -DPROJ_INCLUDE_DIR="${GODEM}/external/libproj/proj/src;${GODEM}/external/libproj/proj/include"
    -DPROJ_LIBRARY=${LIBS}/libproj.a
    -DGEOS_INCLUDE_DIR=${GODEM}/external/libgeos/geos/include
    -DGEOS_LIBRARY=${LIBS}/libgeos.a
    -DSQLITE3_INCLUDE_DIR=${GODEM}/external/libsqlite3
    -DSQLITE3_LIBRARY=${LIBS}/libsqlite3.a
    -DZLIB_LIBRARY=${LIBS}/libzlib.a
    -DZLIB_INCLUDE_DIR=${GODEM}/external/zlib
    -DGMT_EXCLUDE_GEOS=ON
    -DGMT_EXCLUDE_PCRE=ON
    -DGMT_EXCLUDE_PCRE2=ON
    -DGMT_EXCLUDE_FFTW3=ON
    -DCMAKE_C_FLAGS=-DGMT_NO_CURL
    "${GMT_SRC}"
  COMMAND ${CMAKE_COMMAND} --build "${GMT_BLD}" --target gmtlib --target pslib -j4
  COMMAND ${CMAKE_COMMAND} -E copy "${GMT_BLD}/src/libgmt.a" "${LIBS}/libgmt.a"
  COMMAND ${CMAKE_COMMAND} -E copy "${GMT_BLD}/src/libpostscriptlight.a" "${LIBS}/libpsl.a"
  COMMENT "Building GMT (gmtlib + pslib)..."
  BYPRODUCTS "${LIBS}/libgmt.a" "${LIBS}/libpsl.a"
)
