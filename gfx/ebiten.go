package gfx

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func RotateImgOpts(img *ebiten.Image, op ebiten.DrawImageOptions, degrees int) ebiten.DrawImageOptions {
	op.GeoM.Translate(-float64(img.Bounds().Dx())/2, -float64(img.Bounds().Dy())/2)
	op.GeoM.Rotate(float64(degrees%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(float64(img.Bounds().Dx())/2, float64(img.Bounds().Dy())/2)

	return op
}
