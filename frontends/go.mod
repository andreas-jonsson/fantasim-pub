module github.com/andreas-jonsson/fantasim-pub/frontends

go 1.12

replace github.com/andreas-jonsson/fantasim-pub/api => ../api

require (
	github.com/andreas-jonsson/fantasim-pub/api v0.0.0
	github.com/dennwc/dom v0.3.0
	github.com/gopherjs/gopherjs v0.0.0-20190328170749-bb2674552d8f
	github.com/gopherjs/websocket v0.0.0-20170522004412-87ee47603f13
	github.com/ojrac/opensimplex-go v1.0.1
	github.com/ungerik/go3d v0.0.0-20190319220834-8e1a82526839
	github.com/veandco/go-sdl2 v0.3.0
	golang.org/x/mobile v0.0.0-20190327163128-167ebed0ec6d
	golang.org/x/net v0.0.0-20190328230028-74de082e2cca
)
