package vector

import (
	"testing"
)

func TestLength(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Length()
			if result > tt.want*(1+maxError) || result < tt.want*(1-maxError) {
				t.Errorf("result exceeds error threshold, got '%f' want '%f'", result, tt.want)
			}
		})
	}

}

func TestUnit(t *testing.T) {
	maxError := 0.01
	var tests = []struct {
		name string
		v    Vec2
		want Vec2
	}{
		{
			"zero vector", Vec2{0, 0}, Vec2{},
		},
		{
			"normal vector", Vec2{1, 1}, Vec2{0.7071, 0.7071},
		},
		{
			"arbitrary vector", Vec2{5, 1}, Vec2{0.98058, 0.196116},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Unit()
			if result.X > tt.want.X*(1+maxError) || result.X < tt.want.X*(1-maxError) ||
				result.Y > tt.want.Y*(1+maxError) || result.Y < tt.want.Y*(1-maxError) {
				t.Errorf("result exceeds error threshold, got '%f' want '%f'", result, tt.want)
			}
		})
	}

}
