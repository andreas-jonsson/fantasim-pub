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

//go:generate go run ../common/data/generate.go

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/url"
	"strings"
	"time"

	"github.com/andreas-jonsson/fantasim-pub/api"
	"github.com/andreas-jonsson/fantasim-pub/frontends/default/data"
	"github.com/andreas-jonsson/fantasim-pub/frontends/default/renderer"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
)

const fantasimVersionString = "0.0.1"

const (
	imgWidth  = 320
	imgHeight = 200
	imgScale  = 1
)

var idCounter uint64

var (
	keys   = map[int]bool{}
	canvas *js.Object

	playerKey, playerName string
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

func newId() uint64 {
	idCounter++
	return idCounter
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
		Source   string         `json:"source"`
		TileSize int            `json:"tile_size"`
		Mapping  map[string]int `json:"mapping"`
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
		tileSize := tileset.TileSize
		numTiles := sz.X / tileSize

		set := make(map[string]*image.Paletted)
		tilesetRegister[setName] = set

		for tileName, offset := range tileset.Mapping {
			y := offset / numTiles
			x := offset % numTiles

			x *= tileSize
			y *= tileSize

			set[tileName] = img.SubImage(image.Rect(x, y, x+16, y+16)).(*image.Paletted)
		}
	}

	return nil
}

func setupConnection(host string, renderer *renderer.Renderer) {
	ws, err := websocket.Dial(fmt.Sprintf("ws://%s/api", host))
	assert(err)
	defer ws.Close()

	sz := renderer.BackBuffer().Bounds().Size()
	viewportSize := image.Pt(sz.X/16, sz.Y/16)
	cameraPos := image.Pt(0, 0)

	enc := json.NewEncoder(ws)
	dec := json.NewDecoder(ws)

	assert(enc.Encode(&playerKey))

	var version string
	assert(dec.Decode(&version))
	if version != fantasimVersionString {
		throw(fmt.Errorf("invalid version %s, expected %s", version, fantasimVersionString))
	}

	cvr := api.CreateViewRequest{
		X: cameraPos.X,
		Y: cameraPos.Y,
		W: cameraPos.X + viewportSize.X,
		H: cameraPos.Y + viewportSize.Y,
	}

	assert(api.EncodeRequest(enc, &cvr, 0))
	obj, _, err := api.DecodeResponse(dec)
	assert(err)

	viewID := obj.(*api.CreateViewResponse).ViewID

	for range time.Tick(time.Second / 15) {

		cameraPos.X++
		cameraPos.Y++

		uvr := api.UpdateViewRequest{ViewID: viewID, X: cameraPos.X, Y: cameraPos.Y}
		assert(api.EncodeRequest(enc, &uvr, 0))
		_, _, err := api.DecodeResponse(dec)
		assert(err)

		rvr := api.ReadViewRequest{ViewID: viewID}
		assert(api.EncodeRequest(enc, &rvr, 0))
		obj, _, err := api.DecodeResponse(dec)
		assert(err)

		rvresp := obj.(*api.ReadViewResponse)
		tileReg := tilesetRegister["tiles"]

		treeColor := color.RGBA{G: 0xFF, A: 0xFF}
		waterColor := color.RGBA{B: 0xFF, A: 0xFF}

		renderer.Clear()

		for y := 0; y < cvr.H; y++ {
			for x := 0; x < cvr.W; x++ {
				tileData := rvresp.Data[y*cvr.W+x]

				var (
					tile   *image.Paletted
					fg, bg color.RGBA
				)

				switch tileData.Surface {
				case "water":
					tile = tileReg["water"]
					fg = waterColor
					bg = color.RGBA{B: tileData.Height, A: 0xFF}
				case "tree":
					tile = tileReg["tree"]
					fg = treeColor
					bg = color.RGBA{G: tileData.Height, A: 0xFF}
				case "grass":
					tile = tileReg["grass"]
					fg = treeColor
					bg = color.RGBA{G: tileData.Height, A: 0xFF}
				default:
					continue
				}

				renderer.Blit(image.Pt(x*16, y*16), tile, fg, bg)
			}
		}

		renderer.Present()
	}
}

func updateTitle() {
	js.Global.Get("document").Set("title", "Fantasim")
}

func load() {
	document := js.Global.Get("document")

	location := js.Global.Get("location")
	urlStr := strings.TrimPrefix(location.Get("search").String(), "?")
	v, err := url.ParseQuery(urlStr)
	assert(err)

	playerKey = v.Get("key")
	playerName = v.Get("name")

	if playerKey == "" || playerName == "" {
		throw(errors.New("invalid parameters"))
	}

	document.Set("onkeydown", func(e *js.Object) {
		keys[e.Get("keyCode").Int()] = true
	})

	document.Set("onkeyup", func(e *js.Object) {
		keys[e.Get("keyCode").Int()] = false
	})

	renderer := renderer.NewRenderer(640, 360)

	assert(buildTilesets())
	setupConnection(location.Get("host").String(), renderer)
}

func main() {
	js.Global.Call("addEventListener", "load", func() { go load() })
}
