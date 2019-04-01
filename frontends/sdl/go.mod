module github.com/andreas-jonsson/fantasim-pub/frontends/sdl

go 1.12

require (
	github.com/andreas-jonsson/fantasim-pub/api v0.0.0-00010101000000-000000000000 // indirect
	github.com/andreas-jonsson/fantasim-pub/frontends/common v0.0.0-00010101000000-000000000000
	github.com/ojrac/opensimplex-go v1.0.1 // indirect
	github.com/veandco/go-sdl2 v0.3.0
	golang.org/x/net v0.0.0-20190328230028-74de082e2cca
)

replace github.com/andreas-jonsson/fantasim-pub/api => ../../api

replace github.com/andreas-jonsson/fantasim-pub/frontends/common => ../common
