package entities

import (
	"github.com/gentoomaniac/go-arena/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

type Object interface {
	Name() string
	CollisionBox() vector.Rectangle
	Position() vector.Vec2
	SetPosition(vector.Vec2)
	Speed() int
	SetSpeed(int)
	Orientation() float64
	SetOrientation(float64)
	Sprite() *ebiten.Image
}
