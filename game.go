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
	ColisionDamage  = 0 // how much health does a player loose on colisions
	CannonCooldown  = 60
	ShellDamage     = 15
	ViewRange       = 2500
	MaxSpeed        = 25.0
	Acceleration    = 0.05
	Friction        = Acceleration * 2
	RespawnWaitTime = 180 // number ticks
	ShellSpeed      = 30.0
	UpdateSpeed     = 1
	StepMode        = false // update frame on key press only
	NextTick        = false
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
	Pressed        map[ebiten.Key]bool
	PressedBefore  map[ebiten.Key]bool
	frameImage     *ebiten.Image
	gameOver       bool
	statsFrame     *ui.Stats
	tabPressed     bool
	respawns       int
	Tick           int
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
			ID:              index,
			Name:            ai.Name(),
			State:           entities.Alive,
			Position:        vector.Vec2{X: float64(spawnPoint.X), Y: float64(spawnPoint.Y)},
			Velocity:        vector.Vec2{0, 0},
			ImpactVelocity:  vector.Vec2{0, 0},
			Mass:            10,
			Health:          100,
			MaxHealth:       100,
			Energy:          100,
			MaxEnergy:       100,
			MaxSpeed:        MaxSpeed,
			Acceleration:    Acceleration,
			Friction:        Friction,
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
	p.UpdateImpactVelocity()

	enemies := make([]*entities.Enemy, 0)
	for _, e := range g.players {
		if e != p {
			distance := physics.Distance(p.Position, e.Position)

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
			Position:         p.Position,
			TargetSpeed:      p.TargetSpeed,
			MaxSpeed:         p.MaxSpeed,
			CurrentSpeed:     p.CurrentSpeed + p.ImpactVelocity.Length(),
			Orientation:      p.Orientation,
			Collided:         p.Collided,
			CollidedWithTank: p.CollidedWithTank,
			Hit:              p.Hit,
			CannonReady:      p.CannonCooldown <= 0,
			Enemy:            enemies,
		})

		p.Orientation += output.OrientationChange
		p.UpdateSpeed(output.Speed)

		if p.CannonCooldown > 0 {
			p.CannonCooldown--
		} else {
			if output.Shoot {
				p.CannonCooldown = CannonCooldown
				newShell := entities.NewShell()
				newShell.Source = p
				newShell.Movement = vector.Vec2{1, 0}.Rotate(p.Orientation).WithLength(ShellSpeed)
				newShell.Orientation = p.Orientation
				newShell.Position = p.Position //.Sum(vector.Vec2{p.CollisionRadius, 0}.Rotate(p.Orientation))
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
				p.TargetSpeed = 0
				p.Velocity = vector.Vec2{0, 0}
				p.ImpactVelocity = vector.Vec2{0, 0}
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

	p.UpdateVelocity()
	velocity := p.Velocity.Sum(p.ImpactVelocity)
	// Collisions
	for index, e := range g.players {
		if e != p {
			circleDistance := physics.DistanceBetweenCircles(vector.Circle{e.Position, e.CollisionRadius}, vector.Circle{p.Position, p.CollisionRadius})

			// ToDo: needs refactor to break this code into physics module and write tests for it
			// collisions: https://www.youtube.com/watch?v=LPzyNOHY3A4&ab_channel=javidx9
			if circleDistance < 0 {
				// check for static collision
				displaceBy := math.Abs(circleDistance) / 2
				vPlayerEnemy := p.Position.ToPoint(e.Position)

				// ToDo: Multiple collisions happen right after one another which causes hughe spikes in ImpactVelocity
				// ToDo: This can move a tank out of the level boundaries
				vDisplace := vPlayerEnemy.Unit().ScalarProduct(displaceBy)
				p.Position.X += vDisplace.X
				p.Position.Y += vDisplace.Y
				g.players[index].Position.X -= vDisplace.X
				g.players[index].Position.Y -= vDisplace.Y

				p.CollidedWithTank = true
				g.players[index].CollidedWithTank = true

				// vector between center points
				vecPE := p.Position.ToPoint(e.Position)

				// normal vector between balls
				normal := vecPE.Unit()

				// perpendicular vector
				tangent := normal.Perpendicular()

				// Dot Product Tangent
				eVelocity := e.Velocity.Sum(e.ImpactVelocity)
				dpTanP := velocity.DotProduct(tangent)
				dpTanE := eVelocity.DotProduct(tangent)

				// Dot Product Normal
				dpNormP := velocity.DotProduct(normal)
				dpNormE := eVelocity.DotProduct(normal)

				// Conservation of momentum in D
				mP := (dpNormP*(p.Mass-e.Mass) + 2.0*e.Mass*dpNormE) / (p.Mass + e.Mass)
				mE := (dpNormE*(e.Mass-p.Mass) + 2.0*p.Mass*dpNormP) / (p.Mass + e.Mass)

				// Update impact velocity
				p.ImpactVelocity.X += (tangent.X*dpTanP + normal.X*mP) * g.scalingFactor // ToDo: Tweak this magic number a bit more
				p.ImpactVelocity.Y += (tangent.Y*dpTanP + normal.Y*mP) * g.scalingFactor
				g.players[index].ImpactVelocity.X -= (tangent.X*dpTanE + normal.X*mE) * g.scalingFactor
				g.players[index].ImpactVelocity.Y -= (tangent.Y*dpTanE + normal.Y*mE) * g.scalingFactor

			}
		}
	}

	// check arena bounds
	collisionPoint := vector.Vec2{X: p.Position.X + velocity.X, Y: p.Position.Y + velocity.Y}
	p.Collided = false
	if collisionPoint.X < 0.0 || collisionPoint.X > float64(g.arenaMap.PixelWidth) ||
		physics.PointLineDistance(vector.Vec2{0, 0}, vector.Vec2{0, float64(g.arenaMap.PixelHeight)}, collisionPoint) < p.CollisionRadius ||
		physics.PointLineDistance(vector.Vec2{float64(g.arenaMap.PixelWidth), 0}, vector.Vec2{float64(g.arenaMap.PixelWidth), float64(g.arenaMap.PixelHeight)}, collisionPoint) < p.CollisionRadius {
		p.Collided = true
		p.Health -= ColisionDamage
		if p.Health <= 0 && p.State == entities.Alive {
			p.RespawnCooldown = RespawnWaitTime
			p.State = entities.Dead
			log.Info().Str("name", p.Name).Msg("crashed into level boundary")
		}
		// ToDo: p.Position - radius - level.X  =  displacement to left border
		p.Velocity.X = 0
		p.ImpactVelocity.X = 0
	}
	if collisionPoint.Y < 0.0 || collisionPoint.Y > float64(g.arenaMap.PixelHeight) ||
		physics.PointLineDistance(vector.Vec2{0, 0}, vector.Vec2{float64(g.arenaMap.PixelWidth), 0}, collisionPoint) < p.CollisionRadius ||
		physics.PointLineDistance(vector.Vec2{0, float64(g.arenaMap.PixelHeight)}, vector.Vec2{float64(g.arenaMap.PixelWidth), float64(g.arenaMap.PixelHeight)}, collisionPoint) < p.CollisionRadius {
		p.Collided = true
		p.Health -= ColisionDamage
		if p.Health <= 0 && p.State == entities.Alive {
			p.RespawnCooldown = RespawnWaitTime
			p.State = entities.Dead
			log.Info().Str("name", p.Name).Str("axis", "y").Msg("crashed into level boundary")
		}
		p.Velocity.Y = 0
		p.ImpactVelocity.Y = 0
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
				p.ImpactVelocity.Sum(shell.Movement.WithLength(10))
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
	g.Pressed = map[ebiten.Key]bool{}
	g.tabPressed = false
	NextTick = false
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			g.Pressed[k] = true
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
			case ebiten.KeyArrowLeft:
				if UpdateSpeed > 1 {
					UpdateSpeed -= 1
				}
			case ebiten.KeyArrowRight:
				if UpdateSpeed < 10 {
					UpdateSpeed += 1
				}
			case ebiten.KeyS:
				if _, exists := g.PressedBefore[k]; !exists {
					StepMode = !StepMode
					log.Debug().Msg("Toggled single step mode")
				}
			case ebiten.KeyN:
				if _, exists := g.PressedBefore[k]; !exists {
					NextTick = true
				}
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
	g.handleInput()

	if StepMode && !NextTick {
		return nil
	}

	if g.Tick%UpdateSpeed == 0 || UpdateSpeed <= 0 {
		if !g.gameOver {
			// update all player positions
			for _, p := range g.players {
				p.Position.X += p.Velocity.X + p.ImpactVelocity.X
				p.Position.Y += p.Velocity.Y + p.ImpactVelocity.Y

				if p.ID == 0 {
					log.Debug().Str("impactVelocity", p.ImpactVelocity.String()).Msg("")
					//log.Debug().Float64("X", (tangent.X*dpTanP+normal.X*mP)*g.scalingFactor).Float64("Y", (tangent.Y*dpTanP+normal.Y*mP)*g.scalingFactor).Msg("collided")
				}
			}

			for _, p := range g.players {
				g.updatePlayer(p)
			}

			g.updateShells()
		}

		g.isGameOver()

		g.Tick = 1
	} else {
		g.Tick++
	}

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
		//ebitenutil.DrawRect(g.screenBuffer, p.Position.X-p.CollisionRadius, p.Position.Y-p.CollisionRadius, p.CollisionRadius*2, p.CollisionRadius*2, color.Gray{})
	}

	// ======== Draw Shells =========
	for _, s := range g.shells {
		shellOp := ebiten.DrawImageOptions{}
		shellOp = gfx.Rotate(s.Sprite(), shellOp, int(s.Orientation))

		// to move the image
		shellOp.GeoM.Translate(s.Position.X-float64(s.Sprite().Bounds().Dx()/2), s.Position.Y-float64(s.Sprite().Bounds().Dy()/2))

		g.screenBuffer.DrawImage(s.Sprite(), &shellOp)
		//ebitenutil.DrawRect(g.screenBuffer, s.Position.X-s.CollisionRadius, s.Position.Y-s.CollisionRadius, s.CollisionRadius*2, s.CollisionRadius*2, color.Gray{})
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
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Position: %s", g.selectedPlayer.Position), 16, 96)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Speed: %f", g.selectedPlayer.Velocity.Length()), 16, 112)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Velocity: %s", g.selectedPlayer.Velocity), 16, 128)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ImpactVelocity: %s", g.selectedPlayer.ImpactVelocity), 16, 144)
	} else {
		for i, p := range g.players {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("#%d - %s H(%d/%d) S(%.0f/%.0f)", i+1, p.Name, p.Health, p.MaxHealth, math.Round(p.Velocity.Length()), p.MaxSpeed), 16, 48+i*16)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
