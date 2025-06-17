package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"
)

var (
	motor     *cm.SimpleMotor
	segRadius                  = 3.0
	space     *cm.Space        = cm.NewSpace()
	drawer    *ebitencm.Drawer = ebitencm.NewDrawer()
	screen    vec.Vec2         = vec.Vec2{640, 480}
	camTarget vec.Vec2         = vec.Vec2{320, 240}
	cam       *kamera.Camera   = kamera.NewCamera(camTarget.X, camTarget.Y, screen.X, screen.Y)
)

type Game struct {
}

func (g *Game) Update() error {
	space.Step(1 / 60.0)
	drawer.GeoM.Reset()
	cam.ApplyCameraTransform(drawer.GeoM)
	drawer.HandleMouseEvent(space)

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		motor.Rate = -5
		motor.SetMaxForce(80000)
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		motor.Rate = 5
		motor.SetMaxForce(80000)
	} else {
		motor.SetMaxForce(0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		cam.Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		cam.ZoomFactor += 10
	}

	if ebiten.IsKeyPressed(ebiten.KeyX) {
		cam.ZoomFactor -= 10
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		cam.Angle += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		cam.Angle -= 0.1
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		camTarget.X += 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		camTarget.X -= 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		camTarget.Y -= 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		camTarget.Y += 10
	}
	cam.LookAt(camTarget.X, camTarget.Y)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawer.DrawSpace(space, screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(screen.X), int(screen.Y)
}

func main() {
	space.Iterations = 20
	space.SetGravity(vec.Vec2{0, 400})
	cam.SmoothType = kamera.Lerp

	offset := 30.0
	chassisMass := 2.0
	a := vec.Vec2{-offset, 0}
	b := vec.Vec2{offset, 0}

	chassis := cm.NewBody(chassisMass, cm.MomentForSegment(chassisMass, a, b, 0))
	space.AddBody(chassis)
	shape := cm.NewSegmentShape(chassis, a, b, segRadius)
	space.AddShape(shape)
	shape.SetShapeFilter(cm.ShapeFilter{1, cm.AllCategories, cm.AllCategories})

	crankMass := 1.0
	crankRadius := 13.0
	crank := cm.NewBody(crankMass, cm.MomentForCircle(crankMass, crankRadius, 0, vec.Vec2{}))
	space.AddBody(crank)
	shape = space.AddShape(cm.NewCircleShape(crank, crankRadius, vec.Vec2{}))
	shape.SetShapeFilter(cm.ShapeFilter{1, cm.AllCategories, cm.AllCategories})
	space.AddConstraint(cm.NewPivotJoint2(chassis, crank, vec.Vec2{}, vec.Vec2{}))
	side := 30.0
	const numLegs = 2
	for i := 0; i < numLegs; i++ {
		makeLeg(space, side, offset, chassis, crank, vec.ForAngle(float64(2*i+0)/numLegs*math.Pi).Scale(crankRadius))
		makeLeg(space, side, -offset, chassis, crank, vec.ForAngle(float64(2*i+1)/numLegs*math.Pi).Scale(crankRadius))
	}

	motor = space.AddConstraint(cm.NewSimpleMotor(chassis, crank, 6)).Class.(*cm.SimpleMotor)

	// move to center of screen
	space.EachBody(func(b *cm.Body) {
		b.SetPosition(b.Position().Add(screen.Scale(0.5)))
	})

	// Walls
	walls := []vec.Vec2{
		{-screen.X, screen.Y}, {screen.X * 2, screen.Y},
	}
	for i := 0; i < len(walls)-1; i += 2 {
		shape := cm.NewSegmentShape(space.StaticBody, walls[i], walls[i+1], 0)
		space.AddShape(shape)
		shape.SetElasticity(0.9)
		shape.SetFriction(0.9)
	}

	ebiten.SetWindowSize(int(screen.X), int(screen.Y))
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func makeLeg(space *cm.Space, side, offset float64, chassis, crank *cm.Body, anchor vec.Vec2) {
	var a, b vec.Vec2
	var shape *cm.Shape

	legMass := 1.0

	// make a leg
	a = vec.Vec2{}
	b = vec.Vec2{0, -side}

	upperLeg := cm.NewBody(legMass, cm.MomentForSegment(legMass, a, b, 0))
	space.AddBody(upperLeg)
	upperLeg.SetPosition(vec.Vec2{offset, 0})

	shape = space.AddShape(cm.NewSegmentShape(upperLeg, a, b, segRadius))
	shape.SetShapeFilter(cm.ShapeFilter{1, cm.AllCategories, cm.AllCategories})

	space.AddConstraint(cm.NewPivotJoint2(chassis, upperLeg, vec.Vec2{offset, 0}, vec.Vec2{}))

	// lower leg
	a = vec.Vec2{}
	b = vec.Vec2{0, -1 * side}
	b = b.NegY()
	lowerLeg := cm.NewBody(legMass, cm.MomentForSegment(legMass, a, b, 0))
	space.AddBody(lowerLeg)
	lowerLeg.SetPosition(vec.Vec2{offset, side})

	shape = cm.NewSegmentShape(lowerLeg, a, b, segRadius)
	space.AddShape(shape)
	shape.SetShapeFilter(cm.ShapeFilter{1, cm.AllCategories, cm.AllCategories})

	shape = space.AddShape(cm.NewCircleShape(lowerLeg, segRadius*2.0, b))
	shape.SetShapeFilter(cm.ShapeFilter{1, cm.AllCategories, cm.AllCategories})
	shape.SetElasticity(0)
	shape.SetFriction(1)

	space.AddConstraint(cm.NewPinJoint(chassis, lowerLeg, vec.Vec2{offset, 0}, vec.Vec2{}))

	space.AddConstraint(cm.NewGearJoint(upperLeg, lowerLeg, 0, 1))

	var constraint *cm.Constraint
	diag := math.Sqrt(side*side + offset*offset)

	constraint = space.AddConstraint(cm.NewPinJoint(crank, upperLeg, anchor, vec.Vec2{0, -side}))
	constraint.Class.(*cm.PinJoint).Dist = diag

	constraint = space.AddConstraint(cm.NewPinJoint(crank, lowerLeg, anchor, vec.Vec2{}))
	constraint.Class.(*cm.PinJoint).Dist = diag
}
