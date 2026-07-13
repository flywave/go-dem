# Install script for directory: /Users/xuning/Work/go-dem/external/gmt/src

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

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/libgmt.6.7.0.dylib;/Users/xuning/Work/go-dem/libs/darwin_arm/libgmt.6.dylib")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm" TYPE SHARED_LIBRARY FILES
    "/Users/xuning/Work/go-dem/build5/external/gmt/src/libgmt.6.7.0.dylib"
    "/Users/xuning/Work/go-dem/build5/external/gmt/src/libgmt.6.dylib"
    )
  foreach(file
      "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/libgmt.6.7.0.dylib"
      "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/libgmt.6.dylib"
      )
    if(EXISTS "${file}" AND
       NOT IS_SYMLINK "${file}")
      execute_process(COMMAND "/usr/bin/install_name_tool"
        -id "@executable_path/../lib/libgmt.6.dylib"
        -change "@rpath/libpostscriptlight.6.dylib" "@executable_path/../lib/libpostscriptlight.6.dylib"
        "${file}")
      execute_process(COMMAND /usr/bin/install_name_tool
        -delete_rpath "/Users/xuning/Work/go-dem/build5/external/gmt/src"
        -add_rpath "$ORIGIN/../lib"
        "${file}")
      if(CMAKE_INSTALL_DO_STRIP)
        execute_process(COMMAND "/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/strip" -x "${file}")
      endif()
    endif()
  endforeach()
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/libgmt.dylib")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm" TYPE SHARED_LIBRARY FILES "/Users/xuning/Work/go-dem/build5/external/gmt/src/libgmt.dylib")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/libpostscriptlight.6.7.0.dylib;/Users/xuning/Work/go-dem/libs/darwin_arm/libpostscriptlight.6.dylib")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm" TYPE SHARED_LIBRARY FILES
    "/Users/xuning/Work/go-dem/build5/external/gmt/src/libpostscriptlight.6.7.0.dylib"
    "/Users/xuning/Work/go-dem/build5/external/gmt/src/libpostscriptlight.6.dylib"
    )
  foreach(file
      "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/libpostscriptlight.6.7.0.dylib"
      "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/libpostscriptlight.6.dylib"
      )
    if(EXISTS "${file}" AND
       NOT IS_SYMLINK "${file}")
      execute_process(COMMAND "/usr/bin/install_name_tool"
        -id "@executable_path/../lib/libpostscriptlight.6.dylib"
        "${file}")
      execute_process(COMMAND /usr/bin/install_name_tool
        -add_rpath "$ORIGIN/../lib"
        "${file}")
      if(CMAKE_INSTALL_DO_STRIP)
        execute_process(COMMAND "/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/strip" -x "${file}")
      endif()
    endif()
  endforeach()
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/libpostscriptlight.dylib")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm" TYPE SHARED_LIBRARY FILES "/Users/xuning/Work/go-dem/build5/external/gmt/src/libpostscriptlight.dylib")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/gmt")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm" TYPE EXECUTABLE FILES "/Users/xuning/Work/go-dem/build5/external/gmt/src/gmt")
  if(EXISTS "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/gmt" AND
     NOT IS_SYMLINK "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/gmt")
    execute_process(COMMAND "/usr/bin/install_name_tool"
      -change "@rpath/libgmt.6.dylib" "@executable_path/../lib/libgmt.6.dylib"
      -change "@rpath/libpostscriptlight.6.dylib" "@executable_path/../lib/libpostscriptlight.6.dylib"
      "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/gmt")
    execute_process(COMMAND /usr/bin/install_name_tool
      -delete_rpath "/Users/xuning/Work/go-dem/build5/external/gmt/src"
      -add_rpath "$ORIGIN/../lib"
      "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/gmt")
    if(CMAKE_INSTALL_DO_STRIP)
      execute_process(COMMAND "/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/strip" -u -r "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/gmt")
    endif()
  endif()
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/isogmt;/Users/xuning/Work/go-dem/libs/darwin_arm/gmt-config")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm" TYPE PROGRAM FILES
    "/Users/xuning/Work/go-dem/build5/external/gmt/src/isogmt"
    "/Users/xuning/Work/go-dem/build5/external/gmt/src/gmt-config"
    )
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/gmtswitch")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm" TYPE PROGRAM RENAME "gmtswitch" FILES "/Users/xuning/Work/go-dem/external/gmt/src/gmtswitch")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/gmt_shell_functions.sh")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm" TYPE PROGRAM RENAME "gmt_shell_functions.sh" FILES "/Users/xuning/Work/go-dem/external/gmt/src/gmt_shell_functions.sh")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_resources.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/declspec.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include" TYPE FILE FILES
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_resources.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/declspec.h"
    )
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/postscriptlight.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_common_math.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_common_string.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_common.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_constants.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_contour.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_dcw.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_decorate.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_defaults.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_error.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_error_codes.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_fft.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_gdalread.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_glib.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_grd.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_grdio.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_hash.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_io.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_macros.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_memory.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_modern.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_nan.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_notposix.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_plot.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_private.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_project.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_prototypes.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_psl.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_shore.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_common_sighandler.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_symbol.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_synopsis.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_texture.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_time.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_types.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_dev.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_customio.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_hidden.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_mb.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_remote.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_common_byteswap.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include" TYPE FILE FILES
    "/Users/xuning/Work/go-dem/external/gmt/src/postscriptlight.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_common_math.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_common_string.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_common.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_constants.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_contour.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_dcw.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_decorate.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_defaults.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_error.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_error_codes.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_fft.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_gdalread.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_glib.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_grd.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_grdio.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_hash.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_io.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_macros.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_memory.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_modern.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_nan.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_notposix.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_plot.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_private.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_project.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_prototypes.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_psl.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_shore.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_common_sighandler.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_symbol.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_synopsis.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_texture.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_time.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_types.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_dev.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_customio.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_hidden.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_mb.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_remote.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gmt_common_byteswap.h"
    )
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/compat/qsort.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include/compat" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/compat/qsort.h")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/config.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_config.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_dimensions.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gmt_version.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/psl_config.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include" TYPE FILE FILES
    "/Users/xuning/Work/go-dem/build5/external/gmt/src/config.h"
    "/Users/xuning/Work/go-dem/build5/external/gmt/src/gmt_config.h"
    "/Users/xuning/Work/go-dem/build5/external/gmt/src/gmt_dimensions.h"
    "/Users/xuning/Work/go-dem/build5/external/gmt/src/gmt_version.h"
    "/Users/xuning/Work/go-dem/build5/external/gmt/src/psl_config.h"
    )
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/gmt/plugins/supplements.so")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/gmt/plugins" TYPE MODULE FILES "/Users/xuning/Work/go-dem/build5/external/gmt/src/plugins/supplements.so")
  if(EXISTS "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/gmt/plugins/supplements.so" AND
     NOT IS_SYMLINK "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/gmt/plugins/supplements.so")
    execute_process(COMMAND "/usr/bin/install_name_tool"
      -change "@rpath/libgmt.6.dylib" "@executable_path/../lib/libgmt.6.dylib"
      -change "@rpath/libpostscriptlight.6.dylib" "@executable_path/../lib/libpostscriptlight.6.dylib"
      "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/gmt/plugins/supplements.so")
    execute_process(COMMAND /usr/bin/install_name_tool
      -delete_rpath "/Users/xuning/Work/go-dem/build5/external/gmt/src"
      -add_rpath "$ORIGIN/../lib"
      "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/gmt/plugins/supplements.so")
    if(CMAKE_INSTALL_DO_STRIP)
      execute_process(COMMAND "/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/strip" -x "$ENV{DESTDIR}/Users/xuning/Work/go-dem/libs/darwin_arm/gmt/plugins/supplements.so")
    endif()
  endif()
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/geodesy/README.geodesy")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/geodesy" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/geodesy/README.geodesy")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/gsfml/README.gsfml")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/gsfml" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/README.gsfml")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/CK1995n.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/Chron_Normal.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/Chron_Reverse.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/GST2004n.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/GST2012n.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/Geek2007n.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/CK1995r.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/Chron_Normal2.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/Chron_Reverse2.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/GST2004r.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/GST2012r.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/Geek2007r.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml/fz_analysis.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include/gsfml" TYPE FILE FILES
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/CK1995n.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/Chron_Normal.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/Chron_Reverse.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/GST2004n.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/GST2012n.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/Geek2007n.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/CK1995r.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/Chron_Normal2.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/Chron_Reverse2.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/GST2004r.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/GST2012r.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/Geek2007r.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gsfml/fz_analysis.h"
    )
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/gshhg/README.gshhg")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/gshhg" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/gshhg/README.gshhg")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/gshhg/gshhg.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/gshhg/gmt_gshhg.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include/gshhg" TYPE FILE FILES
    "/Users/xuning/Work/go-dem/external/gmt/src/gshhg/gshhg.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/gshhg/gmt_gshhg.h"
    )
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/img/README.img")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/img" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/img/README.img")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/mgd77/README.mgd77")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/mgd77" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/README.mgd77")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/cm4_functions.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77_IGF_coeffs.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77_codes.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77_e77.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77_functions.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77_init.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77_recalc.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77_rls_coeffs.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77defaults.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77magref.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77sniffer.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77/mgd77snifferdefaults.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include/mgd77" TYPE FILE FILES
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/cm4_functions.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77_IGF_coeffs.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77_codes.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77_e77.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77_functions.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77_init.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77_recalc.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77_rls_coeffs.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77defaults.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77magref.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77sniffer.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/mgd77/mgd77snifferdefaults.h"
    )
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/potential/README.potential")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/potential" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/potential/README.potential")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/potential/okbfuns.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/potential/newton.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/potential/modeltime.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/potential/talwani.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include/potential" TYPE FILE FILES
    "/Users/xuning/Work/go-dem/external/gmt/src/potential/okbfuns.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/potential/newton.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/potential/modeltime.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/potential/talwani.h"
    )
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/segy/README.segy")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/segy" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/segy/README.segy")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/segy/segy.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/segy/segy_io.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/segy/segyreel.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include/segy" TYPE FILE FILES
    "/Users/xuning/Work/go-dem/external/gmt/src/segy/segy.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/segy/segy_io.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/segy/segyreel.h"
    )
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/seis/README.seis")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/seis" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/seis/README.seis")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/seis/meca.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/seis/meca_symbol.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/seis/utilmeca.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/seis/seis_defaults.h;/Users/xuning/Work/go-dem/libs/darwin_arm/include/seis/sacio.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include/seis" TYPE FILE FILES
    "/Users/xuning/Work/go-dem/external/gmt/src/seis/meca.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/seis/meca_symbol.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/seis/utilmeca.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/seis/seis_defaults.h"
    "/Users/xuning/Work/go-dem/external/gmt/src/seis/sacio.h"
    )
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/spotter/README.spotter")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/spotter" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/spotter/README.spotter")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/spotter/spotter.sh")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/spotter" TYPE PROGRAM FILES "/Users/xuning/Work/go-dem/external/gmt/src/spotter/spotter.sh")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/spotter/spotter.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include/spotter" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/spotter/spotter.h")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/x2sys/README.x2sys")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/x2sys" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/x2sys/README.x2sys")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/x2sys/test_x2sys.sh")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/x2sys" TYPE PROGRAM FILES "/Users/xuning/Work/go-dem/external/gmt/src/x2sys/test_x2sys.sh")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Runtime" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/Users/xuning/Work/go-dem/libs/darwin_arm/include/x2sys/x2sys.h")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/Users/xuning/Work/go-dem/libs/darwin_arm/include/x2sys" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/x2sys/x2sys.h")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/windbarbs/README.windbarb")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/windbarbs" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/windbarbs/README.windbarb")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Documentation" OR NOT CMAKE_INSTALL_COMPONENT)
  list(APPEND CMAKE_ABSOLUTE_DESTINATION_FILES
   "/supplements/nswing/README.nswing")
  if(CMAKE_WARN_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(WARNING "ABSOLUTE path INSTALL DESTINATION : ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  if(CMAKE_ERROR_ON_ABSOLUTE_INSTALL_DESTINATION)
    message(FATAL_ERROR "ABSOLUTE path INSTALL DESTINATION forbidden (by caller): ${CMAKE_ABSOLUTE_DESTINATION_FILES}")
  endif()
  file(INSTALL DESTINATION "/supplements/nswing" TYPE FILE FILES "/Users/xuning/Work/go-dem/external/gmt/src/nswing/README.nswing")
endif()

if(NOT CMAKE_INSTALL_LOCAL_ONLY)
  # Include the install script for each subdirectory.
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/geodesy/cmake_install.cmake")
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/gsfml/cmake_install.cmake")
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/gshhg/cmake_install.cmake")
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/img/cmake_install.cmake")
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/mgd77/cmake_install.cmake")
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/potential/cmake_install.cmake")
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/segy/cmake_install.cmake")
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/seis/cmake_install.cmake")
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/spotter/cmake_install.cmake")
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/x2sys/cmake_install.cmake")
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/windbarbs/cmake_install.cmake")
  include("/Users/xuning/Work/go-dem/build5/external/gmt/src/nswing/cmake_install.cmake")

endif()

string(REPLACE ";" "\n" CMAKE_INSTALL_MANIFEST_CONTENT
       "${CMAKE_INSTALL_MANIFEST_FILES}")
if(CMAKE_INSTALL_LOCAL_ONLY)
  file(WRITE "/Users/xuning/Work/go-dem/build5/external/gmt/src/install_local_manifest.txt"
     "${CMAKE_INSTALL_MANIFEST_CONTENT}")
endif()
