cmake_minimum_required(VERSION 3.1.3)

project(htdp Fortran C)

enable_language(Fortran)

set(HTDP_SRC_DIR ${CMAKE_CURRENT_SOURCE_DIR})

# Fortran source files
set(HTDP_FSRC
    ${HTDP_SRC_DIR}/htdp.f
    ${HTDP_SRC_DIR}/initvl.f
    ${HTDP_SRC_DIR}/initeq.f
    ${HTDP_SRC_DIR}/initbd.f
    ${HTDP_SRC_DIR}/initps.f
)

# C wrapper source
set(HTDP_CSRC
    ${HTDP_SRC_DIR}/capi/htdp_capi.c
)

# Build as static library
set(CMAKE_Fortran_FLAGS "${CMAKE_Fortran_FLAGS} -fPIC -O2 -ffixed-line-length-132")

# Generate a version of htdp.f without PROGRAM main to avoid duplicate _main
set(HTDP_LIB_F ${CMAKE_CURRENT_BINARY_DIR}/htdp_lib.f)
file(READ ${HTDP_SRC_DIR}/htdp.f HTDP_SRC_CONTENT)
string(REGEX REPLACE "PROGRAM[ ]+HTDP" "SUBROUTINE HTDPMAIN" HTDP_LIB_CONTENT "${HTDP_SRC_CONTENT}")
file(WRITE ${HTDP_LIB_F} "${HTDP_LIB_CONTENT}")
message(STATUS "Generated htdp_lib.f (PROGRAM → SUBROUTINE)")

add_library(htdp STATIC
    ${HTDP_LIB_F}
    ${HTDP_SRC_DIR}/initvl.f
    ${HTDP_SRC_DIR}/initeq.f
    ${HTDP_SRC_DIR}/initbd.f
    ${HTDP_SRC_DIR}/initps.f
    ${HTDP_CSRC}
)

target_include_directories(htdp PUBLIC
    ${HTDP_SRC_DIR}/capi
)

# Fortran compiler flags
if(APPLE)
    set_target_properties(htdp PROPERTIES
        Fortran_FLAGS "-fPIC -O2 -ffixed-line-length-132"
    )
elseif(UNIX)
    set_target_properties(htdp PROPERTIES
        Fortran_FLAGS "-fPIC -O2 -ffixed-line-length-132"
    )
endif()

# Also build the htdp executable (for os/exec use)
add_executable(htdp_exe
    ${HTDP_FSRC}
)

set_target_properties(htdp_exe PROPERTIES
    OUTPUT_NAME htdp
    Fortran_FLAGS "-O2 -ffixed-line-length-132"
)

if(APPLE)
    set_target_properties(htdp_exe PROPERTIES
        Fortran_FLAGS "-O2 -ffixed-line-length-132"
    )
endif()

# Output directories
set(FLYWAVE_BASE_INSTALL_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../../libs")

if(APPLE)
    if(${CMAKE_SYSTEM_PROCESSOR} MATCHES "^arm")
        set(FLYWAVE_LIBRARY_OUTPUT_PATH "${FLYWAVE_BASE_INSTALL_DIR}/darwin_arm")
    else()
        set(FLYWAVE_LIBRARY_OUTPUT_PATH "${FLYWAVE_BASE_INSTALL_DIR}/darwin")
    endif()
elseif(WIN32)
    set(FLYWAVE_LIBRARY_OUTPUT_PATH "${FLYWAVE_BASE_INSTALL_DIR}/windows")
else()
    if(CMAKE_SYSTEM_PROCESSOR MATCHES "aarch64")
        set(FLYWAVE_LIBRARY_OUTPUT_PATH "${FLYWAVE_BASE_INSTALL_DIR}/linux_arm")
    else()
        set(FLYWAVE_LIBRARY_OUTPUT_PATH "${FLYWAVE_BASE_INSTALL_DIR}/linux")
    endif()
endif()

make_directory(${FLYWAVE_LIBRARY_OUTPUT_PATH})

set_target_properties(htdp PROPERTIES
    ARCHIVE_OUTPUT_DIRECTORY_DEBUG ${FLYWAVE_LIBRARY_OUTPUT_PATH}
    ARCHIVE_OUTPUT_DIRECTORY_RELEASE ${FLYWAVE_LIBRARY_OUTPUT_PATH}
    LIBRARY_OUTPUT_DIRECTORY_DEBUG ${FLYWAVE_LIBRARY_OUTPUT_PATH}
    LIBRARY_OUTPUT_DIRECTORY_RELEASE ${FLYWAVE_LIBRARY_OUTPUT_PATH}
)

set_target_properties(htdp_exe PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY_DEBUG ${FLYWAVE_LIBRARY_OUTPUT_PATH}
    RUNTIME_OUTPUT_DIRECTORY_RELEASE ${FLYWAVE_LIBRARY_OUTPUT_PATH}
)

# Install rules
install(TARGETS htdp
    ARCHIVE DESTINATION ${FLYWAVE_LIBRARY_OUTPUT_PATH}
    LIBRARY DESTINATION ${FLYWAVE_LIBRARY_OUTPUT_PATH}
)

install(TARGETS htdp_exe
    RUNTIME DESTINATION ${FLYWAVE_LIBRARY_OUTPUT_PATH}
)

install(DIRECTORY ${HTDP_SRC_DIR}/capi/
    DESTINATION ${FLYWAVE_BASE_INSTALL_DIR}
    FILES_MATCHING PATTERN "*.h"
)
