# go-dem GMT build — wrapper around GMT's own cmake
cmake_minimum_required(VERSION 3.10)
project(gmt C)

# Required by GMT cmake modules and src/CMakeLists.txt
set(GMT_SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}")

# Generate config.h from config.h.in (needed by gmt_config.h)
configure_file("${CMAKE_CURRENT_SOURCE_DIR}/src/config.h.in" "${CMAKE_CURRENT_BINARY_DIR}/src/config.h")

# Set up GMT's cmake module path
list(INSERT CMAKE_MODULE_PATH 0
  "${CMAKE_CURRENT_SOURCE_DIR}/cmake/modules"
)
include(GmtHelperMacros)
include(CreateDebugSym)

# Load GMT's default configuration (version numbers, lib names, etc.)
include("${CMAKE_CURRENT_SOURCE_DIR}/cmake/ConfigDefault.cmake")

# Set install paths (needed by src/CMakeLists.txt)
set(CMAKE_INSTALL_PREFIX "${CMAKE_CURRENT_SOURCE_DIR}/../../libs/darwin_arm")
set(GMT_BINDIR "${CMAKE_INSTALL_PREFIX}")
set(GMT_LIBDIR "${CMAKE_INSTALL_PREFIX}")
set(GMT_INCLUDEDIR "${CMAKE_INSTALL_PREFIX}/include")
set(GMT_DATADIR "${CMAKE_INSTALL_PREFIX}/share/gmt")

# Dummy check target for supplements
add_custom_target(check)

# Build
add_subdirectory(src)
