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
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type gameServer struct {
	name, host, thumbnail string
	timestamp             time.Time
}

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
		srv := v.(gameServer)

		if srv.thumbnail != "" {
			fmt.Fprintf(&buf, `<img width="128" height="128" src="data:image/png;base64,%s"/><br>`, srv.thumbnail)
		}

		fmt.Fprintf(&buf, `<a href="http://%s">%s (%s)</a><br>`, key, srv.name, key)
		return true
	})

	fmt.Fprintln(&buf, `</body></html>`)
	ln := buf.Len()

	w.Header().Set("Content-Length", fmt.Sprint(ln))
	buf.WriteTo(w)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	srv := gameServer{timestamp: time.Now()}
	values := r.URL.Query()

	srv.host = values.Get("host")
	if srv.host == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	srv.name = values.Get("name")

	if v, ok := gameServers.Load(srv.host); ok {
		s := v.(gameServer)
		srv.thumbnail = s.thumbnail

		if srv.name == "" {
			srv.name = s.name
		}
	}

	gameServers.Store(srv.host, srv)
}

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	srv := gameServer{timestamp: time.Now()}
	values := r.URL.Query()

	srv.host = values.Get("host")
	if srv.host == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	srv.thumbnail = string(data)

	if v, ok := gameServers.Load(srv.host); ok {
		srv.name = v.(gameServer).name
	}

	gameServers.Store(srv.host, srv)
}

func closeHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		gameServers.Delete(host)
	}
}

var gameServers sync.Map

func main() {
	go func() {
		for range time.Tick(time.Second) {
			gameServers.Range(func(k, v interface{}) bool {
				srv := v.(gameServer)
				if time.Since(srv.timestamp) > 30*time.Second {
					gameServers.Delete(k.(string))
				}
				return true
			})
		}
	}()

	http.HandleFunc("/", noCacheFunc(rootHandler))
	http.HandleFunc("/ping", noCacheFunc(pingHandler))
	http.HandleFunc("/thumbnail", noCacheFunc(thumbnailHandler))
	http.HandleFunc("/close", noCacheFunc(closeHandler))

	log.Println("Webserver running...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalln(err)
	}
}
