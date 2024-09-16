package main

// This is based on "jakecoffman/cp-examples/bench".

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/vec"
)

var (
	screenSize = vec.Vec2{640, 640}
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

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"FPS: %0.2f",
		ebiten.ActualFPS(),
	))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(screenSize.X), int(screenSize.Y)
}

var (
	ball *cm.Body
)

func main() {
	// Initialising Chipmunk
	space := cm.NewSpace()
	space.SleepTimeThreshold = 0.2
	space.SetGravity(vec.Vec2{0, 50})

	simpleTerrain(space)
	var r float64 = 6.0
	for i := 0; i < 100; i++ {
		pos := vec.Vec2{(float64(i%10) * r * 2), (float64(i/10) * r * 2)}
		pos = pos.Add(screenSize.Scale(0.5)).Add(vec.Vec2{-50, -50})
		addBall(space, pos.X, pos.Y, r)
	}

	// Initialising Ebitengine/v2
	game := &Game{}
	game.space = space
	game.drawer = ebitencm.NewDrawer()
	// game.drawer.StrokeDisabled = true
	// game.drawer.FillDisabled = true
	ebiten.SetWindowSize(int(screenSize.X), int(screenSize.Y))
	ebiten.SetWindowTitle("ebiten-chipmunk - bench")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func addBall(space *cm.Space, x, y, radius float64) {
	mass := 1.
	body := space.AddBody(cm.NewBody(mass, cm.MomentForCircle(mass, 0, radius, vec.Vec2{})))
	ball = body
	body.SetPosition(vec.Vec2{x, y})
	shape := space.AddShape(cm.NewCircle(body, radius, vec.Vec2{}))
	shape.SetElasticity(0.91)
	shape.SetFriction(0.9)
	body.ApplyImpulseAtLocalPoint(vec.Vec2{0, 200}, vec.Vec2{})
}

func simpleTerrain(space *cm.Space) *cm.Space {
	var simpleTerrainVerts = []vec.Vec2{
		{350.00, 425.07}, {336.00, 436.55}, {272.00, 435.39}, {258.00, 427.63}, {225.28, 420.00}, {202.82, 396.00},
		{191.81, 388.00}, {189.00, 381.89}, {173.00, 380.39}, {162.59, 368.00}, {150.47, 319.00}, {128.00, 311.55},
		{119.14, 286.00}, {126.84, 263.00}, {120.56, 227.00}, {141.14, 178.00}, {137.52, 162.00}, {146.51, 142.00},
		{156.23, 136.00}, {158.00, 118.27}, {170.00, 100.77}, {208.43, 84.00}, {224.00, 69.65}, {249.30, 68.00},
		{257.00, 54.77}, {363.00, 45.94}, {374.15, 54.00}, {386.00, 69.60}, {413.00, 70.73}, {456.00, 84.89},
		{468.09, 99.00}, {467.09, 123.00}, {464.92, 135.00}, {469.00, 141.03}, {497.00, 148.67}, {513.85, 180.00},
		{509.56, 223.00}, {523.51, 247.00}, {523.00, 277.00}, {497.79, 311.00}, {478.67, 348.00}, {467.90, 360.00},
		{456.76, 382.00}, {432.95, 389.00}, {417.00, 411.32}, {373.00, 433.19}, {361.00, 430.02}, {350.00, 425.07},
	}

	// terrain offset
	offset := vec.Vec2{0, 0}

	for i := 0; i < len(simpleTerrainVerts)-1; i++ {
		a := simpleTerrainVerts[i]
		b := simpleTerrainVerts[i+1]
		s := cm.NewSegment(space.StaticBody, a.Add(offset), b.Add(offset), 3)
		s.SetElasticity(0.91)
		s.SetFriction(0.9)
		space.AddShape(s)
	}

	return space
}
