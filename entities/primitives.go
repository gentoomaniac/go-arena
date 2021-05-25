package entities

import "fmt"

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

type Color struct {
	R     float64
	G     float64
	B     float64
	Alpha float64
}

func (c Color) String() string {
	return fmt.Sprintf("(%f, %f, %f, %f)", c.R, c.G, c.B, c.Alpha)
}
