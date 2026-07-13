# go-dem GMT build — wrapper around GMT's own cmake
cmake_minimum_required(VERSION 3.10)
project(gmt C)

# Set up GMT's cmake module path
list(INSERT CMAKE_MODULE_PATH 0
  "${CMAKE_CURRENT_SOURCE_DIR}/cmake/modules"
)
include(GmtHelperMacros)
include(CreateDebugSym)

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
