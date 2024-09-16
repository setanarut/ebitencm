package main

// This is based on "jakecoffman/cp-examples/plink".

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/vec"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	pentagonMass   = 0.0
	pentagonMoment = 0.0
)

const numVerts = 5

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

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"FPS: %0.2f",
		ebiten.ActualFPS(),
	))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Initialising Chipmunk

	space := cm.NewSpace()
	space.Iterations = 5
	space.SetGravity(vec.Vec2{X: 0, Y: 500})

	var body *cm.Body
	var shape *cm.Shape

	tris := []vec.Vec2{
		{X: -15, Y: -15},
		{X: 0, Y: 10},
		{X: 15, Y: -15},
	}

	// Triangles
	for i := 0; i < 9; i++ {
		for j := 0; j < 6; j++ {

			stagger := (j % 2) * 40
			offset := vec.Vec2{X: float64(i*80 - 320 + stagger), Y: float64((j * 70) - 170)}
			shape = space.AddShape(cm.NewPolyShape(space.StaticBody, 3, tris, cm.NewTransformTranslate(offset), 0))
			shape.SetElasticity(1)
			shape.SetFriction(1)
		}
	}

	verts := []vec.Vec2{}
	for i := 0; i < numVerts; i++ {
		angle := -2.0 * math.Pi * float64(i) / numVerts
		verts = append(verts, vec.Vec2{X: 10 * math.Cos(angle), Y: 10 * math.Sin(angle)})
	}

	pentagonMass = 0.5
	pentagonMoment = cm.MomentForPoly(1, numVerts, verts, vec.Vec2{}, 0)

	for i := 0; i < 300; i++ {
		body = space.AddBody(cm.NewBody(pentagonMass, pentagonMoment))
		x := rand.Float64()*640 - 320
		body.SetPosition(vec.Vec2{X: x, Y: -300})

		shape = space.AddShape(cm.NewPolyShape(body, numVerts, verts, cm.NewTransformIdentity(), 0))
		shape.SetElasticity(0)
		shape.SetFriction(0.4)
	}

	// Initialising Ebitengine/v2

	game := &Game{}
	game.space = space
	game.drawer = ebitencm.NewDrawer()
	game.drawer.FlipYAxis = true

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ebiten-chipmunk - plink")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
