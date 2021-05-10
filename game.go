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

	// ======== Map ========
	scaledScreenOp := &ebiten.DrawImageOptions{}
	scaledScreenOp.GeoM.Scale(g.scalingFactor, g.scalingFactor)
	for i := range g.arenaMap.Layers {
		g.screenBuffer.DrawImage(g.arenaMap.Layers[i].Render(g.arenaMap, g.scalingFactor, false), scaledScreenOp)
	}

	scaledScreenOp.ColorM.Scale(1, 1, 1, 0.5)
	g.screenBuffer.DrawImage(g.arenaMap.GetObjectGroupByName("collisionmap").DebugRender(g.arenaMap, g.scalingFactor), scaledScreenOp)

	// ======== Screenbuffer ========

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0.0, 0.0)
	screen.DrawImage(g.screenBuffer, op)

	// ======== Info ========
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
