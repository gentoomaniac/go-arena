package ui

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/png"

	"github.com/gentoomaniac/go-arena/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

var (
	HeadlineScaling = 0.5
	TextScaling     = 0.25
	MarginTop       = 50.0
	MarginLeft      = 75.0
	Spacer          = 30.0
)

//go:embed background.png
var background []byte
var bgImage *ebiten.Image

func init() {
	img, err := png.Decode(bytes.NewReader(background))
	if err != nil {
		log.Error().Err(err).Msg("could not load shell image")
	}

	bgImage = ebiten.NewImageFromImage(img)
}

func NewStats(headline string, players []*entities.Player) *Stats {
	return &Stats{headline: NewText(headline), players: players}
}

type Stats struct {
	headline *Text
	cache    *ebiten.Image
	players  []*entities.Player
}

func (s *Stats) Image(refresh bool) *ebiten.Image {
	if s.cache == nil || refresh {
		op := &ebiten.DrawImageOptions{}

		s.cache = ebiten.NewImage(bgImage.Size())
		s.cache.DrawImage(bgImage, op)
		headlineImg := s.headline.Image(false)

		op.GeoM.Scale(HeadlineScaling, HeadlineScaling)
		op.GeoM.Translate(
			float64(s.cache.Bounds().Dx()/2)-float64(headlineImg.Bounds().Dx()/2)*HeadlineScaling,
			MarginTop,
		)
		s.cache.DrawImage(headlineImg, op)

		op.GeoM.Reset()
		op.GeoM.Scale(TextScaling, TextScaling)
		op.GeoM.Translate(MarginLeft, MarginTop+float64(headlineImg.Bounds().Dy())*HeadlineScaling+Spacer)
		for index, p := range s.players {
			text := NewText(fmt.Sprintf("#%d %s %d", index+1, p.Name, p.Health))
			s.cache.DrawImage(text.Image(false), op)
			op.GeoM.Translate(0, float64(text.image.Bounds().Dy())*TextScaling+Spacer)
		}
	}
	return s.cache
}
