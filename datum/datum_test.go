package datum

import (
	"strings"
	"testing"
)

func TestGetFrameByEPSG_Tidal(t *testing.T) {
	f := GetFrameByEPSG(5866)
	if f == nil {
		t.Fatal("MLLW (5866) not found")
	}
	if f.Name != "mllw" {
		t.Errorf("expected mllw, got %s", f.Name)
	}
	if f.Type != FrameTidal {
		t.Errorf("expected tidal, got %v", f.Type)
	}
}

func TestGetFrameByEPSG_HTDP(t *testing.T) {
	f := GetFrameByEPSG(7912)
	if f == nil {
		t.Fatal("ITRF2014 ellipsoid (7912) not found")
	}
	if f.Type != FrameHTDP {
		t.Errorf("expected htdp, got %v", f.Type)
	}
}

func TestGetFrameByEPSG_CDN(t *testing.T) {
	f := GetFrameByEPSG(5773)
	if f == nil {
		t.Fatal("EGM96 (5773) not found")
	}
	if f.Type != FrameCDN {
		t.Errorf("expected cdn, got %v", f.Type)
	}
}

func TestGetFrameByEPSG_NotFound(t *testing.T) {
	f := GetFrameByEPSG(999999)
	if f != nil {
		t.Errorf("expected nil for unknown EPSG, got %v", f)
	}
}

func TestGetFrameByName_Tidal(t *testing.T) {
	f := GetFrameByName("mllw")
	if f == nil {
		t.Fatal("mllw not found by name")
	}
	if f.EPSG != 5866 && f.EPSG != 1089 {
		t.Errorf("expected EPSG 5866 or 1089, got %d", f.EPSG)
	}
}

func TestGetFrameByName_CDN(t *testing.T) {
	f := GetFrameByName("EGM96")
	if f == nil {
		t.Fatal("EGM96 not found by name")
	}
	if f.EPSG != 5773 {
		t.Errorf("expected 5773, got %d", f.EPSG)
	}
}

func TestGetFrameByName_HTDP(t *testing.T) {
	f := GetFrameByName("ITRF2020")
	if f == nil {
		t.Fatal("ITRF2020 not found by name")
	}
	if f.HTDPID != 24 {
		t.Errorf("expected HTDPID 24, got %d", f.HTDPID)
	}
}

func TestFrameTypeName(t *testing.T) {
	if FrameTypeName(5866) != "tidal" {
		t.Errorf("5866: expected tidal, got %s", FrameTypeName(5866))
	}
	if FrameTypeName(7912) != "htdp" {
		t.Errorf("7912: expected htdp, got %s", FrameTypeName(7912))
	}
	if FrameTypeName(5773) != "cdn" {
		t.Errorf("5773: expected cdn, got %s", FrameTypeName(5773))
	}
	if FrameTypeName(999) != "unknown" {
		t.Errorf("999: expected unknown, got %s", FrameTypeName(999))
	}
}

func TestListFrames(t *testing.T) {
	s := ListFrames()
	if !strings.Contains(s, "Tidal") {
		t.Error("expected Tidal header")
	}
	if !strings.Contains(s, "HTDP") {
		t.Error("expected HTDP header")
	}
	if !strings.Contains(s, "CDN") {
		t.Error("expected CDN header")
	}
}

func TestEPSGToVerticalDatum(t *testing.T) {
	if EPSGToVerticalDatum(3855) != 3 {
		t.Errorf("3855: expected EGM2008(3), got %d", EPSGToVerticalDatum(3855))
	}
	if EPSGToVerticalDatum(5773) != 2 {
		t.Errorf("5773: expected EGM96(2), got %d", EPSGToVerticalDatum(5773))
	}
	if EPSGToVerticalDatum(5798) != 1 {
		t.Errorf("5798: expected EGM84(1), got %d", EPSGToVerticalDatum(5798))
	}
	if EPSGToVerticalDatum(5703) != 2 {
		t.Errorf("5703: expected EGM96(2), got %d", EPSGToVerticalDatum(5703))
	}
}

func TestGeoidModels(t *testing.T) {
	if _, ok := GeoidModels["g2018"]; !ok {
		t.Error("g2018 not found")
	}
	if _, ok := GeoidModels["g2012b"]; !ok {
		t.Error("g2012b not found")
	}
	if GeoidModels["g2018"].Uncertainty <= 0 {
		t.Errorf("g2018 uncertainty should be > 0, got %f", GeoidModels["g2018"].Uncertainty)
	}
}

func TestHTDPFrameEpoch(t *testing.T) {
	f := GetFrameByEPSG(7662)
	if f == nil {
		t.Fatal("WGS84(G1674) not found")
	}
	if f.Epoch != 2000.0 {
		t.Errorf("expected epoch 2000.0, got %f", f.Epoch)
	}
}
