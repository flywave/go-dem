package datum

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-geo"
	"github.com/flywave/go-geoid"
)

func testRegion() *dem.Region {
	return dem.NewRegionFromBBox(-122.5, 40.0, -122.0, 40.5, geo.NewProj("EPSG:4326"), 0.1, 0.1)
}

func TestComputeGeoidGridRowOrder(t *testing.T) {
	region := testRegion()
	grid := computeGeoidGrid(region, geoid.EGM96)
	g := geoid.NewGeoid(geoid.EGM96, true)
	bbox := region.BBox()

	northLat := bbox.Max[1] - 0.5*region.YRes
	southLat := bbox.Min[1] + 0.5*region.YRes
	westLon := bbox.Min[0] + 0.5*region.XRes

	firstRow := g.GetHeight(northLat, westLon)
	lastRow := g.GetHeight(southLat, westLon)

	if math.Abs(grid[0]-firstRow) > 1e-9 {
		t.Errorf("row 0 must sample the northernmost pixel center: grid[0]=%f, GetHeight(lat=%f)=%f", grid[0], northLat, firstRow)
	}
	lastIdx := (region.YSize - 1) * region.XSize
	if math.Abs(grid[lastIdx]-lastRow) > 1e-9 {
		t.Errorf("last row must sample the southernmost pixel center: grid[%d]=%f, GetHeight(lat=%f)=%f", lastIdx, grid[lastIdx], southLat, lastRow)
	}
	if math.Abs(northLat-southLat) < 1e-9 {
		t.Fatal("test region has no north-south extent")
	}

	for y := 0; y < region.YSize; y++ {
		lat := bbox.Max[1] - (float64(y)+0.5)*region.YRes
		for x := 0; x < region.XSize; x++ {
			lon := bbox.Min[0] + (float64(x)+0.5)*region.XRes
			want := g.GetHeight(lat, lon)
			if math.IsNaN(want) {
				t.Fatalf("reference geoid height NaN at lat=%f lon=%f", lat, lon)
			}
			got := grid[y*region.XSize+x]
			if math.Abs(got-want) > 1e-9 {
				t.Errorf("grid[%d,%d]=%f, want %f (lat=%f lon=%f)", x, y, got, want, lat, lon)
			}
		}
	}
}

func TestGenerateGeoidGrid(t *testing.T) {
	region := dem.NewRegionFromBBox(-122.5, 40.0, -122.3, 40.2, geo.NewProj("EPSG:4326"), 0.1, 0.1)
	vg, err := GenerateGeoidGrid(region, geoid.EGM96)
	if err != nil {
		t.Fatal(err)
	}
	direct := computeGeoidGrid(region, geoid.EGM96)
	for i := range direct {
		if vg.Data[i] != direct[i] {
			t.Fatalf("GenerateGeoidGrid data mismatch at %d: %f vs %f", i, vg.Data[i], direct[i])
		}
	}
	wantUnc := geoidUncertainty(geoid.EGM96)
	if wantUnc <= 0 {
		t.Fatalf("expected positive geoid uncertainty, got %f", wantUnc)
	}
	for i, u := range vg.Uncertainty {
		if u != wantUnc {
			t.Fatalf("Uncertainty[%d]=%f, want %f", i, u, wantUnc)
		}
	}
}

func TestConvertHeightKnownGeoidOffset(t *testing.T) {
	region := dem.NewRegionFromBBox(-122.5, 40.0, -122.4, 40.1, geo.NewProj("EPSG:4326"), 0.05, 0.05)
	data := []float64{100, 200, 300, 400}

	ortho, err := ConvertHeight(data, region, &ConvertOptions{From: HeightEllipsoidal, To: HeightOrthometric, Model: geoid.EGM96, Cubic: true})
	if err != nil {
		t.Fatal(err)
	}

	g := geoid.NewGeoid(geoid.EGM96, true)
	gt := region.GeoTransform()
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			idx := y*2 + x
			lat := gt[3] + (float64(y)+0.5)*gt[5]
			lon := gt[0] + (float64(x)+0.5)*gt[1]
			n := g.GetHeight(lat, lon)
			want := data[idx] - n
			if math.Abs(ortho[idx]-want) > 1e-6 {
				t.Errorf("ortho[%d]=%f, want %f (N=%f at lat=%f lon=%f)", idx, ortho[idx], want, n, lat, lon)
			}
		}
	}

	back, err := ConvertHeight(ortho, region, &ConvertOptions{From: HeightOrthometric, To: HeightEllipsoidal, Model: geoid.EGM96, Cubic: true})
	if err != nil {
		t.Fatal(err)
	}
	for i := range data {
		if math.Abs(back[i]-data[i]) > 1e-6 {
			t.Errorf("round trip back[%d]=%f, want %f", i, back[i], data[i])
		}
	}
}

