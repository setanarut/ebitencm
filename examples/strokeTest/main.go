package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/v"
)

const (
	screenWidth  = 300
	screenHeight = 300
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
	screen.Fill(color.Gray{100})
	// Drawing with Ebitengine/v2
	vector.DrawFilledRect(screen, 100, 100, 100, 100, color.White, false)
	g.drawer.DrawSpace(g.space, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{}
	space := cm.NewSpace()
	g.space = space
	space.SleepTimeThreshold = 0.5
	space.SetGravity(v.Vec{X: 0, Y: 0})

	g.drawer = ebitencm.NewDrawer()
	g.drawer.SetOpacity(0.8)
	// g.drawer.DrawingOptions.AllStrokesDisabled = true
	g.drawer.DrawingOptions.DynamicBodyStrokeWidth = 10.
	g.drawer.DrawingOptions.StaticBodyStrokeWidth = 10.

	// walls
	walls := []v.Vec{
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

	addBox(space, 150, 150, 100, 100, 0)
	// AddBall(space, 150, 150, 50)

	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ebiten-chipmunk - ball")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func addBox(space *cm.Space, x, y, w, h, r float64) *cm.Shape {
	b := cm.NewBody(1, cm.MomentForBox(1, w, h))
	s := cm.NewBoxShape(b, w, h, r)
	s.SetElasticity(0.5)
	s.SetFriction(0.5)
	b.SetPosition(v.Vec{x, y})
	space.AddBodyWithShapes(b)
	return s
}
func AddTriangle(space *cm.Space, x, y float64) *cm.Shape {
	verts := []v.Vec{{0, -5}, {5, 6}, {4, 7}, {4, 7}, {-4, 7}, {-5, 6}}
	geom := cm.NewTransformTranslate(v.Vec{0, 0})
	geom.Scale(6, 6)

	for i, v := range verts {
		verts[i] = geom.ApplyVector(v)
	}

	b := cm.NewBody(1, cm.MomentForPoly(1, len(verts), verts, v.Vec{}, 0))
	s := cm.NewPolyShape(b, verts, cm.NewTransformIdentity(), 0)
	s.SetElasticity(0.5)
	s.SetFriction(0.5)
	b.SetPosition(v.Vec{x, y})
	space.AddBodyWithShapes(b)
	return s
}

func AddBall(space *cm.Space, x, y, radius float64) *cm.Shape {
	mass := radius * radius / 500.0
	b := cm.NewBody(mass, cm.MomentForCircle(mass, 0, radius, v.Vec{}))
	cm.NewCircleShape(b, radius, v.Vec{})
	b.Shapes[0].SetElasticity(0.5)
	b.Shapes[0].SetFriction(0.5)
	b.SetPosition(v.Vec{x, y})
	space.AddBodyWithShapes(b)
	return b.Shapes[0]
}
