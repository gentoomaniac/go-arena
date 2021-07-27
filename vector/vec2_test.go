package vector

import (
	"fmt"
	"testing"
)

func TestIntersection(t *testing.T) {
	maxError := 0.01
	var tests = []struct {
		name string
		v    Vec2
		want float64
	}{
		{
			"zero vector", Vec2{0, 0}, 0,
		},
		{
			"normal vector", Vec2{1, 1}, 1.41421356237,
		},
		{
			"negative vector", Vec2{-2, -2}, 2.82842712475,
		},
		{
			"arbitrary vector", Vec2{5, 1}, 5.09901951359,
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("vector: %s", tt.v)

		t.Run(testname, func(t *testing.T) {
			result := tt.v.Length()
			if result > tt.want*(1+maxError) || result < tt.want*(1-maxError) {
				t.Errorf("result exceeds error threshold, got '%f' want '%f'", result, tt.want)
			}
		})
	}

}
