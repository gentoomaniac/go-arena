package go_arena

import (
	"image"
	"math/rand"
	"time"

	"github.com/gentoomaniac/ebitmx"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

const (
	screenWidth   = 965
	screenHeight  = 965
	scalingFactor = .15
)

type RunOpts struct {
	Bots     []string
	Respawns int
	MapPath  string
}

func Run(opts RunOpts) {
	tmxMap, error := ebitmx.LoadFromFile(opts.MapPath)
	if error != nil {
		log.Fatal().Err(error).Msg("")
	}

	startPosition := image.Point{0, 0}

	tmxMap.CameraBounds = image.Rect(0, 0, 6400, 6400)
	tmxMap.CameraPosition = startPosition
	log.Debug().Int("width", tmxMap.PixelWidth).Int("height", tmxMap.PixelHeight).Msg("map dimensions")

	rand.Seed(time.Now().UTC().UnixNano())
	game := NewGame().WithMap(tmxMap).WithScalingFactor(scalingFactor).WithRespawns(opts.Respawns).WithBots(opts.Bots)
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
