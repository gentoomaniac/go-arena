package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"plugin"

	_ "embed"

	"github.com/gentoomaniac/ebitmx"
	"github.com/gentoomaniac/go-arena/entities"
	"github.com/gentoomaniac/go-arena/gfx"
	"github.com/gentoomaniac/go-arena/physics"
	"github.com/gentoomaniac/go-arena/ui"
	"github.com/gentoomaniac/go-arena/vector"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/rs/zerolog/log"
)

var (
	ColisionDamage  = 1 // how much health does a player loose on colisions
	CannonCooldown  = 60
	ShellDamage     = 15
	ViewRange       = 2500
	MaxSpeed        = 25.0
	Acceleration    = 0.1
	RespawnWaitTime = 180 // number ticks
	ShellSpeed      = 30.0
)

func NewGame() *Game {
	return &Game{}
}

type Game struct {
	arenaMap       *ebitmx.TmxMap
	scalingFactor  float64
	screenBuffer   *ebiten.Image
	players        []*entities.Player
	shells         []*entities.Shell
	selectedPlayer *entities.Player
	Pressed        []ebiten.Key
	PressedBefore  []ebiten.Key
	frameImage     *ebiten.Image
	gameOver       bool
	statsFrame     *ui.Stats
	tabPressed     bool
	respawns       int
}

