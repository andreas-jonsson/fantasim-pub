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
	"fmt"

	"github.com/andreas-jonsson/fantasim-pub/api"
	"github.com/andreas-jonsson/vsdl-go"
)

type menuOption struct {
	text    string
	key     vsdl.Keycode
	cb      func(api.Encoder) error
	subMenu *menuPage
}

type menuPage struct {
	title   string
	options []*menuOption
}

var rootMenu = &menuPage{
	title: "Menu",
	options: []*menuOption{
		{
			key:  'b',
			text: "Build",
			subMenu: &menuPage{
				title: "Build",
				options: []*menuOption{
					{
						key:  's',
						text: "Stockpile",
						cb:   buildStockpile,
					},
				},
			},
		},
		{
			key:  'd',
			text: "Designate",
			subMenu: &menuPage{
				title: "Designate",
				options: []*menuOption{
					{
						key:  'e',
						text: "Explore",
						cb:   exploreLocation,
					},
					{
						key:  'c',
						text: "Cut trees",
						cb:   designateTreeCutting,
					},
				},
			},
		},
		{
			key:  'h',
			text: "Jump to home location",
			cb:   cameraToHomeLocation,
		},
		{
			key:  'j',
			text: "Print job queue to log",
			cb:   printJobQueue,
		},
	},
}

var menuStack = []*menuPage{rootMenu}

func resetMenuWindow() {
	menuStack = menuStack[:1]
}

func updateCtrlWindow(keysym vsdl.Keysym) func(api.Encoder) error {
	currentMenu := menuStack[len(menuStack)-1]

	ln := len(menuStack)
	if ln > 1 && keysym.Sym == vsdl.BackSpaceKey {
		menuStack = menuStack[:ln-1]
		return nil
	}

	for _, o := range currentMenu.options {
		if o.key == keysym.Sym {
			if o.subMenu != nil {
				menuStack = append(menuStack, o.subMenu)
			} else {
				return o.cb
			}
		}
	}
	return nil
}

func updateCtrlWindowText(text []string) ([]string, string) {
	text = text[:0]
	currentMenu := menuStack[len(menuStack)-1]

	for _, o := range currentMenu.options {
		text = append(text, fmt.Sprintf(" %s: %s", string(o.key), o.text))
		text = append(text, "")
	}

	return text, currentMenu.title
}
