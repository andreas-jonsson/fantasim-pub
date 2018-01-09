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

package renderer

import (
	"image"
	"image/color"
	"image/draw"
	"runtime"
	"strconv"

	"github.com/gopherjs/gopherjs/js"
)

type Renderer struct {
	backBuffer    *image.RGBA
	context       *js.Object
	width, height int
}

func NewRenderer(w, h int) *Renderer {
	r := &Renderer{
		width:  w,
		height: h,
	}

	document := js.Global.Get("document")
	r.backBuffer = image.NewRGBA(image.Rect(0, 0, r.width, r.height))

	canvas := document.Call("createElement", "canvas")
	canvas.Call("setAttribute", "width", strconv.Itoa(r.width))
	canvas.Call("setAttribute", "height", strconv.Itoa(r.height))

	style := canvas.Get("style")
	style.Set("width", strconv.Itoa(w)+"px")
	style.Set("height", strconv.Itoa(h)+"px")
	style.Set("cursor", "none")

	document.Get("body").Call("appendChild", canvas)
	r.context = canvas.Call("getContext", "2d")
	//setupCanvasInput(canvas, w, h, r.width, r.height)

	return r
}

func (r *Renderer) Clear() {
	draw.Draw(r.backBuffer, r.backBuffer.Bounds(), &image.Uniform{color.RGBA{A: 0xFF}}, image.ZP, draw.Src)
}

func (r *Renderer) Present() {
	img := r.context.Call("getImageData", 0, 0, r.width, r.height)
	data := img.Get("data")

	arrBuf := js.Global.Get("ArrayBuffer").New(data.Length())
	buf8 := js.Global.Get("Uint8ClampedArray").New(arrBuf)
	buf32 := js.Global.Get("Uint32Array").New(arrBuf)

	buf := buf32.Interface().([]uint)
	pix := r.backBuffer.Pix

	for offset := 0; offset < len(pix); offset += 4 {
		buf[offset/4] = 0xFF000000 | (uint(pix[offset+2]) << 16) | (uint(pix[offset+1]) << 8) | uint(pix[offset])
	}

	data.Call("set", buf8)
	r.context.Call("putImageData", img, 0, 0)

	runtime.Gosched()
}

func (r *Renderer) BackBuffer() draw.Image {
	return r.backBuffer
}

func (r *Renderer) Blit(dp image.Point, src *image.Paletted, fg, bg color.RGBA) {
	srcImg := *src
	srcImg.Palette = color.Palette{
		bg,
		fg,
	}

	sr := srcImg.Bounds()
	dr := image.Rectangle{dp, dp.Add(sr.Size())}
	draw.Draw(r.backBuffer, dr, &srcImg, sr.Min, draw.Over)
}
