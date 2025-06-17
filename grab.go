package ebitencm

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/cm"
	"github.com/setanarut/v"
)

const GrabableMaskBit uint = 1 << 31

var GrabFilter cm.ShapeFilter = cm.ShapeFilter{
	Group:      cm.NoGroup,
	Categories: GrabableMaskBit,
	Mask:       cm.AllCategories,
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

	cursor := v.Vec{X: float64(x), Y: float64(y)}
	cursor = ScreenToWorld(cursor, *drawer.GeoM)

	if isJuestTouched {
		h.mouseBody.SetVelocityVector(v.Vec{})
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

func (h *mouseEventHandler) onMouseDown(space *cm.Space, cursorPosition v.Vec) {
	// give the mouse click a little radius to make it easier to click small shapes.
	radius := 5.0

	info := space.PointQueryNearest(cursorPosition, radius, GrabFilter)

	if info.Shape != nil && info.Shape.Body.Mass() < math.MaxFloat64 {
		var nearest v.Vec
		if info.Distance > 0 {
			nearest = info.Point
		} else {
			nearest = cursorPosition
		}

		body := info.Shape.Body
		h.mouseJoint = cm.NewPivotJoint2(h.mouseBody, body, v.Vec{}, body.WorldToLocal(nearest))
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
