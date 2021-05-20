package main

import (
	"math/rand"

	"github.com/gentoomaniac/go-arena/player"
)

type TestBot struct{}

func (t *TestBot) Compute(input player.AIInput) player.AIOutput {
	orientation := -0.3
	if input.Collided {
		orientation = -10 - float64(rand.Int()%10)
	}

	return player.AIOutput{Speed: 10, OrientationChange: orientation}
}

func (t TestBot) Name() string {
	return "TestBot"
}

var Bot TestBot