func (g *Game) Init() (err error) {
	log.Debug().Msg("init()")
	g.screenBuffer = ebiten.NewImage(g.arenaMap.PixelWidth, g.arenaMap.PixelHeight)
	g.frameImage, err = gfx.LoadFrameSprite()
	if err != nil {
		return
	}
	g.gameOver = false
	g.statsFrame = ui.NewStats("Stats", g.players)
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

//go:embed gfx/fire_transparent.gif
var fireGif []byte

func (g *Game) WithRespawns(respawns int) *Game {
	g.respawns = respawns
	return g
}

func removeSpawn(s []*ebitmx.Object, i int) []*ebitmx.Object {
	if i >= len(s) || i < 0 {
		return s
	}

	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (g *Game) WithBots(bots []string) *Game {
	var color *gfx.Color
	spawnPoints := g.arenaMap.GetObjectGroupByName("spawn_points").Objects
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
		ai.Init()

		playerSprite, err := gfx.GetPlayerSprite()
		if err != nil {
			return nil
		}

		switch index % 4 {
		case 0:
			color = &gfx.Color{R: 1, G: .7, B: .7, Alpha: 1}
		case 1:
			color = &gfx.Color{R: 1, G: 1, B: .7, Alpha: 1}
		case 2:
			color = &gfx.Color{R: .7, G: 1, B: .7, Alpha: 1}
		case 3:
			color = &gfx.Color{R: .7, G: .7, B: 1, Alpha: 1}
		}

		spawnIndex := rand.Int() % len(spawnPoints)
		spawnPoint := spawnPoints[spawnIndex]
		spawnPoints = removeSpawn(spawnPoints, spawnIndex)

		player := &entities.Player{
			Name:            ai.Name(),
			State:           entities.Alive,
			Position:        vector.Vec2{X: float64(spawnPoint.X), Y: float64(spawnPoint.Y)},
			Health:          100,
			MaxHealth:       100,
			Energy:          100,
			MaxEnergy:       100,
			MaxSpeed:        MaxSpeed,
			Acceleration:    Acceleration,
			Orientation:     float64(rand.Int() % 360),
			Sprite:          playerSprite,
			Color:           color,
			CollisionRadius: (float64(playerSprite.Bounds().Dx()) / 2) * .6,
			Collided:        false,
			AI:              ai,
			MaxRespawns:     g.respawns,
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
	enemies := make([]*entities.Enemy, 0)
	for _, e := range g.players {
		if e != p {
			distance := physics.DistanceBetweenCircles(vector.Circle{p.Position, p.CollisionRadius}, vector.Circle{e.Position, e.CollisionRadius})

			// check for collision and displace if collided
			if distance <= 0 {
				displaceBy := math.Abs(distance) / 2
				displacementVector := vector.Vec2{p.Position.X - e.Position.X, p.Position.Y - e.Position.Y}
				p.Position = p.Position.Sum(displacementVector.Unit().ScalarProduct(-displaceBy))
				e.Position = e.Position.Sum(displacementVector.Unit().ScalarProduct(displaceBy))

				distance = physics.DistanceBetweenCircles(vector.Circle{p.Position, p.CollisionRadius}, vector.Circle{e.Position, e.CollisionRadius})
			}

			// add visible enemies to input data
			if distance <= float64(ViewRange) {
				angle := (math.Atan2(e.Position.Y-p.Position.Y, e.Position.X-p.Position.X) * 180 / math.Pi) - p.Orientation
				enemies = append(enemies, &entities.Enemy{
					Distance: distance,
					Angle:    angle,
					State:    e.State,
				})
			}
		}
	}

	if p.State == entities.Alive {
		output := p.AI.Compute(entities.AIInput{
			Position:     p.Position,
			TargetSpeed:  p.TargetSpeed,
			MaxSpeed:     p.MaxSpeed,
			CurrentSpeed: p.CurrentSpeed,
			Orientation:  p.Orientation,
			Collided:     p.Collided,
			Hit:          p.Hit,
			CannonReady:  p.CannonCooldown <= 0,
			Enemy:        enemies,
		})

		p.UpdateSpeed(output.Speed)

		p.Orientation = p.Orientation + output.OrientationChange
		if p.CannonCooldown > 0 {
			p.CannonCooldown--
		} else {
			if output.Shoot {
				p.CannonCooldown = CannonCooldown
				newShell := entities.NewShell()
				newShell.Source = p
				newShell.Movement = vector.Vec2{
					X: ShellSpeed * math.Cos(p.Orientation*math.Pi/180),
					Y: ShellSpeed * math.Sin(p.Orientation*math.Pi/180),
				}
				newShell.Orientation = p.Orientation
				newShell.Position = p.Position
				newShell.Damage = ShellDamage
				newShell.CollisionRadius = float64(newShell.Sprite().Bounds().Dx()) / 2

				g.shells = append(g.shells, newShell)
			}
		}
	} else if p.State == entities.Dead {
		p.UpdateSpeed(0)
		if p.NumberRespawns < p.MaxRespawns {
			if p.RespawnCooldown > 0 {
				p.RespawnCooldown--
			} else {
				p.CurrentSpeed = 0
				spawnPoints := g.arenaMap.GetObjectGroupByName("spawn_points").Objects
				spawnPoint := spawnPoints[rand.Int()%len(spawnPoints)]
				p.Position.X = float64(spawnPoint.X)
				p.Position.Y = float64(spawnPoint.Y)
				p.State = entities.Alive
				p.Health = p.MaxHealth
				p.NumberRespawns++
			}
		}
	}

	p.Movement = vector.Vec2{
		X: p.CurrentSpeed * math.Cos(p.Orientation*math.Pi/180),
		Y: p.CurrentSpeed * math.Sin(p.Orientation*math.Pi/180),
	}

	// check arena bounds
	collisionPoint := vector.Vec2{X: p.Position.X + p.Movement.X, Y: p.Position.Y + p.Movement.Y}
	p.Collided = false
	if physics.Intersection(p.Position, collisionPoint, vector.Vec2{0, 0}, vector.Vec2{0, float64(g.arenaMap.PixelHeight)}) != nil ||
		physics.Intersection(p.Position, collisionPoint, vector.Vec2{float64(g.arenaMap.PixelWidth), 0}, vector.Vec2{float64(g.arenaMap.PixelWidth), float64(g.arenaMap.PixelHeight)}) != nil {
		p.Collided = true
		p.Health -= ColisionDamage
		if p.Health <= 0 && p.State == entities.Alive {
			p.RespawnCooldown = RespawnWaitTime
			p.State = entities.Dead
			log.Info().Str("name", p.Name).Msg("crashed into level boundary")
		}
		p.Movement.X = 0
	}
	if physics.Intersection(p.Position, collisionPoint, vector.Vec2{0, 0}, vector.Vec2{float64(g.arenaMap.PixelWidth), 0}) != nil ||
		physics.Intersection(p.Position, collisionPoint, vector.Vec2{0, float64(g.arenaMap.PixelHeight)}, vector.Vec2{float64(g.arenaMap.PixelWidth), float64(g.arenaMap.PixelHeight)}) != nil {
		p.Collided = true
		p.Health -= ColisionDamage
		if p.Health <= 0 && p.State == entities.Alive {
			p.RespawnCooldown = RespawnWaitTime
			p.State = entities.Dead
			log.Info().Str("name", p.Name).Str("axis", "y").Msg("crashed into level boundary")
		}
		p.Movement.Y = 0
	}

	// check hit by shell
	p.Hit = false
	for i, shell := range g.shells {
		if shell.Source != p {
			if distance := physics.DistanceBetweenCircles(
				vector.Circle{shell.Position, shell.Source.CollisionRadius},
				vector.Circle{p.Position, p.CollisionRadius}); distance < 0 {

				// ToDo: This makes the shell disappear before it visually hit
				// the shell should get a hit flag and get removed after the next draw
				g.shells = remove(g.shells, i)
				p.Hit = true
				p.Health -= shell.Damage
				if p.Health <= 0 && p.State == entities.Alive {
					p.RespawnCooldown = RespawnWaitTime
					p.State = entities.Dead
					log.Info().Str("target", p.Name).Str("source", shell.Source.Name).Int("max", p.MaxRespawns).Int("spawns", p.NumberRespawns).Msg("killed")
				}
			}
		}
	}

	//ToDo: Object collisions are currently not working
	// mapObjects := g.arenaMap.GetObjectGroupByName("collisionmap").Objects
	// pObject := vector.Rect(
	// 	p.CollisionBox().Min.X+p.Movement.X,
	// 	p.CollisionBox().Min.Y+p.Movement.Y,
	// 	p.CollisionBox().Max.X+p.Movement.X,
	// 	p.CollisionBox().Max.Y+p.Movement.Y,
	// )
	// p.Collided = false
	// for _, object := range mapObjects {
	// 	objectBox := vector.Rect(float64(object.X), float64(object.Y), float64(object.X+object.Width), float64(object.Y+object.Height))
	// 	if checkColisionBox(pObject, objectBox) || checkColisionBox(objectBox, pObject) {
	// 		p.Collided = true
	// 		p.Health -= ColisionDamage
	// 		if p.Health <= 0 {
	// 			p.RespawnCooldown = RespawnWaitTime
	// 			p.State = entities.Dead
	// 			log.Info().Str("name", p.Name).Str("object", object.Name).Int("max", p.MaxRespawns).Int("spawns", p.NumberRespawns).Msgf("crashed into object")
	// 		}
	// 		p.CurrentSpeed = 0.0
	// 		p.Movement.X = 0
	// 		p.Movement.Y = 0
	// 	}
	// }
}

func remove(s []*entities.Shell, i int) []*entities.Shell {
	if i >= len(s) || i < 0 {
		return s
	}

	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (g *Game) handleInput() {
	g.Pressed = nil
	g.tabPressed = false
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			g.Pressed = append(g.Pressed, k)
			switch k {
			case ebiten.Key1:
				g.selectedPlayer = g.players[0]
			case ebiten.Key2:
				g.selectedPlayer = g.players[1]
			case ebiten.Key3:
				g.selectedPlayer = g.players[2]
			case ebiten.Key4:
				g.selectedPlayer = g.players[3]
			case ebiten.KeyEscape:
				g.selectedPlayer = nil
			case ebiten.KeyTab:
				g.tabPressed = true
			}
		}
	}
	g.PressedBefore = g.Pressed

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		pointer := vector.Vec2{X: float64(mx) / g.scalingFactor, Y: float64(my) / g.scalingFactor}
		for _, p := range g.players {
			if physics.DistanceBetweenCircles(vector.Circle{pointer, 1}, vector.Circle{p.Position, p.CollisionRadius}) < 0 {
				g.selectedPlayer = p
				break
			}
		}
	}
}

func (g *Game) updateShells() {
	// calculate shells
	for i, s := range g.shells {
		collisionPoint := vector.Vec2{X: s.Position.X + s.Movement.X, Y: s.Position.Y + s.Movement.Y}
		if collisionPoint.X < 0 || collisionPoint.X > float64(g.arenaMap.PixelWidth) {
			g.shells = remove(g.shells, i)
			continue
		} else {
			s.Position.X += s.Movement.X
		}
		if collisionPoint.Y < 0 || collisionPoint.Y > float64(g.arenaMap.PixelHeight) {
			g.shells = remove(g.shells, i)
			continue
		} else {
			s.Position.Y += s.Movement.Y
		}
	}
}

func (g *Game) isGameOver() {
	var alivePlayers = 0
	for _, p := range g.players {
		if p.State == entities.Alive || p.NumberRespawns < p.MaxRespawns {
			alivePlayers++
		}
	}
	if alivePlayers <= 1 {
		g.gameOver = true
	}
}

func (g *Game) Update() error {
	if !g.gameOver {
		g.handleInput()

		// update all player positions
		for _, p := range g.players {
			p.Position.X += p.Movement.X
			p.Position.Y += p.Movement.Y
		}

		for _, p := range g.players {
			g.updatePlayer(p)
		}

		g.updateShells()
	}

	g.isGameOver()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, layer := range g.arenaMap.Layers {
		g.screenBuffer.DrawImage(layer.Render(g.arenaMap, g.scalingFactor, false), &ebiten.DrawImageOptions{})
	}

	// collisionOp := &ebiten.DrawImageOptions{}
	// collisionOp.ColorM.Scale(1, 0, 0, .75)
	// g.screenBuffer.DrawImage(g.arenaMap.GetObjectGroupByName("collisionmap").DebugRender(g.arenaMap, g.scalingFactor), collisionOp)

	// ======== Draw Player =========
	for _, p := range g.players {
		playerOp := ebiten.DrawImageOptions{}
		playerOp = gfx.Rotate(p.Sprite, playerOp, int(p.Orientation))
		playerOp.ColorM.Scale(p.Color.R, p.Color.G, p.Color.B, p.Color.Alpha)
		playerOp.GeoM.Translate(p.Position.X-float64(p.Sprite.Bounds().Dx()/2), p.Position.Y-float64(p.Sprite.Bounds().Dy()/2))

		g.screenBuffer.DrawImage(p.Sprite, &playerOp)

		if p.State == entities.Dead {
			fireOp := ebiten.DrawImageOptions{}
			fireOp.GeoM.Translate(p.Position.X-float64(p.Animations[gfx.Fire].Width/2), p.Position.Y-float64(p.Animations[gfx.Fire].Height/2))
			g.screenBuffer.DrawImage(p.Animations[gfx.Fire].GetFrame(), &fireOp)
		}

		if p == g.selectedPlayer {
			frameOp := ebiten.DrawImageOptions{}
			frameOp.GeoM.Translate(p.Position.X-float64(g.frameImage.Bounds().Dx())/2, p.Position.Y-float64(g.frameImage.Bounds().Dy())/2)
			g.screenBuffer.DrawImage(g.frameImage, &frameOp)
		}
	}

	// ======== Draw Shells =========
	for _, s := range g.shells {
		shellOp := ebiten.DrawImageOptions{}
		shellOp = gfx.Rotate(s.Sprite(), shellOp, int(s.Orientation))

		// // to move the image
		shellOp.GeoM.Translate(s.Position.X-float64(s.Sprite().Bounds().Dx()/2), s.Position.Y-float64(s.Sprite().Bounds().Dy()/2))

		g.screenBuffer.DrawImage(s.Sprite(), &shellOp)
	}

	// ======== Screenbuffer ========

	scaledScreenOp := &ebiten.DrawImageOptions{}
	scaledScreenOp.GeoM.Scale(g.scalingFactor, g.scalingFactor)
	screen.DrawImage(g.screenBuffer, scaledScreenOp)

	if g.gameOver || g.tabPressed {
		frame := g.statsFrame.Image(true)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(float64(screenWidth/2)-float64(frame.Bounds().Dx()/2)*0.5, float64(screenHeight/2)-float64(frame.Bounds().Dy()/2)*0.5)
		screen.DrawImage(frame, op)
	}

	// ======== Info ========
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()), 16, 16)
	ebitenutil.DebugPrintAt(screen, "----", 16, 32)
	if g.selectedPlayer != nil {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Selected Player: %s", g.selectedPlayer.Name), 16, 64)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Health: %d/%d", g.selectedPlayer.Health, g.selectedPlayer.MaxHealth), 16, 80)
	} else {
		for i, p := range g.players {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("#%d - %s (%d/%d)", i+1, p.Name, p.Health, p.MaxHealth), 16, 48+i*16)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
