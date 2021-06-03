package main

import (
	"math/rand"

	"github.com/gentoomaniac/go-arena/entities"
)

type GentooBot struct {
}

func (g *GentooBot) Init() {}

func (g *GentooBot) Compute(input entities.AIInput) entities.AIOutput {
	shoot := false
	orientation := 0.2
	speed := float64(15 + rand.Int()%10)
	if input.Collided {
		orientation = 10 + float64(rand.Int()%5)
	}

	if input.CannonReady {
		shoot = true
	}

	return entities.AIOutput{Speed: speed, OrientationChange: orientation, Shoot: shoot}
}

func (g *GentooBot) Name() string {
	return "Gentoobot"
}

var Bot GentooBot
