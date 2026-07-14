package dem

import (
	"fmt"
	"strings"

	"github.com/flywave/go-geo"
)

func ParseSRS(input string) (geo.Proj, error) {
	if input == "" {
		return nil, fmt.Errorf("empty SRS string")
	}
	p := geo.NewProj(input)
	if p == nil {
		return nil, fmt.Errorf("unrecognized SRS: %s", input)
	}
	return p, nil
}

func SRSIsLatLong(p geo.Proj) bool {
	if p == nil {
		return false
	}
	return p.IsLatLong()
}

func SRSToWKT(p geo.Proj) string {
	if p == nil {
		return ""
	}
	return p.GetDef()
}

func SRSToProj4(p geo.Proj) string {
	if p == nil {
		return ""
	}
	return p.GetDef()
}

func SRSToEPSG(p geo.Proj) int {
	if p == nil {
		return 0
	}
	code := p.GetSrsCode()
	if code == "" {
		return 0
	}
	var epsg int
	if _, err := fmt.Sscanf(code, "EPSG:%d", &epsg); err != nil {
		return 0
	}
	return epsg
}

func SRSGetAuthorityCode(p geo.Proj) string {
	if p == nil {
		return ""
	}
	return p.GetSrsCode()
}

func SRSGetCSType(input string) string {
	if input == "" {
		return "UNKNOWN"
	}
	p := geo.NewProj(input)
	if p == nil {
		return "UNKNOWN"
	}
	if p.IsLatLong() {
		return "GEOGCS"
	}
	return "PROJCS"
}

func SRSIsProjected(p geo.Proj) bool {
	if p == nil {
		return false
	}
	return !p.IsLatLong()
}

func SRSEquals(a, b geo.Proj) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.GetSrsCode() == b.GetSrsCode()
}

func SRSClone(p geo.Proj) geo.Proj {
	if p == nil {
		return nil
	}
	code := p.GetSrsCode()
	if code != "" {
		return geo.NewProj(code)
	}
	return geo.NewProj(p.GetDef())
}

func SRSVerticalFromCompound(input string) string {
	if !strings.Contains(input, "+") {
		return ""
	}
	parts := strings.Split(input, "+")
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[len(parts)-1])
}

func SRSHorizontalFromCompound(input string) string {
	if !strings.Contains(input, "+") {
		return input
	}
	parts := strings.Split(input, "+")
	if len(parts) < 2 {
		return input
	}
	return strings.TrimSpace(parts[0])
}

func SRSCompound(horiz, vert string) string {
	if vert == "" {
		return horiz
	}
	return fmt.Sprintf("%s+%s", horiz, vert)
}
