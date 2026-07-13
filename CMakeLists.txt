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
# HDF5 (simplified CMake — direct source build)
#===============================================================================
ADD_SUBDIRECTORY("external/hdf5")

#===============================================================================
# netCDF (links against hdf5 targets directly)
#===============================================================================
ADD_SUBDIRECTORY("external/netcdf")

#===============================================================================
# GMT (set variables so find_package finds our locally-built libs)
#===============================================================================
SET(GDAL_INCLUDE_DIR
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal"
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/gcore"
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/port"
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/alg"
  "${CMAKE_CURRENT_SOURCE_DIR}/../../flywave-gdal/libs"
)
SET(GDAL_LIBRARY "${INSTALL_DIR}/libgdal_c.a")
SET(PROJ_INCLUDE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/libproj/proj/src")
SET(PROJ_LIBRARY "${INSTALL_DIR}/libproj.a")
SET(HDF5_INCLUDE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/hdf5/src")
SET(HDF5_LIBRARY "hdf5")
SET(HDF5_HL_LIBRARY "hdf5_hl")
SET(NETCDF_INCLUDE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/netcdf/include")
SET(NETCDF_LIBRARY "netcdf")
SET(NETCDF_LIBRARIES "netcdf")

#===============================================================================
# GMT — build as part of unified cmake
#===============================================================================
SET(GMT_EXCLUDE_GEOS ON CACHE BOOL "" FORCE)
SET(GMT_EXCLUDE_PCRE ON CACHE BOOL "" FORCE)
SET(GMT_EXCLUDE_PCRE2 ON CACHE BOOL "" FORCE)
SET(GMT_EXCLUDE_FFTW3 ON CACHE BOOL "" FORCE)
ADD_DEFINITIONS(-DGMT_NO_CURL)
SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -DGMT_NO_CURL")
ADD_SUBDIRECTORY("external/gmt")
