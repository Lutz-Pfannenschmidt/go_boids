package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type QuadTree struct {
	Bounds   Rect
	Boids    []*Boid
	capacity int
	divided  bool

	northWest *QuadTree
	northEast *QuadTree
	southWest *QuadTree
	southEast *QuadTree
}

func NewQuadTree(bounds Rect, capacity int) *QuadTree {
	return &QuadTree{
		Bounds:   bounds,
		Boids:    nil,
		capacity: capacity,
		divided:  false,
	}
}

func (q *QuadTree) Insert(b *Boid) bool {
	if !q.Bounds.Contains(b.Position) {
		return false
	}

	if len(q.Boids) < q.capacity {
		q.Boids = append(q.Boids, b)
		return true
	}

	if !q.divided {
		q.Subdivide()
	}

	if q.northWest.Insert(b) {
		return true
	}
	if q.northEast.Insert(b) {
		return true
	}
	if q.southWest.Insert(b) {
		return true
	}
	return q.southEast.Insert(b)
}

func (q *QuadTree) Subdivide() {

	q.divided = true

	x := q.Bounds.Position.X
	y := q.Bounds.Position.Y
	w := q.Bounds.Width
	h := q.Bounds.Height

	nw := Rect{Position: Point{X: x, Y: y}, Width: w / 2, Height: h / 2}
	ne := Rect{Position: Point{X: x + w/2, Y: y}, Width: w / 2, Height: h / 2}
	sw := Rect{Position: Point{X: x, Y: y + h/2}, Width: w / 2, Height: h / 2}
	se := Rect{Position: Point{X: x + w/2, Y: y + h/2}, Width: w / 2, Height: h / 2}

	q.northWest = NewQuadTree(nw, q.capacity)
	q.northEast = NewQuadTree(ne, q.capacity)
	q.southWest = NewQuadTree(sw, q.capacity)
	q.southEast = NewQuadTree(se, q.capacity)

	for _, b := range q.Boids {
		if q.northWest.Bounds.Contains(b.Position) {
			q.northWest.Insert(b)
		} else if q.northEast.Bounds.Contains(b.Position) {
			q.northEast.Insert(b)
		} else if q.southWest.Bounds.Contains(b.Position) {
			q.southWest.Insert(b)
		} else if q.southEast.Bounds.Contains(b.Position) {
			q.southEast.Insert(b)
		}
	}

	q.Boids = nil
}

func (q *QuadTree) QueryRanges(circles []*Circle) []*Boid {
	var boids []*Boid

	for _, c := range circles {
		boids = append(boids, q.QueryRange(c)...)
	}

	return boids
}

// QueryRange returns all boids in the given range.
// accepts an optional rectangle to use as the query range (quad tree stuff) to avoid creating a new rectangle every time.
func (q *QuadTree) QueryRange(c *Circle, optRect ...*Rect) []*Boid {
	var rangeRect *Rect

	if len(optRect) > 0 {
		rangeRect = optRect[0]
	} else {
		rangeRect = &Rect{
			Position: Point{X: c.Position.X - c.Radius, Y: c.Position.Y - c.Radius},
			Width:    c.Radius * 2,
			Height:   c.Radius * 2,
		}
	}

	var boids []*Boid

	if !q.Bounds.Intersects(*rangeRect) {
		return boids
	}

	for _, b := range q.Boids {
		if rangeRect.Contains(b.Position) {
			boids = append(boids, b)
		}
	}

	if q.divided {
		boids = append(boids, q.northWest.QueryRange(c, rangeRect)...)
		boids = append(boids, q.northEast.QueryRange(c, rangeRect)...)
		boids = append(boids, q.southWest.QueryRange(c, rangeRect)...)
		boids = append(boids, q.southEast.QueryRange(c, rangeRect)...)
	}

	filteredBoids := []*Boid{}

	for _, b := range boids {
		sqDistance := c.Position.SqDistance(b.Position)
		if sqDistance < math.Pow(c.Radius, 2) {
			filteredBoids = append(filteredBoids, b)
		}
	}

	return filteredBoids
}

func (q *QuadTree) DrawOutline(r *sdl.Renderer, c *Color) {
	q.Bounds.Draw(r, c)

	if q.divided {
		q.northWest.DrawOutline(r, c)
		q.northEast.DrawOutline(r, c)
		q.southWest.DrawOutline(r, c)
		q.southEast.DrawOutline(r, c)
	}
}

func (q *QuadTree) DrawBoids(r *sdl.Renderer, c *Color) {
	for _, b := range q.Boids {
		b.Draw(r, c)
	}

	if q.divided {
		q.northWest.DrawBoids(r, c)
		q.northEast.DrawBoids(r, c)
		q.southWest.DrawBoids(r, c)
		q.southEast.DrawBoids(r, c)
	}
}

func (q *QuadTree) Fill(r *sdl.Renderer, c *Color) {
	q.Bounds.Fill(r, c)

	if q.divided {
		q.northWest.Fill(r, c)
		q.northEast.Fill(r, c)
		q.southWest.Fill(r, c)
		q.southEast.Fill(r, c)
	}
}
