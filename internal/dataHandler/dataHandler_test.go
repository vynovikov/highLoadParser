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

func (s *dataHandlerSuite) TestCheck() {

	tt := []struct {
		name         string
		dataHandler  DataHandler
		dto          DataHandlerDTO
		wantPresence Presence
		wantError    error
	}{
		{
			name: "No ASKG",
			dataHandler: &memoryDataHandlerStruct{
				Map: map[key]map[bool]value{},
			},
			//d:            &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 0, B: repo.False, E: repo.True}, APB: repo.AppPieceBody{B: []byte("azaza")}},
			dto:          &Dat{part: 0, ts: "qqq", body: []byte("azaza")},
			wantPresence: Presence{},
			wantError:    nil,
		},
		/*
			{
				name: "ASKG met, no ASKD",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {},
					},
				},
				d: &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 1, B: repo.False, E: repo.True}, APB: repo.AppPieceBody{B: []byte("azaza")}},
				wantPresence: Presence{
					ASKG: true},
				wantError: nil,
			},

			{
				name: "ASKG met, wrong ASKD",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {},
						},
					},
				},
				d: &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 1, B: repo.False, E: repo.True}, APB: repo.AppPieceBody{B: []byte("azaza")}},
				wantPresence: Presence{
					ASKG: true},
				wantError: nil,
			},

			{
				name: "AppSub, no ASKG",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{},
				},
				d:            &repo.AppSub{ASH: repo.AppSubHeader{TS: "qqq", Part: 1}, ASB: repo.AppSubBody{B: []byte("\r\n")}},
				wantPresence: Presence{},
				wantError:    nil,
			},

			{
				name: "AppSub, ASKG met, no opposite detailed branch",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "alice",
										FileName: "short.txt",
									},
									B: repo.BeginningData{Part: 0},
									E: repo.True,
								},
							},
						},
					},
				},
				d: &repo.AppSub{ASH: repo.AppSubHeader{TS: "qqq", Part: 1}, ASB: repo.AppSubBody{B: []byte("\r\n")}},
				wantPresence: Presence{
					ASKG: true,
				},
				wantError: nil,
			},
			{
				name: "AppSub, ASKG met, opposite detailed branch met",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "alice",
										FileName: "short.txt",
									},
									B: repo.BeginningData{Part: 0},
									E: repo.Probably,
								},
							},
						},
					},
				},
				d: &repo.AppSub{ASH: repo.AppSubHeader{TS: "qqq", Part: 1}, ASB: repo.AppSubBody{B: []byte("\r\n")}},
				wantPresence: Presence{
					ASKG: true,
					ASKD: true,
					OB:   true,
					GMap: map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
						{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
							false: {
								D: repo.Disposition{
									H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.BeginningData{Part: 0},
								E: repo.Probably,
							},
						},
					},
				},
				wantError: nil,
			},

			{
				name: "B() == repo.False && E() == repo.Probably, no askg met",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{},
				},
				d:            &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 0, B: repo.False, E: repo.Probably}, APB: repo.AppPieceBody{B: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazaza")}},
				wantPresence: Presence{},
			},

			{
				name: "B() == repo.False && E() == repo.Probably, askg met but no OB met",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {false: {
								D: repo.Disposition{
									H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.BeginningData{Part: 0},
								E: repo.True,
							}},
						},
					},
				},
				d: &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 2, B: repo.False, E: repo.Probably}, APB: repo.AppPieceBody{B: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazaza")}},
				wantPresence: Presence{
					ASKG: true,
					ASKD: true,
				},
				wantError: nil,
			},

			{
				name: "B() == repo.False && E() == repo.Probably, askg met, OB met",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: true}: {
								true: {
									D: repo.Disposition{
										H: []byte("\r\n"),
									},
									B: repo.BeginningData{Part: 2},
									E: repo.Probably,
								},
							},
						},
					},
				},
				d: &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 2, B: repo.False, E: repo.Probably}, APB: repo.AppPieceBody{B: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazaza")}},
				wantPresence: Presence{
					ASKG: true,
					ASKD: true,
					OB:   true,
					GMap: map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
						{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: true}: {
							true: {
								D: repo.Disposition{
									H: []byte("\r\n"),
								},
								B: repo.BeginningData{Part: 2},
								E: repo.Probably,
							},
						},
					},
				},
				wantError: nil,
			},

			{
				name: "B() == repo.True && E() == repo.Probably, askd met, askd.T() met",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "alice",
										FileName: "short.txt",
									},
									B: repo.BeginningData{Part: 0},
									E: repo.True,
								},
							},
							{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: true}: {
								true: {
									D: repo.Disposition{
										H: []byte("\r\n"),
									},
									B: repo.BeginningData{Part: 2},
									E: repo.Probably,
								},
							},
						},
					},
				},
				d: &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 2, B: repo.True, E: repo.Probably}, APB: repo.AppPieceBody{B: []byte("azaza")}},
				wantPresence: Presence{
					ASKG: true,
					ASKD: true,
					OB:   true,
					GMap: map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
						{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
							false: {
								D: repo.Disposition{
									H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.BeginningData{Part: 0},
								E: repo.True,
							},
						},
						{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: true}: {
							true: {
								D: repo.Disposition{
									H: []byte("\r\n"),
								},
								B: repo.BeginningData{Part: 2},
								E: repo.Probably,
							},
						},
					},
				},
				wantError: nil,
			},

			{
				name: "E() == repo.False, askd met, askd.T() met",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "alice",
										FileName: "short.txt",
									},
									B: repo.BeginningData{Part: 0},
									E: repo.True,
								},
							},
							{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: true}: {
								true: {
									D: repo.Disposition{
										H: []byte("\r\n"),
									},
									B: repo.BeginningData{Part: 2},
									E: repo.Probably,
								},
							},
						},
					},
				},
				d: &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 2, B: repo.True, E: repo.False}, APB: repo.AppPieceBody{B: []byte("azaza")}},
				wantPresence: Presence{
					ASKG: true,
					ASKD: true,
					OB:   false,
					GMap: map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
						{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
							false: {
								D: repo.Disposition{
									H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.BeginningData{Part: 0},
								E: repo.True,
							},
						},
					},
				},
				wantError: nil,
			},

			{
				name: "E() == repo.Probably, askd met, opposite branch met",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "alice",
										FileName: "short.txt",
									},
									B: repo.BeginningData{Part: 0},
									E: repo.True,
								},
							},
							{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "alice",
										FileName: "short.txt",
									},
									B: repo.BeginningData{Part: 1},
									E: repo.True,
								},
							},
							{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: true}: {
								true: {
									D: repo.Disposition{
										H: []byte("\r\n"),
									},
									B: repo.BeginningData{Part: 2},
									E: repo.Probably,
								},
							},
						},
					},
				},
				d: &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 2, B: repo.True, E: repo.Probably}, APB: repo.AppPieceBody{B: []byte("azaza")}},
				wantPresence: Presence{
					ASKG: true,
					ASKD: true,
					OB:   true,
					GMap: map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
						{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
							false: {
								D: repo.Disposition{
									H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.BeginningData{Part: 1},
								E: repo.True,
							},
						},
						{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: true}: {
							true: {
								D: repo.Disposition{
									H: []byte("\r\n"),
								},
								B: repo.BeginningData{Part: 2},
								E: repo.Probably,
							},
						},
					},
				},
				wantError: nil,
			},

			{
				name: "ASKD met, B() == repo.True, E() == repo.False, Cur == 1 && S.C.Blocked == true => enpty Presense",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "alice",
										FileName: "short.txt",
									},
									B: repo.BeginningData{Part: 0},
									E: repo.True,
								},
							},
						},
					},
					C: map[repo.AppStoreKeyGeneral]repo.Counter{
						{TS: "qqq"}: {Max: 3, Cur: 1, Blocked: true},
					},
				},
				d:            &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 2, B: repo.True, E: repo.True}, APB: repo.AppPieceBody{B: []byte("azaza")}},
				wantPresence: Presence{},
				wantError:    errors.New("in store.Presense matched but Cur == 1 && Blocked"),
			},

			{
				name: "ASKD met, B() == repo.True, E() == repo.False, Cur == 1 && Fuse == false => all trues",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "alice",
										FileName: "short.txt",
									},
									B: repo.BeginningData{Part: 0},
									E: repo.True,
								},
							},
						},
					},
					C: map[repo.AppStoreKeyGeneral]repo.Counter{
						{TS: "qqq"}: {Max: 3, Cur: 1, Blocked: false},
					},
				},
				d: &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 2, B: repo.True, E: repo.False}, APB: repo.AppPieceBody{B: []byte("azaza")}},
				wantPresence: Presence{
					ASKG: true,
					ASKD: true,
					GMap: map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
						{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
							false: {
								D: repo.Disposition{
									H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.BeginningData{Part: 0},
								E: repo.True,
							},
						},
					},
				},

				wantError: nil,
			},

			{
				name: "ASKD met, 2 specific branches, E() == repo.True",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 3}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "alice",
										FileName: "short.txt",
									},
									B: repo.BeginningData{Part: 0},
									E: repo.Probably,
								},
								true: {
									D: repo.Disposition{
										H: []byte("\r\n"),
									},
									B: repo.BeginningData{Part: 2},
									E: repo.Probably,
								},
							},
							{SK: repo.StreamKey{TS: "qqq", Part: 4}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "bob",
										FileName: "long.txt",
									},
									B: repo.BeginningData{Part: 3},
									E: repo.True,
								},
							},
						},
					},
				},
				d: &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 3, B: repo.True, E: repo.True}, APB: repo.AppPieceBody{B: []byte("azaza")}},
				wantPresence: Presence{
					ASKG: true,
					ASKD: true,
					OB:   false,
					GMap: map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
						{SK: repo.StreamKey{TS: "qqq", Part: 3}, S: false}: {
							false: {
								D: repo.Disposition{
									H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.BeginningData{Part: 0},
								E: repo.Probably,
							},
							true: {
								D: repo.Disposition{
									H: []byte("\r\n"),
								},
								B: repo.BeginningData{Part: 2},
								E: repo.Probably,
							},
						},
					},
				},
				wantError: nil,
			},

			{
				name: "ASKD met, 2 detailed branches, 2 specific branches, E() == repo.Probably",
				dataHandler: &memoryDataHandlerStruct{
					Map: map[key]map[bool]value{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 3}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "alice",
										FileName: "short.txt",
									},
									B: repo.BeginningData{Part: 0},
									E: repo.Probably,
								},
								true: {
									D: repo.Disposition{
										H: []byte("\r\n"),
									},
									B: repo.BeginningData{Part: 2},
									E: repo.Probably,
								},
							},
							{SK: repo.StreamKey{TS: "qqq", Part: 3}, S: true}: {
								true: {
									D: repo.Disposition{
										H: []byte("\r\n"),
									},
									B: repo.BeginningData{Part: 3},
									E: repo.Probably,
								},
							},
							{SK: repo.StreamKey{TS: "qqq", Part: 6}, S: false}: {
								false: {
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "bob",
										FileName: "long.txt",
									},
									B: repo.BeginningData{Part: 5},
									E: repo.True,
								},
							},
						},
					},
				},
				d: &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 3, B: repo.True, E: repo.Probably}, APB: repo.AppPieceBody{B: []byte("azaza")}},
				wantPresence: Presence{
					ASKG: true,
					ASKD: true,
					OB:   true,
					GMap: map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
						{SK: repo.StreamKey{TS: "qqq", Part: 3}, S: false}: {
							false: {
								D: repo.Disposition{
									H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.BeginningData{Part: 0},
								E: repo.Probably,
							},
							true: {
								D: repo.Disposition{
									H: []byte("\r\n"),
								},
								B: repo.BeginningData{Part: 2},
								E: repo.Probably,
							},
						},
						{SK: repo.StreamKey{TS: "qqq", Part: 3}, S: true}: {
							true: {
								D: repo.Disposition{
									H: []byte("\r\n"),
								},
								B: repo.BeginningData{Part: 3},
								E: repo.Probably,
							},
						},
					},
				},
				wantError: nil,
			},
		*/
	}

	for _, v := range tt {
		s.Run(v.name, func() {
			gotPresense, gotError := v.dataHandler.Check(v.d)
			if v.wantError != nil {
				s.Equal(v.wantError, gotError)
			}

			s.Equal(v.wantPresence, gotPresense)
		})
	}
}
