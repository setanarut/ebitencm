package ebitencm

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

const DrawPointLineScale = 1.0

type Drawer struct {
	Screen      *ebiten.Image
	AntiAlias   bool
	StrokeWidth float32
	FlipYAxis   bool

	// private
	handler    mouseEventHandler
	whiteImage *ebiten.Image
}

func NewDrawer() *Drawer {
	whiteImage := ebiten.NewImage(3, 3)
	whiteImage.Fill(color.White)
	return &Drawer{
		AntiAlias:   true,
		StrokeWidth: 1,
		FlipYAxis:   true,
		whiteImage:  whiteImage,
	}
}

func (d *Drawer) WithScreen(screen *ebiten.Image) *Drawer {
	d.Screen = screen
	return d
}

func (d *Drawer) DrawCircle(pos vec.Vec2, angle, radius float64, outline, fill cm.FColor, data interface{}) {
	var f float64 = 1
	if d.FlipYAxis {
		f = -1
		angle *= f
	}

	path := &vector.Path{}
	path.Arc(float32(pos.X), -float32(pos.Y*f), float32(radius), 0, 2*math.Pi, vector.Clockwise)
	d.drawFill(d.Screen, *path, fill.R, fill.G, fill.B, fill.A)

	path.MoveTo(float32(pos.X), -float32(pos.Y*f))
	path.LineTo(float32(pos.X+math.Cos(angle)*radius), -float32(pos.Y*f+math.Sin(angle)*radius))
	path.Close()
	d.drawOutline(d.Screen, *path, outline.R, outline.G, outline.B, outline.A)
}

func (d *Drawer) DrawSegment(a, b vec.Vec2, fill cm.FColor, data interface{}) {
	var f float64 = 1
	if d.FlipYAxis {
		f = -1
	}

	var path *vector.Path = &vector.Path{}
	path.MoveTo(float32(a.X), -float32(a.Y*f))
	path.LineTo(float32(b.X), -float32(b.Y*f))
	path.Close()

	d.drawFill(d.Screen, *path, fill.R, fill.G, fill.B, fill.A)
	d.drawOutline(d.Screen, *path, fill.R, fill.G, fill.B, fill.A)
}

func (d *Drawer) DrawFatSegment(a, b vec.Vec2, radius float64, outline, fill cm.FColor, data interface{}) {
	var f float64 = 1
	if d.FlipYAxis {
		f = -1
	}

	var path vector.Path = vector.Path{}
	t1 := -float32(math.Atan2(b.Y*f-a.Y*f, b.X-a.X)) + math.Pi/2
	t2 := t1 + math.Pi
	path.Arc(float32(a.X), -float32(a.Y*f), float32(radius), t1, t1+math.Pi, vector.Clockwise)
	path.Arc(float32(b.X), -float32(b.Y*f), float32(radius), t2, t2+math.Pi, vector.Clockwise)
	path.Close()
	d.drawFill(d.Screen, path, fill.R, fill.G, fill.B, fill.A)
	d.drawOutline(d.Screen, path, outline.R, outline.G, outline.B, outline.A)
}

