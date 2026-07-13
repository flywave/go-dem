cmake_minimum_required(VERSION 2.8.12)
project(iconv)

if(UNIX AND NOT APPLE)
  SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -fPIC ")
  SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fPIC ")
ENDIF()

add_definitions(
  -DLIBICONV_PLUG
)
set(TARGET_NAME iconv)
add_subdirectory("${CMAKE_CURRENT_LIST_DIR}/src" src)
