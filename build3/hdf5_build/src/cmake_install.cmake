# Install script for directory: /Users/xuning/Work/go-dem/external/hdf5/src

# Set the install prefix
if(NOT DEFINED CMAKE_INSTALL_PREFIX)
  set(CMAKE_INSTALL_PREFIX "/Users/xuning/Work/go-dem/libs/darwin_arm")
endif()
string(REGEX REPLACE "/$" "" CMAKE_INSTALL_PREFIX "${CMAKE_INSTALL_PREFIX}")

# Set the install configuration name.
if(NOT DEFINED CMAKE_INSTALL_CONFIG_NAME)
  if(BUILD_TYPE)
    string(REGEX REPLACE "^[^A-Za-z0-9_]+" ""
           CMAKE_INSTALL_CONFIG_NAME "${BUILD_TYPE}")
  else()
    set(CMAKE_INSTALL_CONFIG_NAME "Release")
  endif()
  message(STATUS "Install configuration: \"${CMAKE_INSTALL_CONFIG_NAME}\"")
endif()

# Set the component getting installed.
if(NOT CMAKE_INSTALL_COMPONENT)
  if(COMPONENT)
    message(STATUS "Install component: \"${COMPONENT}\"")
    set(CMAKE_INSTALL_COMPONENT "${COMPONENT}")
  else()
    set(CMAKE_INSTALL_COMPONENT)
  endif()
endif()

# Is this installation the result of a crosscompile?
if(NOT DEFINED CMAKE_CROSSCOMPILING)
  set(CMAKE_CROSSCOMPILING "FALSE")
endif()

# Set path to fallback-tool for dependency-resolution.
if(NOT DEFINED CMAKE_OBJDUMP)
  set(CMAKE_OBJDUMP "/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/objdump")
endif()

if(NOT CMAKE_INSTALL_LOCAL_ONLY)
  # Include the install script for the subdirectory.
  include("/Users/xuning/Work/go-dem/build3/hdf5_build/src/H5FDsubfiling/cmake_install.cmake")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "headers" OR NOT CMAKE_INSTALL_COMPONENT)
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include" TYPE FILE FILES
    "/Users/xuning/Work/go-dem/external/hdf5/src/hdf5.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5api_adpt.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5public.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Apublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5ACpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Cpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Dpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Epubgen.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Epublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5ESdevelop.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5ESpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Fpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDcore.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDdevelop.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDdirect.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDfamily.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDhdfs.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDlog.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDmirror.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDmpi.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDmpio.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDmulti.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDonion.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDros3.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDsec2.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDsplitter.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDstdio.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDwindows.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDsubfiling/H5FDsubfiling.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5FDsubfiling/H5FDioc.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Gpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Idevelop.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Ipublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Ldevelop.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Lpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Mpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5MMpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Opublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Ppublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5PLextern.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5PLpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Rpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Spublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Tdevelop.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Tpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5TSdevelop.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5VLconnector.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5VLconnector_passthru.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5VLnative.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5VLpassthru.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5VLpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Zdevelop.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Zpublic.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5Epubgen.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5version.h"
    "/Users/xuning/Work/go-dem/external/hdf5/src/H5overflow.h"
    "/Users/xuning/Work/go-dem/build3/hdf5_build/src/H5pubconf.h"
    )
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "libraries" OR NOT CMAKE_INSTALL_COMPONENT)
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/lib" TYPE STATIC_LIBRARY FILES "/Users/xuning/Work/go-dem/build3/hdf5_build/bin/libhdf5.a")
  if(EXISTS "$ENV{DESTDIR}${CMAKE_INSTALL_PREFIX}/lib/libhdf5.a" AND
     NOT IS_SYMLINK "$ENV{DESTDIR}${CMAKE_INSTALL_PREFIX}/lib/libhdf5.a")
    execute_process(COMMAND "/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/ranlib" "$ENV{DESTDIR}${CMAKE_INSTALL_PREFIX}/lib/libhdf5.a")
  endif()
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "libraries" OR NOT CMAKE_INSTALL_COMPONENT)
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/lib/pkgconfig" TYPE FILE FILES "/Users/xuning/Work/go-dem/build3/hdf5_build/CMakeFiles/hdf5.pc")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "libraries" OR NOT CMAKE_INSTALL_COMPONENT)
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/bin" TYPE FILE PERMISSIONS OWNER_READ OWNER_WRITE OWNER_EXECUTE GROUP_READ GROUP_EXECUTE WORLD_READ WORLD_EXECUTE FILES "/Users/xuning/Work/go-dem/build3/hdf5_build/CMakeFiles/h5cc")
endif()

string(REPLACE ";" "\n" CMAKE_INSTALL_MANIFEST_CONTENT
       "${CMAKE_INSTALL_MANIFEST_FILES}")
if(CMAKE_INSTALL_LOCAL_ONLY)
  file(WRITE "/Users/xuning/Work/go-dem/build3/hdf5_build/src/install_local_manifest.txt"
     "${CMAKE_INSTALL_MANIFEST_CONTENT}")
endif()
