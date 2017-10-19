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

func startGame(conn io.ReadWriter) error {
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	if err := enc.Encode(&playerKey); err != nil {
		return nil
	}

	var version string
	if err := dec.Decode(&version); err != nil {
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

	for range time.Tick(time.Second / 30) {
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

		treeBgColor := color.RGBA{R: 93, G: 184, B: 155, A: 0xFF}
		waterColor := color.RGBA{R: 255, G: 119, B: 15, A: 0xFF}
		waterBgColor := color.RGBA{R: 255, G: 215, B: 15, A: 0xFF}

		treeColor := func(x, y int) color.RGBA {
			if noise.Eval2(float64(x), float64(y)) > 0 {
				return color.RGBA{R: 85, G: 138, B: 0, A: 0xFF}
			}
			return color.RGBA{R: 85, G: 100, B: 0, A: 0xFF}
		}

		waterTile := func(x, y int) *image.Paletted {
			if noise.Eval2(float64(x), float64(y)) > 0 {
				return tileReg["water"]
			}
			return tileReg["water2"]
		}

		grassBgColor := func(x, y int) color.RGBA {
			if noise.Eval2(float64(x), float64(y)) > 0 {
				return color.RGBA{R: 90, G: 178, B: 155, A: 0xFF}
			}
			return treeBgColor
		}

		snowBgColor := func(x, y int) color.RGBA {
			if noise.Eval2(float64(x), float64(y)) > 0 {
				return color.RGBA{R: 255, G: 236, B: 179, A: 0xFF}
			}
			return color.RGBA{R: 240, G: 230, B: 170, A: 0xFF}
		}

		for y := 0; y < cvr.H; y++ {
			for x := 0; x < cvr.W; x++ {
				tileData := rvresp.Data[y*cvr.W+x]
				wx := x + cameraPos.X
				wy := y + cameraPos.Y

				var (
					tile   *image.Paletted
					fg, bg color.RGBA
				)

				f := tileData.Flags
				switch {
				case f.Is(api.Water):
					tile = waterTile(wx, wy)
					fg = waterColor
					bg = waterBgColor
				case f.Is(api.Tree):
					tile = tileReg["tree"]
					fg = treeColor(wx, wy)
					bg = treeBgColor
				case f.Is(api.Bush):
					tile = tileReg["bush"]
					fg = treeColor(wx, wy)
					bg = treeBgColor
				case f.Is(api.Plant):
					tile = tileReg["plant"]
					fg = treeColor(wx, wy)
					bg = treeBgColor
				case f.Is(api.Stone):
					tile = tileReg["stone"]
					fg = color.RGBA{R: 128, G: 128, B: 128, A: 0xFF}
					bg = grassBgColor(wx, wy)
				case f.Is(api.Snow):
					tile = tileReg["none"]
					fg = snowBgColor(wx, wy)
				default:
					tile = tileReg["none"]
					fg = grassBgColor(wx, wy)
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

	ws, err := websocket.Dial(fmt.Sprintf("ws://%s/api", "localhost"), "", "http://localhost/")
	if err != nil {
		log.Fatalln(err)
	}
	defer ws.Close()

	if err := startGame(ws); err != nil {
		log.Fatalln(err)
	}
}
