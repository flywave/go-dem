package waffle

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

type movingAverageWaffle struct {
	baseWaffle
}

func init() {
	Register(dem.MethodMovingAverage, func() Waffle {
		return &movingAverageWaffle{baseWaffle: baseWaffle{name: string(dem.MethodMovingAverage)}}
	})
}

func (w *movingAverageWaffle) Run(points []Point, opts *Options) (*Result, error) {
	if len(points) == 0 {
		return nil, fmt.Errorf("no source points provided")
	}
	region := opts.Region
	if region == nil {
		return nil, fmt.Errorf("region is required")
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

	pts := make([]vec2.T, len(points))
	zs := make([]float64, len(points))
	for i, p := range points {
		pts[i] = p.Position
		zs[i] = p.Z
	}

	kdtree := NewKDTree(pts)

	gt := region.GeoTransform()
	searchRadius := opts.SearchRadius
	if searchRadius <= 0 {
		searchRadius = region.XRes * 10
	}
	minPoints := opts.MinPoints
	if minPoints <= 0 {
		minPoints = 3
	}

	for y := 0; y < region.YSize; y++ {
		for x := 0; x < region.XSize; x++ {
			geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
			geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]

			q := vec2.T{geoX, geoY}
			idxs, _ := kdtree.RadiusSearch(q, searchRadius)

			if len(idxs) < minPoints {
				idxs2, _ := kdtree.KNN(q, minPoints)
				if len(idxs2) >= minPoints {
					demData[y*region.XSize+x] = meanAverage(idxs2, zs)
				}
				continue
			}

			demData[y*region.XSize+x] = meanAverage(idxs, zs)
		}
	}

	return &Result{DEM: demData, Region: region}, nil
}

func meanAverage(idxs []int, zs []float64) float64 {
	var sum float64
	for _, idx := range idxs {
		sum += zs[idx]
	}
	return sum / float64(len(idxs))
}
