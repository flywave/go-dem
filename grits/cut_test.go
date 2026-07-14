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
	gt := reg.GeoTransform()
	c := &cutFilter{}
	res, err := c.Run(data, reg, &Options{
		CutBounds: [4]float64{2, 3, 7, 8},
		NoData:    &nd,
	})
	if err != nil {
		t.Fatalf("cut error: %v", err)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			geoX := gt[0] + float64(x)*gt[1]
			geoY := gt[3] + float64(y)*gt[5]
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

func TestCut_Invert(t *testing.T) {
	w, h := 8, 8
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	gt := reg.GeoTransform()
	c := &cutFilter{}
	res, err := c.Run(data, reg, &Options{
		CutBounds: [4]float64{2, 3, 5, 6},
		CutInvert: true,
		NoData:    &nd,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			geoX := gt[0] + float64(x)*gt[1]
			geoY := gt[3] + float64(y)*gt[5]
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
