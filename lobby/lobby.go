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

package lobby

import "time"

const TcpPort = 49149

const (
	Ping  = "ping"
	Close = "close"
)

type Message struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Host string `json:"host"`
	Data string `json:"data"`

	timestamp time.Time
}

func (msg *Message) Timestamp() time.Time {
	return msg.timestamp
}

func (msg *Message) SetTimestamp(t time.Time) {
	msg.timestamp = t
}
