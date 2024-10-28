package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/cm"
	"github.com/setanarut/ebitencm"
	"github.com/setanarut/vec"
)

// Sabitler
const (
	MIN_SPEED           = 4.453125
	MAX_SPEED           = 153.75
	MAX_WALK_SPEED      = 93.75
	MAX_FALL_SPEED      = 270.0
	MAX_FALL_SPEED_CAP  = 240.0
	MIN_SLOW_DOWN_SPEED = 33.75
	WALK_ACCELERATION   = 133.59375
	RUN_ACCELERATION    = 200.390625
	WALK_FRICTION       = 182.8125
	SKID_FRICTION       = 365.625
	STOMP_SPEED         = 240.0
	STOMP_SPEED_CAP     = -60.0
	COOLDOWN_TIME_SEC   = 3.0
)

var delta = 1 / 60.0

// Global slice tanımlamaları
var JUMP_SPEED = []float64{-240.0, -240.0, -300.0}
var LONG_JUMP_GRAVITY = []float64{450.0, 421.875, 562.5}
var GRAVITY = []float64{1575.0, 1350.0, 2025.0}
var SPEED_THRESHOLDS = []float64{60, 138.75}

// Input
var is_facing_left = false
var is_running = false

var Is_jumping = false
var is_falling = false
var is_skiding = false
var is_crouching = false

// var _old_velocity = vec.Vec2{}

var input_axis = vec.Vec2{}
var Speed_scale = 0.0

var min_speed = MIN_SPEED
var max_speed = MAX_WALK_SPEED
var acceleration = WALK_ACCELERATION

var speed_threshold int = 0
var Screen vec.Vec2 = vec.Vec2{640, 480}

var playerBody *cm.Body
var playerShape *cm.Shape

// var groundNormal vec.Vec2
var is_on_floor bool

type Game struct {
	space  *cm.Space
	drawer *ebitencm.Drawer
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	cm.DrawSpace(g.space, g.drawer.WithScreen(screen))
}

func (g *Game) Update() error {
	g.space.Step(1 / 60.0)
	is_on_floor = on_floor()
	process_input()
	g.drawer.HandleMouseEvent(g.space)
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(Screen.X), int(Screen.Y)
}

func main() {
	game := &Game{}
	space := cm.NewSpace()
	space.SleepTimeThreshold = 0.5
	space.SetGravity(vec.Vec2{X: 0, Y: 10})

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

	playerBody = cm.NewBody(0.0001, math.MaxFloat64)
	playerShape = cm.NewBoxShapeWithBody2(playerBody, cm.BB{-6, -8, 6, 8}, 0)
	playerBody.SetPosition(vec.Vec2{100, 200})
	playerBody.SetVelocityUpdateFunc(VelFunc)
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

func VelFunc(body *cm.Body, gravity vec.Vec2, damping, dt float64) {
	velocity := playerBody.Velocity()

	// process_jump()
	if is_on_floor {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			Is_jumping = true

			var speed = math.Abs(velocity.X)
			speed_threshold = len(SPEED_THRESHOLDS)

			for i := 0; i < len(SPEED_THRESHOLDS); i++ {
				if speed < SPEED_THRESHOLDS[i] {
					speed_threshold = i
					break
				}
			}
			velocity.Y = JUMP_SPEED[speed_threshold]

		}
	} else {
		var gravity = GRAVITY[speed_threshold]
		if ebiten.IsKeyPressed(ebiten.KeySpace) && !is_falling {
			gravity = LONG_JUMP_GRAVITY[speed_threshold]
		}
		velocity.Y = velocity.Y + gravity*delta
		if velocity.Y > MAX_FALL_SPEED {
			velocity.Y = MAX_FALL_SPEED_CAP
		}
	}

	if velocity.Y > 0 {
		Is_jumping = false
		is_falling = true
	} else if is_on_floor {
		is_falling = false
	}

	// process_walk()
	if input_axis.X != 0 {
		if is_on_floor {
			if velocity.X != 0 {
				is_facing_left = input_axis.X < 0.0
				is_skiding = velocity.X < 0.0 != is_facing_left
			}
			if is_skiding {
				min_speed = MIN_SLOW_DOWN_SPEED
				max_speed = MAX_WALK_SPEED
				acceleration = SKID_FRICTION
			} else if is_running {
				min_speed = MIN_SPEED
				max_speed = MAX_SPEED
				acceleration = RUN_ACCELERATION
			} else {
				min_speed = MIN_SPEED
				max_speed = MAX_WALK_SPEED
				acceleration = WALK_ACCELERATION
			}
		} else if is_running && math.Abs(velocity.X) > MAX_WALK_SPEED {
			max_speed = MAX_SPEED
		} else {
			max_speed = MAX_WALK_SPEED
		}
		var target_speed = input_axis.X * max_speed
		velocity.X = moveToward(velocity.X, target_speed, acceleration*delta)
	} else if is_on_floor && velocity.X != 0 {
		if !is_skiding {
			acceleration = WALK_FRICTION
		}
		if input_axis.Y != 0 {
			min_speed = MIN_SLOW_DOWN_SPEED
		} else {
			min_speed = MIN_SPEED
		}
		if math.Abs(velocity.X) < min_speed {
			velocity.X = 0.0
		} else {
			velocity.X = moveToward(velocity.X, 0.0, acceleration*delta)
		}
	}
	if math.Abs(velocity.X) < MIN_SLOW_DOWN_SPEED {
		is_skiding = false
	}

	Speed_scale = math.Abs(velocity.X) / MAX_SPEED

	playerBody.SetVelocityVector(velocity)
}

func on_floor() bool {
	groundNormal := vec.Vec2{}
	playerBody.EachArbiter(func(arb *cm.Arbiter) {
		n := arb.Normal().Neg()
		if n.Y < groundNormal.Y {
			groundNormal = n
		}
	})

	is_on_floor = groundNormal.Y < 0
	return is_on_floor
}

func getAxis() vec.Vec2 {
	axis := vec.Vec2{}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		axis.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		axis.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		axis.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		axis.X += 1
	}
	return axis
}

func moveToward(from, to, delta float64) float64 {
	if math.Abs(to-from) <= delta {
		return to
	}
	if to > from {
		return from + delta
	}
	return from - delta
}

func process_input() {
	input_axis = getAxis()
	if is_on_floor {
		is_running = ebiten.IsKeyPressed(ebiten.KeyShift)
		is_crouching = ebiten.IsKeyPressed(ebiten.KeyDown)
		if is_crouching && input_axis.X != 0 {
			is_crouching = false
			input_axis.X = 0.0
		}
	}
}
