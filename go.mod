module github.com/setanarut/ebitencm

go 1.23.1

require github.com/hajimehoshi/ebiten/v2 v2.7.10

require github.com/setanarut/vec v1.1.0

require github.com/ojrac/opensimplex-go v1.0.2 // indirect

require (
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.7.1 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/setanarut/cm v1.9.0
	github.com/setanarut/kamera/v2 v2.5.2
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

retract [v1.1.0, v1.1.1]
