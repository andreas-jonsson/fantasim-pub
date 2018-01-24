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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/andreas-jonsson/fantasim-pub/frontends/common/game"
	"golang.org/x/net/websocket"
)

var playerKey, playerName string

func main() {
	var (
		fantasimUrl,
		serverAddress string
	)

	flag.StringVar(&fantasimUrl, "url", "", "fantasim URL startup.")
	flag.StringVar(&serverAddress, "host", "localhost", "Server address.")
	flag.StringVar(&playerKey, "key", os.Getenv("fantasim-key"), "The player key assigned by the server.")
	flag.StringVar(&playerName, "name", "Unknown", "The name of the player.")
	flag.Parse()

	if fantasimUrl != "" {
		u, err := url.Parse(fantasimUrl)
		if err != nil {
			log.Fatalln(err)
		}

		q := u.Query()
		playerName = q.Get("name")
		playerKey = q.Get("key")
		serverAddress = u.Hostname()
	}

	if err := game.Initialize(); err != nil {
		log.Fatalln(err)
	}

	origin := fmt.Sprintf("http://%s/", serverAddress)
	apiWs, err := websocket.Dial(fmt.Sprintf("ws://%s/api", serverAddress), "", origin)
	if err != nil {
		log.Fatalln(err)
	}
	defer apiWs.Close()

	infoWs, err := websocket.Dial(fmt.Sprintf("ws://%s/info", serverAddress), "", origin)
	if err != nil {
		log.Fatalln(err)
	}
	defer infoWs.Close()

	enc := json.NewEncoder(io.MultiWriter(apiWs, infoWs))
	if err := enc.Encode(&playerKey); err != nil {
		log.Fatalln(err)
	}

	if err := game.Start(apiWs, infoWs); err != nil {
		log.Fatalln(err)
	}
}
