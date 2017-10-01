/*
Copyright (C) 2017 Andreas T Jonsson

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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
)

const fantasimVersionString = "0.0.1"

const (
	imgWidth  = 320
	imgHeight = 200
	imgScale  = 1
)

var (
	keys   = map[int]bool{}
	canvas *js.Object

	playerKey, playerName string
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

func setupConnection(host string) {
	//ctx := canvas.Call("getContext", "2d")
	//img := ctx.Call("getImageData", 0, 0, imgWidth, imgHeight)

	//if img.Get("data").Length() != len(finalImage.Pix) {
	//	throw(errors.New("data size of images do not match"))
	//}

	ws, err := websocket.Dial(fmt.Sprintf("ws://%s/api", host))
	//ws, err := websocket.Dial("ws://localhost/api")
	assert(err)
	defer ws.Close()

	go func() {
		enc := json.NewEncoder(ws)
		assert(enc.Encode(&playerKey))

		fmt.Println("encodeing done!")
	}()

	dec := json.NewDecoder(ws)

	var version string
	assert(dec.Decode(&version))
	if version != fantasimVersionString {
		throw(fmt.Errorf("invalid version %s, expected %s", version, fantasimVersionString))
	}

	fmt.Println("decodeing done!")

	for {
	}
}

func updateTitle() {
	js.Global.Get("document").Set("title", "Fantasim")
}

func load() {
	document := js.Global.Get("document")

	location := js.Global.Get("location")
	urlStr := strings.TrimPrefix(location.Get("search").String(), "?")
	v, err := url.ParseQuery(urlStr)
	assert(err)

	playerKey = v.Get("key")
	playerName = v.Get("name")

	if playerKey == "" || playerName == "" {
		throw(errors.New("invalid parameters"))
	}

	document.Set("onkeydown", func(e *js.Object) {
		keys[e.Get("keyCode").Int()] = true
	})

	document.Set("onkeyup", func(e *js.Object) {
		keys[e.Get("keyCode").Int()] = false
	})

	canvas = document.Call("createElement", "canvas")
	canvas.Call("setAttribute", "width", strconv.Itoa(imgWidth))
	canvas.Call("setAttribute", "height", strconv.Itoa(imgHeight))
	canvas.Get("style").Set("width", strconv.Itoa(imgWidth*imgScale)+"px")
	canvas.Get("style").Set("height", strconv.Itoa(imgHeight*imgScale)+"px")
	document.Get("body").Call("appendChild", canvas)

	setupConnection(location.Get("host").String())
}

func main() {
	js.Global.Call("addEventListener", "load", func() { go load() })
}
