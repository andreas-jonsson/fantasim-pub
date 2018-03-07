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
	"time"

	"github.com/andreas-jonsson/fantasim-pub/frontends/common/game"
	"github.com/andreas-jonsson/fantasim-pub/frontends/sdl/system"
	"golang.org/x/net/websocket"
)

const lobbyURL = "http://lobby.fantasim.net"

var playerKey, playerName string

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
	flag.StringVar(&playerKey, "key", os.Getenv("fantasim_key"), "The player key assigned by the server.")
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
		serverAddress = u.Host
	}

	origin := fmt.Sprintf("http://%s/", serverAddress)
	apiWs, err := websocket.Dial(fmt.Sprintf("ws://%s/api", serverAddress), "", origin)
	if err != nil {
		open(lobbyURL)
		log.Fatalln(err)
	}
	defer apiWs.Close()
	apiSocket := newTimeoutReadWriter(apiWs, time.Second*3)

	infoWs, err := websocket.Dial(fmt.Sprintf("ws://%s/info", serverAddress), "", origin)
	if err != nil {
		open(lobbyURL)
		log.Fatalln(err)
	}
	defer infoWs.Close()

	if err := json.NewEncoder(io.MultiWriter(apiSocket, infoWs)).Encode(&playerKey); err != nil {
		log.Fatalln(err)
	}

	if err := json.NewEncoder(apiSocket).Encode("gob"); err != nil {
		log.Fatalln(err)
	}

	enc := gob.NewEncoder(apiSocket)
	dec := gob.NewDecoder(apiSocket)
	decInfo := json.NewDecoder(infoWs)

	sz := image.Pt(1280, 720)
	sys, err := system.InitSDL(sz, sz, false)
	if err != nil {
		log.Fatalln(err)
	}
	defer sys.Quit()

	if err := game.Initialize(sys, sys); err != nil {
		log.Fatalln(err)
	}

	if err := game.Start(enc, dec, decInfo); err != nil {
		log.Fatalln(err)
	}
}

type TimeoutReadWriter struct {
	conn *websocket.Conn
	t    time.Duration
}

func newTimeoutReadWriter(conn *websocket.Conn, timeout time.Duration) io.ReadWriter {
	return &TimeoutReadWriter{conn, timeout}
}

func (trw *TimeoutReadWriter) Read(p []byte) (int, error) {
	conn := trw.conn
	conn.SetReadDeadline(time.Now().Add(trw.t))
	return conn.Read(p)
}

func (trw *TimeoutReadWriter) Write(b []byte) (int, error) {
	conn := trw.conn
	conn.SetWriteDeadline(time.Now().Add(trw.t))
	return conn.Write(b)
}

func open(url string) error {
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

	return exec.Command(cmd, append(args, url)...).Start()
}
