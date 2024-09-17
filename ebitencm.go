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
const flipFactor = -1.0

type Drawer struct {
	// Ebitengine screen
	Screen      *ebiten.Image
	StrokeWidth float32
	// GeoM for drawing vertices. Useful for cameras
	GeoM *ebiten.GeoM
	// Disable filling except DrawDot().
	FillDisabled bool
	// Disable strokes
	StrokeDisabled bool

	OptStroke *ebiten.DrawTrianglesOptions
	OptFill   *ebiten.DrawTrianglesOptions

	// private
	handler    mouseEventHandler
	whiteImage *ebiten.Image
}

func NewDrawer() *Drawer {
	whiteImage := ebiten.NewImage(3, 3)
	whiteImage.Fill(color.White)
	return &Drawer{
		StrokeWidth: 1,
		OptStroke:   &ebiten.DrawTrianglesOptions{},
		OptFill:     &ebiten.DrawTrianglesOptions{},
		//private
		whiteImage: whiteImage,
		GeoM:       &ebiten.GeoM{},
	}
}

func (d *Drawer) WithScreen(screen *ebiten.Image) *Drawer {
	d.Screen = screen
	return d
}

func (d *Drawer) SetScreen(screen *ebiten.Image) {
	d.Screen = screen
}
func (d *Drawer) SetStrokeAntialias(antialias bool) {
	d.OptStroke.AntiAlias = antialias
}
func (d *Drawer) SetFillAntialias(antialias bool) {
	d.OptFill.AntiAlias = antialias
}

func (d *Drawer) DrawCircle(pos vec.Vec2, angle, radius float64, outline, fill cm.FColor, data interface{}) {
	angle *= flipFactor
	path := &vector.Path{}
	path.Arc(float32(pos.X), -float32(pos.Y*flipFactor), float32(radius), 0, 2*math.Pi, vector.Clockwise)
	// Fill
	if !d.FillDisabled {
		d.fillPath(d.Screen, *path, fill.R, fill.G, fill.B, fill.A)
	}
	path.MoveTo(float32(pos.X), -float32(pos.Y*flipFactor))
	path.LineTo(float32(pos.X+math.Cos(angle)*radius), -float32(pos.Y*flipFactor+math.Sin(angle)*radius))
	path.Close()
	// Stroke
	if !d.StrokeDisabled {
		d.strokePath(d.Screen, *path, outline.R, outline.G, outline.B, outline.A)
	}
}

func (d *Drawer) DrawSegment(a, b vec.Vec2, fillColor cm.FColor, data interface{}) {

	var path vector.Path = vector.Path{}
	path.MoveTo(float32(a.X), -float32(a.Y*flipFactor))
	path.LineTo(float32(b.X), -float32(b.Y*flipFactor))
	path.Close()
	if !d.FillDisabled {
		d.fillPath(d.Screen, path, fillColor.R, fillColor.G, fillColor.B, fillColor.A)
	}
	if !d.StrokeDisabled {
		d.strokePath(d.Screen, path, fillColor.R, fillColor.G, fillColor.B, fillColor.A)
	}
}

func (d *Drawer) DrawFatSegment(a, b vec.Vec2, radius float64, outline, fillColor cm.FColor, data interface{}) {

	var path vector.Path = vector.Path{}
	t1 := -float32(math.Atan2(b.Y*flipFactor-a.Y*flipFactor, b.X-a.X)) + math.Pi/2
	t2 := t1 + math.Pi
	path.Arc(float32(a.X), -float32(a.Y*flipFactor), float32(radius), t1, t1+math.Pi, vector.Clockwise)
	path.Arc(float32(b.X), -float32(b.Y*flipFactor), float32(radius), t2, t2+math.Pi, vector.Clockwise)
	path.Close()

	if !d.FillDisabled {
		d.fillPath(d.Screen, path, fillColor.R, fillColor.G, fillColor.B, fillColor.A)
	}

	if !d.StrokeDisabled {
		d.strokePath(d.Screen, path, outline.R, outline.G, outline.B, outline.A)
	}
}