func TestConvertHeightNoDataAndNaN(t *testing.T) {
	region := dem.NewRegionFromBBox(-122.5, 40.0, -122.4, 40.1, geo.NewProj("EPSG:4326"), 0.05, 0.05)
	data := []float64{100, dem.DefaultNoData, math.NaN(), 400}
	out, err := ConvertHeight(data, region, &ConvertOptions{From: HeightEllipsoidal, To: HeightOrthometric, Model: geoid.EGM96})
	if err != nil {
		t.Fatal(err)
	}
	if out[1] != dem.DefaultNoData {
		t.Errorf("nodata pixel must stay nodata, got %f", out[1])
	}
	if out[2] != dem.DefaultNoData {
		t.Errorf("NaN pixel must become nodata, got %f", out[2])
	}
}

func TestConvertHeightUnsupportedCombinations(t *testing.T) {
	region := dem.NewRegionFromBBox(-122.5, 40.0, -122.4, 40.1, geo.NewProj("EPSG:4326"), 0.1, 0.1)
	data := []float64{100, 200, 300, 400}

	if _, err := ConvertHeight(data, region, nil); err == nil {
		t.Error("nil options must return an error")
	}
	if _, err := ConvertHeight(data, region, &ConvertOptions{From: HeightEllipsoidal, To: HeightGeoid, Model: geoid.EGM96}); err == nil {
		t.Error("ellipsoidal→geoid must return an error")
	}
	if _, err := ConvertHeight(data, region, &ConvertOptions{From: HeightGeoid, To: HeightEllipsoidal, Model: geoid.EGM96}); err == nil {
		t.Error("geoid→ellipsoidal must return an error")
	}
	if _, err := ConvertHeight(data, region, &ConvertOptions{From: HeightGeoid, To: HeightOrthometric, Model: geoid.EGM96}); err == nil {
		t.Error("geoid→orthometric must return an error")
	}
	if _, err := ConvertHeight(data, region, &ConvertOptions{From: HeightEllipsoidal, To: HeightOrthometric, Model: geoid.HAE}); err == nil {
		t.Error("conversion without a geoid model must return an error")
	}
	if _, err := ConvertHeight(data, region, &ConvertOptions{From: HeightEllipsoidal, To: HeightOrthometric, Model: geoid.UNKNOWN}); err == nil {
		t.Error("conversion with UNKNOWN model must return an error")
	}
	if _, err := ConvertHeight(data[:3], testRegion(), &ConvertOptions{From: HeightEllipsoidal, To: HeightOrthometric, Model: geoid.EGM96}); err == nil {
		t.Error("short data must return an error")
	}

	same, err := ConvertHeight(data, region, &ConvertOptions{From: HeightGeoid, To: HeightGeoid, Model: geoid.EGM96})
	if err != nil {
		t.Errorf("identical from/to must be a no-op without error, got %v", err)
	}
	for i := range data {
		if same[i] != data[i] {
			t.Errorf("identity conversion changed data[%d]: %f → %f", i, data[i], same[i])
		}
	}
}

func TestVerticalTransformMSLToEGM96IsZeroGrid(t *testing.T) {
	vt := NewVerticalTransform(TransformOptions{EpsgIn: mslEPSG, EpsgOut: 5773, Region: testRegion()})
	res, err := vt.Run()
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range res.Grid {
		if math.Abs(v) > 1e-9 {
			t.Fatalf("MSL→EGM96 is a model-identity conversion, grid must be zero, got %f at %d", v, i)
		}
	}
	if res.EpsgOut != 5773 {
		t.Errorf("EpsgOut=%d, want 5773", res.EpsgOut)
	}
}

