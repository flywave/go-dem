# netCDF — build externally, then use as imported target
cmake_minimum_required(VERSION 3.10)
project(netcdf C)

# netCDF must be built separately with its own cmake:
#   cd external/netcdf && mkdir _build && cd _build
#   cmake .. -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=../../libs/darwin_arm [options]
#   cmake --build . -j4
#   cp liblib/libnetcdf.a ../../libs/darwin_arm/
#   cp ../include/netcdf*.h ../../libs/darwin_arm/include/

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

add_library(netcdf STATIC IMPORTED GLOBAL)
set_target_properties(netcdf PROPERTIES
  IMPORTED_LOCATION "${GODEM}/libs/${PLAT}/libnetcdf.a"
  INTERFACE_INCLUDE_DIRECTORIES "${GODEM}/libs/${PLAT}/include"
  INTERFACE_LINK_LIBRARIES "hdf5;hdf5_hl"
)

# Also provide NETCDF_LIBRARIES variable for GMT
set(NETCDF_LIBRARIES "netcdf" PARENT_SCOPE)
