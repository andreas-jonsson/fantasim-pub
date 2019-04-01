module github.com/andreas-jonsson/fantasim-pub/frontends/wasm

go 1.12

require (
	github.com/andreas-jonsson/fantasim-pub/api v0.0.0-00010101000000-000000000000 // indirect
	github.com/andreas-jonsson/fantasim-pub/frontends/common v0.0.0-00010101000000-000000000000
	github.com/ojrac/opensimplex-go v1.0.1 // indirect
)

replace github.com/andreas-jonsson/fantasim-pub/api => ../../api

replace github.com/andreas-jonsson/fantasim-pub/frontends/common => ../common
