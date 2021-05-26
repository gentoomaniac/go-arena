package main

import (
	"math/rand"

	"github.com/gentoomaniac/go-arena/entities"
)

type TestBot struct{}

func (t *TestBot) Compute(input entities.AIInput) entities.AIOutput {
	shoot := false
	orientation := -0.3
	if input.Collided {
		orientation = -10 - float64(rand.Int()%10)
	}
	if len(input.Enemy) > 0 {
		enemy := input.Enemy[0]
		if enemy.State == entities.Alive {
			orientation = enemy.Angle
			shoot = true
		}
	}

	return entities.AIOutput{Speed: 10, OrientationChange: orientation, Shoot: shoot}
}

func (t TestBot) Name() string {
	return "TestBot"
}

var Bot TestBot
