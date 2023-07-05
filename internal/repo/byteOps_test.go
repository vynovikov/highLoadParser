package repo

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

type byteOpsSuite struct {
	suite.Suite
}

func TestByteOps(t *testing.T) {
	suite.Run(t, new(byteOpsSuite))
}

func (s *byteOpsSuite) TestLineRightEndIndexLimit() {
	bs := []byte("012345" + Sep)

	p := LineRightEndIndexLimit(bs, 0, 10)
	s.Equal(5, p)

}

func (s *byteOpsSuite) TestReverse() {
	bs := []byte("012345")

	bbs := Reverse(bs)

	s.Equal([]byte("543210"), bbs)
}

func (s *byteOpsSuite) TestFindBoundary() {
	bs := []byte("1111" + Sep + "2222" + Sep + "3333" + Sep + BoundaryField + "bRoot" + Sep + "4444" + Sep + "bPrefix" + "bRoot")

	boundary := FindBoundary(bs)

	s.True(cmp.Equal(boundary, Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}))
}

func (s *byteOpsSuite) TestGenBoundary() {
	boundaryVoc := Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}

	boundaryCalced := GenBoundary(boundaryVoc)

	s.Equal([]byte("\r\n"+"bPrefix"+"bRoot"), boundaryCalced)
}

func (s *byteOpsSuite) TestIsLastBoundaryEnding() {
	tt := []struct {
		name   string
		bs     []byte
		bou    Boundary
		wanted bool
	}{

		{
			name:   "b is part of suffix",
			bs:     []byte("fix"),
			bou:    Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wanted: true,
		},

		{
			name:   "b is suffix",
			bs:     []byte("bSuffix"),
			bou:    Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wanted: true,
		},

		{
			name:   "b has root ending and suffix",
			bs:     []byte("ootbSuffix"),
			bou:    Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wanted: true,
		},

		{
			name:   "b has wrong root ending and suffix",
			bs:     []byte("dotbSuffix"),
			bou:    Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wanted: false,
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			s.Equal(v.wanted, IsLastBoundaryEnding(v.bs, v.bou))
		})
	}
}

func (s *byteOpsSuite) TestSlicer() {

	tt := []struct {
		name             string
		afu              AppFeederUnit
		wantedDataPieces []DataPiece
	}{

		{
			name: "no full boundary and no partial boundary",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 1,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 1, B: True, E: True}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
			},
		},

		{
			name: "no full boundary, partial boundary present",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 1,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefixbRo"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 1, B: True, E: Probably}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				&AppSub{ASH: AppSubHeader{TS: "qqq", Part: 1}, ASB: AppSubBody{B: []byte("\r\nbPrefixbRo")}},
			},
		},
		{
			name: "no full boundary, last boundary",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefixbRootbSuffix" + Sep),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: Last}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
			},
		},

		{
			name: "no full boundary, CR in the end",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						Part: 1,
						TS:   "qqq",
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + "\r"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 1, B: True, E: Probably}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				&AppSub{ASH: AppSubHeader{TS: "qqq", Part: 1}, ASB: AppSubBody{B: []byte("\r")}},
			},
		},

		{
			name: "no full boundary, CRLF in the end",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						Part: 1,
						TS:   "qqq",
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + "\r\n"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 1, B: True, E: Probably}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				&AppSub{ASH: AppSubHeader{TS: "qqq", Part: 1}, ASB: AppSubBody{B: []byte("\r\n")}},
			},
		},

		{
			name: "1 full boundary",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "bPrefix" + "bRoot" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: False}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza")}},
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: False, E: True}, APB: AppPieceBody{B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
			},
		},

		{
			name: "1 full boundary, partial boundary present",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "bPrefix" + "bRoot" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefixbRo"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: False}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza")}},
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: False, E: Probably}, APB: AppPieceBody{B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				&AppSub{ASH: AppSubHeader{TS: "qqq", Part: 2}, ASB: AppSubBody{B: []byte("\r\nbPrefixbRo")}},
			},
		},

		{
			name: "1 full boundary, CR in the end",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza" + Sep + "bPrefix" + "bRoot" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + "\r"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: False}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazaza")}},
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: False, E: Probably}, APB: AppPieceBody{B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				&AppSub{ASH: AppSubHeader{TS: "qqq", Part: 2}, ASB: AppSubBody{B: []byte("\r")}},
			},
		},

		{
			name: "full last boundary after begin piece",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + Sep + "bPrefix" + "bRoot" + "bSuffix" + Sep),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: Last}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb")}},
			},
		},

		{
			name: "full last boundary after middle piece",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + Sep + "bPrefix" + "bRoot" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefix" + "bRoot" + "bSuffix" + Sep),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: False}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb")}},
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: False, E: Last}, APB: AppPieceBody{B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
			},
		},

		{
			name: "full boundary in the end",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + Sep + "bPrefix" + "bRoot" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefixbRoot"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: False}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb")}},
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: False, E: Probably}, APB: AppPieceBody{B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				&AppSub{ASH: AppSubHeader{TS: "qqq", Part: 2}, ASB: AppSubBody{B: []byte("\r\nbPrefixbRoot")}},
			},
		},

		{
			name: "full boundary in the end with CR",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + Sep + "bPrefix" + "bRoot" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefixbRoot" + "\r"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: False}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb")}},
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: False, E: Probably}, APB: AppPieceBody{B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				&AppSub{ASH: AppSubHeader{TS: "qqq", Part: 2}, ASB: AppSubBody{B: []byte("\r\nbPrefixbRoot\r")}},
			},
		},

		{
			name: "partial last boundary with separated suffix after middle piece",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb" + Sep + "bPrefix" + "bRoot" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefix" + "bRoot" + "bSu"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: False}, APB: AppPieceBody{B: []byte("a1234567890azazazazazazazazazazazazazazazazazazazazazazazazazazabzbzbzbzbzbzbzbzbzbzbzbzbzbzbzb")}},
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: False, E: Last}, APB: AppPieceBody{B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
			},
		},

		{
			name: "3 full lboundary no partial boundary",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefix" + "bRoot" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefix" + "bRoot" + Sep + "c1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefix" + "bRoot" + Sep + "d1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: False}, APB: AppPieceBody{B: []byte("a1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: False, E: False}, APB: AppPieceBody{B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: False, E: False}, APB: AppPieceBody{B: []byte("c1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: False, E: True}, APB: AppPieceBody{B: []byte("d1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
			},
		},

		{
			name: "partial last boundary after begin piece",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefix" + "bRoot" + "bSuf"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: Last}, APB: AppPieceBody{B: []byte("a1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
			},
		},
		{
			name: "partial last boundary after middle piece",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("a1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefix" + "bRoot" + Sep + "b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz" + Sep + "bPrefix" + "bRoot" + "bSuf"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: False}, APB: AppPieceBody{B: []byte("a1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: False, E: Last}, APB: AppPieceBody{B: []byte("b1234567890bzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbzbz")}},
			},
		},

		{
			name: "last part of last boundary",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("uffix"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: Last}},
			},
		},

		{
			name: "part of last boundary",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("ootbSuffix"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: Last}},
			},
		},

		{
			name: "last boundary suffix only",
			afu: AppFeederUnit{
				R: ReceiverUnit{
					H: ReceiverHeader{
						TS:   "qqq",
						Part: 2,
						Bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
					},
					B: ReceiverBody{
						B: []byte("bSuffix"),
					},
				},
			},
			wantedDataPieces: []DataPiece{
				&AppPieceUnit{APH: AppPieceHeader{TS: "qqq", Part: 2, B: True, E: Last}},
			},
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			s.Equal(v.wantedDataPieces, Slicer(v.afu))

		})
	}
}

