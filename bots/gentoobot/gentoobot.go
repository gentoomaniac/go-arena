package main

import (
	"math/rand"

	"github.com/gentoomaniac/go-arena/player"
)

type GentooBot struct {
}

func (g *GentooBot) Compute(input player.AIInput) player.AIOutput {
	orientation := 0.2
	if input.Collided {
		orientation = 10 + float64(rand.Int()%5)
	}

	return player.AIOutput{Speed: 20, OrientationChange: orientation}
}

func (g *GentooBot) Name() string {
	return "Gentoobot"
}

var Bot GentooBot
