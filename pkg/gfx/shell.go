package gfx

import (
	"bytes"
	_ "embed"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

//go:embed shell.png
var shellPng []byte
var shellImage *ebiten.Image

func GetShellImage() *ebiten.Image {
	if shellImage == nil {
		img, err := png.Decode(bytes.NewReader(shellPng))
		if err != nil {
			log.Error().Err(err).Msg("could not load shell image")
		}

		shellImage = ebiten.NewImageFromImage(img)
	}

	return shellImage
}
