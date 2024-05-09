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
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
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
							E:    service.True,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
			},
		},

		{
			name: "4. No full boundary, partial boundary present",
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

		{
			name: "5. No full boundary, last boundary",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + sep + "bPrefixbRootbSuffix" + sep),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				last: true,
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.False,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
			},
		},

		{
			name: "6. One full boundary",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.False,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza"),
						},
					},
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.False,
							E:    service.True,
						},
						PSB: service.ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				},
			},
		},

		{
			name: "7. One full boundary. Partial boundary present",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + sep + "bPrefixbRo"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.False,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza"),
						},
					},
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.False,
							E:    service.Probably,
						},
						PSB: service.ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
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

		{
			name: "8. One full boundary. CR",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz\r"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.False,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza"),
						},
					},
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.False,
							E:    service.Probably,
						},
						PSB: service.ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
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
			name: "9. Full last boundary after begin piece",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + sep + "bPrefix" + "bRoot" + "bSuffix" + sep),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				last: true,
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.False,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb"),
						},
					},
				},
			},
		},

		{
			name: "10. Full last boundary after begin piece",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + sep + "bPrefix" + "bRoot" + "bSuffix" + sep),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				last: true,
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.False,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb"),
						},
					},
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.False,
							E:    service.False,
						},
						PSB: service.ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
			},
		},

		{
			name: "11. Full boundary in the end",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + sep + "bPrefixbRoot"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.False,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb"),
						},
					},
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.False,
							E:    service.Probably,
						},
						PSB: service.ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
				pssu: &service.ParserServiceSub{
					PSSH: service.ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: service.ParserServiceSubBody{
						B: []byte(sep + "bPrefixbRoot"),
					},
				},
			},
		},

		{
			name: "12. Full boundary in the end with CR",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + sep + "bPrefixbRoot" + "\r"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + sep + "bPrefix" + "bRoot" + sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.False,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb"),
						},
					},
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.False,
							E:    service.Probably,
						},
						PSB: service.ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
				pssu: &service.ParserServiceSub{
					PSSH: service.ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: service.ParserServiceSubBody{
						B: []byte(sep + "bPrefixbRoot" + "\r"),
					},
				},
			},
		},

		{
			name: "13. Partial last boundary after begin piece",
			initDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + sep + "bPrefix" + "bRoot" + "bSuf"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
			},
			wantedDTO: &parserServiceInitDTO{
				part: 0,
				ts:   "qqq",
				body: []byte("a1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				bou:  boundary{prefix: []byte("bPrefix"), root: []byte("bRoot")},
				last: true,
				psus: []*service.ParserServiceUnit{
					{
						PSH: service.ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    service.True,
							E:    service.False,
						},
						PSB: service.ParserServiceBody{
							B: []byte("a1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
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
