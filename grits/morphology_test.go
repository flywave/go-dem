package grits

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func region10x10() *dem.Region {
	return dem.NewRegionFromBBox(0, 0, 10, 10, nil, 1, 1)
}

func TestErode_Basic(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[5*w+5] = nd

	e := &erodeFilter{}
	res, err := e.Run(data, region10x10(), &Options{Radius: 1, NoData: &nd})
	if err != nil {
		t.Fatalf("erode error: %v", err)
	}
	if res[5*w+5] != nd {
		t.Errorf("erode: hole should stay noData, got %.2f", res[5*w+5])
	}
	if res[4*w+4] != nd {
		t.Errorf("erode: neighbor within radius of noData should erode, got %.2f", res[4*w+4])
	}
	if res[1*w+1] != 100 {
		t.Errorf("erode: pixel far from noData should survive, got %.2f", res[1*w+1])
	}
}

func TestErode_AllValid(t *testing.T) {
	data := makeFlatDEM(10, 10, 100)
	nd := -9999.0
	e := &erodeFilter{}
	res, err := e.Run(data, region10x10(), &Options{Radius: 1, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	erodedBorder := 0
	for _, v := range res {
		if v == nd || math.IsNaN(v) {
			erodedBorder++
		}
	}
	t.Logf("erode: %d border pixels eroded (expected with radius=1)", erodedBorder)
	if erodedBorder > 40 {
		t.Errorf("too many pixels eroded: %d", erodedBorder)
	}
}

func TestErode_BorderNoData(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[0] = nd

	e := &erodeFilter{}
	res, err := e.Run(data, region5x5(), &Options{Radius: 1, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res[1] != nd || res[w] != nd {
		t.Error("erode: noData should propagate to orthogonal neighbors")
	}
	if res[w+1] != nd {
		t.Errorf("erode: square window includes diagonal of noData, should erode, got %.2f", res[w+1])
	}
	if res[2*w+2] != 100 {
		t.Errorf("erode: pixel beyond window of noData should survive, got %.2f", res[2*w+2])
	}
}

func TestDilate_Basic(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, nd)
	data[5*w+5] = 100

	d := &dilateFilter{}
	res, err := d.Run(data, region10x10(), &Options{Radius: 1, NoData: &nd})
	if err != nil {
		t.Fatalf("dilate error: %v", err)
	}
	dilated := 0
	for i := range res {
		if res[i] != nd {
			dilated++
		}
	}
	if dilated != 9 {
		t.Errorf("dilate should expand seed to 3x3 block, got %d valid pixels", dilated)
	}
	if res[5*w+5] != 100 {
		t.Errorf("original seed changed: %.2f", res[5*w+5])
	}
}

func TestDilate_TakesMax(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, nd)
	data[5*w+5] = 100
	data[5*w+6] = 200

	d := &dilateFilter{}
	res, err := d.Run(data, region10x10(), &Options{Radius: 1, NoData: &nd})
	if err != nil {
		t.Fatalf("dilate error: %v", err)
	}
	if res[4*w+5] != 200 {
		t.Errorf("dilate must take window maximum, expected 200, got %.2f", res[4*w+5])
	}
	if res[6*w+5] != 200 {
		t.Errorf("dilate must take window maximum, expected 200, got %.2f", res[6*w+5])
	}
}

func TestDilate_AllNoData(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, nd)
	d := &dilateFilter{}
	res, err := d.Run(data, region5x5(), &Options{Radius: 1, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i, v := range res {
		if v != nd {
			t.Errorf("all-nodata: pixel %d should remain nodata, got %.2f", i, v)
		}
	}
}

func TestOpen_Basic(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	for y := 3; y <= 6; y++ {
		for x := 3; x <= 6; x++ {
			data[y*w+x] = nd
		}
	}

	o := &openFilter{}
	res, err := o.Run(data, region10x10(), &Options{Radius: 1, NoData: &nd})
	if err != nil {
		t.Fatalf("open error: %v", err)
	}
	if res[5*w+5] != nd {
		t.Errorf("open: hole larger than structuring element should stay noData, got %.2f", res[5*w+5])
	}
	if res[0] != 100 {
		t.Errorf("open: valid area should be restored by the dilate step, got %.2f", res[0])
	}
}

func TestClose_Basic(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, nd)
	data[5*w+5] = 100

	c := &closeFilter{}
	res, err := c.Run(data, region10x10(), &Options{Radius: 1, NoData: &nd})
	if err != nil {
		t.Fatalf("close error: %v", err)
	}
	valid := 0
	for _, v := range res {
		if v != nd {
			valid++
		}
	}
	if valid != 1 {
		t.Errorf("close of single seed should keep exactly the seed pixel, got %d valid", valid)
	}
	if res[5*w+5] != 100 {
		t.Errorf("close should preserve seed value, got %.2f", res[5*w+5])
	}
}

func TestDilate_AllValidStaysValid(t *testing.T) {
	w, h := 8, 8
	nd := -9999.0
	data := makeRampDEM(w, h)

	d := &dilateFilter{}
	res, _ := d.Run(data, dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1), &Options{Radius: 1, NoData: &nd})
	for i := range res {
		if res[i] == nd {
			t.Errorf("all-valid: pixel %d became noData after dilate", i)
		}
		if res[i] < data[i] {
			t.Errorf("dilate should not decrease values: pixel %d %.0f -> %.2f", i, data[i], res[i])
		}
	}
}

func TestOpenClose_Involution(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[5*w+5] = nd
	data[5*w+6] = nd

	o := &openFilter{}
	afterOpen, _ := o.Run(data, region10x10(), &Options{Radius: 1, NoData: &nd})

	c := &closeFilter{}
	afterClose, _ := c.Run(afterOpen, region10x10(), &Options{Radius: 1, NoData: &nd})

	if len(afterClose) != len(data) {
		t.Error("output size changed")
	}
	if afterOpen[5*w+5] != nd || afterOpen[5*w+6] != nd {
		t.Errorf("open should keep holes larger than the structuring element, got %.2f and %.2f",
			afterOpen[5*w+5], afterOpen[5*w+6])
	}
	if afterClose[5*w+5] != 100 {
		t.Errorf("close should fill the small remaining gap, got %.2f", afterClose[5*w+5])
	}
}
