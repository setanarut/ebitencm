package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/vec"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	space  *cm.Space
	drawer *ebitencm.Drawer
)

type Game struct{}

func (g *Game) Update() error {
	// Handling dragging
	drawer.HandleMouseEvent(space)
	space.Step(1 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	cm.DrawSpace(space, drawer.WithScreen(screen))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Initialising Chipmunk
	space = cm.NewSpace()
	space.SetGravity(vec.Vec2{X: 0, Y: -100})
	addWall(space, vec.Vec2{X: -200, Y: -100}, vec.Vec2{X: -10, Y: -150}, 5)
	addWall(space, vec.Vec2{X: 200, Y: -100}, vec.Vec2{X: 10, Y: -150}, 5)
	addBall(space, -50, 0, 50)
	addBall(space, 50, 200, 20)

	// Initialising Ebitengine/v2
	game := &Game{}
	drawer = ebitencm.NewDrawer(screenWidth, screenHeight)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.RunGame(game)
}

func addWall(space *cm.Space, pos1 vec.Vec2, pos2 vec.Vec2, radius float64) {
	shape := space.AddShape(cm.NewSegment(space.StaticBody, pos1, pos2, radius))
	shape.SetElasticity(0.5)
	shape.SetFriction(0.5)
}
func addBall(space *cm.Space, x, y, radius float64) *cm.Body {
	mass := radius * radius / 100.0
	body := space.AddBody(
		cm.NewBody(
			mass,
			cm.MomentForCircle(mass, 0, radius, vec.Vec2{}),
		),
	)
	body.SetPosition(vec.Vec2{X: x, Y: y})

	shape := space.AddShape(
		cm.NewCircle(
			body,
			radius,
			vec.Vec2{},
		),
	)
	shape.SetElasticity(0.5)
	shape.SetFriction(0.5)
	return body
}
