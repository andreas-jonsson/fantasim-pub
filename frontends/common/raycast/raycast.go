/*
Copyright (C) 2017 Andreas T Jonsson

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

package raycast

import (
	"image"
	"image/color"
	"math"

	"github.com/ungerik/go3d/float64/vec2"
)

type Raycaster struct {
	pos, dir, plane vec2.T
	renderTarget    RenderTarget
	world           World
}

type RenderTarget interface {
	Bounds() image.Rectangle
	Set(x, y int, c color.Color)
	SetZ(x int, z float64)
	GetZ(x int) float64
}

type Texture interface {
	Bounds() image.Rectangle
	At(x, y int) color.Color
}

type World interface {
	GetTexture(index, shade int) Texture
	GetTile(x, y int) int
}

func NewRaycaster(rt RenderTarget, w World) *Raycaster {
	return &Raycaster{
		dir:          vec2.T{-1.0, 0.0},
		plane:        vec2.T{0.0, 0.66},
		renderTarget: rt,
		world:        w,
	}
}

func (rc *Raycaster) Rotate(r float64) {
	tmp := rc.dir[0]
	rc.dir[0] = rc.dir[0]*math.Cos(r) - rc.dir[1]*math.Sin(r)
	rc.dir[1] = tmp*math.Sin(r) + rc.dir[1]*math.Cos(r)

	tmp = rc.plane[0]
	rc.plane[0] = rc.plane[0]*math.Cos(r) - rc.plane[1]*math.Sin(r)
	rc.plane[1] = tmp*math.Sin(r) + rc.plane[1]*math.Cos(r)
}

func (rc *Raycaster) Move(v vec2.T) {
	rc.pos.Add(&v)
}

func (rc *Raycaster) Dir() vec2.T {
	return rc.dir
}

func (rc *Raycaster) Pos() vec2.T {
	return rc.pos
}

func (rc *Raycaster) Render() {
	// Reference: http://lodev.org/cgtutor/raycasting.html

	rtSize := rc.renderTarget.Bounds().Size()
	for x := 0; x < rtSize.X; x++ {
		// Calculate ray position and direction.
		cameraX := 2*float64(x)/float64(rtSize.X) - 1 // X coordinate in camera space.
		rayPos := rc.pos

		v := vec2.T{rc.plane[0] * cameraX, rc.plane[1] * cameraX}
		rayDir := vec2.Add(&rc.dir, &v)

		// Which box of the map we're in.
		mapIndex := [2]int{int(rayPos[0]), int(rayPos[1])}

		// Length of ray from one X or Y-side to next X or Y-side.
		raySqr := vec2.Mul(&rayDir, &rayDir)
		deltaDist := vec2.T{
			math.Sqrt(1 + raySqr[1]/raySqr[0]),
			math.Sqrt(1 + raySqr[0]/raySqr[1]),
		}

		var (
			sideDist vec2.T // Length of ray from current position to next X or Y-side.
			step     [2]int // What direction to step in X or Y-direction. (Either +1 or -1)
			side     int    // Was a NS or a EW wall hit?
		)

		// Calculate step and initial sideDist.
		if rayDir[0] < 0 {
			step[0] = -1
			sideDist[0] = (rayPos[0] - float64(mapIndex[0])) * deltaDist[0]
		} else {
			step[0] = 1
			sideDist[0] = (float64(mapIndex[0]) + 1 - rayPos[0]) * deltaDist[0]
		}

		if rayDir[1] < 0 {
			step[1] = -1
			sideDist[1] = (rayPos[1] - float64(mapIndex[1])) * deltaDist[1]
		} else {
			step[1] = 1
			sideDist[1] = (float64(mapIndex[1]) + 1 - rayPos[1]) * deltaDist[1]
		}

		var tileIndex int

		// DDA loop.
		for {
			// Jump to next map square, OR in X-direction, OR in Y-direction.
			if sideDist[0] < sideDist[1] {
				sideDist[0] += deltaDist[0]
				mapIndex[0] += step[0]
				side = 0
			} else {
				sideDist[1] += deltaDist[1]
				mapIndex[1] += step[1]
				side = 1
			}

			// Check if ray has hit a wall.
			tileIndex = rc.world.GetTile(mapIndex[0], mapIndex[1])
			if tileIndex > 0 {
				tileIndex--
				break
			}
		}

		// Calculate distance of perpendicular ray. (Oblique distance will give fisheye effect!)
		var perpWallDist float64
		if side == 0 {
			perpWallDist = (float64(mapIndex[0]) - rayPos[0] + float64(1-step[0])/2) / rayDir[0]
		} else {
			perpWallDist = (float64(mapIndex[1]) - rayPos[1] + float64(1-step[1])/2) / rayDir[1]
		}

		// Calculate height of line to draw on screen.
		lineHeight := int(float64(rtSize.Y) / perpWallDist)

		// Calculate lowest and highest pixel to fill in current stripe.
		drawStart := -lineHeight/2 + rtSize.Y/2
		if drawStart < 0 {
			drawStart = 0
		}

		drawEnd := lineHeight/2 + rtSize.Y/2
		if drawEnd >= rtSize.Y {
			drawEnd = rtSize.Y - 1
		}

		// Calculate value of wallX.
		var wallX float64 // Where exactly the wall was hit.
		if side == 0 {
			wallX = rayPos[1] + perpWallDist*rayDir[1]
		} else {
			wallX = rayPos[0] + perpWallDist*rayDir[0]
		}
		wallX -= math.Floor(wallX)

		texture := rc.world.GetTexture(tileIndex, side)
		texSize := texture.Bounds().Size()

		// X coordinate on the texture.
		texX := int(wallX * float64(texSize.X))
		if side == 0 && rayDir[0] > 0 {
			texX = texSize.X - texX - 1
		}

		if side == 1 && rayDir[1] < 0 {
			texX = texSize.X - texX - 1
		}

		rc.renderTarget.SetZ(x, perpWallDist)

		for y := drawStart; y < drawEnd; y++ {
			d := y - rtSize.Y/2 + lineHeight/2
			texY := int(float64(d*texSize.Y) / float64(lineHeight))

			// Using the color interface is slow... but convenient. :)
			c := texture.At(texX, texY)
			//if side > 0 {
			//	r, g, b, _ := c.RGBA()
			//	rc.renderTarget.Set(x, y, color.RGBA{uint8(r / 2), uint8(g / 2), uint8(b / 2), 255})
			//} else {
			rc.renderTarget.Set(x, y, c)
			//}
		}
	}
}
