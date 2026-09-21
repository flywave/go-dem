package datalist

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/flywave/flywave-gdal"
	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
)

type DataSourceType string

const (
	SourceRaster DataSourceType = "raster"
	SourcePoint  DataSourceType = "pointcloud"
	SourceVector DataSourceType = "vector"
)

type DataEntry struct {
	Path        string
	Type        DataSourceType
	Weight      float64
	Uncertainty float64
	SRS         geo.Proj
	Priority    int
}

type DataList struct {
	Entries []DataEntry
	Region  *dem.Region
}

type Stack struct {
	Elevation   []float64
	Count       []float64
	Weight      []float64
	Uncertainty []float64
	SourceID    []float64
	Region      *dem.Region
	NoData      float64
}

func NewStack(region *dem.Region, noData float64) *Stack {
	size := region.XSize * region.YSize
	s := &Stack{
		Elevation:   make([]float64, size),
		Count:       make([]float64, size),
		Weight:      make([]float64, size),
		Uncertainty: make([]float64, size),
		SourceID:    make([]float64, size),
		Region:      region,
		NoData:      noData,
	}
	for i := range s.Elevation {
		s.Elevation[i] = noData
		s.Weight[i] = 0
		s.Uncertainty[i] = 0
		s.SourceID[i] = 0
	}
	return s
}

type StackMode string

const (
	StackModeMean   StackMode = "mean"
	StackModeMin    StackMode = "min"
	StackModeMax    StackMode = "max"
	StackModeFirst  StackMode = "first"
	StackModeLast   StackMode = "last"
	StackModeWeight StackMode = "weight"
)

func isNilProj(p geo.Proj) bool {
	if p == nil {
		return true
	}
	sp, ok := p.(*geo.SRSProj4)
	return ok && sp == nil
}

func sameSRS(a, b geo.Proj) bool {
	if isNilProj(a) || isNilProj(b) {
		return isNilProj(a) && isNilProj(b)
	}
	return a.Eq(b)
}

func (s *Stack) takeFrom(other *Stack, i int) {
	s.Elevation[i] = other.Elevation[i]
	s.Count[i] = other.Count[i]
	s.Weight[i] = other.Weight[i]
	s.Uncertainty[i] = other.Uncertainty[i]
	s.SourceID[i] = other.SourceID[i]
}

func (s *Stack) Merge(other *Stack, mode StackMode) error {
	if s.Region.XSize != other.Region.XSize || s.Region.YSize != other.Region.YSize {
		return fmt.Errorf("stack size mismatch: %dx%d vs %dx%d",
			s.Region.XSize, s.Region.YSize,
			other.Region.XSize, other.Region.YSize)
	}
	if s.Region.BBox() != other.Region.BBox() {
		return fmt.Errorf("stack bounds mismatch: %v/%v vs %v/%v",
			s.Region.BBox().Min, s.Region.BBox().Max,
			other.Region.BBox().Min, other.Region.BBox().Max)
	}
	if !sameSRS(s.Region.SRS(), other.Region.SRS()) {
		return fmt.Errorf("stack SRS mismatch")
	}

	switch mode {
	case StackModeMean, StackModeMin, StackModeMax, StackModeFirst, StackModeLast, StackModeWeight:
	default:
		return fmt.Errorf("unknown stack mode: %s", mode)
	}

	for i := range s.Elevation {
		otherVal := other.Elevation[i]
		if otherVal == other.NoData || math.IsNaN(otherVal) {
			continue
		}

		currentVal := s.Elevation[i]
		if currentVal == s.NoData || math.IsNaN(currentVal) {
			s.takeFrom(other, i)
			continue
		}

		switch mode {
		case StackModeMean:
			totalCount := s.Count[i] + other.Count[i]
			if totalCount > 0 {
				s.Elevation[i] = (currentVal*s.Count[i] + otherVal*other.Count[i]) / totalCount
				s.Count[i] = totalCount
			}
			s.Uncertainty[i] = math.Sqrt(s.Uncertainty[i]*s.Uncertainty[i] + other.Uncertainty[i]*other.Uncertainty[i])
		case StackModeMin:
			if otherVal < currentVal {
				s.takeFrom(other, i)
			}
		case StackModeMax:
			if otherVal > currentVal {
				s.takeFrom(other, i)
			}
		case StackModeFirst:
		case StackModeLast:
			s.takeFrom(other, i)
		case StackModeWeight:
			totalWeight := s.Weight[i] + other.Weight[i]
			if totalWeight > 0 {
				s.Elevation[i] = (currentVal*s.Weight[i] + otherVal*other.Weight[i]) / totalWeight
				s.Weight[i] = totalWeight
			}
			s.Uncertainty[i] = math.Sqrt(s.Uncertainty[i]*s.Uncertainty[i] + other.Uncertainty[i]*other.Uncertainty[i])
		}
	}

	return nil
}

