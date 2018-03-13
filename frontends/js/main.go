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
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/andreas-jonsson/fantasim-pub/frontends/common/game"
	sys "github.com/andreas-jonsson/fantasim-pub/frontends/js/platform"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
)

var playerKey, playerName string

func throw(err error) {
	js.Global.Call("alert", err.Error())
	panic(err)
}

func assert(err error) {
	if err != nil {
		throw(err)
	}
}

func updateTitle() {
	js.Global.Get("document").Set("title", "Fantasim JS")
}

func load() {
	location := js.Global.Get("location")
	urlStr := strings.TrimPrefix(location.Get("search").String(), "?")
	v, err := url.ParseQuery(urlStr)
	assert(err)

	playerKey = v.Get("key")
	playerName = v.Get("name")

	if playerKey == "" || playerName == "" {
		throw(errors.New("invalid parameters"))
	}

	serverAddress := v.Get("host")
	if serverAddress == "" {
		serverAddress = location.Get("host").String()
	}

	apiWs, err := websocket.Dial(fmt.Sprintf("ws://%s/api", serverAddress))
	if err != nil {
		log.Fatalln(err)
	}
	defer apiWs.Close()

	var apiSocket io.ReadWriter = apiWs
	//if !*noTimeout {
	//	apiSocket = newTimeoutReadWriter(apiWs, time.Second*3)
	//}

	infoWs, err := websocket.Dial(fmt.Sprintf("ws://%s/info", serverAddress))
	if err != nil {
		log.Fatalln(err)
	}
	defer infoWs.Close()

	if err := json.NewEncoder(io.MultiWriter(apiSocket, infoWs)).Encode(&playerKey); err != nil {
		log.Fatalln(err)
	}

	if err := json.NewEncoder(apiSocket).Encode("json"); err != nil {
		log.Fatalln(err)
	}
	enc := json.NewEncoder(apiSocket)
	dec := json.NewDecoder(apiSocket)
	decInfo := json.NewDecoder(infoWs)

	s := sys.InitJS(image.Pt(640, 400))
	defer s.Quit()

	if err := game.Initialize(s, s); err != nil {
		log.Fatalln(err)
	}

	updateTitle()
	if err := game.Start(enc, dec, decInfo); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	js.Global.Call("addEventListener", "load", func() { go load() })
}
