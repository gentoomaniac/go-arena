package entities

import (
	"bytes"
	"image/png"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shell.png
var shellPng []byte

func loadSprite() (*ebiten.Image, error) {
	img, err := png.Decode(bytes.NewReader(shellPng))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func NewShell() (s *Shell, err error) {
	s = &Shell{}
	sprite, err := loadSprite()
	if err != nil {
		return nil, err
	}

	s.SetSprite(sprite)

	return s, nil
}

type Shell struct {
	name         string
	collisionBox CollisionBox
	position     Vector
	speed        int
	orientation  float64
	sprite       *ebiten.Image
}

func (s Shell) Name() string {
	return s.name
}

func (s Shell) CollisionBox() CollisionBox {
	return s.collisionBox
}

func (s Shell) Position() Vector {
	return s.position
}
func (s *Shell) SetPosition(p Vector) {
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
func (s *Shell) SetSprite(sprite *ebiten.Image) {
	s.sprite = sprite
}
