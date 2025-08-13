module github.com/setanarut/ebitencm

go 1.25.0

require (
	github.com/hajimehoshi/ebiten/v2 v2.8.8
	github.com/setanarut/v v1.2.1
)

require github.com/setanarut/fastnoise v1.1.1 // indirect

require (
	github.com/ebitengine/gomobile v0.0.0-20250329061421-6d0a8e981e4c // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/setanarut/cm v1.14.2
	github.com/setanarut/kamera/v2 v2.97.2
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
)

retract [v1.0.0, v1.7.0]
