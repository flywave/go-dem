module github.com/flywave/go-dem

go 1.24

toolchain go1.24.4

require (
	github.com/flywave/flywave-gdal v0.0.0
	github.com/flywave/flywave-pointcloud v0.0.0
	github.com/flywave/go-delaunay v0.0.0
	github.com/flywave/go-geo v0.0.0-20250607132733-46bd30e585ce
	github.com/flywave/go-geoid v0.0.0-20221115021843-2080cff61475
	github.com/flywave/go-kriging v0.0.0
	github.com/flywave/go3d v0.0.0-20250619003741-cab1a6ea6de6
)

require (
	github.com/flywave/go-cog v0.0.0-20250607133043-41acd04eb904 // indirect
	github.com/flywave/go-geom v0.0.0-20250607125323-f685bf20f12c // indirect
	github.com/flywave/go-geos v0.0.0-20250607125930-047054a9f657 // indirect
	github.com/flywave/go-proj v0.0.0-20250607132305-d70d32f5ad2d // indirect
	github.com/google/tiff v0.0.0-20161109161721-4b31f3041d9a // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/hhrutter/lzw v1.0.0 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/image v0.28.0 // indirect
	gonum.org/v1/gonum v0.8.2 // indirect
)

replace github.com/flywave/go-geo => ../go-geo

replace github.com/flywave/go-geos => ../go-geos

replace github.com/flywave/go-proj => ../go-proj

replace github.com/flywave/go-geoid => ../go-geoid

replace github.com/flywave/go-kriging => ../go-kriging

replace github.com/flywave/go-delaunay => ../go-delaunay

replace github.com/flywave/go-cog => ../go-cog

replace github.com/flywave/flywave-gdal => ../flywave-gdal

replace github.com/flywave/flywave-pointcloud => ../flywave-pointcloud

replace github.com/flywave/go-geom => ../go-geom

replace github.com/flywave/go3d => ../go3d
