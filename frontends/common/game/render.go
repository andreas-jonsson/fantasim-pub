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
	"log"
	"time"

	"github.com/andreas-jonsson/fantasim-pub/api"
)

func render(backBuffer *image.RGBA, cvr *api.CreateViewRequest, rvresp *api.ReadViewResponse, cameraPos, currentCameraPos image.Point) error {
	tileReg := tilesetRegister["tiles"]

	treeBgColor := color.RGBA{R: 155, G: 184, B: 93, A: 0xFF}
	stoneColor := color.RGBA{R: 128, G: 128, B: 128, A: 0xFF}
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
				fg = stoneColor
			}

			if tileData.Building != api.InvalidID {
				tile = tileReg["floor"]
				fg = buildingColors[int(tileData.BuildingType)%len(buildingColors)]
				bg = buildingColors[int(^tileData.BuildingType)%len(buildingColors)]
			}

			if tile == nil {
				log.Fatalln("Could not load tile!")
			}

			n := int(noise.Eval2(float64(x), float64(y)) * 10)
			ticks500 := int(time.Since(startTime)/(time.Second/2)) + n
			ticks1000 := int(time.Since(startTime)/time.Second) + n
			numItems := len(tileData.Items)
			showUnit := ticks500%5 != 0 || numItems == 0
			numUnits := len(tileData.Units)

			switch {
			case numUnits > 0 && showUnit:
				if ticks1000 < 0 {
					ticks1000 = ticks1000 * -1
				}

				unit := tileData.Units[ticks1000%numUnits]
				fg := color.RGBA{R: 0xFF, A: 0xFF}

				switch unit.Allegiance {
				case api.Friendly:
					fg = color.RGBA{G: 0xFF, A: 0xFF}
				case api.Neutral:
					fg = color.RGBA{R: 0xFF, G: 0xFF, A: 0xFF}
				}

				blitImage(backBuffer, dp, unitTile(unit), fg, bg)
			case numItems > 1:
				fg = color.RGBA{R: 189, G: 129, B: 60, A: 0xFF}
				tile = tileReg["crate"]
				blitImage(backBuffer, dp, tile, fg, bg)
			case numItems > 0:
				tile, fg, bg = itemTile(tileData.Items[0].Class, bg)
				blitImage(backBuffer, dp, tile, fg, bg)
			default:
				blitImage(backBuffer, dp, tile, fg, bg)
			}
		}
	}

	return nil
}
