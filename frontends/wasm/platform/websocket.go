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

package platform

import (
	"bytes"
	"io"
	"log"
	"syscall/js"
	"time"
)

type WebSocket struct {
	sock     js.Value
	buf      bytes.Buffer
	recvChan chan js.Value
}

func Dial(addr string) (ws *WebSocket, err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if jsErr, ok := e.(*js.Error); ok && jsErr != nil {
			ws = nil
			err = jsErr
		} else {
			panic(e)
		}
	}()

	sock := js.Global().Get("WebSocket").New(addr)
	sock.Set("binaryType", "arraybuffer")

	connChan := make(chan struct{}, 1)
	sock.Set("onopen", js.NewEventCallback(js.PreventDefault, func(_ js.Value) {
		connChan <- struct{}{}
	}))

	jsUint8Array := js.Global().Get("Uint8Array")
	recvChan := make(chan js.Value, 4096)

	sock.Set("onmessage", js.NewEventCallback(js.PreventDefault, func(e js.Value) {
		byteArray := jsUint8Array.New(e.Get("data"))
		select {
		case recvChan <- byteArray:
		default:
			log.Println("websocket queue is full")
			recvChan <- byteArray
		}
	}))

	<-connChan
	ws = &WebSocket{
		sock:     sock,
		recvChan: recvChan,
	}
	return
}

func (ws *WebSocket) readData() error {
	var data js.Value
	select {
	case data = <-ws.recvChan:
	case <-time.After(time.Millisecond):
	}

	ln := data.Length()
	for i := 0; i < ln; i++ {
		if err := ws.buf.WriteByte(byte(data.Index(i).Int())); err != nil {
			return err
		}
	}
	return nil
}

func (ws *WebSocket) Read(p []byte) (n int, err error) {
	if ws.buf.Len() > 0 {
		n, err = ws.buf.Read(p)
		if err == io.EOF {
			return n, nil
		}
		return
	}

	if err := ws.readData(); err != nil {
		return 0, err
	}

	n, err = ws.buf.Read(p)
	if err == io.EOF {
		return n, nil
	}
	return
}

func (ws *WebSocket) Write(p []byte) (n int, err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if jsErr, ok := e.(*js.Error); ok && jsErr != nil {
			err = jsErr
			n = 0
		} else {
			panic(e)
		}
	}()

	array := js.TypedArrayOf(p)
	ws.sock.Call("send", array)
	array.Release()

	n = len(p)
	return
}

func (ws *WebSocket) Close() (err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if jsErr, ok := e.(*js.Error); ok && jsErr != nil {
			err = jsErr
		} else {
			panic(e)
		}
	}()

	// Use close code closeNormalClosure to indicate that the purpose
	// for which the connection was established has been fulfilled.
	// See https://tools.ietf.org/html/rfc6455#section-7.4.
	ws.sock.Call("close", closeNormalClosure)
	return
}

// Close codes defined in RFC 6455, section 11.7.
const (
	// 1000 indicates a normal closure, meaning that the purpose for
	// which the connection was established has been fulfilled.
	closeNormalClosure = 1000
)
