package datum

/*
#cgo CFLAGS: -I${SRCDIR}/../external/HTDP/capi
#cgo darwin CFLAGS: -I/opt/homebrew/include
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../libs/darwin -lhtdp -lm -L/opt/homebrew/lib/gcc/current -lgfortran
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../libs/darwin_arm -lhtdp -lm -L/opt/homebrew/lib/gcc/current -lgfortran
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../libs/linux -lhtdp -lm -lgfortran
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../libs/linux_arm -lhtdp -lm -lgfortran

#include "htdp_capi.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

func HTDPTransformPoint(lat, lon, h float64, srcID, dstID int, srcEpoch, dstEpoch float64) (float64, float64, float64, error) {
	var outLat, outLon, outH C.double
	ret := C.htdp_transform(
		C.double(lat), C.double(lon), C.double(h),
		C.int(srcID), C.double(srcEpoch),
		C.int(dstID), C.double(dstEpoch),
		&outLat, &outLon, &outH,
	)
	if ret != 0 {
		return 0, 0, 0, errHTDPFailed
	}
	return float64(outLat), float64(outLon), float64(outH), nil
}

func HTDPGetVelocity(lat, lon, h float64) (vn, ve, vu float64, err error) {
	var cVn, cVe, cVu C.double
	ret := C.htdp_velocity(
		C.double(lat), C.double(lon), C.double(h),
		&cVn, &cVe, &cVu,
	)
	if ret != 0 {
		return 0, 0, 0, errHTDPFailed
	}
	return float64(cVn), float64(cVe), float64(cVu), nil
}

func HTDPSetGridPath(path string) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.htdp_set_grid_path(cPath)
}

func cHTDPGrid(gridDef [6]float64, srcID, dstID int, srcEpoch, dstEpoch float64) []float64 {
	xCount := int(gridDef[4])
	yCount := int(gridDef[5])
	n := xCount * yCount
	result := make([]float64, n)

	xMin, yMax := gridDef[0], gridDef[1]
	xMax, yMin := gridDef[2], gridDef[3]
	xInc := (xMax - xMin) / float64(xCount-1)
	yInc := (yMin - yMax) / float64(yCount-1)

	idx := 0
	for y := 0; y < yCount; y++ {
		lat := yMax + float64(y)*yInc
		for x := 0; x < xCount; x++ {
			lon := xMin + float64(x)*xInc
			if outLat, outLon, outH, err := HTDPTransformPoint(lat, lon, 0, srcID, dstID, srcEpoch, dstEpoch); err == nil {
				_ = outLat
				_ = outLon
				result[idx] = outH
			}
			idx++
		}
	}
	return result
}
