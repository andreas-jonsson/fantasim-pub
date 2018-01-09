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

package game

import (
	"image"
	"image/color"
)

type window struct {
	title        string
	putch        func(int, int, string)
	tileset      map[string]*image.Paletted
	fg, bg       color.RGBA
	rect, canvas image.Rectangle
}

func newWindow(title string, rect image.Rectangle, tileset map[string]*image.Paletted, putch func(int, int, string)) *window {
	one := image.Pt(1, 1)
	return &window{
		title,
		putch,
		tileset,
		color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
		color.RGBA{A: 0xFF},
		rect,
		image.Rectangle{Min: rect.Add(one).Min, Max: rect.Sub(one).Max},
	}
}

func (w *window) clear() {
	for y := w.canvas.Min.Y; y < w.canvas.Max.Y; y++ {
		for x := w.canvas.Min.X; x < w.canvas.Max.X; x++ {
			w.putch(x, y, string(" "))
		}
	}

	x, y := w.rect.Min.X, w.rect.Min.Y
	sz := w.rect.Size()

	for i := 1; i < sz.X-1; i++ {
		w.putch(x+i, y, "#196")
		w.putch(x+i, y+(sz.Y-1), "#196")
	}

	w.putch(x, y, "#218")
	w.putch(x+(sz.X-1), y, "#191")
	w.putch(x, y+(sz.Y-1), "#192")
	w.putch(x+(sz.X-1), y+(sz.Y-1), "#217")

	for i := 1; i < sz.Y-1; i++ {
		w.putch(x, y+i, "#179")
		w.putch(x+(sz.X-1), y+i, "#179")
	}

	c := sz.X/2 - len(w.title)/2
	for i, r := range w.title {
		w.putch(x+c+i, y, string(r))
	}
}

func (w *window) print(x, y int, text string) {
	x += w.canvas.Min.X
	y += w.canvas.Min.Y

	for i, r := range text {
		p := x + i
		if image.Pt(p, y).In(w.canvas) {
			w.putch(p, y, string(r))
		}
	}
}
