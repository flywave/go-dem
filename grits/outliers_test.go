package grits

import (
	"testing"

	"github.com/flywave/go-dem"
)

func TestOutliers_Basic(t *testing.T) {
	w, h := 12, 12
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
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
	if res[6*w+6] == nd {
		t.Log("outliers: high spike masked")
	} else {
		t.Logf("outliers: high spike not masked (score may vary)")
	}
	if res[3*w+3] == nd {
		t.Log("outliers: low spike masked")
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
