package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vynovikov/highLoadParser/internal/entities"
)

type serviceSuite struct {
	suite.Suite
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(serviceSuite))
}

func (s *serviceSuite) TestEvolve() {

	tt := []struct {
		name      string
		initDTO   *ParserServiceDTO
		wantedDTO *ParserServiceDTO
	}{
		{
			name:      "0. Nil",
			initDTO:   nil,
			wantedDTO: nil,
		},

		{
			name: "1. CR in the end",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + "\r"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    Probably,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
				pssu: &ParserServiceSub{
					PSSH: ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: ParserServiceSubBody{
						B: []byte("\r"),
					},
				},
			},
		},

		{
			name: "2. CRLF in the end",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + "\r\n"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    Probably,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
				pssu: &ParserServiceSub{
					PSSH: ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: ParserServiceSubBody{
						B: []byte("\r\n"),
					},
				},
			},
		},

		{
			name: "3. No full boundary and no partial boundary",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    True,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
			},
		},

		{
			name: "4. No full boundary, partial boundary present",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + entities.Sep + "bPrefixbRo"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    Probably,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
				pssu: &ParserServiceSub{
					PSSH: ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: ParserServiceSubBody{
						B: []byte("\r\nbPrefixbRo"),
					},
				},
			},
		},

		{
			name: "5. No full boundary, last boundary",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + entities.Sep + "bPrefixbRootbSuffix" + entities.Sep),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				last: true,
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
			},
		},

		{
			name: "6. One full boundary",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza"),
						},
					},
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    False,
							E:    True,
						},
						PSB: ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				},
			},
		},

		{
			name: "7. One full boundary. Partial boundary present",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + entities.Sep + "bPrefixbRo"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza"),
						},
					},
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    False,
							E:    Probably,
						},
						PSB: ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				},
				pssu: &ParserServiceSub{
					PSSH: ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: ParserServiceSubBody{
						B: []byte("\r\nbPrefixbRo"),
					},
				},
			},
		},

		{
			name: "8. One full boundary. CR",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz\r"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza"),
						},
					},
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    False,
							E:    Probably,
						},
						PSB: ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
				pssu: &ParserServiceSub{
					PSSH: ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: ParserServiceSubBody{
						B: []byte("\r"),
					},
				},
			},
		},

		{
			name: "9. Full last boundary after begin piece",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + entities.Sep + "bPrefix" + "bRoot" + "bSuffix" + entities.Sep),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				last: true,
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb"),
						},
					},
				},
			},
		},

		{
			name: "10. Full last boundary after begin piece",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + entities.Sep + "bPrefix" + "bRoot" + "bSuffix" + entities.Sep),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				last: true,
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb"),
						},
					},
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    False,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
			},
		},

		{
			name: "11. Full boundary in the end",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + entities.Sep + "bPrefixbRoot"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb"),
						},
					},
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    False,
							E:    Probably,
						},
						PSB: ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
				pssu: &ParserServiceSub{
					PSSH: ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: ParserServiceSubBody{
						B: []byte(entities.Sep + "bPrefixbRoot"),
					},
				},
			},
		},

		{
			name: "12. Full boundary in the end with CR",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + entities.Sep + "bPrefixbRoot" + "\r"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb"),
						},
					},
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    False,
							E:    Probably,
						},
						PSB: ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
				pssu: &ParserServiceSub{
					PSSH: ParserServiceSubHeader{
						Part: 0,
						TS:   "qqq",
					},
					PSSB: ParserServiceSubBody{
						B: []byte(entities.Sep + "bPrefixbRoot" + "\r"),
					},
				},
			},
		},

		{
			name: "13. Partial last boundary with entities.Separated suffix after middle piece",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + entities.Sep + "bPrefix" + "bRoot" + "bSu"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				last: true,
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb"),
						},
					},
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    False,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
			},
		},

		{
			name: "14. Three full boundary no partial boundary",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazaz" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "c1234567890czczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczcz" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "d1234567890dzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazaz" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "c1234567890czczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczcz" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "d1234567890dzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazaz"),
						},
					},
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    False,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    False,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("c1234567890czczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczczcz"),
						},
					},
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    False,
							E:    True,
						},
						PSB: ParserServiceBody{
							B: []byte("d1234567890dzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdzdz"),
						},
					},
				},
			},
		},

		{
			name: "15. Partial last boundary after begin piece",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazaazazazazaazazazazaazazazazaazazazazaazazazazaazazazaza" + entities.Sep + "bPrefix" + "bRoot" + "bSuf"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazaazazazazaazazazazaazazazazaazazazazaazazazazaazazazaza"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				last: true,
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazaazazazazaazazazazaazazazazaazazazazaazazazazaazazazaza"),
						},
					},
				},
			},
		},

		{
			name: "16. Partial last boundary after middle piece",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazaz" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + entities.Sep + "bPrefix" + "bRoot" + "bSuf"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazaz" + entities.Sep + "bPrefix" + "bRoot" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				last: true,
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazazaz"),
						},
					},
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    False,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
						},
					},
				},
			},
		},

		{
			name: "17. Last part of last boundary",
			initDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("uffix" + entities.Sep),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 0,
				TS:   "qqq",
				Body: []byte("uffix" + entities.Sep),
				last: true,
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 0,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("uffix" + entities.Sep),
						},
					},
				},
			},
		},

		{
			name: "18. Intermediate part of last boundary",
			initDTO: &ParserServiceDTO{
				Part: 1,
				TS:   "qqq",
				Body: []byte("ootbSuffix" + entities.Sep),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 1,
				TS:   "qqq",
				Body: []byte("ootbSuffix" + entities.Sep),
				last: true,
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 1,
							TS:   "qqq",
							B:    True,
							E:    False,
						},
						PSB: ParserServiceBody{
							B: []byte("ootbSuffix" + entities.Sep),
						},
					},
				},
			},
		},

		{
			name: "19. No boundary",
			initDTO: &ParserServiceDTO{
				Part: 1,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			wantedDTO: &ParserServiceDTO{
				Part: 1,
				TS:   "qqq",
				Body: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
				Bou:  entities.Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				psus: []*ParserServiceUnit{
					{
						PSH: ParserServiceHeader{
							Part: 1,
							TS:   "qqq",
							B:    True,
							E:    True,
						},
						PSB: ParserServiceBody{
							B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + entities.Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
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
