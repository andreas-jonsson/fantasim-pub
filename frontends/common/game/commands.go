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

	"github.com/andreas-jonsson/fantasim-pub/api"
)

var (
	areaToolStart image.Point
	areaTool      func(api.Encoder, image.Rectangle) error

	pickTool func(api.Encoder, image.Point, image.Point) error

	moveCameraTool func(enc api.Encoder, viewID int, cameraPos *image.Point) error
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

func designateTreeCutting(enc api.Encoder, _ image.Point) error {
	pickTool = func(enc api.Encoder, p, wp image.Point) error {
		pickTool = nil
		areaToolStart = p

		areaTool = func(enc api.Encoder, r image.Rectangle) error {
			defer resetAllTools()

			// TODO: Implement this.

			glogf("Cut trees: %v", r)
			return nil
		}
		return nil
	}
	return nil
}

func designatePile(enc api.Encoder, _ image.Point) error {
	pickTool = func(enc api.Encoder, p, wp image.Point) error {
		pickTool = nil
		areaToolStart = p

		areaTool = func(enc api.Encoder, r image.Rectangle) error {
			defer resetAllTools()

			// TODO: Implement this.

			glogf("Pile in this area: %v", r)
			return nil
		}
		return nil
	}
	return nil
}

func exploreLocation(enc api.Encoder, _ image.Point) error {
	pickTool = func(enc api.Encoder, p, wp image.Point) error {
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

func cameraToHomeLocation(enc api.Encoder, _ image.Point) error {
	moveCameraTool = func(enc api.Encoder, viewID int, cameraPos *image.Point) error {
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

func printJobQueue(enc api.Encoder, _ image.Point) error {
	moveCameraTool = func(enc api.Encoder, viewID int, cameraPos *image.Point) error {
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
