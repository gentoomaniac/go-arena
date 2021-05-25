package main

import (
	"image"

	"github.com/gentoomaniac/ebitmx"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

const mapPath = "maps/test.tmx"
const (
	screenWidth   = 965
	screenHeight  = 965
	scalingFactor = .15
)

func run(bots []string) {
	tmxMap, error := ebitmx.LoadFromFile(mapPath)
	if error != nil {
		log.Fatal().Err(error).Msg("")
	}

	startPosition := image.Point{0, 0}

	tmxMap.CameraBounds = image.Rect(0, 0, 6400, 6400)
	tmxMap.CameraPosition = startPosition
	log.Debug().Int("width", tmxMap.PixelWidth).Int("height", tmxMap.PixelHeight).Msg("map dimensions")

	game := NewGame().WithMap(tmxMap).WithScalingFactor(scalingFactor).WithBots(bots)
	err := game.Init()
	if err != nil {
		log.Error().Err(err).Msg("initialising game failed")
		return
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("go-arena")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
