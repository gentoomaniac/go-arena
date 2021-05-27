package entities

import (
	"fmt"
)

type Vector struct {
	X float64
	Y float64
}

func (v Vector) String() string {
	return fmt.Sprintf("(%f, %f)", v.X, v.Y)
}

type CollisionBox struct {
	Min Vector
	Max Vector
}

func (c CollisionBox) String() string {
	return fmt.Sprintf("[%s, %s]", c.Min, c.Max)
}

func Box(x0, y0, x1, y1 float64) CollisionBox {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return CollisionBox{Vector{x0, y0}, Vector{x1, y1}}
}

type Color struct {
	R     float64
	G     float64
	B     float64
	Alpha float64
}

func (c Color) String() string {
	return fmt.Sprintf("(%f, %f, %f, %f)", c.R, c.G, c.B, c.Alpha)
}
