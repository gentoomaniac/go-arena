package player

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten"
)

type Vector struct {
	X float64
	Y float64
}

func (v Vector) String() string {
	return fmt.Sprintf("(%f, %f)", v.X, v.Y)
}

type CollisionBox struct {
	Min Vector
	Max Vector
}

func (c CollisionBox) String() string {
	return fmt.Sprintf("[%s, %s]", c.Min, c.Max)
}

type Player struct {
	Name           string
	Position       Vector
	Hitbox         image.Rectangle
	Health         int
	MaxHealth      int
	Energy         int
	MaxEnergy      int
	CurrentSpeed   int
	Speed          int
	MaxSpeed       int
	Orientation    float64
	ColisionBounds CollisionBox
	Collided       bool
	Sprite         *ebiten.Image
	AI             AI
}

func (p Player) GetCollisionBox() CollisionBox {
	return CollisionBox{
		Min: Vector{
			p.Position.X + float64(p.ColisionBounds.Min.X),
			p.Position.Y + float64(p.ColisionBounds.Min.Y)},
		Max: Vector{
			p.Position.X + float64(p.ColisionBounds.Max.X),
			p.Position.Y + float64(p.ColisionBounds.Max.Y),
		},
	}
}
