package main

import (
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"
)

var Screen vec.Vec2 = vec.Vec2{640, 480}
var friction, elasticity, gravity float64 = 0.8, 0.9, 100

type Game struct {
	cam    *kamera.Camera
	space  *cm.Space
	drawer *ebitencm.Drawer
	ball   *cm.Body
}

func (g *Game) Update() error {
	g.space.Step(1 / 60.0)

	// Handling dragging
	g.drawer.HandleMouseEvent(g.space)

	pos := g.ball.Position()
	g.cam.LookAt(pos.X, pos.Y)

	// drawer offset
	g.drawer.CameraOffset.X, g.drawer.CameraOffset.Y = g.cam.TopLeft()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	cm.DrawSpace(g.space, g.drawer.WithScreen(screen))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(Screen.X), int(Screen.Y)
}

func addBall(space *cm.Space, pos vec.Vec2, radius float64) *cm.Body {
	mass := radius / space.Gravity.Y
	body := space.AddBody(cm.NewBody(mass, cm.MomentForCircle(mass, 0, radius, vec.Vec2{})))
	body.SetPosition(pos)
	shape := space.AddShape(cm.NewCircle(body, radius, vec.Vec2{}))
	shape.SetElasticity(elasticity)
	shape.SetFriction(friction)
	return body
}

func main() {
	// Initialising Chipmunk
	space := cm.NewSpace()
	space.SleepTimeThreshold = 0.5
	space.SetGravity(vec.Vec2{X: 0, Y: gravity})

	// walls
	walls := []vec.Vec2{
		{0, 0}, {Screen.X, 0},
		{Screen.X, 0}, {Screen.X, Screen.Y},
		{Screen.X, Screen.Y}, {0, Screen.Y},
		{0, Screen.Y}, {0, 0},
	}
	for i := 0; i < len(walls)-1; i += 2 {
		shape := space.AddShape(cm.NewSegment(space.StaticBody, walls[i], walls[i+1], 10))
		shape.SetElasticity(elasticity)
		shape.SetFriction(friction)
	}

	ballPos := Screen.Scale(0.5)
	game := &Game{
		cam: kamera.NewCamera(ballPos.X, ballPos.Y, Screen.X, Screen.Y),
	}
	game.cam.Lerp = true
	// Add ball
	game.ball = addBall(space, ballPos, 50)

	for range 10 {
		addBall(space, Screen.Scale(rand.Float64()), 30)
	}
	// Initialising Ebitengine/v2
	game.space = space
	game.drawer = ebitencm.NewDrawer()
	game.drawer.OptStroke.AntiAlias = true
	game.drawer.OptFill.AntiAlias = true
	// game.drawer.FillDisabled = true
	// game.drawer.StrokeDisabled = true

	ebiten.SetWindowSize(int(Screen.X), int(Screen.Y))
	ebiten.SetWindowTitle("ebiten-chipmunk - ball")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
