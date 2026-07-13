CMAKE_MINIMUM_REQUIRED(VERSION 3.1.3)
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
MAKE_DIRECTORY("${INSTALL_DIR}/lib")
MAKE_DIRECTORY("${INSTALL_DIR}/include")

#===============================================================================
# GDAL ecosystem — all via add_subdirectory (flywave-gdal convention)
#===============================================================================
ADD_SUBDIRECTORY("${CMAKE_CURRENT_SOURCE_DIR}/external/zlib")
ADD_SUBDIRECTORY("${CMAKE_CURRENT_SOURCE_DIR}/external/libpng")
ADD_SUBDIRECTORY("${CMAKE_CURRENT_SOURCE_DIR}/external/libjpeg")
ADD_SUBDIRECTORY("${CMAKE_CURRENT_SOURCE_DIR}/external/libexpat")
ADD_SUBDIRECTORY("${CMAKE_CURRENT_SOURCE_DIR}/external/libiconv")
ADD_SUBDIRECTORY("${CMAKE_CURRENT_SOURCE_DIR}/external/libsqlite3")
ADD_SUBDIRECTORY("${CMAKE_CURRENT_SOURCE_DIR}/external/libgeos")
ADD_SUBDIRECTORY("${CMAKE_CURRENT_SOURCE_DIR}/external/libproj")
ADD_SUBDIRECTORY("${CMAKE_CURRENT_SOURCE_DIR}/external/webp")
ADD_SUBDIRECTORY("${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal")

#===============================================================================
# HTDP (Fortran)
#===============================================================================
ADD_SUBDIRECTORY("${CMAKE_CURRENT_SOURCE_DIR}/external/HTDP")

#===============================================================================
# Install .a files to libs/{platform}/
#===============================================================================
# GDAL libraries
FILE(GLOB GDAL_ARCHIVES
  "${CMAKE_CURRENT_BINARY_DIR}/external/libgdal/libgdal_*.a"
  "${CMAKE_CURRENT_BINARY_DIR}/external/libgdal/libgdal_c.a"
)
INSTALL(FILES ${GDAL_ARCHIVES} DESTINATION "${INSTALL_DIR}")

# GDAL headers
INSTALL(DIRECTORY
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/port/"
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/gcore/"
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/alg/"
  "${CMAKE_CURRENT_SOURCE_DIR}/external/libgdal/gdal/ogr/"
  DESTINATION "${INSTALL_DIR}/include"
  FILES_MATCHING PATTERN "*.h"
)

# Other libraries
FOREACH(lib zlib png jpeg expat iconv sqlite3 geos proj webp)
  INSTALL(FILES "${CMAKE_CURRENT_BINARY_DIR}/external/lib${lib}/lib${lib}.a"
    DESTINATION "${INSTALL_DIR}")
ENDFOREACH()

# HTDP
INSTALL(FILES "${CMAKE_CURRENT_BINARY_DIR}/external/HTDP/libhtdp.a"
  DESTINATION "${INSTALL_DIR}")
INSTALL(DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}/external/HTDP/capi/"
  DESTINATION "${FLYWAVE_BASE_INSTALL_DIR}"
  FILES_MATCHING PATTERN "*.h")
