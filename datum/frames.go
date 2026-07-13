package datum

import (
	"fmt"
	"strings"
)

type FrameType int

const (
	FrameTidal FrameType = iota
	FrameHTDP
	FrameCDN
	FrameGeoid
	FrameEllipsoid
)

type Frame struct {
	EPSG        int
	Name        string
	Description string
	Type        FrameType
	Uncertainty float64
	HTDPID      int
	Epoch       float64
}

var TidalFrames = map[int]Frame{
	1089: {EPSG: 5866, Name: "mllw", Description: "Mean Lower Low Water", Type: FrameTidal, Uncertainty: 0},
	5866: {EPSG: 5866, Name: "mllw", Description: "Mean Lower Low Water", Type: FrameTidal, Uncertainty: 0},
	1091: {EPSG: 1091, Name: "mlw", Description: "Mean Low Water", Type: FrameTidal, Uncertainty: 0},
	1090: {EPSG: 5869, Name: "mhhw", Description: "Mean Higher High Water", Type: FrameTidal, Uncertainty: 0},
	5869: {EPSG: 5869, Name: "mhhw", Description: "Mean Higher High Water", Type: FrameTidal, Uncertainty: 0},
	5868: {EPSG: 5868, Name: "mhw", Description: "Mean High Water", Type: FrameTidal, Uncertainty: 0},
	5714: {EPSG: 5714, Name: "msl", Description: "Mean Sea Level", Type: FrameTidal, Uncertainty: 0},
	5713: {EPSG: 5713, Name: "mtl", Description: "Mean Tide Level", Type: FrameTidal, Uncertainty: 0},
	0:    {EPSG: 0, Name: "crd", Description: "Columbia River Datum", Type: FrameTidal, Uncertainty: 0},
	1:    {EPSG: 1, Name: "xgeoid20b", Description: "xgeoid 2020 B", Type: FrameTidal, Uncertainty: 0},
	7968: {EPSG: 7968, Name: "NGVD", Description: "National Geodetic Vertical Datum", Type: FrameTidal, Uncertainty: 0},
}

