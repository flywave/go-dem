#!/bin/bash
# One-command build all external dependencies
# Usage: ./scripts/build_deps.sh [gmt]
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEM_DIR="$(dirname "$SCRIPT_DIR")"
cd "$DEM_DIR"

JOBS=$(sysctl -n hw.ncpu 2>/dev/null || nproc 2>/dev/null || echo 4)

echo "=== Phase 1: GDAL ecosystem (zlib/png/jpeg/expat/iconv/sqlite3/geos/proj/webp/gdal/htdp) ==="
mkdir -p build && cd build
cmake .. -DCMAKE_BUILD_TYPE=Release
cmake --build . -j$JOBS
cmake --install . --prefix "$DEM_DIR/libs/tmp" 2>/dev/null || true
cd "$DEM_DIR"

echo "=== Phase 2: GMT (zero external args, config in CMakeLists.txt) ==="
cmake --build build --target gmt_build -j$JOBS

echo "=== Done ==="
echo "Libraries:"
ls "$DEM_DIR"/libs/darwin_arm/*.a 2>/dev/null | wc -l
echo "files in libs/darwin_arm/"
