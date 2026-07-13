package uncertainty

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

type Method string

const (
	MethodSplitSample Method = "split_sample"
	MethodProximity   Method = "proximity"
	MethodCombined    Method = "combined"
)

type Options struct {
	Method          Method
	NoData          float64
	SampleFraction  float64
	SearchRadius    float64
	Seed            int64
}

type Result struct {
	TotalUncertainty       []float64
	SourceUncertainty      []float64
	InterpolationUncertainty []float64
	Proximity              []float64
}

func Estimate(data []float64, region *dem.Region, opts *Options) (*Result, error) {
	switch opts.Method {
	case MethodSplitSample:
		return splitSampleUncertainty(data, region, opts)
	case MethodProximity:
		return proximityUncertainty(data, region, opts)
	case MethodCombined:
		return combinedUncertainty(data, region, opts)
	default:
		return nil, fmt.Errorf("unknown uncertainty method: %s", opts.Method)
	}
}

func splitSampleUncertainty(data []float64, region *dem.Region, opts *Options) (*Result, error) {
	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}
	fraction := opts.SampleFraction
	if fraction <= 0 || fraction > 1 {
		fraction = 0.3
	}
	searchRadius := opts.SearchRadius
	if searchRadius <= 0 {
		searchRadius = region.XRes * 5
	}

	var validPts []struct {
		x, y int
		z    float64
		geoX, geoY float64
	}

	gt := region.GeoTransform()
	w, h := region.XSize, region.YSize

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			z := data[idx]
			if z == noData || math.IsNaN(z) {
				continue
			}
			geoX := gt[0] + float64(x)*gt[1] + float64(y)*gt[2]
			geoY := gt[3] + float64(x)*gt[4] + float64(y)*gt[5]
			validPts = append(validPts, struct {
				x, y int
				z    float64
				geoX, geoY float64
			}{x, y, z, geoX, geoY})
		}
	}

	if len(validPts) < 10 {
		return nil, fmt.Errorf("too few valid points for split-sample: %d", len(validPts))
	}

	rng := rand.New(rand.NewSource(opts.Seed))
	rng.Shuffle(len(validPts), func(i, j int) {
		validPts[i], validPts[j] = validPts[j], validPts[i]
	})

	splitIdx := int(float64(len(validPts)) * fraction)
	trainingPts := validPts[splitIdx:]
	validationPts := validPts[:splitIdx]

	interpUnc := make([]float64, w*h)
	sourceUnc := make([]float64, w*h)
	proximity := make([]float64, w*h)

	for i := range interpUnc {
		interpUnc[i] = noData
		sourceUnc[i] = noData
		proximity[i] = noData
	}

	minProximity := make([]float64, w*h)
	for i := range minProximity {
		minProximity[i] = float64(w + h)
	}

	for _, vp := range validationPts {
		var sumWeight, sumVal float64
		var minDist float64 = -1
		neighborCount := 0

		for _, tp := range trainingPts {
			dx := float64(vp.x - tp.x)
			dy := float64(vp.y - tp.y)
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist < 1e-10 {
				sumVal = tp.z
				sumWeight = 1
				minDist = 0
				neighborCount = 1
				break
			}

			if dist <= searchRadius/region.XRes {
				weight := 1.0 / (dist*dist + 1e-15)
				sumWeight += weight
				sumVal += weight * tp.z
				neighborCount++
			}

			if minDist < 0 || dist < minDist {
				minDist = dist
			}
		}

		proximity[vp.y*w+vp.x] = minDist * region.XRes

		if sumWeight > 0 && neighborCount >= 3 {
			predicted := sumVal / sumWeight
			err := vp.z - predicted
			interpUnc[vp.y*w+vp.x] = math.Abs(err)
		}
	}

	interpUncGrid := fillUncertaintyGaps(interpUnc, w, h, noData)
	sourceUncGrid := make([]float64, w*h)
	for i := range sourceUncGrid {
		if interpUncGrid[i] != noData {
			sourceUncGrid[i] = interpUncGrid[i] * 0.5
		} else {
			sourceUncGrid[i] = noData
		}
	}

	totalUnc := make([]float64, w*h)
	for i := range totalUnc {
		if interpUncGrid[i] != noData && sourceUncGrid[i] != noData {
			totalUnc[i] = math.Sqrt(interpUncGrid[i]*interpUncGrid[i] + sourceUncGrid[i]*sourceUncGrid[i])
		} else {
			totalUnc[i] = noData
		}
	}

	return &Result{
		TotalUncertainty:         totalUnc,
		SourceUncertainty:        sourceUncGrid,
		InterpolationUncertainty: interpUncGrid,
		Proximity:                proximity,
	}, nil
}

