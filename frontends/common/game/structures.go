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

	"github.com/andreas-jonsson/fantasim-pub/api"
)

func materialColor(mat api.ItemClass) color.RGBA {
	return color.RGBA{R: 139, G: 69, B: 19, A: 0xFF}
}

func structureTile(ty api.StructureType, mat api.ItemClass, bg color.RGBA) (*image.Paletted, color.RGBA, color.RGBA) {
	asciiReg := tilesetRegister["ascii"]
	//patternsReg := tilesetRegister["patterns"]

	switch ty {
	case api.WallStructure:
		return asciiReg["#"], materialColor(mat), bg
	default:
		return asciiReg["?"], color.RGBA{R: 139, G: 69, B: 19, A: 0xFF}, bg
	}
}
