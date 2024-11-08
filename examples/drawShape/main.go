package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/vec"
)

const (
	screenWidth  = 500
	screenHeight = 500
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
	vector.DrawFilledRect(screen, 100, 100, 100, 100, color.White, false)
	g.drawer.DrawSpace(g.space, screen)
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
		s := cm.NewSegmentShape(space.StaticBody, walls[i], walls[i+1], 10)
		s.SetElasticity(0.5)
		s.SetFriction(0.5)
	}
	space.AddBodyWithShapes(space.StaticBody)

	// game.circ = addBall(space, screenWidth*0.5, screenHeight*0.5, 50)
	game.circ = addBox(space, 150, 150, 100, 100)
	// game.circ = addTriangle(space, 150, 150)
	// game.circ = addBoxStatic(space, screenWidth*0.5, screenHeight*0.5, 0)

	// Initialising Ebitengine/v2
	game.space = space
	game.drawer = ebitencm.NewDrawer()

	// game.drawer.StrokeDisabled = true
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ebiten-chipmunk - ball")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func addBox(space *cm.Space, x, y, w, h float64) *cm.Shape {
	b := cm.NewBody(1, cm.MomentForBox(1, w, h))
	s := cm.NewBoxShape(b, w, h, 50)
	s.SetElasticity(0.5)
	s.SetFriction(0.5)
	b.SetPosition(vec.Vec2{x, y})
	space.AddBodyWithShapes(b)
	return s
}
func addTriangle(space *cm.Space, x, y float64) *cm.Shape {
	verts := []vec.Vec2{{0, -5}, {5, 6}, {4, 7}, {4, 7}, {-4, 7}, {-5, 6}}
	geom := cm.NewTransformTranslate(vec.Vec2{0, 0})
	geom.Scale(6, 6)

	for i, v := range verts {
		verts[i] = geom.ApplyVector(v)
	}

	b := cm.NewBody(1, cm.MomentForPoly(1, len(verts), verts, vec.Vec2{}, 0))
	s := cm.NewPolyShape(b, verts, cm.NewTransformIdentity(), 0)
	s.SetElasticity(0.5)
	s.SetFriction(0.5)
	b.SetPosition(vec.Vec2{x, y})
	space.AddBodyWithShapes(b)
	return s
}

// func addBall(space *cm.Space, x, y, radius float64) *cm.Shape {
// 	mass := radius * radius / 500.0
// 	b := cm.NewBody(mass, cm.MomentForCircle(mass, 0, radius, vec.Vec2{}))
// 	cm.NewCircleShape(b, radius, vec.Vec2{})
// 	b.Shapes[0].SetElasticity(0.5)
// 	b.Shapes[0].SetFriction(0.5)
// 	b.SetPosition(vec.Vec2{x, y})
// 	space.AddBodyWithShapes(b)
// 	return b.Shapes[0]
// }
