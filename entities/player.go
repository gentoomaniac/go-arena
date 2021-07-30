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
	ID             int
	Name           string
	State          State
	Position       vector.Vec2
	Acceleration   float64
	Friction       float64
	Velocity       vector.Vec2
	CurrentSpeed   float64
	ImpactVelocity vector.Vec2
	Mass           float64
	Health         int
	MaxHealth      int
	Energy         int
	MaxEnergy      int
	//CurrentSpeed     float64
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

	if p.CurrentSpeed > p.TargetSpeed {
		p.CurrentSpeed -= p.Acceleration
	} else if p.CurrentSpeed < p.TargetSpeed {
		p.CurrentSpeed += p.Acceleration
	}
}

func (p *Player) UpdateVelocity() {
	p.Velocity = vector.Vec2{
		X: p.CurrentSpeed * math.Cos(p.Orientation*math.Pi/180),
		Y: p.CurrentSpeed * math.Sin(p.Orientation*math.Pi/180),
	}
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