func proximityUncertainty(data []float64, region *dem.Region, opts *Options) (*Result, error) {
	noData := opts.NoData
	if noData == 0 {
		noData = dem.DefaultNoData
	}

	w, h := region.XSize, region.YSize
	prox := make([]float64, w*h)
	interpUnc := make([]float64, w*h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] == noData || math.IsNaN(data[idx]) {
				prox[idx] = noData
				interpUnc[idx] = noData
				continue
			}

			minDist := float64(w + h)
			for dy := -10; dy <= 10; dy++ {
				for dx := -10; dx <= 10; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx, ny := x+dx, y+dy
					if nx < 0 || nx >= w || ny < 0 || ny >= h {
						continue
					}
					if data[ny*w+nx] != noData {
						dist := math.Sqrt(float64(dx*dx + dy*dy))
						if dist < minDist {
							minDist = dist
						}
					}
				}
			}
			prox[idx] = minDist
			interpUnc[idx] = minDist * region.XRes * 0.1
		}
	}

	return &Result{
		TotalUncertainty:         interpUnc,
		InterpolationUncertainty: interpUnc,
		Proximity:                prox,
		SourceUncertainty:        make([]float64, w*h),
	}, nil
}

func combinedUncertainty(data []float64, region *dem.Region, opts *Options) (*Result, error) {
	ss, err := splitSampleUncertainty(data, region, opts)
	if err != nil {
		return nil, err
	}

	prox, err := proximityUncertainty(data, region, opts)
	if err != nil {
		return nil, err
	}

	w, h := region.XSize, region.YSize
	total := make([]float64, w*h)
	for i := range total {
		if ss.TotalUncertainty[i] != opts.NoData && prox.TotalUncertainty[i] != opts.NoData {
			total[i] = math.Sqrt(ss.TotalUncertainty[i]*ss.TotalUncertainty[i] +
				prox.TotalUncertainty[i]*prox.TotalUncertainty[i])
		} else {
			total[i] = opts.NoData
		}
	}

	return &Result{
		TotalUncertainty:         total,
		SourceUncertainty:        ss.SourceUncertainty,
		InterpolationUncertainty: ss.InterpolationUncertainty,
		Proximity:                prox.Proximity,
	}, nil
}

func fillUncertaintyGaps(data []float64, w, h int, noData float64) []float64 {
	result := make([]float64, len(data))
	copy(result, data)

	var validPts []vec2.T
	var validVals []float64

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] != noData && !math.IsNaN(data[idx]) {
				validPts = append(validPts, vec2.T{float64(x), float64(y)})
				validVals = append(validVals, data[idx])
			}
		}
	}

	if len(validPts) == 0 {
		return result
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			if data[idx] != noData && !math.IsNaN(data[idx]) {
				continue
			}
			var sumWeight, sumVal float64
			for i, vp := range validPts {
				dx := float64(x) - vp[0]
				dy := float64(y) - vp[1]
				dist := dx*dx + dy*dy
				if dist < 1e-10 {
					sumVal = validVals[i]
					sumWeight = 1
					break
				}
				weight := 1.0 / dist
				sumWeight += weight
				sumVal += weight * validVals[i]
			}
			if sumWeight > 0 {
				result[idx] = sumVal / sumWeight
			}
		}
	}

	return result
}

func WriteUncertainty(unc *Result, region *dem.Region, demPath string, noData float64) error {
	if err := dem.CreateDEM(unc.TotalUncertainty, region, demPath+"_tv u.tif", noData); err != nil {
		return fmt.Errorf("total uncertainty: %v", err)
	}
	if err := dem.CreateDEM(unc.InterpolationUncertainty, region, demPath+"_interp_u.tif", noData); err != nil {
		return fmt.Errorf("interpolation uncertainty: %v", err)
	}
	if err := dem.CreateDEM(unc.SourceUncertainty, region, demPath+"_src_u.tif", noData); err != nil {
		return fmt.Errorf("source uncertainty: %v", err)
	}
	if err := dem.CreateDEM(unc.Proximity, region, demPath+"_prox.tif", noData); err != nil {
		return fmt.Errorf("proximity: %v", err)
	}
	return nil
}
