package main

import (
	"bytes"
	"fmt"
	"image/png"
	"math"
	"plugin"

	_ "embed"

	"github.com/gentoomaniac/ebitmx"
	"github.com/gentoomaniac/go-arena/entities"
	"github.com/gentoomaniac/go-arena/gfx"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/rs/zerolog/log"
)

var (
	ColisionDamage = 5 // how much health does a player loose on colisions
	CannonCooldown = 60
)

func NewGame() *Game {
	return &Game{}
}

//go:embed tank.png
var tankImage []byte

func getPlayerSprite() (*ebiten.Image, error) {
	img, err := png.Decode(bytes.NewReader(tankImage))
	if err != nil {
		return nil, err
	}
	eimg := ebiten.NewImageFromImage(img)

	scalingFactor := 4.0
	playerOp := &ebiten.DrawImageOptions{}
	// to scale the imageplayer
	playerOp.GeoM.Translate(float64(-eimg.Bounds().Dx()/2), float64(-eimg.Bounds().Dy()/2))
	playerOp.GeoM.Scale(scalingFactor, scalingFactor)
	playerOp.GeoM.Rotate(90 * math.Pi / 180)
	playerOp.GeoM.Translate(float64(eimg.Bounds().Dx()/2*int(scalingFactor)), float64(eimg.Bounds().Dy()/2*int(scalingFactor)))

	playerSprite := ebiten.NewImage(eimg.Bounds().Dx()*int(scalingFactor), eimg.Bounds().Dy()*int(scalingFactor))
	log.Debug().Msgf("playerSprite: %s", playerSprite.Bounds())

	playerSprite.DrawImage(eimg, playerOp)

	return playerSprite, nil
}

type Game struct {
	arenaMap      *ebitmx.TmxMap
	scalingFactor float64
	screenBuffer  *ebiten.Image
	players       []*entities.Player
	shells        []*entities.Shell
}

