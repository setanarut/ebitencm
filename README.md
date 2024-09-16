# ebitencm

**ebitencm** (Ebitengine Chipmunk drawer) is an implementation of the `cm.IDrawer` Interface from [setanarut/cm](https://github.com/setanarut/cm). This implementation utilizes [hajimehoshi/ebiten/v2](https://github.com/hajimehoshi/ebiten), making it possible to run across multiple platforms. Coordinate system is top-left by default, same as Ebitengine


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

## Changing coordinate system

Coordinate system is top-left by default, same as Ebitengine. You can correct the coordinate system by setting FlipYAxis to true.

```Go
  func main() {
  	game.drawer = ebitencm.NewDrawer()
	game.drawer.FlipYAxis = false
  }
```

## Examples

Run `basic` example on your local machine
```
go run github.com/setanarut/ebitencm/examples/basic@latest
```

Browse to the [examples](./examples/) folder for all examples.


