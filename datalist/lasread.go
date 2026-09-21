package datalist

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flywave/flywave-pointcloud/pdal"
	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/waffle"
	"github.com/flywave/go3d/float64/vec2"
)

const lasProgressBatch = 100000

func ReadLASPoints(path string) ([]waffle.Point, error) {
	pts, _, err := readLASFile(path, nil, nil)
	return pts, err
}

func ReadLASPointsWithSRS(path string) ([]waffle.Point, string, error) {
	return readLASFile(path, nil, nil)
}

func ReadLASPointsWithProgress(path string, progress dem.ProgressFunc, ctx context.Context) ([]waffle.Point, error) {
	pts, _, err := readLASFile(path, progress, ctx)
	return pts, err
}

func readLASFile(path string, progress dem.ProgressFunc, ctx context.Context) ([]waffle.Point, string, error) {
	if path == "" {
		return nil, "", fmt.Errorf("las path is empty")
	}
	if _, err := os.Stat(path); err != nil {
		return nil, "", fmt.Errorf("open las %s: %v", path, err)
	}

	pdal.PdalInitStage()

	stage := pdal.NewStage("readers.las")
	if stage == nil {
		return nil, "", fmt.Errorf("create las reader stage failed")
	}
	defer stage.Free()

	opts := pdal.NewOptions()
	defer opts.Free()
	opts.Add("filename", path)
	stage.SetOptions(opts)

	table := pdal.NewPointTable()
	defer table.Free()
	if err := stage.Prepare(table); err != nil {
		return nil, "", fmt.Errorf("prepare las reader %s: %v", path, err)
	}
	viewSet := stage.Execute(table)
	if viewSet == nil {
		return nil, "", fmt.Errorf("read las %s: %v", path, pdal.LastError())
	}
	defer viewSet.Free()

	total := 0
	cit := viewSet.Iterator()
	for {
		has, view := cit.Next()
		if !has {
			break
		}
		if view != nil {
			total += int(view.Size())
		}
	}
	cit.Free()
	if total == 0 {
		return nil, "", fmt.Errorf("no points in las file: %s", path)
	}

	pts := make([]waffle.Point, 0, total)
	srsWKT := ""
	done := 0

	it := viewSet.Iterator()
	defer it.Free()
	for {
		if ctx != nil {
			if err := ctx.Err(); err != nil {
				return nil, "", fmt.Errorf("read las %s canceled: %v", path, err)
			}
		}
		has, view := it.Next()
		if !has {
			break
		}
		if view == nil {
			continue
		}
		n := int(view.Size())
		xs := make([]float64, n)
		ys := make([]float64, n)
		zs := make([]float64, n)
		view.GetFieldValues(pdal.X, view.DimType(int(pdal.X)), xs)
		view.GetFieldValues(pdal.Y, view.DimType(int(pdal.Y)), ys)
		view.GetFieldValues(pdal.Z, view.DimType(int(pdal.Z)), zs)
		if srsWKT == "" {
			if srs := view.GetSpatialReference(); srs != nil {
				if !srs.Empty() {
					srsWKT = srs.GetWkt()
				}
				srs.Free()
			}
		}
		for start := 0; start < n; start += lasProgressBatch {
			end := start + lasProgressBatch
			if end > n {
				end = n
			}
			for i := start; i < end; i++ {
				pts = append(pts, waffle.Point{
					Position: vec2.T{xs[i], ys[i]},
					Z:        zs[i],
				})
			}
			done += end - start
			if progress != nil {
				progress("las_read", done, total)
			}
		}
		view.Free()
	}

	if srsWKT == "" {
		if srs := stage.GetSpatialReference(); srs != nil {
			if !srs.Empty() {
				srsWKT = srs.GetWkt()
			}
			srs.Free()
		}
	}

	return pts, srsWKT, nil
}

func PointsFromDataList(dl *DataList) ([]waffle.Point, error) {
	if dl == nil || len(dl.Entries) == 0 {
		return nil, fmt.Errorf("datalist has no entries")
	}
	var all []waffle.Point
	for i := range dl.Entries {
		pts, err := pointsFromEntry(&dl.Entries[i])
		if err != nil {
			return nil, err
		}
		all = append(all, pts...)
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("no valid points found in datalist")
	}
	return all, nil
}

func pointsFromEntry(e *DataEntry) ([]waffle.Point, error) {
	switch e.Type {
	case SourceRaster:
		return waffle.PointsFromRasterWithNoData(e.Path, dem.DefaultNoData)
	case SourcePoint:
		switch strings.ToLower(filepath.Ext(e.Path)) {
		case ".las", ".laz":
			return ReadLASPoints(e.Path)
		default:
			xf, err := ParseXYZFile(e.Path, e.SRS)
			if err != nil {
				return nil, err
			}
			pts := make([]waffle.Point, len(xf.Points))
			for i, p := range xf.Points {
				pts[i] = waffle.Point{Position: vec2.T{p.X, p.Y}, Z: p.Z}
			}
			return pts, nil
		}
	default:
		return nil, fmt.Errorf("unsupported datalist source type %s: %s", e.Type, e.Path)
	}
}