func (g *Game) Init() (err error) {
	log.Debug().Msg("init()")
	g.screenBuffer = ebiten.NewImage(g.arenaMap.PixelWidth, g.arenaMap.PixelHeight)
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

//go:embed fire_transparent.gif
var fireGif []byte

func (g *Game) WithBots(bots []string) *Game {
	var color *entities.Color
	for index, botModulePath := range bots {
		botPlugin, err := plugin.Open(botModulePath)
		if err != nil {
			log.Error().Err(err).Msg("failed loading bot")
			return nil
		}
		botObj, err := botPlugin.Lookup("Bot")
		if err != nil {
			log.Error().Err(err).Msg("no object called 'Bot' found")
			return nil
		}
		ai, ok := botObj.(entities.AI)
		if !ok {
			log.Error().Err(err).Msg("bot object doesn't implement the AI interface")
			return nil
		}

		playerSprite, err := getPlayerSprite()
		if err != nil {
			return nil
		}

		switch index % 4 {
		case 0:
			color = &entities.Color{R: 1, G: .7, B: .7, Alpha: 1}
		case 1:
			color = &entities.Color{R: 1, G: 1, B: .7, Alpha: 1}
		case 2:
			color = &entities.Color{R: .7, G: 1, B: .7, Alpha: 1}
		case 3:
			color = &entities.Color{R: .7, G: .7, B: 1, Alpha: 1}
		}
		player := &entities.Player{
			Name:        ai.Name(),
			State:       entities.Alive,
			Position:    entities.Vector{X: 1000 * float64(index), Y: 1000 * float64(index)},
			Health:      100,
			MaxHealth:   100,
			Energy:      100,
			MaxEnergy:   100,
			Speed:       10,
			MaxSpeed:    20,
			Orientation: 0,
			Sprite:      playerSprite,
			Color:       color,
			ColisionBounds: entities.CollisionBox{
				Min: entities.Vector{X: 0, Y: 0},
				Max: entities.Vector{X: float64(playerSprite.Bounds().Dx()), Y: float64(playerSprite.Bounds().Dy())},
			},
			Collided: false,
			AI:       ai,
		}
		player.Animations = make(map[gfx.AnimationType]*gfx.Animation)

		fireAnimation, err := gfx.AnimationFromGIF(bytes.NewReader(fireGif))
		if err != nil {
			log.Error().Err(err).Msg("could not load fire animation")
			return nil
		}
		fireAnimation.AnimationSpeed = 5
		player.Animations[gfx.Fire] = fireAnimation

		g.players = append(g.players, player)
	}
	return g
}

func (g *Game) updatePlayer(p *entities.Player) {
	output := p.AI.Compute(entities.AIInput{
		Position:     p.Position,
		Speed:        p.Speed,
		CurrentSpeed: p.CurrentSpeed,
		Orientation:  p.Orientation,
		Collided:     p.Collided,
		CannonReady:  p.CannonCooldown <= 0,
	})
	//jsonOutput, _ := json.Marshal(output)
	//log.Debug().RawJSON("output", jsonOutput).Msg("bot Compute() result")
	p.Speed = output.Speed
	p.Orientation = p.Orientation + output.OrientationChange
	if p.CannonCooldown > 0 {
		p.CannonCooldown--
	} else {
		if output.Shoot {
			p.CannonCooldown = CannonCooldown
			newShell, err := entities.NewShell()
			newShell.SetOrientation(p.Orientation)
			newShell.SetPosition(p.Position)
			newShell.SetSpeed(30)
			if err != nil {
				log.Error().Err(err).Msg("failed adding shell")
			}
			g.shells = append(g.shells, newShell)
			log.Debug().Str("tank", p.Name).Msg("tank fired")
		}
	}

	playerVector := entities.Vector{
		X: float64(p.Speed) * math.Cos(p.Orientation*math.Pi/180),
		Y: float64(p.Speed) * math.Sin(p.Orientation*math.Pi/180),
	}

	//oldPos := p.Position

	collisionPoint := entities.Vector{X: p.Position.X + playerVector.X, Y: p.Position.Y + playerVector.Y}
	p.Collided = false
	if int(collisionPoint.X) < 0+p.Hitbox.Dx() || int(collisionPoint.X) > g.arenaMap.PixelWidth-p.Hitbox.Dx() {
		p.Collided = true
		p.Health -= ColisionDamage
		if p.Health <= 0 {
			p.State = entities.Dead
		}
	} else {
		if p.State == entities.Alive {
			p.Position.X += playerVector.X
		}
	}
	if int(collisionPoint.Y) < 0+p.Hitbox.Dy() || int(collisionPoint.Y) > g.arenaMap.PixelHeight-p.Hitbox.Dy() {
		p.Collided = true
		p.Health -= ColisionDamage
		if p.Health <= 0 {
			p.State = entities.Dead
		}
	} else {
		if p.State == entities.Alive {
			p.Position.Y += playerVector.Y
		}
	}
}

func remove(s []*entities.Shell, i int) []*entities.Shell {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (g *Game) Update() error {
	for _, p := range g.players {
		if p.State == entities.Alive {
			g.updatePlayer(p)
		}
	}

	for i, s := range g.shells {
		shellVector := entities.Vector{
			X: float64(s.Speed()) * math.Cos(s.Orientation()*math.Pi/180),
			Y: float64(s.Speed()) * math.Sin(s.Orientation()*math.Pi/180),
		}

		position := s.Position()
		collisionPoint := entities.Vector{X: position.X + shellVector.X, Y: position.Y + shellVector.Y}
		if collisionPoint.X < 0+s.CollisionBox().Max.X || collisionPoint.X > float64(g.arenaMap.PixelWidth)-s.CollisionBox().Max.X {
			g.shells = remove(g.shells, i)
			continue
		} else {
			position.X += shellVector.X
		}
		if collisionPoint.Y < 0+s.CollisionBox().Max.Y || collisionPoint.Y > float64(g.arenaMap.PixelHeight)-s.CollisionBox().Max.Y {
			g.shells = remove(g.shells, i)
			continue
		} else {
			position.Y += shellVector.Y
		}
		s.SetPosition(position)
	}
	return nil
}

func RotateImgOpts(img *ebiten.Image, op ebiten.DrawImageOptions, degrees int) ebiten.DrawImageOptions {
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
		g.screenBuffer.DrawImage(layer.Render(g.arenaMap, g.scalingFactor, false), &ebiten.DrawImageOptions{})
	}

	// collisionOp := &ebiten.DrawImageOptions{}
	// collisionOp.ColorM.Scale(1, 0, 0, .75)
	// err := g.screenBuffer.DrawImage(g.arenaMap.GetObjectGroupByName("collisionmap").DebugRender(g.arenaMap, g.scalingFactor), collisionOp)
	// if err != nil {
	// 	log.Debug().Err(err).Msg("rendering collisionmap failed")
	// }

	// ======== Draw Player =========
	for _, p := range g.players {
		playerOp := ebiten.DrawImageOptions{}
		playerOp = RotateImgOpts(p.Sprite, playerOp, int(p.Orientation))
		playerOp.ColorM.Scale(p.Color.R, p.Color.G, p.Color.B, p.Color.Alpha)

		// // entities.Position is absolute on the Map, the coordinates here need to be relative to the camera 0/0
		// // Here is an edge case when the Camera is bigger than the map, stuff breaks
		//playerProjectedX := int(float64(entities.Position.X-g.arenaMap.CameraPosition.X)*g.scalingFactor + float64(g.arenaMap.CameraBounds.Max.X)/2)
		//playerProjectedY := int(float64(entities.Position.Y-g.arenaMap.CameraPosition.Y)*g.scalingFactor + float64(g.arenaMap.CameraBounds.Max.Y)/2)

		// // to move the image
		playerOp.GeoM.Translate(p.Position.X-float64(p.Sprite.Bounds().Dx()/2), p.Position.Y-float64(p.Sprite.Bounds().Dy()/2))

		g.screenBuffer.DrawImage(p.Sprite, &playerOp)

		if p.State == entities.Dead {
			fireOp := ebiten.DrawImageOptions{}
			fireOp.GeoM.Translate(p.Position.X-float64(p.Animations[gfx.Fire].Width/2), p.Position.Y-float64(p.Animations[gfx.Fire].Height/2))
			g.screenBuffer.DrawImage(p.Animations[gfx.Fire].GetFrame(), &fireOp)
		}
		//log.Debug().Str("name", entities.Name).Str("pos", entities.Position.String()).Float64("orientation", entities.Orientation).Msg("draw player")
	}

	// ======== Draw Shells =========
	for _, s := range g.shells {
		shellOp := ebiten.DrawImageOptions{}
		shellOp = RotateImgOpts(s.Sprite(), shellOp, int(s.Orientation()))

		// // to move the image
		shellOp.GeoM.Translate(s.Position().X-float64(s.Sprite().Bounds().Dx()/2), s.Position().Y-float64(s.Sprite().Bounds().Dy()/2))

		g.screenBuffer.DrawImage(s.Sprite(), &shellOp)
	}

	// ======== Screenbuffer ========

	scaledScreenOp := &ebiten.DrawImageOptions{}
	scaledScreenOp.GeoM.Scale(g.scalingFactor, g.scalingFactor)
	screen.DrawImage(g.screenBuffer, scaledScreenOp)

	// ======== Info ========
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
