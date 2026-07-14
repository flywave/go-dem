package pointz

import (
	"math"
	"sort"
)

type SMRFGroundClassificationOptions struct {
	CellSize  float64
	Slope     float64
	Window    float64
	Scalar    float64
	Threshold float64
}

func DefaultSMRFOptions() *SMRFGroundClassificationOptions {
	return &SMRFGroundClassificationOptions{
		CellSize:  1.0,
		Slope:     0.15,
		Window:    18.0,
		Scalar:    1.0,
		Threshold: 0.5,
	}
}

func ClassifyGroundSMRF(points []Point3D, opts *SMRFGroundClassificationOptions) []uint8 {
	if opts == nil {
		opts = DefaultSMRFOptions()
	}
	n := len(points)
	if n == 0 {
		return nil
	}
	classification := make([]uint8, n)

	minX, minY := points[0].X, points[0].Y
	maxX, maxY := points[0].X, points[0].Y
	for _, p := range points[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	cell := opts.CellSize
	if cell <= 0 {
		cell = 1.0
	}
	cols := int((maxX-minX)/cell) + 1
	rows := int((maxY-minY)/cell) + 1
	if cols < 1 || rows < 1 {
		return classification
	}

	ZImin := make([]float64, rows*cols)
	for i := range ZImin {
		ZImin[i] = math.NaN()
	}
	for _, p := range points {
		c := int(math.Floor((p.X - minX) / cell))
		r := int(math.Floor((p.Y - minY) / cell))
		if c < 0 {
			c = 0
		}
		if c >= cols {
			c = cols - 1
		}
		if r < 0 {
			r = 0
		}
		if r >= rows {
			r = rows - 1
		}
		idx := r*cols + c
		if math.IsNaN(ZImin[idx]) || p.Z < ZImin[idx] {
			ZImin[idx] = p.Z
		}
	}

	ZImin = knnfillGrid(ZImin, points, cols, rows, minX, minY, cell, 8)

	negZImin := make([]float64, len(ZImin))
	for i, v := range ZImin {
		negZImin[i] = -v
	}
	Low := progressiveFilter(negZImin, 5.0, 1.0, cols, rows, cell)

	slope := opts.Slope
	if slope <= 0 {
		slope = 0.15
	}
	window := opts.Window
	if window <= 0 {
		window = 18.0
	}
	Obj := progressiveFilter(ZImin, slope, window, cols, rows, cell)

	ZIpro := make([]float64, len(ZImin))
	copy(ZIpro, ZImin)
	for i := range ZIpro {
		if Obj[i] == 1 || Low[i] == 1 {
			ZIpro[i] = math.NaN()
		}
	}
	ZIpro = knnfillGrid(ZIpro, points, cols, rows, minX, minY, cell, 8)

	gsurfs := make([]float64, rows*cols)
	for r := 1; r < rows-1; r++ {
		for c := 1; c < cols-1; c++ {
			dx := (ZIpro[r*cols+min(c+1, cols-1)] - ZIpro[r*cols+max(c-1, 0)]) / (2 * cell * cell)
			dy := (ZIpro[min(r+1, rows-1)*cols+c] - ZIpro[max(r-1, 0)*cols+c]) / (2 * cell * cell)
			gsurfs[r*cols+c] = math.Sqrt(dx*dx + dy*dy)
		}
	}
	gsurfsFill := knnfillGrid(gsurfs, points, cols, rows, minX, minY, cell, 8)

	scalar := opts.Scalar
	if scalar <= 0 {
		scalar = 1.0
	}
	threshold := opts.Threshold
	if threshold <= 0 {
		threshold = 0.5
	}

	for i, p := range points {
		c := int(math.Floor((p.X - minX) / cell))
		r := int(math.Floor((p.Y - minY) / cell))
		if c < 0 || c >= cols || r < 0 || r >= rows {
			continue
		}
		idx := r*cols + c
		if math.IsNaN(ZIpro[idx]) || math.IsNaN(gsurfsFill[idx]) {
			continue
		}
		t := threshold + scalar*gsurfsFill[idx]
		if math.Abs(ZIpro[idx]-p.Z) <= t {
			classification[i] = 2
		} else {
			classification[i] = 1
		}
	}

	return classification
}

func knnfillGrid(grid []float64, points []Point3D, cols, rows int, minX, minY, cell float64, k int) []float64 {
	type ptval struct {
		x, y, z float64
	}
	var valid []ptval
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			if !math.IsNaN(grid[r*cols+c]) {
				valid = append(valid, ptval{
					x: minX + (float64(c)+0.5)*cell,
					y: minY + (float64(r)+0.5)*cell,
					z: grid[r*cols+c],
				})
			}
		}
	}
	if len(valid) == 0 {
		return grid
	}

	out := make([]float64, len(grid))
	copy(out, grid)

	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			idx := r*cols + c
			if !math.IsNaN(out[idx]) {
				continue
			}
			x := minX + (float64(c)+0.5)*cell
			y := minY + (float64(r)+0.5)*cell

			type kv struct {
				dist float64
				val  float64
			}
			kvs := make([]kv, len(valid))
			for i, v := range valid {
				dx := x - v.x
				dy := y - v.y
				kvs[i] = kv{dist: dx*dx + dy*dy, val: v.z}
			}
			sort.Slice(kvs, func(i, j int) bool { return kvs[i].dist < kvs[j].dist })
			kn := k
			if kn > len(kvs) {
				kn = len(kvs)
			}

			var sum float64
			for i := 0; i < kn; i++ {
				sum += kvs[i].val
			}
			out[idx] = sum / float64(kn)
		}
	}
	return out
}

