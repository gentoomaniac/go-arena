package ui

import (
	"bytes"
	_ "embed"
	"image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

var (
	Fade     = true
	FadeTime = 1000 * time.Millisecond
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

func NewStats(headline string) *Stats {
	s := &Stats{headline: NewText(headline)}

	return s
}

type Stats struct {
	headline *Text
	cache    *ebiten.Image
}

func (s *Stats) Image(refresh bool) *ebiten.Image {
	if s.cache == nil || refresh {
		op := &ebiten.DrawImageOptions{}

		s.cache = ebiten.NewImage(bgImage.Size())
		s.cache.DrawImage(bgImage, op)
		headlineImg := s.headline.Image(false)

		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(
			float64(s.cache.Bounds().Dx()/2)-float64(headlineImg.Bounds().Dx()/2)*0.5,
			50,
		)
		s.cache.DrawImage(headlineImg, op)
	}
	return s.cache
}
