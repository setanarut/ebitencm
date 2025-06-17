package ebitencm

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/cm"
	"github.com/setanarut/v"
)

type Drawer struct {
	// Ebitengine screen
	Screen *ebiten.Image

	// Drawing colors
	Theme          *Theme
	DrawingOptions *DrawingOptions
	// GeoM for drawing vertices. Useful for cameras
	GeoM *ebiten.GeoM

	DrawTriangleStrokeOpt *ebiten.DrawTrianglesOptions
	DrawTriagleFillOpt    *ebiten.DrawTrianglesOptions
	// private
	handler    mouseEventHandler
	whiteImage *ebiten.Image
}

func NewDrawer() *Drawer {
	whiteImage := ebiten.NewImage(3, 3)
	whiteImage.Fill(color.White)
	return &Drawer{
		DrawingOptions:        DefaultDrawingOptions(),
		DrawTriangleStrokeOpt: &ebiten.DrawTrianglesOptions{AntiAlias: true},
		DrawTriagleFillOpt:    &ebiten.DrawTrianglesOptions{AntiAlias: true},
		whiteImage:            whiteImage,
		GeoM:                  &ebiten.GeoM{},
		Theme:                 DefaultTheme(),
	}
}

func (d *Drawer) SetScreen(screen *ebiten.Image) {
	d.Screen = screen
}
func (d *Drawer) SetStrokeAntialias(antialias bool) {
	d.DrawTriangleStrokeOpt.AntiAlias = antialias
}

func (d *Drawer) SetFillAntialias(antialias bool) {
	d.DrawTriagleFillOpt.AntiAlias = antialias
}

func (d *Drawer) drawCircle(
	pos v.Vec,
	angle, radius float64,
	outline, fill cm.FColor,
	strokeWidth float32,
) {
	// angle *= flipFactor
	path := &vector.Path{}
	path.Arc(
		float32(pos.X),
		float32(pos.Y),
		float32(radius),
		0,
		2*math.Pi,
		vector.Clockwise,
	)
	// Fill
	if !d.DrawingOptions.AllFillsDisabled {
		d.fillPath(d.Screen, *path, fill.R, fill.G, fill.B, fill.A)
	}
	path.MoveTo(
		float32(pos.X),
		float32(pos.Y),
	)
	path.LineTo(
		float32(pos.X+math.Cos(angle)*radius),
		float32(pos.Y+math.Sin(angle)*radius),
	)
	path.Close()
	// Stroke
	if !d.DrawingOptions.AllStrokesDisabled {
		d.strokePath(d.Screen, *path, outline.R, outline.G, outline.B, outline.A, strokeWidth)
	}
}

func (d *Drawer) drawSegment(a, b v.Vec, clr cm.FColor, strokeWidth float32) {

	var path vector.Path = vector.Path{}
	path.MoveTo(float32(a.X), float32(a.Y))
	path.LineTo(float32(b.X), float32(b.Y))
	path.Close()
	if !d.DrawingOptions.AllFillsDisabled {
		d.fillPath(d.Screen, path, clr.R, clr.G, clr.B, clr.A)
	}
	if !d.DrawingOptions.AllStrokesDisabled {
		d.strokePath(d.Screen, path, clr.R, clr.G, clr.B, clr.A, strokeWidth)
	}
}

func (d *Drawer) drawFatSegment(
	a, b v.Vec,
	radius float64,
	outline, fillColor cm.FColor,
	strokeWidth float32,
) {

	var path vector.Path = vector.Path{}
	t1 := float32(math.Atan2(b.Y-a.Y, b.X-a.X)) + math.Pi/2
	t2 := t1 + math.Pi
	path.Arc(
		float32(a.X),
		float32(a.Y),
		float32(radius),
		t1,
		t1+math.Pi,
		vector.Clockwise,
	)
	path.Arc(
		float32(b.X),
		float32(b.Y),
		float32(radius),
		t2,
		t2+math.Pi,
		vector.Clockwise,
	)
	path.Close()

	if !d.DrawingOptions.AllFillsDisabled {
		d.fillPath(d.Screen, path, fillColor.R, fillColor.G, fillColor.B, fillColor.A)
	}

	if !d.DrawingOptions.AllStrokesDisabled {
		d.strokePath(d.Screen, path, outline.R, outline.G, outline.B, outline.A, strokeWidth)
	}
}

