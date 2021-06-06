package main

import "github.com/gentoomaniac/go-arena/entities"

func checkColisionPoint(a entities.Vector, b entities.CollisionBox) bool {
	if a.X >= b.Min.X &&
		a.X <= b.Max.X &&
		a.Y >= b.Min.Y &&
		a.Y <= b.Max.Y {
		return true
	}
	return false
}
func checkColisionBox(a entities.CollisionBox, b entities.CollisionBox) bool {
	return checkColisionPoint(a.Min, b) ||
		checkColisionPoint(entities.Vector{X: a.Min.X, Y: a.Max.Y}, b) ||
		checkColisionPoint(a.Max, b) ||
		checkColisionPoint(entities.Vector{X: a.Max.X, Y: a.Min.Y}, b)
}