func (s *byteOpsSuite) TestSingleLineRightTrimmed() {
	tt := []struct {
		name      string
		bs        []byte
		limit     int
		wantC     int
		wantL     []byte
		wantError error
	}{
		{
			name:      "err zero length",
			bs:        []byte(""),
			limit:     3,
			wantC:     0,
			wantL:     nil,
			wantError: errors.New("passed byte slice with zero length"),
		},
		{
			name:      "err no actual characters",
			bs:        []byte("\r\n\r\n\r\n\r\n"),
			limit:     3,
			wantC:     0,
			wantL:     nil,
			wantError: errors.New("no actual characters before limit"),
		},
		{
			name:      "happy without Sep",
			bs:        []byte("11111" + Sep + "222222" + Sep + "33333333333" + Sep + "444444444444"),
			limit:     10,
			wantL:     []byte("11111"),
			wantError: nil,
		},
		{
			name:      "happy with Sep",
			bs:        []byte("\r\n11111" + Sep + "222222" + Sep + "33333333333" + Sep + "444444444444"),
			limit:     12,
			wantC:     2,
			wantL:     []byte("11111"),
			wantError: nil,
		},
		{
			name:      "happy last boundary part",
			bs:        []byte("-"),
			limit:     12,
			wantL:     []byte("-"),
			wantError: nil,
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			c, l, err := SingleLineRightTrimmed(v.bs, v.limit)

			s.Equal(v.wantL, l)
			s.Equal(v.wantC, c)

			if err != nil {
				s.Equal(v.wantError, err)
			} else {
				s.NoError(err)
			}

		})
	}
}

func (s *byteOpsSuite) TestSingleLineRight() {
	bs := []byte("11111" + Sep + "222222" + Sep + "33333333333" + Sep + "444444444444")
	l, err := SingleLineRightUnchanged(bs, 28)

	s.NoError(err)
	s.Equal([]byte("11111"), l)
}

func (s *byteOpsSuite) TestIsPartlyBoundaryRight() {
	bs := []byte("oot")
	boundary := Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}

	s.True(IsPartlyBoundaryRight(bs, boundary))
}

func (s *byteOpsSuite) TestNoDigits() {
	tt := []struct {
		name string
		bs   []byte
	}{
		{
			name: "absent",
			bs:   []byte("afahf _=\""),
		},
		{
			name: "present",
			bs:   []byte("sadjfd345_\""),
		},
	}
	for _, tc := range tt {
		s.Run(tc.name, func() {
			if tc.name == "absent" {
				s.True(NoDigits(tc.bs))
			}
			if tc.name == "present" {
				s.False(NoDigits(tc.bs))
			}
		})
	}
}

func (s *byteOpsSuite) TestAllPrintalbe() {
	tt := []struct {
		name string
		bs   []byte
		want bool
	}{
		{
			name: "all printable",
			bs:   []byte("afahf _="),
			want: true,
		},
		{
			name: "have at least 1 NonPrintable",
			bs:   []byte("sadjfd345_\r777"),
			want: false,
		},
	}
	for _, tc := range tt {
		s.Run(tc.name, func() {
			if tc.want {
				s.True(AllPrintalbe(tc.bs))
			}
			if !tc.want {
				s.False(AllPrintalbe(tc.bs))
			}

		})
	}
}
func (s *byteOpsSuite) TestCurrentLineFirstPrintIndexLeft() {
	tt := []struct {
		name      string
		bs        []byte
		wantValue int
		wantError error
	}{
		{
			name:      "happy",
			bs:        []byte("12345\r\n"),
			wantValue: 4,
			wantError: nil,
		},
		{
			name:      "unhappy no printable",
			bs:        []byte("\n\r\n\r\n\r\n\r\r\n"),
			wantValue: -1,
			wantError: errors.New("in repo.CurrentLineFirstPrintIndexLeft no actual characters before limit"),
		},
		{
			name:      "unhappy zero lenght",
			bs:        []byte(""),
			wantValue: -1,
			wantError: errors.New("in repo.CurrentLineFirstPrintIndexLeft passed byte slice with zero length"),
		},
	}

	for _, v := range tt {
		s.Run(v.name, func() {
			i, err := CurrentLineFirstPrintIndexLeft(v.bs, len(v.bs)-2)
			if err != nil {
				s.Equal(v.wantError, err)

			} else {
				s.NoError(err)
			}
			s.Equal(v.wantValue, i)
		})
	}

}

