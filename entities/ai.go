package entities

import "github.com/gentoomaniac/go-arena/vector"

type AIInput struct {
	Position     vector.Vec2
	CurrentSpeed float64
	TargetSpeed  float64
	MaxSpeed     float64
	Orientation  float64
	Collided     bool
	Hit          bool
	CannonReady  bool
	Enemy        []*Enemy
}

type AIOutput struct {
	Speed             float64 `json:"speed"`
	OrientationChange float64 `json:"orientationChange"`
	Shoot             bool    `json:"shoot"`
}

type AI interface {
	Compute(AIInput) AIOutput
	Init()
	Name() string
}
