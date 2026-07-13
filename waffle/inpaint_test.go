package waffle

import (
	"math"
	"testing"
)

func TestFMMInpaint_NoHole(t *testing.T) {
	w, h := 10, 10
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = float64(y*10 + x)
		}
	}

	result := inpaintFMM(data, w, h, -9999)
	for i := range data {
		if result[i] != data[i] {
			t.Errorf("no-hole: expected %.0f at %d, got %.0f", data[i], i, result[i])
		}
	}
}

func TestFMMInpaint_SmallHole(t *testing.T) {
	w, h := 5, 5
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = 100.0
		}
	}

	center := 2*w + 2
	data[center] = -9999

	result := inpaintFMM(data, w, h, -9999)
	if result[center] == -9999 || math.IsNaN(result[center]) {
		t.Errorf("small hole: center not filled, got %v", result[center])
	}
	if result[center] < 90 || result[center] > 110 {
		t.Errorf("small hole: center should be ~100, got %.2f", result[center])
	}
}

func TestFMMInpaint_EdgeHole(t *testing.T) {
	w, h := 5, 5
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = 50.0
		}
	}

	corner := 0
	data[corner] = -9999

	result := inpaintFMM(data, w, h, -9999)
	if result[corner] != data[corner] {
		t.Logf("edge hole: corner filler value = %.2f", result[corner])
	}
}

func TestFMMInpaint_LinearGradientFill(t *testing.T) {
	w, h := 10, 10
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = float64(y)
		}
	}

	for y := 3; y <= 6; y++ {
		for x := 3; x <= 6; x++ {
			data[y*w+x] = -9999
		}
	}

	result := inpaintFMM(data, w, h, -9999)

	for y := 3; y <= 6; y++ {
		for x := 3; x <= 6; x++ {
			val := result[y*w+x]
			if val == -9999 || math.IsNaN(val) {
				t.Errorf("gradient: hole pixel (%d,%d) not filled", x, y)
				continue
			}
			expected := float64(y)
			if math.Abs(val-expected) > 8.0 {
				t.Errorf("gradient: at (%d,%d) expected ~%.0f, got %.2f", x, y, expected, val)
			}
		}
	}
}

func TestFMMInpaint_LargeHole(t *testing.T) {
	w, h := 20, 20
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = 100.0
		}
	}

	for y := 5; y <= 14; y++ {
		for x := 5; x <= 14; x++ {
			data[y*w+x] = -9999
		}
	}

	result := inpaintFMM(data, w, h, -9999)

	filled := 0
	for y := 5; y <= 14; y++ {
		for x := 5; x <= 14; x++ {
			val := result[y*w+x]
			if val != -9999 && !math.IsNaN(val) {
				filled++
			}
		}
	}
	holePixels := 100
	if filled < holePixels/2 {
		t.Errorf("large hole: only %d/%d filled", filled, holePixels)
	}
}

func TestFMMInpaint_AllNoData(t *testing.T) {
	w, h := 5, 5
	data := make([]float64, w*h)
	for i := range data {
		data[i] = -9999
	}
	result := inpaintFMM(data, w, h, -9999)
	for i, v := range result {
		if v != -9999 && !math.IsNaN(v) {
			t.Errorf("all-nodata: index %d should remain nodata, got %.2f", i, v)
		}
	}
}

func TestFMMInpaint_SinglePixelHole(t *testing.T) {
	w, h := 7, 7
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = 200.0
		}
	}

	data[3*w+3] = -9999

	result := inpaintFMM(data, w, h, -9999)
	val := result[3*w+3]
	if val == -9999 || math.IsNaN(val) {
		t.Fatalf("single pixel hole not filled")
	}
	if val < 180 || val > 220 {
		t.Errorf("single pixel: expected ~200, got %.2f", val)
	}
}

func TestFMMInpaint_CheckboardHoles(t *testing.T) {
	w, h := 8, 8
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if (x+y)%2 == 0 {
				data[y*w+x] = float64(y*10 + x)
			} else {
				data[y*w+x] = -9999
			}
		}
	}

	result := inpaintFMM(data, w, h, -9999)
	filledCount := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			val := result[y*w+x]
			if (x+y)%2 != 0 && val != -9999 && !math.IsNaN(val) {
				filledCount++
			}
		}
	}
	if filledCount < 20 {
		t.Errorf("checkboard: only %d/32 holes filled", filledCount)
	}
}

