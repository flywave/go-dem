package grits

import (
	"testing"

	"github.com/flywave/go-dem"
)

func TestCut_Basic(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	c := &cutFilter{}
	res, err := c.Run(data, reg, &Options{
		CutBounds: []float64{2, 3, 7, 8},
		NoData:    &nd,
	})
	if err != nil {
		t.Fatalf("cut error: %v", err)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			geoX := float64(x) + 0.5
			geoY := float64(h) - (float64(y) + 0.5)
			inside := geoX >= 2 && geoX <= 7 && geoY >= 3 && geoY <= 8
			if inside && res[idx] == nd {
				t.Errorf("pixel (%d,%d) geo=(%.1f,%.1f) inside cut should be valid", x, y, geoX, geoY)
			}
			if !inside && res[idx] != nd {
				t.Errorf("pixel (%d,%d) geo=(%.1f,%.1f) outside cut should be noData", x, y, geoX, geoY)
			}
		}
	}
}

func TestCut_NoBounds(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	c := &cutFilter{}
	res, err := c.Run(data, reg, &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i, v := range res {
		if v == nd {
			t.Errorf("pixel %d should remain valid without cut bounds", i)
		}
	}
}

func TestCut_EmptyBounds(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	c := &cutFilter{}
	for _, bounds := range [][]float64{nil, {0}, {1, 2, 3}} {
		res, err := c.Run(data, reg, &Options{CutBounds: bounds, NoData: &nd})
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		for i, v := range res {
			if v != data[i] {
				t.Errorf("bounds %v: pixel %d should remain valid, got %.2f", bounds, i, v)
			}
		}
	}
}

func TestCut_OriginBounds(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	c := &cutFilter{}
	res, err := c.Run(data, reg, &Options{CutBounds: []float64{0, 0, 1, 5}, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res[0] != 100 {
		t.Errorf("pixel (0,0) center (0.5,4.5) should be inside x=[0,1] y=[0,5], got %.2f", res[0])
	}
	if res[1] != nd {
		t.Errorf("pixel outside zero-origin bounds should be noData, got %.2f", res[1])
	}
	if res[2*w+2] != nd {
		t.Errorf("pixel outside zero-origin bounds should be noData, got %.2f", res[2*w+2])
	}
}

func TestCut_Invert(t *testing.T) {
	w, h := 8, 8
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	c := &cutFilter{}
	res, err := c.Run(data, reg, &Options{
		CutBounds: []float64{2, 3, 5, 6},
		CutInvert: true,
		NoData:    &nd,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			geoX := float64(x) + 0.5
			geoY := float64(h) - (float64(y) + 0.5)
			inside := geoX >= 2 && geoX <= 5 && geoY >= 3 && geoY <= 6
			if inside && res[idx] != nd {
				t.Errorf("invert: pixel (%d,%d) geo=(%.1f,%.1f) inside should be noData", x, y, geoX, geoY)
			}
			if !inside && res[idx] == nd {
				t.Errorf("invert: pixel (%d,%d) geo=(%.1f,%.1f) outside should be valid", x, y, geoX, geoY)
			}
		}
	}
}
