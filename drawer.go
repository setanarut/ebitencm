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

var springVerts = []vec.Vec2{
	{0.00, 0.0}, {0.20, 0.0}, {0.25, 3.0}, {0.30, -6.0}, {0.35, 6.0},
	{0.40, -6.0}, {0.45, 6.0}, {0.50, -6.0}, {0.55, 6.0}, {0.6, -6.0},
	{0.65, 6.0}, {0.70, -3.0}, {0.75, 6.0}, {0.80, 0.0}, {1.0, 0.0}}

// DrawShape draws shapes with the drawer implementation
func (drw *Drawer) DrawShape(shape *cm.Shape, outline, fill cm.FColor, strokeWidth float32) {
	body := shape.Body

	switch shape.Class.(type) {
	case *cm.Circle:
		circle := shape.Class.(*cm.Circle)
		drw.drawCircle(circle.TransformC(), body.Angle(), circle.Radius(), outline, fill, strokeWidth)
	case *cm.Segment:
		seg := shape.Class.(*cm.Segment)
		drw.drawFatSegment(seg.TransformA(), seg.TransformB(), seg.Radius(), outline, fill, strokeWidth)
	case *cm.PolyShape:
		poly := shape.Class.(*cm.PolyShape)

		count := poly.Count()
		planes := poly.Planes
		verts := make([]vec.Vec2, count)

		for i := 0; i < count; i++ {
			verts[i] = planes[i].V0
		}
		drw.drawPolygon(count, verts, poly.Radius, outline, fill, strokeWidth)
	default:
		panic("Unknown shape type")
	}
}

// DrawConstraint draws constraints with the drawer implementation
func (drw *Drawer) DrawConstraint(constraint *cm.Constraint, strokeWidth float32) {

	bodyA := constraint.BodyA()
	bodyB := constraint.BodyB()

	switch constraint.Class.(type) {

	case *cm.PinJoint:

		joint := constraint.Class.(*cm.PinJoint)
		a := bodyA.Transform().Apply(joint.AnchorA)
		b := bodyB.Transform().Apply(joint.AnchorB)
		drw.drawDot(drw.DrawingOptions.ConstraintsDotRadius, a, drw.Theme.ConstraintPinJointDot)
		drw.drawDot(drw.DrawingOptions.ConstraintsDotRadius, b, drw.Theme.ConstraintPinJointDot)
		drw.drawSegment(a, b, drw.Theme.ConstraintPinJointSegment, strokeWidth)

	case *cm.SlideJoint:

		joint := constraint.Class.(*cm.SlideJoint)
		a := bodyA.Transform().Apply(joint.AnchorA)
		b := bodyB.Transform().Apply(joint.AnchorB)
		drw.drawDot(drw.DrawingOptions.ConstraintsDotRadius, a, drw.Theme.ConstraintSlideJointDot)
		drw.drawDot(drw.DrawingOptions.ConstraintsDotRadius, b, drw.Theme.ConstraintSlideJointDot)
		drw.drawSegment(a, b, drw.Theme.ConstraintSlideJointSegment, strokeWidth)

	case *cm.PivotJoint:

		joint := constraint.Class.(*cm.PivotJoint)
		a := bodyA.Transform().Apply(joint.AnchorA)
		b := bodyB.Transform().Apply(joint.AnchorB)
		drw.drawDot(drw.DrawingOptions.ConstraintsDotRadius, a, drw.Theme.ConstraintPinJointDot)
		drw.drawDot(drw.DrawingOptions.ConstraintsDotRadius, b, drw.Theme.ConstraintPinJointDot)

	case *cm.GrooveJoint:

		joint := constraint.Class.(*cm.GrooveJoint)
		a := bodyA.Transform().Apply(joint.GrooveA)
		b := bodyA.Transform().Apply(joint.GrooveB)
		c := bodyB.Transform().Apply(joint.AnchorB)
		drw.drawDot(drw.DrawingOptions.ConstraintsDotRadius, c, drw.Theme.ConstraintGrooveJointDot)
		drw.drawSegment(a, b, drw.Theme.ConstraintGrooveJointSegment, strokeWidth)

	case *cm.DampedSpring:

		spring := constraint.Class.(*cm.DampedSpring)
		a := bodyA.Transform().Apply(spring.AnchorA)
		b := bodyB.Transform().Apply(spring.AnchorB)
		drw.drawDot(drw.DrawingOptions.ConstraintsDotRadius, a, drw.Theme.ConstraintDampedSpringDot)
		drw.drawDot(drw.DrawingOptions.ConstraintsDotRadius, b, drw.Theme.ConstraintDampedSpringDot)
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
			drw.drawSegment(verts[i], verts[i+1], drw.Theme.ConstraintDampedSpringSegment, strokeWidth)
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

// DrawSpace draws all shapes in space with the drawer implementation
func (drw *Drawer) DrawSpace(space *cm.Space, screen *ebiten.Image) {
	drw.Screen = screen

	if !drw.DrawingOptions.StaticBodyDisabled {
		space.EachStaticShape(func(shape *cm.Shape) {
			drw.DrawShape(shape, drw.Theme.StaticBodyStroke, drw.Theme.StaticBodyFill, drw.DrawingOptions.StaticBodyStrokeWidth)
		})
	}

	if !drw.DrawingOptions.DynamicBodyDisabled {
		space.EachDynamicShape(func(shape *cm.Shape) {

			var clr cm.FColor

			if shape.Body.IsSleeping() {
				clr = drw.Theme.DynamicBodySleepingFill
			} else if shape.Body.IdleTime() > shape.Space.SleepTimeThreshold {
				clr = drw.Theme.DynamicBodyIdleFill
			} else {
				clr = drw.Theme.DynamicBodyFill
			}

			drw.DrawShape(shape, drw.Theme.DynamicBodyStroke, clr, drw.DrawingOptions.DynamicBodyStrokeWidth)
		})
	}

	if !drw.DrawingOptions.ConstraintDisabled {
		space.EachConstraint(func(c *cm.Constraint) {
			drw.DrawConstraint(c, drw.DrawingOptions.ConstraintsStrokeWidth)
		})

	}
	// Draw Collision Point
	if !drw.DrawingOptions.CollisionNormalDisabled {
		for _, arb := range space.Arbiters {

			bodyA, bodyB := arb.Bodies()

			n := arb.Normal()
			for j := 0; j < arb.Count(); j++ {
				p1 := bodyA.Position().Add(arb.Contacts[j].R1)
				p2 := bodyB.Position().Add(arb.Contacts[j].R2)
				a := p1.Add(n.Scale(-drw.DrawingOptions.CollisionNormalLength / 2))
				b := p2.Add(n.Scale(drw.DrawingOptions.CollisionNormalLength / 2))
				drw.drawSegment(a, b, drw.Theme.CollisionNormal, drw.DrawingOptions.CollisionNormalStrokeWidth)
			}
		}

	}
}
