package dem

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/flywave/go-geo"
	vec2d "github.com/flywave/go3d/float64/vec2"
)

type Region struct {
	Extent *geo.MapExtent
	XRes   float64
	YRes   float64
	XSize  int
	YSize  int
	ZMin   float64
	ZMax   float64
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
	parts := strings.FieldsFunc(regionStr, func(r rune) bool {
		return r == '/' || r == ',' || r == ' '
	})
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid region format, expected xmin/xmax/ymin/ymax: %s", regionStr)
	}
	minX, _ := strconv.ParseFloat(parts[0], 64)
	maxX, _ := strconv.ParseFloat(parts[1], 64)
	minY, _ := strconv.ParseFloat(parts[2], 64)
	maxY, _ := strconv.ParseFloat(parts[3], 64)
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
		ZMin:   r.ZMin,
		ZMax:   r.ZMax,
	}
}

func (r *Region) String() string {
	return fmt.Sprintf("Region[%s] %.6f/%.6f/%.6f/%.6f res=%.6f x %d y %d",
		r.SRS().GetSrsCode(),
		r.Extent.BBox.Min[0], r.Extent.BBox.Min[1],
		r.Extent.BBox.Max[0], r.Extent.BBox.Max[1],
		r.XRes, r.XSize, r.YSize)
}

func (r *Region) IsValid() bool {
	b := r.Extent.BBox
	return b.Max[0] > b.Min[0] && b.Max[1] > b.Min[1] && r.XRes > 0 && r.YRes > 0 && r.XSize > 0 && r.YSize > 0
}

func (r *Region) Copy() *Region {
	bbox := vec2d.Rect{Min: r.Extent.BBox.Min, Max: r.Extent.BBox.Max}
	extent := &geo.MapExtent{BBox: bbox, Srs: r.Extent.Srs}
	return &Region{
		Extent: extent,
		XRes:   r.XRes,
		YRes:   r.YRes,
		XSize:  r.XSize,
		YSize:  r.YSize,
		ZMin:   r.ZMin,
		ZMax:   r.ZMax,
	}
}

func (r *Region) Center() (float64, float64) {
	b := r.Extent.BBox
	return (b.Min[0] + b.Max[0]) / 2, (b.Min[1] + b.Max[1]) / 2
}

func (r *Region) Width() float64 {
	return r.Extent.BBox.Max[0] - r.Extent.BBox.Min[0]
}

func (r *Region) Height() float64 {
	return r.Extent.BBox.Max[1] - r.Extent.BBox.Min[1]
}

func (r *Region) ZRange() (float64, float64) {
	return r.ZMin, r.ZMax
}

func (r *Region) SetZRange(zMin, zMax float64) {
	r.ZMin, r.ZMax = zMin, zMax
}

func (r *Region) Buffer(xBuf, yBuf float64) *Region {
	if xBuf == 0 && yBuf == 0 {
		return r.Copy()
	}
	if yBuf == 0 {
		yBuf = xBuf
	}
	n := r.Copy()
	n.Extent.BBox.Min[0] -= xBuf
	n.Extent.BBox.Max[0] += xBuf
	n.Extent.BBox.Min[1] -= yBuf
	n.Extent.BBox.Max[1] += yBuf
	return n
}

func (r *Region) BufferPct(pct float64) *Region {
	xBuf := r.Width() * pct / 100
	yBuf := r.Height() * pct / 100
	return r.Buffer(xBuf, yBuf)
}

func (r *Region) Intersection(other *Region) *Region {
	b := r.Extent.BBox
	ob := other.Extent.BBox
	n := r.Copy()
	n.Extent.BBox.Min[0] = math.Max(b.Min[0], ob.Min[0])
	n.Extent.BBox.Max[0] = math.Min(b.Max[0], ob.Max[0])
	n.Extent.BBox.Min[1] = math.Max(b.Min[1], ob.Min[1])
	n.Extent.BBox.Max[1] = math.Min(b.Max[1], ob.Max[1])
	return n
}

func (r *Region) Union(other *Region) *Region {
	b := r.Extent.BBox
	ob := other.Extent.BBox
	n := r.Copy()
	n.Extent.BBox.Min[0] = math.Min(b.Min[0], ob.Min[0])
	n.Extent.BBox.Max[0] = math.Max(b.Max[0], ob.Max[0])
	n.Extent.BBox.Min[1] = math.Min(b.Min[1], ob.Min[1])
	n.Extent.BBox.Max[1] = math.Max(b.Max[1], ob.Max[1])
	return n
}

func (r *Region) Intersects(other *Region) bool {
	b, ob := r.Extent.BBox, other.Extent.BBox
	return b.Min[0] <= ob.Max[0] && b.Max[0] >= ob.Min[0] &&
		b.Min[1] <= ob.Max[1] && b.Max[1] >= ob.Min[1]
}

