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
	CurrentSpeed   float64
	TargetSpeed    float64
	MaxSpeed       float64
	Acceleration   float64
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

	//log.Debug().Float64("targetSpeed", p.TargetSpeed).Float64("currentSpeed", p.CurrentSpeed).Str("name", p.Name).Msg("")

}
