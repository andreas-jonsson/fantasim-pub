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
	"fmt"
	"image"
	"strings"

	sys "github.com/andreas-jonsson/fantasim-pub/frontends/common/platform"

	"github.com/veandco/go-sdl2/sdl"
)

var keyMapping = map[sdl.Keycode]int{
	sdl.K_UP:        sys.KeyUp,
	sdl.K_DOWN:      sys.KeyDown,
	sdl.K_LEFT:      sys.KeyLeft,
	sdl.K_RIGHT:     sys.KeyRight,
	sdl.K_ESCAPE:    sys.KeyEsc,
	sdl.K_SPACE:     sys.KeySpace,
	sdl.K_BACKSPACE: sys.KeyBackSpace,
}

var mouseMapping = map[int]int{
	sdl.MOUSEBUTTONDOWN: sys.MouseButtonDown,
	sdl.MOUSEBUTTONUP:   sys.MouseButtonUp,
	sdl.MOUSEWHEEL:      sys.MouseWheel,
}

type SDL struct {
	window     *sdl.Window
	backBuffer *sdl.Texture
	renderer   *sdl.Renderer
}

const fullscreenMode = sdl.WINDOW_FULLSCREEN_DESKTOP

func InitSDL(windowSize, resolution image.Point, fullscreen bool) (*SDL, error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, err
	}

	flags := uint32(sdl.WINDOW_SHOWN)
	if fullscreen {
		flags |= fullscreenMode
	}

	window, err := sdl.CreateWindow("Fantasim SDL", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(windowSize.X), int32(windowSize.Y), flags)
	if err != nil {
		sdl.Quit()
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		sdl.Quit()
		window.Destroy()
		return nil, err
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")
	renderer.SetLogicalSize(int32(resolution.X), int32(resolution.Y))

	backBuffer, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(resolution.X), int32(resolution.Y))
	if err != nil {
		sdl.Quit()
		window.Destroy()
		renderer.Destroy()
		return nil, err
	}

	sdl.ShowCursor(0)
	return &SDL{window, backBuffer, renderer}, nil
}

func (s *SDL) Quit() {
	s.backBuffer.Destroy()
	s.renderer.Destroy()
	s.window.Destroy()
	sdl.Quit()
}

func (s *SDL) Present(img image.Image) error {
	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		return fmt.Errorf("invalid image format: %T", img)
	}

	s.backBuffer.Update(nil, rgbaImg.Pix, rgbaImg.Stride)
	s.renderer.Copy(s.backBuffer, nil, nil)
	s.renderer.Present()
	return nil
}

func (s *SDL) ToggleFullscreen() (bool, error) {
	isFullscreen := (s.window.GetFlags() & sdl.WINDOW_FULLSCREEN) != 0
	if isFullscreen {
		s.window.SetFullscreen(0)
	} else {
		s.window.SetFullscreen(fullscreenMode)
	}
	return !isFullscreen, nil
}

func (s *SDL) Resolution() image.Point {
	w, h := s.renderer.GetLogicalSize()
	return image.Pt(int(w), int(h))
}

func (*SDL) PollEvent() sys.Event {
	event := sdl.PollEvent()
	if event == nil {
		return nil
	}

	switch t := event.(type) {
	case *sdl.QuitEvent:
		return &sys.QuitEvent{}
	case *sdl.KeyboardEvent:
		ev := &sys.KeyboardEvent{}
		if key, ok := keyMapping[t.Keysym.Sym]; ok {
			ev.Key = key
		} else {
			ev.Key = sys.KeyUnknown
		}

		ev.Mod = int(t.Keysym.Mod)
		ev.Name = strings.ToLower(sdl.GetKeyName(t.Keysym.Sym))

		if t.Type == sdl.KEYUP {
			ev.Type = sys.KeyboardUp
		} else if t.Type == sdl.KEYDOWN {
			ev.Type = sys.KeyboardDown
		} else {
			return &sys.UnknownEvent{}
		}
		return ev
	case *sdl.MouseButtonEvent:
		ev := &sys.MouseButtonEvent{}
		ev.Button = int(t.Button)
		ev.X = int(t.X)
		ev.Y = int(t.Y)

		switch t.Type {
		case sdl.MOUSEBUTTONDOWN:
			ev.Type = sys.MouseButtonDown
		case sdl.MOUSEBUTTONUP:
			ev.Type = sys.MouseButtonUp
		case sdl.MOUSEWHEEL:
			ev.Type = sys.MouseWheel
		default:
			return &sys.UnknownEvent{}
		}
		return ev
	case *sdl.MouseMotionEvent:
		ev := &sys.MouseMotionEvent{}
		ev.X = int(t.X)
		ev.Y = int(t.Y)
		ev.XRel = int(t.XRel)
		ev.YRel = int(t.YRel)
		return ev
	default:
		return &sys.UnknownEvent{}
	}

	return nil
}
