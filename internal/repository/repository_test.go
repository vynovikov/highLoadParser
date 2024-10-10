package repository

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"
	"github.com/vynovikov/highLoadParser/internal/dataHandler"
)

type repositorySuite struct {
	suite.Suite
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

type mockRedisDataHandler struct {
	calledNum  int
	parameters []interface{}
}

func (m *mockRedisDataHandler) Create(dataHandler.DataHandlerDTO, dataHandler.Boundary) (dataHandler.ProducerUnit, error) {
	return &dataHandler.ProducerUnitStruct{}, nil
}

func (m *mockRedisDataHandler) Read(dataHandler.DataHandlerDTO) (dataHandler.Value, error) {
	return dataHandler.Value{}, nil
}

func (m *mockRedisDataHandler) Updade(dataHandler.DataHandlerDTO, dataHandler.Boundary) (dataHandler.ProducerUnit, error) {
	return &dataHandler.ProducerUnitStruct{}, nil
}

func (m *mockRedisDataHandler) Delete(dataHandler.KeyDetailed) error {
	return nil
}

func (m *mockRedisDataHandler) Set(key dataHandler.KeyDetailed, val dataHandler.Value) error {
	m.calledNum++
	m.parameters = append(m.parameters, key, val)

	if strings.Contains(key.Ts, "error") {
		return fmt.Errorf("lalala")
	}

	return nil
}

func (m *mockRedisDataHandler) Get(key dataHandler.KeyDetailed) (dataHandler.Value, error) {
	m.calledNum++
	m.parameters = append(m.parameters, key)

	if strings.Contains(key.Ts, "error") {
		switch key.Ts[strings.Index(key.Ts, "_error_")+len("_error_"):] {
		case "redis_nil":
			return dataHandler.Value{}, errors.Join(dataHandler.ErrKeyNotFound, redis.Nil)
		}
		return dataHandler.Value{}, fmt.Errorf("other error")
	}

	return dataHandler.Value{}, nil
}

func (s *repositorySuite) TestRegister() {
	tt := []struct {
		name         string
		dto          RepositoryDTO
		initRepo     ParserRepository
		wantedRepo   ParserRepository
		wantedResult dataHandler.ProducerUnit
		wantedError  error
	}{
		{
			name: "1",
			dto: &RepositoryUnit{
				R_ts:   "qqq_error_redis_nil",
				R_part: 0,
				R_b:    0,
				R_e:    0,
				R_body: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nazazaza"),
				R_bou:  Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")},
			},
			initRepo: &repositoryStruct{
				dataHandler: &mockRedisDataHandler{},
			},

			wantedRepo: &repositoryStruct{
				dataHandler: &mockRedisDataHandler{
					calledNum: 2,

					parameters: []interface{}{
						dataHandler.KeyDetailed{
							Ts: "qqq_error_redis_nil",
						},
						dataHandler.KeyDetailed{
							Ts: "qqq_error_redis_nil",
						},
						dataHandler.Value{
							H: dataHandler.HeaderData{
								FormName: "alice",
								FileName: "",
								Header:   []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
							},
							E: 0,
						},
					},
				},
			},
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {

			_, gotErr := v.initRepo.Register(v.dto)
			if gotErr != nil {
				s.Equal(v.wantedError, gotErr)
			}
			s.Equal(v.wantedRepo, v.initRepo)
		})
	}
}

func (s *repositorySuite) TestGetHeaderLines() {
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

func (s *repositorySuite) TestNewValue() {

	tt := []struct {
		name        string
		dto         RepositoryDTO
		bou         Boundary
		wantedValue dataHandler.Value
		wantedError error
	}{

		{
			name: "1. Full header, name only",
			dto:  &RepositoryUnit{R_ts: "qqq", R_part: 0, R_body: []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\nazazaza"), R_b: 0, R_e: 1, R_isSub: false, R_last: false, R_bou: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}},
			wantedValue: dataHandler.Value{
				E: 1,
				H: dataHandler.HeaderData{
					FormName: "alice",
					FileName: "",
					Header:   []byte("Content-Disposition: form-data; name=\"alice\"\r\n\r\n"),
				}},
		},

		{
			name: "2. Full header, name + filename",
			dto:  &RepositoryUnit{R_ts: "qqq", R_part: 0, R_body: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\nazazaza"), R_b: 0, R_e: 1, R_isSub: false, R_last: false, R_bou: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}},
			wantedValue: dataHandler.Value{
				E: 1,
				H: dataHandler.HeaderData{
					FormName: "alice",
					FileName: "short.txt",
					Header:   []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r\n"),
				}},
		},

		{
			name: "3. Partial header",
			dto:  &RepositoryUnit{R_ts: "qqq", R_part: 0, R_body: []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r"), R_b: 0, R_e: 1, R_isSub: false, R_last: false, R_bou: Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}},
			wantedValue: dataHandler.Value{
				E: 1,
				H: dataHandler.HeaderData{
					FormName: "",
					FileName: "",
					Header:   []byte("Content-Disposition: form-data; name=\"alice\"; filename=\"short.txt\"\r\nContent-Type: text/plain\r\n\r"),
				}},
		},
	}

	for _, v := range tt {

		s.Run(v.name, func() {

			gotValue, err := newValue(v.dto)

			if v.wantedError != nil {

				s.Equal(v.wantedError, err)
			}

			s.Equal(v.wantedValue, gotValue)

		})
	}
}

func (s *repositorySuite) TestGenBoundary() {
	boundaryVoc := Boundary{Prefix: []byte("bPrefix"), Root: []byte("bRoot")}

	boundaryCalced := genBoundary(boundaryVoc)

	s.Equal([]byte("\r\nbPrefix"+"bRoot"), boundaryCalced)
}

func (s *repositorySuite) TestIsLastBoundaryPart() {
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
