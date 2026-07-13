package grits

import (
	"math"
	"testing"

	"github.com/flywave/go-dem"
)

func TestParsePolygonWKT_Basic(t *testing.T) {
	wkt := "POLYGON ((0 0, 10 0, 10 10, 0 10, 0 0))"
	poly, xmin, xmax, ymin, ymax, err := parsePolygonWKT(wkt)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(poly) < 3 {
		t.Errorf("polygon too short: %d points", len(poly))
	}
	if math.Abs(xmin) > 1e-10 || math.Abs(ymin) > 1e-10 {
		t.Errorf("expected (0,0), got (%.2f,%.2f)", xmin, ymin)
	}
	if math.Abs(xmax-10) > 1e-10 || math.Abs(ymax-10) > 1e-10 {
		t.Errorf("expected (10,10), got (%.2f,%.2f)", xmax, ymax)
	}
}

func TestPointInPolygon_Inside(t *testing.T) {
	poly, _, _, _, _, _ := parsePolygonWKT("POLYGON ((0 0, 10 0, 10 10, 0 10, 0 0))")
	if !pointInPolygon(5, 5, poly) {
		t.Error("(5,5) should be inside")
	}
}

func TestPointInPolygon_Outside(t *testing.T) {
	poly, _, _, _, _, _ := parsePolygonWKT("POLYGON ((0 0, 10 0, 10 10, 0 10, 0 0))")
	if pointInPolygon(15, 5, poly) {
		t.Error("(15,5) should be outside")
	}
}

func TestPointInPolygon_OnEdge(t *testing.T) {
	poly, _, _, _, _, _ := parsePolygonWKT("POLYGON ((0 0, 10 0, 10 10, 0 10, 0 0))")
	inside := pointInPolygon(5, 0, poly)
	t.Logf("point on edge: inside=%v", inside)
}

func TestClipFilter_NoPolygon(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeRampDEM(w, h)
	cf := &clipFilter{}
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	res, err := cf.Run(data, reg, &Options{NoData: &nd})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for i, v := range res {
		if v != data[i] {
			t.Errorf("no-polygon: pixel %d changed %.0f->%.2f", i, data[i], v)
		}
	}
}

func TestClipFilter_WithPolygon(t *testing.T) {
	w, h := 10, 10
	nd := -9999.0
	data := makeFlatDEM(w, h, 100)
	reg := dem.NewRegionFromBBox(0, 0, float64(w), float64(h), nil, 1, 1)
	cf := &clipFilter{}
	res, err := cf.Run(data, reg, &Options{
		PolygonWKT: "POLYGON ((2 2, 8 2, 8 8, 2 8, 2 2))",
		NoData:     &nd,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if res[0] != nd {
		t.Log("corner (0,0) outside polygon -> masked")
	}
	if res[5*w+5] == nd {
		t.Error("center (5,5) inside polygon but masked")
	}
}

func TestParsePolygonWKT_Complex(t *testing.T) {
	wkt := "POLYGON ((-122.5 37.5, -122.0 37.5, -122.0 38.0, -122.5 38.0, -122.5 37.5))"
	poly, xmin, xmax, ymin, ymax, err := parsePolygonWKT(wkt)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(poly) < 3 {
		t.Errorf("too few points: %d", len(poly))
	}
	if math.Abs(xmin+122.5) > 1e-6 || math.Abs(ymin-37.5) > 1e-6 {
		t.Errorf("bbox mismatch: (%.4f,%.4f)-(%.4f,%.4f)", xmin, ymin, xmax, ymax)
	}
}

func TestParsePolygonWKT_Invalid(t *testing.T) {
	_, _, _, _, _, err := parsePolygonWKT("not a polygon")
	if err == nil {
		t.Error("expected error for invalid WKT")
	}
}
