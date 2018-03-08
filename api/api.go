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

	// Wildlife
	Dear
	Boar
	Wolf
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
	case Dear:
		return "Dear"
	case Boar:
		return "Boar"
	case Wolf:
		return "Wolf"
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
	NoItem ItemClass = iota
	PartialItem
	LogItem
	FirewoodItem
	PlankItem
	StoneItem
	MeatItem
	BonesItem
	SeedsItem
	CropItem

	// Corpses
	HumanCorpseItem
	DwarfCorpseItem
	GoblinCorpseItem
	OrcCorpseItem
	TrollCorpseItem
	ElvenCorpseItem
	DeamonCorpseItem

	DearCorpseItem
	BoarCorpseItem
	WolfCorpseItem
)

func (c ItemClass) String() string {
	switch c {
	case PartialItem:
		return "Partial"
	case LogItem:
		return "Log"
	case FirewoodItem:
		return "Firewood"
	case PlankItem:
		return "Plank"
	case StoneItem:
		return "Stone"
	case MeatItem:
		return "Meat"
	case BonesItem:
		return "Bones"
	case SeedsItem:
		return "Seeds"
	case CropItem:
		return "Crop"
	case HumanCorpseItem, DwarfCorpseItem, GoblinCorpseItem, OrcCorpseItem, TrollCorpseItem, ElvenCorpseItem, DeamonCorpseItem:
		return "Corpse"
	case DearCorpseItem, BoarCorpseItem, WolfCorpseItem:
		return "Animal Corpse"
	default:
		panic("invalid item type")
	}
}

type BuildingType uint8

const (
	NoBuilding BuildingType = iota
	StockpileBuilding
	SawmillBuilding
	ButcherShoppBuilding
	FarmBuilding
)

func (b BuildingType) String() string {
	switch b {
	case StockpileBuilding:
		return "Stockpile"
	case SawmillBuilding:
		return "Sawmill"
	case ButcherShoppBuilding:
		return "Butcher Shopp"
	case FarmBuilding:
		return "Farm"
	default:
		panic("invalid building type")
	}
}

type StructureType uint8

const (
	NoStructure StructureType = iota
	WallStructure
)

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
	Flags             TileFlag       `json:"flags"`
	Height            uint8          `json:"height"`
	BuildingType      BuildingType   `json:"building_type"`
	StructureType     StructureType  `json:"structure_type"`
	StructureMaterial ItemClass      `json:"structure_material"`
	Building          uint64         `json:"building"`
	Units             []UnitViewData `json:"units"`
	Items             []ItemViewData `json:"items"`
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

type MineLocationRequest struct {
	Location Point `json:"location"`
}

type MineLocationResponse struct {
	Error string `json:"error"`
}

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
	Hunger float32  `json:"hunger"`
	Debug  []string `json:"debug"`
}

type CutTreesRequest struct {
	Trees []Point `json:"trees"`
}

type CutTreesResponse Empty

type AttackUnitsRequest struct {
	Units []uint64 `json:"units"`
}

type AttackUnitsResponse Empty

type CollectItemsRequest struct {
	Items []Point `json:"items"`
}

type CollectItemsResponse Empty

type GatherSeedsRequest struct {
	Seeds []Point `json:"seeds"`
}

type GatherSeedsResponse Empty

type SeedFarmRequest struct {
	BuildingID uint64 `json:"building_id"`
}

type SeedFarmResponse Empty

type DesignateRequest struct {
	Building BuildingType `json:"building"`
	Location Rect         `json:"location"`
}

type DesignateResponse struct {
	Error string `json:"error"`
}

type BuildRequest struct {
	Structure StructureType `json:"structure"`
	Material  ItemClass     `json:"material"`
	Location  Rect          `json:"location"`
}

type BuildResponse Empty

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

func register(req, resp interface{}) {
	registerType(requestTypeRegistry, req)
	registerType(responseTypeRegistry, resp)
}

func init() {
	registerType(requestTypeRegistry, CloseRequest{})
	registerType(requestTypeRegistry, DebugCommandRequest{})
	registerType(requestTypeRegistry, CreateViewRequest{})
	registerType(requestTypeRegistry, DestroyViewRequest{})
	registerType(requestTypeRegistry, UpdateViewRequest{})
	registerType(requestTypeRegistry, ReadViewRequest{})
	registerType(requestTypeRegistry, ExploreLocationRequest{})
	registerType(requestTypeRegistry, MineLocationRequest{})
	registerType(requestTypeRegistry, JobQueueRequest{})
	registerType(requestTypeRegistry, ViewHomeRequest{})
	registerType(requestTypeRegistry, UnitStatsRequest{})
	registerType(requestTypeRegistry, CutTreesRequest{})
	registerType(requestTypeRegistry, CollectItemsRequest{})
	registerType(requestTypeRegistry, BuildRequest{})
	registerType(requestTypeRegistry, AttackUnitsRequest{})
	registerType(requestTypeRegistry, GatherSeedsRequest{})

	registerType(responseTypeRegistry, Empty{})
	registerType(responseTypeRegistry, DebugCommandResponse{})
	registerType(responseTypeRegistry, CreateViewResponse{})
	registerType(responseTypeRegistry, DestroyViewResponse{})
	registerType(responseTypeRegistry, UpdateViewResponse{})
	registerType(responseTypeRegistry, ReadViewResponse{})
	registerType(responseTypeRegistry, ExploreLocationResponse{})
	registerType(responseTypeRegistry, MineLocationResponse{})
	registerType(responseTypeRegistry, JobQueueResponse{})
	registerType(responseTypeRegistry, ViewHomeResponse{})
	registerType(responseTypeRegistry, UnitStatsResponse{})
	registerType(responseTypeRegistry, CutTreesResponse{})
	registerType(responseTypeRegistry, CollectItemsResponse{})
	registerType(responseTypeRegistry, BuildResponse{})
	registerType(responseTypeRegistry, AttackUnitsResponse{})
	registerType(responseTypeRegistry, GatherSeedsResponse{})

	register(SeedFarmRequest{}, SeedFarmResponse{})
	register(DesignateRequest{}, DesignateResponse{})
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
