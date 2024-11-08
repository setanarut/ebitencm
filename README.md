[![GoDoc](https://godoc.org/github.com/setanarut/ebitencm?status.svg)](https://pkg.go.dev/github.com/setanarut/ebitencm)

# ebitencm

**ebitencm** is *Drawer* for [setanarut/cm](https://github.com/setanarut/cm) Chipmunk physics space.

# Features

- Color theme customization for *Fill* and *Stroke* colors.
- Detailed drawing options.
- `drawer.GeoM{}` structure is provided for screen transformation. (Cameras)

## Usage

First create a drawer

```Go
var drawer *ebitencm.Drawer = ebitencm.NewDrawer()
```

Then

```go
func (g *Game) Draw(screen *ebiten.Image) {
	drawer.DrawSpace(space, screen)
}
```

## Dragging

If you want to enable dragging, call the `HandleMouseEvent()` function within the `Update` method, passing the `*cm.Space` object. This will allow objects to be dragged using a mouse or touch device.

```go
func (g *Game) Update() error {
	drawer.HandleMouseEvent(space)
```

## Camera transform

Use `Drawer.GeoM{}` for all vertices transform. The cursor position is calculated according to this matrix.

```Go
// move the all space objects 100 pixels to the left (move camera to right)
drawer.GeoM.Translate(-100, 0)
```

Here is an example with the [setanarut/kamera](https://github.com/setanarut/kamera) package.

```Go
func (g *Game) Update() error {
	g.space.Step(1 / 60.0)
	g.cam.LookAt(x, y)
	// Apply camera transform to drawer
	g.drawer.GeoM.Reset()
	g.cam.ApplyCameraTransform(g.drawer.GeoM)
	// Enable cursor dragging
	g.drawer.HandleMouseEvent(g.space)
```

## Examples

Browse to the [examples](./examples/) folder for all examples.