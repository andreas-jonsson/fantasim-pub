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
	context       js.Value
	events        chan sys.Event
	width, height int
	mousePos      image.Point
}

func InitWASM(sz image.Point, scale int) *WASM {
	s := &WASM{
		events: make(chan sys.Event, 64),
		width:  sz.X,
		height: sz.Y,
	}

	document := js.Global().Get("document")
	canvas := document.Call("createElement", "canvas")
	canvas.Call("setAttribute", "width", strconv.Itoa(s.width))
	canvas.Call("setAttribute", "height", strconv.Itoa(s.height))
	canvas.Set("imageSmoothingEnabled", false)
	canvas.Set("oncontextmenu", js.NewEventCallback(js.PreventDefault, func(js.Value) {}))

	style := canvas.Get("style")
	style.Set("width", strconv.Itoa(sz.X*scale)+"px")
	style.Set("height", strconv.Itoa(sz.Y*scale)+"px")
	style.Set("cursor", "none")

	document.Get("body").Call("appendChild", canvas)
	s.context = canvas.Call("getContext", "2d")

	document.Set("onkeydown", js.NewEventCallback(js.PreventDefault, func(e js.Value) {
		key := e.Get("key").String()
		select {
		case s.events <- &sys.KeyboardEvent{
			Key:  keyMapping[key],
			Type: sys.KeyboardDown,
			Name: key,
		}:
		default:
		}
	}))

	document.Set("onkeyup", js.NewEventCallback(js.PreventDefault, func(e js.Value) {
		key := e.Get("key").String()
		select {
		case s.events <- &sys.KeyboardEvent{
			Key:  keyMapping[key],
			Type: sys.KeyboardUp,
			Name: key,
		}:
		default:
		}
	}))

	document.Set("onmousemove", js.NewEventCallback(js.PreventDefault, func(e js.Value) {
		s.mousePos.X = e.Get("clientX").Int()
		s.mousePos.Y = e.Get("clientY").Int()

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
	}))

	canvas.Set("onmousedown", js.NewEventCallback(js.PreventDefault, func(e js.Value) {
		select {
		case s.events <- &sys.MouseButtonEvent{
			Type:   sys.MouseButtonDown,
			X:      s.mousePos.X,
			Y:      s.mousePos.Y,
			Button: e.Get("button").Int() + 1,
		}:
		default:
		}
	}))

	canvas.Set("onmouseup", js.NewEventCallback(js.PreventDefault, func(e js.Value) {
		select {
		case s.events <- &sys.MouseButtonEvent{
			Type:   sys.MouseButtonUp,
			X:      s.mousePos.X,
			Y:      s.mousePos.Y,
			Button: e.Get("button").Int() + 1,
		}:
		default:
		}
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

func (s *WASM) ToggleFullscreen() (bool, error) {
	return false, errors.New("not supported")
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
