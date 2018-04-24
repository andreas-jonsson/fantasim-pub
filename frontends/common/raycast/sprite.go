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
	"math"
	"sort"

	"github.com/ungerik/go3d/float64/vec2"
)

type SpriteInstance struct {
	Pos vec2.T
	Tex Texture
	ln  float64
}

type SpriteInstances []SpriteInstance

func (s SpriteInstances) Len() int           { return len(s) }
func (s SpriteInstances) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SpriteInstances) Less(i, j int) bool { return s[i].ln < s[j].ln }

type Spritecaster struct {
	sprites SpriteInstances
}

func NewSpritecaster(sprites SpriteInstances) *Spritecaster {
	return &Spritecaster{sprites}
}

func (sc *Spritecaster) Render(rc *Raycaster) {
	// Reference: http://lodev.org/cgtutor/raycasting3.html

	pos := rc.pos
	dir := rc.dir
	plane := rc.plane

	// Calculate the distance to all sprites.
	for _, si := range sc.sprites {
		si.ln = si.Pos.Sub(&pos).Length()
	}

	// Sort sprites. (Back to front.)
	sort.Sort(sc.sprites)

	rtSize := rc.renderTarget.Bounds().Size()
	for _, s := range sc.sprites {
		// Translate sprite position to relative to camera.
		spritePos := s.Pos.Sub(&pos)

		// Transform sprite with the inverse camera matrix.
		invDet := 1.0 / (plane[0]*dir[1] - dir[0]*plane[1])

		transformX := invDet * (dir[1]*spritePos[0] - dir[0]*spritePos[1])
		transformY := invDet * (-plane[1]*spritePos[0] + plane[0]*spritePos[1])

		spriteScreenX := int(float64(rtSize.X/2) * (1 + transformX/transformY))

		// Calculate height of the sprite on screen.
		spriteHeight := int(math.Abs(float64(rtSize.Y) / transformY)) // Using transformY instead of the real distance prevents fisheye.

		// Calculate lowest and highest pixel to fill in current stripe.
		drawStartY := -spriteHeight/2 + rtSize.Y/2
		if drawStartY < 0 {
			drawStartY = 0
		}

		drawEndY := spriteHeight/2 + rtSize.Y/2
		if drawEndY >= rtSize.Y {
			drawEndY = rtSize.Y - 1
		}

		// Calculate width of the sprite.
		spriteWidth := int(math.Abs(float64(rtSize.Y) / (transformY)))
		drawStartX := -spriteWidth/2 + spriteScreenX
		if drawStartX < 0 {
			drawStartX = 0
		}

		drawEndX := spriteWidth/2 + spriteScreenX
		if drawEndX >= rtSize.X {
			drawEndX = rtSize.X - 1
		}

		texSize := s.Tex.Bounds().Size()

		// Loop through every vertical stripe of the sprite on screen.
		for x := drawStartX; x < drawEndX; x++ {
			texX := (x - (-spriteWidth/2 + spriteScreenX)) * texSize.X / spriteWidth

			if transformY > 0 && x > 0 && x < rtSize.X && transformY < rc.renderTarget.GetZ(x) {
				for y := drawStartY; y < drawEndY; y++ {
					d := y - rtSize.Y/2 + spriteHeight/2
					texY := int(float64(d*texSize.Y) / float64(spriteHeight))

					c := s.Tex.At(texX, texY)
					r, g, b, a := c.RGBA()

					if a < 255 || r+g+b > 0 {
						rc.renderTarget.Set(x, y, c)
					}
				}
			}
		}
	}
}
