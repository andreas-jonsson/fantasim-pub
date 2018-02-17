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
	"strings"
	"sync/atomic"
	"time"

	"github.com/andreas-jonsson/fantasim-pub/api"
	"github.com/andreas-jonsson/fantasim-pub/frontends/common/data"
	"github.com/andreas-jonsson/vsdl-go"
	"github.com/ojrac/opensimplex-go"
)

const fantasimVersionString = "0.0.1"

var noise = opensimplex.NewWithSeed(time.Now().UnixNano())

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

func updateLogWithServerInfo(lines []string) []string {
	for {
		select {
		case s, ok := <-infoChan:
			if !ok {
				return lines
			}
			lines = append(lines, s)
		default:
			return lines
		}
	}
}

func itemClassToString(it api.ItemClass) string {
	switch it {
	case api.Log:
		return "Log"
	default:
		return "Unidentified"
	}
}

func Initialize() error {
	return buildTilesets()
}

func update(backBuffer *image.RGBA, cvr *api.CreateViewRequest, rvresp *api.ReadViewResponse, cameraPos, currentCameraPos image.Point) error {
	tileReg := tilesetRegister["tiles"]
	asciiReg := tilesetRegister["ascii"]

	treeBgColor := color.RGBA{R: 155, G: 184, B: 93, A: 0xFF}
	brookColor := color.RGBA{R: 45, G: 169, B: 220, A: 0xFF}
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

	cameraOffset := currentCameraPos.Sub(cameraPos)

	for y := 0; y < cvr.H; y++ {
		for x := 0; x < cvr.W; x++ {

			tileData := index(x, y)
			wx := x + cameraPos.X
			wy := y + cameraPos.Y
			dp := image.Pt((x-cameraOffset.X)*16, (y-cameraOffset.Y)*16)

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
			case f.Is(api.Brook):
				tile = tileReg["none"]
				fg = brookColor
				bg = fg
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

			switch {
			case len(tileData.Units) > 0:
				unit := tileData.Units[0]
				fg := color.RGBA{R: 0xFF, A: 0xFF}

				switch unit.Allegiance {
				case api.Friendly:
					fg = color.RGBA{G: 0xFF, A: 0xFF}
				case api.Neutral:
					fg = color.RGBA{R: 0xFF, G: 0xFF, A: 0xFF}
				}

				blitImage(backBuffer, dp, tileReg["deamon"], fg, bg)
			case len(tileData.Items) > 0:
				fg := color.RGBA{R: 139, G: 69, B: 19, A: 0xFF}
				blitImage(backBuffer, dp, asciiReg["-"], fg, bg)
			default:
				blitImage(backBuffer, dp, tile, fg, bg)
			}
		}
	}

	return nil
}

var logLines []string

func glog(s ...interface{}) {
	logLines = append(logLines, fmt.Sprint(s...))
}

func glogf(f string, s ...interface{}) {
	glog(fmt.Sprintf(f, s...))
}

