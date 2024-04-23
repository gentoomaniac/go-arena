package entities

type Enemy struct {
	Angle    float64
	Distance float64
	Health   int
	Speed    int
	State    State
}
