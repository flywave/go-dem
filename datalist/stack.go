package datalist

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
)

type DataSourceType string

const (
	SourceRaster  DataSourceType = "raster"
	SourcePoint   DataSourceType = "pointcloud"
	SourceVector  DataSourceType = "vector"
)

type DataEntry struct {
	Path     string
	Type     DataSourceType
	Weight   float64
	Uncertainty float64
	SRS      geo.Proj
	Priority int
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

func (s *Stack) Merge(other *Stack, mode StackMode) error {
	if s.Region.XSize != other.Region.XSize || s.Region.YSize != other.Region.YSize {
		return fmt.Errorf("stack size mismatch: %dx%d vs %dx%d",
			s.Region.XSize, s.Region.YSize,
			other.Region.XSize, other.Region.YSize)
	}

	for i := range s.Elevation {
		otherVal := other.Elevation[i]
		if otherVal == other.NoData || math.IsNaN(otherVal) {
			continue
		}

		currentVal := s.Elevation[i]
		if currentVal == s.NoData || math.IsNaN(currentVal) {
			s.Elevation[i] = otherVal
			s.Count[i] = other.Count[i]
			s.Weight[i] = other.Weight[i]
			s.Uncertainty[i] = other.Uncertainty[i]
			s.SourceID[i] = other.SourceID[i]
			continue
		}

		switch mode {
		case StackModeMean:
			totalCount := s.Count[i] + other.Count[i]
			if totalCount > 0 {
				s.Elevation[i] = (currentVal*s.Count[i] + otherVal*other.Count[i]) / totalCount
				s.Count[i] = totalCount
			}
		case StackModeMin:
			if otherVal < currentVal {
				s.Elevation[i] = otherVal
			}
		case StackModeMax:
			if otherVal > currentVal {
				s.Elevation[i] = otherVal
			}
		case StackModeWeight:
			totalWeight := s.Weight[i] + other.Weight[i]
			if totalWeight > 0 {
				s.Elevation[i] = (currentVal*s.Weight[i] + otherVal*other.Weight[i]) / totalWeight
				s.Weight[i] = totalWeight
			}
		}

		s.Uncertainty[i] = math.Sqrt(s.Uncertainty[i]*s.Uncertainty[i] + other.Uncertainty[i]*other.Uncertainty[i])
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
	size := 1
	region := &dem.Region{XSize: 0, YSize: 0}
	data, reg, err := dem.ReadDEM(path)
	if err != nil {
		return nil, err
	}
	region = reg
	size = region.XSize * region.YSize

	stack := &Stack{
		Elevation:   data,
		Count:       make([]float64, size),
		Weight:      make([]float64, size),
		Uncertainty: make([]float64, size),
		SourceID:    make([]float64, size),
		Region:      region,
		NoData:      dem.DefaultNoData,
	}

	if band2, _, err := dem.ReadDEMBand(path, 2); err == nil {
		stack.Count = band2
	}
	if band3, _, err := dem.ReadDEMBand(path, 3); err == nil {
		stack.Weight = band3
	}
	if band4, _, err := dem.ReadDEMBand(path, 4); err == nil {
		stack.Uncertainty = band4
	}
	if band5, _, err := dem.ReadDEMBand(path, 5); err == nil {
		stack.SourceID = band5
	}

	return stack, nil
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

func detectType(path string) DataSourceType {
	ext := filepath.Ext(path)
	switch ext {
	case ".tif", ".tiff", ".img", ".asc", ".hgt":
		return SourceRaster
	case ".las", ".laz", ".xyz", ".csv", ".txt":
		return SourcePoint
	case ".shp", ".geojson", ".json", ".gpkg":
		return SourceVector
	default:
		return SourceRaster
	}
}

func validateDataList(dl *DataList) error {
	if len(dl.Entries) == 0 {
		return fmt.Errorf("no valid data entries")
	}
	for _, entry := range dl.Entries {
		if _, err := os.Stat(entry.Path); os.IsNotExist(err) {
			return fmt.Errorf("data source not found: %s", entry.Path)
		}
	}
	return nil
}
