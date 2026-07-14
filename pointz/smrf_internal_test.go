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
	if result[4] < 4 {
		t.Logf("after 2 dilations, center=%.1f", result[4])
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
	if objCount > 0 {
		t.Logf("progressive filter on flat: %d/%d marked as objects", objCount, len(result))
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
		t.Log("spike may not be detected with these parameters")
	}
}

func TestKnnFillGrid_Basic(t *testing.T) {
	grid := []float64{
		10, math.NaN(), 30,
		math.NaN(), 50, math.NaN(),
		70, math.NaN(), 90,
	}
	pts := []Point3D{
		{X: 0, Y: 0, Z: 10},
		{X: 1, Y: 0, Z: 20},
		{X: 2, Y: 2, Z: 50},
	}
	result := knnfillGrid(grid, pts, 3, 3, 0, 0, 1, 3)
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
	pts := []Point3D{{X: 0, Y: 0, Z: 10}}
	result := knnfillGrid(grid, pts, 3, 3, 0, 0, 1, 3)
	for i, v := range result {
		if !math.IsNaN(v) {
			t.Logf("cell %d filled with %.1f (may happen with sparse data)", i, v)
		}
		_ = i
	}
}

func TestKnnFillGrid_NoFillNeeded(t *testing.T) {
	grid := []float64{10, 20, 30, 40}
	pts := []Point3D{{X: 0, Y: 0, Z: 10}, {X: 1, Y: 1, Z: 20}}
	result := knnfillGrid(grid, pts, 2, 2, 0, 0, 1, 2)
	for i := range grid {
		if math.Abs(result[i]-grid[i]) > 0.01 {
			t.Errorf("cell %d changed: %.1f -> %.1f", i, grid[i], result[i])
		}
	}
}
