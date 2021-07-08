package main

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
