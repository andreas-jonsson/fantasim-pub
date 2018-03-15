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
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

type ReadViewRequest struct {
	ViewID int  `json:"view_id"`
	RLE    bool `json:"rle"`
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

type ReadViewBase struct {
	Flags             TileFlag      `json:"flags"`
	UserFlags         UserFlag      `json:"usrflags"`
	Height            uint8         `json:"height"`
	BuildingType      BuildingType  `json:"building_type"`
	StructureType     StructureType `json:"structure_type"`
	StructureMaterial ItemClass     `json:"structure_material"`
	Building          uint64        `json:"building"`
}

type ReadViewData struct {
	ReadViewBase
	Units []UnitViewData `json:"units"`
	Items []ItemViewData `json:"items"`
	RLE   uint8          `json:"rle"`
}

type ReadViewResponse struct {
	W       uint16         `json:"w"`
	H       uint16         `json:"h"`
	RLESize uint32         `json:"rle_size"`
	Data    []ReadViewData `json:"data"`
}

func (rvr *ReadViewResponse) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	dec.Token() // {
	dec.Token() // w
	if err := dec.Decode(&rvr.W); err != nil {
		return err
	}

	dec.Token() // h
	if err := dec.Decode(&rvr.H); err != nil {
		return err
	}

	dec.Token() // rle_size
	if err := dec.Decode(&rvr.RLESize); err != nil {
		return err
	}

	dec.Token() // data
	idx := 0

	// Fastpath if there is no compression
	if rvr.RLESize == 0 {
		if err := dec.Decode(&rvr.Data); err != nil {
			return err
		}
		goto end
	}
	rvr.Data = make([]ReadViewData, int(rvr.W)*int(rvr.H))

	dec.Token() // {

	for i := 0; i < int(rvr.RLESize); i++ {
		rvd := &rvr.Data[idx]
		idx++

		if err := dec.Decode(rvd); err != nil {
			return err
		}

		for j := 0; j < int(rvd.RLE); j++ {
			rvr.Data[idx] = *rvd
			idx++
		}
	}

	dec.Token() // }

end: // Verify the last token
	if tok, err := dec.Token(); err != nil {
		return err
	} else {
		if delim, ok := tok.(json.Delim); !ok || delim != '}' {
			return errors.New("invalid delimiter")
		}
		return nil
	}
}

func (rvr *ReadViewResponse) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(rvr.W)
	enc.Encode(rvr.H)
	enc.Encode(rvr.RLESize)

	if rvr.RLESize == 0 {
		enc.Encode(rvr.Data)
	} else {
		for i := 0; i < int(rvr.RLESize); i++ {
			enc.Encode(rvr.Data[i])
		}
	}
	return buf.Bytes(), nil
}

func (rvr *ReadViewResponse) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&rvr.W); err != nil {
		return err
	}
	if err := dec.Decode(&rvr.H); err != nil {
		return err
	}
	if err := dec.Decode(&rvr.RLESize); err != nil {
		return err
	}

	// Fastpath if there is no compression
	if rvr.RLESize == 0 {
		return dec.Decode(&rvr.Data)
	}
	rvr.Data = make([]ReadViewData, int(rvr.W)*int(rvr.H))

	idx := 0
	for i := 0; i < int(rvr.RLESize); i++ {
		rvd := &rvr.Data[idx]
		idx++

		if err := dec.Decode(rvd); err != nil {
			return err
		}

		for j := 0; j < int(rvd.RLE); j++ {
			rvr.Data[idx] = *rvd
			idx++
		}
	}
	return nil
}