func TestFMMInpaint_IrregularHole(t *testing.T) {
	w, h := 10, 10
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = 500.0
		}
	}

	holePixels := [][2]int{
		{4, 4}, {5, 4}, {4, 5}, {5, 5},
		{3, 4}, {6, 4}, {4, 3}, {5, 6},
	}
	for _, p := range holePixels {
		data[p[1]*w+p[0]] = -9999
	}

	result := inpaintFMM(data, w, h, -9999)
	for _, p := range holePixels {
		val := result[p[1]*w+p[0]]
		if val == -9999 || math.IsNaN(val) {
			t.Errorf("irregular: pixel (%d,%d) not filled", p[0], p[1])
		}
	}
}

func TestFMMInpaint_HoleAtCorner(t *testing.T) {
	w, h := 5, 5
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = 100.0
		}
	}

	cornerPixels := [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	for _, p := range cornerPixels {
		data[p[1]*w+p[0]] = -9999
	}

	result := inpaintFMM(data, w, h, -9999)
	for _, p := range cornerPixels {
		val := result[p[1]*w+p[0]]
		if val == -9999 || math.IsNaN(val) {
			t.Errorf("corner hole: pixel (%d,%d) not filled", p[0], p[1])
		}
	}
}

func TestFMMInpaint_StripHole(t *testing.T) {
	w, h := 8, 8
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = float64(x)
		}
	}

	for x := 3; x <= 4; x++ {
		for y := 0; y < h; y++ {
			data[y*w+x] = -9999
		}
	}

	result := inpaintFMM(data, w, h, -9999)
	for x := 3; x <= 4; x++ {
		for y := 0; y < h; y++ {
			val := result[y*w+x]
			if val == -9999 || math.IsNaN(val) {
				t.Errorf("strip: column %d row %d not filled", x, y)
				continue
			}
			expected := float64(x)
			if math.Abs(val-expected) > 5.0 {
				t.Errorf("strip: at (%d,%d) expected %.0f, got %.2f", x, y, expected, val)
			}
		}
	}
}

func TestFMMInpaint_EdgePreservation(t *testing.T) {
	w, h := 12, 12
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x < 6 {
				data[y*w+x] = 100.0
			} else {
				data[y*w+x] = 500.0
			}
		}
	}

	for y := 4; y <= 7; y++ {
		for x := 4; x <= 7; x++ {
			data[y*w+x] = -9999
		}
	}

	result := inpaintFMM(data, w, h, -9999)

	leftEdge := result[5*w+5]
	rightEdge := result[5*w+6]
	t.Logf("edge L: %.2f, edge R: %.2f, diff: %.2f", leftEdge, rightEdge, rightEdge-leftEdge)

	if leftEdge != -9999 && rightEdge != -9999 {
		diff := math.Abs(rightEdge - leftEdge)
		if diff > 0 {
			t.Logf("edge contrast preserved: diff=%.2f (goal >0)", diff)
		}
	}
}

func TestFMMInpaint_MountainRidge(t *testing.T) {
	w, h := 15, 15
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			cx, cy := 7.0, 7.0
			dist := math.Sqrt((float64(x)-cx)*(float64(x)-cx) + (float64(y)-cy)*(float64(y)-cy))
			data[y*w+x] = 1000*math.Exp(-dist*dist/20) + 200
		}
	}

	for y := 5; y <= 9; y++ {
		for x := 5; x <= 9; x++ {
			data[y*w+x] = -9999
		}
	}

	result := inpaintFMM(data, w, h, -9999)

	centerVal := result[7*w+7]
	expected := 1000*math.Exp(0) + 200
	if centerVal != -9999 && !math.IsNaN(centerVal) {
		t.Logf("mountain peak: expected %.0f, got %.2f (diff %.0f)", expected, centerVal, expected-centerVal)
	}
}

func benchmarkFMMInpaint(b *testing.B, size int, holeRatio float64) {
	w, h := size, size
	data := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y*w+x] = math.Sin(float64(x)*0.1) + math.Cos(float64(y)*0.1)
		}
	}

	holeSize := int(float64(size) * holeRatio)
	offset := (size - holeSize) / 2
	for y := offset; y < offset+holeSize; y++ {
		for x := offset; x < offset+holeSize; x++ {
			data[y*w+x] = -9999
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testData := make([]float64, len(data))
		copy(testData, data)
		_ = inpaintFMM(testData, w, h, -9999)
	}
}

func BenchmarkFMMInpaint_50x50(b *testing.B)  { benchmarkFMMInpaint(b, 50, 0.3) }
func BenchmarkFMMInpaint_100x100(b *testing.B) { benchmarkFMMInpaint(b, 100, 0.2) }
func BenchmarkFMMInpaint_200x200(b *testing.B) { benchmarkFMMInpaint(b, 200, 0.1) }


