package waffle

import (
	"fmt"
	"math"

	"github.com/flywave/flywave-gdal"
	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

type idwWaffle struct {
	baseWaffle
	power float64
}

func NewIDW(power float64) *idwWaffle {
	return &idwWaffle{baseWaffle: baseWaffle{name: string(dem.MethodIDW)}, power: power}
}

func init() {
	Register(dem.MethodIDW, func() Waffle {
		return &idwWaffle{baseWaffle: baseWaffle{name: string(dem.MethodIDW)}, power: 2.0}
	})
}

func (w *idwWaffle) Run(sources []string, opts *Options) (*Result, error) {
	if len(sources) == 0 {
		return nil, fmt.Errorf("no source data provided")
	}
	region := opts.Region
	if region == nil {
		return nil, fmt.Errorf("region is required")
	}

	if len(sources) == 1 && isLASFile(sources[0]) {
		return w.runWithPointsToGrid(sources[0], region, opts)
	}
	return w.runMemoryIDW(sources, region, opts)
}

func (w *idwWaffle) runWithPointsToGrid(source string, region *dem.Region, opts *Options) (*Result, error) {
	bbox := [4]float64{
		region.BBox().Min[0], region.BBox().Min[1],
		region.BBox().Max[0], region.BBox().Max[1],
	}

	gridOpts := &gdal.PointsToGridOptions{
		InputFilePath:  source,
		OutputFilePath: fmt.Sprintf("%s_idw.tif", source),
		Resolution:     region.XRes,
		EpsgCode:       4326,
		Knn:            opts.MinPoints,
		BBox:           &bbox,
		BBoxEPSGCode:   4326,
	}
	if opts.UpperZ != nil {
		gridOpts.MaxZ = opts.UpperZ
	}
	if opts.LowerZ != nil {
		gridOpts.MinZ = opts.LowerZ
	}
	if opts.MinPoints <= 0 {
		gridOpts.Knn = 3
	}

	err := gdal.PointsToGrid(gridOpts)
	if err != nil {
		return nil, fmt.Errorf("points to grid failed: %v", err)
	}

	data, _, err := dem.ReadDEM(gridOpts.OutputFilePath)
	if err != nil {
		return nil, fmt.Errorf("read grid result: %v", err)
	}
	return &Result{DEM: data, Region: region}, nil
}

func (w *idwWaffle) runMemoryIDW(sources []string, region *dem.Region, opts *Options) (*Result, error) {
	pts, zs, err := collectPoints(sources)
	if err != nil {
		return nil, err
	}
	if len(pts) == 0 {
		return nil, fmt.Errorf("no valid data points found in sources")
	}

	if region.XSize <= 0 || region.YSize <= 0 {
		region.XSize = int(math.Round((region.BBox().Max[0] - region.BBox().Min[0]) / region.XRes))
		region.YSize = int(math.Round((region.BBox().Max[1] - region.BBox().Min[1]) / region.YRes))
	}

	demData := make([]float64, region.XSize*region.YSize)
	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}
	for i := range demData {
		demData[i] = noData
	}

	gt := region.GeoTransform()
	searchRadius := opts.SearchRadius
	if searchRadius <= 0 {
		searchRadius = region.XRes * 10
	}
	power := w.power
	minPoints := opts.MinPoints
	if minPoints <= 0 {
		minPoints = 3
	}

	for y := 0; y < region.YSize; y++ {
		for x := 0; x < region.XSize; x++ {
			geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
			geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]

			var sumWeight, sumValue float64
			count := 0

			for i, pt := range pts {
				dx := geoX - pt[0]
				dy := geoY - pt[1]
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > searchRadius {
					continue
				}
				if dist < 1e-10 {
					sumValue = zs[i]
					sumWeight = 1
					count = 1
					break
				}
				weight := 1.0 / math.Pow(dist, power)
				sumWeight += weight
				sumValue += weight * zs[i]
				count++
			}

			if count >= minPoints && sumWeight > 0 {
				demData[y*region.XSize+x] = sumValue / sumWeight
			}
		}
	}

	return &Result{DEM: demData, Region: region}, nil
}

func collectPoints(sources []string) ([]vec2.T, []float64, error) {
	var pts []vec2.T
	var zs []float64

	for _, src := range sources {
		data, srcRegion, err := dem.ReadDEM(src)
		if err != nil {
			return nil, nil, fmt.Errorf("read %s: %v", src, err)
		}
		gt := srcRegion.GeoTransform()
		noData := dem.DefaultNoData
		for y := 0; y < srcRegion.YSize; y++ {
			for x := 0; x < srcRegion.XSize; x++ {
				z := data[y*srcRegion.XSize+x]
				if z == noData || math.IsNaN(z) {
					continue
				}
				geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
				geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]
				pts = append(pts, vec2.T{geoX, geoY})
				zs = append(zs, z)
			}
		}
	}
	return pts, zs, nil
}

func isLASFile(path string) bool {
	ext := path[len(path)-4:]
	return ext == ".las" || ext == ".laz"
}