func (r *Region) Contains(other *Region) bool {
	b, ob := r.Extent.BBox, other.Extent.BBox
	return b.Min[0] <= ob.Min[0] && b.Max[0] >= ob.Max[0] &&
		b.Min[1] <= ob.Min[1] && b.Max[1] >= ob.Max[1]
}

func (r *Region) SrcWin(gt [6]float64, xSize, ySize int) (int, int, int, int) {
	pxMin := int(math.Floor((r.Extent.BBox.Min[0] - gt[0]) / gt[1]))
	pxMax := int(math.Ceil((r.Extent.BBox.Max[0] - gt[0]) / gt[1]))
	pyMin := int(math.Floor((r.Extent.BBox.Max[1] - gt[3]) / gt[5]))
	pyMax := int(math.Ceil((r.Extent.BBox.Min[1] - gt[3]) / gt[5]))

	if pxMin < 0 {
		pxMin = 0
	}
	if pyMin < 0 {
		pyMin = 0
	}
	if pxMax > xSize {
		pxMax = xSize
	}
	if pyMax > ySize {
		pyMax = ySize
	}
	return pxMin, pyMin, pxMax - pxMin, pyMax - pyMin
}

func (r *Region) GeoTransformFromCount(xCount, yCount int) [6]float64 {
	b := r.Extent.BBox
	xInc := (b.Max[0] - b.Min[0]) / float64(xCount)
	yInc := (b.Min[1] - b.Max[1]) / float64(yCount)
	return [6]float64{b.Min[0], xInc, 0, b.Max[1], 0, yInc}
}

func (r *Region) Format(fmtStr string) string {
	b := r.Extent.BBox
	switch fmtStr {
	case "gmt", "str":
		return fmt.Sprintf("%.6f/%.6f/%.6f/%.6f", b.Min[0], b.Max[0], b.Min[1], b.Max[1])
	case "sstr":
		return fmt.Sprintf("%.10f/%.10f/%.10f/%.10f", b.Min[0], b.Max[0], b.Min[1], b.Max[1])
	case "fstr":
		return fmt.Sprintf("%f/%f/%f/%f", b.Min[0], b.Max[0], b.Min[1], b.Max[1])
	case "bbox":
		return fmt.Sprintf("%.6f,%.6f,%.6f,%.6f", b.Min[0], b.Min[1], b.Max[0], b.Max[1])
	case "te":
		return fmt.Sprintf("%.6f %.6f %.6f %.6f", b.Min[0], b.Min[1], b.Max[0], b.Max[1])
	case "fn":
		return fmt.Sprintf("%.2f_%.2f_%.2f_%.2f", b.Min[0], b.Max[0], b.Min[1], b.Max[1])
	case "fn_full":
		return fmt.Sprintf("%.6f_%.6f_%.6f_%.6f", b.Min[0], b.Max[0], b.Min[1], b.Max[1])
	case "wkt":
		return fmt.Sprintf("POLYGON((%.6f %.6f,%.6f %.6f,%.6f %.6f,%.6f %.6f,%.6f %.6f))",
			b.Min[0], b.Min[1], b.Max[0], b.Min[1],
			b.Max[0], b.Max[1], b.Min[0], b.Max[1],
			b.Min[0], b.Min[1])
	default:
		return fmt.Sprintf("%.6f/%.6f/%.6f/%.6f", b.Min[0], b.Max[0], b.Min[1], b.Max[1])
	}
}

func (r *Region) Round(decimals int) *Region {
	n := r.Copy()
	factor := math.Pow(10, float64(decimals))
	n.Extent.BBox.Min[0] = math.Floor(n.Extent.BBox.Min[0]*factor) / factor
	n.Extent.BBox.Max[0] = math.Ceil(n.Extent.BBox.Max[0]*factor) / factor
	n.Extent.BBox.Min[1] = math.Floor(n.Extent.BBox.Min[1]*factor) / factor
	n.Extent.BBox.Max[1] = math.Ceil(n.Extent.BBox.Max[1]*factor) / factor
	return n
}

func (r *Region) Chunk(xChunks, yChunks int) []*Region {
	if xChunks <= 0 {
		xChunks = 1
	}
	if yChunks <= 0 {
		yChunks = 1
	}
	b := r.Extent.BBox
	xStep := (b.Max[0] - b.Min[0]) / float64(xChunks)
	yStep := (b.Max[1] - b.Min[1]) / float64(yChunks)

	var chunks []*Region
	for yi := 0; yi < yChunks; yi++ {
		for xi := 0; xi < xChunks; xi++ {
			sub := r.Copy()
			sub.Extent.BBox.Min[0] = b.Min[0] + float64(xi)*xStep
			sub.Extent.BBox.Max[0] = b.Min[0] + float64(xi+1)*xStep
			sub.Extent.BBox.Min[1] = b.Min[1] + float64(yi)*yStep
			sub.Extent.BBox.Max[1] = b.Min[1] + float64(yi+1)*yStep
			chunks = append(chunks, sub)
		}
	}
	return chunks
}

func (r *Region) ExportAsWKT() string {
	return r.Format("wkt")
}
