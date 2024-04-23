package gfx

import (
	"bytes"
	_ "embed"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed frame.png
var frameRawImage []byte

func LoadFrameSprite() (*ebiten.Image, error) {
	img, err := png.Decode(bytes.NewReader(frameRawImage))
	if err != nil {
		return nil, err
	}
	eimg := ebiten.NewImageFromImage(img)

	frameOp := &ebiten.DrawImageOptions{}
	frameOp.GeoM.Translate(float64(-eimg.Bounds().Dx()/2), float64(-eimg.Bounds().Dy()/2))
	frameOp.GeoM.Scale(tankScalingFactor, tankScalingFactor)
	frameOp.GeoM.Translate(float64(eimg.Bounds().Dx()/2*int(tankScalingFactor)), float64(eimg.Bounds().Dy()/2*int(tankScalingFactor)))

	frameImage := ebiten.NewImage(eimg.Bounds().Dx()*int(tankScalingFactor), eimg.Bounds().Dy()*int(tankScalingFactor))

	frameImage.DrawImage(eimg, frameOp)
	return frameImage, nil
}
