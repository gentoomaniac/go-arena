package entities

import (
	"image"

	"github.com/gentoomaniac/go-arena/gfx"
	"github.com/hajimehoshi/ebiten/v2"
)

type State int

const (
	Alive State = iota
	Dead
)

func (s State) String() string {
	return [...]string{"Alive", "Dead"}[s]
}

type Player struct {
	Name           string
	State          State
	Position       Vector
	Movement       Vector
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
	CannonCooldown int
	Sprite         *ebiten.Image
	Color          *Color
	AI             AI
	Animations     map[gfx.AnimationType]*gfx.Animation
}

func (p Player) CollisionBox() CollisionBox {
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
