package pointz

import (
	"math"
	"testing"
)

func TestErodeDiamond_Basic(t *testing.T) {
	data := []float64{
		5, 8, 7,
		3, 9, 6,
		4, 2, 1,
	}
	result := erodeDiamond(data, 3, 3)
	if result[4] > 3 {
		t.Errorf("center (idx 4) should be min of neighbors, got %.1f", result[4])
	}
}

func TestErodeDiamond_NaN(t *testing.T) {
	data := []float64{
		5, 8, 7,
		3, math.NaN(), 6,
		4, 2, 1,
	}
	result := erodeDiamond(data, 3, 3)
	if !math.IsNaN(result[4]) {
		t.Error("NaN should propagate")
	}
}

func TestDilateDiamond_Basic(t *testing.T) {
	data := []float64{
		1, 2, 1,
		2, 1, 2,
		1, 2, 1,
	}
	result := dilateDiamond(data, 3, 3, 1)
	if result[4] < 2 {
		t.Errorf("center should be max of neighbors, got %.1f", result[4])
	}
}

func TestDilateDiamond_MultipleIterations(t *testing.T) {
	data := []float64{
		5, 5, 5,
		5, 1, 5,
		5, 5, 5,
	}
	result := dilateDiamond(data, 3, 3, 2)
	if result[4] != 5 {
		t.Errorf("after 2 dilations, center should be 5, got %.1f", result[4])
	}
}

func TestProgressiveFilter_Flat(t *testing.T) {
	data := make([]float64, 25)
	for i := 0; i < 25; i++ {
		data[i] = 100
	}
	result := progressiveFilter(data, 0.5, 3.0, 5, 5, 1.0)
	objCount := 0
	for _, v := range result {
		if v == 1 {
			objCount++
		}
	}
	if objCount != 0 {
		t.Errorf("progressive filter on flat surface: %d/%d marked as objects, expected 0", objCount, len(result))
	}
}

func TestProgressiveFilter_Spike(t *testing.T) {
	data := make([]float64, 25)
	for i := 0; i < 25; i++ {
		data[i] = 100
	}
	data[12] = 200

	result := progressiveFilter(data, 0.5, 3.0, 5, 5, 1.0)
	if result[12] != 1 {
		t.Errorf("spike at center should be detected as object, got %d", result[12])
	}
}

func TestKnnFillGrid_Basic(t *testing.T) {
	grid := []float64{
		10, math.NaN(), 30,
		math.NaN(), 50, math.NaN(),
		70, math.NaN(), 90,
	}
	result := knnfillGrid(grid, 3, 3, 0, 0, 1, 3)
	for i, v := range result {
		if math.IsNaN(v) {
			t.Errorf("cell %d should be filled, got NaN", i)
		}
	}
}

func TestKnnFillGrid_AllNaN(t *testing.T) {
	grid := make([]float64, 9)
	for i := range grid {
		grid[i] = math.NaN()
	}
	result := knnfillGrid(grid, 3, 3, 0, 0, 1, 3)
	for i, v := range result {
		if !math.IsNaN(v) {
			t.Errorf("all-NaN grid should stay NaN, cell %d filled with %.1f", i, v)
		}
	}
}

func TestKnnFillGrid_NoFillNeeded(t *testing.T) {
	grid := []float64{10, 20, 30, 40}
	result := knnfillGrid(grid, 2, 2, 0, 0, 1, 2)
	for i := range grid {
		if math.Abs(result[i]-grid[i]) > 0.01 {
			t.Errorf("cell %d changed: %.1f -> %.1f", i, grid[i], result[i])
		}
	}
}

func TestSurfaceSlope_CellSize(t *testing.T) {
	grid := []float64{
		0, 2, 4,
		0, 2, 4,
		0, 2, 4,
	}
	slopes1 := surfaceSlope(grid, 3, 3, 1)
	if math.Abs(slopes1[4]-2) > 1e-9 {
		t.Errorf("cell=1: expected slope 2 at center, got %f", slopes1[4])
	}
	slopes2 := surfaceSlope(grid, 3, 3, 2)
	if math.Abs(slopes2[4]-1) > 1e-9 {
		t.Errorf("cell=2: expected slope 1 at center, got %f", slopes2[4])
	}
}

func TestSurfaceSlope_BoundaryFilled(t *testing.T) {
	grid := []float64{
		0, 2, 4,
		0, 2, 4,
		0, 2, 4,
	}
	slopes := surfaceSlope(grid, 3, 3, 1)
	for i, v := range slopes {
		if math.IsNaN(v) {
			t.Errorf("boundary cell %d should have slope, got NaN", i)
		}
	}
	if math.Abs(slopes[0]-1) > 1e-9 {
		t.Errorf("corner cell: expected slope 1, got %f", slopes[0])
	}
}

func TestSurfaceSlope_NaNGuard(t *testing.T) {
	grid := []float64{
		0, math.NaN(), 4,
		0, 2, 4,
		0, 2, 4,
	}
	slopes := surfaceSlope(grid, 3, 3, 1)
	if !math.IsNaN(slopes[0]) {
		t.Errorf("cell adjacent to NaN should stay NaN, got %f", slopes[0])
	}
}
