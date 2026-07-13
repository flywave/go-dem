package waffle

import (
	"fmt"

	"github.com/flywave/go-dem"
	"github.com/flywave/go3d/float64/vec2"
)

type Point struct {
	Position    vec2.T
	Z           float64
	Uncertainty float64
}

type Options struct {
	Region       *dem.Region
	NoData       float64
	SearchRadius float64
	MinPoints    int
	UpperZ       *float64
	LowerZ       *float64
	ChunkX       int
	ChunkY       int
	Threads      int
}

type Waffle interface {
	Name() string
	Run(points []Point, opts *Options) (*Result, error)
}

type Result struct {
	DEM         []float64
	Stack       [][]float64
	Mask        []float64
	Uncertainty []float64
	Region      *dem.Region
}

type Constructor func() Waffle

var registry = make(map[dem.InterpMethod]Constructor)

func Register(method dem.InterpMethod, c Constructor) {
	registry[method] = c
}

func New(method dem.InterpMethod) (Waffle, error) {
	c, ok := registry[method]
	if !ok {
		return nil, fmt.Errorf("unknown interpolation method: %s", method)
	}
	return c(), nil
}

func ListMethods() []dem.InterpMethod {
	methods := make([]dem.InterpMethod, 0, len(registry))
	for m := range registry {
		methods = append(methods, m)
	}
	return methods
}

type baseWaffle struct {
	name string
}

func (b *baseWaffle) Name() string { return b.name }
