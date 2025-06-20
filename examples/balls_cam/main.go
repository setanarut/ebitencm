package main

import (
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/v"
)

var Screen v.Vec = v.Vec{640, 480}
var friction, elasticity, gravity float64 = 0.8, 0.9, 100
var targetOffset = v.Vec{}

type Game struct {
	cam    *kamera.Camera
	space  *cm.Space
	drawer *ebitencm.Drawer
	ball   *cm.Body
}

func (g *Game) Update() error {
	g.space.Step(1 / 60.0)

	// Apply camera transform to drawer
	g.drawer.GeoM.Reset()
	g.cam.ApplyCameraTransform(g.drawer.GeoM)

	// Enable cursor dragging
	g.drawer.HandleMouseEvent(g.space)

	pos := g.ball.Position()

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		targetOffset = v.Vec{}
		g.cam.Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.cam.ZoomFactor += 10
	}

	if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.cam.ZoomFactor -= 10
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.cam.Angle += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.cam.Angle -= 0.1
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		targetOffset.X += 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		targetOffset.X -= 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		targetOffset.Y -= 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		targetOffset.Y += 10
	}

	// g.cam.LookAt(target.X, target.Y)
	pos = pos.Add(targetOffset)
	g.cam.LookAt(pos.X, pos.Y)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	g.drawer.DrawSpace(g.space, screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(Screen.X), int(Screen.Y)
}

func addBall(space *cm.Space, pos v.Vec, radius float64) *cm.Body {
	mass := radius / space.Gravity.Y
	b := cm.NewBody(mass, cm.MomentForCircle(mass, 0, radius, v.Vec{}))
	cm.NewCircleShape(b, radius, v.Vec{})
	b.Shapes[0].SetElasticity(elasticity)
	b.Shapes[0].SetFriction(friction)
	space.AddBodyWithShapes(b)
	b.SetPosition(pos)
	return b
}

func main() {
	// Initialising Chipmunk
	space := cm.NewSpace()
	space.SleepTimeThreshold = 0.5
	space.SetGravity(v.Vec{X: 0, Y: gravity})

	// walls
	walls := []v.Vec{
		{0, 0}, {Screen.X, 0},
		{Screen.X, 0}, {Screen.X, Screen.Y},
		{Screen.X, Screen.Y}, {0, Screen.Y},
		{0, Screen.Y}, {0, 0},
	}
	// sbwall := cm.NewStaticBody()
	for i := 0; i < len(walls)-1; i += 2 {
		s := cm.NewSegmentShape(space.StaticBody, walls[i], walls[i+1], 10)
		s.SetElasticity(elasticity)
		s.SetFriction(friction)
	}
	space.AddBodyWithShapes(space.StaticBody)

	ballPos := Screen.Scale(0.5)
	game := &Game{
		cam: kamera.NewCamera(ballPos.X, ballPos.Y, Screen.X, Screen.Y),
	}

	// Add ball
	game.ball = addBall(space, ballPos, 50)

	for range 10 {
		addBall(space, Screen.Scale(rand.Float64()), 30)
	}

	// Initialising Ebitengine/v2
	game.cam.SmoothType = kamera.Lerp
	game.space = space

	// Init drawer
	game.drawer = ebitencm.NewDrawer()
	game.drawer.DrawTriangleStrokeOpt.AntiAlias = true
	game.drawer.DrawTriagleFillOpt.AntiAlias = true
	// game.drawer.FillDisabled = true
	// game.drawer.StrokeDisabled = true

	ebiten.SetWindowSize(int(Screen.X), int(Screen.Y))
	ebiten.SetWindowTitle("ebitencm camera example")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
