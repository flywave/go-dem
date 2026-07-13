package gmt

import (
	"os"
	"path/filepath"
	"runtime"
)

var gmtDataDir string

// SetGMTDataDir sets the GMT share data directory path.
// GMT modules need access to this directory for coastlines,
// CPT files, and other runtime data.  Call before any GMT operation.
func SetGMTDataDir(path string) {
	gmtDataDir = path
	if path != "" {
		os.Setenv("GMT_SHAREDIR", path)
	}
}

// GMTDataDir returns the current GMT share data directory.
func GMTDataDir() string {
	return gmtDataDir
}

func init() {
	// Priority:
	// 1. GMT_DATA environment variable
	// 2. GMT_SHAREDIR environment variable (standard GMT)
	// 3. Default: gmt_data/ relative to this source file
	dir := os.Getenv("GMT_DATA")
	if dir != "" {
		SetGMTDataDir(dir)
		return
	}
	dir = os.Getenv("GMT_SHAREDIR")
	if dir != "" {
		SetGMTDataDir(dir)
		return
	}
	dir = defaultGMTDataDir()
	if dir != "" {
		SetGMTDataDir(dir)
	}
}

func defaultGMTDataDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	// dem/gmt/gmt_data.go -> ../../gmt_data
	base := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	candidate := filepath.Join(base, "gmt_data")
	if info, err := os.Stat(candidate); err == nil && info.IsDir() {
		return candidate
	}
	return ""
}
