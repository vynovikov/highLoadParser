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
			name: "",
			initDataHandler: &memoryDataHandlerStruct{
				Map:    map[key]map[bool]value{},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("azazaza"), isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[key]map[bool]value{
					{
						TS:   "qqq",
						Part: 1,
					}: {
						false: {},
					},
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
