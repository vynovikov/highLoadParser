package dataHandler

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type dataHandlerSuite struct {
	suite.Suite
}

func TestDataHandlerSuite(t *testing.T) {

	suite.Run(t, new(dataHandlerSuite))
}

func (s *dataHandlerSuite) TestCreate() {
	tt := []struct {
		name              string
		initDataHandler   DataHandler
		dto               DataHandlerDTO
		wantedDataHandler DataHandler
		wantedError       error
	}{

		{
			name: "1. Empty Map, dto.B() == False, !isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("azazaza"), isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {false: {}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "2. Empty Map, dto.B() == False, isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("azazaza"), isSub: true, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 0}: {true: {}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "3. Map has other TS, dto.B() == False, !isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {false: {}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {false: {}}},
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {false: {}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "4. Map has other TS, dto.B() == False, isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {false: {}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), isSub: true, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {false: {}}},
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {true: {}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "5. Map has same key and different part, dto.B() == False",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 3}: {false: {}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {false: {}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "6. Map has same key and same part, dto.B() == False",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {false: {}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {false: {}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},
	}

	for _, v := range tt {

		s.Run(v.name, func() {

			err := v.initDataHandler.Create(v.dto)

			if v.wantedError != nil {

				s.Equal(v.wantedError, err)
			}

			s.Equal(v.wantedDataHandler, v.initDataHandler)
		})
	}
}
