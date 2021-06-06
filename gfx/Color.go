package gfx

import "fmt"

type Color struct {
	R     float64
	G     float64
	B     float64
	Alpha float64
}

func (c Color) String() string {
	return fmt.Sprintf("(%f, %f, %f, %f)", c.R, c.G, c.B, c.Alpha)
}