func TestVerticalTransformMSLToEGM2008IsModelDifference(t *testing.T) {
	region := testRegion()
	vt := NewVerticalTransform(TransformOptions{EpsgIn: mslEPSG, EpsgOut: 3855, Region: region})
	res, err := vt.Run()
	if err != nil {
		t.Fatal(err)
	}

	n2008 := computeGeoidGrid(region, geoid.EGM2008)
	n96 := computeGeoidGrid(region, geoid.EGM96)
	maxAbs := 0.0
	for i := range n96 {
		want := n2008[i] - n96[i]
		if math.Abs(res.Grid[i]-want) > 1e-9 {
			t.Fatalf("grid[%d]=%f, want model difference %f", i, res.Grid[i], want)
		}
		if math.Abs(want) > maxAbs {
			maxAbs = math.Abs(want)
		}
	}
	if maxAbs > 10 {
		t.Errorf("MSL→EGM2008 offset must be a small inter-model difference, got max |ΔN|=%f", maxAbs)
	}
}

func TestVerticalTransformTidalCombinationsReturnError(t *testing.T) {
	region := testRegion()
	pairs := [][2]int{
		{5866, mslEPSG},
		{mslEPSG, 5866},
		{5866, 5868},
		{5866, 5703},
		{5703, 5866},
		{5866, 4269},
		{4269, 5866},
	}
	for _, p := range pairs {
		vt := NewVerticalTransform(TransformOptions{EpsgIn: p[0], EpsgOut: p[1], Region: region})
		res, err := vt.Run()
		if err == nil {
			t.Errorf("%d→%d must return an error (no tidal offset grids), got grid[0]=%f", p[0], p[1], res.Grid[0])
		}
		if res != nil {
			t.Errorf("%d→%d must not return a result on error", p[0], p[1])
		}
	}
}

