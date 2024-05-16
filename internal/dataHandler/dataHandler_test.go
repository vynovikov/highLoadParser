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
			name: "1. Empty Map, !isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("azazaza"), b: False, e: True, isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {false: {e: True}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "2. Empty Map, isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("azazaza"), b: False, e: Probably, isSub: true, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 0}: {true: {e: Probably}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "3. Map has other TS, !isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {false: {e: True}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), b: False, e: True, isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {false: {e: True}}},
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {false: {e: True}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "4. Map has other TS, isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {false: {e: Probably}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), b: False, e: Probably, isSub: true, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {false: {e: Probably}}},
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {true: {e: Probably}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "5. Map has same key and different part, !isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 3}: {false: {e: True}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), b: False, e: True, isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {false: {e: True}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "6. Map has same key and different part, isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 3}: {false: {e: True}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), b: False, e: Probably, isSub: true, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {true: {e: Probably}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "7. Map has same key and same part, !isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {false: {e: True}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), b: False, e: True, isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {false: {e: True}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "8. Map has same key and same part, isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {false: {e: True}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), b: False, e: Probably, isSub: true, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {false: {e: Probably}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "9. Map has same key and same part, value.e = True, !d.IsSub, d.E() == Probably",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {false: {e: True}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), b: False, e: Probably, isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {false: {e: Probably}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "10. Map has same key and same part, value.e = Probably, !d.IsSub, d.E() == Probably",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {true: {e: Probably}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), b: False, e: Probably, isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {true: {e: Probably}, false: {e: Probably}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "11. Map has same key and same part, value.e == Probably, d.IsSub, d.E() == Probably",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {false: {e: Probably}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("azazaza"), b: False, e: Probably, isSub: true, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {false: {e: Probably}, true: {e: Probably}}},
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
