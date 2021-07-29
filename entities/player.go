package entities

import (
	"math"

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

func (p *Player) UpdateSpeed(newSpeed float64) {
	p.TargetSpeed = newSpeed
	if p.TargetSpeed > p.MaxSpeed {
		p.TargetSpeed = p.MaxSpeed
	} else if p.TargetSpeed < -p.MaxSpeed {
		p.TargetSpeed = -p.MaxSpeed
	}

	if diff := p.CurrentSpeed - p.TargetSpeed; diff > p.Friction {
		p.CurrentSpeed += p.Friction
	} else if diff > p.Acceleration {
		p.CurrentSpeed += p.Acceleration
	} else {
		p.CurrentSpeed = p.TargetSpeed
	}

	p.Velocity.X = p.CurrentSpeed * math.Cos(p.Orientation*math.Pi/180)
	p.Velocity.Y = p.CurrentSpeed * math.Sin(p.Orientation*math.Pi/180)
}
