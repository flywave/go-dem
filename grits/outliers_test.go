package grits

import (
	"testing"

	"github.com/flywave/go-dem"
)

func TestOutliers_Basic(t *testing.T) {
	w, h := 12, 12
	nd := -9999.0
	data := makeRampDEM(w, h)
	data[6*w+6] = 9999
	data[3*w+3] = -9998

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	o := &outliersFilter{}
	res, err := o.Run(data, reg, &Options{
		Percentile: 75,
		Threshold:  1.5,
		Iterations: 1,
		NoData:     &nd,
	})
	if err != nil {
		t.Fatalf("outliers error: %v", err)
	}
	if res[6*w+6] != nd {
		t.Errorf("high spike not masked: %.2f", res[6*w+6])
	}
	if res[3*w+3] != nd {
		t.Errorf("low spike not masked: %.2f", res[3*w+3])
	}
}

func TestOutliers_EdgeNoData(t *testing.T) {
	w, h := 12, 12
	nd := -9999.0
	data := makeRampDEM(w, h)
	for y := 0; y < h; y++ {
		data[y*w] = nd
		data[y*w+1] = nd
	}
	data[5*w+2] = 9999

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	o := &outliersFilter{}
	res, err := o.Run(data, reg, &Options{
		Percentile: 75,
		Threshold:  1.5,
		Iterations: 1,
		NoData:     &nd,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res[5*w+2] != nd {
		t.Errorf("spike near noData boundary not masked: %.2f", res[5*w+2])
	}
	if res[2*w+2] == nd {
		t.Errorf("normal pixel near noData boundary masked: %.2f", res[2*w+2])
	}
	for y := 0; y < h; y++ {
		if res[y*w] != nd {
			t.Errorf("noData pixel (%d,0) should stay noData, got %.2f", y, res[y*w])
		}
	}
}

func TestOutliers_Multipass(t *testing.T) {
	w, h := 15, 15
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[7*w+7] = 9999

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	o := &outliersFilter{}
	res, err := o.Run(data, reg, &Options{
		Percentile: 75,
		Threshold:  1.5,
		Iterations: 3,
		NoData:     &nd,
		Method:     "aggressive",
	})
	if err != nil {
		t.Fatalf("multipass error: %v", err)
	}
	if len(res) != len(data) {
		t.Errorf("output size mismatch")
	}
}

func TestOutliers_AllFlat(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	o := &outliersFilter{}
	res, err := o.Run(data, reg, &Options{
		Percentile: 75,
		Threshold:  1.5,
		NoData:     &nd,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i, v := range res {
		if v == nd {
			t.Errorf("flat pixel %d should not be outlier", i)
		}
	}
}

func TestOutliers_EdgeCluster(t *testing.T) {
	w, h := 8, 8
	nd := -9999.0
	data := makeRampDEM(w, h)

	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	o := &outliersFilter{}
	_, err := o.Run(data, reg, &Options{
		Percentile: 75,
		Threshold:  1.5,
		NoData:     &nd,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
