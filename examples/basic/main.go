package main

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/vec"
)

var Screen = vec.Vec2{640, 480}

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
	return int(Screen.X), int(Screen.Y)
}

func main() {
	// Initialising Chipmunk
	space = cm.NewSpace()
	space.SetGravity(vec.Vec2{0, 100})

	// Walls
	center := Screen.Scale(0.5)
	a := vec.ForAngle(11 * math.Pi / 6).Scale(100).Add(center)
	b := vec.ForAngle(7 * math.Pi / 6).Scale(100).Add(center)
	addWall(space, center, a, 5)
	addWall(space, center, b, 5)

	// Balls
	addBall(space, vec.Vec2{center.X, 0}, 50)
	addBall(space, vec.Vec2{center.X, 0}, 30)
	// Initialising Ebitengine/v2
	game := &Game{}
	drawer = ebitencm.NewDrawer()
	ebiten.SetWindowSize(int(Screen.X), int(Screen.Y))
	ebiten.RunGame(game)
}

func addWall(space *cm.Space, pos1 vec.Vec2, pos2 vec.Vec2, radius float64) {
	shape := space.AddShape(cm.NewSegment(space.StaticBody, pos1, pos2, radius))
	shape.SetElasticity(0.5)
	shape.SetFriction(0.5)
}
func addBall(space *cm.Space, pos vec.Vec2, radius float64) *cm.Body {
	mass := radius * radius / 100.0
	body := space.AddBody(
		cm.NewBody(
			mass,
			cm.MomentForCircle(mass, 0, radius, vec.Vec2{}),
		),
	)
	body.SetPosition(pos)
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
