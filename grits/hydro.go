package grits

import (
	"math"

	"github.com/flywave/go-dem"
)

type hydroFilter struct{ baseGrits }

func init() {
	Register("hydro", func() Grits { return &hydroFilter{baseGrits{name: "hydro"}} })
}

func (f *hydroFilter) Run(data []float64, region *dem.Region, opts *Options) ([]float64, error) {
	noData := opts.GetNoData()

	result := fillSinks(data, region.XSize, region.YSize, noData)
	return result, nil
}

type hydrologyEdge struct {
	x, y int
	z    float64
}

func fillSinks(data []float64, w, h int, noData float64) []float64 {
	result := make([]float64, len(data))
	copy(result, data)

	for iteration := 0; iteration < 10; iteration++ {
		changed := 0
		for y := 1; y < h-1; y++ {
			for x := 1; x < w-1; x++ {
				idx := y*w + x
				z := result[idx]
				if z == noData || math.IsNaN(z) {
					continue
				}

				minNeighbor := math.MaxFloat64
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						if dx == 0 && dy == 0 {
							continue
						}
						nx, ny := x+dx, y+dy
						nidx := ny*w + nx
						nval := result[nidx]
						if nval != noData && !math.IsNaN(nval) && nval < minNeighbor {
							minNeighbor = nval
						}
					}
				}

				if minNeighbor < math.MaxFloat64 && z < minNeighbor-1e-10 {

					flowTo := -1
					minSlope := math.MaxFloat64
					for dy := -1; dy <= 1; dy++ {
						for dx := -1; dx <= 1; dx++ {
							if dx == 0 && dy == 0 {
								continue
							}
							nx, ny := x+dx, y+dy
							nidx := ny*w + nx
							nval := result[nidx]
							if nval == noData || math.IsNaN(nval) {
								continue
							}
							dist := math.Sqrt(float64(dx*dx + dy*dy))
							slope := (z - nval) / dist
							if slope > minSlope {
								minSlope = slope
								flowTo = idx
							}
						}
					}

					if flowTo >= 0 {
						result[idx] = result[flowTo] + 0.0001
						changed++
					}
				}
			}
		}
		if changed == 0 {
			break
		}
	}

	result = fillFlatAreas(result, w, h, noData)

	return result
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func fillFlatAreas(data []float64, w, h int, noData float64) []float64 {
	result := make([]float64, len(data))
	copy(result, data)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if result[idx] == noData || math.IsNaN(result[idx]) {
				continue
			}

			edge := false
			var borderZ float64
			borderFound := false
			distance := 0

			for d := 1; d <= 10; d++ {
				for dy := -d; dy <= d; dy++ {
					for dx := -d; dx <= d; dx++ {
						if absInt(dx) != d && absInt(dy) != d {
							continue
						}
						nx, ny := x+dx, y+dy
						if nx < 0 || nx >= w || ny < 0 || ny >= h {
							edge = true
							break
						}
						nidx := ny*w + nx
						nval := result[nidx]
						if nval == noData || math.IsNaN(nval) {
							edge = true
							break
						}
						if math.Abs(nval-result[idx]) > 1e-6 {
							borderZ = nval
							borderFound = true
							distance = d
							break
						}
					}
					if edge || borderFound {
						break
					}
				}
				if edge || borderFound {
					break
				}
			}

			if borderFound && distance > 1 {
				gradient := (borderZ - result[idx]) / float64(distance)
				for d := 1; d < distance; d++ {
					for dy := -d; dy <= d; dy++ {
						for dx := -d; dx <= d; dx++ {
							if absInt(dx) != d && absInt(dy) != d {
								continue
							}
							nx, ny := x+dx, y+dy
							if nx < 0 || nx >= w || ny < 0 || ny >= h {
								continue
							}
							nidx := ny*w + nx
							if math.Abs(result[nidx]-result[idx]) < 1e-6 {
								result[nidx] = result[idx] + gradient*float64(d)
							}
						}
					}
				}
			}
		}
	}

	return result
}
