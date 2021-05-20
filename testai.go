package main

import "math/rand"

type TestAI struct {
}

func (t *TestAI) Compute(input AIInput) AIOutput {
	orientation := 0.2
	if input.Collided {
		orientation = 10 + float64(rand.Int()%10)
	}

	return AIOutput{Speed: 20, OrientationChange: orientation}
}