func (d *Drawer) DrawPolygon(count int, verts []vec.Vec2, radius float64, outline, fill cm.FColor, data interface{}) {
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

		path.MoveTo(float32(v0.X), -float32(v0.Y*flipFactor))
		path.LineTo(float32(v1.X), -float32(v1.Y*flipFactor))
		path.LineTo(float32(v2.X), -float32(v2.Y*flipFactor))
		path.LineTo(float32(v0.X), -float32(v0.Y*flipFactor))
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

		path.MoveTo(float32(inner0.X), -float32(inner0.Y*flipFactor))
		path.LineTo(float32(inner1.X), -float32(inner1.Y*flipFactor))
		path.LineTo(float32(outer1.X), -float32(outer1.Y*flipFactor))
		path.LineTo(float32(inner0.X), -float32(inner0.Y*flipFactor))

		path.MoveTo(float32(inner0.X), -float32(inner0.Y*flipFactor))
		path.LineTo(float32(outer0.X), -float32(outer0.Y*flipFactor))
		path.LineTo(float32(outer1.X), -float32(outer1.Y*flipFactor))
		path.LineTo(float32(inner0.X), -float32(inner0.Y*flipFactor))

		path.MoveTo(float32(inner0.X), -float32(inner0.Y*flipFactor))
		path.LineTo(float32(outer0.X), -float32(outer0.Y*flipFactor))
		path.LineTo(float32(outer2.X), -float32(outer2.Y*flipFactor))
		path.LineTo(float32(inner0.X), -float32(inner0.Y*flipFactor))

		path.MoveTo(float32(inner0.X), -float32(inner0.Y*flipFactor))
		path.LineTo(float32(outer2.X), -float32(outer2.Y*flipFactor))
		path.LineTo(float32(outer3.X), -float32(outer3.Y*flipFactor))
		path.LineTo(float32(inner0.X), -float32(inner0.Y*flipFactor))
		j = i
		i++
	}
	if !d.FillDisabled {
		d.fillPath(d.Screen, *path, fill.R, fill.G, fill.B, fill.A)
	}
}
func (d *Drawer) DrawDot(size float64, pos vec.Vec2, fill cm.FColor, data interface{}) {
	var path *vector.Path = &vector.Path{}
	path.Arc(float32(pos.X), -float32(pos.Y*flipFactor), float32(2), 0, 2*math.Pi, vector.Clockwise)
	path.Close()
	d.fillPath(d.Screen, *path, fill.R, fill.G, fill.B, fill.A)
}

func (d *Drawer) Flags() uint {
	return 0
}

func (d *Drawer) OutlineColor() cm.FColor {
	return cm.FColor{1, 1, 1, 0.8}
}

func (d *Drawer) ShapeColor(shape *cm.Shape, data interface{}) cm.FColor {
	body := shape.Body()

	if body.IsSleeping() {
		return cm.FColor{1, 1, 1, 0.35}
	}

	if body.IdleTime() > shape.Space().SleepTimeThreshold {
		return cm.FColor{1, 1, 1, 0.4}
	}
	return cm.FColor{1, 1, 1, 0.5}
}

func (d *Drawer) ConstraintColor() cm.FColor {
	return cm.FColor{0, 0.75, 0, 1}
}

func (d *Drawer) CollisionPointColor() cm.FColor {
	return cm.FColor{1, 0, 0, 1}
}

func (d *Drawer) Data() interface{} {
	return nil
}

func (d *Drawer) HandleMouseEvent(space *cm.Space) {
	d.handler.handleMouseEvent(d, space)
}
func (d *Drawer) strokePath(screen *ebiten.Image, path vector.Path, r, g, b, a float32) {
	sop := &vector.StrokeOptions{}
	sop.Width = d.StrokeWidth
	sop.LineJoin = vector.LineJoinRound
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, sop)
	applyMatrixToVertices(vs, d.GeoM, r, g, b, a)

	screen.DrawTriangles(vs, is, d.whiteImage, d.OptStroke)
}

func (d *Drawer) fillPath(screen *ebiten.Image, path vector.Path, r, g, b, a float32) {
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	applyMatrixToVertices(vs, d.GeoM, r, g, b, a)
	screen.DrawTriangles(vs, is, d.whiteImage, d.OptFill)
}

// RotateAbout rotates point about origin
func RotateAbout(point vec.Vec2, angle float64, origin vec.Vec2) vec.Vec2 {
	b := vec.Vec2{}
	b.X = math.Cos(angle)*(point.X-origin.X) - math.Sin(angle)*(point.Y-origin.Y) + origin.X
	b.Y = math.Sin(angle)*(point.X-origin.X) + math.Cos(angle)*(point.Y-origin.Y) + origin.Y
	return b
}

func applyMatrixToVertices(vs []ebiten.Vertex, matrix *ebiten.GeoM, r, g, b, a float32) {
	for i := range vs {
		x, y := matrix.Apply(float64(vs[i].DstX), float64(vs[i].DstY))
		vs[i].DstX, vs[i].DstY = float32(x), float32(y)
		vs[i].SrcX, vs[i].SrcY = 1, 1
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = r, g, b, a
	}
}
