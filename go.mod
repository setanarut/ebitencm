module github.com/setanarut/ebitencm

go 1.23.2

require github.com/hajimehoshi/ebiten/v2 v2.8.3

require github.com/setanarut/vec v1.1.1

require github.com/setanarut/fastnoise v1.1.1 // indirect

require (
	github.com/ebitengine/gomobile v0.0.0-20241016134836-cc2e38a7c0ee // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.1 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/setanarut/cm v1.12.0
	github.com/setanarut/kamera/v2 v2.7.0
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
)

retract [v1.0.0, v1.5.1]
