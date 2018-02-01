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
	"image"

	"github.com/andreas-jonsson/fantasim-pub/api"
)

var (
	areaToolStart image.Point
	areaTool      func(*json.Encoder, image.Rectangle) error

	pickTool func(*json.Encoder, image.Point) error

	moveCameraTool func(enc *json.Encoder, viewID int, cameraPos *image.Point) error
)

func resetAllTools() {
	areaToolStart = image.ZP
	areaTool = nil
	pickTool = nil
	moveCameraTool = nil
}

func designateTreeCutting(enc *json.Encoder) error {
	return nil
}

func exploreLocation(enc *json.Encoder) error {
	pickTool = func(enc *json.Encoder, p image.Point) error {
		id, err := encodeRequest(enc, &api.ExploreLocationRequest{p.X, p.Y})
		if err != nil {
			return err
		}
		discardResponse(id)

		glogf("Explore location: %d,%d", p.X, p.Y)
		return nil
	}
	return nil
}

func cameraToHomeLocation(enc *json.Encoder) error {
	moveCameraTool = func(enc *json.Encoder, viewID int, cameraPos *image.Point) error {
		id, err := encodeRequest(enc, &api.ViewHomeRequest{viewID})
		if err != nil {
			return err
		}

		if resp, err := decodeResponse(id); err != nil {
			return err
		} else {
			r := resp.(*api.ViewHomeResponse)
			*cameraPos = image.Pt(r.X, r.Y)
			glogf("Jump to home location: %d,%d", r.X, r.Y)
			return nil
		}
	}
	return nil
}
