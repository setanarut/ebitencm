package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/vec"
)

var Screen vec.Vec2 = vec.Vec2{640, 480}

const (
	PlayerVelocity        = 500.0
	PlayerGroundAccelTime = 0.1
	PlayerAirAccelTime    = 0.25
	JumpHeight            = 50.0
	JumpBoostHeight       = 200.0
	FallVelocity          = 1000.0
	Gravity               = 2000.0

	PlayerGroundAccel = PlayerVelocity / PlayerGroundAccelTime
	PlayerAirAccel    = PlayerVelocity / PlayerAirAccelTime
)

var playerBody *cm.Body
var playerShape *cm.Shape

var remainingBoost float64
var grounded, lastJumpState bool

type Game struct {
	space  *cm.Space
	drawer *ebitencm.Drawer
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	cm.DrawSpace(g.space, g.drawer.WithScreen(screen))
}
func (g *Game) Update() error {

	// Handling dragging
	g.drawer.HandleMouseEvent(g.space)
	g.space.Step(1 / 60.0)

	jumpState := ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp)

	// If the jump key was just pressed this frame, jump!
	if jumpState && !lastJumpState && grounded {
		jumpV := math.Sqrt(2.0 * JumpHeight * Gravity)
		playerBody.SetVelocityVector(playerBody.Velocity().Add(vec.Vec2{0, -jumpV}))

		remainingBoost = JumpBoostHeight / jumpV
	}

	remainingBoost -= 1. / 60.
	lastJumpState = jumpState

	return nil

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(Screen.X), int(Screen.Y)
}

func main() {
	game := &Game{}
	space := cm.NewSpace()
	space.SleepTimeThreshold = 0.5
	space.SetGravity(vec.Vec2{X: 0, Y: Gravity})

	// walls
	walls := []vec.Vec2{
		{0, 0}, {Screen.X, 0},
		{Screen.X, 0}, {Screen.X, Screen.Y},
		{Screen.X, Screen.Y}, {0, Screen.Y},
		{0, Screen.Y}, {0, 0},
	}
	for i := 0; i < len(walls)-1; i += 2 {
		s := cm.NewSegmentShapeWithBody(space.StaticBody, walls[i], walls[i+1], 40)
		s.SetElasticity(0.5)
		s.SetFriction(0.5)
	}
	space.AddBodyWithShapes(space.StaticBody)

	playerBody = cm.NewBody(1, math.MaxFloat64)
	playerShape = cm.NewBoxShapeWithBody2(playerBody, cm.BB{-15, -27.5, 15, 27.5}, 10)
	playerBody.SetPosition(vec.Vec2{100, 200})
	playerBody.SetVelocityUpdateFunc(playerUpdateVelocity)
	playerShape.SetElasticity(0)
	playerShape.SetFriction(0)
	playerShape.SetCollisionType(1)
	space.AddBodyWithShapes(playerBody)

	// Initialising Ebitengine/v2
	game.space = space
	game.drawer = ebitencm.NewDrawer()

	game.drawer.OptStroke.AntiAlias = true
	game.drawer.OptFill.AntiAlias = true

	ebiten.SetWindowSize(int(Screen.X), int(Screen.Y))
	ebiten.SetWindowTitle("ebiten-chipmunk - ball")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func lerpConst(f1, f2, d float64) float64 {
	return f1 + clamp(f2-f1, -d, d)
}

func clamp(f, min, max float64) float64 {
	if f > min {
		return math.Min(f, max)
	} else {
		return math.Min(min, max)
	}
}

func playerUpdateVelocity(body *cm.Body, gravity vec.Vec2, damping, dt float64) {
	jumpState := ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp)

	// Grab the grounding normal from last frame
	groundNormal := vec.Vec2{}
	playerBody.EachArbiter(func(arb *cm.Arbiter) {
		n := arb.Normal().Neg()

		if n.Y < groundNormal.Y {
			groundNormal = n
		}
	})

	grounded = groundNormal.Y < 0
	if groundNormal.Y > 0 {
		remainingBoost = 0
	}

	// Do a normal-ish update
	boost := jumpState && remainingBoost > 0
	var g vec.Vec2
	if !boost {
		g = gravity
	}
	body.UpdateVelocity(g, damping, dt)

	// Target horizontal speed for air/ground control
	var targetVx float64
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		targetVx -= PlayerVelocity
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		targetVx += PlayerVelocity
	}

	// Update the surface velocity and friction
	// Note that the "feet" move in the opposite direction of the player.
	surfaceV := vec.Vec2{-targetVx, 0}
	playerShape.SurfaceVelocity = surfaceV
	if grounded {
		playerShape.SetFriction(PlayerGroundAccel / Gravity)
	} else {
		playerShape.SetFriction(0)
	}

	// Apply air control if not grounded
	if !grounded {
		v := playerBody.Velocity()
		playerBody.SetVelocity(lerpConst(v.X, targetVx, PlayerAirAccel*dt), v.Y)
	}

	v := body.Velocity()

	body.SetVelocity(v.X, clamp(v.Y, -FallVelocity, math.MaxFloat64))
}
