package controllers

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vynovikov/highLoadParser/internal/service"
)

type controllersSuite struct {
	suite.Suite
}

func TestControllersSuite(t *testing.T) {
	suite.Run(t, new(controllersSuite))
}

func (s *controllersSuite) TestEvolve() {

	tt := []struct {
		name      string
		initDTO   *parserServiceInitDTO
		wantedDTO *parserServiceInitDTO
	}{
		{
			name:      "0. Nil",
			initDTO:   nil,
			wantedDTO: nil,
		},

		{
			name: "1. CR in the end",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + "\r"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.Probably,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
				pssu: &service.ParserServiceSub{
					PSSH: service.ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: service.ParserServiceSubBody{
						B: []byte("\r"),
					},
				},
			},
		},

		{
			name: "2. CRLF in the end",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + "\r\n"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.Probably,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
				pssu: &service.ParserServiceSub{
					PSSH: service.ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: service.ParserServiceSubBody{
						B: []byte("\r\n"),
					},
				},
			},
		},

		{
			name: "3. No full boundary and no partial boundary",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + sep + "bPrefixbRo"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.Probably,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
				pssu: &service.ParserServiceSub{
					PSSH: service.ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: service.ParserServiceSubBody{
						B: []byte("\r\nbPrefixbRo"),
					},
				},
			},
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {

			v.initDTO.Evolve(0)

			s.Equal(v.wantedDTO, v.initDTO)
		})
	}
}
