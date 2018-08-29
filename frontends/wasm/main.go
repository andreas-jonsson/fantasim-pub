/*
Copyright (C) 2017-2018 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

//go:generate go run ../common/data/generate.go

package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"net/url"
	"strings"
	"syscall/js"

	"github.com/andreas-jonsson/fantasim-pub/frontends/common/game"
	sys "github.com/andreas-jonsson/fantasim-pub/frontends/wasm/platform"
)

func throw(err error) {
	js.Global().Call("alert", err.Error())
	panic(err)
}

func assert(err error) {
	if err != nil {
		throw(err)
	}
}

func main() {
	location := js.Global().Get("location")
	urlStr := strings.TrimPrefix(location.Get("search").String(), "?")
	v, err := url.ParseQuery(urlStr)
	assert(err)

	playerKey := v.Get("key")
	if playerKey == "" {
		throw(errors.New("invalid parameters"))
	}

	serverAddress := v.Get("host")
	if serverAddress == "" {
		serverAddress = location.Get("host").String()
	}

	apiWs, err := sys.Dial(fmt.Sprintf("ws://%s/api", serverAddress))
	assert(err)

	infoWs, err := sys.Dial(fmt.Sprintf("ws://%s/info", serverAddress))
	assert(err)

	assert(json.NewEncoder(io.MultiWriter(apiWs, infoWs)).Encode(&playerKey))

	if err := json.NewEncoder(apiWs).Encode("gob"); err != nil {
		log.Fatalln(err)
	}

	enc := gob.NewEncoder(apiWs)
	dec := gob.NewDecoder(apiWs)
	decInfo := gob.NewDecoder(infoWs)

	s := sys.InitWASM(image.Pt(640, 400))
	defer s.Quit()

	if err := game.Initialize(s, s); err != nil {
		log.Fatalln(err)
	}

	game.GameFps = 10
	game.RequestFps = 1

	assert(game.Start(enc, dec, decInfo))
}
