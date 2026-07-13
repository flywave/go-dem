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

func cHTDPTransformPoint(lat, lon, h float64, srcID, dstID int, srcEpoch, dstEpoch float64) (float64, float64, float64, error) {
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

func cHTDPGetVelocity(lat, lon, h float64) (vn, ve, vu float64, err error) {
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

func cHTDPSetGridPath(path string) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.htdp_set_grid_path(cPath)
}