func (d *Drawer) drawPolygon(count int, verts []v.Vec, radius float64, outline, fill cm.FColor, strokeWidth float32) {
	type ExtrudeVerts struct {
		offset, n v.Vec
	}
	extrude := make([]ExtrudeVerts, count)

	for i := 0; i < count; i++ {
		v0 := verts[(i-1+count)%count]
		v1 := verts[i]
		v2 := verts[(i+1)%count]

		n1 := reversePerp(v1.Sub(v0)).Unit()
		n2 := reversePerp(v2.Sub(v1)).Unit()

		offset := n1.Add(n2).Scale(1.0 / (n1.Dot(n2) + 1.0))
		extrude[i] = ExtrudeVerts{offset, n2}
	}

	var path *vector.Path = &vector.Path{}

	// insetScaleFactor := -math.Max(0, 1.0/drawPointLineScale-radius) // neg
	// outsetScaleFactor := 1.0/drawPointLineScale + radius - insetScaleFactor
	// outset2ScaleFactor := 1.0/drawPointLineScale + radius - insetScaleFactor
	insetScaleFactor := -1.
	outsetScaleFactor := 1.
	outset2ScaleFactor := 1.

	j := count - 1
	for i := 0; i < count; {
		vA := verts[i]
		vB := verts[j]
		nA := extrude[i].n
		nB := extrude[j].n
		offsetA := extrude[i].offset
		offsetB := extrude[j].offset
		innerA := vA.Add(offsetA.Scale(insetScaleFactor))
		innerB := vB.Add(offsetB.Scale(insetScaleFactor))
		outer0 := innerA.Add(nB.Scale(outsetScaleFactor))
		outer1 := innerB.Add(nB.Scale(outsetScaleFactor))
		outer2 := innerA.Add(offsetA.Scale(outsetScaleFactor))
		outer3 := innerA.Add(offsetA.Scale(outset2ScaleFactor))
		outer4 := innerA.Add(nA.Scale(outsetScaleFactor))
		path.LineTo(float32(outer1.X), float32(outer1.Y))
		path.LineTo(float32(outer0.X), float32(outer0.Y))

		if radius != 0 {
			path.ArcTo(float32(outer3.X), float32(outer3.Y), float32(outer4.X), float32(outer4.Y), float32(radius))
		} else {
			// ArcTo() and Arc() are very computationally expensive, so use LineTo()
			path.LineTo(
				float32(outer2.X),
				float32(outer2.Y))
		}
		j = i
		i++
	}
	path.Close()

	if !d.DrawingOptions.AllFillsDisabled {
		d.fillPath(d.Screen, *path, fill.R, fill.G, fill.B, fill.A)
	}
	if !d.DrawingOptions.AllStrokesDisabled {
		d.strokePath(d.Screen, *path, outline.R, outline.G, outline.B, outline.A, strokeWidth)
	}
}

func (d *Drawer) drawDot(radius float64, pos v.Vec, fill cm.FColor) {
	if !d.DrawingOptions.AllDotsDisabled {
		var path *vector.Path = &vector.Path{}
		path.Arc(
			float32(pos.X),
			float32(pos.Y),
			float32(radius),
			0,
			2*math.Pi,
			vector.Clockwise,
		)
		path.Close()
		d.fillPath(d.Screen, *path, fill.R, fill.G, fill.B, fill.A)
	}
}

func (d *Drawer) HandleMouseEvent(space *cm.Space) {
	d.handler.handleMouseEvent(d, space)
}

func (d *Drawer) strokePath(screen *ebiten.Image, path vector.Path, r, g, b, a float32, w float32) {
	sop := &vector.StrokeOptions{}
	sop.Width = w
	sop.LineJoin = vector.LineJoinRound
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, sop)
	applyMatrixToVertices(vs, d.GeoM, r, g, b, a)
	screen.DrawTriangles(vs, is, d.whiteImage, d.DrawTriangleStrokeOpt)
}

func (d *Drawer) fillPath(screen *ebiten.Image, path vector.Path, r, g, b, a float32) {
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	applyMatrixToVertices(vs, d.GeoM, r, g, b, a)
	screen.DrawTriangles(vs, is, d.whiteImage, d.DrawTriagleFillOpt)
}

func applyMatrixToVertices(vs []ebiten.Vertex, matrix *ebiten.GeoM, r, g, b, a float32) {
	for i := range vs {
		x, y := matrix.Apply(float64(vs[i].DstX), float64(vs[i].DstY))
		vs[i].DstX, vs[i].DstY = float32(x), float32(y)
		vs[i].SrcX, vs[i].SrcY = 1, 1
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = r, g, b, a
	}
}

// ScreenToWorld converts screen-space coordinates to world-space
func ScreenToWorld(screenPoint v.Vec, cameraGeoM ebiten.GeoM) v.Vec {
	if cameraGeoM.IsInvertible() {
		cameraGeoM.Invert()
		worldX, worldY := cameraGeoM.Apply(screenPoint.X, screenPoint.Y)
		return v.Vec{worldX, worldY}
	} else {
		// When scaling it can happened that matrix is not invertable
		return v.Vec{math.NaN(), math.NaN()}
	}
}

func ToFColor(c color.RGBA) cm.FColor {
	r := float32(c.R) / 255.0
	g := float32(c.G) / 255.0
	b := float32(c.B) / 255.0
	a := float32(c.A) / 255.0
	return cm.FColor{r, g, b, a}
}

