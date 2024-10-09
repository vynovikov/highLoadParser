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

/*
	func (s *dataHandlerSuite) TestCreate() {
		tt := []struct {
			name              string
			initDataHandler   DataHandler
			dto               DataHandlerDTO
			bou               Boundary
			wantedDataHandler DataHandler
			wantedResult      ProducerUnit
			wantedError       error
		}{

			{
				name: "1. Empty Map, !isSub, full header, name only",
				initDataHandler: &memoryDataHandlerStruct{
					Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 0, Dh_body: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nazazaza"), Dh_b: 0, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     0,
					Dh_FormName: "alice",
					Dh_FileName: "",
					Dh_Body:     []byte("azazaza"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
				},
			},

			{
				name: "2. Empty Map, !isSub, full header, name + filename",
				initDataHandler: &memoryDataHandlerStruct{
					Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 0, Dh_body: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), Dh_b: 0, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     0,
					Dh_FormName: "alice",
					Dh_FileName: "short.txt",
					Dh_Body:     []byte("azazaza"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
				},
			},

			{
				name: "3. Empty Map, !isSub, not full header",
				initDataHandler: &memoryDataHandlerStruct{
					Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 0, Dh_body: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"), Dh_b: 0, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     0,
					Dh_FormName: "",
					Dh_FileName: "",
					Dh_Body:     make([]byte, 0),
					Dh_Start:    false,
					Dh_IsSub:    false,
					Dh_End:      false,
				},
			},

			{
				name: "4. Empty Map, isSub",
				initDataHandler: &memoryDataHandlerStruct{
					Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 0, Dh_body: []byte("\r\n----"), Dh_b: 0, Dh_e: 2, Dh_isSub: true, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 0}: {true: {
							E: 2,
							H: headerData{
								formName:    "",
								fileName:    "",
								headerBytes: []byte("\r\n----"),
							}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     0,
					Dh_FormName: "",
					Dh_FileName: "",
					Dh_Body:     make([]byte, 0),
					Dh_Start:    false,
					Dh_IsSub:    true,
					Dh_End:      false,
				},
			},

			{
				name: "5. Map has other TS, !isSub",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "www"}: {{ts: "www", part: 4}: {false: {
							E: 1,
							H: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 0, Dh_body: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), Dh_b: 0, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "www"}: {{ts: "www", part: 4}: {false: {
							E: 1,
							H: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {false: {
							E: 1,
							H: headerData{
								formName:    "bob",
								fileName:    "long.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     0,
					Dh_FormName: "bob",
					Dh_FileName: "long.txt",
					Dh_Body:     []byte("bzbzbzbzbz"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "6. Map has other TS, isSub",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "www"}: {{ts: "www", part: 4}: {false: {
							E: 1,
							H: headerData{
								formName:    "alice",
								fileName:    "short.txt",
								headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
							}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 0, Dh_body: []byte("\r\n-----"), Dh_b: 0, Dh_e: 2, Dh_isSub: true, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "www"}: {{ts: "www", part: 4}: {
							false: {
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
						{ts: "qqq"}: {{ts: "qqq", part: 0}: {
							true: {
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     0,
					Dh_FormName: "",
					Dh_FileName: "",
					Dh_Body:     make([]byte, 0),
					Dh_Start:    false,
					Dh_IsSub:    true,
					Dh_End:      false,
				},
			},

			{
				name: "7. Map has same key and different part, !isSub",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 3}: {
							false: {
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 4, Dh_body: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), Dh_b: 0, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 5}: {
							false: {
								E: 1,
								H: headerData{
									formName:    "bob",
									fileName:    "long.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     4,
					Dh_FormName: "bob",
					Dh_FileName: "long.txt",
					Dh_Body:     []byte("bzbzbzbzbz"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "8. Map has same key and different part, isSub",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 3}: {
							false: {
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 4, Dh_body: []byte("\r\n-----"), Dh_b: 0, Dh_e: 2, Dh_isSub: true, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 4}: {

							true: {
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     4,
					Dh_FormName: "",
					Dh_FileName: "",
					Dh_Body:     make([]byte, 0),
					Dh_Start:    false,
					Dh_IsSub:    true,
					Dh_End:      false,
				},
			},

			{
				name: "9. Map has same key and same part, !isSub",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 4}: {
							false: {
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 4, Dh_body: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), Dh_b: 0, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 5}: {
							false: {
								E: 1,
								H: headerData{
									formName:    "bob",
									fileName:    "long.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     4,
					Dh_FormName: "bob",
					Dh_FileName: "long.txt",
					Dh_Body:     []byte("bzbzbzbzbz"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "10. len(Map[keyGeneral]) == 0, dto.B() == 0, dto.E() == 0, last",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 4, Dh_body: []byte("Content-Disposition: form-data; name=\"bob\"\r\n\r\nbzbzbzbzbz"), Dh_b: 0, Dh_e: 0, Dh_isSub: false, Dh_last: true},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     4,
					Dh_FormName: "bob",
					Dh_FileName: "",
					Dh_Body:     []byte("bzbzbzbzbz"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      true,
					Dh_Final:    true,
				},
			},

			{
				name: "11. Map has same key and same part, isSub",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 4}: {
							false: {
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 4, Dh_body: []byte("\r\n-----"), Dh_b: 0, Dh_e: 2, Dh_isSub: true, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 4}: {
							false: {
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}},
							true: {
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}}}}},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     4,
					Dh_FormName: "",
					Dh_FileName: "",
					Dh_Body:     make([]byte, 0),
					Dh_Start:    false,
					Dh_IsSub:    true,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "12. Map has same key and same part, value.e = True, !d.IsSub, d.E() == 2",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 4}: {
							false: {
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 4, Dh_body: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), Dh_b: 0, Dh_e: 2, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 4}: {
							false: {
								E: 2,
								H: headerData{
									formName:    "bob",
									fileName:    "long.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     4,
					Dh_FormName: "bob",
					Dh_FileName: "long.txt",
					Dh_Body:     []byte("bzbzbzbzbz"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "13. Map has same key and same part, value.e = 2, !d.IsSub, d.E() == 2",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 4}: {
							true: {
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 4, Dh_body: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), Dh_b: 0, Dh_e: 2, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 5}: {
							true: {
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}},
							false: {
								E: 2,
								H: headerData{
									formName:    "bob",
									fileName:    "long.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								},
							}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     4,
					Dh_FormName: "bob",
					Dh_FileName: "long.txt",
					Dh_Body:     []byte("bzbzbzbzbz"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "14. Map has same key and same part, value.e == 2, d.IsSub, d.E() == 2",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 4}: {
							false: {
								E: 2,
								H: headerData{
									formName:    "bob",
									fileName:    "long.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								},
							}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 4, Dh_body: []byte("\r\n-----"), Dh_b: 0, Dh_e: 2, Dh_isSub: true, Dh_last: false},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 5}: {
							false: {
								E: 2,
								H: headerData{
									formName:    "bob",
									fileName:    "long.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								},
							},
							true: {
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}}}}},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     4,
					Dh_FormName: "",
					Dh_FileName: "",
					Dh_Body:     make([]byte, 0),
					Dh_Start:    false,
					Dh_IsSub:    true,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},
		}

		for _, v := range tt {

			s.Run(v.name, func() {

				gotResult, err := v.initDataHandler.Create(v.dto, v.bou)

				if v.wantedError != nil {

					s.Equal(v.wantedError, err)
				}

				s.Equal(v.wantedDataHandler, v.initDataHandler)

				s.Equal(v.wantedResult, gotResult)
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
			wantedResult      ProducerUnit
			wantedError       error
		}{

			{
				name: "1. Part matched, header is full",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("azazaza"), Dh_b: 1, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 2}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "alice",
					Dh_FileName: "short.txt",
					Dh_Body:     []byte("azazaza"),
					Dh_Start:    false,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "1.1. Part matched, header is full, last == true",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("azazaza"), Dh_b: 1, Dh_e: 1, Dh_isSub: false, Dh_last: true},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 2}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "alice",
					Dh_FileName: "short.txt",
					Dh_Body:     []byte("azazaza"),
					Dh_Start:    false,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    true,
				},
			},

			{
				name: "2. Part matched, header is not full",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("Content-Dispos"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), Dh_b: 1, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 2}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "alice",
					Dh_FileName: "short.txt",
					Dh_Body:     []byte("azazaza"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "2.2. Part matched, header is not full, last",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("Content-Dispos"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), Dh_b: 1, Dh_e: 1, Dh_isSub: false, Dh_last: true},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 2}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "alice",
					Dh_FileName: "short.txt",
					Dh_Body:     []byte("azazaza"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    true,
				},
			},

			{
				name: "3. Part is not matched",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("Content-Dispos"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 2, Dh_body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), Dh_b: 1, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("Content-Dispos"),
								}}}},
					},
					Buffer: []DataHandlerDTO{
						&DataHandlerUnit{Dh_ts: "qqq", Dh_part: 2, Dh_body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), Dh_b: 1, Dh_e: 1, Dh_isSub: false, Dh_last: false},
					},
				},
			},

			{
				name: "4. dto.E() == False",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("Content-Dispos"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), Dh_b: 1, Dh_e: 0, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "alice",
					Dh_FileName: "short.txt",
					Dh_Body:     []byte("azazaza"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      true,
					Dh_Final:    false,
				},
			},

			{
				name: "5. dto.E() == 2, header is not full, value.e = True",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("Content-Dispos"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), Dh_b: 1, Dh_e: 2, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 2,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "alice",
					Dh_FileName: "short.txt",
					Dh_Body:     []byte("azazaza"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "6. dto.E() == 2, header is full, value.e = True",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("azazaza"), Dh_b: 1, Dh_e: 2, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 2,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "alice",
					Dh_FileName: "short.txt",
					Dh_Body:     []byte("azazaza"),
					Dh_Start:    false,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "7. dto.E() == 2, header is not full, len(value) > 1",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("Content-Dispos"),
								}},
							true: value{
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}},
						}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("ition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), Dh_b: 1, Dh_e: 2, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 2}: {
							false: value{
								E: 2,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}},
							true: value{
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}},
						}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "alice",
					Dh_FileName: "short.txt",
					Dh_Body:     []byte("azazaza"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "8. dto.E() == 2, header is full, len(value) > 1",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}},
							true: value{
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}},
						}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("azazaza"), Dh_b: 1, Dh_e: 2, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 2}: {
							false: value{
								E: 2,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}},
							true: value{
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}},
						}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "alice",
					Dh_FileName: "short.txt",
					Dh_Body:     []byte("azazaza"),
					Dh_Start:    false,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "9. DataUnit next to fork, new value, d.E() == True, new boundary",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}},
							true: value{
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}},
						}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("----------bRoot\r\nContent-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), Dh_b: 1, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 2}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "bob",
									fileName:    "long.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}},
						}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "bob",
					Dh_FileName: "long.txt",
					Dh_Body:     []byte("bzbzbzbzbz"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "10. DataUnit next to fork, d.E() == 2, new boundary => new map value",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}},
							true: value{
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}},
						}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("----------bRoot\r\nContent-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzbz"), Dh_b: 1, Dh_e: 2, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 2}: {
							false: value{
								E: 2,
								H: headerData{
									formName:    "bob",
									fileName:    "long.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}},
						}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "bob",
					Dh_FileName: "long.txt",
					Dh_Body:     []byte("bzbzbzbzbz"),
					Dh_Start:    true,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},

			{
				name: "11. DataUnit next to fork,  d.E() == True, boundary-like bytes => old map value remains",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 2,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}},
							true: value{
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}},
						}},
					},
					Buffer: []DataHandlerDTO{},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 1, Dh_body: []byte("------\r\nbzbzbzbzbz"), Dh_b: 1, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("---------------"), Root: []byte("bRoot")},
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 2}: {
							false: value{
								E: 1,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}},
						}},
					},
					Buffer: []DataHandlerDTO{},
				},
				wantedResult: &ProducerUnitStruct{
					Dh_TS:       "qqq",
					Dh_Part:     1,
					Dh_FormName: "alice",
					Dh_FileName: "short.txt",
					Dh_Body:     []byte("\r\n-----------\r\nbzbzbzbzbz"),
					Dh_Start:    false,
					Dh_IsSub:    false,
					Dh_End:      false,
					Dh_Final:    false,
				},
			},
		}

		for _, v := range tt {

			s.Run(v.name, func() {

				gotResult, err := v.initDataHandler.Updade(v.dto, v.bou)

				if v.wantedError != nil {

					s.Equal(v.wantedError, err)
				}

				s.Equal(v.wantedDataHandler, v.initDataHandler)

				s.Equal(v.wantedResult, gotResult)
			})
		}
	}

func (s *dataHandlerSuite) TestDelete() {

		tt := []struct {
			name              string
			initDataHandler   DataHandler
			ts                string
			wantedDataHandler DataHandler
			wantedError       error
		}{
			{
				name: "1. Delete key ",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 2,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
						{ts: "www"}: {{ts: "www", part: 1}: {
							false: value{
								E: 2,
								H: headerData{
									formName:    "bob",
									fileName:    "long.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
				ts: "qqq",
				wantedDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "www"}: {{ts: "www", part: 1}: {
							false: value{
								E: 2,
								H: headerData{
									formName:    "bob",
									fileName:    "long.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}}}},
					},
					Buffer: []DataHandlerDTO{},
				},
			},

			{
				name: "2. Delete key and remake Map",
				initDataHandler: &memoryDataHandlerStruct{
					Map: map[keyGeneral]map[keyDetailed]map[bool]value{
						{ts: "qqq"}: {{ts: "qqq", part: 1}: {
							false: value{
								E: 2,
								H: headerData{
									formName:    "alice",
									fileName:    "short.txt",
									headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
								}},
							true: value{
								E: 2,
								H: headerData{
									formName:    "",
									fileName:    "",
									headerBytes: []byte("\r\n-----"),
								}},
						}},
					},
					Buffer: []DataHandlerDTO{},
				},
				ts: "qqq",
				wantedDataHandler: &memoryDataHandlerStruct{
					Map:    map[keyGeneral]map[keyDetailed]map[bool]value{},
					Buffer: []DataHandlerDTO{},
				},
			},
		}

		for _, v := range tt {

			s.Run(v.name, func() {

				err := v.initDataHandler.Delete(v.ts)

				if v.wantedError != nil {

					s.Equal(v.wantedError, err)
				}

				s.Equal(v.wantedDataHandler, v.initDataHandler)
			})
		}
	}

func (s *dataHandlerSuite) TestUpdateValue() {

		tt := []struct {
			name              string
			initValue         value
			dto               DataHandlerDTO
			bou               Boundary
			headerEndingIndex int
			wantedValue       value
			wantedIndex       int
			wantedError       error
		}{

			{
				name: "1.Dh_b == 0. Value is full",
				initValue: value{
					E: 1,
					H: headerData{
						formName:    "",
						fileName:    "",
						headerBytes: []byte(""),
					},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 0, Dh_body: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nazazaza"), Dh_b: 0, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				wantedValue: value{
					E: 1,
					H: headerData{
						formName:    "alice",
						fileName:    "",
						headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
					}},
				wantedIndex: 48,
			},

			{
				name: "2. Dh_b == 1. Value is not full, dto has header with no filename",
				initValue: value{
					E: 1,
					H: headerData{
						formName:    "",
						fileName:    "",
						headerBytes: []byte("Content-Dispo"),
					},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 0, Dh_body: []byte("sition: form-data; name=\"alice\"\r\n\r\nazazaza"), Dh_b: 1, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				wantedValue: value{
					E: 1,
					H: headerData{
						formName:    "alice",
						fileName:    "",
						headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
					}},
				wantedIndex: 35,
			},

			{
				name: "3. Dh_b == True. Value is not full, dto has header has both name and filename",
				initValue: value{
					E: 1,
					H: headerData{
						formName:    "",
						fileName:    "",
						headerBytes: []byte("Content-Dispo"),
					},
				},
				dto: &DataHandlerUnit{Dh_ts: "qqq", Dh_part: 0, Dh_body: []byte("sition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), Dh_b: 1, Dh_e: 1, Dh_isSub: false, Dh_last: false},
				bou: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
				wantedValue: value{
					E: 1,
					H: headerData{
						formName:    "alice",
						fileName:    "short.txt",
						headerBytes: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
					}},
				wantedIndex: 83,
			},
		}

		for _, v := range tt {

			s.Run(v.name, func() {

				gotValue, gotIndex, err := updateValue(v.initValue, v.dto, v.bou)

				if v.wantedError != nil {

					s.Equal(v.wantedError, err)
				}

				s.Equal(v.wantedValue, gotValue)
				s.Equal(v.wantedIndex, gotIndex)

			})
		}
	}
*/
func (s *dataHandlerSuite) TestGetHeaderEndingIndex() {

	tt := []struct {
		name        string
		body        []byte
		header      []byte
		wantedIndex int
		wantedError error
	}{

		{
			name:        "1. No header ending in body",
			body:        []byte("azaza"),
			header:      []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedIndex: -1,
		},

		{
			name:        "2. Body beginning is header ending ",
			body:        []byte("sition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazaza"),
			header:      []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedIndex: 83,
		},

		{
			name:        "3. CR in the beginning of body",
			body:        []byte("\nazaza"),
			header:      []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
			wantedIndex: -1,
		},
	}
	for _, v := range tt {

		s.Run(v.name, func() {

			s.Equal(v.wantedIndex, getHeaderEndingIndex(v.body, v.header))
		})
	}
}
