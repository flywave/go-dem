package gmt

/*
#cgo CFLAGS: -I${SRCDIR}/../../external/gmt/capi -I${SRCDIR}/../../external/gmt/src
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../../libs/darwin -lgmt -lpsl -lm -L/opt/homebrew/lib -lnetcdf -L${SRCDIR}/../../libs/darwin -lhdf5 -lhdf5_hl -lgdal_c -lproj -lgeos -lsqlite3 -lzlib -lpng -ljpeg -lexpat -liconv -framework Accelerate
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../../libs/darwin_arm -lgmt -lpsl -lm -L/opt/homebrew/lib -lnetcdf -L${SRCDIR}/../../libs/darwin_arm -lhdf5 -lhdf5_hl -lgdal_c -lproj -lgeos -lsqlite3 -lzlib -lpng -ljpeg -lexpat -liconv -lc++ -framework Accelerate
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../../libs/linux -lgmt -lpsl -lm -lnetcdf -lhdf5 -lhdf5_hl -lgdal_c -lproj -lgeos -lsqlite3 -lz -lpng -ljpeg -lexpat -liconv
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../../libs/linux_arm -lgmt -lpsl -lm -lnetcdf -lhdf5 -lhdf5_hl -lgdal_c -lproj -lgeos -lsqlite3 -lz -lpng -ljpeg -lexpat -liconv

#include <stdlib.h>
#include "gmt_capi.h"
#include "../../external/gmt/capi/gmt_capi.c"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type GridConfig struct {
	XInc, YInc   float64
	XMin, XMax   float64
	YMin, YMax   float64
	Tension      float64
	SearchRadius float64
	EmptyValue   int
}

func init() {
	C.gmt_begin()
}

func Surface(inputPath, outputPath string, cfg *GridConfig) error {
	if cfg.Tension <= 0 {
		cfg.Tension = 0.25
	}
	cIn := C.CString(inputPath)
	cOut := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cIn))
	defer C.free(unsafe.Pointer(cOut))
	ret := C.gmt_surface(cIn, cOut,
		C.double(cfg.Tension),
		C.double(cfg.XInc), C.double(cfg.YInc),
		C.double(cfg.XMin), C.double(cfg.XMax),
		C.double(cfg.YMin), C.double(cfg.YMax))
	if ret != 0 {
		return fmt.Errorf("gmt surface failed with code %d", int(ret))
	}
	return nil
}

func Grdfilter(inputPath, outputPath, filterType, distFlag string) error {
	cIn := C.CString(inputPath)
	cOut := C.CString(outputPath)
	cFilt := C.CString(filterType)
	cDist := C.CString(distFlag)
	defer C.free(unsafe.Pointer(cIn))
	defer C.free(unsafe.Pointer(cOut))
	defer C.free(unsafe.Pointer(cFilt))
	defer C.free(unsafe.Pointer(cDist))
	ret := C.gmt_grdfilter(cIn, cOut, cFilt, cDist)
	if ret != 0 {
		return fmt.Errorf("gmt grdfilter failed with code %d", int(ret))
	}
	return nil
}

func Triangulate(inputPath, outputPath string, cfg *GridConfig) error {
	cIn := C.CString(inputPath)
	cOut := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cIn))
	defer C.free(unsafe.Pointer(cOut))
	ret := C.gmt_triangulate(cIn, cOut,
		C.double(cfg.XInc), C.double(cfg.YInc),
		C.double(cfg.XMin), C.double(cfg.XMax),
		C.double(cfg.YMin), C.double(cfg.YMax))
	if ret != 0 {
		return fmt.Errorf("gmt triangulate failed with code %d", int(ret))
	}
	return nil
}

func Blockmean(inputPath, outputPath string, cfg *GridConfig) error {
	cIn := C.CString(inputPath)
	cOut := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cIn))
	defer C.free(unsafe.Pointer(cOut))
	ret := C.gmt_blockmean(cIn, cOut,
		C.double(cfg.XInc), C.double(cfg.YInc),
		C.double(cfg.XMin), C.double(cfg.XMax),
		C.double(cfg.YMin), C.double(cfg.YMax))
	if ret != 0 {
		return fmt.Errorf("gmt blockmean failed with code %d", int(ret))
	}
	return nil
}

func Nearneighbor(inputPath, outputPath string, cfg *GridConfig) error {
	cIn := C.CString(inputPath)
	cOut := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cIn))
	defer C.free(unsafe.Pointer(cOut))
	sr := cfg.SearchRadius
	if sr <= 0 {
		sr = cfg.XInc * 5
	}
	if cfg.EmptyValue == 0 {
		cfg.EmptyValue = -9999
	}
	ret := C.gmt_nearneighbor(cIn, cOut,
		C.double(cfg.XInc), C.double(cfg.YInc),
		C.double(cfg.XMin), C.double(cfg.XMax),
		C.double(cfg.YMin), C.double(cfg.YMax),
		C.double(sr), C.int(cfg.EmptyValue))
	if ret != 0 {
		return fmt.Errorf("gmt nearneighbor failed with code %d", int(ret))
	}
	return nil
}
