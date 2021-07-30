package vector

import (
	"fmt"
	"math"
)

type Vec2 struct {
	X float64
	Y float64
}

func (v Vec2) String() string {
	return fmt.Sprintf("(%f, %f)", v.X, v.Y)
}

func (v Vec2) Length() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

func (v Vec2) Unit() Vec2 {
	magnitude := math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))

	return Vec2{v.X / magnitude, v.Y / magnitude}
}

func (v Vec2) ScalarProduct(m float64) Vec2 {
	return Vec2{v.X * m, v.Y * m}
}

func (v Vec2) Sum(v2 Vec2) Vec2 {
	return Vec2{v.X + v2.X, v.Y + v2.X}
}

func (v Vec2) DotProduct(v2 Vec2) float64 {
	return v.X*v2.X + v.Y*v2.Y
}

func (v Vec2) WithLength(l float64) Vec2 {
	if v.Length() == 0 {
		return Vec2{}
	}
	return v.ScalarProduct(l / v.Length())
}

func (v Vec2) Negative() Vec2 {
	return v.ScalarProduct(-1)
}

func (v Vec2) Rotate(b float64) Vec2 {
	rad := b * (math.Pi / 180)
	return Vec2{
		X: v.X*math.Cos(rad) - v.Y*math.Sin(rad),
		Y: v.X*math.Sin(rad) + v.Y*math.Cos(rad),
	}
}