func TestVerticalTransformUnknownFramesReturnError(t *testing.T) {
	region := testRegion()
	for _, p := range [][2]int{{4326, 5773}, {5773, 4326}, {4269, 4326}, {0, 5773}} {
		vt := NewVerticalTransform(TransformOptions{EpsgIn: p[0], EpsgOut: p[1], Region: region})
		if _, err := vt.Run(); err == nil {
			t.Errorf("%d→%d must return an error for unregistered frames", p[0], p[1])
		}
	}
	if _, err := TransformDEM([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, region, 4326, 5773); err == nil {
		t.Error("TransformDEM must propagate the unknown-frame error")
	}
}

func TestVerticalTransformUnresolvableCDNReturnsError(t *testing.T) {
	vt := NewVerticalTransform(TransformOptions{EpsgIn: 6647, EpsgOut: 5703, Region: testRegion()})
	if _, err := vt.Run(); err == nil {
		t.Error("CGVD2013(6647)→NAVD88(5703) has no resolvable models and must return an error")
	}
}

func TestVerticalTransformCDNPairIsModelDifference(t *testing.T) {
	region := testRegion()
	vt := NewVerticalTransform(TransformOptions{EpsgIn: 5773, EpsgOut: 3855, Region: region})
	res, err := vt.Run()
	if err != nil {
		t.Fatal(err)
	}
	n2008 := computeGeoidGrid(region, geoid.EGM2008)
	n96 := computeGeoidGrid(region, geoid.EGM96)
	for i := range n96 {
		if math.Abs(res.Grid[i]-(n2008[i]-n96[i])) > 1e-9 {
			t.Fatalf("EGM96→EGM2008 grid[%d]=%f, want N_to−N_from=%f", i, res.Grid[i], n2008[i]-n96[i])
		}
	}
}

func TestVerticalTransformGeoidStepsUseSignedUndulation(t *testing.T) {
	region := testRegion()
	vt := NewVerticalTransform(TransformOptions{EpsgIn: 5773, EpsgOut: ellipsoidEPSG, Region: region})
	n96 := computeGeoidGrid(region, geoid.EGM96)

	cdn2ell, _, err := vt.executeStep(transformStep{from: 5773, to: ellipsoidEPSG, via: "cdn2ellipsoid"})
	if err != nil {
		t.Fatal(err)
	}
	for i := range n96 {
		if math.Abs(cdn2ell[i]-(-n96[i])) > 1e-9 {
			t.Fatalf("cdn2ellipsoid grid[%d]=%f, want -N=%f", i, cdn2ell[i], -n96[i])
		}
	}

	ell2cdn, _, err := vt.executeStep(transformStep{from: ellipsoidEPSG, to: 5773, via: "ellipsoid2cdn"})
	if err != nil {
		t.Fatal(err)
	}
	for i := range n96 {
		if math.Abs(ell2cdn[i]-n96[i]) > 1e-9 {
			t.Fatalf("ellipsoid2cdn grid[%d]=%f, want +N=%f", i, ell2cdn[i], n96[i])
		}
	}

	msl2ell, _, err := vt.executeStep(transformStep{from: mslEPSG, to: ellipsoidEPSG, via: "msl2ellipsoid"})
	if err != nil {
		t.Fatal(err)
	}
	for i := range n96 {
		if math.Abs(msl2ell[i]-(-n96[i])) > 1e-9 {
			t.Fatalf("msl2ellipsoid grid[%d]=%f, want -N=%f", i, msl2ell[i], -n96[i])
		}
	}
}

func TestPlanStepsChains(t *testing.T) {
	cases := []struct {
		in  int
		out int
		via []string
	}{
		{4269, mslEPSG, []string{"htdp2ellipsoid", "ellipsoid2msl"}},
		{mslEPSG, 4269, []string{"msl2ellipsoid", "ellipsoid2htdp"}},
		{5773, 4269, []string{"cdn2ellipsoid", "ellipsoid2htdp"}},
		{4269, 5773, []string{"htdp2ellipsoid", "ellipsoid2cdn"}},
		{4269, 7912, []string{"htdp2htdp"}},
		{5773, 3855, []string{"cdn2cdn"}},
		{mslEPSG, 3855, []string{"msl2cdn"}},
		{3855, mslEPSG, []string{"cdn2msl"}},
	}
	for _, c := range cases {
		steps, err := planSteps(c.in, c.out, GetFrameByEPSG(c.in), GetFrameByEPSG(c.out))
		if err != nil {
			t.Fatalf("%d→%d: unexpected planning error: %v", c.in, c.out, err)
		}
		if len(steps) != len(c.via) {
			t.Fatalf("%d→%d: got %d steps, want %d", c.in, c.out, len(steps), len(c.via))
		}
		for i, s := range steps {
			if s.via != c.via[i] {
				t.Errorf("%d→%d step %d via=%q, want %q", c.in, c.out, i, s.via, c.via[i])
			}
			if i > 0 && s.from != steps[i-1].to {
				t.Errorf("%d→%d broken chain at step %d: from=%d, previous to=%d", c.in, c.out, i, s.from, steps[i-1].to)
			}
		}
		if steps[len(steps)-1].to != c.out {
			t.Errorf("%d→%d chain ends at %d, want %d", c.in, c.out, steps[len(steps)-1].to, c.out)
		}
	}

	if _, err := planSteps(5866, mslEPSG, GetFrameByEPSG(5866), GetFrameByEPSG(mslEPSG)); err == nil {
		t.Error("MLLW→MSL planning must return an error")
	}
	if _, err := planSteps(mslEPSG, 5866, GetFrameByEPSG(mslEPSG), GetFrameByEPSG(5866)); err == nil {
		t.Error("MSL→MLLW planning must return an error")
	}
}

func TestVerticalTransformUncertaintyStacking(t *testing.T) {
	region := testRegion()

	vt := NewVerticalTransform(TransformOptions{EpsgIn: 5773, EpsgOut: 5703, Region: region})
	res, err := vt.Run()
	if err != nil {
		t.Fatal(err)
	}
	want := math.Hypot(FrameUncertainty(5703), geoidUncertainty(geoid.EGM96, geoid.EGM96))
	if math.Abs(res.Uncertainty[0]-want) > 1e-9 {
		t.Errorf("uncertainty=%f, want frame(0.05) RSS geoid(%f)=%f", res.Uncertainty[0], geoidUncertainty(geoid.EGM96, geoid.EGM96), want)
	}
	if res.Uncertainty[0] <= FrameUncertainty(5703) {
		t.Errorf("geoid uncertainty must be stacked on the frame uncertainty: %f <= %f", res.Uncertainty[0], FrameUncertainty(5703))
	}

	vt2 := NewVerticalTransform(TransformOptions{EpsgIn: mslEPSG, EpsgOut: 5703, Region: region})
	res2, err := vt2.Run()
	if err != nil {
		t.Fatal(err)
	}
	want2 := math.Hypot(FrameUncertainty(5703), geoidUncertainty(geoid.EGM96))
	if math.Abs(res2.Uncertainty[0]-want2) > 1e-9 {
		t.Errorf("uncertainty=%f, want %f (base must come from frameIn only, output frame counted once)", res2.Uncertainty[0], want2)
	}
}

func TestVerticalTransformHTDPStubGuard(t *testing.T) {
	region := testRegion()
	vt := NewVerticalTransform(TransformOptions{EpsgIn: 4269, EpsgOut: 7912, Region: region})
	t.Logf("htdpIsStub=%v", htdpIsStub())

	if htdpIsStub() {
		grid, err := vt.htdpGrid(4269, 7912)
		if err == nil {
			t.Fatal("stub htdp implementation must return an error instead of a zero grid")
		}
		if grid != nil {
			t.Fatal("stub htdp implementation must not return a grid")
		}
		if _, err := vt.Run(); err == nil {
			t.Fatal("stub htdp chain must fail with an error, not silently return a zero grid")
		}
	} else if _, err := vt.htdpGrid(4269, 7912); err != nil {
		t.Fatalf("real htdp implementation should not fail: %v", err)
	}
}

func TestGenerateTransformGridModelWiring(t *testing.T) {
	region := testRegion()

	vg, err := GenerateTransformGrid(region, 6647, 5703, geoid.EGM96)
	if err != nil {
		t.Fatalf("explicit model must resolve the unresolvable datum 6647: %v", err)
	}
	for i, v := range vg.Data {
		if math.Abs(v) > 1e-9 {
			t.Fatalf("EGM96-resolved 6647→5703 must be a zero grid, got %f at %d", v, i)
		}
	}

	if _, err := GenerateTransformGrid(region, 6647, 5703, geoid.HAE); err == nil {
		t.Error("6647→5703 without a usable model must return an error")
	}

	if _, err := GenerateTransformGrid(region, 4326, 5773, geoid.EGM96); err == nil {
		t.Error("unknown frame must return an error even with a model")
	}
}

func TestGetFrameByNameDeterministic(t *testing.T) {
	first := GetFrameByName("mhw")
	if first == nil {
		t.Fatal("mhw not found")
	}
	if first.EPSG != 5868 || first.Name != "mhw" {
		t.Errorf(`exact match "mhw" must return EPSG 5868, got %d (%s)`, first.EPSG, first.Name)
	}
	for i := 0; i < 100; i++ {
		f := GetFrameByName("mhw")
		if f.EPSG != first.EPSG || f.Name != first.Name {
			t.Fatalf("GetFrameByName is not deterministic: %d(%s) then %d(%s)", first.EPSG, first.Name, f.EPSG, f.Name)
		}
	}

	if f := GetFrameByName("egm96"); f == nil || f.EPSG != 5773 {
		t.Errorf(`substring "egm96" must resolve to EPSG 5773, got %v`, f)
	}
	if f := GetFrameByName("no-such-frame"); f != nil {
		t.Errorf("unknown name must return nil, got %v", f)
	}
	if GetFrameByEPSG(0) != nil {
		t.Error("EPSG 0 is not a registered datum and must return nil")
	}
}

func TestVerticalTransformSameEpsgNoOp(t *testing.T) {
	vt := NewVerticalTransform(TransformOptions{EpsgIn: 5773, EpsgOut: 5773, Region: testRegion()})
	res, err := vt.Run()
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range res.Grid {
		if v != 0 {
			t.Fatalf("identity transform must return a zero grid, got %f at %d", v, i)
		}
	}
}

func TestVerticalTransformNilRegion(t *testing.T) {
	vt := NewVerticalTransform(TransformOptions{EpsgIn: 5773, EpsgOut: 3855})
	if _, err := vt.Run(); err == nil || !strings.Contains(fmt.Sprint(err), "region") {
		t.Errorf("nil region must return an error, got %v", err)
	}
}

func TestGeoidUncertaintyTable(t *testing.T) {
	if geoidUncertainty(geoid.EGM96) != GetGeoidUncertainty("g1999") {
		t.Error("EGM96 must map to its GeoidModels entry")
	}
	if geoidUncertainty(geoid.EGM2008) != GetGeoidUncertainty("g2012a") {
		t.Error("EGM2008 must map to its GeoidModels entry")
	}
	if geoidUncertainty(geoid.EGM84) != GetGeoidUncertainty("geoid09") {
		t.Error("EGM84 must map to its GeoidModels entry")
	}
	if GetGeoidUncertainty("unknown-model") != 0 {
		t.Error("unknown model must have zero uncertainty")
	}
}
