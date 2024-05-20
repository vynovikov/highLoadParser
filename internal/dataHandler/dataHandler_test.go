package dataHandler

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type dataHandlerSuite struct {
	suite.Suite
}

func TestDataHandlerSuite(t *testing.T) {

	suite.Run(t, new(dataHandlerSuite))
}

func (s *dataHandlerSuite) TestGenBoundary() {
	boundaryVoc := Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}

	boundaryCalced := genBoundary(boundaryVoc)

	s.Equal([]byte("\r\nbPrefix"+"bRoot"), boundaryCalced)
}

func (s *dataHandlerSuite) TestIsLastBoundaryPart() {
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
			s.Equal(v.want, isLastBoundaryPart(v.b, v.bou))
		})
	}
}

func (s *dataHandlerSuite) TestCreate() {
	tt := []struct {
		name              string
		initDataHandler   DataHandler
		dto               DataHandlerDTO
		bou               Boundary
		wantedDataHandler DataHandler
		wantedError       error
	}{

		{
			name: "1. Empty Map, !isSub, full header, name only",
			initDataHandler: &memoryDataHandlerStruct{
				Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nazazaza"), b: False, e: True, isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "2. Empty Map, !isSub, full header, name + filename",
			initDataHandler: &memoryDataHandlerStruct{
				Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), b: False, e: True, isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "3. Empty Map, !isSub, not full header",
			initDataHandler: &memoryDataHandlerStruct{
				Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"), b: False, e: True, isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "4. Empty Map, isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("\r\n----"), b: False, e: Probably, isSub: true, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 0}: {true: {
						e: Probably,
						h: headerData{
							formName:    "",
							fileName:    "",
							headerBytes: []byte("\r\n----"),
						}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "5. Map has other TS, !isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {false: {
						e: True,
						h: headerData{
							formName:    "alice",
							fileName:    "short.txt",
							headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
						}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), b: False, e: True, isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {false: {
						e: True,
						h: headerData{
							formName:    "alice",
							fileName:    "short.txt",
							headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
						}}}},
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {false: {
						e: True,
						h: headerData{
							formName:    "bob",
							fileName:    "long.txt",
							headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
						}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "6. Map has other TS, isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {false: {
						e: True,
						h: headerData{
							formName:    "alice",
							fileName:    "short.txt",
							headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
						}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("\r\n-----"), b: False, e: Probably, isSub: true, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "www"}: {{ts: "www", part: 4}: {
						false: {
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
					{ts: "qqq"}: {{ts: "qqq", part: 0}: {
						true: {
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "7. Map has same key and different part, !isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 3}: {
						false: {
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), b: False, e: True, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {
						false: {
							e: True,
							h: headerData{
								formName:    "bob",
								fileName:    "long.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "8. Map has same key and different part, isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 3}: {
						false: {
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("\r\n-----"), b: False, e: Probably, isSub: true, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {
						true: {
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "9. Map has same key and same part, !isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {
						false: {
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), b: False, e: True, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {
						false: {
							e: True,
							h: headerData{
								formName:    "bob",
								fileName:    "long.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "10. Map has same key and same part, isSub",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {
						false: {
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("\r\n-----"), b: False, e: Probably, isSub: true, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {
						false: {
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}},
						true: {
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}}}}},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "11. Map has same key and same part, value.e = True, !d.IsSub, d.E() == Probably",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {
						false: {
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), b: False, e: Probably, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {
						false: {
							e: Probably,
							h: headerData{
								formName:    "bob",
								fileName:    "long.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "12. Map has same key and same part, value.e = Probably, !d.IsSub, d.E() == Probably",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {
						true: {
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), b: False, e: Probably, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {
						true: {
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}},
						false: {
							e: Probably,
							h: headerData{
								formName:    "bob",
								fileName:    "long.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							},
						}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "13. Map has same key and same part, value.e == Probably, d.IsSub, d.E() == Probably",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 4}: {
						false: {
							e: Probably,
							h: headerData{
								formName:    "bob",
								fileName:    "long.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							},
						}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 4, body: []byte("\r\n-----"), b: False, e: Probably, isSub: true, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 5}: {
						false: {
							e: Probably,
							h: headerData{
								formName:    "bob",
								fileName:    "long.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							},
						},
						true: {
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}}}}},
				Buffer: []DataHandlerDTO{},
			},
		},
	}

	for _, v := range tt {

		s.Run(v.name, func() {

			err := v.initDataHandler.Create(v.dto, v.bou)

			if v.wantedError != nil {

				s.Equal(v.wantedError, err)
			}

			s.Equal(v.wantedDataHandler, v.initDataHandler)
		})
	}
}

func (s *dataHandlerSuite) TestUpdate() {
	tt := []struct {
		name              string
		initDataHandler   DataHandler
		dto               DataHandlerDTO
		bou               Boundary
		wantedDataHandler DataHandler
		wantedError       error
	}{

		{
			name: "1. Part matched, header is full",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 1, body: []byte("azazaza"), b: True, e: True, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 2}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "2. Part matched, header is not full",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("Content-Dispos"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 1, body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), b: True, e: True, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 2}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "3. Part is not matched",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("Content-Dispos"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 2, body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), b: True, e: True, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("Content-Dispos"),
							}}}},
				},
				Buffer: []DataHandlerDTO{
					&DataHandlerUnit{ts: "qqq", part: 2, body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), b: True, e: True, isSub: false, last: false},
				},
			},
		},

		{
			name: "4. dto.E() == False",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("Content-Dispos"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 1, body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), b: True, e: False, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "5. dto.E() == Probably, header is not full, value.e = True",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("Content-Dispos"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 1, body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), b: True, e: Probably, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: Probably,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "5. dto.E() == Probably, header is full, value.e = True",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 1, body: []byte("azazaza"), b: True, e: Probably, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: Probably,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "7. dto.E() == Probably, header is not full, len(value) > 1",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("Content-Dispos"),
							}},
						true: value{
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}},
					}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 1, body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), b: True, e: Probably, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 2}: {
						false: value{
							e: Probably,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}},
						true: value{
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}},
					}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "8. dto.E() == Probably, header is full, len(value) > 1",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}},
						true: value{
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}},
					}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 1, body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), b: True, e: Probably, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 2}: {
						false: value{
							e: Probably,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}},
						true: value{
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}},
					}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "9. DataUnit next to fork, new value, d.E() == True, new boundary",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}},
						true: value{
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}},
					}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 1, body: []byte("----------bRoot\r\nContent-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), b: True, e: True, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 2}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "bob",
								fileName:    "long.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}},
					}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "10. DataUnit next to fork, d.E() == Probably, new boundary => new map value",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}},
						true: value{
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}},
					}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 1, body: []byte("----------bRoot\r\nContent-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), b: True, e: Probably, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 2}: {
						false: value{
							e: Probably,
							h: headerData{
								formName:    "bob",
								fileName:    "long.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}},
					}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		{
			name: "11. DataUnit next to fork,  d.E() == True, boundary-like bytes => old map value remains",
			initDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}},
						true: value{
							e: Probably,
							h: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n-----"),
							}},
					}},
				},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 1, body: []byte("------\r\nbzbzbzbzbz"), b: True, e: True, isSub: false, last: false},
			bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 2}: {
						false: value{
							e: True,
							h: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}},
					}},
				},
				Buffer: []DataHandlerDTO{},
			},
		},

		// TODO:
		// DataUnit after isSub 2 cases

	}

	for _, v := range tt {

		s.Run(v.name, func() {

			err := v.initDataHandler.Updade(v.dto, v.bou)

			if v.wantedError != nil {

				s.Equal(v.wantedError, err)
			}

			s.Equal(v.wantedDataHandler, v.initDataHandler)
		})
	}
}

