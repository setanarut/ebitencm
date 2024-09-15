package main

// This is based on "jakecoffman/cp-examples/theojansen".

import (
	"fmt"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/vec"
)

var (
	motor *cm.SimpleMotor
)

const (
	screenWidth  = 640
	screenHeight = 480
	hwidth       = screenWidth / 2
	hheight      = screenHeight / 2
)

type Game struct {
	space    *cm.Space
	drawer   *ebitencm.Drawer
	touchIDs []ebiten.TouchID
}

func (g *Game) Update() error {

	clickLeft, clickRight := false, false
	mouseX, _ := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if mouseX < screenWidth/2 {
			clickLeft = true
		} else {
			clickRight = true
		}
	}

	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || g.leftTouched() || clickLeft {
		motor.Rate = -5
		motor.SetMaxForce(100000)
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || g.rightTouched() || clickRight {
		motor.Rate = 5
		motor.SetMaxForce(100000)
	} else {
		motor.SetMaxForce(0)
	}

	g.space.Step(1 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	cm.DrawSpace(g.space, g.drawer.WithScreen(screen))

	msg := fmt.Sprintf(
		"FPS: %0.2f\n"+
			"Press left or right arrow key to rotate the motor\n"+
			"Press left or right half of the screen to rotate the motor",
		ebiten.ActualFPS(),
	)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	// Initialising Chipmunk

	space := cm.NewSpace()
	space.Iterations = 20
	space.SetGravity(vec.Vec2{X: 0, Y: -500})

	var shape *cm.Shape
	var a, b vec.Vec2

	walls := []vec.Vec2{
		{X: -320, Y: -240}, {X: -320, Y: 240},
		{X: 320, Y: -240}, {X: 320, Y: 240},
		{X: -320, Y: -240}, {X: 320, Y: -240},
	}

	for i := 0; i < len(walls)-1; i += 2 {
		shape = space.AddShape(cm.NewSegment(space.StaticBody, walls[i], walls[i+1], 0))
		shape.SetElasticity(0.9)
		shape.SetFriction(0.9)
	}

	offset := 30.0
	chassisMass := 2.0
	a = vec.Vec2{X: -offset, Y: 0}
	b = vec.Vec2{X: offset, Y: 0}
	chassis := space.AddBody(cm.NewBody(chassisMass, cm.MomentForSegment(chassisMass, a, b, 0)))

	shape = space.AddShape(cm.NewSegment(chassis, a, b, segRadius))
	shape.SetShapeFilter(cm.NewShapeFilter(1, cm.AllCategories, cm.AllCategories))

	crankMass := 1.0
	crankRadius := 13.0
	crank := space.AddBody(cm.NewBody(crankMass, cm.MomentForCircle(crankMass, crankRadius, 0, vec.Vec2{})))

	shape = space.AddShape(cm.NewCircle(crank, crankRadius, vec.Vec2{}))
	shape.SetShapeFilter(cm.NewShapeFilter(1, cm.AllCategories, cm.AllCategories))

	space.AddConstraint(cm.NewPivotJoint2(chassis, crank, vec.Vec2{}, vec.Vec2{}))

	side := 30.0

	const numLegs = 2
	for i := 0; i < numLegs; i++ {
		makeLeg(space, side, offset, chassis, crank, vec.ForAngle(float64(2*i+0)/numLegs*math.Pi).Scale(crankRadius))
		makeLeg(space, side, -offset, chassis, crank, vec.ForAngle(float64(2*i+1)/numLegs*math.Pi).Scale(crankRadius))
	}

	motor = space.AddConstraint(cm.NewSimpleMotor(chassis, crank, 6)).Class.(*cm.SimpleMotor)

	// Initialising Ebitengine/v2

	game := &Game{}
	game.space = space
	game.drawer = ebitencm.NewDrawer(screenWidth, screenHeight)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ebiten-chipmunk - theojansen")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

const segRadius = 3.0

func makeLeg(space *cm.Space, side, offset float64, chassis, crank *cm.Body, anchor vec.Vec2) {
	var a, b vec.Vec2
	var shape *cm.Shape

	legMass := 1.0

	// make a leg
	a = vec.Vec2{}
	b = vec.Vec2{X: 0, Y: side}
	upperLeg := space.AddBody(cm.NewBody(legMass, cm.MomentForSegment(legMass, a, b, 0)))
	upperLeg.SetPosition(vec.Vec2{X: offset, Y: 0})

	shape = space.AddShape(cm.NewSegment(upperLeg, a, b, segRadius))
	shape.SetShapeFilter(cm.NewShapeFilter(1, cm.AllCategories, cm.AllCategories))

	space.AddConstraint(cm.NewPivotJoint2(chassis, upperLeg, vec.Vec2{X: offset, Y: 0}, vec.Vec2{}))

	// lower leg
	a = vec.Vec2{}
	b = vec.Vec2{X: 0, Y: -1 * side}
	lowerLeg := space.AddBody(cm.NewBody(legMass, cm.MomentForSegment(legMass, a, b, 0)))
	lowerLeg.SetPosition(vec.Vec2{X: offset, Y: -side})

	shape = space.AddShape(cm.NewSegment(lowerLeg, a, b, segRadius))
	shape.SetShapeFilter(cm.NewShapeFilter(1, cm.AllCategories, cm.AllCategories))

	shape = space.AddShape(cm.NewCircle(lowerLeg, segRadius*2.0, b))
	shape.SetShapeFilter(cm.NewShapeFilter(1, cm.AllCategories, cm.AllCategories))
	shape.SetElasticity(0)
	shape.SetFriction(1)

	space.AddConstraint(cm.NewPinJoint(chassis, lowerLeg, vec.Vec2{X: offset, Y: 0}, vec.Vec2{}))

	space.AddConstraint(cm.NewGearJoint(upperLeg, lowerLeg, 0, 1))

	var constraint *cm.Constraint
	diag := math.Sqrt(side*side + offset*offset)

	constraint = space.AddConstraint(cm.NewPinJoint(crank, upperLeg, anchor, vec.Vec2{X: 0, Y: side}))
	constraint.Class.(*cm.PinJoint).Dist = diag

	constraint = space.AddConstraint(cm.NewPinJoint(crank, lowerLeg, anchor, vec.Vec2{}))
	constraint.Class.(*cm.PinJoint).Dist = diag
}

func (g *Game) leftTouched() bool {
	for _, id := range g.touchIDs {
		x, _ := ebiten.TouchPosition(id)
		if x < screenWidth/2 {
			return true
		}
	}
	return false
}

func (g *Game) rightTouched() bool {
	for _, id := range g.touchIDs {
		x, _ := ebiten.TouchPosition(id)
		if x >= screenWidth/2 {
			return true
		}
	}
	return false
}
