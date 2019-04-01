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

package platform

import (
	"errors"
	"image"
	"log"
	"strconv"

	"syscall/js"

	sys "github.com/andreas-jonsson/fantasim-pub/frontends/common/platform"
)

var keyMapping = map[string]int{
	"Arrow Up":    sys.KeyUp,
	"Arrow Down":  sys.KeyDown,
	"Arrow Left":  sys.KeyLeft,
	"Arrow Right": sys.KeyRight,
	"Escape":      sys.KeyEsc,
	" ":           sys.KeySpace,
	"Backspace":   sys.KeyBackSpace,
}

type keyEvent struct {
	key  string
	down bool
}

type WASM struct {
	context, canvas js.Value
	events          chan sys.Event
	width, height   int
	mousePos        image.Point
}

func InitWASM(sz image.Point, scale int) *WASM {
	s := &WASM{
		events: make(chan sys.Event, 64),
		width:  sz.X,
		height: sz.Y,
	}

	document := js.Global().Get("document")
	body := document.Get("body")

	loadingText := document.Call("getElementById", "loadingText")
	body.Call("removeChild", loadingText)

	canvas := document.Call("createElement", "canvas")
	s.canvas = canvas

	canvas.Call("setAttribute", "width", strconv.Itoa(s.width))
	canvas.Call("setAttribute", "height", strconv.Itoa(s.height))
	canvas.Set("imageSmoothingEnabled", true)
	canvas.Set("oncontextmenu", js.FuncOf(func(js.Value, []js.Value) interface{} { return nil }))

	style := canvas.Get("style")
	//style.Set("width", strconv.Itoa(sz.X*scale)+"px")
	//style.Set("height", strconv.Itoa(sz.Y*scale)+"px")
	//style.Set("width", "100%")
	style.Set("height", "100%")
	style.Set("cursor", "none")

	body.Call("appendChild", canvas)
	s.context = canvas.Call("getContext", "2d")

	body.Call("appendChild", document.Call("createElement", "br"))

	fsButton := document.Call("createElement", "input")
	body.Call("appendChild", fsButton)

	fsButton.Set("type", "button")
	fsButton.Set("id", "fsButton")
	fsButton.Set("value", "Fullscreen")
	fsButton.Set("onclick", js.FuncOf(func(js.Value, []js.Value) interface{} {
		if _, err := s.ToggleFullscreen(); err != nil {
			log.Println(err)
		}
		return nil
	}))

	document.Set("onkeydown", js.FuncOf(func(_ js.Value, e []js.Value) interface{} {
		key := e[0].Get("key").String()
		select {
		case s.events <- &sys.KeyboardEvent{
			Key:  keyMapping[key],
			Type: sys.KeyboardDown,
			Name: key,
		}:
		default:
		}
		return nil
	}))

	document.Set("onkeyup", js.FuncOf(func(_ js.Value, e []js.Value) interface{} {
		key := e[0].Get("key").String()
		select {
		case s.events <- &sys.KeyboardEvent{
			Key:  keyMapping[key],
			Type: sys.KeyboardUp,
			Name: key,
		}:
		default:
		}
		return nil
	}))

	document.Set("onmousemove", js.FuncOf(func(_ js.Value, e []js.Value) interface{} {
		s.mousePos.X = e[0].Get("clientX").Int()
		s.mousePos.Y = e[0].Get("clientY").Int()

		if s.mousePos.X > s.width {
			s.mousePos.X = s.width
		}
		if s.mousePos.Y > s.height {
			s.mousePos.Y = s.height
		}

		select {
		case s.events <- &sys.MouseMotionEvent{
			X: s.mousePos.X,
			Y: s.mousePos.Y,
		}:
		default:
		}
		return nil
	}))

	canvas.Set("onmousedown", js.FuncOf(func(_ js.Value, e []js.Value) interface{} {
		select {
		case s.events <- &sys.MouseButtonEvent{
			Type:   sys.MouseButtonDown,
			X:      s.mousePos.X,
			Y:      s.mousePos.Y,
			Button: e[0].Get("button").Int() + 1,
		}:
		default:
		}
		return nil
	}))

	canvas.Set("onmouseup", js.FuncOf(func(_ js.Value, e []js.Value) interface{} {
		select {
		case s.events <- &sys.MouseButtonEvent{
			Type:   sys.MouseButtonUp,
			X:      s.mousePos.X,
			Y:      s.mousePos.Y,
			Button: e[0].Get("button").Int() + 1,
		}:
		default:
		}
		return nil
	}))

	return s
}

func (s *WASM) Quit() {
}

func (s *WASM) Present(screen image.Image) error {
	rgbImg, ok := screen.(*image.RGBA)
	if !ok {
		return errors.New("invalid image format")
	}

	img := s.context.Call("getImageData", 0, 0, s.width, s.height)
	data := img.Get("data")
	array := js.TypedArrayOf(rgbImg.Pix)

	data.Call("set", array)
	s.context.Call("putImageData", img, 0, 0)
	array.Release()

	return nil
}

// ToggleFullscreen is only partialy supported because it does not report the state back correctly.
func (s *WASM) ToggleFullscreen() (bool, error) {
	c := s.canvas
	if c.Get("webkitRequestFullScreen").Type() == js.TypeFunction {
		c.Call("webkitRequestFullScreen")
	} else if c.Get("mozRequestFullScreen").Type() == js.TypeFunction {
		c.Call("mozRequestFullScreen")
	} else {
		return false, errors.New("not supported")
	}
	return true, nil
}

func (s *WASM) Resolution() image.Point {
	return image.Pt(int(s.width), int(s.height))
}

func (s *WASM) PollEvent() sys.Event {
	select {
	case ev := <-s.events:
		return ev
	default:
		return nil
	}
}
