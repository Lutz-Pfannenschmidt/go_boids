package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var renderer *sdl.Renderer
var tree QuadTree
var boids []*Boid

var addBoids int

var mouse Point

const (
	BOID_COUNT    = 100
	TREE_CAP      = 1
	INIT_SCREEN_W = 900
	INIT_SCREEN_H = 900

	DRAW_QUADTREE = false
	DRAW_BOIDS    = true
	DRAW_CURSOR   = false

	BOID_MAX_SPEED = 2.5
	BOID_MIN_SPEED = 2.2

	BOID_RANGE           = 150
	BOID_COLLISION_RANGE = 55

	BOID_TRI_W = 10.0
	BOID_TRI_H = 30.0
)

var (
	BOID_ALIGN_FACTOR      = 0.073
	BOID_COHESION_FACTOR   = 0.015
	BOID_SEPARATION_FACTOR = 0.05
)

var SCREEN_W = INIT_SCREEN_W
var SCREEN_H = INIT_SCREEN_H

func main() {
	mouse = Point{X: 0, Y: 0}

	var errs []error
	errs = append(errs, sdl.Init(sdl.INIT_EVENTS))
	errs = append(errs, sdl.Init(sdl.INIT_VIDEO))

	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}

	defer sdl.Quit()

	window, err := sdl.CreateWindow("Boids", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, INIT_SCREEN_W, INIT_SCREEN_H, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.SetResizable(true)

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	tree = *NewQuadTree(Rect{Position: Point{X: 0, Y: 0}, Width: float64(SCREEN_W), Height: float64(SCREEN_H)}, TREE_CAP)

	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < BOID_COUNT; i++ {
		// Random position
		pos := Point{X: float64(gen.Intn(SCREEN_W)), Y: float64(gen.Intn(SCREEN_H))}
		// Random velocity
		vel := gen.Float64()*(BOID_MAX_SPEED-BOID_MIN_SPEED) + BOID_MIN_SPEED
		// Random direction (in radians)
		dir := gen.Float64() * 2 * math.Pi

		boid := Boid{Position: pos, Velocity: vel, Direction: dir}

		tree.Insert(&boid)
		boids = append(boids, &boid)
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			case *sdl.MouseButtonEvent:
				if ev.Type == sdl.MOUSEBUTTONDOWN {
					if ev.Button == sdl.BUTTON_LEFT {
						pos := Point{X: float64(ev.X), Y: float64(ev.Y)}
						vel := rand.Float64()*(BOID_MAX_SPEED-BOID_MIN_SPEED) + BOID_MIN_SPEED
						dir := rand.Float64() * 2 * math.Pi

						boid := Boid{Position: pos, Velocity: vel, Direction: dir}

						tree.Insert(&boid)
						boids = append(boids, &boid)
						addBoids++
					}
				}
			case *sdl.MouseMotionEvent:
				mouse = Point{X: float64(ev.X), Y: float64(ev.Y)}
			case *sdl.WindowEvent:
				windowEvent := event.(*sdl.WindowEvent)
				if windowEvent.Event == sdl.WINDOWEVENT_RESIZED {
					SCREEN_W = int(windowEvent.Data1)
					SCREEN_H = int(windowEvent.Data2)
					renderer.Destroy()
					renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
					if err != nil {
						panic(err)
					}
				}
			case *sdl.KeyboardEvent:
				if ev.Type == sdl.KEYDOWN {
					switch ev.Keysym.Sym {
					case sdl.K_q:
						BOID_ALIGN_FACTOR += 0.001
						fmt.Println("New align factor: ", BOID_ALIGN_FACTOR)
					case sdl.K_a:
						BOID_ALIGN_FACTOR -= 0.001
						fmt.Println("New align factor: ", BOID_ALIGN_FACTOR)
					case sdl.K_w:
						BOID_SEPARATION_FACTOR += 0.001
						fmt.Println("New separation factor: ", BOID_SEPARATION_FACTOR)
					case sdl.K_s:
						BOID_SEPARATION_FACTOR -= 0.001
						fmt.Println("New separation factor: ", BOID_SEPARATION_FACTOR)
					case sdl.K_e:
						BOID_COHESION_FACTOR += 0.001
						fmt.Println("New cohesion factor: ", BOID_COHESION_FACTOR)
					case sdl.K_d:
						BOID_COHESION_FACTOR -= 0.001
						fmt.Println("New cohesion factor: ", BOID_COHESION_FACTOR)
					case sdl.K_ESCAPE:
						running = false
						fmt.Println("\nvar (" +
							"\n\tBOID_ALIGN_FACTOR=" + strconv.FormatFloat(BOID_ALIGN_FACTOR, 'f', -1, 64) +
							"\n\tBOID_COHESION_FACTOR=" + strconv.FormatFloat(BOID_COHESION_FACTOR, 'f', -1, 64) +
							"\n\tBOID_SEPARATION_FACTOR=" + strconv.FormatFloat(BOID_SEPARATION_FACTOR, 'f', -1, 64) +
							"\n)")
					case sdl.K_F11:
						if window.GetFlags()&sdl.WINDOW_FULLSCREEN == 0 {
							window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
						} else {
							window.SetFullscreen(0)
						}
					}
				}
			}
		}

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()

		update()
		draw()

		renderer.Present()
		time.Sleep(16 * time.Millisecond)
	}
}

func update() {
	for _, b := range boids {
		b.Update()
		b.Flock(&tree)
	}

	// wrap around
	for _, b := range boids {
		if b.Position.X < 0 {
			b.Position.X = float64(SCREEN_W)
		}
		if b.Position.X > float64(SCREEN_W) {
			b.Position.X = 0
		}
		if b.Position.Y < 0 {
			b.Position.Y = float64(SCREEN_H)
		}
		if b.Position.Y > float64(SCREEN_H) {
			b.Position.Y = 0
		}
	}

	c := 0
	tree = *NewQuadTree(Rect{Position: Point{X: 0, Y: 0}, Width: float64(SCREEN_W), Height: float64(SCREEN_H)}, TREE_CAP)
	for _, b := range boids {
		tree.Insert(b)
		c++
	}

	if c != BOID_COUNT+addBoids {
		panic("Boids count mismatch!")
	}
}

func draw() {

	if DRAW_CURSOR {
		mouseCircle := Circle{Position: mouse, Radius: 20}
		mouseCircle.Fill(renderer, Color{R: 0, G: 0, B: 255, A: 50})
		mouseCircle.Draw(renderer, Color{R: 0, G: 0, B: 255, A: 255})
	}

	if DRAW_QUADTREE {
		tree.DrawOutline(renderer, &Color{R: 0, G: 0, B: 0, A: 255})
		// tree.Fill(renderer, &Color{R: 0, G: 0, B: 0, A: 255})
	}

	if DRAW_BOIDS {
		tree.DrawBoids(renderer, &Color{R: 0, G: 0, B: 0, A: 255})
	}

}
