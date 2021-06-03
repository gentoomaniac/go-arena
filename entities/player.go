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
	Name            string
	State           State
	Position        Vector
	Movement        Vector
	Hitbox          image.Rectangle
	Health          int
	MaxHealth       int
	Energy          int
	MaxEnergy       int
	CurrentSpeed    float64
	TargetSpeed     float64
	MaxSpeed        float64
	Acceleration    float64
	Orientation     float64
	CollisionBounds CollisionBox
	Collided        bool
	CannonCooldown  int
	Hit             bool
	Sprite          *ebiten.Image
	Color           *Color
	AI              AI
	Animations      map[gfx.AnimationType]*gfx.Animation
	NumberRespawns  int
}

func (p Player) CollisionBox() CollisionBox {
	return CollisionBox{
		Min: Vector{
			p.Position.X + p.CollisionBounds.Min.X,
			p.Position.Y + p.CollisionBounds.Min.Y,
		},
		Max: Vector{
			p.Position.X + p.CollisionBounds.Max.X,
			p.Position.Y + p.CollisionBounds.Max.Y,
		},
	}
}

func (p *Player) UpdateSpeed(newSpeed float64) {
	p.TargetSpeed = newSpeed
	if p.TargetSpeed > p.MaxSpeed {
		p.TargetSpeed = p.MaxSpeed
	} else if p.TargetSpeed < -p.MaxSpeed {
		p.TargetSpeed = -p.MaxSpeed
	}

	if p.CurrentSpeed > p.TargetSpeed {
		p.CurrentSpeed -= p.Acceleration
	} else if p.CurrentSpeed < p.TargetSpeed {
		p.CurrentSpeed += p.Acceleration
	}
}
