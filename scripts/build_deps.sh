#!/bin/bash
# One-click build all external dependencies for go-dem
# Phase 1: GDAL ecosystem (zlib/png/jpeg/expat/iconv/sqlite3/geos/proj/webp/gdal)
# Phase 2: HDF5 → netCDF → GMT
# All output to libs/{platform}/
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEM_DIR="$(dirname "$SCRIPT_DIR")"
cd "$DEM_DIR"

case "$(uname -s)" in
    Darwin) [ "$(uname -m)" = "arm64" ] && PLATFORM="darwin_arm" || PLATFORM="darwin" ;;
    Linux)  [ "$(uname -m)" = "aarch64" ] && PLATFORM="linux_arm" || PLATFORM="linux" ;;
esac

INSTALL_DIR="$DEM_DIR/libs/$PLATFORM"
JOBS=$(sysctl -n hw.ncpu 2>/dev/null || nproc 2>/dev/null || echo 4)
echo "Platform: $PLATFORM  Install: $INSTALL_DIR  Jobs: $JOBS"

#===============================================================================
# Phase 1: GDAL ecosystem (root CMake builds all)
#===============================================================================
echo "=== Phase 1: GDAL ecosystem (zlib/png/jpeg/expat/iconv/sqlite3/geos/proj/webp/gdal) ==="
mkdir -p build && cd build
cmake .. -DCMAKE_BUILD_TYPE=Release
cmake --build . -j$JOBS
cmake --install . --prefix "$INSTALL_DIR" 2>/dev/null || true
cd "$DEM_DIR"

# Ensure .a files are in the right place
echo "=== Installing .a files ==="
find build -name "*.a" | while read f; do
    cp "$f" "$INSTALL_DIR/" 2>/dev/null || true
done
echo "  $(ls "$INSTALL_DIR"/*.a 2>/dev/null | wc -l) .a files"

# Ensure headers are in place
echo "=== Installing headers ==="
mkdir -p "$INSTALL_DIR/include"
for dir in external/libgdal/gdal/port external/libgdal/gdal/gcore \
           external/libgdal/gdal/alg external/libgdal/gdal/ogr \
           external/libproj/proj/src external/libproj/proj/include \
           external/libgeos/geos/include external/libgeos/geos/capi; do
    [ -d "$dir" ] && find "$dir" -name "*.h" -exec cp {} "$INSTALL_DIR/include/" \; 2>/dev/null
done
echo "  $(ls "$INSTALL_DIR/include"/*.h 2>/dev/null | wc -l) headers"

#===============================================================================
# Phase 2: HDF5
#===============================================================================
echo "=== Phase 2: HDF5 ==="
cd external/hdf5 && mkdir -p build && cd build
cmake .. -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_PREFIX="$INSTALL_DIR" \
    -DBUILD_SHARED_LIBS=OFF -DBUILD_TESTING=OFF \
    -DHDF5_BUILD_TOOLS=OFF -DHDF5_BUILD_EXAMPLES=OFF \
    -DHDF5_BUILD_FORTRAN=OFF -DHDF5_BUILD_CPP_LIB=OFF \
    -DHDF5_BUILD_JAVA=OFF -DHDF5_BUILD_HL_LIB=OFF \
    -DHDF5_ENABLE_ZLIB_SUPPORT=OFF -DHDF5_ENABLE_PARALLEL=OFF \
    -DHDF5_ALLOW_EXTERNAL_SUPPORT=NO
cmake --build . --target install -j$JOBS
cd "$DEM_DIR"

#===============================================================================
# Phase 3: netCDF
#===============================================================================
echo "=== Phase 3: netCDF ==="
cd external/netcdf && mkdir -p build && cd build
cmake .. -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_PREFIX="$INSTALL_DIR" \
    -DBUILD_SHARED_LIBS=OFF -DBUILD_TESTING=OFF \
    -DENABLE_DAP=OFF -DENABLE_BYTERANGE=OFF \
    -DNETCDF_ENABLE_HDF5=ON -DNETCDF_ENABLE_NCZARR=OFF \
    -DHDF5_ROOT="$INSTALL_DIR" -DHDF5_USE_STATIC_LIBRARIES=ON
cmake --build . --target install -j$JOBS
cd "$DEM_DIR"

#===============================================================================
# Phase 4: GMT
#===============================================================================
echo "=== Phase 4: GMT ==="
cd external/gmt && mkdir -p build && cd build
cmake .. -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_PREFIX="$INSTALL_DIR" -DBUILD_SHARED_LIBS=OFF \
    -DBUILD_TESTING=OFF \
    -DNETCDF_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DNETCDF_LIBRARY="$INSTALL_DIR/lib/libnetcdf.a" \
    -DHDF5_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DHDF5_LIBRARY="$INSTALL_DIR/lib/libhdf5.a" \
    -DGDAL_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DGDAL_LIBRARY="$INSTALL_DIR/libgdal_c.a" \
    -DPROJ_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DPROJ_LIBRARY="$INSTALL_DIR/libproj.a" \
    -DGEOS_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DGEOS_LIBRARY="$INSTALL_DIR/libgeos.a" \
    -DSQLITE3_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DSQLITE3_LIBRARY="$INSTALL_DIR/libsqlite3.a" \
    -DZLIB_LIBRARY="$INSTALL_DIR/libzlib.a" -DZLIB_INCLUDE_DIR="$INSTALL_DIR/include" \
    -DCURL_LIBRARY="" -DCURL_INCLUDE_DIR="" -DGMT_NO_CURL=ON
cmake --build . --target install -j$JOBS 2>&1 | tail -20
cd "$DEM_DIR"

echo "=== Build complete ==="
ls "$INSTALL_DIR"/*.a
