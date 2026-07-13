//go:build htdp

package datum

/*
#cgo CFLAGS: -I${SRCDIR}/../external/HTDP/capi
#cgo darwin CFLAGS: -I/opt/homebrew/include
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../libs/darwin -lhtdp -lm -lc++ -L/opt/homebrew/lib/gcc/current -lgfortran -lquadmath
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../libs/darwin_arm -lhtdp -lm -lc++ -L/opt/homebrew/lib/gcc/current -lgfortran -lquadmath
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../libs/linux -lhtdp -lm -lgfortran -lquadmath
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../libs/linux_arm -lhtdp -lm -lgfortran -lquadmath

#include "htdp_capi.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"unsafe"
)

func htdpAvailable() bool {
	_, err := exec.LookPath("htdp")
	if err == nil {
		return true
	}
	libPath := filepath.Join("libs", "darwin_arm", "htdp")
	if _, err := os.Stat(libPath); err == nil {
		return true
	}
	libPath = filepath.Join("libs", "darwin", "htdp")
	if _, err := os.Stat(libPath); err == nil {
		return true
	}
	return false
}

type HTDPResult struct {
	Lat, Lon, H float64
}

func HTDPTransformPoint(lat, lon, h float64, srcID, dstID int, srcEpoch, dstEpoch float64) (*HTDPResult, error) {
	if !htdpAvailable() {
		return nil, fmt.Errorf("HTDP not available: install htdp or set PATH")
	}

	if gridPath := os.Getenv("HTDP_GRID_PATH"); gridPath != "" {
		cPath := C.CString(gridPath)
		defer C.free(unsafe.Pointer(cPath))
		C.htdp_set_grid_path(cPath)
	}

	var outLat, outLon, outH C.double
	ret := C.htdp_transform(
		C.double(lat), C.double(lon), C.double(h),
		C.int(srcID), C.double(srcEpoch),
		C.int(dstID), C.double(dstEpoch),
		&outLat, &outLon, &outH,
	)
	if ret != 0 {
		return nil, fmt.Errorf("HTDP transform failed with code %d", int(ret))
	}

	return &HTDPResult{
		Lat: float64(outLat),
		Lon: float64(outLon),
		H:   float64(outH),
	}, nil
}

func HTDPGetVelocity(lat, lon, h float64) (vn, ve, vu float64, err error) {
	if !htdpAvailable() {
		return 0, 0, 0, fmt.Errorf("HTDP not available")
	}

	var cVn, cVe, cVu C.double
	ret := C.htdp_velocity(
		C.double(lat), C.double(lon), C.double(h),
		&cVn, &cVe, &cVu,
	)
	if ret != 0 {
		return 0, 0, 0, fmt.Errorf("HTDP velocity failed")
	}

	return float64(cVn), float64(cVe), float64(cVu), nil
}

func HTDPExecGrid(homeDir string, srcHTDPID, dstHTDPID int, srcEpoch, dstEpoch float64, gridDef [6]float64) ([]float64, error) {
	tmpInput := filepath.Join(homeDir, "_tmp_input.xyz")
	tmpControl := filepath.Join(homeDir, "_tmp_control.txt")
	tmpOutput := filepath.Join(homeDir, "_tmp_output.xyz")

	writeHTDPInput(gridDef, tmpInput)
	writeHTDPControl(tmpControl, tmpOutput, tmpInput, srcHTDPID, srcEpoch, dstHTDPID, dstEpoch)

	htdpPath := "htdp"
	if _, err := exec.LookPath("htdp"); err != nil {
		htdpPath = filepath.Join("libs", "darwin_arm", "htdp")
	}

	cmd := exec.Command(htdpPath, tmpControl)
	cmd.Dir = homeDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("htdp run: %v, output: %s", err, string(out))
	}

	return readHTDPOutput(tmpOutput, gridDef)
}

func writeHTDPInput(gridDef [6]float64, path string) {
	xCount := int(gridDef[4])
	yCount := int(gridDef[5])

	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	xMin, yMax := gridDef[0], gridDef[1]
	xMax, yMin := gridDef[2], gridDef[3]
	xInc := (xMax - xMin) / float64(xCount-1)
	yInc := (yMin - yMax) / float64(yCount-1)

	for y := 0; y < yCount; y++ {
		lat := yMax + float64(y)*yInc
		for x := 0; x < xCount; x++ {
			lon := xMin + float64(x)*xInc
			fmt.Fprintf(f, "%.8f %.8f 0.0\n", lat, lon)
		}
	}
}

func writeHTDPControl(ctrlPath, outPath, inPath string, srcID int, srcEpoch float64, dstID int, dstEpoch float64) {
	content := fmt.Sprintf(
		`I
%s
%s
1
3
%.4f %d %.4f %d
5
0
0
`,
		inPath, outPath,
		srcEpoch, srcID, dstEpoch, dstID,
	)
	os.WriteFile(ctrlPath, []byte(content), 0644)
}

func readHTDPOutput(path string, gridDef [6]float64) ([]float64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	xCount := int(gridDef[4])
	yCount := int(gridDef[5])
	n := xCount * yCount
	result := make([]float64, n)

	line := 0
	for y := 0; y < yCount; y++ {
		for x := 0; x < xCount; x++ {
			var lat, lon, h float64
			var fn string
			for line < len(data) && data[line] != '\n' {
				line++
			}
			if line >= len(data) {
				continue
			}
			fmt.Sscanf(string(data[line:]), "%s %f %f %f", &fn, &lat, &lon, &h)
			result[y*xCount+x] = h
		}
	}

	return result, nil
}

func init() {
	htdpAvailable()
}