func progressiveFilter(ZImin []float64, slope, maxWindow float64, cols, rows int, cell float64) []int {
	maxRadius := int(math.Ceil(maxWindow / cell))
	prevSurface := make([]float64, len(ZImin))
	copy(prevSurface, ZImin)
	prevErosion := make([]float64, len(ZImin))
	copy(prevErosion, ZImin)

	Obj := make([]int, len(ZImin))

	for radius := 1; radius <= maxRadius; radius++ {
		curErosion := erodeDiamond(prevErosion, rows, cols)
		curOpening := dilateDiamond(curErosion, rows, cols, radius)
		copy(prevErosion, curErosion)

		threshold := slope * cell * float64(radius)
		for i := range prevSurface {
			diff := math.Abs(prevSurface[i] - curOpening[i])
			if diff > threshold {
				Obj[i] = 1
			}
		}
		copy(prevSurface, curOpening)
	}
	return Obj
}

func erodeDiamond(src []float64, rows, cols int) []float64 {
	dst := make([]float64, len(src))
	copy(dst, src)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			idx := r*cols + c
			if math.IsNaN(src[idx]) {
				continue
			}
			neighbors := []int{idx}
			if c > 0 {
				neighbors = append(neighbors, r*cols+(c-1))
			}
			if c < cols-1 {
				neighbors = append(neighbors, r*cols+(c+1))
			}
			if r > 0 {
				neighbors = append(neighbors, (r-1)*cols+c)
			}
			if r < rows-1 {
				neighbors = append(neighbors, (r+1)*cols+c)
			}
			minVal := src[idx]
			for _, ni := range neighbors {
				if !math.IsNaN(src[ni]) && src[ni] < minVal {
					minVal = src[ni]
				}
			}
			dst[idx] = minVal
		}
	}
	return dst
}

func dilateDiamond(src []float64, rows, cols int, radius int) []float64 {
	dst := make([]float64, len(src))
	copy(dst, src)
	for iter := 0; iter < radius; iter++ {
		tmp := make([]float64, len(dst))
		copy(tmp, dst)
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				idx := r*cols + c
				neighbors := []int{idx}
				if c > 0 {
					neighbors = append(neighbors, r*cols+(c-1))
				}
				if c < cols-1 {
					neighbors = append(neighbors, r*cols+(c+1))
				}
				if r > 0 {
					neighbors = append(neighbors, (r-1)*cols+c)
				}
				if r < rows-1 {
					neighbors = append(neighbors, (r+1)*cols+c)
				}
				maxVal := tmp[idx]
				for _, ni := range neighbors {
					if !math.IsNaN(tmp[ni]) && tmp[ni] > maxVal {
						maxVal = tmp[ni]
					}
				}
				dst[idx] = maxVal
			}
		}
	}
	return dst
}
