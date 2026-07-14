package perspecto

import (
	"math"
	"os"
	"testing"
)

func TestHistogramPNG_Basic(t *testing.T) {
	data := make([]float64, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = float64(i)
	}
	img, err := HistogramPNG(data, &HistogramOptions{
		Bins:   50,
		Width:  400,
		Height: 200,
	})
	if err != nil {
		t.Fatalf("histogram error: %v", err)
	}
	if img == nil {
		t.Fatal("nil image")
	}
	bounds := img.Bounds()
	if bounds.Dx() != 400 || bounds.Dy() != 200 {
		t.Errorf("expected 400x200, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestHistogramPNG_CDF(t *testing.T) {
	data := make([]float64, 500)
	for i := 0; i < 500; i++ {
		data[i] = float64(i)
	}
	img, err := HistogramPNG(data, &HistogramOptions{
		Bins:   50,
		Type:   "cdf",
		Width:  300,
		Height: 200,
	})
	if err != nil {
		t.Fatalf("cdf error: %v", err)
	}
	if img == nil {
		t.Fatal("nil image")
	}
}

func TestHistogramPNG_NoData(t *testing.T) {
	data := make([]float64, 100)
	for i := 0; i < 100; i++ {
		data[i] = -9999
	}
	_, err := HistogramPNG(data, &HistogramOptions{NoData: -9999})
	if err == nil {
		t.Error("expected error for all-nodata input")
	}
}

func TestHistogramPNG_AllSame(t *testing.T) {
	data := make([]float64, 100)
	for i := 0; i < 100; i++ {
		data[i] = 42
	}
	img, err := HistogramPNG(data, &HistogramOptions{Bins: 10, Width: 200, Height: 100})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if img == nil {
		t.Fatal("nil image")
	}
}

func TestHistogramPNG_WithInvalid(t *testing.T) {
	data := []float64{1, 2, 3, math.NaN(), 4, 5, -9999, 6}
	img, err := HistogramPNG(data, &HistogramOptions{Bins: 10, NoData: -9999})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if img == nil {
		t.Fatal("nil image")
	}
}

func TestHistogramPNG_DefaultBins(t *testing.T) {
	data := make([]float64, 100)
	for i := 0; i < 100; i++ {
		data[i] = float64(i)
	}
	img, err := HistogramPNG(data, &HistogramOptions{Width: 100, Height: 50})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if img == nil {
		t.Fatal("nil image")
	}
}

func TestHistogramPNG_WithStats(t *testing.T) {
	data := make([]float64, 200)
	for i := 0; i < 200; i++ {
		data[i] = float64(i)
	}
	img, err := HistogramPNG(data, &HistogramOptions{
		Bins:      20,
		Width:     400,
		Height:    300,
		ShowStats: true,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if img == nil {
		t.Fatal("nil image")
	}
}

func TestWriteHistogramPNG(t *testing.T) {
	data := make([]float64, 100)
	for i := 0; i < 100; i++ {
		data[i] = float64(i)
	}
	path := t.TempDir() + "/hist.png"
	err := WriteHistogramPNG(data, path, &HistogramOptions{Bins: 10, Width: 200, Height: 100})
	if err != nil {
		t.Fatalf("write error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("file not created")
	}
}
