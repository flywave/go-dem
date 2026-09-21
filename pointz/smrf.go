package pointz

import (
	"math"
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

	ZImin = knnfillGrid(ZImin, cols, rows, minX, minY, cell, 8)

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
	ZIpro = knnfillGrid(ZIpro, cols, rows, minX, minY, cell, 8)

	gsurfs := surfaceSlope(ZIpro, rows, cols, cell)
	gsurfsFill := knnfillGrid(gsurfs, cols, rows, minX, minY, cell, 8)

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

func surfaceSlope(ZIpro []float64, rows, cols int, cell float64) []float64 {
	gsurfs := make([]float64, rows*cols)
	for i := range gsurfs {
		gsurfs[i] = math.NaN()
	}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			zl := ZIpro[r*cols+max(c-1, 0)]
			zr := ZIpro[r*cols+min(c+1, cols-1)]
			zd := ZIpro[max(r-1, 0)*cols+c]
			zu := ZIpro[min(r+1, rows-1)*cols+c]
			if math.IsNaN(zl) || math.IsNaN(zr) || math.IsNaN(zd) || math.IsNaN(zu) {
				continue
			}
			dx := (zr - zl) / (2 * cell)
			dy := (zu - zd) / (2 * cell)
			gsurfs[r*cols+c] = math.Sqrt(dx*dx + dy*dy)
		}
	}
	return gsurfs
}

type kv struct {
	dist float64
	val  float64
}

func siftDownKV(h []kv, i int) {
	for {
		l := 2*i + 1
		if l >= len(h) {
			return
		}
		m := l
		if rr := l + 1; rr < len(h) && h[rr].dist > h[l].dist {
			m = rr
		}
		if h[i].dist >= h[m].dist {
			return
		}
		h[i], h[m] = h[m], h[i]
		i = m
	}
}

func knnfillGrid(grid []float64, cols, rows int, minX, minY, cell float64, k int) []float64 {
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

	kn := k
	if kn > len(valid) {
		kn = len(valid)
	}

	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			idx := r*cols + c
			if !math.IsNaN(out[idx]) {
				continue
			}
			x := minX + (float64(c)+0.5)*cell
			y := minY + (float64(r)+0.5)*cell

			h := make([]kv, 0, kn)
			for _, v := range valid {
				dx := x - v.x
				dy := y - v.y
				d := dx*dx + dy*dy
				if len(h) < kn {
					h = append(h, kv{dist: d, val: v.z})
					if len(h) == kn {
						for i := len(h)/2 - 1; i >= 0; i-- {
							siftDownKV(h, i)
						}
					}
					continue
				}
				if d < h[0].dist {
					h[0] = kv{dist: d, val: v.z}
					siftDownKV(h, 0)
				}
			}

			var sum float64
			for _, e := range h {
				sum += e.val
			}
			out[idx] = sum / float64(len(h))
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
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			idx := r*cols + c
			v := src[idx]
			if math.IsNaN(v) {
				dst[idx] = v
				continue
			}
			minVal := v
			if c > 0 {
				if n := src[idx-1]; !math.IsNaN(n) && n < minVal {
					minVal = n
				}
			}
			if c < cols-1 {
				if n := src[idx+1]; !math.IsNaN(n) && n < minVal {
					minVal = n
				}
			}
			if r > 0 {
				if n := src[idx-cols]; !math.IsNaN(n) && n < minVal {
					minVal = n
				}
			}
			if r < rows-1 {
				if n := src[idx+cols]; !math.IsNaN(n) && n < minVal {
					minVal = n
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
				maxVal := tmp[idx]
				if c > 0 {
					if n := tmp[idx-1]; !math.IsNaN(n) && n > maxVal {
						maxVal = n
					}
				}
				if c < cols-1 {
					if n := tmp[idx+1]; !math.IsNaN(n) && n > maxVal {
						maxVal = n
					}
				}
				if r > 0 {
					if n := tmp[idx-cols]; !math.IsNaN(n) && n > maxVal {
						maxVal = n
					}
				}
				if r < rows-1 {
					if n := tmp[idx+cols]; !math.IsNaN(n) && n > maxVal {
						maxVal = n
					}
				}
				dst[idx] = maxVal
			}
		}
	}
	return dst
}
