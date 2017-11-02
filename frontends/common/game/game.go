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

package game

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"sync/atomic"
	"time"

	"github.com/andreas-jonsson/fantasim-pub/api"
	"github.com/andreas-jonsson/fantasim-pub/frontends/common/data"
	"github.com/andreas-jonsson/vsdl-go"
	"github.com/ojrac/opensimplex-go"
)

const fantasimVersionString = "0.0.1"

var idCounter uint64

func newId() uint64 {
	return atomic.AddUint64(&idCounter, 1)
}

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
	srcImg := *src
	srcImg.Palette = color.Palette{
		bg,
		fg,
	}

	sr := srcImg.Bounds()
	dr := image.Rectangle{dp, dp.Add(sr.Size())}
	draw.Draw(dst, dr, &srcImg, sr.Min, draw.Over)
}

func Initialize() error {
	return buildTilesets()
}

func update(backBuffer *image.RGBA, cvr *api.CreateViewRequest, rvresp *api.ReadViewResponse, noise *opensimplex.Noise, cameraPos image.Point) error {
	tileReg := tilesetRegister["tiles"]

	treeBgColor := color.RGBA{R: 155, G: 184, B: 93, A: 0xFF}
	waterColor := color.RGBA{R: 15, G: 119, B: 255, A: 0xFF}
	waterBgColor := color.RGBA{R: 15, G: 215, B: 255, A: 0xFF}

	treeColor := func(x, y int) color.RGBA {
		if noise.Eval2(float64(x), float64(y)) > 0 {
			return color.RGBA{R: 0, G: 138, B: 85, A: 0xFF}
		}
		return color.RGBA{R: 0, G: 100, B: 85, A: 0xFF}
	}

	snowTreeColor := func(x, y int) color.RGBA {
		if noise.Eval2(float64(x), float64(y)) > 0 {
			return color.RGBA{R: 127, G: 137, B: 127, A: 0xFF}
		}
		return color.RGBA{R: 175, G: 185, B: 175, A: 0xFF}
	}

	waterTile := func(x, y int) *image.Paletted {
		if noise.Eval2(float64(x), float64(y)) > 0 {
			return tileReg["water"]
		}
		return tileReg["water2"]
	}

	grassBgColor := func(x, y int) color.RGBA {
		if noise.Eval2(float64(x), float64(y)) > 0 {
			return color.RGBA{R: 155, G: 178, B: 90, A: 0xFF}
		}
		return treeBgColor
	}

	snowBgColor := func(x, y int) color.RGBA {
		if noise.Eval2(float64(x), float64(y)) > 0 {
			return color.RGBA{R: 179, G: 236, B: 255, A: 0xFF}
		}
		return color.RGBA{R: 170, G: 230, B: 240, A: 0xFF}
	}

	sandBgColor := func(x, y int) color.RGBA {
		if noise.Eval2(float64(x), float64(y)) > 0 {
			return color.RGBA{R: 189, G: 183, B: 107, A: 0xFF}
		}
		return color.RGBA{R: 180, G: 175, B: 100, A: 0xFF}
	}

	index := func(x, y int) api.ReadViewData {
		return rvresp.Data[y*cvr.W+x]
	}

	for y := 0; y < cvr.H; y++ {
		for x := 0; x < cvr.W; x++ {
			tileData := index(x, y)
			wx := x + cameraPos.X
			wy := y + cameraPos.Y

			var (
				tile *image.Paletted
				fg   = color.RGBA{A: 0xFF}
				bg   = fg
			)

			f := tileData.Flags
			switch {
			case f.Is(api.Water):
				tile = waterTile(wx, wy)
				fg = waterColor
				bg = waterBgColor
			case f.Is(api.Snow):
				tile = tileReg["none"]
				fg = snowBgColor(wx, wy)
				bg = fg
			case f.Is(api.Sand):
				tile = tileReg["none"]
				fg = sandBgColor(wx, wy)
				bg = fg
			default:
				tile = tileReg["none"]
				fg = grassBgColor(wx, wy)
				bg = fg
			}

			switch {
			case f.Is(api.Tree):
				if f.Is(api.Sand) {
					tile = tileReg["palm"]
				} else if f.Is(api.Snow) {
					tile = tileReg["pine"]
				} else {
					tile = tileReg["tree"]
				}
				if f.Is(api.Snow) {
					fg = snowTreeColor(wx, wy)
				} else {
					fg = treeColor(wx, wy)
				}
			case f.Is(api.Bush):
				if f.Is(api.Sand) {
					tile = tileReg["cactus"]
				} else {
					tile = tileReg["bush"]
				}
				if f.Is(api.Snow) {
					fg = snowTreeColor(wx, wy)
				} else {
					fg = treeColor(wx, wy)
				}
			case f.Is(api.Plant):
				tile = tileReg["plant"]
				fg = treeColor(wx, wy)
			case f.Is(api.Stone):
				tile = tileReg["stone"]
				fg = color.RGBA{R: 128, G: 128, B: 128, A: 0xFF}
			}

			if tile == nil {
				log.Fatalln("Could not load tile!")
			}

			blitImage(backBuffer, image.Pt(x*16, y*16), tile, fg, bg)
		}
	}

	return nil
}

