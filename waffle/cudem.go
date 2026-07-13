package waffle

import (
	"fmt"
	"math"
	"sort"

	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

type cudemWaffle struct {
	baseWaffle
}

func init() {
	Register(dem.MethodCUDEM, func() Waffle {
		return &cudemWaffle{baseWaffle: baseWaffle{name: string(dem.MethodCUDEM)}}
	})
}

type cudemLevel struct {
	Scale      float64
	Resolution float64
}

func (cw *cudemWaffle) Run(sources []string, opts *Options) (*Result, error) {
	pts, zs, err := collectPoints(sources)
	if err != nil {
		return nil, err
	}
	if len(pts) == 0 {
		return nil, fmt.Errorf("no data points")
	}

	region := opts.Region
	if region.XSize <= 0 || region.YSize <= 0 {
		region.XSize = int(math.Round((region.BBox().Max[0] - region.BBox().Min[0]) / region.XRes))
		region.YSize = int(math.Round((region.BBox().Max[1] - region.BBox().Min[1]) / region.YRes))
	}

	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	levels := computeStepDownLevels(pts, region.XRes)

	gt := region.GeoTransform()
	w, h := region.XSize, region.YSize

	demData := make([]float64, w*h)
	for i := range demData {
		demData[i] = noData
	}

	result := make([]float64, w*h)
	count := make([]int, w*h)

	sort.Slice(levels, func(i, j int) bool {
		return levels[i].Resolution > levels[j].Resolution
	})

	for li, level := range levels {

		scale := level.Scale

		gridW := int(math.Ceil(float64(w) / scale))
		gridH := int(math.Ceil(float64(h) / scale))

		gridData := make([]float64, gridW*gridH)
		gridCount := make([]int, gridW*gridH)
		for i := range gridData {
			gridData[i] = noData
		}

		res := level.Resolution
		searchRadius := res * 3

		for i, pt := range pts {
			gx := int(math.Floor((pt[0] - gt[0]) / res))
			gy := int(math.Floor((pt[1] - gt[3]) / res))
			if gx < 0 || gx >= gridW || gy < 0 || gy >= gridH {
				continue
			}
			idx := gy*gridW + gx
			if gridData[idx] == noData || math.IsNaN(gridData[idx]) {
				gridData[idx] = zs[i]
				gridCount[idx] = 1
			} else {
				gridData[idx] = (gridData[idx]*float64(gridCount[idx]) + zs[i]) / float64(gridCount[idx]+1)
				gridCount[idx]++
			}
		}

		filledGrid := make([]float64, gridW*gridH)
		copy(filledGrid, gridData)

		for gy := 0; gy < gridH; gy++ {
			for gx := 0; gx < gridW; gx++ {
				idx := gy*gridW + gx
				if filledGrid[idx] != noData && !math.IsNaN(filledGrid[idx]) {
					continue
				}

				var sum, wSum float64
				for dy := -3; dy <= 3; dy++ {
					for dx := -3; dx <= 3; dx++ {
						if dx == 0 && dy == 0 {
							continue
						}
						nx, ny := gx+dx, gy+dy
						if nx < 0 || nx >= gridW || ny < 0 || ny >= gridH {
							continue
						}
						val := gridData[ny*gridW+nx]
						if val == noData || math.IsNaN(val) {
							continue
						}
						dist := math.Sqrt(float64(dx*dx + dy*dy))
						if dist <= 3 {
							weight := 1.0 / (dist*dist + 1e-15)
							sum += val * weight
							wSum += weight
						}
					}
				}
				if wSum > 0 {
					filledGrid[idx] = sum / wSum
				}
			}
		}

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				idx := y*w + x
				if li == 0 {
					geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
					geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]

					gx := int(math.Floor((geoX - gt[0]) / res))
					gy := int(math.Floor((geoY - gt[3]) / res))

					if gx >= 0 && gx < gridW && gy >= 0 && gy < gridH {
						val := filledGrid[gy*gridW+gx]
						if val != noData && !math.IsNaN(val) {
							result[idx] = val
							count[idx] = 1
						}
					}
				}

				if count[idx] > 0 && li > 0 {
					geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
					geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]

					ptsInRadius := findPointsInRadius(pts, geoX, geoY, searchRadius)
					if len(ptsInRadius) >= 3 {
						continue
					}

					gx := int(math.Floor((geoX - gt[0]) / res))
					gy := int(math.Floor((geoY - gt[3]) / res))

					if gx >= 0 && gx < gridW && gy >= 0 && gy < gridH {
						val := filledGrid[gy*gridW+gx]
						if val != noData && !math.IsNaN(val) {
							result[idx] = val
						}
					}
				}
			}
		}
	}

	for i := range result {
		if result[i] == 0 && demData[i] != noData {
			result[i] = demData[i]
		}
		if result[i] == noData || math.IsNaN(result[i]) {
			result[i] = noData
		}
	}

	return &Result{DEM: result, Region: region}, nil
}

func computeStepDownLevels(pts []vec2.T, baseResolution float64) []cudemLevel {
	xMin, xMax := pts[0][0], pts[0][0]
	yMin, yMax := pts[0][1], pts[0][1]
	for _, pt := range pts {
		if pt[0] < xMin {
			xMin = pt[0]
		}
		if pt[0] > xMax {
			xMax = pt[0]
		}
		if pt[1] < yMin {
			yMin = pt[1]
		}
		if pt[1] > yMax {
			yMax = pt[1]
		}
	}
	area := (xMax - xMin) * (yMax - yMin)
	pointDensity := float64(len(pts)) / area

	nLevels := 3
	if pointDensity < 100 {
		nLevels = 4
	}
	if pointDensity < 10 {
		nLevels = 5
	}

	levels := make([]cudemLevel, nLevels)
	for i := 0; i < nLevels; i++ {
		levels[i] = cudemLevel{
			Scale:      math.Pow(2, float64(i)),
			Resolution: baseResolution * math.Pow(2, float64(i)),
		}
	}
	return levels
}

func findPointsInRadius(pts []vec2.T, x, y, radius float64) []vec2.T {
	var result []vec2.T
	radiusSq := radius * radius
	for _, pt := range pts {
		dx := x - pt[0]
		dy := y - pt[1]
		if dx*dx+dy*dy <= radiusSq {
			result = append(result, pt)
		}
	}
	return result
}


