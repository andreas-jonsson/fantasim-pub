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
	areaTool      func(api.Encoder, image.Rectangle, image.Point, image.Point, *api.ReadViewResponse) error

	pickTool func(api.Encoder, image.Point, image.Point, image.Point, *api.ReadViewResponse) error

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

func designateTreeCutting(enc api.Encoder) error {
	pickTool = func(enc api.Encoder, p, _, _ image.Point, _ *api.ReadViewResponse) error {
		pickTool = nil
		areaToolStart = p

		areaTool = func(enc api.Encoder, r image.Rectangle, camPos, vp image.Point, resp *api.ReadViewResponse) error {
			defer resetAllTools()
			if resp == nil {
				return nil
			}

			r.Min.X /= 2
			r.Max.X /= 2

			var trees []api.CutTreeData

			for y := r.Min.Y; y <= r.Max.Y; y++ {
				for x := r.Min.X; x <= r.Max.X; x++ {
					if resp.Data[y*vp.X+x].Flags.Is(api.Tree) {
						trees = append(trees, api.CutTreeData{camPos.X + x, camPos.Y + y})
					}
				}
			}

			if len(trees) == 0 {
				return nil
			}

			id, err := encodeRequest(enc, &api.CutTreesRequest{trees})
			if err != nil {
				return err
			}
			discardResponse(id)

			glogf("Cut %d tree(s) in area: %v", len(trees), r)
			return nil
		}
		return nil
	}
	return nil
}

func buildStockpile(enc api.Encoder) error {
	pickTool = func(enc api.Encoder, p, _, _ image.Point, _ *api.ReadViewResponse) error {
		pickTool = nil
		areaToolStart = p

		areaTool = func(enc api.Encoder, r image.Rectangle, camPos, _ image.Point, _ *api.ReadViewResponse) error {
			defer resetAllTools()

			r = r.Add(camPos)
			id, err := encodeRequest(enc, &api.BuildRequest{
				Building: api.StockpileBuilding,
				Location: api.Rect{Min: api.Point{r.Min.X, r.Min.Y}, Max: api.Point{r.Max.X, r.Max.Y}},
			})
			if err != nil {
				return err
			}

			resp, err := decodeResponse(id)
			if err != nil {
				return err
			}

			buildResp := resp.(*api.BuildResponse)
			if buildResp.Error == "" {
				glogf("Stockpile is planed for area: %v", r)
				return nil
			}

			glogf("Could not build stockpile: %v", buildResp.Error)
			return nil
		}
		return nil
	}
	return nil
}

func exploreLocation(enc api.Encoder) error {
	pickTool = func(enc api.Encoder, p, wp, _ image.Point, _ *api.ReadViewResponse) error {
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

func cameraToHomeLocation(_ api.Encoder) error {
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

func printJobQueue(_ api.Encoder) error {
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
