# HDF5 — use prebuilt libraries from libs/
cmake_minimum_required(VERSION 3.10)
project(hdf5 C)

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

# Create imported targets pointing to prebuilt .a files
add_library(hdf5 STATIC IMPORTED GLOBAL)
set_target_properties(hdf5 PROPERTIES
  IMPORTED_LOCATION "${GODEM}/libs/${PLAT}/libhdf5.a"
  INTERFACE_INCLUDE_DIRECTORIES "${CMAKE_CURRENT_SOURCE_DIR}/src"
)

add_library(hdf5_hl STATIC IMPORTED GLOBAL)
set_target_properties(hdf5_hl PROPERTIES
  IMPORTED_LOCATION "${GODEM}/libs/${PLAT}/libhdf5_hl.a"
  INTERFACE_INCLUDE_DIRECTORIES "${CMAKE_CURRENT_SOURCE_DIR}/hl/src"
  INTERFACE_LINK_LIBRARIES hdf5
)