type Theme struct {
	CollisionNormal               cm.FColor
	ConstraintDampedSpringDot     cm.FColor
	ConstraintDampedSpringSegment cm.FColor
	ConstraintGrooveJointDot      cm.FColor
	ConstraintGrooveJointSegment  cm.FColor
	ConstraintPinJointDot         cm.FColor
	ConstraintPinJointSegment     cm.FColor
	ConstraintPivotJointDot       cm.FColor
	ConstraintSlideJointDot       cm.FColor
	ConstraintSlideJointSegment   cm.FColor
	DynamicBodyFill               cm.FColor
	DynamicBodyIdleFill           cm.FColor
	DynamicBodySleepingFill       cm.FColor
	DynamicBodyStroke             cm.FColor
	StaticBodyFill                cm.FColor
	StaticBodyStroke              cm.FColor
}

// SetOpacity overwrites all Theme color alphas [0-1}]
func (d *Drawer) SetOpacity(alpha float32) {
	d.Theme.CollisionNormal.A = alpha
	d.Theme.ConstraintDampedSpringDot.A = alpha
	d.Theme.ConstraintDampedSpringSegment.A = alpha
	d.Theme.ConstraintGrooveJointDot.A = alpha
	d.Theme.ConstraintGrooveJointSegment.A = alpha
	d.Theme.ConstraintPinJointDot.A = alpha
	d.Theme.ConstraintPinJointSegment.A = alpha
	d.Theme.ConstraintPivotJointDot.A = alpha
	d.Theme.ConstraintSlideJointDot.A = alpha
	d.Theme.ConstraintSlideJointSegment.A = alpha
	d.Theme.DynamicBodyFill.A = alpha
	d.Theme.DynamicBodyIdleFill.A = alpha
	d.Theme.DynamicBodySleepingFill.A = alpha
	d.Theme.DynamicBodyStroke.A = alpha
	d.Theme.StaticBodyFill.A = alpha
	d.Theme.StaticBodyStroke.A = alpha
}

func DefaultTheme() *Theme {
	return &Theme{
		CollisionNormal:               cm.FColor{1, 1, 0, 1},
		ConstraintDampedSpringDot:     cm.FColor{1, 0.7, 0.7, 1},
		ConstraintDampedSpringSegment: cm.FColor{1, 0.7, 0.7, 1},
		ConstraintGrooveJointDot:      cm.FColor{0, 0.7, 0.7, 1},
		ConstraintGrooveJointSegment:  cm.FColor{0, 0.7, 0.7, 1},
		ConstraintPinJointDot:         cm.FColor{0, 0.7, 0.7, 1},
		ConstraintPinJointSegment:     cm.FColor{0, 0.7, 0.7, 1},
		ConstraintPivotJointDot:       cm.FColor{0, 0.7, 0.7, 1},
		ConstraintSlideJointDot:       cm.FColor{0, 0.7, 0.7, 1},
		ConstraintSlideJointSegment:   cm.FColor{0, 0.7, 0.7, 1},
		DynamicBodyFill:               cm.FColor{0, 0, 1, 1},
		DynamicBodyIdleFill:           cm.FColor{0.5, 0.5, 0.5, 1},
		DynamicBodySleepingFill:       cm.FColor{0.5, 0.5, 0.5, 1},
		DynamicBodyStroke:             cm.FColor{0.69, 0.165, 0.537, 1},
		StaticBodyFill:                cm.FColor{0.6, 0.3, 0.5, 1},
		StaticBodyStroke:              cm.FColor{0.69, 0.165, 0.537, 1},
	}
}

type DrawingOptions struct {
	AllDotsDisabled            bool
	AllFillsDisabled           bool
	AllStrokesDisabled         bool
	CollisionNormalDisabled    bool
	CollisionNormalLength      float64
	CollisionNormalStrokeWidth float32
	ConstraintDisabled         bool
	ConstraintsDotRadius       float64
	ConstraintsStrokeWidth     float32
	DynamicBodyDisabled        bool
	DynamicBodyStrokeWidth     float32
	StaticBodyDisabled         bool
	StaticBodyStrokeWidth      float32
}

func DefaultDrawingOptions() *DrawingOptions {
	return &DrawingOptions{
		AllDotsDisabled:            false,
		AllFillsDisabled:           false,
		AllStrokesDisabled:         false,
		CollisionNormalDisabled:    false,
		CollisionNormalLength:      12,
		CollisionNormalStrokeWidth: 2,
		ConstraintDisabled:         false,
		ConstraintsDotRadius:       2,
		ConstraintsStrokeWidth:     2,
		DynamicBodyDisabled:        false,
		DynamicBodyStrokeWidth:     2,
		StaticBodyDisabled:         false,
		StaticBodyStrokeWidth:      2,
	}
}

// ReversePerp returns a perpendicular vector. (-90 degree rotation)
func reversePerp(a v.Vec) v.Vec {
	return v.Vec{a.Y, -a.X}
}
