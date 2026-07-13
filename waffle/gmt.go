package waffle

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/dem/gmt"
)

type gmtSurfaceWaffle struct{ baseWaffle }
type gmtTriangulateWaffle struct{ baseWaffle }
type gmtNearneighborWaffle struct{ baseWaffle }

func init() {
	Register("gmt_surface", func() Waffle {
		return &gmtSurfaceWaffle{baseWaffle: baseWaffle{name: "gmt_surface"}}
	})
	Register("gmt_triangulate", func() Waffle {
		return &gmtTriangulateWaffle{baseWaffle: baseWaffle{name: "gmt_triangulate"}}
	})
	Register("gmt_nearneighbor", func() Waffle {
		return &gmtNearneighborWaffle{baseWaffle: baseWaffle{name: "gmt_nearneighbor"}}
	})
}

func writeXYZ(points []Point, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, p := range points {
		if _, err := fmt.Fprintf(f, "%.10f %.10f %.6f\n", p.Position[0], p.Position[1], p.Z); err != nil {
			return err
		}
	}
	return nil
}

func readGRD(path string) ([]float64, *dem.Region, error) {
	return dem.ReadDEM(path)
}

func (w *gmtSurfaceWaffle) Run(points []Point, opts *Options) (*Result, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("need at least 3 points")
	}
	region := opts.Region
	if region.XSize <= 0 || region.YSize <= 0 {
		region.XSize = int(math.Round((region.BBox().Max[0] - region.BBox().Min[0]) / region.XRes))
		region.YSize = int(math.Round((region.BBox().Max[1] - region.BBox().Min[1]) / region.YRes))
	}

	tmpDir, err := os.MkdirTemp("", "gmt_surface_*")
	if err != nil {
		return nil, fmt.Errorf("temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	xyzPath := filepath.Join(tmpDir, "input.xyz")
	grdPath := filepath.Join(tmpDir, "output.grd")

	if err := writeXYZ(points, xyzPath); err != nil {
		return nil, fmt.Errorf("write xyz: %v", err)
	}

	tension := 0.25
	cfg := &gmt.GridConfig{
		XInc: region.XRes, YInc: region.YRes,
		XMin: region.BBox().Min[0], XMax: region.BBox().Max[0],
		YMin: region.BBox().Min[1], YMax: region.BBox().Max[1],
		Tension: tension,
	}
	if err := gmt.Surface(xyzPath, grdPath, cfg); err != nil {
		return nil, fmt.Errorf("gmt surface: %v", err)
	}

	data, _, err := readGRD(grdPath)
	if err != nil {
		return nil, fmt.Errorf("read grd: %v", err)
	}

	return &Result{DEM: data, Region: region}, nil
}

func (w *gmtTriangulateWaffle) Run(points []Point, opts *Options) (*Result, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("need at least 3 points")
	}
	region := opts.Region
	if region.XSize <= 0 || region.YSize <= 0 {
		region.XSize = int(math.Round((region.BBox().Max[0] - region.BBox().Min[0]) / region.XRes))
		region.YSize = int(math.Round((region.BBox().Max[1] - region.BBox().Min[1]) / region.YRes))
	}

	tmpDir, err := os.MkdirTemp("", "gmt_tri_*")
	if err != nil {
		return nil, fmt.Errorf("temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	xyzPath := filepath.Join(tmpDir, "input.xyz")
	grdPath := filepath.Join(tmpDir, "output.grd")

	if err := writeXYZ(points, xyzPath); err != nil {
		return nil, fmt.Errorf("write xyz: %v", err)
	}

	cfg := &gmt.GridConfig{
		XInc: region.XRes, YInc: region.YRes,
		XMin: region.BBox().Min[0], XMax: region.BBox().Max[0],
		YMin: region.BBox().Min[1], YMax: region.BBox().Max[1],
	}
	if err := gmt.Triangulate(xyzPath, grdPath, cfg); err != nil {
		return nil, fmt.Errorf("gmt triangulate: %v", err)
	}

	data, _, err := readGRD(grdPath)
	if err != nil {
		return nil, fmt.Errorf("read grd: %v", err)
	}

	return &Result{DEM: data, Region: region}, nil
}

func (w *gmtNearneighborWaffle) Run(points []Point, opts *Options) (*Result, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("need at least 3 points")
	}
	region := opts.Region
	if region.XSize <= 0 || region.YSize <= 0 {
		region.XSize = int(math.Round((region.BBox().Max[0] - region.BBox().Min[0]) / region.XRes))
		region.YSize = int(math.Round((region.BBox().Max[1] - region.BBox().Min[1]) / region.YRes))
	}

	tmpDir, err := os.MkdirTemp("", "gmt_nn_*")
	if err != nil {
		return nil, fmt.Errorf("temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	xyzPath := filepath.Join(tmpDir, "input.xyz")
	grdPath := filepath.Join(tmpDir, "output.grd")

	if err := writeXYZ(points, xyzPath); err != nil {
		return nil, fmt.Errorf("write xyz: %v", err)
	}

	sr := opts.SearchRadius
	if sr <= 0 {
		sr = region.XRes * 5
	}

	cfg := &gmt.GridConfig{
		XInc: region.XRes, YInc: region.YRes,
		XMin: region.BBox().Min[0], XMax: region.BBox().Max[0],
		YMin: region.BBox().Min[1], YMax: region.BBox().Max[1],
		SearchRadius: sr,
		EmptyValue:   -9999,
	}
	if err := gmt.Nearneighbor(xyzPath, grdPath, cfg); err != nil {
		return nil, fmt.Errorf("gmt nearneighbor: %v", err)
	}

	data, _, err := readGRD(grdPath)
	if err != nil {
		return nil, fmt.Errorf("read grd: %v", err)
	}

	return &Result{DEM: data, Region: region}, nil
}
