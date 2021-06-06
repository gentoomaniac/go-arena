package entities

import (
	"bytes"
	"image/png"

	_ "embed"

	"github.com/gentoomaniac/go-arena/vector"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

//go:embed shell.png
var shellPng []byte
var shellImage *ebiten.Image

func init() {
	img, err := png.Decode(bytes.NewReader(shellPng))
	if err != nil {
		log.Error().Err(err).Msg("could not load shell image")
	}

	shellImage = ebiten.NewImageFromImage(img)
}

func NewShell() *Shell {
	s := &Shell{}
	s.sprite = shellImage
	s.collisionBox = vector.Rectangle{
		Min: vector.Vec2{X: -float64(shellImage.Bounds().Max.X) / 2, Y: -float64(shellImage.Bounds().Max.Y) / 2},
		Max: vector.Vec2{X: float64(shellImage.Bounds().Max.X) / 2, Y: float64(shellImage.Bounds().Max.Y) / 2},
	}

	return s
}

type Shell struct {
	name         string
	collisionBox vector.Rectangle
	position     vector.Vec2
	speed        int
	orientation  float64
	sprite       *ebiten.Image
	Damage       int
	Source       *Player
}

func (s Shell) Name() string {
	return s.name
}

func (s Shell) CollisionBox() vector.Rectangle {
	return vector.Rectangle{
		Min: vector.Vec2{
			X: s.position.X + s.collisionBox.Min.X,
			Y: s.position.Y + s.collisionBox.Min.Y,
		},
		Max: vector.Vec2{
			X: s.position.X + s.collisionBox.Max.X,
			Y: s.position.Y + s.collisionBox.Max.Y,
		},
	}
}

func (s Shell) Position() vector.Vec2 {
	return s.position
}
func (s *Shell) SetPosition(p vector.Vec2) {
	s.position = p
}

func (s Shell) Speed() int {
	return s.speed
}
func (s *Shell) SetSpeed(speed int) {
	s.speed = speed
}

func (s Shell) Orientation() float64 {

	return s.orientation
}
func (s *Shell) SetOrientation(o float64) {
	s.orientation = o
}

func (s Shell) Sprite() *ebiten.Image {
	return s.sprite
}
