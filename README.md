# ebitencm

**ebitencm** (Ebitengine Chipmunk drawer) is an implementation of the `cm.IDrawer` Interface from [setanarut/cm](https://github.com/setanarut/cm). This implementation utilizes [hajimehoshi/ebiten/v2](https://github.com/hajimehoshi/ebiten), making it possible to run across multiple platforms. Coordinate system is top-left by default, same as Ebitengine

![scr](https://github.com/user-attachments/assets/ca27ad36-509e-4f33-b526-372598d3144c)

## Usage

Within your `Draw()` method, invoke the `cm.DrawSpace()` function, passing in both a `*cm.Space` and a `*ebitencm.Drawer` as parameters.

```go
type Game struct {
	space  *cm.Space
	drawer *ebitencm.Drawer
}
func main() {
	game := &Game{}
	game.space = space
	game.drawer = ebitencm.NewDrawer()
}
func (g *Game) Draw(screen *ebiten.Image) {
	// Drawing with Ebitengine/v2
	cm.DrawSpace(g.space, g.drawer.WithScreen(screen))
}
```

## Dragging

If you want to enable dragging, call the `HandleMouseEvent()` function within the `Update` method, passing the `*cm.Space` object. This will allow objects to be dragged using a mouse or touch device.

```go
func (g *Game) Update() error {
	// Handling dragging
	g.drawer.HandleMouseEvent(g.space)
	g.space.Step(1 / 60.0)
	return nil
}
```

## Camera transform

Use `Drawer.GeoM{}` for camera transform. The cursor position is calculated according to this matrix. Don't forget to reset GeoM with `GeoM.Reset()` in `Update()`. Here is an example with the [setanarut/kamera](https://github.com/setanarut/kamera) package.

```Go
func (g *Game) Update() error {
	g.space.Step(1 / 60.0)

	// Apply camera transform to drawer
	g.drawer.GeoM.Reset()
	g.cam.ApplyCameraTransform(g.drawer.GeoM)

	// Enable cursor dragging
	g.drawer.HandleMouseEvent(g.space)
```

### Camera demo

Run camera example on your local machine

```
go run github.com/setanarut/ebitencm/examples/camera@latest
```

- Camera Position = WASD
- Camera Rotation = Q / E
- Camera Zoom = Z / X
- Drag object = Cursor
- Reset camera = Backspace

## Examples

Browse to the [examples](./examples/) folder for all examples.