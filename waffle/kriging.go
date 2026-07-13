package waffle

import (
	"fmt"
	"math"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-kriging"
	vec3d "github.com/flywave/go3d/float64/vec3"
)

type krigingWaffle struct {
	baseWaffle
	modelType kriging.ModelType
}

func init() {
	Register(dem.MethodKriging, func() Waffle {
		return &krigingWaffle{baseWaffle: baseWaffle{name: string(dem.MethodKriging)}}
	})
}

func (w *krigingWaffle) Run(sources []string, opts *Options) (*Result, error) {
	pts, zs, err := collectPoints(sources)
	if err != nil {
		return nil, err
	}
	if len(pts) == 0 {
		return nil, fmt.Errorf("no valid data points found")
	}

	region := opts.Region
	if region.XSize <= 0 || region.YSize <= 0 {
		region.XSize = int(math.Round((region.BBox().Max[0] - region.BBox().Min[0]) / region.XRes))
		region.YSize = int(math.Round((region.BBox().Max[1] - region.BBox().Min[1]) / region.YRes))
	}

	pos := make([]vec3d.T, len(pts))
	for i, pt := range pts {
		pos[i] = vec3d.T{pt[0], pt[1], zs[i]}
	}

	model := kriging.New(pos)
	modelType := w.modelType
	if modelType == "" {
		modelType = kriging.Gaussian
	}
	_, err = model.Train(modelType, 0, 100)
	if err != nil {
		return nil, fmt.Errorf("kriging training failed: %v", err)
	}

	demData := make([]float64, region.XSize*region.YSize)
	gt := region.GeoTransform()
	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	for y := 0; y < region.YSize; y++ {
		for x := 0; x < region.XSize; x++ {
			geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
			geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]

			val := model.Predict(geoX, geoY)
			if math.IsNaN(val) || math.IsInf(val, 0) {
				demData[y*region.XSize+x] = noData
			} else {
				demData[y*region.XSize+x] = val
			}
		}
	}

	return &Result{DEM: demData, Region: region}, nil
}