func Start(enc api.Encoder, dec, decInfo api.Decoder) error {
	var version string
	if err := decInfo.Decode(&version); err != nil {
		return err
	}

	if version != fantasimVersionString {
		return fmt.Errorf("invalid version %s, expected %s", version, fantasimVersionString)
	}

	startAsyncDecoder(dec)
	startInfoDecode(decInfo)

	sz := image.Pt(1280, 720)
	backBuffer := image.NewRGBA(image.Rectangle{Max: sz})

	sdlMain := func() error {
		viewportSize := image.Pt(sz.X/16, sz.Y/16)
		cameraPos := image.Pt(0, 0)
		requestedCameraPos := cameraPos
		responseCameraPos := cameraPos
		mousePos := sz.Div(2)

		cvr := api.CreateViewRequest{
			X: cameraPos.X,
			Y: cameraPos.Y,
			W: cameraPos.X + viewportSize.X,
			H: cameraPos.Y + viewportSize.Y,
		}

		id, err := encodeRequest(enc, &cvr)
		if err != nil {
			return err
		}

		obj, err := decodeResponse(id)
		if err != nil {
			return err
		}

		viewID := obj.(*api.CreateViewResponse).ViewID

		id, err = encodeRequest(enc, &api.ViewHomeRequest{viewID})
		if err != nil {
			return err
		}

		if resp, err := decodeResponse(id); err != nil {
			return err
		} else {
			r := resp.(*api.ViewHomeResponse)
			cameraPos = image.Pt(r.X, r.Y)
		}

		putch := func(x, y int, ch string) {
			blitImage(backBuffer, image.Pt(x*8, y*16), tilesetRegister["default"][ch], color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, color.RGBA{A: 0xFF})
		}

		print := func(x, y int, text string) {
			for i, r := range text {
				putch(x+i, y, string(r))
			}
		}

		var (
			readViewRequestID = invalidRequestID
			rvresp            *api.ReadViewResponse
		)

		var (
			ctrlWindowText,
			contextMenuText []string
		)

		var (
			contextMenu,
			logWindow,
			ctrlWindow *window
		)

		for range time.Tick(time.Second / 15) {
			for ev := range vsdl.Events() {
				switch t := ev.(type) {
				case *vsdl.QuitEvent:
					return nil
				case *vsdl.KeyDownEvent:
					switch {
					case t.Keysym.IsKey(vsdl.EscKey):
						if hasAnyTool() {
							resetAllTools()
						} else if logWindow != nil || ctrlWindow != nil {
							logWindow = nil
							ctrlWindow = nil
						} else {
							return nil
						}
					case t.Keysym.Sym == vsdl.Keycode('a') && t.Keysym.IsMod(vsdl.CtrlMod):
						cameraPos.X -= viewportSize.X
					case t.Keysym.Sym == vsdl.Keycode('d') && t.Keysym.IsMod(vsdl.CtrlMod):
						cameraPos.X += viewportSize.X
					case t.Keysym.Sym == vsdl.Keycode('w') && t.Keysym.IsMod(vsdl.CtrlMod):
						cameraPos.Y -= viewportSize.Y
					case t.Keysym.Sym == vsdl.Keycode('s') && t.Keysym.IsMod(vsdl.CtrlMod):
						cameraPos.Y += viewportSize.Y
					case t.Keysym.Sym == vsdl.Keycode('l'):
						if logWindow == nil {
							logWindow = newWindow(
								"Log",
								image.Rect(2, viewportSize.Y-12, viewportSize.X*2-2, viewportSize.Y-2),
								tilesetRegister["default"],
								putch,
							)
						} else {
							logWindow = nil
						}
					case t.Keysym.IsKey(vsdl.SpaceKey):
						if ctrlWindow == nil {
							resetAllTools()
							resetMenuWindow()

							ctrlWindow = newWindow(
								"Menu",
								image.Rect(viewportSize.X*2-48, 2, viewportSize.X*2-2, viewportSize.Y-2),
								tilesetRegister["default"],
								putch,
							)
						} else {
							ctrlWindow = nil
						}
					default:
						if ctrlWindow != nil {
							if cb := updateCtrlWindow(t.Keysym); cb != nil {
								resetMenuWindow()
								ctrlWindow = nil

								if err := cb(enc); err != nil {
									return nil
								}

								if moveCameraTool != nil {
									if err := moveCameraTool(enc, viewID, &cameraPos); err != nil {
										return err
									}
								}
							}
						}
					}
				case *vsdl.MouseMotionEvent:
					mousePos = image.Pt(int(t.X), int(t.Y))
				case *vsdl.MouseButtonEvent:
					if t.State == 1 {
						switch t.Button {
						case 1:
							pX, pY := float64(mousePos.X)/float64(sz.X), float64(mousePos.Y)/float64(sz.Y)
							mX, mY := float64(viewportSize.X)*pX, float64(viewportSize.Y)*pY
							mp := image.Pt(mousePos.X/8, mousePos.Y/16)
							mouseWorldPos := cameraPos.Add(image.Pt(int(mX), int(mY)))

							switch {
							case pickTool != nil:
								if err := pickTool(enc, mp, mouseWorldPos, cameraPos, rvresp); err != nil {
									return err
								}
							case areaTool != nil:
								r := image.Rect(areaToolStart.X, areaToolStart.Y, mousePos.X/8, mousePos.Y/16)
								if err := areaTool(enc, r, cameraPos, viewportSize, rvresp); err != nil {
									return err
								}
							}
						case 3:
							p := image.Pt(mousePos.X/8, mousePos.Y/16).Add(image.Pt(2, 1))
							contextMenu = newWindow(
								"Info",
								image.Rectangle{Min: p, Max: p.Add(image.Pt(32, 16))},
								tilesetRegister["default"],
								putch,
							)

							if rvresp != nil {
								p := image.Pt(mousePos.X/16, mousePos.Y/16)
								i := viewportSize.X*p.Y + p.X
								t := rvresp.Data[i]
								units := t.Units

								if len(units) > 0 {
									u := units[0]

									id, err := encodeRequest(enc, &api.UnitStatsRequest{int(u.ID)})
									if err != nil {
										return err
									}

									if resp, err := decodeResponse(id); err != nil {
										return err
									} else {
										r := resp.(*api.UnitStatsResponse)
										contextMenuText = []string{
											fmt.Sprintf("Unit ID: %v", u.ID),
											fmt.Sprintf("Health:  %v", r.Health),
											fmt.Sprintf("Thirst:  %v", r.Thirst),
											"",
										}
										for _, s := range r.Debug {
											contextMenuText = append(contextMenuText, s)
										}
									}
								} else {
									switch {
									case t.Flags.Is(api.Water):
										contextMenuText = append(contextMenuText, "Tile: Water")
									case t.Flags.Is(api.Sand):
										contextMenuText = append(contextMenuText, "Tile: Sand")
									case t.Flags.Is(api.Snow):
										contextMenuText = append(contextMenuText, "Tile: Snow")
									default:
										contextMenuText = append(contextMenuText, "Tile: Grass")
									}

									switch {
									case t.Flags.Is(api.Tree):
										contextMenuText = append(contextMenuText, "Object: Tree")
									case t.Flags.Is(api.Bush):
										contextMenuText = append(contextMenuText, "Object: Bush")
									case t.Flags.Is(api.Bush):
										contextMenuText = append(contextMenuText, "Object: Plant")
									case t.Flags.Is(api.Bush):
										contextMenuText = append(contextMenuText, "Object: Stone")
									}

									contextMenuText = append(contextMenuText, fmt.Sprintf("Height: %v", t.Height))
								}

								if len(t.Items) > 0 {
									s := "Item(s):"
									for _, it := range t.Items {
										s += " " + itemClassToString(it.Class)
									}
									contextMenuText = append(contextMenuText, s)
								}
							}
						}
					} else {
						contextMenu = nil
						contextMenuText = nil
					}
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

			if readViewRequestID == invalidRequestID {
				requestedCameraPos = cameraPos
				uvr := api.UpdateViewRequest{ViewID: viewID, X: cameraPos.X, Y: cameraPos.Y}

				id, err := encodeRequest(enc, &uvr)
				if err != nil {
					return err
				}
				discardResponse(id)

				rvr := api.ReadViewRequest{ViewID: viewID}
				readViewRequestID, err = encodeRequest(enc, &rvr)
				if err != nil {
					return err
				}
			}

			if obj, err := decodeResponseTimeout(readViewRequestID, time.Millisecond); err == nil {
				readViewRequestID = invalidRequestID
				responseCameraPos = requestedCameraPos
				rvresp = obj.(*api.ReadViewResponse)
			}

			if rvresp != nil {
				update(backBuffer, &cvr, rvresp, responseCameraPos, cameraPos)
			}

			if areaTool != nil {
				r := image.Rect(areaToolStart.X, areaToolStart.Y, mousePos.X/8, mousePos.Y/16)
				top := "+"
				szX := r.Size().X

				if szX > 0 {
					top = "+" + strings.Repeat("-", szX-1) + "+"
				}
				print(r.Min.X, r.Min.Y, top)

				for y := r.Min.Y + 1; y < r.Max.Y; y++ {
					print(r.Min.X, y, "|")
					print(r.Max.X, y, "|")
				}
				print(r.Min.X, r.Max.Y, top)
			}

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

			title := fmt.Sprintf(" Fantasim - [%d:%d] ", cameraPos.X, cameraPos.Y)
			print(sz.X/2-len(title)/2, 0, title)

			logLines = updateLogWithServerInfo(logLines)

			if logWindow != nil {
				logWindow.clear()

				n := 1
				for y := logWindow.canvas.Size().Y - 1; y >= 0; y-- {
					i := len(logLines) - n
					n++
					if i < 0 {
						continue
					}
					logWindow.print(0, y, logLines[i])
				}
			}

			if ctrlWindow != nil {
				ctrlWindowText, ctrlWindow.title = updateCtrlWindowText(ctrlWindowText)
				ctrlWindow.clear()
				for i, line := range ctrlWindowText {
					ctrlWindow.print(0, i, line)
				}

				if len(menuStack) > 1 {
					ctrlWindow.print(0, ctrlWindow.canvas.Max.Y-4, " Backspace: Go Back")
				}
			}

			if contextMenu != nil {
				contextMenu.clear()
				for i, s := range contextMenuText {
					contextMenu.print(0, i, s)
				}
			}

			switch {
			case pickTool != nil:
				putch(mousePos.X/8, mousePos.Y/16, "X")
			default:
				putch(mousePos.X/8, mousePos.Y/16, "#219")
			}

			vsdl.Present(backBuffer)
		}

		return nil
	}

	return vsdl.Initialize(sdlMain, vsdl.ConfigWithRenderer(sz, image.ZP))
}
