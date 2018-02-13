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

	"github.com/andreas-jonsson/fantasim-pub/api"
)

var (
	areaToolStart image.Point
	areaTool      func(*json.Encoder, image.Rectangle, image.Point, image.Point, *api.ReadViewResponse) error

	pickTool func(*json.Encoder, image.Point, image.Point, image.Point, *api.ReadViewResponse) error

	moveCameraTool func(enc *json.Encoder, viewID int, cameraPos *image.Point) error
)

func resetAllTools() {
	areaToolStart = image.ZP
	areaTool = nil
	pickTool = nil
	moveCameraTool = nil
}

func hasAnyTool() bool {
	if areaTool != nil || pickTool != nil || moveCameraTool != nil {
		return true
	}
	return false
}

func designateTreeCutting(_ *json.Encoder) error {
	pickTool = func(enc *json.Encoder, p, wp, camPos image.Point, _ *api.ReadViewResponse) error {
		pickTool = nil
		areaToolStart = p

		areaTool = func(enc *json.Encoder, r image.Rectangle, camPos, vp image.Point, resp *api.ReadViewResponse) error {
			defer resetAllTools()
			if resp == nil {
				return nil
			}

			for y := 0; y < vp.Y; y++ {
				for x := 0; x < vp.X; x++ {
					if resp.Data[y*vp.X+x].Flags.Is(api.Tree) {
						fmt.Println("Treee!", camPos.X+x, camPos.Y+y)
					}
				}
			}

			glogf("Cut trees: %v", r)
			return nil
		}
		return nil
	}
	return nil
}

func exploreLocation(_ *json.Encoder) error {
	pickTool = func(enc *json.Encoder, p, wp, camPos image.Point, _ *api.ReadViewResponse) error {
		defer resetAllTools()

		id, err := encodeRequest(enc, &api.ExploreLocationRequest{wp.X, wp.Y})
		if err != nil {
			return err
		}
		discardResponse(id)

		glogf("Explore location: %d,%d", wp.X, wp.Y)
		return nil
	}
	return nil
}

func cameraToHomeLocation(_ *json.Encoder) error {
	moveCameraTool = func(enc *json.Encoder, viewID int, cameraPos *image.Point) error {
		defer resetAllTools()

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

func printJobQueue(_ *json.Encoder) error {
	moveCameraTool = func(enc *json.Encoder, viewID int, cameraPos *image.Point) error {
		defer resetAllTools()

		id, err := encodeRequest(enc, &api.JobQueueRequest{})
		if err != nil {
			return err
		}

		if resp, err := decodeResponse(id); err != nil {
			return err
		} else {
			r := resp.(*api.JobQueueResponse)
			glog("Job queue:")
			for _, s := range r.Jobs {
				glog(s)
			}
			return nil
		}
	}
	return nil
}
