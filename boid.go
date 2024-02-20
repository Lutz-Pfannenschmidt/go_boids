package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Boid struct {
	Position Point

	// direction is the angle in radians
	Direction float64
	TurnSpeed float64
	// velocity is the speed and direction of the boid
	Velocity float64
}

func (b *Boid) Update() {
	b.Position.X += b.Velocity * float64(math.Cos(b.Direction))
	b.Position.Y += b.Velocity * float64(math.Sin(b.Direction))
}

func (b *Boid) Draw(renderer *sdl.Renderer, color *Color) {
	// Calculate the coordinates of the vertices
	poly := Poly{
		Vertices: []Point{
			{X: b.Position.X + BOID_TRI_H*math.Cos(b.Direction), Y: b.Position.Y + BOID_TRI_H*math.Sin(b.Direction)},
			{X: b.Position.X + BOID_TRI_W*math.Cos(b.Direction-math.Pi/2), Y: b.Position.Y + BOID_TRI_W*math.Sin(b.Direction-math.Pi/2)},
			{X: b.Position.X + BOID_TRI_W*math.Cos(b.Direction+math.Pi/2), Y: b.Position.Y + BOID_TRI_W*math.Sin(b.Direction+math.Pi/2)},
		},
	}

	// Draw the boid
	poly.Draw(renderer, color)

	// Draw the velocity vector multiplied by 10
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.DrawLine(int32(b.Position.X), int32(b.Position.Y), int32(b.Position.X+b.Velocity*10*math.Cos(b.Direction)), int32(b.Position.Y+b.Velocity*10*math.Sin(b.Direction)))
}

func (b *Boid) Flock(tree *QuadTree) {

	var fovs []*Circle

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			fovs = append(fovs, &Circle{
				Position: Point{X: b.Position.X + float64(i*SCREEN_W), Y: b.Position.Y + float64(j*SCREEN_H)},
				Radius:   BOID_RANGE,
			})
		}
	}

	neighborsAndMe := tree.QueryRanges(fovs)
	neighbors := []*Boid{}

	for _, n := range neighborsAndMe {
		if n != b {
			neighbors = append(neighbors, n)
		}
	}

	b.TurnSpeed *= 0.025

	b.TurnSpeed += b.Seperation(neighbors) * BOID_SEPARATION_FACTOR
	b.TurnSpeed += b.Alignment(neighbors) * BOID_ALIGN_FACTOR
	b.TurnSpeed += b.Cohesion(neighbors) * BOID_COHESION_FACTOR

	b.Direction += b.TurnSpeed
}

// Seperation returns the recommendet new direction to avoid other boids in radians
func (b *Boid) Seperation(neighbors []*Boid) float64 {
	runAwayVec := Vec2{X: 0, Y: 0}

	for _, n := range neighbors {
		sqDistance := b.Position.SqDistance(n.Position)
		if sqDistance < math.Pow(BOID_COLLISION_RANGE, 2) {
			dVec := Vec2{
				X: modularDifference(b.Position.X, n.Position.X, float64(SCREEN_W)),
				Y: modularDifference(b.Position.Y, n.Position.Y, float64(SCREEN_H)),
			}
			//dVec.Mult((BOID_COLLISION_RANGE - dVec.Mag()) / dVec.Mag())
			dVec.Div(math.Pow(dVec.Mag(), 2))
			runAwayVec.Add(dVec)
		}
	}

	if runAwayVec.X == 0 && runAwayVec.Y == 0 {
		return 0
	}

	avgAngle := vec2ToAngle(runAwayVec)

	return getAngleDifference(avgAngle, b.Direction)
}

func (b *Boid) Alignment(neighbors []*Boid) float64 {
	dirs := []float64{}
	for _, boid := range neighbors {
		dirs = append(dirs, boid.Direction)
	}

	avgAngle := avgAngle(dirs)

	if math.IsNaN(avgAngle) {
		return 0
	}

	return getAngleDifference(avgAngle, b.Direction)
}

func (b *Boid) Cohesion(neighbors []*Boid) float64 {
	avgRelPos := Vec2{X: 0, Y: 0}

	for _, boid := range neighbors {
		avgRelPos.Add(Vec2{X: modularDifference(boid.Position.X, b.Position.X, float64(SCREEN_W)), Y: modularDifference(boid.Position.Y, b.Position.Y, float64(SCREEN_H))})
	}

	if avgRelPos.X == 0 && avgRelPos.Y == 0 {
		return 0
	}

	return getAngleDifference(vec2ToAngle(avgRelPos), b.Direction)
}

func angleToVec2(angle float64) Vec2 {
	y, x := math.Sincos(angle)
	return Vec2{X: x, Y: y}
}

func vec2ToAngle(vec Vec2) float64 {
	return math.Atan2(vec.Y, vec.X)
}

func avgAngle(angles []float64) float64 {
	sum := Vec2{X: 0, Y: 0}
	for _, angle := range angles {
		sum.Add(angleToVec2(angle))
	}
	if sum.X == 0 && sum.Y == 0 {
		return math.NaN()
	}
	return vec2ToAngle(sum)
}

func getAngleDifference(a, b float64) float64 {
	diff := a - b
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}
	for diff < -math.Pi {
		diff += 2 * math.Pi
	}
	return diff
}

func modularDifference(a, b, m float64) float64 {
	diff := a - b
	for diff > m {
		diff -= 2 * m
	}
	for diff < -m {
		diff += 2 * m
	}
	return diff
}
