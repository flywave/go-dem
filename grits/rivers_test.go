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
	if flowCount != w*h-1 {
		t.Errorf("ramp: every cell except the outflow corner should flow, got %d/%d", flowCount, w*h-1)
	}
}

func TestFlowDirection_Flat(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	dir := computeFlowDirection(data, w, h, nd)
	for _, d := range dir {
		if d >= 0 {
			t.Error("flat dem: no cell should have a flow direction")
			break
		}
	}
}

func TestFlowDirection_SinglePit(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[2*w+2] = 0

	dir := computeFlowDirection(data, w, h, nd)
	if dir[2*w+2] != -1 {
		t.Errorf("pit is the lowest cell, should not flow, got %d", dir[2*w+2])
	}
	flowsToCenter := 0
	for i := range dir {
		if dir[i] == 2*w+2 {
			flowsToCenter++
		}
	}
	if flowsToCenter != 8 {
		t.Errorf("all 8 neighbors of the pit should flow into it, got %d", flowsToCenter)
	}
}

func TestFlowDirection_SteepestDescent(t *testing.T) {
	w, h := 3, 3
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	data[1*w+1] = 100
	data[1*w+2] = 60
	data[2] = 50

	dir := computeFlowDirection(data, w, h, nd)
	if dir[1*w+1] != 1*w+2 {
		t.Errorf("east slope 40 should beat diagonal slope %.2f, got dir %d", 50/math.Sqrt2, dir[1*w+1])
	}
}

func TestFlowAccumulation_Basic(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeRampDEM(w, h)
	acc := computeFlowAccumulationIterative(data, w, h, nd)

	maxAcc := 0.0
	for _, v := range acc {
		if v < 1 {
			t.Errorf("every valid cell accumulates at least itself, got %.2f", v)
		}
		if v > maxAcc {
			maxAcc = v
		}
	}
	if maxAcc != float64(w*h) {
		t.Errorf("outflow corner should accumulate the whole grid: %.0f, want %d", maxAcc, w*h)
	}
	if acc[0] != maxAcc {
		t.Errorf("outflow corner (0,0) should hold the maximum, got %.0f vs %.0f", acc[0], maxAcc)
	}
}

func TestFlowAccumulation_Flat(t *testing.T) {
	w, h := 5, 5
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	acc := computeFlowAccumulationIterative(data, w, h, nd)
	for i, v := range acc {
		if v != 1 {
			t.Errorf("flat dem: cell %d has no inflow, accumulation should be 1, got %.2f", i, v)
		}
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
