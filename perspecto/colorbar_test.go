package perspecto

import (
	"os"
	"testing"
)

func TestColorbarPNG_Basic(t *testing.T) {
	cmap := DefaultTerrainColormap()
	img := ColorbarPNG(cmap, &ColorbarOptions{Width: 200, Height: 50})
	if img == nil {
		t.Fatal("nil image")
	}
	bounds := img.Bounds()
	if bounds.Dx() != 200 || bounds.Dy() != 50 {
		t.Errorf("expected 200x50, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestColorbarPNG_DefaultSize(t *testing.T) {
	img := ColorbarPNG(DefaultTerrainColormap(), &ColorbarOptions{})
	bounds := img.Bounds()
	if bounds.Dx() != 600 || bounds.Dy() != 60 {
		t.Errorf("expected 600x60, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestColorbarPNG_EmptyCmap(t *testing.T) {
	img := ColorbarPNG([]ColorStop{}, &ColorbarOptions{Width: 100, Height: 30})
	if img == nil {
		t.Fatal("nil image")
	}
}

func TestColorbarPNG_NonWhitePixels(t *testing.T) {
	img := ColorbarPNG(DefaultTerrainColormap(), &ColorbarOptions{Width: 200, Height: 50})
	hasColor := false
	for y := 15; y < 40; y++ {
		for x := 50; x < 180; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r != 65535 || g != 65535 || b != 65535 {
				hasColor = true
				break
			}
		}
	}
	if !hasColor {
		t.Error("colorbar should have colored pixels")
	}
}

func TestWriteColorbarPNG(t *testing.T) {
	path := t.TempDir() + "/colorbar.png"
	err := WriteColorbarPNG(DefaultTerrainColormap(), path, &ColorbarOptions{Width: 100, Height: 30})
	if err != nil {
		t.Fatalf("write error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("file not created")
	}
}
