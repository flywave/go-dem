# go-dem GMT build — wrapper around GMT's own cmake
cmake_minimum_required(VERSION 3.10)
project(gmt C)

# In the original GMT build, CMAKE_SOURCE_DIR equals GMT_SOURCE_DIR.
# Since we are add_subdirectory'd, save the original and override locally.
set(_saved_CMAKE_SOURCE_DIR "${CMAKE_SOURCE_DIR}")
set(CMAKE_SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}")
set(GMT_SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}")

# Set up GMT's cmake module path
list(INSERT CMAKE_MODULE_PATH 0
  "${CMAKE_CURRENT_SOURCE_DIR}/cmake/modules"
)
include(GmtHelperMacros)
include(CreateDebugSym)

# Load GMT's full cmake configuration (uses CMAKE_SOURCE_DIR internally)
include("${CMAKE_CURRENT_SOURCE_DIR}/cmake/ConfigDefault.cmake")
include("${CMAKE_CURRENT_SOURCE_DIR}/cmake/modules/ConfigCMake.cmake")
include("${CMAKE_CURRENT_SOURCE_DIR}/cmake/modules/ConfigureChecks.cmake")

# Manually set variables that FindGSHHG.cmake would normally provide
set(GSHHG_MIN_REQUIRED_VERSION_MAJOR 2 CACHE INTERNAL "")
set(GSHHG_MIN_REQUIRED_VERSION_MINOR 2 CACHE INTERNAL "")
set(GSHHG_MIN_REQUIRED_VERSION_PATCH 0 CACHE INTERNAL "")

# Generate config.h from config.h.in (uses variables set by ConfigDefault+ConfigureChecks)
configure_file("${CMAKE_CURRENT_SOURCE_DIR}/src/config.h.in" "${CMAKE_CURRENT_BINARY_DIR}/src/config.h")

# Restore original CMAKE_SOURCE_DIR for parent project
set(CMAKE_SOURCE_DIR "${_saved_CMAKE_SOURCE_DIR}")

# Set install paths (needed by src/CMakeLists.txt)
set(CMAKE_INSTALL_PREFIX "${CMAKE_CURRENT_SOURCE_DIR}/../../libs/darwin_arm")
set(GMT_BINDIR "${CMAKE_INSTALL_PREFIX}")
set(GMT_LIBDIR "${CMAKE_INSTALL_PREFIX}")
set(GMT_INCLUDEDIR "${CMAKE_INSTALL_PREFIX}/include")
set(GMT_DATADIR "${CMAKE_INSTALL_PREFIX}/share/gmt")

# Dummy check target for supplements
add_custom_target(check)

# Disable supplement plugins (they need transitive deps from GDAL static lib)
set(BUILD_SUPPLEMENTS OFF)

# Build
add_subdirectory(src)
