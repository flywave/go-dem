#!/bin/bash
# Build all external dependencies for go-dem
# GDAL + PROJ + GEOS + ZLIB + PNG + JPEG + EXPAT + SQLite3 + WEBP + ICONV + HDF5 + netCDF + GMT
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEM_DIR="$(dirname "$SCRIPT_DIR")"
cd "$DEM_DIR"

case "$(uname -s)" in
    Darwin)
        [ "$(uname -m)" = "arm64" ] && PLATFORM="darwin_arm" || PLATFORM="darwin"
        ;;
    Linux)
        [ "$(uname -m)" = "aarch64" ] && PLATFORM="linux_arm" || PLATFORM="linux"
        ;;
esac

INSTALL_DIR="$DEM_DIR/libs/$PLATFORM"
JOBS=$(sysctl -n hw.ncpu 2>/dev/null || nproc 2>/dev/null || echo 4)
echo "Platform: $PLATFORM  Install: $INSTALL_DIR  Jobs: $JOBS"

#===============================================================================
# Phase 1: GDAL ecosystem (uses flywave-gdal's CMake externally)
#===============================================================================
echo "=== Phase 1: Building GDAL + dependencies ==="
mkdir -p build && cd build
cmake .. \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_PREFIX="$INSTALL_DIR" \
    -DBUILD_SHARED_LIBS=OFF
cmake --build . -j$JOBS
cd "$DEM_DIR"

# Collect .a files
echo "=== Installing GDAL ecosystem libraries ==="
mkdir -p "$INSTALL_DIR/lib" "$INSTALL_DIR/include"
find build -name "*.a" -exec cp {} "$INSTALL_DIR/lib/" \; 2>/dev/null
# Copy headers
for dir in external/libgdal/gdal/port external/libgdal/gdal/gcore external/libgdal/gdal/alg external/libgdal/gdal/ogr; do
    [ -d "$dir" ] && cp "$dir"/*.h "$INSTALL_DIR/include/" 2>/dev/null || true
done

#===============================================================================
# Phase 2: HDF5
#===============================================================================
echo "=== Phase 2: Building HDF5 ==="
cd external/hdf5
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
    -DHDF5_ENABLE_PARALLEL=OFF \
    -DHDF5_ALLOW_EXTERNAL_SUPPORT=NO
cmake --build . --target install -j$JOBS
cd "$DEM_DIR"

#===============================================================================
# Phase 3: netCDF
#===============================================================================
echo "=== Phase 3: Building netCDF ==="
cd external/netcdf
mkdir -p build && cd build
cmake .. \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_PREFIX="$INSTALL_DIR" \
    -DBUILD_SHARED_LIBS=OFF \
    -DBUILD_TESTING=OFF \
    -DENABLE_DAP=OFF \
    -DENABLE_BYTERANGE=OFF \
    -DNETCDF_ENABLE_HDF5=ON \
    -DNETCDF_ENABLE_NCZARR=OFF \
    -DNETCDF_ENABLE_FILTER_SZIP=OFF \
    -DNETCDF_ENABLE_PARALLEL4=OFF \
    -DHDF5_ROOT="$INSTALL_DIR" \
    -DHDF5_USE_STATIC_LIBRARIES=ON
cmake --build . --target install -j$JOBS
cd "$DEM_DIR"

#===============================================================================
# Phase 4: GMT
#===============================================================================
echo "=== Phase 4: Building GMT ==="
cd external/gmt
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
    -DGDAL_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DGDAL_LIBRARY="$INSTALL_DIR/lib/libgdal_c.a" \
    -DPROJ_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DPROJ_LIBRARY="$INSTALL_DIR/lib/libproj.a" \
    -DGEOS_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DGEOS_LIBRARY="$INSTALL_DIR/lib/libgeos.a" \
    -DSQLITE3_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DSQLITE3_LIBRARY="$INSTALL_DIR/lib/libsqlite3.a" \
    -DZLIB_LIBRARY="$INSTALL_DIR/lib/libzlib.a" \
    -DZLIB_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DCURL_LIBRARY="" \
    -DCURL_INCLUDE_DIR="" \
    -DGMT_NO_CURL=ON
cmake --build . --target install -j$JOBS
cd "$DEM_DIR"

echo "=== Done ==="
echo "Libraries installed to: $INSTALL_DIR"
ls -la "$INSTALL_DIR/lib/" | grep -E "\.a"
