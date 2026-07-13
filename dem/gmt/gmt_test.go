package gmt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGMTBegin(t *testing.T) {
	// gmt_begin() is called in init() — just verify no panic
	t.Log("GMT session initialized")
}

func TestSurfaceSmall(t *testing.T) {
	dir, err := os.MkdirTemp("", "gmt_test_*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	xyzPath := filepath.Join(dir, "test.xyz")
	grdPath := filepath.Join(dir, "test.grd")

	// Create a small XYZ file with known data
	xyz := []string{
		"0.0 0.0 100.0\n",
		"1.0 0.0 200.0\n",
		"0.0 1.0 300.0\n",
		"1.0 1.0 400.0\n",
		"0.5 0.5 250.0\n",
	}
	f, err := os.Create(xyzPath)
	if err != nil {
		t.Fatalf("create xyz: %v", err)
	}
	for _, line := range xyz {
		f.WriteString(line)
	}
	f.Close()

	cfg := &GridConfig{
		XInc: 0.5, YInc: 0.5,
		XMin: 0, XMax: 1,
		YMin: 0, YMax: 1,
		Tension: 0.25,
	}

	err = Surface(xyzPath, grdPath, cfg)
	t.Logf("xyzPath=%s grdPath=%s", xyzPath, grdPath)
	if err != nil {
		t.Fatalf("gmt surface: %v", err)
	}

	info, err := os.Stat(grdPath)
	if err != nil {
		t.Fatalf("output not created: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("output file is empty")
	}
	t.Logf("output size: %d bytes", info.Size())
}

func TestBlockmeanSmall(t *testing.T) {
	dir, err := os.MkdirTemp("", "gmt_test_*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	xyzPath := filepath.Join(dir, "test.xyz")
	outPath := filepath.Join(dir, "blockmean.grd")

	xyz := []string{
		"0.0 0.0 100.0\n",
		"0.2 0.0 102.0\n",
		"0.8 0.0 200.0\n",
		"0.0 0.0 101.0\n",
	}
	f, _ := os.Create(xyzPath)
	for _, line := range xyz {
		f.WriteString(line)
	}
	f.Close()

	cfg := &GridConfig{
		XInc: 0.5, YInc: 0.5,
		XMin: 0, XMax: 1,
		YMin: 0, YMax: 1,
	}

	err = Blockmean(xyzPath, outPath, cfg)
	if err != nil {
		t.Fatalf("gmt blockmean: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("output not created: %v", err)
	}
	t.Logf("output size: %d bytes", info.Size())
}

func TestNearneighborSmall(t *testing.T) {
	dir, err := os.MkdirTemp("", "gmt_test_*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	xyzPath := filepath.Join(dir, "test.xyz")
	grdPath := filepath.Join(dir, "nn.grd")

	xyz := []string{
		"0.0 0.0 100.0\n",
		"1.0 0.0 200.0\n",
		"0.0 1.0 300.0\n",
		"1.0 1.0 400.0\n",
	}
	f, _ := os.Create(xyzPath)
	for _, line := range xyz {
		f.WriteString(line)
	}
	f.Close()

	cfg := &GridConfig{
		XInc: 0.5, YInc: 0.5,
		XMin: 0, XMax: 1,
		YMin: 0, YMax: 1,
		SearchRadius: 0.6,
		EmptyValue:   -9999,
	}

	err = Nearneighbor(xyzPath, grdPath, cfg)
	if err != nil {
		t.Fatalf("gmt nearneighbor: %v", err)
	}

	info, err := os.Stat(grdPath)
	if err != nil {
		t.Fatalf("output not created: %v", err)
	}
	t.Logf("output size: %d bytes", info.Size())
}

func TestTriangulateSmall(t *testing.T) {
	dir, err := os.MkdirTemp("", "gmt_test_*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	xyzPath := filepath.Join(dir, "test.xyz")
	grdPath := filepath.Join(dir, "tri.grd")

	xyz := []string{
		"0.0 0.0 100.0\n",
		"1.0 0.0 200.0\n",
		"0.0 1.0 300.0\n",
		"1.0 1.0 400.0\n",
	}
	f, _ := os.Create(xyzPath)
	for _, line := range xyz {
		f.WriteString(line)
	}
	f.Close()

	cfg := &GridConfig{
		XInc: 0.5, YInc: 0.5,
		XMin: 0, XMax: 1,
		YMin: 0, YMax: 1,
	}

	err = Triangulate(xyzPath, grdPath, cfg)
	if err != nil {
		t.Fatalf("gmt triangulate: %v", err)
	}

	info, err := os.Stat(grdPath)
	if err != nil {
		t.Fatalf("output not created: %v", err)
	}
	t.Logf("output size: %d bytes", info.Size())
}

func TestGrdfilterSmall(t *testing.T) {
	dir, err := os.MkdirTemp("", "gmt_test_*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// First create a surface to filter
	xyzPath := filepath.Join(dir, "test.xyz")
	grdPath := filepath.Join(dir, "input.grd")
	filtPath := filepath.Join(dir, "filtered.grd")

	xyz := []string{
		"0.0 0.0 100.0\n",
		"1.0 0.0 200.0\n",
		"0.0 1.0 300.0\n",
		"1.0 1.0 400.0\n",
		"0.5 0.5 250.0\n",
	}
	f, _ := os.Create(xyzPath)
	for _, line := range xyz {
		f.WriteString(line)
	}
	f.Close()

	cfg := &GridConfig{
		XInc: 0.5, YInc: 0.5,
		XMin: 0, XMax: 1,
		YMin: 0, YMax: 1,
		Tension: 0.25,
	}

	if err := Surface(xyzPath, grdPath, cfg); err != nil {
		t.Fatalf("surface for filter test: %v", err)
	}

	err = Grdfilter(grdPath, filtPath, "c100", "0.5")
	if err != nil {
		t.Fatalf("gmt grdfilter: %v", err)
	}

	info, err := os.Stat(filtPath)
	if err != nil {
		t.Fatalf("filtered output not created: %v", err)
	}
	t.Logf("filtered output size: %d bytes", info.Size())
}
