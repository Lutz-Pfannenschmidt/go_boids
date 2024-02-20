package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Color struct {
	sdl.Color

	R, G, B, A uint8
}

// Uint32 returns the color as a 32-bit value.
// fixes the issue with the sdl.Color.Uint32() method, where the alpha channel is at the wrong position.
func (c Color) Uint32() uint32 {
	var v uint32
	v |= uint32(c.A) << 24
	v |= uint32(c.R) << 16
	v |= uint32(c.G) << 8
	v |= uint32(c.B)
	return v
}

func (c Color) ToSdl() sdl.Color {
	return sdl.Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

func RainbowGradient(x, y float64) Color {
	nx := (x + 1) / 2
	ny := (y + 1) / 2

	hue := nx * 360

	saturation := 1 - ny
	value := ny

	r, g, b := hsvToRgb(hue, saturation, value)

	return Color{R: r, G: g, B: b, A: 255}
}

func hsvToRgb(h, s, v float64) (uint8, uint8, uint8) {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	switch {
	case 0 <= h && h < 60:
		return uint8((c + m) * 255), uint8((x + m) * 255), uint8(m * 255)
	case 60 <= h && h < 120:
		return uint8((x + m) * 255), uint8((c + m) * 255), uint8(m * 255)
	case 120 <= h && h < 180:
		return uint8(m * 255), uint8((c + m) * 255), uint8((x + m) * 255)
	case 180 <= h && h < 240:
		return uint8((m * 255)), uint8((x + m) * 255), uint8((c + m) * 255)
	case 240 <= h && h < 300:
		return uint8((x + m) * 255), uint8(m * 255), uint8((c + m) * 255)
	default:
		return uint8(c * 255), uint8(x * 255), uint8(m * 255)
	}
}
