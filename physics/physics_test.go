package physics

import (
	"fmt"
	"testing"

	"github.com/gentoomaniac/go-arena/vector"
)

func TestIntersection(t *testing.T) {
	var tests = []struct {
		name           string
		v1, v2, v3, v4 vector.Vec2
		want           *vector.Vec2
	}{
		{
			"crossing in the middle",
			vector.Vec2{2, 0}, vector.Vec2{2, 4},
			vector.Vec2{0, 2}, vector.Vec2{4, 2},
			&vector.Vec2{2, 2},
		},
		{
			"crossing on outer edge",
			vector.Vec2{4, 0}, vector.Vec2{4, 4},
			vector.Vec2{0, 2}, vector.Vec2{4, 2},
			&vector.Vec2{4, 2},
		},
		{
			"parallels",
			vector.Vec2{0, 0}, vector.Vec2{4, 0},
			vector.Vec2{0, 2}, vector.Vec2{4, 2},
			nil,
		},
		{
			"lines cross outside section",
			vector.Vec2{5, 0}, vector.Vec2{9, 4},
			vector.Vec2{0, 4}, vector.Vec2{4, 0},
			nil,
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("a[%s, %s], b[%s, %s]", tt.v1, tt.v2, tt.v3, tt.v4)

		t.Run(testname, func(t *testing.T) {
			result := Intersection(tt.v1, tt.v2, tt.v3, tt.v4)
			if result == nil && tt.want != nil {
				t.Errorf("got nil, want %s", tt.want)
			}
			if result != nil && tt.want == nil {
				t.Errorf("got %s, want nil", result)
			}
			if result != nil && tt.want != nil {
				if *result != *tt.want {
					t.Errorf("got %s, want %s", result, tt.want)
				}
			}
		})
	}

}

func TestDistance(t *testing.T) {
	maxError := 0.01
	var tests = []struct {
		name string
		a, b vector.Vec2
		want float64
	}{
		{
			"same point",
			vector.Vec2{0, 0},
			vector.Vec2{0, 0},
			0,
		},
		{
			"from origin",
			vector.Vec2{0, 0},
			vector.Vec2{1, 1},
			1.41421356237,
		},
		{
			"negative to positive",
			vector.Vec2{-1, -1},
			vector.Vec2{1, 1},
			2.82842712475,
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s: %s, %s", tt.name, tt.a, tt.b)

		t.Run(testname, func(t *testing.T) {
			result := Distance(tt.a, tt.b)
			if result > tt.want*(1+maxError) || result < tt.want*(1-maxError) {
				t.Errorf("result exceeds error threshold, got '%f' want '%f'", result, tt.want)
			}
		})
	}
}

func TestDoCirclesOverlap(t *testing.T) {
	maxError := 0.01
	var tests = []struct {
		name string
		a, b vector.Circle
		want float64
	}{
		{
			"same position",
			vector.Circle{vector.Vec2{0, 0}, 1},
			vector.Circle{vector.Vec2{0, 0}, 1},
			-2.0,
		},
		{
			"intersecting",
			vector.Circle{vector.Vec2{0, 0}, 1},
			vector.Circle{vector.Vec2{1, 0}, 1},
			-1.0,
		},
		{
			"touching",
			vector.Circle{vector.Vec2{0, 0}, 1},
			vector.Circle{vector.Vec2{2, 0}, 1},
			0.0,
		},
		{
			"not touching",
			vector.Circle{vector.Vec2{0, 0}, 1},
			vector.Circle{vector.Vec2{2.1, 0}, 1},
			0.1,
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s: %s, %s", tt.name, tt.a, tt.b)

		t.Run(testname, func(t *testing.T) {
			result := DistanceBetweenCircles(tt.a, tt.b)
			if result > tt.want*(1+maxError) || result < tt.want*(1-maxError) {
				t.Errorf("result exceeds error threshold, got '%f' want '%f'", result, tt.want)
			}
		})
	}
}

func TestPointLineDistance(t *testing.T) {
	maxError := 0.01
	var tests = []struct {
		name      string
		v1, v2, p vector.Vec2
		want      float64
	}{
		{
			"1/1",
			vector.Vec2{0, 0}, vector.Vec2{1, 1},
			vector.Vec2{1, 0},
			vector.Vec2{1, 1}.Length() / 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PointLineDistance(tt.v1, tt.v2, tt.p)
			if result > tt.want*(1+maxError) || result < tt.want*(1-maxError) {
				t.Errorf("result exceeds error threshold, got '%f' want '%f'", result, tt.want)
			}
		})
	}

}
