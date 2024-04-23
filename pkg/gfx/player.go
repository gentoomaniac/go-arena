package gfx

import (
	"bytes"
	_ "embed"
	"image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

var (
	tankScalingFactor = 4.0
)

//go:embed tank.png
var tankImage []byte

func GetPlayerSprite() (*ebiten.Image, error) {
	img, err := png.Decode(bytes.NewReader(tankImage))
	if err != nil {
		return nil, err
	}
	eimg := ebiten.NewImageFromImage(img)

	playerOp := &ebiten.DrawImageOptions{}
	// to scale the imageplayer
	playerOp.GeoM.Translate(float64(-eimg.Bounds().Dx()/2), float64(-eimg.Bounds().Dy()/2))
	playerOp.GeoM.Scale(tankScalingFactor, tankScalingFactor)
	playerOp.GeoM.Rotate(90 * math.Pi / 180)
	playerOp.GeoM.Translate(float64(eimg.Bounds().Dx()/2*int(tankScalingFactor)), float64(eimg.Bounds().Dy()/2*int(tankScalingFactor)))

	playerSprite := ebiten.NewImage(eimg.Bounds().Dx()*int(tankScalingFactor), eimg.Bounds().Dy()*int(tankScalingFactor))
	log.Debug().Msgf("playerSprite: %s", playerSprite.Bounds())

	playerSprite.DrawImage(eimg, playerOp)

	return playerSprite, nil
}
