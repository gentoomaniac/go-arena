package entities

import (
	"image"
)

type Tank struct {
	Name           string
	Position       image.Point
	hitbox         image.Rectangle
	Health         int
	MaxHealth      int
	Energy         int
	MaxEnergy      int
	Speed          int
	Orientation    float64
	ColisionBounds image.Rectangle
}

func (t Tank) CollisionBox() image.Rectangle {
	return image.Rect(
		t.Position.X+t.ColisionBounds.Min.X,
		t.Position.Y+t.ColisionBounds.Min.Y,
		t.Position.X+t.ColisionBounds.Max.X,
		t.Position.Y+t.ColisionBounds.Max.Y,
	)
}

func (t Tank) Hitbox() image.Rectangle {
	return t.hitbox
}
func (t *Tank) SetHitbox(hitbox image.Rectangle) {
	t.hitbox = hitbox
}
