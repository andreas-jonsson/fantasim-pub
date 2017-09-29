/*
Copyright (C) 2016 Andreas T Jonsson

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
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/andreas-jonsson/fantasim-pub/lobby"
)

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

var gameServers sync.Map

func main() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", lobby.TcpPort))
	if err != nil {
		log.Fatalln(err)
	}
	defer ln.Close()

	go func() {
		for range time.Tick(time.Second) {
			gameServers.Range(func(k, v interface{}) bool {
				msg := v.(lobby.Message)
				if time.Since(msg.Timestamp()) > time.Minute {
					gameServers.Delete(k.(string))
				}
				return true
			})
		}
	}()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			dec := json.NewDecoder(conn)

			go func() {
				var (
					err  error
					host string
				)

				defer func() {
					if err != nil {
						log.Println(err)
					}
					conn.Close()
				}()

				now := time.Now()
				if err = conn.SetDeadline(now.Add(10 * time.Second)); err != nil {
					gameServers.Delete(host)
					return
				}

				var msg lobby.Message
				msg.SetTimestamp(now)

				if err := dec.Decode(&msg); err != nil {
					gameServers.Delete(host)
					return
				}

				switch msg.Type {
				case "ping":
					host = msg.Host
					if v, ok := gameServers.Load(host); ok {
						data := v.(lobby.Message).Data
						if data != "" && msg.Data == "" {
							msg.Data = data
						}
					}
					gameServers.Store(host, msg)
				case "close":
					err = errors.New("close message was sent from game server")
					gameServers.Delete(host)
					return
				default:
					err = fmt.Errorf("invalid message type: %s", msg.Type)
					gameServers.Delete(host)
					return
				}
			}()
		}
	}()

	http.HandleFunc("/", rootHandler)
	log.Println("Webserver running...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalln(err)
	}
}
