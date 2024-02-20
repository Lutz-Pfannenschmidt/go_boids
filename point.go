package main

import "math"

type Point struct {
	X, Y float64
}

func (p Point) SqDistance(p2 Point) float64 {
	return ((p.X-p2.X)*(p.X-p2.X) + (p.Y-p2.Y)*(p.Y-p2.Y))
}

type Vec2 struct {
	X, Y float64
}

func (v *Vec2) Add(other Vec2) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vec2) Sub(other Vec2) {
	v.X -= other.X
	v.Y -= other.Y
}

func (v *Vec2) Mult(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Vec2) Div(scalar float64) {
	v.X /= scalar
	v.Y /= scalar
}

func (v *Vec2) Mag() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