var HTDPFrames = map[int]Frame{
	4269: {EPSG: 4269, Name: "NAD_83(2011/CORS96/2007)", Description: "(North American plate fixed)", Type: FrameHTDP, Uncertainty: 0.02, HTDPID: 1, Epoch: 1997.0},
	6781: {EPSG: 6781, Name: "NAD_83(2011/CORS96/2007)", Description: "(North American plate fixed)", Type: FrameHTDP, Uncertainty: 0.02, HTDPID: 1, Epoch: 1997.0},
	6319: {EPSG: 6319, Name: "NAD_83(2011/CORS96/2007)", Description: "(North American plate fixed)", Type: FrameHTDP, Uncertainty: 0.02, HTDPID: 1, Epoch: 1997.0},
	6321: {EPSG: 6321, Name: "NAD_83(PA11/PACP00)", Description: "(Pacific plate fixed)", Type: FrameHTDP, Uncertainty: 0.02, HTDPID: 2, Epoch: 1997.0},
	6324: {EPSG: 6324, Name: "NAD_83(MA11/MARP00)", Description: "(Mariana plate fixed)", Type: FrameHTDP, Uncertainty: 0.02, HTDPID: 3, Epoch: 1997.0},
	4979: {EPSG: 4979, Name: "WGS_84(original)", Description: "(NAD_83(2011) used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 4, Epoch: 1997.0},
	7815: {EPSG: 7815, Name: "WGS_84(original)", Description: "(NAD_83(2011) used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 4, Epoch: 1997.0},
	7816: {EPSG: 7816, Name: "WGS_84(original)", Description: "(NAD_83(2011) used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 4, Epoch: 1997.0},
	7656: {EPSG: 7656, Name: "WGS_84(G730)", Description: "(ITRF91 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 5, Epoch: 1997.0},
	7657: {EPSG: 7657, Name: "WGS_84(G730)", Description: "(ITRF91 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 5, Epoch: 1997.0},
	7658: {EPSG: 7658, Name: "WGS_84(G873)", Description: "(ITRF94 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 6, Epoch: 1997.0},
	7659: {EPSG: 7659, Name: "WGS_84(G873)", Description: "(ITRF94 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 6, Epoch: 1997.0},
	7660: {EPSG: 7660, Name: "WGS_84(G1150)", Description: "(ITRF2000 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 7, Epoch: 1997.0},
	7661: {EPSG: 7661, Name: "WGS_84(G1150)", Description: "(ITRF2000 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 7, Epoch: 1997.0},
	7662: {EPSG: 7662, Name: "WGS_84(G1674)", Description: "(ITRF2008 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 8, Epoch: 2000.0},
	7663: {EPSG: 7663, Name: "WGS_84(G1674)", Description: "(ITRF2008 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 8, Epoch: 2000.0},
	7664: {EPSG: 7664, Name: "WGS_84(G1762)", Description: "(IGb08 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 9, Epoch: 2000.0},
	7665: {EPSG: 7665, Name: "WGS_84(G1762)", Description: "(IGb08 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 9, Epoch: 2000.0},
	7666: {EPSG: 7666, Name: "WGS_84(G2139)", Description: "(ITRF2014=IGS14=IGb14 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 10, Epoch: 1997.0},
	7667: {EPSG: 7667, Name: "WGS_84(G2139)", Description: "(ITRF2014=IGS14=IGb14 used)", Type: FrameHTDP, Uncertainty: 0, HTDPID: 10, Epoch: 1997.0},
	4910: {EPSG: 4910, Name: "ITRF88", Type: FrameHTDP, HTDPID: 11, Epoch: 1988.0},
	4911: {EPSG: 4911, Name: "ITRF89", Type: FrameHTDP, HTDPID: 12, Epoch: 1988.0},
	7901: {EPSG: 7901, Name: "ITRF89", Type: FrameHTDP, HTDPID: 12, Epoch: 1988.0},
	7902: {EPSG: 7902, Name: "ITRF90", Type: FrameHTDP, HTDPID: 13, Epoch: 1988.0},
	7903: {EPSG: 7903, Name: "ITRF91", Type: FrameHTDP, HTDPID: 14, Epoch: 1988.0},
	7904: {EPSG: 7904, Name: "ITRF92", Type: FrameHTDP, HTDPID: 15, Epoch: 1988.0},
	7905: {EPSG: 7905, Name: "ITRF93", Type: FrameHTDP, HTDPID: 16, Epoch: 1988.0},
	7906: {EPSG: 7906, Name: "ITRF94", Type: FrameHTDP, HTDPID: 17, Epoch: 1988.0},
	7907: {EPSG: 7907, Name: "ITRF96", Type: FrameHTDP, HTDPID: 18, Epoch: 1996.0},
	7908: {EPSG: 7908, Name: "ITRF97", Type: FrameHTDP, HTDPID: 19, Epoch: 1997.0},
	7909: {EPSG: 7909, Name: "ITRF2000", Type: FrameHTDP, HTDPID: 20, Epoch: 2000.0},
	7910: {EPSG: 7910, Name: "ITRF2005", Type: FrameHTDP, HTDPID: 21, Epoch: 2000.0},
	7911: {EPSG: 7911, Name: "ITRF2008", Type: FrameHTDP, HTDPID: 22, Epoch: 2000.0},
	7912: {EPSG: 7912, Name: "ELLIPSOID", Description: "IGS14/IGb14/WGS84/ITRF2014 Ellipsoid", Type: FrameHTDP, HTDPID: 23, Epoch: 2000.0},
	1322: {EPSG: 1322, Name: "ITRF2020", Description: "IGS20", Type: FrameHTDP, HTDPID: 24, Epoch: 2000.0},
}

var CDNFrames = map[int]Frame{
	9245: {EPSG: 9245, Name: "CGVD2013(CGG2013a) height", Type: FrameCDN, Uncertainty: 0},
	6647: {EPSG: 6647, Name: "CGVD2013(CGG2013) height", Type: FrameCDN, Uncertainty: 0},
	3855: {EPSG: 3855, Name: "EGM2008 height", Type: FrameCDN, Uncertainty: 0},
	5773: {EPSG: 5773, Name: "EGM96 height", Type: FrameCDN, Uncertainty: 0},
	5703: {EPSG: 5703, Name: "NAVD88 height", Type: FrameCDN, Uncertainty: 0.05},
	6360: {EPSG: 6360, Name: "NAVD88 height (usFt)", Type: FrameCDN, Uncertainty: 0.05},
	8228: {EPSG: 8228, Name: "NAVD88 height (Ft)", Type: FrameCDN, Uncertainty: 0.05},
	6644: {EPSG: 6644, Name: "GUVD04 height", Type: FrameCDN, Uncertainty: 0},
	6641: {EPSG: 6641, Name: "PRVD02 height", Type: FrameCDN, Uncertainty: 0},
	6643: {EPSG: 6643, Name: "ASVD02 height", Type: FrameCDN, Uncertainty: 0},
	9279: {EPSG: 9279, Name: "SA LLD height", Type: FrameCDN, Uncertainty: 0},
}

var GeoidModels = map[string]struct {
	Name        string
	Uncertainty float64
}{
	"g2018":  {Name: "geoid 2018", Uncertainty: 0.0127},
	"g2012b": {Name: "geoid 2012b", Uncertainty: 0.017},
	"g2012a": {Name: "geoid 2012a", Uncertainty: 0.017},
	"g1999":  {Name: "geoid 1999", Uncertainty: 0.046},
	"geoid09": {Name: "geoid 2009", Uncertainty: 0.05},
	"geoid03": {Name: "geoid 2003", Uncertainty: 0.046},
}

func GetFrameByEPSG(epsg int) *Frame {
	if f, ok := TidalFrames[epsg]; ok {
		return &f
	}
	if f, ok := HTDPFrames[epsg]; ok {
		return &f
	}
	if f, ok := CDNFrames[epsg]; ok {
		return &f
	}
	return nil
}

func GetFrameByName(name string) *Frame {
	name = strings.ToLower(name)
	for _, f := range TidalFrames {
		if strings.Contains(strings.ToLower(f.Name), name) {
			return &f
		}
	}
	for _, f := range HTDPFrames {
		if strings.Contains(strings.ToLower(f.Name), name) {
			return &f
		}
		if strings.Contains(strings.ToLower(f.Description), name) {
			return &f
		}
	}
	for _, f := range CDNFrames {
		if strings.Contains(strings.ToLower(f.Name), name) {
			return &f
		}
	}
	return nil
}

func FrameTypeName(epsg int) string {
	f := GetFrameByEPSG(epsg)
	if f == nil {
		return "unknown"
	}
	switch f.Type {
	case FrameTidal:
		return "tidal"
	case FrameHTDP:
		return "htdp"
	case FrameCDN:
		return "cdn"
	case FrameGeoid:
		return "geoid"
	default:
		return "unknown"
	}
}

func ListFrames() string {
	var s strings.Builder
	s.WriteString("Tidal:\n")
	for _, f := range TidalFrames {
		s.WriteString(fmt.Sprintf("  %d\t%s\n", f.EPSG, f.Name))
	}
	s.WriteString("HTDP:\n")
	for _, f := range HTDPFrames {
		s.WriteString(fmt.Sprintf("  %d\t%s\n", f.EPSG, f.Name))
	}
	s.WriteString("CDN:\n")
	for _, f := range CDNFrames {
		s.WriteString(fmt.Sprintf("  %d\t%s\n", f.EPSG, f.Name))
	}
	return s.String()
}
