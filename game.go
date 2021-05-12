package main

import (
	"fmt"

	"github.com/gentoomaniac/ebitmx"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func NewGame() *Game {
	return &Game{}
}

type Game struct {
	arenaMap      *ebitmx.TmxMap
	scalingFactor float64
	screenBuffer  *ebiten.Image
}

func (g *Game) Init() (err error) {
	g.screenBuffer, err = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	return
}

func (g *Game) WithMap(tmxMap *ebitmx.TmxMap) *Game {
	g.arenaMap = tmxMap
	return g
}

func (g *Game) WithScalingFactor(s float64) *Game {
	g.scalingFactor = s
	return g
}

func (g *Game) Update(screen *ebiten.Image) error {

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	scaledScreenOp := &ebiten.DrawImageOptions{}
	scaledScreenOp.GeoM.Scale(g.scalingFactor, g.scalingFactor)
	for _, layer := range g.arenaMap.Layers {
		err := g.screenBuffer.DrawImage(layer.Render(g.arenaMap, g.scalingFactor, false), scaledScreenOp)
		if err != nil {
			log.Error().Err(err).Str("layer", layer.Name).Msg("rendering layer failed")
		}
	}

	scaledScreenOp.ColorM.Scale(1, 1, 1, 0.5)
	err := g.screenBuffer.DrawImage(g.arenaMap.GetObjectGroupByName("collisionmap").DebugRender(g.arenaMap, g.scalingFactor), scaledScreenOp)
	if err != nil {
		log.Debug().Err(err).Msg("rendering collisionmap failed")
	}

	// ======== Screenbuffer ========
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0.0, 0.0)
	err = screen.DrawImage(g.screenBuffer, op)
	if err != nil {
		log.Debug().Err(err).Msg("rendering screen failed")
	}

	// ======== Info ========
	err = ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	if err != nil {
		log.Debug().Err(err).Msg("writing Debug message failed")
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
