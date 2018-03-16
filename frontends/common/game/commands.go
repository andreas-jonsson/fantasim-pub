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
	updateHighlights = true
}

func hasAnyTool() bool {
	if areaTool != nil || pickTool != nil || moveCameraTool != nil {
		return true
	}
	return false
}

func orderTreeCutting(enc api.Encoder) error {
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

			var trees []api.Point
			for y := r.Min.Y; y <= r.Max.Y; y++ {
				for x := r.Min.X; x <= r.Max.X; x++ {
					if resp.Data[y*vp.X+x].Flags.Is(api.Tree) {
						trees = append(trees, api.Point{camPos.X + x, camPos.Y + y})
					}
				}
			}

			if len(trees) == 0 {
				alert()
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

func collectItems(enc api.Encoder) error {
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

			var items []api.Point
			for y := r.Min.Y; y <= r.Max.Y; y++ {
				for x := r.Min.X; x <= r.Max.X; x++ {
					if len(resp.Data[y*vp.X+x].Items) > 0 {
						items = append(items, api.Point{camPos.X + x, camPos.Y + y})
					}
				}
			}

			if len(items) == 0 {
				alert()
				return nil
			}

			id, err := encodeRequest(enc, &api.CollectItemsRequest{items})
			if err != nil {
				return err
			}
			discardResponse(id)

			glogf("Collect %d items(s) in area: %v", len(items), r)
			return nil
		}
		return nil
	}
	return nil
}

func gatherSeeds(enc api.Encoder) error {
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

			var plants []api.Point
			for y := r.Min.Y; y <= r.Max.Y; y++ {
				for x := r.Min.X; x <= r.Max.X; x++ {
					if resp.Data[y*vp.X+x].Flags.Is(api.Plant) {
						plants = append(plants, api.Point{camPos.X + x, camPos.Y + y})
					}
				}
			}

			if len(plants) == 0 {
				alert()
				return nil
			}

			id, err := encodeRequest(enc, &api.GatherSeedsRequest{plants})
			if err != nil {
				return err
			}
			discardResponse(id)

			glogf("Geather seeds in area: %v", r)
			return nil
		}
		return nil
	}
	return nil
}

func buildStructure(enc api.Encoder, s api.StructureType, c api.ItemClass) error {
	pickTool = func(enc api.Encoder, p, _, _ image.Point, _ *api.ReadViewResponse) error {
		pickTool = nil
		areaToolStart = p

		areaTool = func(enc api.Encoder, r image.Rectangle, camPos, _ image.Point, _ *api.ReadViewResponse) error {
			defer resetAllTools()

			r.Min.X /= 2
			r.Max.X /= 2
			r = r.Add(camPos)

			id, err := encodeRequest(enc, &api.BuildRequest{
				Structure: s,
				Material:  c,
				Location:  api.Rect{Min: api.Point{r.Min.X, r.Min.Y}, Max: api.Point{r.Max.X + 1, r.Max.Y + 1}},
			})
			if err != nil {
				return err
			}
			discardResponse(id)

			glogf("Building structure: %v", r)
			return nil
		}
		return nil
	}
	return nil
}

func designateBuilding(enc api.Encoder, b api.BuildingType) error {
	pickTool = func(enc api.Encoder, p, _, _ image.Point, _ *api.ReadViewResponse) error {
		pickTool = nil
		areaToolStart = p

		areaTool = func(enc api.Encoder, r image.Rectangle, camPos, _ image.Point, _ *api.ReadViewResponse) error {
			defer resetAllTools()

			r.Min.X /= 2
			r.Max.X /= 2
			r = r.Add(camPos)

			sz := r.Size()
			if sz.X < 3 || sz.Y < 3 || sz.X > 30 || sz.Y > 30 {
				alert()
				return nil
			}

			id, err := encodeRequest(enc, &api.DesignateRequest{
				Building: b,
				Location: api.Rect{Min: api.Point{r.Min.X, r.Min.Y}, Max: api.Point{r.Max.X, r.Max.Y}},
			})
			if err != nil {
				return err
			}

			resp, err := decodeResponse(id)
			if err != nil {
				return err
			}

			buildResp := resp.(*api.DesignateResponse)
			if buildResp.Error == "" {
				glogf("%s is planed for area: %v", b, r)
				return nil
			}

			alert()
			glogf("Could not build %s: %v", b, buildResp.Error)
			return nil
		}
		return nil
	}
	return nil
}

