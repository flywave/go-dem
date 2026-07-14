package perspecto

import (
	"math"
	"testing"
)

func TestRescaleColormap(t *testing.T) {
	cmap := []ColorStop{
		{Value: 0, R: 0, G: 0, B: 255},
		{Value: 100, R: 0, G: 255, B: 0},
		{Value: 200, R: 255, G: 0, B: 0},
	}
	rescaled := RescaleColormap(cmap, 1000, 2000)
	if rescaled[0].Value != 1000 {
		t.Errorf("min: expected 1000, got %.0f", rescaled[0].Value)
	}
	if rescaled[2].Value != 2000 {
		t.Errorf("max: expected 2000, got %.0f", rescaled[2].Value)
	}
	if math.Abs(rescaled[1].Value-1500) > 0.01 {
		t.Errorf("mid: expected 1500, got %.0f", rescaled[1].Value)
	}
}

func TestRescaleColormap_Empty(t *testing.T) {
	r := RescaleColormap(nil, 0, 100)
	if r != nil {
		t.Error("nil input should return nil")
	}
}

func TestRescaleColormap_ZeroRange(t *testing.T) {
	cmap := []ColorStop{{Value: 0, R: 255, G: 0, B: 0}}
	r := RescaleColormap(cmap, 100, 100)
	if len(r) != 1 {
		t.Error("zero range should return original")
	}
}

func TestInterpolateColor_BelowMin(t *testing.T) {
	cmap := []ColorStop{{Value: 10, R: 255, G: 0, B: 0}, {Value: 20, R: 0, G: 255, B: 0}}
	r, g, b := InterpolateColor(cmap, 5)
	if r != 255 || g != 0 || b != 0 {
		t.Errorf("below min: expected (255,0,0), got (%d,%d,%d)", r, g, b)
	}
}

func TestInterpolateColor_AboveMax(t *testing.T) {
	cmap := []ColorStop{{Value: 10, R: 255, G: 0, B: 0}, {Value: 20, R: 0, G: 255, B: 0}}
	r, g, b := InterpolateColor(cmap, 30)
	if r != 0 || g != 255 || b != 0 {
		t.Errorf("above max: expected (0,255,0), got (%d,%d,%d)", r, g, b)
	}
}

func TestInterpolateColor_Exact(t *testing.T) {
	cmap := []ColorStop{{Value: 10, R: 255, G: 0, B: 0}, {Value: 20, R: 0, G: 255, B: 0}}
	r, g, b := InterpolateColor(cmap, 10)
	if r != 255 || g != 0 || b != 0 {
		t.Errorf("exact: expected (255,0,0), got (%d,%d,%d)", r, g, b)
	}
}

func TestInterpolateColor_Interpolated(t *testing.T) {
	cmap := []ColorStop{{Value: 0, R: 0, G: 0, B: 0}, {Value: 10, R: 100, G: 100, B: 100}}
	r, g, b := InterpolateColor(cmap, 5)
	if r != 50 || g != 50 || b != 50 {
		t.Errorf("mid: expected (50,50,50), got (%d,%d,%d)", r, g, b)
	}
}

func TestInterpolateColor_Empty(t *testing.T) {
	r, g, b := InterpolateColor(nil, 5)
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("empty: expected (0,0,0), got (%d,%d,%d)", r, g, b)
	}
}

func TestDefaultColormaps(t *testing.T) {
	terrain := DefaultTerrainColormap()
	if len(terrain) == 0 {
		t.Error("terrain colormap should not be empty")
	}
	bathy := DefaultBathymetryColormap()
	if len(bathy) == 0 {
		t.Error("bathymetry colormap should not be empty")
	}
}
