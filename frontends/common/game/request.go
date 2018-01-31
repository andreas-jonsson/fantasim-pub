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
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/andreas-jonsson/fantasim-pub/api"
)

const (
	invalidRequestID = -1
	decodeTimeout    = time.Second * 10
)

var (
	requestID    int
	responseChan chan asyncResponse
	infoChan     chan string
)

type asyncResponse struct {
	obj interface{}
	id  int
}

func encodeRequest(enc *json.Encoder, obj interface{}) (int, error) {
	if requestID++; requestID == invalidRequestID {
		requestID++
	}
	if err := api.EncodeRequest(enc, obj, requestID); err != nil {
		return requestID, err
	}
	return requestID, nil
}

func decodeResponseTimeout(id int, timeout time.Duration) (interface{}, error) {
	t := time.NewTimer(timeout)
	defer t.Stop()
	for {
		select {
		case resp, ok := <-responseChan:
			if !ok {
				return nil, errors.New("response channel was closed")
			}
			if resp.id == id {
				return resp.obj, nil
			}
			responseChan <- resp
		case <-t.C:
			return nil, errors.New("decode timeout")
		}
	}
}

func decodeResponse(id int) (interface{}, error) {
	return decodeResponseTimeout(id, decodeTimeout)
}

func discardResponse(id int) {
	go func() {
		decodeResponseTimeout(id, time.Minute)
	}()
}

func startAsyncDecoder(dec *json.Decoder) {
	responseChan = make(chan asyncResponse, 1024)
	go func() {
		for {
			resp, id, err := api.DecodeResponse(dec)
			if err != nil {
				log.Println(err)
				close(responseChan)
				return
			}
			responseChan <- asyncResponse{resp, id}
		}
	}()
}

func startInfoDecode(dec *json.Decoder) {
	infoChan = make(chan string, 1024)
	go func() {
		for {
			var str string
			if err := dec.Decode(&str); err != nil {
				log.Println(err)
				close(infoChan)
				return
			}
			select {
			case infoChan <- str:
			default:
				log.Println("info channel is full")
			}
		}
	}()
}
