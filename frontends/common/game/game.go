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
	"strings"
	"time"

	"github.com/andreas-jonsson/fantasim-pub/api"
	"github.com/andreas-jonsson/fantasim-pub/frontends/common/data"
	sys "github.com/andreas-jonsson/fantasim-pub/frontends/common/platform"
	"github.com/ojrac/opensimplex-go"
)

const (
	fantasimVersionString = "0.0.1"
)

var (
	GameFps    = 30
	RequestFps = 10
)

var (
	startTime = time.Now()
	noise     = opensimplex.NewWithSeed(time.Now().UnixNano())
)

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

			offset := dst.PixOffset(x, y)
			dst.Pix[offset] = col.R
			dst.Pix[offset+1] = col.G
			dst.Pix[offset+2] = col.B
			dst.Pix[offset+3] = col.A
			sx++
		}
		sy++
	}
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

func unitTile(u api.UnitViewData) *image.Paletted {
	tileReg := tilesetRegister["tiles"]
	asciiReg := tilesetRegister["ascii"]
	miscReg := tilesetRegister["misc"]

	switch u.Race {
	case api.Human:
		return asciiReg["h"]
	case api.Dwarf:
		return asciiReg["d"]
	case api.Goblin:
		return miscReg["goblin"]
	case api.Orc:
		return asciiReg["o"]
	case api.Troll:
		return asciiReg["T"]
	case api.Elven:
		return asciiReg["e"]
	case api.Deamon:
		return tileReg["deamon"]

	// Wildlife
	case api.Dear:
		return asciiReg["d"]
	case api.Boar:
		return asciiReg["b"]
	case api.Wolf:
		return asciiReg["w"]
	default:
		return asciiReg["?"]
	}
}

func itemTile(it api.ItemClass, bg color.RGBA) (*image.Paletted, color.RGBA, color.RGBA) {
	tileReg := tilesetRegister["tiles"]
	asciiReg := tilesetRegister["ascii"]
	miscReg := tilesetRegister["misc"]
	variousReg := tilesetRegister["various"]

	fg := color.RGBA{R: 139, G: 69, B: 19, A: 0xFF}
	red := color.RGBA{R: 200, G: 0, B: 0, A: 0xFF}
	white := color.RGBA{R: 220, G: 220, B: 220, A: 0xFF}

	switch it {
	case api.PartialItem:
		return asciiReg["&"], color.RGBA{R: 129, G: 129, B: 220, A: 0xFF}, bg
	case api.LogItem:
		return asciiReg["-"], fg, bg
	case api.FirewoodItem:
		return miscReg["firewood"], fg, bg
	case api.PlankItem:
		return asciiReg["="], fg, bg
	case api.StoneItem:
		return tileReg["stone"], color.RGBA{R: 128, G: 128, B: 128, A: 0xFF}, bg
	case api.MeatItem:
		return variousReg["meat"], color.RGBA{R: 210, G: 128, B: 128, A: 0xFF}, bg
	case api.BonesItem:
		return variousReg["bone"], color.RGBA{R: 200, G: 200, B: 200, A: 0xFF}, bg
	case api.SeedsItem:
		return miscReg["seeds"], color.RGBA{R: 190, G: 180, B: 19, A: 0xFF}, bg
	case api.CropItem:
		return miscReg["crop"], color.RGBA{R: 100, G: 190, B: 100, A: 0xFF}, bg
	case api.HumanCorpseItem:
		return asciiReg["h"], white, red
	case api.DwarfCorpseItem:
		return asciiReg["d"], white, red
	case api.GoblinCorpseItem:
		return miscReg["goblin"], white, red
	case api.OrcCorpseItem:
		return asciiReg["o"], white, red
	case api.TrollCorpseItem:
		return asciiReg["T"], white, red
	case api.ElvenCorpseItem:
		return asciiReg["e"], white, red
	case api.DeamonCorpseItem:
		return tileReg["deamon"], white, red
	case api.DearCorpseItem:
		return asciiReg["d"], white, red
	case api.BoarCorpseItem:
		return asciiReg["b"], white, red
	case api.WolfCorpseItem:
		return asciiReg["w"], white, red
	default:
		return asciiReg["?"], fg, bg
	}
}

var (
	renderer sys.Renderer
	input    sys.Input
)

