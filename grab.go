package ebitencm

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

var GrabableMaskBit uint = 1 << 31
var grabFilter cm.ShapeFilter = cm.ShapeFilter{
	Group:      cm.NoGroup,
	Categories: GrabableMaskBit,
	Mask:       GrabableMaskBit,
}

type mouseEventHandler struct {
	mouseJoint *cm.Constraint
	mouseBody  *cm.Body
	touchIDs   []ebiten.TouchID
}

func (h *mouseEventHandler) handleMouseEvent(drawer *Drawer, space *cm.Space) {
	if h.mouseBody == nil {
		h.mouseBody = cm.NewKinematicBody()
	}

	var x, y int

	// touch position
	for _, id := range h.touchIDs {
		x, y = ebiten.TouchPosition(id)
		if x == 0 && y == 0 || inpututil.IsTouchJustReleased(id) {
			h.onMouseUp(space)
			h.touchIDs = []ebiten.TouchID{}
			break
		}
	}
	isJuestTouched := false
	touchIDs := inpututil.AppendJustPressedTouchIDs(h.touchIDs[:0])
	for _, id := range touchIDs {
		isJuestTouched = true
		h.touchIDs = []ebiten.TouchID{id}
		x, y = ebiten.TouchPosition(id)
		break
	}
	// mouse position
	if len(h.touchIDs) == 0 {
		x, y = ebiten.CursorPosition()
	}

	// ! GeoM uygulanacak
	cursor := vec.Vec2{X: float64(x), Y: float64(y)}

	// cursor.X, cursor.Y = drawer.GeoM.Apply(cursor.X, cursor.Y)
	cursor = ScreenToWorld(cursor, *drawer.GeoM)

	// offX, offY := drawer.Camera.TopLeft()
	// cursorPosition = cursorPosition.Add(vec.Vec2{offX, offY})

	if isJuestTouched {
		h.mouseBody.SetVelocityVector(vec.Vec2{})
		h.mouseBody.SetPosition(cursor)
	} else {
		newPoint := h.mouseBody.Position().Lerp(cursor, 0.25)
		h.mouseBody.SetVelocityVector(newPoint.Sub(h.mouseBody.Position()).Scale(60.0))
		h.mouseBody.SetPosition(newPoint)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || isJuestTouched {
		h.onMouseDown(space, cursor)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		h.onMouseUp(space)
	}
}

func (h *mouseEventHandler) onMouseDown(space *cm.Space, cursorPosition vec.Vec2) {
	// give the mouse click a little radius to make it easier to click small shapes.
	radius := 5.0

	info := space.PointQueryNearest(cursorPosition, radius, grabFilter)

	if info.Shape != nil && info.Shape.Body().Mass() < cm.Infinity {
		var nearest vec.Vec2
		if info.Distance > 0 {
			nearest = info.Point
		} else {
			nearest = cursorPosition
		}

		body := info.Shape.Body()
		h.mouseJoint = cm.NewPivotJoint2(h.mouseBody, body, vec.Vec2{}, body.WorldToLocal(nearest))
		h.mouseJoint.SetMaxForce(50000)
		h.mouseJoint.SetErrorBias(math.Pow(1.0-0.15, 60.0))
		space.AddConstraint(h.mouseJoint)
	}
}

func (h *mouseEventHandler) onMouseUp(space *cm.Space) {
	if h.mouseJoint == nil {
		return
	}
	space.RemoveConstraint(h.mouseJoint)
	h.mouseJoint = nil
}

// ScreenToWorld converts screen-space coordinates to world-space
func ScreenToWorld(screenPoint vec.Vec2, g ebiten.GeoM) vec.Vec2 {
	if g.IsInvertible() {
		g.Invert()
		worldX, worldY := g.Apply(screenPoint.X, screenPoint.Y)
		return vec.Vec2{worldX, worldY}
	} else {
		// When scaling it can happened that matrix is not invertable
		return vec.Vec2{math.NaN(), math.NaN()}
	}
}