func seedFarm(enc api.Encoder) error {
	pickTool = func(enc api.Encoder, p, wp, vp image.Point, rvresp *api.ReadViewResponse) error {
		defer resetAllTools()

		p.X /= 2
		tile := rvresp.Data[p.Y*vp.X+p.X]

		if tile.Building == api.InvalidID || tile.BuildingType != api.FarmBuilding {
			glog("No farm at location:", wp)
			return nil
		}

		id, err := encodeRequest(enc, &api.SeedFarmRequest{tile.Building})
		if err != nil {
			return err
		}
		discardResponse(id)
		return nil
	}
	return nil
}

func mineLocation(enc api.Encoder) error {
	pickTool = func(enc api.Encoder, p, wp, _ image.Point, _ *api.ReadViewResponse) error {
		defer resetAllTools()

		id, err := encodeRequest(enc, &api.MineLocationRequest{api.Point{wp.X, wp.Y}})
		if err != nil {
			return err
		}

		resp, err := decodeResponse(id)
		if err != nil {
			return err
		}

		mineResp := resp.(*api.MineLocationResponse)
		if mineResp.Error == "" {
			glog("Mining location:", wp)
			return nil
		}

		alert()
		glog("Could not mine location:", mineResp.Error)
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

func attackUnits(enc api.Encoder) error {
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

			var units []uint64
			for y := r.Min.Y; y <= r.Max.Y; y++ {
				for x := r.Min.X; x <= r.Max.X; x++ {
					if tu := resp.Data[y*vp.X+x].Units; len(tu) > 0 {
						for _, u := range tu {
							if u.Allegiance != api.Friendly {
								units = append(units, u.ID)
							}
						}
					}
				}
			}

			if len(units) == 0 {
				alert()
				return nil
			}

			id, err := encodeRequest(enc, &api.AttackUnitsRequest{units})
			if err != nil {
				return err
			}

			discardResponse(id)
			return nil
		}
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

func debugCommand(enc api.Encoder, cmd string) {
	moveCameraTool = func(enc api.Encoder, _ int, _ *image.Point) error {
		defer resetAllTools()

		id, err := encodeRequest(enc, &api.DebugCommandRequest{cmd})
		if err != nil {
			return err
		}

		if resp, err := decodeResponse(id); err != nil {
			return err
		} else {
			s := resp.(*api.DebugCommandResponse).Error
			if s != "" {
				glog(s)
			}
			return nil
		}
	}
}

func listResources(_ api.Encoder) error {
	moveCameraTool = func(enc api.Encoder, viewID int, cameraPos *image.Point) error {
		defer resetAllTools()

		id, err := encodeRequest(enc, &api.DebugCommandRequest{"resources"})
		if err != nil {
			return err
		}

		if resp, err := decodeResponse(id); err != nil {
			return err
		} else {
			glog("Resources:", resp.(*api.DebugCommandResponse).Error)
			return nil
		}
	}
	return nil
}

func listPlayers(_ api.Encoder) error {
	moveCameraTool = func(enc api.Encoder, viewID int, cameraPos *image.Point) error {
		defer resetAllTools()

		id, err := encodeRequest(enc, &api.DebugCommandRequest{"players"})
		if err != nil {
			return err
		}

		if resp, err := decodeResponse(id); err != nil {
			return err
		} else {
			glog("Players:", resp.(*api.DebugCommandResponse).Error)
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
