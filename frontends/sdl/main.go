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

//go:generate go run ../common/data/generate.go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/andreas-jonsson/vsdl-go"
	"github.com/ojrac/opensimplex-go"

	"github.com/andreas-jonsson/fantasim-pub/api"
	"github.com/andreas-jonsson/fantasim-pub/frontends/common/data"
	"golang.org/x/net/websocket"
)

const fantasimVersionString = "0.0.1"

const (
	imgWidth  = 320
	imgHeight = 200
	imgScale  = 1
)

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

func startGame(apiConn io.ReadWriter, infoConn io.Reader) error {
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

		const screenEdge = 10
		if mousePos.X < screenEdge {
			cameraPos.X--
		} else if mousePos.X > sz.X-screenEdge {
			cameraPos.X++
		}
		if mousePos.Y < screenEdge {
			cameraPos.Y--
		} else if mousePos.Y > sz.Y-screenEdge {
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

		vsdl.Present(backBuffer)
	}

	return nil
}

var playerKey, playerName string

func main() {
	flag.StringVar(&playerKey, "key", os.Getenv("fantasim-key"), "The player key assigned by the server.")
	flag.StringVar(&playerName, "name", "Unknown", "The name of the player.")
	flag.Parse()

	if err := buildTilesets(); err != nil {
		log.Fatalln(err)
	}

	apiWs, err := websocket.Dial(fmt.Sprintf("ws://%s/api", "localhost"), "", "http://localhost/")
	if err != nil {
		log.Fatalln(err)
	}
	defer apiWs.Close()

	infoWs, err := websocket.Dial(fmt.Sprintf("ws://%s/info", "localhost"), "", "http://localhost/")
	if err != nil {
		log.Fatalln(err)
	}
	defer infoWs.Close()

	enc := json.NewEncoder(io.MultiWriter(apiWs, infoWs))
	if err := enc.Encode(&playerKey); err != nil {
		log.Fatalln(err)
	}

	if err := startGame(apiWs, infoWs); err != nil {
		log.Fatalln(err)
	}
}
