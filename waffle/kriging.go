package waffle

import (
	"fmt"
	"math"
	"sort"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-kriging"
	vec3d "github.com/flywave/go3d/float64/vec3"
)

type krigingWaffle struct {
	baseWaffle
	modelType kriging.ModelType
}

const maxKrigingPoints = 5000

func init() {
	Register(dem.MethodKriging, func() Waffle {
		return &krigingWaffle{baseWaffle: baseWaffle{name: string(dem.MethodKriging)}}
	})
}

func (w *krigingWaffle) Run(points []Point, opts *Options) (*Result, error) {
	if len(points) == 0 {
		return nil, fmt.Errorf("no data points")
	}
	if opts == nil || opts.Region == nil {
		return nil, fmt.Errorf("region is required")
	}

	region := opts.Region
	if region.XSize <= 0 || region.YSize <= 0 {
		region.XSize = int(math.Round((region.BBox().Max[0] - region.BBox().Min[0]) / region.XRes))
		region.YSize = int(math.Round((region.BBox().Max[1] - region.BBox().Min[1]) / region.YRes))
	}

	pos := make([]vec3d.T, len(points))
	for i, p := range points {
		pos[i] = vec3d.T{p.Position[0], p.Position[1], p.Z}
	}
	pos = decimateForKriging(pos, maxKrigingPoints)

	model := kriging.New(pos)
	modelType := w.modelType
	if modelType == "" {
		modelType = kriging.Gaussian
	}
	dem.ReportProgress(opts.Progress, w.Name(), 0, 1)
	if err := dem.CheckCtx(opts.Ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", w.Name(), err)
	}
	_, err := model.Train(modelType, 0, 100)
	if err != nil {
		return nil, fmt.Errorf("kriging training failed: %v", err)
	}
	dem.ReportProgress(opts.Progress, w.Name(), 1, 1)

	demData := make([]float64, region.XSize*region.YSize)
	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	for y := 0; y < region.YSize; y++ {
		for x := 0; x < region.XSize; x++ {
			geoX, geoY := region.PixelCenterGeo(x, y)

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

func decimateForKriging(pos []vec3d.T, maxPoints int) []vec3d.T {
	if len(pos) <= maxPoints {
		return pos
	}

	minX, minY := pos[0][0], pos[0][1]
	maxX, maxY := minX, minY
	for _, p := range pos {
		if p[0] < minX {
			minX = p[0]
		}
		if p[0] > maxX {
			maxX = p[0]
		}
		if p[1] < minY {
			minY = p[1]
		}
		if p[1] > maxY {
			maxY = p[1]
		}
	}

	cells := int(math.Ceil(math.Sqrt(float64(maxPoints))))
	span := math.Max(maxX-minX, maxY-minY)
	cellSize := span / float64(cells)
	if cellSize <= 0 {
		return pos[:1]
	}

	type acc struct {
		x, y, z float64
		n       int
	}
	accs := make(map[[2]int]*acc)
	for _, p := range pos {
		key := [2]int{int((p[0] - minX) / cellSize), int((p[1] - minY) / cellSize)}
		a := accs[key]
		if a == nil {
			a = &acc{}
			accs[key] = a
		}
		a.x += p[0]
		a.y += p[1]
		a.z += p[2]
		a.n++
	}

	keys := make([][2]int, 0, len(accs))
	for k := range accs {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i][0] != keys[j][0] {
			return keys[i][0] < keys[j][0]
		}
		return keys[i][1] < keys[j][1]
	})

	out := make([]vec3d.T, 0, len(keys))
	for _, k := range keys {
		a := accs[k]
		n := float64(a.n)
		out = append(out, vec3d.T{a.x / n, a.y / n, a.z / n})
	}
	return out
}