func (s *byteOpsSuite) TestGetCurrentLineLeft() {
	tt := []struct {
		name      string
		bs        []byte
		fromIndex int
		limit     int
		wantValue []byte
		wantError error
	}{
		{
			name:      "happy",
			bs:        []byte("\r\n12345\r\n"),
			fromIndex: 6,
			limit:     9,
			wantValue: []byte("12345"),
			wantError: nil,
		},
		{
			name:      "unhappy limit exceeded",
			bs:        []byte("\r\n0123456789\r\n"),
			fromIndex: 11,
			limit:     8,
			wantValue: nil,
			wantError: errors.New("in repo.GetCurrentLineLeft line limit exceeded. No separator met"),
		},
		{
			name:      "unhappy zero length",
			bs:        []byte(""),
			fromIndex: 11,
			limit:     8,
			wantValue: []byte{},
			wantError: errors.New("in repo.GetCurrentLineLeft passed byte slice with zero length"),
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			got, err := GetCurrentLineLeft(v.bs, v.fromIndex, v.limit)
			if err != nil {
				s.Equal(v.wantError, err)
			} else {
				s.NoError(err)
			}
			s.Equal(v.wantValue, got)
		})
	}
}
func (s *byteOpsSuite) TestSingleLineLeftTrimmed() {
	tt := []struct {
		name      string
		bs        []byte
		limit     int
		wantValue []byte
		wantError error
	}{
		{
			name:      "happy 1-character separator",
			bs:        []byte("\r\n12345\r"),
			limit:     9,
			wantValue: []byte("12345"),
			wantError: nil,
		},
		{
			name:      "happy 2-character separator",
			bs:        []byte("\r\n12345\r\n"),
			limit:     9,
			wantValue: []byte("12345"),
			wantError: nil,
		},
		{
			name:      "unhappy limit exceeded",
			bs:        []byte("\r\n0123456789\r\n"),
			limit:     8,
			wantValue: nil,
			wantError: errors.New("in repo.GetCurrentLineLeft line limit exceeded. No separator met"),
		},
		{
			name:      "unhappy zero length",
			bs:        []byte{},
			limit:     8,
			wantValue: nil,
			wantError: errors.New("in repo.CurrentLineFirstPrintIndexLeft passed byte slice with zero length"),
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			got, err := SingleLineLeftTrimmed(v.bs, v.limit)
			if err != nil {
				s.Equal(v.wantError, err)
			} else {
				s.NoError(err)
			}
			s.Equal(v.wantValue, got)
		})
	}
}
func (s *byteOpsSuite) TestGetLinesLeft() {
	tt := []struct {
		name      string
		bs        []byte
		limit     int
		voc       Vocabulaty
		wantValue [][]byte
		wantError error
	}{
		{
			name:  "happy 1 line separated",
			bs:    []byte("\r\nbPrefi"),
			limit: 9,
			voc:   Vocabulaty{Boundary: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}, CType: "Content-Type"},
			wantValue: [][]byte{
				[]byte("bPrefi"),
			},
			wantError: nil,
		},
		{
			name:  "happy 2 line separated",
			bs:    []byte("\r\nbPrefixbRoot\r\nContent-Disposit"),
			limit: 40,
			voc:   Vocabulaty{Boundary: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}, CDisposition: "Content-Disposition", CType: "Content-Type"},
			wantValue: [][]byte{
				[]byte("Content-Disposit"),
				[]byte("bPrefixbRoot"),
			},
			wantError: nil,
		},
		{
			name:  "happy 3 line separated",
			bs:    []byte("\r\nbPrefixbRoot\r\nContent-Disposition :\r\nContent-Type"),
			limit: 60,
			voc:   Vocabulaty{Boundary: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}, CDisposition: "Content-Disposition", CType: "Content-Type"},
			wantValue: [][]byte{
				[]byte("Content-Type"),
				[]byte("Content-Disposition :"),
				[]byte("bPrefixbRoot"),
			},
			wantError: nil,
		},
		{
			name:  "happy third line separator separated",
			bs:    []byte("\r\nbPrefixbRoot\r\nContent-Disposition :\r\nContent-Type :\r"),
			limit: 60,
			voc:   Vocabulaty{Boundary: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}, CDisposition: "Content-Disposition", CType: "Content-Type"},
			wantValue: [][]byte{
				[]byte("Content-Type :"),
				[]byte("Content-Disposition :"),
				[]byte("bPrefixbRoot"),
			},
			wantError: nil,
		},
		{
			name:  "happy second line separator separated",
			bs:    []byte("\r\nbPrefixbRoot\r\nContent-Disposition :\r"),
			limit: 60,
			voc:   Vocabulaty{Boundary: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}, CDisposition: "Content-Disposition", CType: "Content-Type"},
			wantValue: [][]byte{
				[]byte("Content-Disposition :"),
				[]byte("bPrefixbRoot"),
			},
			wantError: nil,
		},
		{
			name:  "happy first line separator separated",
			bs:    []byte("\r\nbPrefixbRoot\r"),
			limit: 60,
			voc:   Vocabulaty{Boundary: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}, CDisposition: "Content-Disposition", CType: "Content-Type"},
			wantValue: [][]byte{
				[]byte("bPrefixbRoot"),
			},
			wantError: nil,
		},
	}

	for _, v := range tt {
		s.Run(v.name, func() {
			got, err := GetLinesLeft(v.bs, v.limit, v.voc)
			if err != nil {
				s.Equal(v.wantError, err)
			} else {
				s.NoError(err)
			}
			s.Equal(v.wantValue, got)
		})
	}
}
func (s *byteOpsSuite) TestGetLinesRightMiddle() {
	tt := []struct {
		name      string
		bs        []byte
		limit     int
		voc       Vocabulaty
		wantValue [][]byte
		wantError error
	}{
		{
			name:  "happy 1 line not full",
			bs:    []byte("Content-Disposition: form-data; name=\"claire\"; filename=\"sho"),
			limit: MaxHeaderLimit,
			voc:   Vocabulaty{CDisposition: "Content-Disposition", CType: "Content-Type", FormName: "name=\"", FileName: " filename=\""},
			wantValue: [][]byte{
				[]byte("Content-Disposition: form-data; name=\"claire\"; filename=\"sho"),
			},
			wantError: errors.New("in GetLinesRightMiddle header \"Content-Disposition: form-data; name=\"claire\"; filename=\"sho\" is not full"),
		},

		{
			name:      "unhappy 1 line only",
			bs:        []byte("Content-Disposition: form-data; name=\"claire\"\r\naaaaaaaaaaaaassssssssssss\r\n"),
			limit:     MaxHeaderLimit,
			voc:       Vocabulaty{CDisposition: "Content-Disposition", CType: "Content-Type", FormName: "name=\"", FileName: " filename=\""},
			wantValue: [][]byte{},
			wantError: errors.New("in GetLinesRightMiddle second line \"aaaaaaaaaaaaassssssssssss\" is unexpected"),
		},

		{
			name:      "unhappy 1 line unexpected",
			bs:        []byte("Content-Disposition: form-data; name=\"cla ire\"\r"),
			limit:     MaxHeaderLimit,
			voc:       Vocabulaty{CDisposition: "Content-Disposition", CType: "Content-Type", FormName: "name=\"", FileName: " filename=\""},
			wantValue: [][]byte{},
			wantError: errors.New("in GetLinesRightMiddle first line \"Content-Disposition: form-data; name=\"cla ire\"\" is unexpected"),
		},

		{
			name:  "happy 2 lines not full",
			bs:    []byte("Content-Disposition: form-data; name=\"claire\"; filename=\"short.txt\"\r\nContent-Type: text/pla"),
			limit: MaxHeaderLimit,
			voc:   Vocabulaty{CDisposition: "Content-Disposition", CType: "Content-Type", FormName: "name=\"", FileName: " filename=\""},
			wantValue: [][]byte{
				[]byte("Content-Disposition: form-data; name=\"claire\"; filename=\"short.txt\""),
				[]byte("Content-Type: text/pla"),
			},
			wantError: errors.New("in GetLinesRightMiddle header \"Content-Disposition: form-data; name=\"claire\"; filename=\"short.txt\"Content-Type: text/pla\" is not full"),
		},

		{
			name:      "unhappy 2 lines unexpected",
			bs:        []byte("Content-Disposition: form-data; name=\"claire\"; filename=\"short.txt\"\r\nCnotent-Type: text"),
			limit:     MaxHeaderLimit,
			voc:       Vocabulaty{CDisposition: "Content-Disposition", CType: "Content-Type", FormName: "name=\"", FileName: " filename=\""},
			wantValue: [][]byte{},
			wantError: errors.New("in GetLinesRightMiddle second line \"Cnotent-Type: text\" is unexpected"),
		},

		{
			name:  "happy 2 lines full",
			bs:    []byte("Content-Disposition: form-data; name=\"claire\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazazaa"),
			limit: MaxHeaderLimit,
			voc:   Vocabulaty{CDisposition: "Content-Disposition", CType: "Content-Type", FormName: "name=\"", FileName: " filename=\""},
			wantValue: [][]byte{
				[]byte("Content-Disposition: form-data; name=\"claire\"; filename=\"short.txt\""),
				[]byte("Content-Type: text/plain"),
			},
			wantError: nil,
		},
	}

	for _, v := range tt {
		s.Run(v.name, func() {
			lines, _, err := GetLinesRightMiddle(v.bs, v.limit)
			if err != nil {
				s.Equal(v.wantError, err)
			}
			s.Equal(v.wantValue, lines)
		})
	}
}

