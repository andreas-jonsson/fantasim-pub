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
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/andreas-jonsson/fantasim-pub/frontends/common/game"
	sys "github.com/andreas-jonsson/fantasim-pub/frontends/js/platform"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
)

func throw(err error) {
	js.Global.Call("alert", err.Error())
	panic(err)
}

func assert(err error) {
	if err != nil {
		throw(err)
	}
}

func load() {
	location := js.Global.Get("location")
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

	apiWs, err := websocket.Dial(fmt.Sprintf("ws://%s/api", serverAddress))
	assert(err)
	defer apiWs.Close()

	apiSocket := newTimeoutReadWriter(apiWs, time.Second*3)

	infoWs, err := websocket.Dial(fmt.Sprintf("ws://%s/info", serverAddress))
	assert(err)
	defer infoWs.Close()

	assert(json.NewEncoder(io.MultiWriter(apiSocket, infoWs)).Encode(&playerKey))

	if err := json.NewEncoder(apiSocket).Encode("gob"); err != nil {
		log.Fatalln(err)
	}

	enc := gob.NewEncoder(apiSocket)
	dec := gob.NewDecoder(apiSocket)
	decInfo := gob.NewDecoder(infoWs)

	s := sys.InitJS(image.Pt(640, 400))
	defer s.Quit()

	// Clear screen at exit
	defer s.Present(image.NewRGBA(image.Rectangle{Max: s.Resolution()}))

	game.GameFps = 15
	game.RequestFps = 1

	assert(game.Initialize(s, s))
	js.Global.Get("document").Set("title", "Fantasim JS")
	assert(game.Start(enc, dec, decInfo))
}

func main() {
	js.Global.Call("addEventListener", "load", func() { go load() })
}

type TimeoutReadWriter struct {
	conn net.Conn
	t    time.Duration
}

func newTimeoutReadWriter(conn net.Conn, timeout time.Duration) io.ReadWriter {
	return &TimeoutReadWriter{conn, timeout}
}

func (trw *TimeoutReadWriter) Read(p []byte) (int, error) {
	conn := trw.conn
	conn.SetReadDeadline(time.Now().Add(trw.t))
	return conn.Read(p)
}

func (trw *TimeoutReadWriter) Write(b []byte) (int, error) {
	conn := trw.conn
	conn.SetWriteDeadline(time.Now().Add(trw.t))
	return conn.Write(b)
}
