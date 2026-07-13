package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/datalist"
	"github.com/flywave/go-dem/waffle"
	"github.com/flywave/go-geo"
)

func main() {
	regionStr := flag.String("R", "", "Region: xmin/xmax/ymin/ymax")
	resolution := flag.Float64("E", 0, "Resolution in degrees/meters")
	method := flag.String("M", "idw", "Interpolation method: idw, kriging, linear, nearest, natural_neighbor, cudem, cube, inpaint")
	output := flag.String("O", "output.tif", "Output DEM file")
	srsCode := flag.String("srs", "EPSG:4326", "Output SRS")
	noData := flag.Float64("nodata", -9999, "NoData value")
	searchRadius := flag.Float64("radius", 0, "Search radius for interpolation (0=auto)")
	minPoints := flag.Int("minpts", 0, "Minimum points for interpolation (0=auto)")
	idwPower := flag.Float64("power", 2.0, "IDW power parameter")
	datalistFlag := flag.Bool("datalist", false, "Treat input as datalist file")
	flag.Parse()

	if *regionStr == "" || *resolution <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: dem -R xmin/xmax/ymin/ymax -E resolution [-M method] [-O output] [sources...]\n")
		fmt.Fprintf(os.Stderr, "Methods: idw, kriging, linear, nearest, natural_neighbor, cudem, cube, inpaint\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	sources := flag.Args()
	if len(sources) == 0 {
		fmt.Fprintf(os.Stderr, "no source files provided\n")
		os.Exit(1)
	}

	srs := geo.NewProj(*srsCode)
	if srs == nil {
		fmt.Fprintf(os.Stderr, "invalid SRS: %s\n", *srsCode)
		os.Exit(1)
	}

	region, err := dem.NewRegionFromString(*regionStr, srs, *resolution, *resolution)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid region: %v\n", err)
		os.Exit(1)
	}

	var points []waffle.Point
	if *datalistFlag {
		dl, err := datalist.BuildDataList(sources)
		if err != nil {
			fmt.Fprintf(os.Stderr, "build datalist: %v\n", err)
			os.Exit(1)
		}
		var paths []string
		for _, e := range dl.Entries {
			paths = append(paths, e.Path)
		}
		points, err = waffle.PointsFromMultiple(paths)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load points: %v\n", err)
			os.Exit(1)
		}
	} else {
		points, err = waffle.PointsFromMultiple(sources)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load points: %v\n", err)
			os.Exit(1)
		}
	}

	interpMethod := dem.InterpMethod(*method)
	w, err := waffle.New(interpMethod)
	if err != nil {
		methods := waffle.ListMethods()
		strs := make([]string, len(methods))
		for i, m := range methods {
			strs[i] = string(m)
		}
		allMethods := strings.Join(strs, ", ")
		fmt.Fprintf(os.Stderr, "unknown method '%s'. Available: %s\n", *method, allMethods)
		os.Exit(1)
	}

	opts := &waffle.Options{
		Region:       region,
		NoData:       *noData,
		SearchRadius: *searchRadius,
		MinPoints:    *minPoints,
	}

	result, err := w.Run(points, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "interpolation failed: %v\n", err)
		os.Exit(1)
	}

	if err := dem.CreateDEM(result.DEM, region, *output, *noData); err != nil {
		fmt.Fprintf(os.Stderr, "write DEM failed: %v\n", err)
		os.Exit(1)
	}

	_ = idwPower

	fmt.Printf("DEM written to %s (%d points -> %dx%d grid, method=%s)\n",
		*output, len(points), region.XSize, region.YSize, *method)
}
