package entities

import (
	"github.com/gentoomaniac/go-arena/gfx"
	"github.com/gentoomaniac/go-arena/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

func NewShell() *Shell {
	s := &Shell{}
	s.sprite = gfx.GetShellImage()
	s.collisionRadius = float64(s.sprite.Bounds().Dx()) / 2

	return s
}

type Shell struct {
	name            string
	collisionRadius float64
	Position        vector.Vec2
	Movement        vector.Vec2
	Orientation     float64
	sprite          *ebiten.Image
	Damage          int
	Source          *Player
}

func (s Shell) Name() string {
	return s.name
}

func (s Shell) Sprite() *ebiten.Image {
	return s.sprite
}