func (s *byteOpsSuite) TestGetLinesRightBegin() {
	tt := []struct {
		name      string
		bs        []byte
		limit     int
		bou       Boundary
		wantValue [][]byte
		wantError error
	}{
		{
			name:  "happy 1 line",
			bs:    []byte("Content-Type: text/plain\r\n"),
			limit: MaxHeaderLimit,
			bou:   Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantValue: [][]byte{
				[]byte("Content-Type: text/plain"),
			},
			wantError: nil,
		},

		{
			name:      "unhappy 1 line",
			bs:        []byte("Content-Typp: text/plain\r\n"),
			limit:     MaxHeaderLimit,
			bou:       Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantValue: [][]byte{},
			wantError: errors.New("first line \"Content-Typp: text/plain\" is unexpected"),
		},
		{
			name:  "happy 2 lines",
			bs:    []byte("rm-data; name=\"claire\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"),
			limit: MaxHeaderLimit,
			bou:   Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantValue: [][]byte{
				[]byte("rm-data; name=\"claire\"; filename=\"short.txt\""),
				[]byte("Content-Type: text/plain"),
			},
			wantError: nil,
		},

		{
			name:      "unhappy_1 2 lines",
			bs:        []byte("rm-data; nsme=\"claire\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"),
			limit:     MaxHeaderLimit,
			bou:       Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantValue: [][]byte{},
			wantError: errors.New("first line \"rm-data; nsme=\"claire\"; filename=\"short.txt\"\" is unexpected"),
		},

		{
			name:      "unhappy_2 2 lines",
			bs:        []byte("rm-data; name=\"claire\"; filename=\"short.txt\"\r\nCintent-Type: text/plain\r\n"),
			limit:     MaxHeaderLimit,
			bou:       Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantValue: [][]byte{},
			wantError: errors.New("second line \"Cintent-Type: text/plain\" is unexpected"),
		},

		{
			name:      "unhappy_3 2 lines",
			bs:        []byte("rm-data; name=\"claire\"; filename=\"short.txt\"\r\nContent-Type: text:plain\r\n"),
			limit:     MaxHeaderLimit,
			bou:       Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantValue: [][]byte{},
			wantError: errors.New("second line \"Content-Type: text:plain\" is unexpected"),
		},

		{
			name:  "happy 3 lines",
			bs:    []byte("oot\r\nContent-Disposition: form-data; name=\"claire\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"),
			limit: MaxHeaderLimit,
			bou:   Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantValue: [][]byte{
				[]byte("oot"),
				[]byte("Content-Disposition: form-data; name=\"claire\"; filename=\"short.txt\""),
				[]byte("Content-Type: text/plain"),
			},
			wantError: nil,
		},

		{
			name:      "unhappy unexpected boundary ending",
			bs:        []byte("abc\r\nContent-Disposition: form-data; name=\"claire\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"),
			limit:     MaxHeaderLimit,
			bou:       Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantValue: [][]byte{},
			wantError: errors.New("first line \"abc\" is unexpected"),
		},

		{
			name:      "unhappy unexpected Content-Disposition line ending",
			bs:        []byte("ion: form-data; name=\"claire\"; filename=short.txt\r\nContent-Type: text/plain\r\n"),
			limit:     MaxHeaderLimit,
			bou:       Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantValue: [][]byte{},
			wantError: errors.New("first line \"ion: form-data; name=\"claire\"; filename=short.txt\" is unexpected"),
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			lines, _, err := GetLinesRightBegin(v.bs, v.limit, v.bou)
			if err != nil {
				s.Equal(v.wantError, err)
			}
			s.Equal(v.wantValue, lines)
		})
	}
}

func (s *byteOpsSuite) TestPartlyBoundaryLeft() {

	tt := []struct {
		name string
		bs   []byte
		bou  Boundary
		want []byte
	}{
		{
			name: "bPrefix separated",
			bs:   []byte("1111111111" + Sep + "bPref"),
			bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			want: []byte("bPref"),
		},
		{
			name: "bRoot separated",
			bs:   []byte("1111111111" + Sep + "bPrefix" + "bRo"),
			bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			want: []byte("bPrefix" + "bRo"),
		},
		{
			name: "bSuffix separated",
			bs:   []byte("1111111111" + Sep + "bPrefix" + "bRoot" + "bSuf"),
			bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			want: []byte("bPrefix" + "bRoot" + "bSuf"),
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			got, err := PartlyBoundaryLeft(v.bs, v.bou)

			s.NoError(err)
			s.Equal(v.want, got)
		})

	}

}

