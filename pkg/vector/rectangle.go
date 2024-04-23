package vector

import "fmt"

type Rectangle struct {
	Min Vec2
	Max Vec2
}

func (c Rectangle) String() string {
	return fmt.Sprintf("[%s, %s]", c.Min, c.Max)
}

func Rect(x0, y0, x1, y1 float64) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle{Vec2{x0, y0}, Vec2{x1, y1}}
}
