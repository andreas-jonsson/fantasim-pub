/*
Copyright (c) 2017-2018 Andreas T Jonsson

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

package api

import (
	"fmt"
	"reflect"
)

const InvalidID uint64 = 0

type TileFlag uint8

const (
	Water TileFlag = 1 << iota
	Brook
	Sand
	Snow
	Tree
	Bush
	Plant
	Stone
)

func (flags TileFlag) Is(f TileFlag) bool {
	return flags&f != 0
}

type Allegiance uint8

const (
	Friendly Allegiance = iota
	Neutral
	Hostile
)

type UnitRace uint8

const (
	Human UnitRace = iota
	Dwarf
	Goblin
	Orc
	Troll
	Elven
	Deamon
)

func (r UnitRace) String() string {
	switch r {
	case Human:
		return "Human"
	case Dwarf:
		return "Dwarf"
	case Goblin:
		return "Goblin"
	case Orc:
		return "Orc"
	case Troll:
		return "Troll"
	case Elven:
		return "Elven"
	case Deamon:
		return "Deamon"
	default:
		return "Unknown Race"
	}
}

type UnitClass uint8

const (
	UnknownUnitClass UnitClass = iota
)

type ItemClass uint8

const (
	Log ItemClass = iota
	Firewood
	Plank
)

func (c ItemClass) String() string {
	switch c {
	case Log:
		return "Log"
	case Firewood:
		return "Firewood"
	case Plank:
		return "Plank"
	default:
		panic("invalid building type")
	}
}

type BuildingType uint8

const (
	StockpileBuilding BuildingType = iota
	SawmillBuilding
)

func (b BuildingType) String() string {
	switch b {
	case StockpileBuilding:
		return "Stockpile"
	case SawmillBuilding:
		return "Sawmill"
	default:
		panic("invalid building type")
	}
}

type Header struct {
	Type string `json:"type"`
	Id   int    `json:"id"`
}

type (
	Empty struct{}
	Any   map[interface{}]interface{}
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Rect struct {
	Min Point `json:"min"`
	Max Point `json:"max"`
}

type CloseRequest Empty

type CreateViewRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type CreateViewResponse struct {
	ViewID int `json:"view_id"`
}

type DestroyViewRequest struct {
	ViewID int `json:"view_id"`
}

type DestroyViewResponse Empty

type UpdateViewRequest struct {
	ViewID int `json:"view_id"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type UpdateViewResponse Empty

type ReadViewRequest struct {
	ViewID int `json:"view_id"`
}

type UnitViewData struct {
	ID         uint64     `json:"unit_id"`
	Allegiance Allegiance `json:"allegiance"`
	Race       UnitRace   `json:"race"`
	Class      UnitClass  `json:"class"`
}

type ItemViewData struct {
	ID    uint64    `json:"item_id"`
	Class ItemClass `json:"class"`
}

type ReadViewData struct {
	Flags        TileFlag       `json:"flags"`
	Height       uint8          `json:"height"`
	Building     uint64         `json:"building"`
	BuildingType BuildingType   `json:"building_type"`
	Units        []UnitViewData `json:"units"`
	Items        []ItemViewData `json:"items"`
}

type ReadViewResponse struct {
	Data []ReadViewData `json:"data"`
}

type DebugCommandRequest struct {
	Command string `json:"command"`
}

type DebugCommandResponse struct {
	Error string `json:"error"`
}

type ExploreLocationRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ExploreLocationResponse Empty

type JobQueueRequest Empty

type JobQueueResponse struct {
	Jobs []string `json:"jobs"`
}

type ViewHomeRequest struct {
	ViewID int `json:"view_id"`
}

type ViewHomeResponse struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type UnitStatsRequest struct {
	UnitID int `json:"unit_id"`
}

type UnitStatsResponse struct {
	Name   string   `json:"name"`
	Health float32  `json:"health"`
	Thirst float32  `json:"thirst"`
	Debug  []string `json:"debug"`
}

type CutTreeData struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type CutTreesRequest struct {
	Trees []CutTreeData `json:"trees"`
}

type CutTreesResponse Empty

type BuildRequest struct {
	Building BuildingType `json:"building"`
	Location Rect         `json:"location"`
}

type BuildResponse struct {
	Error string `json:"error"`
}

var (
	requestTypeRegistry  = make(map[string]reflect.Type)
	responseTypeRegistry = make(map[string]reflect.Type)
)

func getTypeName(v interface{}) (string, reflect.Type) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		return t.Name(), t
	}
	return t.Name(), t
}

func registerType(m map[string]reflect.Type, v interface{}) {
	n, t := getTypeName(v)
	m[n] = t
}

func init() {
	registerType(requestTypeRegistry, CloseRequest{})
	registerType(requestTypeRegistry, DebugCommandRequest{})
	registerType(requestTypeRegistry, CreateViewRequest{})
	registerType(requestTypeRegistry, DestroyViewRequest{})
	registerType(requestTypeRegistry, UpdateViewRequest{})
	registerType(requestTypeRegistry, ReadViewRequest{})
	registerType(requestTypeRegistry, ExploreLocationRequest{})
	registerType(requestTypeRegistry, JobQueueRequest{})
	registerType(requestTypeRegistry, ViewHomeRequest{})
	registerType(requestTypeRegistry, UnitStatsRequest{})
	registerType(requestTypeRegistry, CutTreesRequest{})
	registerType(requestTypeRegistry, BuildRequest{})

	registerType(responseTypeRegistry, Empty{})
	registerType(responseTypeRegistry, DebugCommandResponse{})
	registerType(responseTypeRegistry, CreateViewResponse{})
	registerType(responseTypeRegistry, DestroyViewResponse{})
	registerType(responseTypeRegistry, UpdateViewResponse{})
	registerType(responseTypeRegistry, ReadViewResponse{})
	registerType(responseTypeRegistry, ExploreLocationResponse{})
	registerType(responseTypeRegistry, JobQueueResponse{})
	registerType(responseTypeRegistry, ViewHomeResponse{})
	registerType(responseTypeRegistry, UnitStatsResponse{})
	registerType(responseTypeRegistry, CutTreesResponse{})
	registerType(responseTypeRegistry, BuildResponse{})
}

type (
	Encoder interface {
		Encode(interface{}) error
	}

	Decoder interface {
		Decode(interface{}) error
	}
)

func decode(dec Decoder, m map[string]reflect.Type, op string) (interface{}, int, error) {
	var header Header
	if err := dec.Decode(&header); err != nil {
		return nil, 0, err
	}

	t, ok := m[header.Type]
	if !ok {
		err := fmt.Errorf("message is not a %s: %s", op, header.Type)
		obj := &Any{}
		dec.Decode(&obj)
		return obj, header.Id, err
	}

	obj := reflect.New(t).Interface()
	return obj, header.Id, dec.Decode(obj)
}

func DecodeRequest(dec Decoder) (interface{}, int, error) {
	return decode(dec, requestTypeRegistry, "request")
}

func DecodeResponse(dec Decoder) (interface{}, int, error) {
	return decode(dec, responseTypeRegistry, "response")
}

func encode(enc Encoder, m map[string]reflect.Type, op string, obj interface{}, id int) error {
	name, _ := getTypeName(obj)
	if _, ok := m[name]; !ok {
		return fmt.Errorf("object is not a %s: %s", op, name)
	}

	header := Header{Type: name, Id: id}
	if err := enc.Encode(&header); err != nil {
		return err
	}

	if err := enc.Encode(obj); err != nil {
		return err
	}
	return nil
}

func EncodeRequest(enc Encoder, obj interface{}, id int) error {
	return encode(enc, requestTypeRegistry, "request", obj, id)
}

func EncodeResponse(enc Encoder, obj interface{}, id int) error {
	return encode(enc, responseTypeRegistry, "response", obj, id)
}