func (d *Drawer) DrawPolygon(count int, verts []vec.Vec2, radius float64, outline, fill cm.FColor, data interface{}) {
	var flipYFactor float64 = 1
	if d.FlipYAxis {
		flipYFactor = -1
	}

	type ExtrudeVerts struct {
		offset, n vec.Vec2
	}
	extrude := make([]ExtrudeVerts, count)

	for i := 0; i < count; i++ {
		v0 := verts[(i-1+count)%count]
		v1 := verts[i]
		v2 := verts[(i+1)%count]

		n1 := v1.Sub(v0).ReversePerp().Unit()
		n2 := v2.Sub(v1).ReversePerp().Unit()

		offset := n1.Add(n2).Scale(1.0 / (n1.Dot(n2) + 1.0))
		extrude[i] = ExtrudeVerts{offset, n2}
	}

	var path *vector.Path = &vector.Path{}

	inset := -math.Max(0, 1.0/DrawPointLineScale-radius)
	for i := 0; i < count-2; i++ {
		v0 := verts[0].Add(extrude[0].offset.Scale(inset))
		v1 := verts[i+1].Add(extrude[i+1].offset.Scale(inset))
		v2 := verts[i+2].Add(extrude[i+2].offset.Scale(inset))

		path.MoveTo(float32(v0.X), -float32(v0.Y*flipYFactor))
		path.LineTo(float32(v1.X), -float32(v1.Y*flipYFactor))
		path.LineTo(float32(v2.X), -float32(v2.Y*flipYFactor))
		path.LineTo(float32(v0.X), -float32(v0.Y*flipYFactor))
	}

	outset := 1.0/DrawPointLineScale + radius - inset
	j := count - 1
	for i := 0; i < count; {
		vA := verts[i]
		vB := verts[j]

		nA := extrude[i].n
		nB := extrude[j].n

		offsetA := extrude[i].offset
		offsetB := extrude[j].offset

		innerA := vA.Add(offsetA.Scale(inset))
		innerB := vB.Add(offsetB.Scale(inset))

		inner0 := innerA
		inner1 := innerB
		outer0 := innerA.Add(nB.Scale(outset))
		outer1 := innerB.Add(nB.Scale(outset))
		outer2 := innerA.Add(offsetA.Scale(outset))
		outer3 := innerA.Add(nA.Scale(outset))

		path.MoveTo(float32(inner0.X), -float32(inner0.Y*flipYFactor))
		path.LineTo(float32(inner1.X), -float32(inner1.Y*flipYFactor))
		path.LineTo(float32(outer1.X), -float32(outer1.Y*flipYFactor))
		path.LineTo(float32(inner0.X), -float32(inner0.Y*flipYFactor))

		path.MoveTo(float32(inner0.X), -float32(inner0.Y*flipYFactor))
		path.LineTo(float32(outer0.X), -float32(outer0.Y*flipYFactor))
		path.LineTo(float32(outer1.X), -float32(outer1.Y*flipYFactor))
		path.LineTo(float32(inner0.X), -float32(inner0.Y*flipYFactor))

		path.MoveTo(float32(inner0.X), -float32(inner0.Y*flipYFactor))
		path.LineTo(float32(outer0.X), -float32(outer0.Y*flipYFactor))
		path.LineTo(float32(outer2.X), -float32(outer2.Y*flipYFactor))
		path.LineTo(float32(inner0.X), -float32(inner0.Y*flipYFactor))

		path.MoveTo(float32(inner0.X), -float32(inner0.Y*flipYFactor))
		path.LineTo(float32(outer2.X), -float32(outer2.Y*flipYFactor))
		path.LineTo(float32(outer3.X), -float32(outer3.Y*flipYFactor))
		path.LineTo(float32(inner0.X), -float32(inner0.Y*flipYFactor))

		j = i
		i++
	}

	d.drawFill(d.Screen, *path, fill.R, fill.G, fill.B, fill.A)
}
func (d *Drawer) DrawDot(size float64, pos vec.Vec2, fill cm.FColor, data interface{}) {
	var f float64 = 1
	if d.FlipYAxis {
		f = -1
	}

	var path *vector.Path = &vector.Path{}
	path.Arc(
		float32(pos.X),
		-float32(pos.Y*f),
		float32(2),
		0, 2*math.Pi, vector.Clockwise)
	path.Close()

	d.drawFill(d.Screen, *path, fill.R, fill.G, fill.B, fill.A)
}

func (d *Drawer) Flags() uint {
	return 0
}

func (d *Drawer) OutlineColor() cm.FColor {
	return cm.FColor{R: 200.0 / 255.0, G: 210.0 / 255.0, B: 230.0 / 255.0, A: 1}
}

func (d *Drawer) ShapeColor(shape *cm.Shape, data interface{}) cm.FColor {
	body := shape.Body()
	if body.IsSleeping() {
		return cm.FColor{R: .2, G: .2, B: .2, A: 0.5}
	}

	if body.IdleTime() > shape.Space().SleepTimeThreshold {
		return cm.FColor{R: .66, G: .66, B: .66, A: 0.5}
	}
	return cm.FColor{R: 0.7, G: 0.3, B: 0.6, A: 0.5}
}

func (d *Drawer) ConstraintColor() cm.FColor {
	return cm.FColor{R: 0, G: 0.75, B: 0, A: 1}
}

func (d *Drawer) CollisionPointColor() cm.FColor {
	return cm.FColor{R: 1, G: 0.1, B: 0.2, A: 1}
}

func (d *Drawer) Data() interface{} {
	return nil
}

func (d *Drawer) HandleMouseEvent(space *cm.Space) {
	d.handler.handleMouseEvent(
		d,
		space,
	)
}
func (d *Drawer) drawOutline(
	screen *ebiten.Image,
	path vector.Path,
	r, g, b, a float32,
) {
	sop := &vector.StrokeOptions{}
	sop.Width = d.StrokeWidth
	sop.LineJoin = vector.LineJoinRound
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, sop)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = r
		vs[i].ColorG = g
		vs[i].ColorB = b
		vs[i].ColorA = a
	}
	op := &ebiten.DrawTrianglesOptions{}
	op.FillRule = ebiten.FillAll
	op.AntiAlias = d.AntiAlias
	screen.DrawTriangles(vs, is, d.whiteImage, op)
}

func (d *Drawer) drawFill(
	screen *ebiten.Image,
	path vector.Path,
	r, g, b, a float32,
) {
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = r
		vs[i].ColorG = g
		vs[i].ColorB = b
		vs[i].ColorA = a
	}
	op := &ebiten.DrawTrianglesOptions{}
	op.FillRule = ebiten.FillAll
	op.AntiAlias = d.AntiAlias
	screen.DrawTriangles(vs, is, d.whiteImage, op)
}
