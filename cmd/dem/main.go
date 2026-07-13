package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/flywave/go-dem"
	"github.com/flywave/go-dem/datalist"
	"github.com/flywave/go-dem/waffle"
	"github.com/flywave/go-geo"
)

func main() {
	regionStr := flag.String("R", "", "Region: xmin/xmax/ymin/ymax")
	resolution := flag.Float64("E", 0, "Resolution in degrees/meters")
	method := flag.String("M", "idw", "Interpolation method: idw, kriging, linear, cubic, nearest")
	output := flag.String("O", "output.tif", "Output DEM file")
	srsCode := flag.String("srs", "EPSG:4326", "Output SRS")
	noData := flag.Float64("nodata", -9999, "NoData value")
	flag.Parse()

	if *regionStr == "" || *resolution <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: dem -R xmin/xmax/ymin/ymax -E resolution [-M method] [-O output] [sources...]\n")
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

	dataList, err := datalist.BuildDataList(sources)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build datalist: %v\n", err)
		os.Exit(1)
	}

	interpMethod := dem.InterpMethod(*method)
	w, err := waffle.New(interpMethod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create waffle: %v\n", err)
		os.Exit(1)
	}

	opts := &waffle.Options{
		Region: region,
		NoData: *noData,
	}

	var inputPaths []string
	for _, entry := range dataList.Entries {
		inputPaths = append(inputPaths, entry.Path)
	}

	result, err := w.Run(inputPaths, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "waffle run failed: %v\n", err)
		os.Exit(1)
	}

	if err := dem.CreateDEM(result.DEM, region, *output, *noData); err != nil {
		fmt.Fprintf(os.Stderr, "write DEM failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("DEM written to %s\n", *output)
}
