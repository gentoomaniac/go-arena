package entities

import (
	"github.com/gentoomaniac/go-arena/gfx"
	"github.com/gentoomaniac/go-arena/vector"
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
	ID           int
	Name         string
	State        State
	Position     vector.Vec2
	Acceleration float64
	Friction     float64
	Velocity     vector.Vec2
	Orientation  vector.Vec2
	Mass         float64
	Health       int
	MaxHealth    int
	Energy       int
	MaxEnergy    int
	//CurrentSpeed     float64
	TargetSpeed      float64
	MaxSpeed         float64
	CollisionRadius  float64
	Collided         bool
	CollidedWithTank bool
	CannonCooldown   int
	Hit              bool
	Sprite           *ebiten.Image
	Color            *gfx.Color
	AI               AI
	Animations       map[gfx.AnimationType]*gfx.Animation
	NumberRespawns   int
	MaxRespawns      int
	RespawnCooldown  int
}

func (p *Player) UpdateSpeed(newSpeed float64) {
	p.TargetSpeed = newSpeed
	if p.TargetSpeed > p.MaxSpeed {
		p.TargetSpeed = p.MaxSpeed
	} else if p.TargetSpeed < -p.MaxSpeed {
		p.TargetSpeed = -p.MaxSpeed
	}

	if p.Velocity.Length() == 0 {
		if newSpeed > p.Acceleration {
			p.Velocity = vector.FromAngle(p.Orientation.Angle(), p.Acceleration)
		} else {
			p.Velocity = vector.FromAngle(p.Orientation.Angle(), newSpeed)
		}
	} else if p.Velocity.Length() > p.TargetSpeed {
		p.Velocity = p.Velocity.Unit().WithLength(p.Velocity.Length() - p.Acceleration)
	} else if p.Velocity.Length() < p.TargetSpeed {
		p.Velocity = p.Velocity.Unit().WithLength(p.Velocity.Length() + p.Acceleration)
	}
}

func (p *Player) UpdateOrientation(angle float64) {
	p.Orientation = p.Velocity.Rotate(angle)
	p.Velocity = p.Velocity.Rotate(angle)
}
