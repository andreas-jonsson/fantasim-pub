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
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"time"

	"github.com/andreas-jonsson/fantasim-pub/api"
	"github.com/andreas-jonsson/fantasim-pub/frontends/common/data"
	sys "github.com/andreas-jonsson/fantasim-pub/frontends/common/platform"
	"github.com/andreas-jonsson/fantasim-pub/frontends/common/raycast"
)

const (
	fantasimVersionString = "0.0.1"
)

var (
	GameFps    = 30
	RequestFps = 10
)

var startTime = time.Now()

func decodePNG(name string) (*image.Paletted, error) {
	fp, err := data.FS.Open(name)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	img, err := png.Decode(fp)
	if err != nil {
		return nil, err
	}

	pimg, ok := img.(*image.Paletted)
	if !ok {
		return nil, fmt.Errorf("%s was not a paletted image", name)
	}

	return pimg, nil
}

var tilesetRegister = make(map[string]map[string]*image.Paletted)

func buildTilesets() error {
	fp, err := data.FS.Open("tiles.json")
	if err != nil {
		return err
	}
	defer fp.Close()

	var tilesets map[string]struct {
		Source     string         `json:"source"`
		TileWidth  int            `json:"tile_width"`
		TileHeight int            `json:"tile_height"`
		Mapping    map[string]int `json:"mapping"`
	}

	dec := json.NewDecoder(fp)
	if err := dec.Decode(&tilesets); err != nil {
		return err
	}

	for setName, tileset := range tilesets {
		img, err := decodePNG(tileset.Source)
		if err != nil {
			return err
		}

		sz := img.Bounds().Size()
		numW := sz.X / tileset.TileWidth

		set := make(map[string]*image.Paletted)
		tilesetRegister[setName] = set

		if len(tileset.Mapping) == 0 {
			tileset.Mapping = make(map[string]int)
			i := 0
			for ; i < 128; i++ {
				tileset.Mapping[string(i)] = i
			}
			for ; i < 256; i++ {
				tileset.Mapping[fmt.Sprintf("#%d", i)] = i
			}
		}

		for tileName, offset := range tileset.Mapping {
			y := offset / numW
			x := offset % numW

			x *= tileset.TileWidth
			y *= tileset.TileHeight

			set[tileName] = img.SubImage(image.Rect(x, y, x+tileset.TileWidth, y+tileset.TileHeight)).(*image.Paletted)
		}
	}

	return nil
}

func blitImage(dst *image.RGBA, dp image.Point, src *image.Paletted, fg, bg color.RGBA) {
	sr := src.Bounds()
	dr := dst.Bounds().Intersect(image.Rectangle{dp, dp.Add(sr.Size())})
	pal := [2]color.RGBA{bg, fg}

	sy := sr.Min.Y
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		sx := sr.Min.X
		for x := dr.Min.X; x < dr.Max.X; x++ {
			i := src.PixOffset(sx, sy)
			col := pal[src.Pix[i]]
			if col.A > 0 {
				offset := dst.PixOffset(x, y)
				dst.Pix[offset] = col.R
				dst.Pix[offset+1] = col.G
				dst.Pix[offset+2] = col.B
				dst.Pix[offset+3] = col.A
			}
			sx++
		}
		sy++
	}
}

var (
	renderer sys.Renderer
	input    sys.Input
)

func Initialize(in sys.Input, out sys.Renderer) error {
	input = in
	renderer = out

	w, err := newWorld()
	if err != nil {
		log.Fatalln(err)
	}

	sprites, err := loadSprites(level)
	if err != nil {
		log.Fatalln(err)
	}

	rc := raycast.NewRaycaster(rt, w)
	sc := raycast.NewSpritecaster(sprites)

	return buildTilesets()
}

func Start(enc api.Encoder, dec, decInfo api.Decoder) error {
	var version string
	if err := decInfo.Decode(&version); err != nil {
		return err
	}

	if version != fantasimVersionString {
		return fmt.Errorf("invalid version %s, expected %s", version, fantasimVersionString)
	}

	var gameType api.GameType
	if err := decInfo.Decode(&gameType); err != nil {
		return err
	}

	//startAsyncDecoder(dec)
	//startInfoDecode(decInfo)

	sz := renderer.Resolution()
	backBuffer := image.NewRGBA(image.Rectangle{Max: sz})
	textFgColor := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	textBgColor := color.RGBA{A: 0xFF}

	putch := func(x, y int, ch string) {
		src, ok := tilesetRegister["default"][ch]
		if !ok {
			src = tilesetRegister["default"]["?"]
		}
		blitImage(backBuffer, image.Pt(x*8, y*16), src, textFgColor, textBgColor)
	}

	print := func(x, y int, text string) {
		for i, r := range text {
			putch(x+i, y, string(r))
		}
	}

	//if err := gameMain(); err != nil {
	for range time.Tick(time.Second / time.Duration(GameFps)) {
		for ev := input.PollEvent(); ev != nil; ev = input.PollEvent() {
			switch ev.(type) {
			case *sys.QuitEvent, *sys.KeyboardEvent, *sys.MouseButtonEvent:
				return nil
			}
		}

		draw.Draw(backBuffer, backBuffer.Bounds(), image.NewUniform(textBgColor), image.ZP, draw.Over)

		pos := image.Pt(sz.X/8, sz.Y/16)
		const msg = "Server was disconnected..."
		print(pos.X/2-len(msg)/2, pos.Y/2, msg)

		renderer.Present(backBuffer)
	}
	//}

	return nil
}
