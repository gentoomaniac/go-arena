package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Object interface {
	Name() string
	CollisionBox() CollisionBox
	Position() Vector
	SetPosition(Vector)
	Speed() int
	SetSpeed(int)
	Orientation() float64
	SetOrientation(float64)
	Sprite() *ebiten.Image
}