func (s *byteOpsSuite) TestPartlyBoundaryRight() {
	tt := []struct {
		name  string
		bs    []byte
		limit int
		want  []byte
	}{
		{
			name:  "happy last bPrefix separated",
			bs:    []byte("efix" + "bRoot" + "bSuffix"),
			limit: len("efix" + "bRoot" + "bSuffix"),
			want:  []byte("efix" + "bRoot" + "bSuffix"),
		},
		{
			name:  "happy last bRoot separated ",
			bs:    []byte("oot" + "bSuffix"),
			limit: len("oot" + "bSuffix"),
			want:  []byte("oot" + "bSuffix"),
		},
		{
			name:  "happy not last bRoot separated ",
			bs:    []byte("oot" + Sep + "111111111"),
			limit: 8,
			want:  []byte("oot"),
		},
		{
			name:  "happy last bSuffix separated",
			bs:    []byte("fix"),
			limit: len("fix"),
			want:  []byte("fix"),
		},
	}

	for _, v := range tt {
		s.Run(v.name, func() {
			got, err := PartlyBoundaryRight(v.bs, v.limit)
			s.NoError(err)
			s.Equal(v.want, got)
		})
	}
}

func (s *byteOpsSuite) TestLineRightLimit() {

	tt := []struct {
		name      string
		bs        []byte
		fromIndex int
		limit     int
		want      []byte
	}{
		{
			name:      "happy",
			bs:        []byte("11111" + Sep + "222222" + Sep + "33333333333" + Sep + "444444444444"),
			fromIndex: 7,
			limit:     7,
			want:      []byte("222222"),
		},
		{
			name:      "unhappy",
			bs:        []byte("11111" + Sep + "222222" + Sep + "33333333333" + Sep + "444444444444"),
			fromIndex: 7,
			limit:     4,
			want:      nil,
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			l := LineRightLimit(v.bs, v.fromIndex, v.limit)
			s.Equal(v.want, l)
		})
	}

}

func (s *byteOpsSuite) TestGetLastLine() {
	tt := []struct {
		name      string
		bs        []byte
		boundary  []byte
		wantValue []byte
	}{
		{
			name:      "CR after boundary",
			bs:        []byte("11111111111111111" + Sep + "bPrefixbRoot" + "\r"),
			boundary:  []byte("bPrefix" + "bRoot"),
			wantValue: []byte(Sep + "bPrefixbRoot" + "\r"),
		},

		{
			name:      "CR after random line",
			bs:        []byte("11111111111111111" + Sep + "22222222" + "\r"),
			boundary:  []byte("bPrefix" + "bRoot"),
			wantValue: []byte(Sep + "22222222" + "\r"),
		},
		{
			name:      "happy CRLF",
			bs:        []byte("11111111111111111" + Sep + "2222222222222222" + "\r\n"),
			boundary:  []byte("bPrefix" + "bRoot"),
			wantValue: []byte("\r\n"),
		},
		{
			name:      "happy default",
			bs:        []byte("11111111111111111" + Sep + "2222222222222222" + "\r\n" + "3333333333"),
			boundary:  []byte("bPrefix" + "bRoot"),
			wantValue: []byte("\r\n" + "3333333333"),
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			b := GetLastLine(v.bs, v.boundary)
			s.Equal(v.wantValue, b)
		})
	}
}

func (s *byteOpsSuite) TestWordRightBorderLimit() {
	tt := []struct {
		name        string
		bs          []byte
		beg         []byte
		end         []byte
		limit       int
		wantedWord  []byte
		wantedError error
	}{
		{
			name:        "happy name occ",
			bs:          []byte("Content-Disposition: form-data; name=\"claire\"; filename=\"short.txt\""),
			beg:         []byte("name=\""),
			end:         []byte("\""),
			limit:       25,
			wantedWord:  []byte("claire"),
			wantedError: errors.New(""),
		},

		{
			name:        "haooy filename occ",
			bs:          []byte("Content-Disposition: form-data; name=\"claire\"; filename=\"short.txt\""),
			beg:         []byte("filename=\""),
			end:         []byte("\""),
			limit:       25,
			wantedWord:  []byte("short.txt"),
			wantedError: errors.New(""),
		},

		{
			name:        "unhappy beginning not found",
			bs:          []byte("Content-Disposition: form-data; name=\"claire\"; filename=\"short.txt\""),
			beg:         []byte("silename=\""),
			end:         []byte("\""),
			limit:       25,
			wantedWord:  []byte(""),
			wantedError: errors.New("beginning not found"),
		},

		{
			name:        "unhaooy limit exceeded",
			bs:          []byte("Content-Disposition: form-data; name=\"claire\"; filename=\"short.txt\""),
			beg:         []byte("filename=\""),
			end:         []byte("1"),
			limit:       3,
			wantedWord:  []byte(""),
			wantedError: errors.New("limit exceeded"),
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			line, err := WordRightBorderLimit(v.bs, v.beg, v.end, v.limit)
			if err != nil {
				s.Equal(v.wantedError, err)
			}
			s.Equal(v.wantedWord, line)
		})
	}
}

func (s *byteOpsSuite) TestRepeatedIntex() {
	tt := []struct {
		name      string
		bs        []byte
		occ       []byte
		i         int
		wantIndex int
	}{
		{
			name:      "happy first",
			bs:        []byte("1aa234aa567aa890aa111"),
			occ:       []byte("aa"),
			i:         1,
			wantIndex: 1,
		},
		{
			name:      "happy second",
			bs:        []byte("1aa234aa567aa890aa111"),
			occ:       []byte("aa"),
			i:         2,
			wantIndex: 6,
		},
		{
			name:      "happy third",
			bs:        []byte("1aa234aa567aa890aa111"),
			occ:       []byte("aa"),
			i:         3,
			wantIndex: 11,
		},

		{
			name:      "happy fourth",
			bs:        []byte("1aa234aa567aa890aa111"),
			occ:       []byte("aa"),
			i:         4,
			wantIndex: 16,
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			s.Equal(v.wantIndex, RepeatedIntex(v.bs, v.occ, v.i))
		})

	}
}

func (s *byteOpsSuite) TestEndingOf() {
	tt := []struct {
		name   string
		long   []byte
		short  []byte
		wanted bool
	}{
		{
			name:   "happy",
			long:   []byte("1234567890"),
			short:  []byte("67890"),
			wanted: true,
		},

		{
			name:   "unhappy",
			long:   []byte("1234567890"),
			short:  []byte("68790"),
			wanted: false,
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			s.Equal(v.wanted, EndingOf(v.long, v.short))
		})
	}
}

