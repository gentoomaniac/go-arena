package physics

import (
	"math"

	"github.com/gentoomaniac/go-arena/vector"
)

func PointInRectangle(a vector.Vec2, b vector.Rectangle) bool {
	if a.X >= b.Min.X &&
		a.X <= b.Max.X &&
		a.Y >= b.Min.Y &&
		a.Y <= b.Max.Y {
		return true
	}
	return false
}
func DoBoxesCollide(a vector.Rectangle, b vector.Rectangle) bool {
	return PointInRectangle(a.Min, b) ||
		PointInRectangle(vector.Vec2{X: a.Min.X, Y: a.Max.Y}, b) ||
		PointInRectangle(a.Max, b) ||
		PointInRectangle(vector.Vec2{X: a.Max.X, Y: a.Min.Y}, b)
}

// calculates the intersection point of the line segments defined by (v1, v2) and (v3, v4)
func Intersection(v1, v2, v3, v4 vector.Vec2) *vector.Vec2 {
	// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line_segment
	t := ((v1.X-v3.X)*(v3.Y-v4.Y) - (v1.Y-v3.Y)*(v3.X-v4.X)) / ((v1.X-v2.X)*(v3.Y-v4.Y) - (v1.Y-v2.Y)*(v3.X-v4.X))
	u := ((v2.X-v1.X)*(v1.Y-v3.Y) - (v2.Y-v1.Y)*(v1.X-v3.X)) / ((v1.X-v2.X)*(v3.Y-v4.Y) - (v1.Y-v2.Y)*(v3.X-v4.X))

	if 0.0 <= t && t <= 1.0 && 0.0 <= u && u <= 1.0 {
		return &vector.Vec2{
			X: v1.X + t*(v2.X-v1.X),
			Y: v1.Y + t*(v2.Y-v1.Y),
		}
	}

	return nil
}

func Distance(a, b vector.Vec2) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

func DistanceBetweenCircles(a vector.Circle, b vector.Circle) float64 {
	return math.Sqrt(math.Pow(a.Position.X-b.Position.X, 2)+math.Pow(a.Position.Y-b.Position.Y, 2)) - (a.Radius + b.Radius)
}

func DoMovingCirclesCollide(aCircle vector.Circle, aVelocity vector.Vec2, bCircle vector.Circle, bVelocity vector.Vec2) float64 {

	return 0
}

func PointLineDistance(v1, v2, p vector.Vec2) float64 {
	slope := v2.Y - v1.Y/v2.X - v1.X
	b := v1.Y - slope*v1.X
	n := v1.Perpendicular().Unit()

	return math.Abs(n.X*p.X+n.Y*p.Y+b) / math.Sqrt(n.X*n.X+n.Y*n.Y)
}
