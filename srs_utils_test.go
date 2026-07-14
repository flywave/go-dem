package dem

import (
	"testing"

	"github.com/flywave/go-geo"
)

func TestParseSRS_EPSG(t *testing.T) {
	p, err := ParseSRS("EPSG:4326")
	if err != nil {
		t.Fatalf("parse EPSG:4326: %v", err)
	}
	if p == nil {
		t.Fatal("nil proj")
	}
	if !SRSIsLatLong(p) {
		t.Error("EPSG:4326 should be lat/long")
	}
}

func TestParseSRS_Empty(t *testing.T) {
	_, err := ParseSRS("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestParseSRS_Invalid(t *testing.T) {
	p := geo.NewProj("INVALID")
	if p == nil {
		t.Log("invalid SRS returns nil (expected)")
	}
}

func TestSRSToEPSG(t *testing.T) {
	p, _ := ParseSRS("EPSG:4326")
	code := SRSToEPSG(p)
	if code != 4326 {
		t.Errorf("expected 4326, got %d", code)
	}
}

func TestSRSToEPSG_Zero(t *testing.T) {
	code := SRSToEPSG(nil)
	if code != 0 {
		t.Errorf("nil: expected 0, got %d", code)
	}
}

func TestSRSGetAuthorityCode(t *testing.T) {
	p, _ := ParseSRS("EPSG:4326")
	code := SRSGetAuthorityCode(p)
	if code != "EPSG:4326" {
		t.Errorf("expected EPSG:4326, got %s", code)
	}
}

func TestSRSGetCSType(t *testing.T) {
	if tp := SRSGetCSType("EPSG:4326"); tp != "GEOGCS" {
		t.Errorf("EPSG:4326: expected GEOGCS, got %s", tp)
	}
	if tp := SRSGetCSType("EPSG:3857"); tp != "PROJCS" {
		t.Errorf("EPSG:3857: expected PROJCS, got %s", tp)
	}
	if tp := SRSGetCSType(""); tp != "UNKNOWN" {
		t.Errorf("empty: expected UNKNOWN, got %s", tp)
	}
}

func TestSRSIsProjected(t *testing.T) {
	p, _ := ParseSRS("EPSG:3857")
	if !SRSIsProjected(p) {
		t.Error("EPSG:3857 should be projected")
	}
	p2, _ := ParseSRS("EPSG:4326")
	if SRSIsProjected(p2) {
		t.Error("EPSG:4326 should not be projected")
	}
}

func TestSRSEquals(t *testing.T) {
	p1, _ := ParseSRS("EPSG:4326")
	p2, _ := ParseSRS("EPSG:4326")
	p3, _ := ParseSRS("EPSG:3857")
	if !SRSEquals(p1, p2) {
		t.Error("same SRS should be equal")
	}
	if SRSEquals(p1, p3) {
		t.Error("different SRS should not be equal")
	}
	if !SRSEquals(nil, nil) {
		t.Error("nil nil should be equal")
	}
	if SRSEquals(p1, nil) {
		t.Error("non-nil and nil should not be equal")
	}
}

func TestSRSClone(t *testing.T) {
	p, _ := ParseSRS("EPSG:4326")
	c := SRSClone(p)
	if c == nil {
		t.Fatal("clone returned nil")
	}
	if !SRSEquals(p, c) {
		t.Error("clone should equal original")
	}
	code := SRSToEPSG(c)
	if code != 4326 {
		t.Errorf("clone epsg: expected 4326, got %d", code)
	}
}

func TestSRSClone_Nil(t *testing.T) {
	c := SRSClone(nil)
	if c != nil {
		t.Error("nil clone should return nil")
	}
}

func TestSRSIsLatLong(t *testing.T) {
	if SRSIsLatLong(nil) {
		t.Error("nil should not be latlong")
	}
	p, _ := ParseSRS("EPSG:4326")
	if !SRSIsLatLong(p) {
		t.Error("EPSG:4326 should be latlong")
	}
}

func TestSRSToWKT(t *testing.T) {
	p, _ := ParseSRS("EPSG:4326")
	wkt := SRSToWKT(p)
	if wkt == "" {
		t.Error("WKT should not be empty")
	}
	t.Logf("WKT: %s", wkt[:min(len(wkt), 80)])
}

func TestSRSToProj4(t *testing.T) {
	p, _ := ParseSRS("EPSG:4326")
	p4 := SRSToProj4(p)
	if p4 == "" {
		t.Error("Proj4 should not be empty")
	}
}

func TestSRSHorizontalFromCompound(t *testing.T) {
	h := SRSHorizontalFromCompound("EPSG:4326+EPSG:5773")
	if h != "EPSG:4326" {
		t.Errorf("expected EPSG:4326, got %s", h)
	}
}

func TestSRSVerticalFromCompound(t *testing.T) {
	v := SRSVerticalFromCompound("EPSG:4326+EPSG:5773")
	if v != "EPSG:5773" {
		t.Errorf("expected EPSG:5773, got %s", v)
	}
}

func TestSRSCompound(t *testing.T) {
	c := SRSCompound("EPSG:4326", "EPSG:5773")
	if c != "EPSG:4326+EPSG:5773" {
		t.Errorf("expected EPSG:4326+EPSG:5773, got %s", c)
	}
	c2 := SRSCompound("EPSG:4326", "")
	if c2 != "EPSG:4326" {
		t.Errorf("empty vert: expected EPSG:4326, got %s", c2)
	}
}
