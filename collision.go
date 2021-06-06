package main

import (
	"github.com/gentoomaniac/go-arena/vector"
)

func checkColisionPoint(a vector.Vec2, b vector.Rectangle) bool {
	if a.X >= b.Min.X &&
		a.X <= b.Max.X &&
		a.Y >= b.Min.Y &&
		a.Y <= b.Max.Y {
		return true
	}
	return false
}
func checkColisionBox(a vector.Rectangle, b vector.Rectangle) bool {
	return checkColisionPoint(a.Min, b) ||
		checkColisionPoint(vector.Vec2{X: a.Min.X, Y: a.Max.Y}, b) ||
		checkColisionPoint(a.Max, b) ||
		checkColisionPoint(vector.Vec2{X: a.Max.X, Y: a.Min.Y}, b)
}