func Initialize(in sys.Input, out sys.Renderer) error {
	input = in
	renderer = out
	return buildTilesets()
}

var logLines []string

func glog(s ...interface{}) {
	logLines = append(logLines, fmt.Sprint(s...))
}

func glogf(f string, s ...interface{}) {
	glog(fmt.Sprintf(f, s...))
}

func alert() {
	textBgColor = color.RGBA{R: 0xFF, A: 0xFF}
	go func() {
		fmt.Print("\a")
	}()
}

var (
	defaultTextFgColor = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	textFgColor        = defaultTextFgColor
	defaultTextBgColor = color.RGBA{A: 0xFF}
	textBgColor        = defaultTextBgColor
)

var buildingColors = []color.RGBA{
	{90, 39, 41, 255},
	{118, 57, 49, 255},
	{145, 85, 77, 255},
	{126, 46, 31, 255},
	{152, 80, 60, 255},
	{165, 113, 78, 255},
	{133, 87, 35, 255},
	{185, 156, 107, 255},
	{213, 196, 161, 255},
	{87, 65, 47, 255},
	{121, 96, 76, 255},
	{171, 149, 132, 255},
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

	sz := renderer.Resolution()
	backBuffer := image.NewRGBA(image.Rectangle{Max: sz})

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

	gameMain := func() error {
		viewportSize := image.Pt(sz.X/16, sz.Y/16)
		cameraPos := image.Pt(0, 0)
		lastReadView := time.Now().Add(-time.Second)
		lastHighlightRequest := lastReadView
		requestedCameraPos := cameraPos
		responseCameraPos := cameraPos
		mousePos := sz.Div(2)
		logUpdated := time.Now()
		logSize := 0
		renderHeight := false

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

		for range time.Tick(time.Second / time.Duration(GameFps)) {
			for ev := input.PollEvent(); ev != nil; ev = input.PollEvent() {
				switch t := ev.(type) {
				case *sys.QuitEvent:
					return nil
				case *sys.KeyboardEvent:
					if t.Type != sys.KeyboardDown {
						continue
					}

					switch {
					case t.Key == sys.KeyEsc:
						if hasAnyTool() {
							resetAllTools()
						} else if logWindow != nil || ctrlWindow != nil {
							logWindow = nil
							ctrlWindow = nil
						} else {
							return nil
						}
					case t.Name == "f11":
						if b, err := renderer.ToggleFullscreen(); err != nil {
							fmt.Println("Could not toggle fullscreen:", err)
						} else {
							fmt.Println("Toggle fullscreen:", b)
						}
					case t.Name == "left shift" || t.Name == "right shift":
						renderHeight = !renderHeight
					case t.Name == "a" && t.IsMod(sys.KeyModCtrl):
						cameraPos.X -= viewportSize.X
					case t.Name == "d" && t.IsMod(sys.KeyModCtrl):
						cameraPos.X += viewportSize.X
					case t.Name == "w" && t.IsMod(sys.KeyModCtrl):
						cameraPos.Y -= viewportSize.Y
					case t.Name == "s" && t.IsMod(sys.KeyModCtrl):
						cameraPos.Y += viewportSize.Y
					case t.Name == "l":
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
					case t.Key == sys.KeySpace:
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
							if cb := updateCtrlWindow(t); cb != nil {
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
				case *sys.MouseMotionEvent:
					if t.X >= sz.X {
						mousePos.X = sz.X - 1
					} else {
						mousePos.X = t.X
					}
					if t.Y >= sz.Y {
						mousePos.Y = sz.Y - 1
					} else {
						mousePos.Y = t.Y
					}
				case *sys.MouseButtonEvent:
					if t.Type == sys.MouseButtonDown {
						switch t.Button {
						case 1:
							pX, pY := float64(mousePos.X)/float64(sz.X), float64(mousePos.Y)/float64(sz.Y)
							mX, mY := float64(viewportSize.X)*pX, float64(viewportSize.Y)*pY
							mp := image.Pt(mousePos.X/8, mousePos.Y/16)
							mouseWorldPos := cameraPos.Add(image.Pt(int(mX), int(mY)))

							if pickTool != nil {
								if err := pickTool(enc, mp, mouseWorldPos, viewportSize, rvresp); err != nil {
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
											r.Name,
											"",
											"Race: " + u.Race.String(),
											fmt.Sprintf("Unit ID: %v", u.ID),
											fmt.Sprintf("Health:  %v", r.Health),
											fmt.Sprintf("Thirst:  %v", r.Thirst),
											fmt.Sprintf("Hunger:  %v", r.Hunger),
											"",
										}
										for _, s := range r.Debug {
											contextMenuText = append(contextMenuText, s)
										}
									}
								} else {
									if t.Building != api.InvalidID {
										contextMenuText = append(contextMenuText, "Building: "+t.BuildingType.String())
									}

									if t.StructureType != api.NoStructure {
										contextMenuText = append(contextMenuText, fmt.Sprintf("Structure: %s %s", t.StructureMaterial, t.StructureType))
									}

									switch {
									case t.Flags.Is(api.Water):
										contextMenuText = append(contextMenuText, "Tile: Water")
									case t.Flags.Is(api.Sand):
										contextMenuText = append(contextMenuText, "Tile: Sand")
									case t.Flags.Is(api.Snow):
										contextMenuText = append(contextMenuText, "Tile: Snow")
									case t.Flags.Is(api.Brook):
										contextMenuText = append(contextMenuText, "Tile: Brook")
									default:
										contextMenuText = append(contextMenuText, "Tile: Grass")
									}

									switch {
									case t.Flags.Is(api.Tree):
										contextMenuText = append(contextMenuText, "Object: Tree")
									case t.Flags.Is(api.Bush):
										contextMenuText = append(contextMenuText, "Object: Bush")
									case t.Flags.Is(api.Plant):
										contextMenuText = append(contextMenuText, "Object: Plant")
									case t.Flags.Is(api.Stone):
										contextMenuText = append(contextMenuText, "Object: Stone")
									}

									contextMenuText = append(contextMenuText, fmt.Sprintf("Height: %v", t.Height))
								}

								if len(t.Items) > 0 {
									s := "Item(s):"
									for _, it := range t.Items {
										s += " " + it.Class.String()
									}
									contextMenuText = append(contextMenuText, s)
								}
							}
						}
					} else {
						if t.Button == 1 && areaTool != nil {
							r := image.Rect(areaToolStart.X, areaToolStart.Y, mousePos.X/8, mousePos.Y/16)
							if err := areaTool(enc, r, cameraPos, viewportSize, rvresp); err != nil {
								return err
							}
						}

						contextMenu = nil
						contextMenuText = nil
					}
				}
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

			if updateHighlights || time.Since(lastHighlightRequest) > time.Second*2 {
				lastHighlightRequest = time.Now()
				updateHighlights = false

				id, err := encodeRequest(enc, &api.HighlightRequest{})
				if err != nil {
					return err
				}

				resp, err := decodeResponse(id)
				if err != nil {
					return err
				}
				highlights = resp.(*api.HighlightResponse).Highlight
			}

			if readViewRequestID == invalidRequestID && time.Since(lastReadView) >= time.Second/time.Duration(RequestFps) {
				lastReadView = time.Now()
				requestedCameraPos = cameraPos
				uvr := api.UpdateViewRequest{ViewID: viewID, X: cameraPos.X, Y: cameraPos.Y}

				id, err := encodeRequest(enc, &uvr)
				if err != nil {
					return err
				}
				discardResponse(id)

				rvr := api.ReadViewRequest{ViewID: viewID, RLE: true}
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
				render(backBuffer, &cvr, rvresp, responseCameraPos, cameraPos, renderHeight)
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
			logLen := len(logLines)

			if logLen > logSize {
				logUpdated = time.Now()
				logSize = logLen
			}

			if logWindow != nil {
				logWindow.clear()

				n := 1
				for y := logWindow.canvas.Size().Y - 1; y >= 0; y-- {
					i := logLen - n
					n++
					if i < 0 {
						continue
					}
					logWindow.print(0, y, logLines[i])
				}
			} else if logLen > 0 && time.Since(logUpdated) < time.Second*5 {
				print(2, sz.Y-1, " "+logLines[logLen-1]+" ")
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

			renderer.Present(backBuffer)

			textFgColor = defaultTextFgColor
			textBgColor = defaultTextBgColor
		}

		return nil
	}

	if err := gameMain(); err != nil {
		textFgColor = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
		textBgColor = color.RGBA{A: 0xFF}

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
	}

	return nil
}
