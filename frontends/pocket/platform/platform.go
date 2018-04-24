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
	"image"
	"image/draw"

	"golang.org/x/mobile/event/size"

	sys "github.com/andreas-jonsson/fantasim-pub/frontends/common/platform"
	"golang.org/x/mobile/event/touch"
)

type Pocket struct{}

const maxEvents = 128

var (
	sizeEvent size.Event
	mousePosX, mousePosY,
	mouseBeginX, mouseBeginY float32
)

var (
	ExternalBackBuffer *image.RGBA
	PaintEventChan     = make(chan chan struct{})
	InputEventChan     = make(chan interface{}, maxEvents)
)

func InitPocket() (*Pocket, error) {
	return &Pocket{}, nil
}

func (p *Pocket) Quit() {
}

func (p *Pocket) Present(img image.Image) error {
	cb := <-PaintEventChan
	draw.Draw(ExternalBackBuffer, img.Bounds(), img, image.ZP, draw.Src)
	cb <- struct{}{}
	return nil
}

func (p *Pocket) ToggleFullscreen() (bool, error) {
	return false, nil
}

func (p *Pocket) Resolution() image.Point {
	return ExternalBackBuffer.Bounds().Size()
}

func (p *Pocket) PollEvent() sys.Event {
	select {
	case ev := <-InputEventChan:
		switch ev := ev.(type) {
		case nil:
			return &sys.QuitEvent{}
		case size.Event:
			ze := size.Event{}
			if sizeEvent == ze {
				sizeEvent = ev
			}
		case touch.Event:
			res := p.Resolution()
			switch ev.Type {
			case touch.TypeBegin:
				mouseBeginX = (ev.X / float32(sizeEvent.WidthPx)) * float32(res.X)
				mouseBeginY = (ev.Y / float32(sizeEvent.HeightPx)) * float32(res.Y)
			case touch.TypeMove:
				newPosX := (ev.X / float32(sizeEvent.WidthPx)) * float32(res.X)
				newPosY := (ev.Y / float32(sizeEvent.HeightPx)) * float32(res.Y)

				deltaX := newPosX - mouseBeginX
				deltaY := newPosY - mouseBeginY

				mouseBeginX = newPosX
				mouseBeginY = newPosY

				mousePosX += deltaX
				mousePosY += deltaY

				return &sys.MouseMotionEvent{int(mousePosX), int(mousePosY), int(deltaX), int(deltaY)}
			}
		}
	default:
	}
	return nil
}
