package dataHandler

import (
	"errors"
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
		/*
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
		*/
		{
			name: "1.1. Empty Map, !isSub, full header",
			initDataHandler: &memoryDataHandlerStruct{
				Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
				Buffer: []DataHandlerDTO{},
			},
			dto: &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), b: False, e: True, isSub: false, last: false},
			wantedDataHandler: &memoryDataHandlerStruct{
				Map: map[keyGeneral]map[keyDetailed]map[bool]value{
					{ts: "qqq"}: {{ts: "qqq", part: 1}: {false: value{
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
		/*
			{
				name: "1.2. Empty Map, !isSub, not full header",
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
		*/
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

func (s *dataHandlerSuite) TestNewValue() {

	tt := []struct {
		name        string
		dto         DataHandlerDTO
		bou         Boundary
		wantedValue value
		wantedError error
	}{
		{
			name: "1. Full header",
			dto:  &DataHandlerUnit{ts: "qqq", part: 0, body: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), b: False, e: True, isSub: false, last: false},
			wantedValue: value{
				e: True,
				h: headerData{
					formName:    "alice",
					fileName:    "short.txt",
					headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
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
			name:        "1 CRLF last Boundary",
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
			name:        "4 CRLF 1 Boundary ending 2 header lines, 1 random line",
			bs:          []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedError: errors.New("in repo.GetHeaderLines header \"fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n\" is ending part"),
		},

		{
			name:        "5 CRLF 1 Boundary ending 2 header lines, 1 random line",
			bs:          []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\naaaaaaaaaaaaaaaaaaaaaaaaa\r\nbbbbbbbbbbbbbbbb\r\ndsfguigdfa"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedError: errors.New("in repo.GetHeaderLines header \"fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n\" is ending part"),
		},

		{
			name:        "4 CRLF 1 Boundary ending 1 header line, 2 random lines",
			bs:          []byte("fixbRoot\r\nContent-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\ndsfguigdfa\r\nf2r7fr27fr2f7r2"),
			bou:         Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			wantedL:     nil,
			wantedError: errors.New("in repo.GetHeaderLines no header found"),
		},

		{
			name:        "4 CRLF 1 Boundary ending rest random lines",
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