func (s *Stack) ToBands() [][]float64 {
	return [][]float64{
		s.Elevation,
		s.Count,
		s.Weight,
		s.Uncertainty,
		s.SourceID,
	}
}

func (s *Stack) Write(outputPath string) error {
	bandData := s.ToBands()
	return dem.CreateStack(bandData, s.Region, outputPath, s.NoData)
}

func ReadStack(path string) (*Stack, error) {
	data, region, err := dem.ReadDEM(path)
	if err != nil {
		return nil, err
	}

	size := region.XSize * region.YSize
	stack := &Stack{
		Elevation:   data,
		Count:       make([]float64, size),
		Weight:      make([]float64, size),
		Uncertainty: make([]float64, size),
		SourceID:    make([]float64, size),
		Region:      region,
		NoData:      dem.DefaultNoData,
	}
	if noData, ok := readFileNoData(path); ok {
		stack.NoData = noData
	}

	var bandErrs []error
	for _, band := range []struct {
		idx int
		dst *[]float64
	}{
		{2, &stack.Count},
		{3, &stack.Weight},
		{4, &stack.Uncertainty},
		{5, &stack.SourceID},
	} {
		bandData, _, err := dem.ReadDEMBand(path, band.idx)
		if err != nil {
			bandErrs = append(bandErrs, fmt.Errorf("band %d: %v", band.idx, err))
			continue
		}
		*band.dst = bandData
	}
	if len(bandErrs) > 0 {
		return nil, errors.Join(bandErrs...)
	}

	return stack, nil
}

func readFileNoData(path string) (float64, bool) {
	var noData float64
	valid := false
	err := gdal.WithDatasetReadonly(path, func(ds gdal.Dataset) error {
		v, ok := ds.RasterBand(1).NoDataValue()
		if ok {
			noData = v
			valid = true
		}
		return nil
	})
	if err != nil || !valid {
		return dem.DefaultNoData, false
	}
	return noData, true
}

func BuildDataList(paths []string) (*DataList, error) {
	dl := &DataList{}
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("stat %s: %v", p, err)
		}
		if info.IsDir() {
			entries, err := os.ReadDir(p)
			if err != nil {
				return nil, fmt.Errorf("read dir %s: %v", p, err)
			}
			for _, entry := range entries {
				if entry.IsDir() || !isKnownSource(entry.Name()) {
					continue
				}
				entryPath := filepath.Join(p, entry.Name())
				dl.Entries = append(dl.Entries, DataEntry{
					Path: entryPath,
					Type: detectType(entryPath),
				})
			}
		} else {
			dl.Entries = append(dl.Entries, DataEntry{
				Path: p,
				Type: detectType(p),
			})
		}
	}
	return dl, nil
}

var sourceTypes = map[string]DataSourceType{
	".tif":     SourceRaster,
	".tiff":    SourceRaster,
	".img":     SourceRaster,
	".asc":     SourceRaster,
	".hgt":     SourceRaster,
	".las":     SourcePoint,
	".laz":     SourcePoint,
	".xyz":     SourcePoint,
	".csv":     SourcePoint,
	".txt":     SourcePoint,
	".shp":     SourceVector,
	".geojson": SourceVector,
	".json":    SourceVector,
	".gpkg":    SourceVector,
}

func detectType(path string) DataSourceType {
	if t, ok := sourceTypes[strings.ToLower(filepath.Ext(path))]; ok {
		return t
	}
	return SourceRaster
}

func isKnownSource(path string) bool {
	_, ok := sourceTypes[strings.ToLower(filepath.Ext(path))]
	return ok
}
