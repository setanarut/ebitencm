package main

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/v"
)

var Screen = v.Vec{640, 480}
var drawer *ebitencm.Drawer = ebitencm.NewDrawer()
var (
	space *cm.Space
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
	drawer.DrawSpace(space, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(Screen.X), int(Screen.Y)
}

func main() {

	// Initialising Chipmunk
	space = cm.NewSpace()
	space.SetGravity(v.Vec{0, 100})

	// Walls
	center := Screen.Scale(0.5)
	a := v.FromAngle(11 * math.Pi / 6).Scale(100).Add(center)
	b := v.FromAngle(7 * math.Pi / 6).Scale(100).Add(center)
	addWall(space, center, a, 5)
	addWall(space, center, b, 5)

	// Balls
	addBall(space, v.Vec{center.X, 0}, 50)
	addBall(space, v.Vec{center.X, 0}, 30)

	// Initialising Ebitengine
	game := &Game{}
	ebiten.SetWindowSize(int(Screen.X), int(Screen.Y))
	ebiten.RunGame(game)
}

func addWall(space *cm.Space, pos1 v.Vec, pos2 v.Vec, radius float64) {
	sb := cm.NewStaticBody()
	cm.NewSegmentShape(sb, pos1, pos2, radius)
	shape := sb.Shapes[0]
	shape.SetElasticity(0.5)
	shape.SetFriction(0.5)
	space.AddBodyWithShapes(sb)
}
func addBall(space *cm.Space, pos v.Vec, radius float64) *cm.Body {
	mass := radius * radius / 100.0
	body := cm.NewBody(mass, cm.MomentForCircle(mass, 0, radius, v.Vec{}))
	cm.NewCircleShape(body, radius, v.Vec{})
	body.Shapes[0].SetElasticity(0.5)
	body.Shapes[0].SetFriction(0.5)
	body.SetPosition(pos)
	space.AddBodyWithShapes(body)
	return body
}
