package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Poly struct {
	Vertices []Point
}

func (p Poly) Draw(r *sdl.Renderer, c *Color) {
	for i := 0; i < len(p.Vertices); i++ {
		v1 := p.Vertices[i]
		v2 := p.Vertices[(i+1)%len(p.Vertices)]

		r.SetDrawColor(c.R, c.G, c.B, c.A)
		r.DrawLine(int32(v1.X), int32(v1.Y), int32(v2.X), int32(v2.Y))
	}
}

type Rect struct {
	Position Point
	Width    float64
	Height   float64
}

func (r Rect) Contains(p Point) bool {
	return p.X >= r.Position.X && p.X <= r.Position.X+r.Width && p.Y >= r.Position.Y && p.Y <= r.Position.Y+r.Height
}

func (r Rect) Intersects(other Rect) bool {
	return r.Position.X < other.Position.X+other.Width &&
		r.Position.X+r.Width > other.Position.X &&
		r.Position.Y < other.Position.Y+other.Height &&
		r.Position.Y+r.Height > other.Position.Y
}

func (r Rect) Draw(renderer *sdl.Renderer, color *Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.DrawRect(&sdl.Rect{X: int32(r.Position.X), Y: int32(r.Position.Y), W: int32(r.Width), H: int32(r.Height)})
}

func (r Rect) Fill(renderer *sdl.Renderer, color *Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.FillRect(&sdl.Rect{X: int32(r.Position.X), Y: int32(r.Position.Y), W: int32(r.Width), H: int32(r.Height)})
}

type Circle struct {
	Position Point
	Radius   float64
}

func (c Circle) Draw(r *sdl.Renderer, color Color) {
	r.SetDrawColor(color.R, color.G, color.B, color.A)
	for i := 0; i < 360; i++ {
		x := c.Position.X + c.Radius*math.Cos(float64(i))
		y := c.Position.Y + c.Radius*math.Sin(float64(i))
		r.DrawPoint(int32(x), int32(y))
	}
}

func (c Circle) Fill(r *sdl.Renderer, color Color) {
	r.SetDrawColor(color.R, color.G, color.B, color.A)
	r.DrawPoint(int32(c.Position.X), int32(c.Position.Y))
	for i := 0; i < int(c.Radius); i++ {
		r.DrawPoint(int32(c.Position.X+float64(i)), int32(c.Position.Y))
		r.DrawPoint(int32(c.Position.X-float64(i)), int32(c.Position.Y))
		r.DrawPoint(int32(c.Position.X), int32(c.Position.Y+float64(i)))
		r.DrawPoint(int32(c.Position.X), int32(c.Position.Y-float64(i)))
	}
}
