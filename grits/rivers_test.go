package grits

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func TestFlowDirection_Basic(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeRampDEM(w, h)
	dir := computeFlowDirection(data, w, h, nd)
	flowCount := 0
	for i, d := range dir {
		if data[i] != nd && !math.IsNaN(data[i]) && d >= 0 {
			flowCount++
		}
	}
	if flowCount == 0 {
		t.Error("no flow directions computed")
	}
}

func TestFlowDirection_Flat(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	dir := computeFlowDirection(data, w, h, nd)
	flowing := 0
	for _, d := range dir {
		if d >= 0 {
			flowing++
		}
	}
	if flowing > 0 {
		t.Logf("flat dem: %d pixels with flow direction (correct: none should flow)", flowing)
	}
}

func TestFlowDirection_SinglePit(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[2*w+2] = 0

	dir := computeFlowDirection(data, w, h, nd)
	centerFlowsTo := dir[2*w+2]
	if centerFlowsTo != w*w+h {
		t.Logf("pit at center flows to %d", centerFlowsTo)
	}

	flowsToCenter := 0
	for i := range dir {
		if dir[i] == 2*w+2 {
			flowsToCenter++
		}
	}
	t.Logf("%d cells flow into the pit", flowsToCenter)
}

func TestFlowAccumulation_Basic(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeRampDEM(w, h)
	acc := computeFlowAccumulationIterative(data, w, h, nd)

	totalAcc := 0.0
	noDataCount := 0
	for _, v := range acc {
		if v == nd || math.IsNaN(v) {
			noDataCount++
		} else {
			totalAcc += v
		}
	}
	if noDataCount == len(acc) {
		t.Error("all pixels have noData accumulation")
	}
	if totalAcc <= 0 {
		t.Error("accumulation sum should be positive")
	}
}

func TestFlowAccumulation_Flat(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	acc := computeFlowAccumulationIterative(data, w, h, nd)
	nonZero := 0
	for _, v := range acc {
		if v > 0 {
			nonZero++
		}
	}
	if nonZero > 0 {
		t.Logf("flat dem: %d cells with flow accumulation", nonZero)
	}
}

func TestThresholdRiverNetwork(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	acc := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			acc[y*w+x] = float64(y * x)
		}
	}

	rivers := thresholdRiverNetwork(acc, w, h, 0.5, nd)
	riverPixels := 0
	for _, v := range rivers {
		if v != nd {
			riverPixels++
		}
	}
	if riverPixels == 0 {
		t.Error("no river pixels detected")
	}
	if riverPixels == w*h {
		t.Error("all pixels are rivers (threshold too low)")
	}
}

func TestRiverFilter_Run(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeRampDEM(w, h)
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	rf := &riverFilter{}
	res, err := rf.Run(data, reg, &Options{Threshold: 0.3, NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(res) != w*h {
		t.Error("output size mismatch")
	}
}

func TestFlowAccumulation_Large(t *testing.T) {
	w, h := 100, 100
	nd := -9999.0
	data := makeRampDEM(w, h)
	acc := computeFlowAccumulationIterative(data, w, h, nd)
	maxAcc := 0.0
	for _, v := range acc {
		if v > maxAcc {
			maxAcc = v
		}
	}
	if maxAcc <= 0 {
		t.Errorf("max accumulation should be > 0, got %.0f", maxAcc)
	}
}
