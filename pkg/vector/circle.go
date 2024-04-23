package vector

import "fmt"

type Circle struct {
	Position Vec2
	Radius   float64
}

func (c Circle) String() string {
	return fmt.Sprintf("(%s, %f)", c.Position, c.Radius)
}