func (s *dataHandlerSuite) TestNewValue() {

	tt := []struct {
		name        string
		dto         DataHandlerDTO
		bou         Boundary
		wantedValue value
		wantedError error
	}{

		{
			name: "1. Full header, name only",
			dto:  &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nazazaza"), b: False, e: True, isSub: false, last: false},
			bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedValue: value{
				e: True,
				h: headerData{
					formName:    "alice",
					fileName:    "",
					headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
				}},
		},

		{
			name: "2. Full header, name + filename",
			dto:  &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), b: False, e: True, isSub: false, last: false},
			bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedValue: value{
				e: True,
				h: headerData{
					formName:    "alice",
					fileName:    "short.txt",
					headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
				}},
		},

		{
			name: "3. Partial header",
			dto:  &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r"), b: False, e: True, isSub: false, last: false},
			bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedValue: value{
				e: True,
				h: headerData{
					formName:    "",
					fileName:    "",
					headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r"),
				}},
		},
	}

	for _, v := range tt {

		s.Run(v.name, func() {

			gotValue, err := newValue(v.dto, v.bou)

			if v.wantedError != nil {

				s.Equal(v.wantedError, err)
			}

			s.Equal(v.wantedValue, gotValue)

		})
	}
}

func (s *dataHandlerSuite) TestGetHeaderLines() {
	tt := []struct {
		name        string
		bs          []byte
		bou         Boundary
		wantedL     []byte
		wantedError error
	}{

		{
			name:        "0 CRLF <1 line right",
			bs:          []byte("Content-Disposition: form-data; name=\"al"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"al"),
			wantedError: fmt.Errorf("\"%s\" %w", string([]byte("Content-Disposition: form-data; name=\"al")), errHeaderNotFull),
		},

		{
			name:        "0 CRLF no header",
			bs:          []byte("azazazazazaza"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("header is not found"),
		},

		{
			name:        "1 CRLF whole sufficient 1-st line",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n"),
			wantedError: fmt.Errorf("\"%s\" %w", string([]byte("Content-Disposition: form-data; name=\"alice\"\r\n")), errHeaderNotFull),
		},

		{
			name:        "1 CRLF whole 1-st line, partyal 2-d",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short\"\r\nCon"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short\"\r\nCon"),
			wantedError: fmt.Errorf("\"%s\" %w", string([]byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short\"\r\nCon")), errHeaderNotFull),
		},
		{
			name:        "1 CRLF just CRLF, random line",
			bs:          []byte("\r\nr23hjrb23hrbj23hbrh23"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\r\n")), errHeaderEnding),
		},

		{
			name:        "1 CRLF last Boundary",
			bs:          []byte("azzsdfgsdhfdsfhsjdfhs\r\n"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("azzsdfgsdhfdsfhsjdfhs\r\n"),
			wantedError: errors.New("\"azzsdfgsdhfdsfhsjdfhs\r\n\" is the last"),
		},

		{
			name:        "1 CRLF no header_2",
			bs:          []byte("azzsdfgsdhfdsfhsjdfhs\r\nfskjfghsjfhgfjkhgjdfhgfd"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("header is not found"),
		},

		{
			name:        "2 CRLF 1 line CD sufficient",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
			wantedError: nil,
		},

		{
			name:        "2 CRLF 1 line CD sufficient + random line",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nhdsghdsvhsdvgshdgvsdv"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
			wantedError: nil,
		},

		{
			name:        "2 CRLF 1 line CD insufficient + CRLF + CT + CTLF",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"),
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"),
			wantedError: fmt.Errorf("\"%s\" %w", string([]byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n")), errHeaderNotFull),
		},

		{
			name:        "2 CRLF 1 line CD sufficient + random line",
			bs:          []byte("position: form-data; name=\"alice\"\r\n\r\nhdsghdsvhsdvgshdgvsdv"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("position: form-data; name=\"alice\"\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("position: form-data; name=\"alice\"\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "2 CRLF 1 header line right 2 random lines",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nsajkfdga\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("header is not found"),
		},

		{
			name:        "2 CRLF all lines random",
			bs:          []byte("we6fwfef6gewfgewfg7efge\r\nsajkfdga\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("header is not found"),
		},

		{
			name:        "2 CRLF CRLF, random line, CRLF, random line",
			bs:          []byte("\r\n2f3hg4f32ghf423gf324\r\nr23hjrb23hrbj23hbrh23"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\r\n")), errHeaderEnding),
		},

		{
			name:        "2 CRLF just 2 * CRLF, random line",
			bs:          []byte("\r\n\r\nr23hjrb23hrbj23hbrh23"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "3 CRLF 2 header lines (CD insufficient), 1 random line",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedError: nil,
		},

		{
			name:        "3 CRLF 1 header line (CD sufficient), 2 random lines",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\ndsfguigdfa6fhgf55\r\nggf8723g723gf823"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "3 CRLF: CRLF + CDsuf + 2*CRLF + rand",
			bs:          []byte("\r\nContent-Disposition: form-data; name=\"alice\"\r\n\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\nContent-Disposition: form-data; name=\"alice\"\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\r\nContent-Disposition: form-data; name=\"alice\"\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "3 CRLF: CRLF + CT + 2*CRLF + rand",
			bs:          []byte("\r\nContent-Type: text/plain\r\n\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\nContent-Type: text/plain\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\r\nContent-Type: text/plain\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "3 CRLF 1 header line, 2 random lines",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nCshdgfhsdgfhsdjf\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("header is not found"),
		},

		{
			name:        "3 CRLF all random lines",
			bs:          []byte("azazzazazzazazaz\r\nCzbbzbzbbzbzbbzbzbzbzbz\r\ndsfguigdfa\r\nf2r7fr27fr2f7r2"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("header is not found"),
		},

		{
			name:        "3 CRLF: CRLF, next random lines",
			bs:          []byte("\r\nr23hjrb23hrbj23hbrh23\r\nsgdhgsdwef6fr6632\r\n438ry34grg438rg438gr43"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\r\n")), errHeaderEnding),
		},

		{
			name:        "3 CRLF just  2 * CRLF, next random lines",
			bs:          []byte("\r\n\r\nr23hjrb23hrbj23hbrh23\r\nsgdhgsdwef6fr6632"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "4 CRLF 1 line CD sufficient + 3 random lines",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nhdsghdsvhsdvgshdgvsdv\r\nhjgvfjhdgvjhfdkgftv87dfvdfv\r\nsoiehfwoefhwefdgvjhsdv"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "4 CRLF 1 line CD insufficient + 2 random line",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nhdsghdsvhsdvgshdgvsdv\r\nhjgvfjhdgvjhfdkgftv87dfvdfv"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedError: nil,
		},

		{
			name:        "4 CRLF: CRLF + CDinsuf + CRLF + CT + 2*CRLF + rand",
			bs:          []byte("\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "4 CRLF 1 Boundary ending 2 header lines, 1 random line",
			bs:          []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "5 CRLF 1 Boundary ending 2 header lines, 1 random line",
			bs:          []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\naaaaaaaaaaaaaaaaaaaaaaaaa\r\nbbbbbbbbbbbbbbbb\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "4 CRLF 1 Boundary ending 1 header line, 2 random lines",
			bs:          []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\ndsfguigdfa\r\nf2r7fr27fr2f7r2"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("header is not found"),
		},

		{
			name:        "4 CRLF 1 Boundary ending rest random lines",
			bs:          []byte("fixbRoot\r\nCzbbzbzbbzbzbbzbzbzbzbz\r\ndsfguigdfa\r\nf2r7fr27fr2f7r2"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("header is not found"),
		},

		{
			name:        "4 CRLF just  2 * CRLF, next random lines",
			bs:          []byte("\r\n\r\nr23hjrb23hrbj23hbrh23\r\nsgdhgsdwef6fr6632\r\n3fd72fd73fd3727df23"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "5 CRLF 1 CRLF 2 header lines, 1 random line",
			bs:          []byte("\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\naaaaaaaaaaaaaaaaaaaaaaaaa\r\nbbbbbbbbbbbbbbbb\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "Precending LF, 0 CRLF. LF + rand",
			bs:          []byte("\nsdjkchdjhcskdhcdsjhckjsdhcjdsk"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\n")), errHeaderEnding),
		},

		{
			name:        "Precending LF, 3 CRLF. LF + rand",
			bs:          []byte("\nsdjkchdjhcskdhcdsjhckjsdhcjdsk\r\nsdjhfjdshjfsd\r\ngruihgeruhguerhguerg\r\n121312j412jk4g1jk4gjkg"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\n")), errHeaderEnding),
		},

		{
			name:        "Precending LF, 1 CRLF. CRLF + LF + rand",
			bs:          []byte("\n\r\nsdjkchdjhcskdhcdsjhckjsdhcjdsk"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\n\r\n")), errHeaderEnding),
		},

		{
			name:        "Precending LF, 2 CRLF. LF + CT + 2*CRLF + rand",
			bs:          []byte("\nContent-Type: text/plain\r\n\r\nsdjkchdjhcskdhcdsjhckjsdhcjdsk"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\nContent-Type: text/plain\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\nContent-Type: text/plain\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "Precending LF, 2 CRLF. LF + CDSuff + 2*CRLF + rand",
			bs:          []byte("\nContent-Disposition: form-data; name=\"alice\"\r\n\r\nsdjkch2323232djhcskdhcdsjhckjsdhcjdsk"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\nContent-Disposition: form-data; name=\"alice\"\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\nContent-Disposition: form-data; name=\"alice\"\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "Precending LF, 3 CRLF. LF + CDinsuf + CRLF + CT + 2*CRLF + rand",
			bs:          []byte("\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nsdjkchdjhcskdhcdsjhckjsdhcjdsk"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedError: fmt.Errorf("\"%s\" is %w", string([]byte("\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n")), errHeaderEnding),
		},

		{
			name:        "Succeeding LF, 0 CRLF. CD full + CR",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r"),
			wantedError: fmt.Errorf("\"%s\" %w", string([]byte("Content-Disposition: form-data; name=\"alice\"\r")), errHeaderNotFull),
		},

		{
			name:        "Succeeding LF, 1 CRLF. CDsuf + CRLF + CR",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r"),
			wantedError: fmt.Errorf("\"%s\" %w", string([]byte("Content-Disposition: form-data; name=\"alice\"\r\n\r")), errHeaderNotFull),
		},

		{
			name:        "Succeeding LF, 1 CRLF. CDinsuf + CRLF + CT + CR",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r"),
			wantedError: fmt.Errorf("\"%s\" %w", string([]byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r")), errHeaderNotFull),
		},

		{
			name:        "Succeeding LF, 2 CRLF. CDinsuf + CRLF + CT + CRLF + CR",
			bs:          []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r"),
			wantedError: fmt.Errorf("\"%s\" %w", string([]byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r")), errHeaderNotFull),
		},

		{
			name:        "Succeeding LF, 3 CRLF. rand + CR",
			bs:          []byte("sdjkchdjhcskdhcdsjhckjsdhcjdsk\r\nsdjhfjdshjfsd\r\ngruihgeruhguerhguerg\r\n121312j412jk4g1jk4gjkg\r"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("header is not found"),
		},

		{
			name:        "1 CRLF,Boundary prefix",
			bs:          []byte("\r\n----------"),
			bou:         Boundary{Prefix: []byte("-------------------"), Root: []byte("bRoot")},
			wantedL:     []byte("\r\n----------"),
			wantedError: fmt.Errorf("\"%s\" %w", string([]byte("\r\n----------")), errHeaderNotFull),
		},
	}

	for _, v := range tt {

		s.Run(v.name, func() {

			got, err := getHeaderLines(v.bs, v.bou)
			if v.wantedError != nil || err != nil {

				s.Equal(v.wantedError, err)
			}

			s.Equal(v.wantedL, got)
		})
	}
}
