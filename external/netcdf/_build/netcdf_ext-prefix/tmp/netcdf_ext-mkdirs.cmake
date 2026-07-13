# Distributed under the OSI-approved BSD 3-Clause License.  See accompanying
# file Copyright.txt or https://cmake.org/licensing for details.

cmake_minimum_required(VERSION ${CMAKE_VERSION}) # this file comes with cmake

# If CMAKE_DISABLE_SOURCE_CHANGES is set to true and the source directory is an
# existing directory in our source tree, calling file(MAKE_DIRECTORY) on it
# would cause a fatal error, even though it would be a no-op.
if(NOT EXISTS "/Users/xuning/Work/go-dem/external/netcdf")
  file(MAKE_DIRECTORY "/Users/xuning/Work/go-dem/external/netcdf")
endif()
file(MAKE_DIRECTORY
  "/Users/xuning/Work/go-dem/external/netcdf/_build/build"
  "/Users/xuning/Work/go-dem/libs/darwin_arm"
  "/Users/xuning/Work/go-dem/external/netcdf/_build/netcdf_ext-prefix/tmp"
  "/Users/xuning/Work/go-dem/external/netcdf/_build/netcdf_ext-prefix/src/netcdf_ext-stamp"
  "/Users/xuning/Work/go-dem/external/netcdf/_build/netcdf_ext-prefix/src"
  "/Users/xuning/Work/go-dem/external/netcdf/_build/netcdf_ext-prefix/src/netcdf_ext-stamp"
)

set(configSubDirs )
foreach(subDir IN LISTS configSubDirs)
    file(MAKE_DIRECTORY "/Users/xuning/Work/go-dem/external/netcdf/_build/netcdf_ext-prefix/src/netcdf_ext-stamp/${subDir}")
endforeach()
if(cfgdir)
  file(MAKE_DIRECTORY "/Users/xuning/Work/go-dem/external/netcdf/_build/netcdf_ext-prefix/src/netcdf_ext-stamp${cfgdir}") # cfgdir has leading slash
endif()