func (s *byteOpsSuite) TestBoundaryPartInLastLine() {
	tt := []struct {
		name        string
		bs          []byte
		bou         Boundary
		wantedL     []byte
		wantedError error
	}{

		{
			name:        "CR in the end",
			bs:          []byte("sdfdsfdsf\r"),
			wantedL:     []byte("\r"),
			wantedError: nil,
		},

		{
			name:        "CRLF in the end",
			bs:          []byte("sdfdsfdsf\r\n"),
			wantedL:     []byte("\r\n"),
			wantedError: nil,
		},

		{
			name:        "random in the end",
			bs:          []byte("\r\n1111111111"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("in repo.BoundaryPartInLastLine no boundary"),
		},

		{
			name:        "boundary part in the end",
			bs:          []byte("azaza\r\nbPrefixb"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\nbPrefixb"),
			wantedError: nil,
		},

		{
			name:        "boundary-like line",
			bs:          []byte("azaza\r\nifx"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("in repo.BoundaryPartInLastLine no boundary"),
		},

		{
			name:        "boundary + CR",
			bs:          []byte("\r\nbPrefixbRoot\r"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\nbPrefixbRoot\r"),
			wantedError: nil,
		},

		{
			name:        "last boundary part",
			bs:          []byte("azaza\r\nbPrefixbRootbSu"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\nbPrefixbRootbSu"),
			wantedError: nil,
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			got, err := BoundaryPartInLastLine(v.bs, v.bou)
			if err != nil {
				s.Equal(v.wantedError, err)
			}
			s.Equal(v.wantedL, got)
		})
	}
}

func (s *byteOpsSuite) TestGetHeaderLines() {
	tt := []struct {
		name        string
		bs          []byte
		bou         Boundary
		wantedL     []byte
		wantedError error
	}{
		/*
			{
				name:        "0 CRLF <1 line right",
				bs:          []byte("Content-Disposition: form-data; name=\"al"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"al"),
				wantedError: errors.New("in repo.GetHeaderLines header \"Content-Disposition: form-data; name=\"al\" is not full"),
			},

			{
				name:        "0 CRLF no header",
				bs:          []byte("azazazazazaza"),
				wantedL:     nil,
				wantedError: errors.New("in repo.GetHeaderLines no header found"),
			},

			{
				name:        "1 CRLF whole sufficient 1-st line",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"Content-Disposition: form-data; name=\"alice\"\r\n\" is not full"),
			},

			{
				name:        "1 CRLF whole 1-st line, partyal 2-d",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short\"\r\nCon"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short\"\r\nCon"),
				wantedError: errors.New("in repo.GetHeaderLines header \"Content-Disposition: form-data; name=\"alice\"; filename=\"short\"\r\nCon\" is not full"),
			},
			{
				name:        "1 CRLF just CRLF, random line",
				bs:          []byte("\r\nr23hjrb23hrbj23hbrh23"),
				wantedL:     []byte("\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\r\n\" is ending part"),
			},

			{
				name:        "1 CRLF last boundary",
				bs:          []byte("azzsdfgsdhfdsfhsjdfhs\r\n"),
				wantedL:     []byte("azzsdfgsdhfdsfhsjdfhs\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"azzsdfgsdhfdsfhsjdfhs\r\n\" is the last"),
			},

			{
				name:        "1 CRLF no header_2",
				bs:          []byte("azzsdfgsdhfdsfhsjdfhs\r\nfskjfghsjfhgfjkhgjdfhgfd"),
				wantedL:     nil,
				wantedError: errors.New("in repo.GetHeaderLines no header found"),
			},

			{
				name:        "2 CRLF 1 line CD sufficient",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
				wantedError: nil,
			},

			{
				name:        "2 CRLF 1 line CD sufficient + random line",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nhdsghdsvhsdvgshdgvsdv"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
				wantedError: nil,
			},

			{
				name:        "2 CRLF 1 line CD insufficient + CRLF + CT + CTLF",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\" is not full"),
			},

			{
				name:        "2 CRLF 1 line CD sufficient + random line",
				bs:          []byte("position: form-data; name=\"alice\"\r\n\r\nhdsghdsvhsdvgshdgvsdv"),
				wantedL:     []byte("position: form-data; name=\"alice\"\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"position: form-data; name=\"alice\"\r\n\r\n\" is ending part"),
			},

			{
				name:        "2 CRLF 1 header line right 2 random lines",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nsajkfdga\r\ndsfguigdfa"),
				wantedL:     nil,
				wantedError: errors.New("in repo.GetHeaderLines no header found"),
			},

			{
				name:        "2 CRLF all lines random",
				bs:          []byte("we6fwfef6gewfgewfg7efge\r\nsajkfdga\r\ndsfguigdfa"),
				wantedL:     nil,
				wantedError: errors.New("in repo.GetHeaderLines no header found"),
			},

			{
				name:        "2 CRLF CRLF, random line, CRLF, random line",
				bs:          []byte("\r\n2f3hg4f32ghf423gf324\r\nr23hjrb23hrbj23hbrh23"),
				wantedL:     []byte("\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\r\n\" is ending part"),
			},

			{
				name:        "2 CRLF just 2 * CRLF, random line",
				bs:          []byte("\r\n\r\nr23hjrb23hrbj23hbrh23"),
				wantedL:     []byte("\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\r\n\r\n\" is ending part"),
			},

			{
				name:        "3 CRLF 2 header lines (CD insufficient), 1 random line",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\ndsfguigdfa"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
				wantedError: nil,
			},

			{
				name:        "3 CRLF 1 header line (CD sufficient), 2 random lines",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\ndsfguigdfa6fhgf55\r\nggf8723g723gf823"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"Content-Disposition: form-data; name=\"alice\"\r\n\r\n\" is ending part"),
			},

			{
				name:        "3 CRLF: CRLF + CDsuf + 2*CRLF + rand",
				bs:          []byte("\r\nContent-Disposition: form-data; name=\"alice\"\r\n\r\ndsfguigdfa"),
				wantedL:     []byte("\r\nContent-Disposition: form-data; name=\"alice\"\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\r\nContent-Disposition: form-data; name=\"alice\"\r\n\r\n\" is ending part"),
			},

			{
				name:        "3 CRLF: CRLF + CT + 2*CRLF + rand",
				bs:          []byte("\r\nContent-Type: text/plain\r\n\r\ndsfguigdfa"),
				wantedL:     []byte("\r\nContent-Type: text/plain\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\r\nContent-Type: text/plain\r\n\r\n\" is ending part"),
			},

			{
				name:        "3 CRLF 1 header line, 2 random lines",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nCshdgfhsdgfhsdjf\r\ndsfguigdfa"),
				wantedL:     nil,
				wantedError: errors.New("in repo.GetHeaderLines no header found"),
			},

			{
				name:        "3 CRLF all random lines",
				bs:          []byte("azazzazazzazazaz\r\nCzbbzbzbbzbzbbzbzbzbzbz\r\ndsfguigdfa\r\nf2r7fr27fr2f7r2"),
				wantedL:     nil,
				wantedError: errors.New("in repo.GetHeaderLines no header found"),
			},

			{
				name:        "3 CRLF: CRLF, next random lines",
				bs:          []byte("\r\nr23hjrb23hrbj23hbrh23\r\nsgdhgsdwef6fr6632\r\n438ry34grg438rg438gr43"),
				wantedL:     []byte("\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\r\n\" is ending part"),
			},

			{
				name:        "3 CRLF just  2 * CRLF, next random lines",
				bs:          []byte("\r\n\r\nr23hjrb23hrbj23hbrh23\r\nsgdhgsdwef6fr6632"),
				wantedL:     []byte("\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\r\n\r\n\" is ending part"),
			},

			{
				name:        "4 CRLF 1 line CD sufficient + 3 random lines",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nhdsghdsvhsdvgshdgvsdv\r\nhjgvfjhdgvjhfdkgftv87dfvdfv\r\nsoiehfwoefhwefdgvjhsdv"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"Content-Disposition: form-data; name=\"alice\"\r\n\r\n\" is ending part"),
			},

			{
				name:        "4 CRLF 1 line CD insufficient + 2 random line",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nhdsghdsvhsdvgshdgvsdv\r\nhjgvfjhdgvjhfdkgftv87dfvdfv"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
				wantedError: nil,
			},

			{
				name:        "4 CRLF: CRLF + CDinsuf + CRLF + CT + 2*CRLF + rand",
				bs:          []byte("\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\ndsfguigdfa"),
				wantedL:     []byte("\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n\" is ending part"),
			},

			{
				name:        "4 CRLF 1 boundary ending 2 header lines, 1 random line",
				bs:          []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\ndsfguigdfa"),
				bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				wantedL:     []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n\" is ending part"),
			},

			{
				name:        "5 CRLF 1 boundary ending 2 header lines, 1 random line",
				bs:          []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\naaaaaaaaaaaaaaaaaaaaaaaaa\r\nbbbbbbbbbbbbbbbb\r\ndsfguigdfa"),
				bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				wantedL:     []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n\" is ending part"),
			},

			{
				name:        "4 CRLF 1 boundary ending 1 header line, 2 random lines",
				bs:          []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\ndsfguigdfa\r\nf2r7fr27fr2f7r2"),
				bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				wantedL:     nil,
				wantedError: errors.New("in repo.GetHeaderLines no header found"),
			},

			{
				name:        "4 CRLF 1 boundary ending rest random lines",
				bs:          []byte("fixbRoot\r\nCzbbzbzbbzbzbbzbzbzbzbz\r\ndsfguigdfa\r\nf2r7fr27fr2f7r2"),
				bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				wantedL:     nil,
				wantedError: errors.New("in repo.GetHeaderLines no header found"),
			},

			{
				name:        "4 CRLF just  2 * CRLF, next random lines",
				bs:          []byte("\r\n\r\nr23hjrb23hrbj23hbrh23\r\nsgdhgsdwef6fr6632\r\n3fd72fd73fd3727df23"),
				wantedL:     []byte("\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\r\n\r\n\" is ending part"),
			},

			{
				name:        "5 CRLF 1 CRLF 2 header lines, 1 random line",
				bs:          []byte("\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\naaaaaaaaaaaaaaaaaaaaaaaaa\r\nbbbbbbbbbbbbbbbb\r\ndsfguigdfa"),
				bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				wantedL:     []byte("\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n\" is ending part"),
			},

			{
				name:        "Precending LF, 0 CRLF. LF + rand",
				bs:          []byte("\nsdjkchdjhcskdhcdsjhckjsdhcjdsk"),
				wantedL:     []byte("\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\n\" is ending part"),
			},

			{
				name:        "Precending LF, 3 CRLF. LF + rand",
				bs:          []byte("\nsdjkchdjhcskdhcdsjhckjsdhcjdsk\r\nsdjhfjdshjfsd\r\ngruihgeruhguerhguerg\r\n121312j412jk4g1jk4gjkg"),
				bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				wantedL:     []byte("\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\n\" is ending part"),
			},

			{
				name:        "Precending LF, 1 CRLF. CRLF + LF + rand",
				bs:          []byte("\n\r\nsdjkchdjhcskdhcdsjhckjsdhcjdsk"),
				wantedL:     []byte("\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\n\r\n\" is ending part"),
			},

			{
				name:        "Precending LF, 2 CRLF. LF + CT + 2*CRLF + rand",
				bs:          []byte("\nContent-Type: text/plain\r\n\r\nsdjkchdjhcskdhcdsjhckjsdhcjdsk"),
				wantedL:     []byte("\nContent-Type: text/plain\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\nContent-Type: text/plain\r\n\r\n\" is ending part"),
			},

			{
				name:        "Precending LF, 2 CRLF. LF + CDSuff + 2*CRLF + rand",
				bs:          []byte("\nContent-Disposition: form-data; name=\"alice\"\r\n\r\nsdjkch2323232djhcskdhcdsjhckjsdhcjdsk"),
				wantedL:     []byte("\nContent-Disposition: form-data; name=\"alice\"\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\nContent-Disposition: form-data; name=\"alice\"\r\n\r\n\" is ending part"),
			},

			{
				name:        "Precending LF, 3 CRLF. LF + CDinsuf + CRLF + CT + 2*CRLF + rand",
				bs:          []byte("\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nsdjkchdjhcskdhcdsjhckjsdhcjdsk"),
				wantedL:     []byte("\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
				wantedError: errors.New("in repo.GetHeaderLines header \"\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n\" is ending part"),
			},
		*/
		{
			name:        "Precending LF, 4 CRLF. LF + Boundary + CRLF+ CDinsuf + CRLF + CT + 2*CRLF + rand",
			bs:          []byte("\nbPrefixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nsdjkchdjhcskdhcdsjhckjsdhcjdsk"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\nbPrefixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedError: errors.New("in repo.GetHeaderLines header \"\nbPrefixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n\" is ending part"),
		},
		/*
			{
				name:        "Succeeding LF, 0 CRLF. CD full + CR",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r"),
				wantedError: errors.New("in repo.GetHeaderLines header \"Content-Disposition: form-data; name=\"alice\"\r\" is not full"),
			},

			{
				name:        "Succeeding LF, 1 CRLF. CDsuf + CRLF + CR",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r"),
				wantedError: errors.New("in repo.GetHeaderLines header \"Content-Disposition: form-data; name=\"alice\"\r\n\r\" is not full"),
			},

			{
				name:        "Succeeding LF, 1 CRLF. CDinsuf + CRLF + CT + CR",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r"),
				wantedError: errors.New("in repo.GetHeaderLines header \"Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\" is not full"),
			},

			{
				name:        "Succeeding LF, 2 CRLF. CDinsuf + CRLF + CT + CRLF + CR",
				bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r"),
				wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r"),
				wantedError: errors.New("in repo.GetHeaderLines header \"Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\" is not full"),
			},

			{
				name:        "Succeeding LF, 3 CRLF. rand + CR",
				bs:          []byte("sdjkchdjhcskdhcdsjhckjsdhcjdsk\r\nsdjhfjdshjfsd\r\ngruihgeruhguerhguerg\r\n121312j412jk4g1jk4gjkg\r"),
				wantedL:     nil,
				wantedError: errors.New("in repo.GetHeaderLines no header found"),
			},
		*/
	}

	for _, v := range tt {
		s.Run(v.name, func() {
			got, err := GetHeaderLines(v.bs, v.bou)
			if v.wantedError != nil || err != nil {
				s.Equal(v.wantedError, err)
			}
			s.Equal(v.wantedL, got)
		})
	}
}

func (s *byteOpsSuite) TestKnownBoundaryPart() {
	tt := []struct {
		name   string
		bs     []byte
		bou    Boundary
		wanted []byte
	}{
		{
			name:   "",
			bs:     []byte("RootbSuffix\r\n"),
			bou:    Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wanted: []byte("Root"),
		},
	}
	for _, v := range tt {
		s.Equal(v.wanted, KnownBoundaryPart(v.bs, v.bou))
	}
}

func (s *byteOpsSuite) TestGetLineFeftWithCRLFLeft() {
	tt := []struct {
		name      string
		bs        []byte
		fromIndex int
		limit     int
		bou       Boundary
		wanted    []byte
	}{
		{
			name:      "good",
			bs:        []byte("azazazazazaza\r\nbzbzbzbbz"),
			fromIndex: 23,
			limit:     20,
			wanted:    []byte("\r\nbzbzbzbbz"),
		},

		{
			name:      "good LF or CR met",
			bs:        []byte("azazaza\r\nzaz\nazab\r\rzbzbzbbz"),
			fromIndex: 23,
			limit:     20,
			wanted:    []byte("\r\nzaz\nazab\r\rzbzbzbbz"),
		},

		{
			name:      "good LF in the end",
			bs:        []byte("azazazazazazabzbzbzbbz\r"),
			fromIndex: 22,
			limit:     20,
			wanted:    []byte("\r"),
		},

		{
			name:      "good last boundary",
			bs:        []byte("azazazazazazabzbzbzbbz\r\nbPrefixbRootbSuffix\r\n"),
			fromIndex: 44,
			limit:     30,
			bou:       Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wanted:    []byte("\r\nbPrefixbRootbSuffix\r\n"),
		},

		{
			name:      "wrong last boundary",
			bs:        []byte("azazazazazazabzbzbzbbz\r\nbPdefixbBoombSuggix\r\n"),
			fromIndex: 44,
			limit:     30,
			bou:       Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wanted:    []byte("\r\n"),
		},

		{
			name:      "limit exceeded",
			bs:        []byte("azazazazazaza\r\nbzbzbzbbz"),
			fromIndex: 23,
			limit:     120,
			wanted:    []byte("\r\nbzbzbzbbz"),
		},

		{
			name:      "fromindex exceeded",
			bs:        []byte("azazazazazaza\r\nbzbzbzbbz"),
			fromIndex: 123,
			limit:     20,
			wanted:    []byte("\r\nbzbzbzbbz"),
		},

		{
			name:      "no CR at all",
			bs:        []byte("zazazazazazaza\rbzbzbzb"),
			fromIndex: 23,
			limit:     20,
			wanted:    []byte("zazazazazazaza\rbzbzbzb"),
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() { s.Equal(v.wanted, GetLineWithCRLFLeft(v.bs, v.fromIndex, v.limit, v.bou)) })
	}
}

func (s *byteOpsSuite) TestBeginningEqual() {
	tt := []struct {
		name   string
		s1     []byte
		s2     []byte
		wanted bool
	}{
		{
			name:   "s1 > s2",
			s1:     []byte("abrahamhfjsdhfjksdhfjksd"),
			s2:     []byte("abraham"),
			wanted: true,
		},

		{
			name:   "s2 > s1",
			s2:     []byte("richardohfjsdhfjksdhfjksd"),
			s1:     []byte("richardo"),
			wanted: true,
		},

		{
			name:   "unhappy",
			s1:     []byte("abrahamhfjsdhfjksdhfjksd"),
			s2:     []byte("abrahamX"),
			wanted: false,
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			s.Equal(v.wanted, BeginningEqual(v.s1, v.s2))
		})
	}
}

func (s *byteOpsSuite) TestContainsBouEnding() {
	tt := []struct {
		name string
		bs   []byte
		bou  Boundary
		want bool
	}{
		{
			name: "happy",
			bs:   []byte("efixbRootbSuffix"),
			bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			want: true,
		},

		{
			name: "unhappy",
			bs:   []byte("defixbRootbSuffix"),
			bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			want: false,
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			s.Equal(v.want, ContainsBouEnding(v.bs, v.bou))
		})
	}
}

func (s *byteOpsSuite) TestIsBoudary() {
	tt := []struct {
		name string
		b    []byte
		n    []byte
		bou  Boundary
		want bool
	}{

		{
			name: "case 1",
			b:    []byte("\r\n"),
			n:    []byte("bPrefixbRoot\r\nContent-Disposition"),
			bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			want: true,
		},

		{
			name: "case 2",
			b:    []byte(""),
			n:    []byte("bPrefixbRoot"),
			bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			want: true,
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			s.Equal(v.want, IsBoundary(v.b, v.n, v.bou))
		})
	}
}

func (s *byteOpsSuite) TestIsLastBoundaryPart() {
	tt := []struct {
		name string
		b    []byte
		bou  Boundary
		want bool
	}{
		{
			name: "ordinary len(b)",
			b:    []byte("---63643643643--"),
			bou:  Boundary{Prefix: []byte("--"), Root: []byte("-------------63643643643")},
			want: true,
		},

		{
			name: "len(b) == 1",
			b:    []byte("-"),
			bou:  Boundary{Prefix: []byte("--"), Root: []byte("-------------63643643643")},
			want: true,
		},
		{
			name: "len(b) == 2",
			b:    []byte("--"),
			bou:  Boundary{Prefix: []byte("--"), Root: []byte("-------------63643643643")},
			want: true,
		},
		{
			name: "wrong 1",
			b:    []byte("---63643643642--"),
			bou:  Boundary{Prefix: []byte("--"), Root: []byte("-------------63643643643")},
			want: false,
		},
		{
			name: "wrong 2",
			b:    []byte("---73643643643--"),
			bou:  Boundary{Prefix: []byte("--"), Root: []byte("-------------63643643643")},
			want: false,
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			s.Equal(v.want, IsLastBoundaryPart(v.b, v.bou))
		})
	}
}
