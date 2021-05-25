package main

import (
	"math/rand"

	"github.com/gentoomaniac/go-arena/entities"
)

type TestBot struct{}

func (t *TestBot) Compute(input entities.AIInput) entities.AIOutput {
	orientation := -0.3
	if input.Collided {
		orientation = -10 - float64(rand.Int()%10)
	}

	return entities.AIOutput{Speed: 10, OrientationChange: orientation, Shoot: true}
}

func (t TestBot) Name() string {
	return "TestBot"
}

var Bot TestBot
