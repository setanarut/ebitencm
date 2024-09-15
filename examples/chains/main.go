package main

// This is based on "jakecoffman/cp-examples/chains".

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

const (
	screenWidth  = 640
	screenHeight = 480
	hwidth       = screenWidth / 2
	hheight      = screenHeight / 2

	ChainCount = 8
	LinkCount  = 10
)

type Game struct {
	count  int
	space  *cm.Space
	drawer *ebitencm.Drawer
}

func (g *Game) Update() error {
	g.drawer.HandleMouseEvent(g.space)

	g.space.Step(1 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.drawer.Screen = screen
	cm.DrawSpace(g.space, g.drawer)

	msg := fmt.Sprintf(
		"FPS: %0.2f",
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
	space.Iterations = 30
	space.SetGravity(vec.Vec2{X: 0, Y: -100})
	space.SleepTimeThreshold = 0.5

	walls := []vec.Vec2{
		{X: -320, Y: -240}, {X: -320, Y: 240},
		{X: 320, Y: -240}, {X: 320, Y: 240},
		{X: -320, Y: -240}, {X: 320, Y: -240},
		{X: -320, Y: 240}, {X: 320, Y: 240},
	}
	for i := 0; i < len(walls)-1; i += 2 {
		shape := space.AddShape(cm.NewSegment(space.StaticBody, walls[i], walls[i+1], 0))
		shape.SetElasticity(1)
		shape.SetFriction(1)
	}

	mass := 1.0
	width := 20.0
	height := 30.0

	spacing := width * 0.3

	var i, j float64
	for i = 0; i < ChainCount; i++ {
		var prev *cm.Body

		for j = 0; j < LinkCount; j++ {
			pos := vec.Vec2{X: 40 * (i - (ChainCount-1)/2.0), Y: 240 - (j+0.5)*height - (j+1)*spacing}

			body := space.AddBody(cm.NewBody(mass, cm.MomentForBox(mass, width, height)))
			body.SetPosition(pos)

			shape := space.AddShape(cm.NewSegment(body, vec.Vec2{X: 0, Y: (height - width) / 2}, vec.Vec2{X: 0, Y: (width - height) / 2}, width/2))
			shape.SetFriction(0.8)

			breakingForce := 80000.0

			var constraint *cm.Constraint
			if prev == nil {
				constraint = space.AddConstraint(cm.NewSlideJoint(body, space.StaticBody, vec.Vec2{X: 0, Y: height / 2}, vec.Vec2{X: pos.X, Y: 240}, 0, spacing))
			} else {
				constraint = space.AddConstraint(cm.NewSlideJoint(body, prev, vec.Vec2{X: 0, Y: height / 2}, vec.Vec2{X: 0, Y: -height / 2}, 0, spacing))
			}

			constraint.SetMaxForce(breakingForce)
			constraint.PostSolve = BreakableJointPostSolve
			constraint.SetCollideBodies(false)

			prev = body
		}
	}

	radius := 15.0
	body := space.AddBody(cm.NewBody(10, cm.MomentForCircle(10, 0, radius, vec.Vec2{})))
	body.SetPosition(vec.Vec2{X: 0, Y: -240 + radius + 5})
	body.SetVelocity(0, 300)

	shape := space.AddShape(cm.NewCircle(body, radius, vec.Vec2{}))
	shape.SetElasticity(0)
	shape.SetFriction(0.9)

	// Initialising Ebitengine/v2

	game := &Game{}
	game.space = space
	game.drawer = ebitencm.NewDrawer(screenWidth, screenHeight)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ebiten-chipmunk - chains")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func BreakableJointPostStepRemove(space *cm.Space, joint interface{}, _ interface{}) {
	space.RemoveConstraint(joint.(*cm.Constraint))
}

func BreakableJointPostSolve(joint *cm.Constraint, space *cm.Space) {
	dt := space.TimeStep()

	// Convert the impulse to a force by dividing it by the timestep.
	force := joint.Class.GetImpulse() / dt
	maxForce := joint.MaxForce()

	// If the force is almost as big as the joint's max force, break it.
	if force > 0.9*maxForce {
		space.AddPostStepCallback(BreakableJointPostStepRemove, joint, nil)
	}
}
