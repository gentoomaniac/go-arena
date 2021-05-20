package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

type Entity interface {
	Hitbox() image.Rectangle
	CollisionBox() image.Rectangle
	Sprite() *ebiten.Image
}
