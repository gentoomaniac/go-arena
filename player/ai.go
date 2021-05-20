package player

type AIInput struct {
	Position     Vector
	Speed        int
	CurrentSpeed int
	Orientation  float64
	Collided     bool
}

type AIOutput struct {
	Speed             int     `json:"speed"`
	OrientationChange float64 `json:"orientationChange"`
}

type AI interface {
	Compute(AIInput) AIOutput
	Name() string
}
