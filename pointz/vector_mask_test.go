package pointz

import (
	"os"
	"testing"
)

func TestVectorMaskFilter_Basic(t *testing.T) {
	geoJSON := `{
		"type":"FeatureCollection",
		"features":[{
			"type":"Feature",
			"properties":{},
			"geometry":{
				"type":"Polygon",
				"coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]
			}
		}]
	}`

	tmp, err := os.CreateTemp("", "mask_*.geojson")
	if err != nil {
		t.Fatal(err)
	}
	tmp.WriteString(geoJSON)
	tmp.Close()
	defer os.Remove(tmp.Name())

	pts := []Point3D{
		{X: 5, Y: 5, Z: 100},
		{X: -5, Y: 5, Z: 100},
		{X: 15, Y: 15, Z: 100},
	}

	mask, err := VectorMaskFilter(pts, &VectorMaskOptions{
		MaskPath: tmp.Name(),
		Invert:   false,
	})
	if err != nil {
		t.Fatalf("vector mask error: %v", err)
	}
	if mask[0] {
		t.Error("point inside polygon should not be masked (invert=false)")
	}
	if !mask[1] || !mask[2] {
		t.Error("points outside polygon should be masked (invert=false)")
	}
}

func TestVectorMaskFilter_Invert(t *testing.T) {
	geoJSON := `{
		"type":"FeatureCollection",
		"features":[{
			"type":"Feature",
			"properties":{},
			"geometry":{
				"type":"Polygon",
				"coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]
			}
		}]
	}`

	tmp, err := os.CreateTemp("", "mask_*.geojson")
	if err != nil {
		t.Fatal(err)
	}
	tmp.WriteString(geoJSON)
	tmp.Close()
	defer os.Remove(tmp.Name())

	pts := []Point3D{
		{X: 5, Y: 5, Z: 100},
		{X: -5, Y: 5, Z: 100},
	}

	mask, err := VectorMaskFilter(pts, &VectorMaskOptions{
		MaskPath: tmp.Name(),
		Invert:   true,
	})
	if err != nil {
		t.Fatalf("vector mask error: %v", err)
	}
	if !mask[0] {
		t.Error("invert: point inside polygon should be masked")
	}
	if mask[1] {
		t.Error("invert: point outside polygon should not be masked")
	}
}

func TestVectorMaskFilter_NilOpts(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0}}
	mask, err := VectorMaskFilter(pts, nil)
	if err != nil {
		t.Fatal(err)
	}
	if mask[0] {
		t.Error("nil opts should not mask")
	}
}

func TestVectorMaskFilter_EmptyPath(t *testing.T) {
	pts := []Point3D{{X: 0, Y: 0}}
	mask, err := VectorMaskFilter(pts, &VectorMaskOptions{MaskPath: ""})
	if err != nil {
		t.Fatal(err)
	}
	if mask[0] {
		t.Error("empty path should not mask")
	}
}

func TestVectorMaskFilter_EmptyPoints(t *testing.T) {
	mask, err := VectorMaskFilter(nil, &VectorMaskOptions{MaskPath: "test.geojson"})
	if err != nil {
		t.Fatal(err)
	}
	if mask != nil {
		t.Error("nil points should return nil")
	}
}
