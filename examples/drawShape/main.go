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
	circ   *cm.Shape
}

func (g *Game) Update() error {
	// Handling dragging
	g.drawer.HandleMouseEvent(g.space)
	g.space.Step(1 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	cm.DrawShape(g.circ, g.drawer.WithScreen(screen))
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
		s := cm.NewSegmentShapeWithBody(space.StaticBody, walls[i], walls[i+1], 10)
		s.SetElasticity(0.5)
		s.SetFriction(0.5)
	}
	space.AddBodyWithShapes(space.StaticBody)

	game.circ = addBall(space, screenWidth*0.5, screenHeight*0.5, 50)

	// Initialising Ebitengine/v2
	game.space = space
	game.drawer = ebitencm.NewDrawer()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ebiten-chipmunk - ball")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func addBall(space *cm.Space, x, y, radius float64) *cm.Shape {
	mass := radius * radius / 500.0
	b := cm.NewBody(mass, cm.MomentForCircle(mass, 0, radius, vec.Vec2{}))
	cm.NewCircleShapeWithBody(b, radius, vec.Vec2{})
	b.Shapes[0].SetElasticity(0.5)
	b.Shapes[0].SetFriction(0.5)
	b.SetPosition(vec.Vec2{x, y})
	space.AddBodyWithShapes(b)
	return b.Shapes[0]
}
