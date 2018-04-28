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
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const (
	oldMinor = 0
	oldPatch = 1
)

const (
	newMinor = 0
	newPatch = 2
)

func replaceInFile(file, oldStr, newStr string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	if !strings.Contains(string(data), oldStr) {
		log.Fatalln("Could not find \"", oldStr, "\" in", file)
	}

	newData := []byte(strings.Replace(string(data), oldStr, newStr, 1))
	if err := ioutil.WriteFile(file, newData, 0644); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	format := "VersionString = \"0.%d.%d\""
	replaceInFile("api/api.go", fmt.Sprintf(format, oldMinor, oldPatch), fmt.Sprintf(format, newMinor, newPatch))

	format = "version: 0.%d.%d"
	replaceInFile("snapcraft.yaml", fmt.Sprintf(format, oldMinor, oldPatch), fmt.Sprintf(format, newMinor, newPatch))
	replaceInFile("appveyor.yml", fmt.Sprintf(format, oldMinor, oldPatch), fmt.Sprintf(format, newMinor, newPatch))

	format = "export FANTASIM_SDL_SHORT_VERSION=0.%d"
	replaceInFile(".travis.yml", fmt.Sprintf(format, oldMinor), fmt.Sprintf(format, newMinor))

	format = "export FANTASIM_SDL_VERSION=${FANTASIM_SDL_SHORT_VERSION}.%d"
	replaceInFile(".travis.yml", fmt.Sprintf(format, oldPatch), fmt.Sprintf(format, newPatch))

	format = "#define MyAppVersion \"0.%d.%d\""
	replaceInFile("frontends/tools/package/setup.iss", fmt.Sprintf(format, oldMinor, oldPatch), fmt.Sprintf(format, newMinor, newPatch))
}
