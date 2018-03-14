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
	"runtime"
	"strconv"

	sys "github.com/andreas-jonsson/fantasim-pub/frontends/common/platform"
	"github.com/gopherjs/gopherjs/js"
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

type JS struct {
	context       *js.Object
	events        chan sys.Event
	width, height int
	mousePos      image.Point
}

func InitJS(sz image.Point) *JS {
	s := &JS{
		events: make(chan sys.Event, 64),
		width:  sz.X,
		height: sz.Y,
	}

	document := js.Global.Get("document")
	canvas := document.Call("createElement", "canvas")
	canvas.Call("setAttribute", "width", strconv.Itoa(s.width))
	canvas.Call("setAttribute", "height", strconv.Itoa(s.height))
	canvas.Set("imageSmoothingEnabled", false)
	canvas.Set("oncontextmenu", func(e *js.Object) {
		e.Call("preventDefault")
	})

	style := canvas.Get("style")
	style.Set("width", strconv.Itoa(sz.X*2)+"px")
	style.Set("height", strconv.Itoa(sz.Y*2)+"px")
	style.Set("cursor", "none")

	document.Get("body").Call("appendChild", canvas)
	s.context = canvas.Call("getContext", "2d")

	document.Set("onkeydown", func(e *js.Object) {
		key := e.Get("key").String()
		select {
		case s.events <- &sys.KeyboardEvent{
			Key:  keyMapping[key],
			Type: sys.KeyboardDown,
			Name: key,
		}:
		default:
		}
	})

	document.Set("onkeyup", func(e *js.Object) {
		key := e.Get("key").String()
		select {
		case s.events <- &sys.KeyboardEvent{
			Key:  keyMapping[key],
			Type: sys.KeyboardUp,
			Name: key,
		}:
		default:
		}
	})

	document.Set("onmousemove", func(e *js.Object) {
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
	})

	canvas.Set("onmousedown", func(e *js.Object) {
		select {
		case s.events <- &sys.MouseButtonEvent{
			Type:   sys.MouseButtonDown,
			X:      s.mousePos.X,
			Y:      s.mousePos.Y,
			Button: e.Get("button").Int() + 1,
		}:
		default:
		}
	})

	canvas.Set("onmouseup", func(e *js.Object) {
		select {
		case s.events <- &sys.MouseButtonEvent{
			Type:   sys.MouseButtonUp,
			X:      s.mousePos.X,
			Y:      s.mousePos.Y,
			Button: e.Get("button").Int() + 1,
		}:
		default:
		}
	})

	return s
}

func (s *JS) Quit() {
}

func (s *JS) Present(screen image.Image) error {
	rgbImg, ok := screen.(*image.RGBA)
	if !ok {
		return errors.New("invalid image format")
	}

	img := s.context.Call("getImageData", 0, 0, s.width, s.height)
	data := img.Get("data")

	arrBuf := js.Global.Get("ArrayBuffer").New(data.Length())
	buf8 := js.Global.Get("Uint8ClampedArray").New(arrBuf)
	buf32 := js.Global.Get("Uint32Array").New(arrBuf)

	buf := buf32.Interface().([]uint)
	pix := rgbImg.Pix

	for offset := 0; offset < len(pix); offset += 4 {
		buf[offset/4] = 0xFF000000 | (uint(pix[offset+2]) << 16) | (uint(pix[offset+1]) << 8) | uint(pix[offset])
	}

	data.Call("set", buf8)
	s.context.Call("putImageData", img, 0, 0)

	runtime.Gosched()
	return nil
}

func (s *JS) ToggleFullscreen() (bool, error) {
	return false, errors.New("not supported")
}

func (s *JS) Resolution() image.Point {
	return image.Pt(int(s.width), int(s.height))
}

func (s *JS) PollEvent() sys.Event {
	select {
	case ev := <-s.events:
		return ev
	default:
		return nil
	}
}
