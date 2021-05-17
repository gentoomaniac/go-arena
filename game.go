package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math"

	_ "embed"

	"github.com/gentoomaniac/ebitmx"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/rs/zerolog/log"
)

func NewGame() *Game {
	return &Game{}
}

type Game struct {
	arenaMap      *ebitmx.TmxMap
	scalingFactor float64
	screenBuffer  *ebiten.Image
	players       []*Player
	tick          int
}

//go:embed tank.png
var tankImage []byte

func (g *Game) Init() (err error) {
	log.Debug().Msg("init()")
	g.screenBuffer, err = ebiten.NewImage(g.arenaMap.PixelWidth, g.arenaMap.PixelHeight, ebiten.FilterDefault)
	img, err := png.Decode(bytes.NewReader(tankImage))
	if err != nil {
		return
	}
	eimg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		return
	}

	scalingFactor := 4.0
	playerOp := &ebiten.DrawImageOptions{}
	// to scale the imageplayer
	playerOp.GeoM.Translate(float64(-eimg.Bounds().Dx()/2), float64(-eimg.Bounds().Dy()/2))
	playerOp.GeoM.Scale(scalingFactor, scalingFactor)
	playerOp.GeoM.Translate(float64(eimg.Bounds().Dx()/2), float64(eimg.Bounds().Dy()/2))

	playerSprite, err := ebiten.NewImage(eimg.Bounds().Dx()*int(scalingFactor), eimg.Bounds().Dy()*int(scalingFactor), ebiten.FilterDefault)
	log.Debug().Msgf("%s", playerSprite.Bounds())
	if err != nil {
		return
	}
	playerSprite.DrawImage(eimg, playerOp)
	g.players = append(g.players, &Player{
		Name:           "TestBot",
		Position:       image.Point{X: 1000, Y: 1000},
		Health:         100,
		MaxHealth:      100,
		Energy:         100,
		MaxEnergy:      100,
		Speed:          10,
		Orientation:    0,
		Sprite:         playerSprite,
		ColisionBounds: image.Rect(0, 0, playerSprite.Bounds().Dx(), playerSprite.Bounds().Dy()),
		AI:             &TestAI{},
	})
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

func (g *Game) updatePlayer(p *Player) {
	output := p.AI.Compute(AIInput{
		Position:     p.Position,
		Speed:        p.Speed,
		CurrentSpeed: p.CurrentSpeed,
		Orientation:  p.Orientation,
	})
	//jsonOutput, _ := json.Marshal(output)
	//log.Debug().RawJSON("output", jsonOutput).Msg("bot Compute() result")
	p.Speed = output.Speed
	p.Orientation = (p.Orientation + output.Orientation) % 360
	log.Debug().Int("orientation", p.Orientation).Msg("")

	var playerVector image.Point
	if p.Speed > 0 {
		playerVector = image.Point{
			X: int(float64(p.Speed) * math.Cos(float64(p.Orientation))),
			Y: int(float64(p.Speed) * math.Sin(float64(p.Orientation))),
		}
	}

	collisionPointX := image.Point{p.Position.X + playerVector.X, p.Position.Y}
	collisionPointY := image.Point{p.Position.X, p.Position.Y + playerVector.Y}
	collisionPoint := image.Point{p.Position.X + playerVector.X, p.Position.Y + playerVector.Y}
	if !g.arenaMap.CheckColisionPoint(collisionPoint) {
		p.Position.X += playerVector.X
		p.Position.Y += playerVector.Y
	} else if !g.arenaMap.CheckColisionPoint(collisionPointX) {
		p.Position.X += playerVector.X
	} else if !g.arenaMap.CheckColisionPoint(collisionPointY) {
		p.Position.Y += playerVector.Y
	}
	//log.Debug().Str("name", p.Name).Str("oldPos", p.Position.String()).Str("newPos", newPosition.String()).Msg("position update")
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.tick%10 == 0 {
		for _, player := range g.players {
			g.updatePlayer(player)
		}
		g.tick = 1
	} else {
		g.tick++
	}
	return nil
}

func RotateImgOpts(img *ebiten.Image, op *ebiten.DrawImageOptions, degrees int) *ebiten.DrawImageOptions {
	// Move the image's center to the screen's upper-left corner.
	// This is a preparation for rotating. When geometry matrices are applied,
	// the origin point is the upper-left corner.
	op.GeoM.Translate(-float64(img.Bounds().Dx())/2, -float64(img.Bounds().Dy())/2)

	// Rotate the image. As a result, the anchor point of this rotate is
	// the center of the image.
	op.GeoM.Rotate(float64(degrees%360) * 2 * math.Pi / 360)

	// Move the image to the screen's center.
	op.GeoM.Translate(float64(img.Bounds().Dx())/2, float64(img.Bounds().Dy())/2)

	return op
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, layer := range g.arenaMap.Layers {
		err := g.screenBuffer.DrawImage(layer.Render(g.arenaMap, g.scalingFactor, false), &ebiten.DrawImageOptions{})
		if err != nil {
			log.Error().Err(err).Str("layer", layer.Name).Msg("rendering layer failed")
		}
	}

	collisionOp := &ebiten.DrawImageOptions{}
	collisionOp.ColorM.Scale(1, 0, 0, .75)
	err := g.screenBuffer.DrawImage(g.arenaMap.GetObjectGroupByName("collisionmap").DebugRender(g.arenaMap, g.scalingFactor), collisionOp)
	if err != nil {
		log.Debug().Err(err).Msg("rendering collisionmap failed")
	}

	// ======== Draw Player =========
	for _, player := range g.players {
		playerOp := ebiten.DrawImageOptions{}
		//RotateImgOpts(player.Sprite, &playerOp, int(player.Orientation))
		// // Player.Position is absolute on the Map, the coordinates here need to be relative to the camera 0/0
		// // Here is an edge case when the Camera is bigger than the map, stuff breaks
		//playerProjectedX := int(float64(player.Position.X-g.arenaMap.CameraPosition.X)*g.scalingFactor + float64(g.arenaMap.CameraBounds.Max.X)/2)
		//playerProjectedY := int(float64(player.Position.Y-g.arenaMap.CameraPosition.Y)*g.scalingFactor + float64(g.arenaMap.CameraBounds.Max.Y)/2)

		// // to move the image
		playerOp.GeoM.Translate((float64(player.Position.X) - float64(player.Sprite.Bounds().Dx()/2)), (float64(player.Position.Y - player.Sprite.Bounds().Dy()/2)))

		if err = g.screenBuffer.DrawImage(player.Sprite, &playerOp); err != nil {
			log.Error().Err(err).Msg("failed drawing player sprite")
			return
		}
	}

	// ======== Screenbuffer ========

	scaledScreenOp := &ebiten.DrawImageOptions{}
	scaledScreenOp.GeoM.Scale(g.scalingFactor, g.scalingFactor)
	err = screen.DrawImage(g.screenBuffer, scaledScreenOp)
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
