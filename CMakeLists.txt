CMAKE_MINIMUM_REQUIRED(VERSION 3.10)
PROJECT(go-dem C CXX Fortran)

IF(NOT CMAKE_BUILD_TYPE)
  SET(CMAKE_BUILD_TYPE Release CACHE STRING "" FORCE)
ENDIF()

SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fPIC -Wall")
SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -fPIC -Wall")
IF(APPLE)
  SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=gnu++14")
ELSE()
  SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++14")
ENDIF()

# Output directory
SET(FLYWAVE_BASE_INSTALL_DIR "${CMAKE_CURRENT_SOURCE_DIR}/libs")
IF(APPLE)
  IF(${CMAKE_SYSTEM_PROCESSOR} MATCHES "^arm")
    SET(PLATFORM_DIR "darwin_arm")
  ELSE()
    SET(PLATFORM_DIR "darwin")
  ENDIF()
ELSEIF(WIN32)
  SET(PLATFORM_DIR "windows")
ELSE()
  IF(CMAKE_SYSTEM_PROCESSOR MATCHES "aarch64")
    SET(PLATFORM_DIR "linux_arm")
  ELSE()
    SET(PLATFORM_DIR "linux")
  ENDIF()
ENDIF()
SET(INSTALL_DIR "${FLYWAVE_BASE_INSTALL_DIR}/${PLATFORM_DIR}")

#===============================================================================
# GDAL ecosystem
#===============================================================================
ADD_SUBDIRECTORY("external/zlib")
# Create ZLIB::ZLIB alias so HDF5/netCDF cmake configs can find it
if(NOT TARGET ZLIB::ZLIB)
  add_library(ZLIB::ZLIB ALIAS zlib)
endif()
ADD_SUBDIRECTORY("external/libpng")
ADD_SUBDIRECTORY("external/libjpeg")
ADD_SUBDIRECTORY("external/libexpat")
ADD_SUBDIRECTORY("external/libiconv")
ADD_SUBDIRECTORY("external/libsqlite3")
ADD_SUBDIRECTORY("external/libgeos")
ADD_SUBDIRECTORY("external/libproj")
ADD_SUBDIRECTORY("external/webp")
ADD_SUBDIRECTORY("external/libgdal")

#===============================================================================
# HTDP
#===============================================================================
ADD_SUBDIRECTORY("external/HTDP")

#===============================================================================
# GMT (set variables so find_package finds our locally-built libs)
#===============================================================================
IF(APPLE)
  SET(GDAL_ARCH_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/arch/osx")
ELSEIF(WIN32)
  SET(GDAL_ARCH_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/arch/win")
ELSE()
  SET(GDAL_ARCH_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/arch/unix")
ENDIF()
SET(GDAL_INCLUDE_DIR
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal"
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/gcore"
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/port"
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/alg"
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/ogr"
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/apps"
  "${GDAL_ARCH_DIR}"
  "${CMAKE_CURRENT_SOURCE_DIR}/../../flywave-gdal/libs"
)
SET(GDAL_LIBRARY "${INSTALL_DIR}/libgdal_c.a")
SET(PROJ_INCLUDE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/libproj/proj/src")
SET(PROJ_LIBRARY "${INSTALL_DIR}/libproj.a")
# Pre-populate GMT_OPTIONAL_LIBRARIES with PROJ (transitive dep of GDAL static lib)
SET(GMT_OPTIONAL_LIBRARIES "${PROJ_LIBRARY};${INSTALL_DIR}/libsqlite3.a")
SET(HDF5_INCLUDE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/hdf5/src")
SET(HDF5_LIBRARY "hdf5")
SET(HDF5_HL_LIBRARY "hdf5_hl")

#===============================================================================
# netCDF — imported target (pre-built lib in libs/)
#===============================================================================
# HDF5 imported targets (pre-built libs in libs/)
add_library(HDF5::HDF5 STATIC IMPORTED GLOBAL)
set_target_properties(HDF5::HDF5 PROPERTIES
  IMPORTED_LOCATION "${INSTALL_DIR}/lib/libhdf5.a"
  INTERFACE_INCLUDE_DIRECTORIES "${CMAKE_CURRENT_SOURCE_DIR}/external/hdf5/src;${INSTALL_DIR}/include/hdf5"
  INTERFACE_LINK_LIBRARIES "ZLIB::ZLIB"
)
add_library(hdf5::hdf5_hl STATIC IMPORTED GLOBAL)
set_target_properties(hdf5::hdf5_hl PROPERTIES
  IMPORTED_LOCATION "${INSTALL_DIR}/lib/libhdf5_hl.a"
  INTERFACE_INCLUDE_DIRECTORIES "${CMAKE_CURRENT_SOURCE_DIR}/external/hdf5/src;${INSTALL_DIR}/include/hdf5;${CMAKE_CURRENT_SOURCE_DIR}/external/hdf5/hl/src"
  INTERFACE_LINK_LIBRARIES "HDF5::HDF5"
)

# netCDF imported target (pre-built lib in libs/)
add_library(netcdf STATIC IMPORTED GLOBAL)
set_target_properties(netcdf PROPERTIES
  IMPORTED_LOCATION "${INSTALL_DIR}/lib/libnetcdf.a"
  INTERFACE_INCLUDE_DIRECTORIES "${CMAKE_CURRENT_SOURCE_DIR}/external/netcdf/include"
  INTERFACE_LINK_LIBRARIES "HDF5::HDF5;hdf5::hdf5_hl"
)
set(NETCDF_INCLUDE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/netcdf/include")
set(NETCDF_LIBRARY "netcdf")
set(NETCDF_LIBRARIES "netcdf")

#===============================================================================
# GMT — build as part of unified cmake
#===============================================================================
SET(GMT_EXCLUDE_GEOS ON CACHE BOOL "" FORCE)
SET(GMT_EXCLUDE_PCRE ON CACHE BOOL "" FORCE)
SET(GMT_EXCLUDE_PCRE2 ON CACHE BOOL "" FORCE)
SET(GMT_EXCLUDE_FFTW3 ON CACHE BOOL "" FORCE)
ADD_DEFINITIONS(-DGMT_NO_CURL)
SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -DGMT_NO_CURL")
# Build GMT as static library to avoid transitive-dep resolution issues
SET(BUILD_SHARED_LIBS OFF CACHE BOOL "" FORCE)
# Provide all GDAL transitive deps for GMT (needed for linking static libgdal_c.a)
# Collect all GDAL format libraries
FILE(GLOB GDAL_FMT_LIBS "${INSTALL_DIR}/libgdal_*.a")
SET(GMT_OPTIONAL_LIBRARIES
  "${INSTALL_DIR}/libproj.a"
  "${INSTALL_DIR}/libsqlite3.a"
  "${INSTALL_DIR}/libgeos.a"
  "${INSTALL_DIR}/libiconv.a"
  "${INSTALL_DIR}/libexpat.a"
  "${INSTALL_DIR}/libzlib.a"
  "${INSTALL_DIR}/libpng.a"
  "${INSTALL_DIR}/libjpeg.a"
  "${INSTALL_DIR}/libwebp.a"
  ${GDAL_FMT_LIBS}
  "-lc++"
  "-lm"
)
ADD_SUBDIRECTORY("external/gmt")
SET(BUILD_SHARED_LIBS "${_saved_BUILD_SHARED_LIBS}")
