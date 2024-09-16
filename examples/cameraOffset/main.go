package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/vec"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	space  *cm.Space
	drawer *ebitencm.Drawer
}

func (g *Game) Update() error {
	// Handling dragging
	g.drawer.HandleMouseEvent(g.space)
	g.space.Step(1 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	cm.DrawSpace(g.space, g.drawer.WithScreen(screen))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}
	// Initialising Chipmunk
	space := cm.NewSpace()
	space.SleepTimeThreshold = 0.5
	space.SetGravity(vec.Vec2{X: 0, Y: 0})

	// walls
	walls := []vec.Vec2{
		{0, 0}, {screenWidth, 0},
		{screenWidth, 0}, {screenWidth, screenHeight},
		{screenWidth, screenHeight}, {0, screenHeight},
		{0, screenHeight}, {0, 0},
	}
	for i := 0; i < len(walls)-1; i += 2 {
		shape := space.AddShape(cm.NewSegment(space.StaticBody, walls[i], walls[i+1], 10))
		shape.SetElasticity(0.5)
		shape.SetFriction(0.5)
	}

	// balls
	addBall(space, screenWidth*0.5, screenHeight*0.5, 50)

	// Initialising Ebitengine/v2
	game.space = space
	game.drawer = ebitencm.NewDrawer()

	// Camera offset
	game.drawer.CameraOffset = vec.Vec2{100, 40}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ebiten-chipmunk - ball")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func addBall(space *cm.Space, x, y, radius float64) {
	mass := radius * radius / 500.0
	body := space.AddBody(cm.NewBody(mass, cm.MomentForCircle(mass, 0, radius, vec.Vec2{})))
	body.SetPosition(vec.Vec2{X: x, Y: y})
	shape := space.AddShape(cm.NewCircle(body, radius, vec.Vec2{}))
	shape.SetElasticity(0.5)
	shape.SetFriction(0.5)
}
