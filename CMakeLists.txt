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
# HTDP
#===============================================================================
ADD_SUBDIRECTORY("external/HTDP")
