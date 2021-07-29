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
	ID               int
	Name             string
	State            State
	Position         vector.Vec2
	Acceleration     float64
	Friction         float64
	Velocity         vector.Vec2
	ImpactVelocity   vector.Vec2
	Mass             float64
	Health           int
	MaxHealth        int
	Energy           int
	MaxEnergy        int
	CurrentSpeed     float64
	TargetSpeed      float64
	MaxSpeed         float64
	Orientation      float64
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

func (p *Player) UpdateVelocity(newSpeed float64) {
	p.TargetSpeed = newSpeed
	if p.TargetSpeed > p.MaxSpeed {
		p.TargetSpeed = p.MaxSpeed
	} else if p.TargetSpeed < -p.MaxSpeed {
		p.TargetSpeed = -p.MaxSpeed
	}

	// ToDo: When the tanks start the velocity is 0 so all calculations will be 0 and tanks not move

	if diff := p.Velocity.Length() - p.TargetSpeed; diff > p.TargetSpeed {
		if diff > p.Friction {
			p.Velocity = p.Velocity.WithLength(p.Velocity.Length() - p.Friction)
		} else {
			p.Velocity = p.Velocity.WithLength(p.TargetSpeed)
		}
	} else {
		if diff > p.Acceleration {
			p.Velocity = p.Velocity.WithLength(p.Velocity.Length() + p.Acceleration)
		} else {
			p.Velocity = p.Velocity.WithLength(p.TargetSpeed)
		}
	}

	//ToDo: need orientation change
}

func (p *Player) UpdateImpactVelocity() {
	len := p.ImpactVelocity.Length()

	if len > p.Friction {
		p.ImpactVelocity = p.ImpactVelocity.WithLength(p.ImpactVelocity.Length() - p.Friction)
	} else {
		p.ImpactVelocity.X = 0
		p.ImpactVelocity.Y = 0
	}
}
