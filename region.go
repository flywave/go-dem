package dem

import (
	"fmt"
	"math"

	"github.com/flywave/go-geo"
	vec2d "github.com/flywave/go3d/float64/vec2"
)

type Region struct {
	Extent *geo.MapExtent
	XRes   float64
	YRes   float64
	XSize  int
	YSize  int
}

func NewRegionFromBBox(minX, minY, maxX, maxY float64, srs geo.Proj, xRes, yRes float64) *Region {
	if yRes == 0 {
		yRes = xRes
	}
	bbox := vec2d.Rect{Min: vec2d.T{minX, minY}, Max: vec2d.T{maxX, maxY}}
	extent := &geo.MapExtent{BBox: bbox, Srs: srs}

	xSize := int(math.Round((maxX - minX) / xRes))
	ySize := int(math.Round((maxY - minY) / yRes))

	return &Region{
		Extent: extent,
		XRes:   xRes,
		YRes:   yRes,
		XSize:  xSize,
		YSize:  ySize,
	}
}

func NewRegionFromString(regionStr string, srs geo.Proj, xRes, yRes float64) (*Region, error) {
	var minX, minY, maxX, maxY float64
	n, err := fmt.Sscanf(regionStr, "%f/%f/%f/%f", &minX, &maxX, &minY, &maxY)
	if err != nil || n != 4 {
		return nil, fmt.Errorf("invalid region format, expected xmin/xmax/ymin/ymax: %s", regionStr)
	}
	return NewRegionFromBBox(minX, minY, maxX, maxY, srs, xRes, yRes), nil
}

func (r *Region) BBox() vec2d.Rect {
	return r.Extent.BBox
}

func (r *Region) SRS() geo.Proj {
	return r.Extent.Srs
}

func (r *Region) GeoTransform() [6]float64 {
	return [6]float64{
		r.Extent.BBox.Min[0],
		r.XRes,
		0,
		r.Extent.BBox.Max[1],
		0,
		-r.YRes,
	}
}

func (r *Region) TransformTo(srs geo.Proj) *Region {
	newBBox := r.Extent.BBoxFor(srs)
	newExtent := &geo.MapExtent{BBox: newBBox, Srs: srs}

	xRes := r.XRes
	yRes := r.YRes

	if srs != nil && !srs.IsLatLong() && r.SRS().IsLatLong() {
		midY := (newBBox.Min[1] + newBBox.Max[1]) / 2
		midYRad := midY * math.Pi / 180
		cosLat := math.Cos(midYRad)
		if cosLat > 0 {
			xRes = r.XRes * 111320 * cosLat
		}
		yRes = r.YRes * 111320
	}

	xSize := int(math.Round((newBBox.Max[0] - newBBox.Min[0]) / xRes))
	ySize := int(math.Round((newBBox.Max[1] - newBBox.Min[1]) / yRes))

	return &Region{
		Extent: newExtent,
		XRes:   xRes,
		YRes:   yRes,
		XSize:  xSize,
		YSize:  ySize,
	}
}

func (r *Region) String() string {
	return fmt.Sprintf("Region[%s] %.6f/%.6f/%.6f/%.6f res=%.6f x %d y %d",
		r.SRS().GetSrsCode(),
		r.Extent.BBox.Min[0], r.Extent.BBox.Min[1],
		r.Extent.BBox.Max[0], r.Extent.BBox.Max[1],
		r.XRes, r.XSize, r.YSize)
}