func Start(apiConn io.ReadWriter, infoConn io.Reader) error {
	enc := json.NewEncoder(apiConn)
	dec := json.NewDecoder(apiConn)
	decInfo := json.NewDecoder(infoConn)

	var version string
	if err := decInfo.Decode(&version); err != nil {
		return err
	}

	if version != fantasimVersionString {
		return fmt.Errorf("invalid version %s, expected %s", version, fantasimVersionString)
	}

	sz := image.Pt(1280, 720)
	backBuffer := image.NewRGBA(image.Rectangle{Max: sz})
	if err := vsdl.Initialize(vsdl.ConfigWithRenderer(sz, image.ZP)); err != nil {
		return err
	}
	defer vsdl.Shutdown()

	viewportSize := image.Pt(sz.X/16, sz.Y/16)
	cameraPos := image.Pt(0, 0)
	mousePos := sz.Div(2)

	cvr := api.CreateViewRequest{
		X: cameraPos.X,
		Y: cameraPos.Y,
		W: cameraPos.X + viewportSize.X,
		H: cameraPos.Y + viewportSize.Y,
	}

	if err := api.EncodeRequest(enc, &cvr, 0); err != nil {
		return err
	}

	obj, _, err := api.DecodeResponse(dec)
	if err != nil {
		return err
	}

	viewID := obj.(*api.CreateViewResponse).ViewID
	noise := opensimplex.NewWithSeed(time.Now().UnixNano())

	for range time.Tick(time.Second / 15) {
		for ev := range vsdl.Events() {
			switch t := ev.(type) {
			case *vsdl.QuitEvent:
				return nil
			case *vsdl.KeyDownEvent:
				if t.Keysym.IsKey(vsdl.EscKey) {
					return nil
				}

				log.Printf("%s: %X\n", t.Keysym, t.Keysym.Sym)
			case *vsdl.MouseMotionEvent:
				mousePos = image.Pt(int(t.X), int(t.Y))
			}
			ev.Release()
		}

		const (
			screenEdgeX = 8
			screenEdgeY = 16
		)

		if mousePos.X < screenEdgeX {
			cameraPos.X--
		} else if mousePos.X > sz.X-screenEdgeX {
			cameraPos.X++
		}
		if mousePos.Y < screenEdgeY {
			cameraPos.Y--
		} else if mousePos.Y > sz.Y-screenEdgeY {
			cameraPos.Y++
		}

		uvr := api.UpdateViewRequest{ViewID: viewID, X: cameraPos.X, Y: cameraPos.Y}
		if err := api.EncodeRequest(enc, &uvr, 0); err != nil {
			return err
		}
		_, _, err := api.DecodeResponse(dec)
		if err != nil {
			return err
		}

		rvr := api.ReadViewRequest{ViewID: viewID}
		if err := api.EncodeRequest(enc, &rvr, 0); err != nil {
			return err
		}
		obj, _, err := api.DecodeResponse(dec)
		if err != nil {
			return err
		}

		rvresp := obj.(*api.ReadViewResponse)

		putch := func(x, y int, ch string) {
			blitImage(backBuffer, image.Pt(x*8, y*16), tilesetRegister["default"][ch], color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, color.RGBA{A: 0xFF})
		}

		print := func(x, y int, text string) {
			for i, r := range text {
				putch(x+i, y, string(r))
			}
		}

		update(backBuffer, &cvr, rvresp, noise, cameraPos)

		sz := image.Pt(sz.X/8, sz.Y/16)

		for i := 1; i < sz.X-1; i++ {
			putch(i, 0, "#196")
			putch(i, sz.Y-1, "#196")
		}

		putch(0, 0, "#218")
		putch(sz.X-1, 0, "#191")
		putch(0, sz.Y-1, "#192")
		putch(sz.X-1, sz.Y-1, "#217")

		for i := 1; i < sz.Y-1; i++ {
			putch(0, i, "#179")
			putch(sz.X-1, i, "#179")
		}

		const title = " Fantasim "
		print(sz.X/2-len(title)/2, 0, title)

		vsdl.Present(backBuffer)
	}

	return nil
}
