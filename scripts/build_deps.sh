#!/bin/bash
# Build external dependencies for go-dem
# HDF5 → netCDF → GMT
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEM_DIR="$(dirname "$SCRIPT_DIR")"
PLATFORM="darwin_arm"  # change for other platforms: darwin, linux, linux_arm

case "$(uname -s)" in
    Darwin)
        if [ "$(uname -m)" = "arm64" ]; then
            PLATFORM="darwin_arm"
        else
            PLATFORM="darwin"
        fi
        ;;
    Linux)
        if [ "$(uname -m)" = "aarch64" ]; then
            PLATFORM="linux_arm"
        else
            PLATFORM="linux"
        fi
        ;;
esac

INSTALL_DIR="$DEM_DIR/libs/$PLATFORM"
FLYWAVE_GDAL_DIR="$DEM_DIR/../../flywave-gdal"
FLYWAVE_LIB="$FLYWAVE_GDAL_DIR/libs/$PLATFORM"
JOBS=$(sysctl -n hw.ncpu 2>/dev/null || nproc 2>/dev/null || echo 4)

echo "Building for $PLATFORM (install: $INSTALL_DIR, jobs: $JOBS)"

#===============================================================================
# 1. HDF5
#===============================================================================
echo "=== Building HDF5 ==="
cd "$DEM_DIR/external/hdf5"
mkdir -p build && cd build
cmake .. \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_PREFIX="$INSTALL_DIR" \
    -DBUILD_SHARED_LIBS=OFF \
    -DBUILD_TESTING=OFF \
    -DHDF5_BUILD_TOOLS=OFF \
    -DHDF5_BUILD_EXAMPLES=OFF \
    -DHDF5_BUILD_FORTRAN=OFF \
    -DHDF5_BUILD_CPP_LIB=OFF \
    -DHDF5_BUILD_JAVA=OFF \
    -DHDF5_BUILD_HL_LIB=OFF \
    -DHDF5_ENABLE_ZLIB_SUPPORT=OFF \
    -DHDF5_ENABLE_SZIP_SUPPORT=OFF \
    -DHDF5_ENABLE_PARALLEL=OFF \
    -DHDF5_ENABLE_THREADSAFE=OFF \
    -DHDF5_ALLOW_EXTERNAL_SUPPORT=NO
cmake --build . --target install -j$JOBS
cd "$DEM_DIR"

#===============================================================================
# 2. netCDF
#===============================================================================
echo "=== Building netCDF ==="
cd "$DEM_DIR/external/netcdf"
mkdir -p build && cd build
cmake .. \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_PREFIX="$INSTALL_DIR" \
    -DBUILD_SHARED_LIBS=OFF \
    -DBUILD_TESTING=OFF \
    -DENABLE_TESTS=OFF \
    -DENABLE_DAP=OFF \
    -DENABLE_BYTERANGE=OFF \
    -DNETCDF_ENABLE_HDF5=ON \
    -DNETCDF_ENABLE_NCZARR=OFF \
    -DNETCDF_ENABLE_FILTER_SZIP=OFF \
    -DNETCDF_ENABLE_PARALLEL4=OFF \
    -DNETCDF_ENABLE_EXAMPLE_TESTS=OFF \
    -DHDF5_ROOT="$INSTALL_DIR" \
    -DHDF5_USE_STATIC_LIBRARIES=ON
cmake --build . --target install -j$JOBS
cd "$DEM_DIR"

#===============================================================================
# 3. GMT
#===============================================================================
echo "=== Building GMT ==="
cd "$DEM_DIR/external/gmt"
mkdir -p build && cd build
cmake .. \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_PREFIX="$INSTALL_DIR" \
    -DBUILD_SHARED_LIBS=OFF \
    -DBUILD_TESTING=OFF \
    -DNETCDF_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DNETCDF_LIBRARY="$INSTALL_DIR/lib/libnetcdf.a" \
    -DHDF5_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DHDF5_LIBRARY="$INSTALL_DIR/lib/libhdf5.a" \
    -DGDAL_INCLUDE_DIR="$FLYWAVE_GDAL_DIR/external/libgdal/include" \
    -DGDAL_LIBRARY="$FLYWAVE_LIB/libgdal_c.a" \
    -DPROJ_INCLUDE_DIR="$FLYWAVE_GDAL_DIR/external/libproj/include" \
    -DPROJ_LIBRARY="$FLYWAVE_LIB/libproj.a" \
    -DGEOS_INCLUDE_DIR="$DEM_DIR/../../go-geos/libs" \
    -DGEOS_LIBRARY="$DEM_DIR/../../go-geos/libs/$PLATFORM/libgeos.a" \
    -DPCRE_ROOT="/opt/homebrew" \
    -DZLIB_LIBRARY="$FLYWAVE_LIB/libzlib.a" \
    -DZLIB_INCLUDE_DIR="$FLYWAVE_GDAL_DIR/external/zlib" \
    -DCURL_LIBRARY="" \
    -DCURL_INCLUDE_DIR="" \
    -DGMT_NO_CURL=ON
cmake --build . --target install -j$JOBS
cd "$DEM_DIR"

echo "=== Done ==="
echo "Libraries installed to: $INSTALL_DIR"
ls -la "$INSTALL_DIR/lib/" | grep -E "hdf5|netcdf|gmt"
