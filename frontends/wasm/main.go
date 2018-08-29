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
	"fmt"
	"image"
	"io"
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
		js.Global().Get("location").Set("href", "http://lobby.fantasim.net")
		return
	}

	serverAddress := v.Get("host")
	if serverAddress == "" {
		serverAddress = location.Get("host").String()
	}

	protocol := "ws"
	if location.Get("protocol").String() == "https" {
		protocol = "wss"
	}

	apiWs, err := sys.Dial(fmt.Sprintf("%s://%s/api", protocol, serverAddress))
	assert(err)

	infoWs, err := sys.Dial(fmt.Sprintf("%s://%s/info", protocol, serverAddress))
	assert(err)

	assert(json.NewEncoder(io.MultiWriter(apiWs, infoWs)).Encode(&playerKey))
	assert(json.NewEncoder(apiWs).Encode("gob"))

	enc := gob.NewEncoder(apiWs)
	dec := gob.NewDecoder(apiWs)
	decInfo := gob.NewDecoder(infoWs)

	s := sys.InitWASM(image.Pt(1280, 720), 1)
	defer s.Quit()

	assert(game.Initialize(s, s))

	game.GameFps = 15
	game.RequestFps = 1

	assert(game.Start(enc, dec, decInfo))
}
