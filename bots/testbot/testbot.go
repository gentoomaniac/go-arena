package main

import (
	"fmt"
	"math/rand"

	"github.com/gentoomaniac/go-arena/entities"
)

type TestBot struct {
	orientation float64
	speed       float64
}

func (t *TestBot) Init() {
	t.orientation = 0.3
	t.speed = 5
}

func (t *TestBot) Compute(input entities.AIInput) entities.AIOutput {
	shoot := false
	orientation := t.orientation
	speed := t.speed

	if input.Collided {
		orientation = -10 - float64(rand.Int()%10)
		speed = 5
	}

	if input.CollidedWithTank {
		speed = 0
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
			speed = 10
		}
	}

	return entities.AIOutput{
		Speed:             speed,
		OrientationChange: orientation,
		Shoot:             shoot,
	}
}

func (t TestBot) Name() string {
	return fmt.Sprintf("TestBot %d", rand.Int()%10)
}

var Bot TestBot
