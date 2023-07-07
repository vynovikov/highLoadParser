package application

import (
	"sync"
	"testing"

	"github.com/vynovikov/highLoadParser/internal/adapters/driven/store"
	"github.com/vynovikov/highLoadParser/internal/repo"

	"github.com/stretchr/testify/suite"
)

var (
	a AppService
)

type applicationSuite struct {
	suite.Suite
}

func TestApplicationSuite(t *testing.T) {
	suite.Run(t, new(applicationSuite))
}

func (s *applicationSuite) SetupTest() {
	a = NewAppService(make(chan struct{}))
	//a.MountLogger(NewDistributorSpyLogger())

}

func (s *applicationSuite) TestHandle() {
	go func() {
		for {
			select {
			case <-a.C.ChanOut:
			case <-a.C.ChanLog:

			}
		}
	}()
	tt := []struct {
		name  string
		a     *App
		d     repo.DataPiece
		bou   repo.Boundary
		wg    sync.WaitGroup
		wantA *App
	}{
		/*
			{
				name: "B() == repo.False, E() == repo.False fileName absent",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{},
				},
				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 3, TS: "qqq", B: repo.False, E: repo.False}, APB: repo.AppPieceBody{B: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nazazaza")},
				},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "B() == repo.False, E() == repo.Last fileName absent",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{},
				},
				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 3, TS: "qqq", B: repo.False, E: repo.Last}, APB: repo.AppPieceBody{B: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nazazaza")},
				},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
									IsLast:   true,
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "B() == repo.False, E() == repo.Last fileName present",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{},
				},
				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 3, TS: "qqq", B: repo.False, E: repo.Last}, APB: repo.AppPieceBody{B: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nazazaza")},
				},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
									IsLast:   true,
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "B() == repo.False, E() == repo.False fileName present",
				a: &App{
					A: a,
					L: &DistributorSpyLogger{},
				},

				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 3, TS: "qqq", B: repo.False, E: repo.False}, APB: repo.AppPieceBody{B: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza")},
				},
				wantA: &App{
					A: a,
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "B() == repo.False, E() == repo.True, header full, intermediate dataPiece",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{},
				},

				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 0, TS: "qqq", B: repo.False, E: repo.True}, APB: repo.AppPieceBody{B: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
									false: repo.AppStoreValue{
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
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "B() == repo.False, E() == repo.True, header not full",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{},
				},

				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 0, TS: "qqq", B: repo.False, E: repo.True}, APB: repo.AppPieceBody{B: []byte("Content-Disposition: form-data; name=\"alice\"; file")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H: []byte("Content-Disposition: form-data; name=\"alice\"; file"),
										},
										B: repo.BeginningData{Part: 0},
										E: repo.True,
									},
								},
							},
						},
					},
					L: &DistributorSpyLogger{},
				},
			},

			{
				name: "B() == repo.True, E() == repo.True, ending part of header",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H: []byte("Content-Disposition: form-data; name=\"alice\"; file"),
										},
										B: repo.BeginningData{Part: 0},
										E: repo.True,
									},
								},
							},
						},
					},
					L: &DistributorSpyLogger{},
				},

				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 1, TS: "qqq", B: repo.True, E: repo.True}, APB: repo.AppPieceBody{B: []byte("name=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
									false: repo.AppStoreValue{
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
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "B() == repo.True, E() == repo.False, header was full",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
									false: repo.AppStoreValue{
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
					L: &DistributorSpyLogger{},
				},

				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 1, TS: "qqq", B: repo.True, E: repo.False}, APB: repo.AppPieceBody{B: []byte("azazaza")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "B() == repo.True, E() == repo.Last, ending part of header",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H: []byte("Content-Disposition: form-data; name=\"alice\"; file"),
										},
										B: repo.BeginningData{Part: 0},
										E: repo.True,
									},
								},
							},
						},
					},
					L: &DistributorSpyLogger{},
				},

				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 1, TS: "qqq", B: repo.True, E: repo.Last}, APB: repo.AppPieceBody{B: []byte("name=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
									FileName: "short.txt",
									IsLast:   true,
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "B() == repo.True, E() == repo.Last, header was full",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
									false: repo.AppStoreValue{
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
					L: &DistributorSpyLogger{},
				},

				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 1, TS: "qqq", B: repo.True, E: repo.Last}, APB: repo.AppPieceBody{B: []byte("azazaza")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
									FileName: "short.txt",
									IsLast:   true,
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "APU B() == repo.True, E() == repo.Probably => adding to store, sending valid body part",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{},
				},
				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 1, TS: "qqq", B: repo.True, E: repo.Probably}, APB: repo.AppPieceBody{B: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
									false: repo.AppStoreValue{
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
								},
							},
						},
					},
					L: &DistributorSpyLogger{

						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "AppSub, APU absent => adding to store, part remains",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{},
				},
				d:   &repo.AppSub{ASH: repo.AppSubHeader{TS: "qqq", Part: 1}, ASB: repo.AppSubBody{B: []byte("\r")}},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: true}: {
									true: repo.AppStoreValue{
										D: repo.Disposition{
											H: []byte("\r"),
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
								},
							},
						},
					},
					L: &DistributorSpyLogger{},
				},
			},

			{
				name: "AppSub, APU present => adding to store, part increments",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 1}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
											FormName: "alice",
											FileName: "short.txt",
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
								},
							},
						},
					},

					L: &DistributorSpyLogger{},
				},
				d:   &repo.AppSub{ASH: repo.AppSubHeader{TS: "qqq", Part: 1}, ASB: repo.AppSubBody{B: []byte("\r")}},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
											FormName: "alice",
											FileName: "short.txt",
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
									true: repo.AppStoreValue{
										D: repo.Disposition{
											H: []byte("\r"),
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
								},
							},
						},
					},
					L: &DistributorSpyLogger{},
				},
			},

			{
				name: "Confirming DataPiece, boundary present, E() == repo.True => updating store, new ADU data header",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
											FormName: "alice",
											FileName: "short.txt",
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
									true: repo.AppStoreValue{
										D: repo.Disposition{
											H: []byte("\r"),
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
								},
							},
						},
					},

					L: &DistributorSpyLogger{},
				},
				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 2, TS: "qqq", B: repo.True, E: repo.True}, APB: repo.AppPieceBody{B: []byte("\n--bRoot\r\nContent-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 3}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H:        []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
											FormName: "bob",
											FileName: "long.txt",
										},
										B: repo.BeginningData{Part: 2},
										E: repo.True,
									},
								},
							},
						},
					},
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "bob",
									FileName: "long.txt",
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "Confirming DataPiece, boundary present, E() == repo.Last => updating store, new ADU data header",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
											FormName: "alice",
											FileName: "short.txt",
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
									true: repo.AppStoreValue{
										D: repo.Disposition{
											H: []byte("\r"),
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
								},
							},
						},
					},

					L: &DistributorSpyLogger{},
				},
				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 2, TS: "qqq", B: repo.True, E: repo.Last}, APB: repo.AppPieceBody{B: []byte("\n--bRoot\r\nContent-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "bob",
									FileName: "long.txt",
									IsLast:   true,
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},
		*/
		{
			name: "Confirming DataPiece, last boundary present, root separated",
			a: &App{
				A: a,
				S: &store.StoreStruct{
					R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
						{TS: "qqq"}: {
							{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
								false: repo.AppStoreValue{
									D: repo.Disposition{
										H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
										FormName: "alice",
										FileName: "short.txt",
									},
									B: repo.BeginningData{Part: 1},
									E: repo.Probably,
								},
								true: repo.AppStoreValue{
									D: repo.Disposition{
										H: []byte("\r\n--bRo"),
									},
									B: repo.BeginningData{Part: 1},
									E: repo.Probably,
								},
							},
						},
					},
				},

				L: &DistributorSpyLogger{},
			},
			d: &repo.AppPieceUnit{
				APH: repo.AppPieceHeader{Part: 2, TS: "qqq", B: repo.True, E: repo.Last}, APB: repo.AppPieceBody{B: []byte("ot--\r\n")},
			},
			bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
			wantA: &App{
				A: a,
				S: &store.StoreStruct{
					R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
				},
				L: &DistributorSpyLogger{
					calls: 1,
					params: []repo.AppUnit{
						repo.AppDistributorUnit{
							H: repo.AppDistributorHeader{
								TS:     "qqq",
								IsLast: true,
							},
						},
					},
				},
			},
		},
		/*
			{
				name: "Confirming DataPiece, boundary present and separated => updating store, ADU new header",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
											FormName: "alice",
											FileName: "short.txt",
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
									true: repo.AppStoreValue{
										D: repo.Disposition{
											H: []byte("\r\n-"),
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
								},
							},
						},
					},

					L: &DistributorSpyLogger{},
				},
				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 2, TS: "qqq", B: repo.True, E: repo.True}, APB: repo.AppPieceBody{B: []byte("-bRoot\r\nContent-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 3}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H:        []byte("Content-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
											FormName: "bob",
											FileName: "long.txt",
										},
										B: repo.BeginningData{Part: 2},
										E: repo.True,
									},
								},
							},
						},
					},
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "bob",
									FileName: "long.txt",
								},
								B: repo.AppDistributorBody{
									B: []byte("azazaza"),
								},
							},
						},
					},
				},
			},

			{
				name: "Confirming DataPiece, boundary not present => updating store, ADU header unchanged",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
											FormName: "alice",
											FileName: "short.txt",
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
									true: repo.AppStoreValue{
										D: repo.Disposition{
											H: []byte("\r\n-"),
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
								},
							},
						},
					},

					L: &DistributorSpyLogger{},
				},
				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 2, TS: "qqq", B: repo.True, E: repo.True}, APB: repo.AppPieceBody{B: []byte("-razaza\r\nbzbzb")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 3}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
											FormName: "alice",
											FileName: "short.txt",
										},
										B: repo.BeginningData{Part: 1},
										E: repo.True,
									},
								},
							},
						},
					},
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
									FileName: "short.txt",
								},
								B: repo.AppDistributorBody{
									B: []byte("\r\n--razaza\r\nbzbzb"),
								},
							},
						},
					},
				},
			},

			{
				name: "Confirming DataPiece, E() == repo.Last, boundary suffix => sending ADU with empty body and IsLast",
				a: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{
							{TS: "qqq"}: {
								{SK: repo.StreamKey{TS: "qqq", Part: 2}, S: false}: {
									false: repo.AppStoreValue{
										D: repo.Disposition{
											H:        []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
											FormName: "alice",
											FileName: "short.txt",
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
									true: repo.AppStoreValue{
										D: repo.Disposition{
											H: []byte("\r\n--bRoot"),
										},
										B: repo.BeginningData{Part: 1},
										E: repo.Probably,
									},
								},
							},
						},
					},

					L: &DistributorSpyLogger{},
				},
				d: &repo.AppPieceUnit{
					APH: repo.AppPieceHeader{Part: 2, TS: "qqq", B: repo.True, E: repo.Last}, APB: repo.AppPieceBody{B: []byte("--")},
				},
				bou: repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
				wantA: &App{
					A: a,
					S: &store.StoreStruct{
						R: map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{},
					},
					L: &DistributorSpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.AppDistributorUnit{
								H: repo.AppDistributorHeader{
									TS:       "qqq",
									FormName: "alice",
									FileName: "short.txt",
									IsLast:   true,
								},
							},
						},
					},
				},
			},
		*/
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			v.a.Handle(v.d, v.bou)
			//logger.L.Infoln("in application.TestHandle waiting...")

			s.Equal(v.wantA, v.a)
		})
	}
}
func (s *applicationSuite) TestCalcBody() {
	tt := []struct {
		name string

		d          repo.DataPiece
		bou        repo.Boundary
		wantADUB   repo.AppDistributorBody
		wantHeader []byte
		wantError  error
	}{
		{
			name:     "no header present",
			d:        &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 0, B: repo.True, E: repo.False}, APB: repo.AppPieceBody{B: []byte("azaza")}},
			bou:      repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
			wantADUB: repo.AppDistributorBody{B: []byte("azaza")},
		},

		{
			name:       "header present",
			d:          &repo.AppPieceUnit{APH: repo.AppPieceHeader{TS: "qqq", Part: 0, B: repo.True, E: repo.False}, APB: repo.AppPieceBody{B: []byte("--bRoot\r\nContent-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\nbzbzbzbzb")}},
			bou:        repo.Boundary{Prefix: []byte("--"), Root: []byte("bRoot")},
			wantADUB:   repo.AppDistributorBody{B: []byte("bzbzbzbzb")},
			wantHeader: []byte("--bRoot\r\nContent-Disposition: form-data; name=\"bob\"; filename=\"long.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			gotADUB, gotHeader, gotErr := CalcBody(v.d, v.bou)
			if v.wantError != nil {
				s.Equal(v.wantError, gotErr)
			}
			s.Equal(v.wantADUB, gotADUB)
			s.Equal(v.wantHeader, gotHeader)

		})
	}
}

type DistributorSpyLogger struct {
	calls  int
	params []repo.AppUnit
}

func (d *DistributorSpyLogger) LogStuff(au repo.AppUnit) {
	d.calls++
	d.params = append(d.params, au)
}

func NewDistributorSpyLogger() *DistributorSpyLogger {
	return &DistributorSpyLogger{}
}
