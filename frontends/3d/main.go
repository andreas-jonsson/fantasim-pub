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

//go:generate go run ../common/data/generate.go
//go:generate go run ../common/data3d/generate.go

package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/andreas-jonsson/fantasim-pub/frontends/3d/game"
	sys "github.com/andreas-jonsson/fantasim-pub/frontends/sdl/platform"
	"golang.org/x/net/websocket"
)

const lobbyURL = "http://lobby.fantasim.net"

func init() {
	runtime.LockOSThread()
}

func main() {
	var (
		fantasimUrl,
		serverAddress string
	)

	flag.StringVar(&fantasimUrl, "url", "", "fantasim URL startup.")
	flag.StringVar(&serverAddress, "host", "localhost", "Server address.")
	playerKey := flag.String("key", os.Getenv("fantasim_key"), "The player key assigned by the server.")
	noTimeout := flag.Bool("notimeout", false, "Disable socket timout")
	flag.Parse()

	if fantasimUrl != "" {
		u, err := url.Parse(fantasimUrl)
		if err != nil {
			log.Fatalln(err)
		}

		q := u.Query()
		*playerKey = q.Get("key")
		serverAddress = u.Host
	}

	origin := fmt.Sprintf("http://%s/", serverAddress)
	apiWs, err := websocket.Dial(fmt.Sprintf("ws://%s/api", serverAddress), "", origin)
	if err != nil {
		open(lobbyURL)
		log.Fatalln(err)
	}
	defer apiWs.Close()

	infoWs, err := websocket.Dial(fmt.Sprintf("ws://%s/info", serverAddress), "", origin)
	if err != nil {
		open(lobbyURL)
		log.Fatalln(err)
	}
	defer infoWs.Close()

	if err := json.NewEncoder(io.MultiWriter(apiWs, infoWs)).Encode(playerKey); err != nil {
		log.Fatalln(err)
	}

	if err := json.NewEncoder(apiWs).Encode("gob"); err != nil {
		log.Fatalln(err)
	}

	enc := gob.NewEncoder(apiWs)
	dec := gob.NewDecoder(apiWs)
	decInfo := gob.NewDecoder(infoWs)

	sz := image.Pt(1280, 720)
	s, err := sys.InitSDL(sz, sz, false)
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Quit()

	if err := game.Initialize(s, s); err != nil {
		log.Fatalln(err)
	}

	if err := game.Start(enc, dec, decInfo); err != nil {
		log.Fatalln(err)
	}
}

func open(url string) {
	var (
		cmd  string
		args []string
	)

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	if err := exec.Command(cmd, append(args, url)...).Start(); err != nil {
		log.Println(err)
	}
}
