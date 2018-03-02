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

package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const (
	srcFolder  = "../../common/data/src/tiles"
	destFilder = "../../common/data/src"
)

type tileSpec struct {
	numTiles, tileSize image.Point
	border, space      int
}

var tileMaps = map[string]tileSpec{
	"16x16_accesories_books_diamons_hair.png": {image.Pt(16, 13), image.Pt(16, 16), 3, 6},
	"16x16_arrows_pointers_aim.png":           {image.Pt(16, 12), image.Pt(16, 16), 3, 6},
	"16x16_misc.png":                          {image.Pt(16, 16), image.Pt(16, 16), 3, 6},
	"16x16_tiles_ditherPatterns.png":          {image.Pt(16, 11), image.Pt(16, 16), 3, 6},
	"16x16_various.png":                       {image.Pt(16, 15), image.Pt(16, 16), 3, 6},
	"16x16_weapons.png":                       {image.Pt(16, 16), image.Pt(16, 16), 3, 6},
}

func main() {
	for filename, tilemap := range tileMaps {
		srcFile := path.Join(srcFolder, filename)

		data, err := ioutil.ReadFile(srcFile)
		if err != nil {
			log.Fatalln(err)
		}

		src, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			log.Fatalln(err)
		}

		pal := color.Palette{color.RGBA{A: 0xFF}, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}}
		dest := image.NewPaletted(image.Rectangle{Max: image.Pt(tilemap.numTiles.X*tilemap.tileSize.X, tilemap.numTiles.Y*tilemap.tileSize.Y)}, pal)

		for yt := 0; yt < tilemap.numTiles.Y; yt++ {
			for xt := 0; xt < tilemap.numTiles.X; xt++ {
				srcPos := image.Pt(tilemap.border+xt*tilemap.tileSize.X+xt*tilemap.space, tilemap.border+yt*tilemap.tileSize.Y+yt*tilemap.space)
				destPos := image.Pt(xt*tilemap.tileSize.X, yt*tilemap.tileSize.Y)
				destRect := image.Rectangle{Min: destPos, Max: destPos.Add(tilemap.tileSize)}

				draw.Draw(dest, destRect, src, srcPos, draw.Over)
			}
		}

		fp, err := os.Create(path.Join(destFilder, filename))
		if err != nil {
			log.Fatalln(err)
		}

		if err := png.Encode(fp, dest); err != nil {
			log.Fatalln(err)
		}

		fp.Close()
	}
}
