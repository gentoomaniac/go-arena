package entities

type AIInput struct {
	Position     Vector
	Speed        int
	CurrentSpeed int
	Orientation  float64
	Collided     bool
	CannonReady  bool
	Enemy        []*Enemy
}

type AIOutput struct {
	Speed             int     `json:"speed"`
	OrientationChange float64 `json:"orientationChange"`
	Shoot             bool    `json:"shoot"`
}

type AI interface {
	Compute(AIInput) AIOutput
	Init()
	Name() string
}
