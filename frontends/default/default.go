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
	"log"
	"strconv"
	"time"

	"github.com/gopherjs/gopherjs/js"
	websocket "github.com/gopherjs/websocket"
)

const (
	imgWidth  = 320
	imgHeight = 200
	imgScale  = 1
)

var (
	keys   = map[int]bool{}
	canvas *js.Object
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

func setupConnection() {
	//ctx := canvas.Call("getContext", "2d")
	//img := ctx.Call("getImageData", 0, 0, imgWidth, imgHeight)

	//if img.Get("data").Length() != len(finalImage.Pix) {
	//	throw(errors.New("data size of images do not match"))
	//}

	//document := js.Global.Get("document")
	//location := document.Get("location")

	//ws, err := websocket.Dial(fmt.Sprintf("ws://%s/game", location.Get("host")))
	ws, err := websocket.Dial("ws://localhost/game")
	assert(err)

	go func() {
		enc := json.NewEncoder(ws)
		for {
			err := enc.Encode("12345")
			assert(err)

			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		dec := json.NewDecoder(ws)
		for {
			tok, err := dec.Token()
			assert(err)

			switch v := tok.(type) {
			case string:
				log.Println(v)
			default:
				log.Println("Unknown type.")
			}
		}
	}()
}

func updateTitle() {
	js.Global.Get("document").Set("title", "Fantasim")
}

func load() {
	document := js.Global.Get("document")

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

	setupConnection()
}

func main() {
	js.Global.Call("addEventListener", "load", func() { go load() })
}
