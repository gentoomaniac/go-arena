package main

import (
	"math/rand"

	"github.com/gentoomaniac/go-arena/entities"
)

type TestBot struct {
	orientation float64
}

func (t *TestBot) Init() {
	t.orientation = 0.3
}

func (t *TestBot) Compute(input entities.AIInput) entities.AIOutput {
	shoot := false
	orientation := t.orientation

	if input.Collided {
		orientation = -10 - float64(rand.Int()%10)
	}

	if len(input.Enemy) > 0 {
		enemy := input.Enemy[0]
		for _, e := range input.Enemy[1:] {
			if e.Distance < enemy.Distance {
				enemy = e
			}
		}
		if enemy.State == entities.Alive {
			orientation = enemy.Angle
			shoot = true
		}
	}

	return entities.AIOutput{
		Speed:             10,
		OrientationChange: orientation,
		Shoot:             shoot,
	}
}

func (t TestBot) Name() string {
	return "TestBot"
}

var Bot TestBot
