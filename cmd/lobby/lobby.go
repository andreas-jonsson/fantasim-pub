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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/andreas-jonsson/fantasim-pub/lobby"
	"golang.org/x/net/websocket"
)

func noCacheFunc(handle func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
		header.Set("Cache-Control", "no-cache")
		header.Set("Cache-Control", "no-store")
		header.Set("Cache-Control", "must-revalidate")
		handle(w, r)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	var buf bytes.Buffer
	fmt.Fprintln(&buf, `<!DOCTYPE html><html><head><title>Fantasim Servers</title></head><body>`)
	fmt.Fprintln(&buf, `<h1>Fantasim Servers</h1>`)

	gameServers.Range(func(k, v interface{}) bool {
		key := k.(string)
		msg := v.(lobby.Message)

		if msg.Data != "" {
			fmt.Fprintf(&buf, `<img width="128" height="128" src="data:image/png;base64,%s"/><br>`, msg.Data)
		}

		fmt.Fprintf(&buf, `<a href="http://%s">%s (%s)</a><br>`, key, msg.Name, key)
		return true
	})

	fmt.Fprintln(&buf, `</body></html>`)
	ln := buf.Len()

	w.Header().Set("Content-Length", fmt.Sprint(ln))
	buf.WriteTo(w)
}

func wsHandler(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()

	var msg lobby.Message
	msg.SetTimestamp(time.Now())

	dec := json.NewDecoder(ws)
	if err := dec.Decode(&msg); err != nil {
		return
	}

	switch msg.Type {
	case "ping":
		if v, ok := gameServers.Load(msg.Host); ok {
			data := v.(lobby.Message).Data
			if data != "" && msg.Data == "" {
				msg.Data = data
			}
		}
		gameServers.Store(msg.Host, msg)
	case "close":
		fmt.Println("close message was sent from game server")
		gameServers.Delete(msg.Host)
	default:
		fmt.Printf("invalid message type: %s\n", msg.Type)
		gameServers.Delete(msg.Host)
	}
}

var gameServers sync.Map

func main() {
	go func() {
		for range time.Tick(time.Second) {
			gameServers.Range(func(k, v interface{}) bool {
				msg := v.(lobby.Message)
				if time.Since(msg.Timestamp()) > 30*time.Second {
					gameServers.Delete(k.(string))
				}
				return true
			})
		}
	}()

	http.HandleFunc("/", noCacheFunc(rootHandler))
	http.Handle("/ws", websocket.Handler(wsHandler))

	log.Println("Webserver running...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalln(err)
	}
}
