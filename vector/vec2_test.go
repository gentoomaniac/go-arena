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

func TestRotate(t *testing.T) {
	maxError := 0.01
	var tests = []struct {
		name  string
		v     Vec2
		angle float64
		want  Vec2
	}{
		{
			"0", Vec2{1, 1}, 0, Vec2{1, 1},
		},
		{
			"90", Vec2{1, 1}, 90, Vec2{-1, 1},
		},
		{
			"180", Vec2{1, 1}, 180, Vec2{-1, -1},
		},
		{
			"360", Vec2{1, 1}, 360, Vec2{1, 1},
		},
		{
			"123", Vec2{45, 16}, 123, Vec2{-37.9274856628, 29.0259509973},
		},
	}

	error := func(g Vec2, w Vec2) {
		t.Errorf("result exceeds error threshold, got '%f' want '%f'", g, w)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Rotate(tt.angle)
			if tt.want.X >= 0 {
				if result.X > tt.want.X*(1+maxError) || result.X < tt.want.X*(1-maxError) {
					error(result, tt.want)
				}
			} else {
				if result.X < tt.want.X*(1+maxError) || result.X > tt.want.X*(1-maxError) {
					error(result, tt.want)
				}
			}

			if tt.want.Y >= 0 {
				if result.Y > tt.want.Y*(1+maxError) || result.Y < tt.want.Y*(1-maxError) {
					error(result, tt.want)
				}
			} else {
				if result.Y < tt.want.Y*(1+maxError) || result.Y > tt.want.Y*(1-maxError) {
					error(result, tt.want)
				}
			}
		})
	}
}

func TestAngle(t *testing.T) {
	maxError := 0.01
	var tests = []struct {
		name string
		v    Vec2
		want float64
	}{
		{
			"zero vector", Vec2{0, 0}, 0.0,
		},
		{
			"45", Vec2{1, 1}, 45.0,
		},
		{
			"arbitrary", Vec2{5, 1}, 11.309932474020195,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Angle()
			if result > tt.want*(1+maxError) {
				t.Errorf("result exceeds error threshold, got '%f' want '%f'", result, tt.want)
			}
		})
	}
}
func TestFromAngle(t *testing.T) {
	maxError := 0.01
	var tests = []struct {
		name   string
		angle  float64
		length float64
		want   Vec2
	}{
		{
			"zero vector", 45, 0, Vec2{},
		},
		{
			"1/1 vector", 45, 1.4142, Vec2{1, 1},
		},
		{
			"zero degree", 0, 1, Vec2{1, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromAngle(tt.angle, tt.length)
			if result.X > tt.want.X*(1+maxError) || result.X < tt.want.X*(1-maxError) ||
				result.Y > tt.want.Y*(1+maxError) || result.Y < tt.want.Y*(1-maxError) {
				t.Errorf("result exceeds error threshold, got '%f' want '%f'", result, tt.want)
			}
		})
	}
}
