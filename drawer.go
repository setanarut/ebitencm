package ebitencm

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// 16 bytes
type FColor struct {
	R, G, B, A float32
}

// DrawShape draws shapes with the drawer implementation
func (drw *Drawer) DrawShape(shape *cm.Shape, outline, fill cm.FColor) {
	var shapeColor cm.FColor

	body := shape.Body

	switch shape.Class.(type) {
	case *cm.Circle:
		circle := shape.Class.(*cm.Circle)
		drw.DrawCircle(circle.TransformC(), body.Angle(), circle.Radius(), drw.Theme.Outline, shapeColor)
	case *cm.Segment:
		seg := shape.Class.(*cm.Segment)
		drw.DrawFatSegment(seg.TransformA(), seg.TransformB(), seg.Radius(), drw.Theme.Outline, shapeColor)
	case *cm.PolyShape:
		poly := shape.Class.(*cm.PolyShape)

		count := poly.Count()
		planes := poly.Planes
		verts := make([]vec.Vec2, count)

		for i := 0; i < count; i++ {
			verts[i] = planes[i].V0
		}
		drw.DrawPolygon(count, verts, poly.Radius, drw.Theme.Outline, shapeColor)
	default:
		panic("Unknown shape type")
	}
}

var springVerts = []vec.Vec2{
	{0.00, 0.0},
	{0.20, 0.0},
	{0.25, 3.0},
	{0.30, -6.0},
	{0.35, 6.0},
	{0.40, -6.0},
	{0.45, 6.0},
	{0.50, -6.0},
	{0.55, 6.0},
	{0.60, -6.0},
	{0.65, 6.0},
	{0.70, -3.0},
	{0.75, 6.0},
	{0.80, 0.0},
	{1.00, 0.0},
}

// DrawConstraint draws constraints with the drawer implementation
func (drw *Drawer) DrawConstraint(constraint *cm.Constraint, fillColor cm.FColor) {
	bodyA := constraint.BodyA()
	bodyB := constraint.BodyB()

	switch constraint.Class.(type) {
	case *cm.PinJoint:
		joint := constraint.Class.(*cm.PinJoint)

		a := bodyA.Transform().Apply(joint.AnchorA)
		b := bodyB.Transform().Apply(joint.AnchorB)

		drw.DrawDot(5, a, fillColor)
		drw.DrawDot(5, b, fillColor)
		drw.DrawSegment(a, b, fillColor)
	case *cm.SlideJoint:
		joint := constraint.Class.(*cm.SlideJoint)

		a := bodyA.Transform().Apply(joint.AnchorA)
		b := bodyB.Transform().Apply(joint.AnchorB)

		drw.DrawDot(5, a, fillColor)
		drw.DrawDot(5, b, fillColor)
		drw.DrawSegment(a, b, fillColor)
	case *cm.PivotJoint:
		joint := constraint.Class.(*cm.PivotJoint)

		a := bodyA.Transform().Apply(joint.AnchorA)
		b := bodyB.Transform().Apply(joint.AnchorB)

		drw.DrawDot(5, a, fillColor)
		drw.DrawDot(5, b, fillColor)
	case *cm.GrooveJoint:
		joint := constraint.Class.(*cm.GrooveJoint)

		a := bodyA.Transform().Apply(joint.GrooveA)
		b := bodyA.Transform().Apply(joint.GrooveB)
		c := bodyB.Transform().Apply(joint.AnchorB)

		drw.DrawDot(5, c, fillColor)
		drw.DrawSegment(a, b, fillColor)
	case *cm.DampedSpring:
		spring := constraint.Class.(*cm.DampedSpring)
		a := bodyA.Transform().Apply(spring.AnchorA)
		b := bodyB.Transform().Apply(spring.AnchorB)

		drw.DrawDot(5, a, fillColor)
		drw.DrawDot(5, b, fillColor)

		delta := b.Sub(a)
		cos := delta.X
		sin := delta.Y
		s := 1.0 / delta.Mag()

		r1 := vec.Vec2{cos, -sin * s}
		r2 := vec.Vec2{sin, cos * s}

		verts := []vec.Vec2{}
		for i := 0; i < len(springVerts); i++ {
			v := springVerts[i]
			verts = append(verts, vec.Vec2{v.Dot(r1) + a.X, v.Dot(r2) + a.Y})
		}

		for i := 0; i < len(springVerts)-1; i++ {
			drw.DrawSegment(verts[i], verts[i+1], fillColor)
		}
	// these aren't drawn in Chipmunk, so they aren't drawn here
	case *cm.GearJoint:
	case *cm.SimpleMotor:
	case *cm.DampedRotarySpring:
	case *cm.RotaryLimitJoint:
	case *cm.RatchetJoint:
	default:
		panic(fmt.Sprintf("Implement me: %#v", constraint.Class))
	}
}

var StaticDrawingDisabled bool
var DynamicDrawingDisabled bool
var ConstraintDrawingDisabled bool
var CollisionPointDrawingDisabled bool

// DrawSpace draws all shapes in space with the drawer implementation
func (drw *Drawer) DrawSpace(space *cm.Space, screen *ebiten.Image) {
	drw.Screen = screen

	if !StaticDrawingDisabled {
		space.EachStaticShape(func(shape *cm.Shape) {
			drw.DrawShape(shape, drw.Theme.Outline, drw.Theme.ShapeSleepingFill)
		})
	}

	if !DynamicDrawingDisabled {
		space.EachDynamicShape(func(shape *cm.Shape) {
			var clr cm.FColor
			if shape.Body.IsSleeping() {
				clr = drw.Theme.ShapeSleepingFill
			}
			if shape.Body.IdleTime() > shape.Space.SleepTimeThreshold {
				clr = drw.Theme.ShapeIdleFill
			}
			drw.DrawShape(shape, drw.Theme.Outline, clr)
		})
	}

	if !ConstraintDrawingDisabled {
		space.EachConstraint(func(c *cm.Constraint) {
			drw.DrawConstraint(c, drw.Theme.Constraint)
		})

	}
	// Draw Collision Point
	if !CollisionPointDrawingDisabled {
		for _, arb := range space.Arbiters {
			bodyA, bodyB := arb.Bodies()
			arb.ContactPointSet()
			n := arb.Normal()
			for j := 0; j < arb.Count(); j++ {
				p1 := bodyA.Position().Add(arb.Contacts[j].R1)
				p2 := bodyB.Position().Add(arb.Contacts[j].R2)
				a := p1.Add(n.Scale(-2))
				b := p2.Add(n.Scale(2))
				drw.DrawSegment(a, b, drw.Theme.CollisionPoint)
			}
		}
	}
}
