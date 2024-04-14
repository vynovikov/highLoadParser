package tps

import (
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vynovikov/study/highLoadParser/internal/adapters/application"
	"github.com/vynovikov/study/highLoadParser/internal/repo"
)

type tpsSuite struct {
	suite.Suite
}

func TestTpsSuite(t *testing.T) {
	suite.Run(t, new(tpsSuite))
}

var a application.AppService

func (s *tpsSuite) SetupTest() {
	a = application.NewAppService(make(chan struct{}))
}

// TestHandleRequest tests work of tps recriver. Testdouble spy is used to evaluate corectness of reciever operation
func (s *tpsSuite) TestHandleRequest() {

	tt := []struct {
		name    string
		R       TpsReceiver
		cl      net.Conn
		sr      net.Conn
		req     string
		wg      sync.WaitGroup
		TS      string
		wantR   TpsReceiver
		wantRes []byte
	}{
		{
			name: "len(req) < 512",
			R: &tpsReceiverStruct{
				A: &application.App{
					A: a,
					L: &SpyLogger{},
				},
			},
			req: "POST / HTTP/1.1\r\n" +
				"Host: localhost\r\n" +
				"User-Agent: curl/7.75.0\r\n" +
				"Accept: */*\r\n" +
				"Content-Length: 5250\r\n" +
				"Content-Type: multipart/form-data; boundary=------------------------c61fd8e07a9d3f9b\r\n" +
				"\r\n" +
				"--------------------------c61fd8e07a9d3f9b\r\n" +
				"Content-Disposition: form-data; name=\"alice\"\r\n" +
				"\r\n" +
				"azaza\r\n" +
				"--------------------------c61fd8e07a9d3f9b--",
			TS: "qqq",
			wantR: &tpsReceiverStruct{
				A: &application.App{
					A: a,
					L: &SpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.ReceiverUnit{
								H: repo.ReceiverHeader{
									TS:   "qqq",
									Part: 0,
									Bou:  repo.Boundary{Prefix: []byte("--"), Root: []byte("------------------------c61fd8e07a9d3f9b")},
								},
								B: repo.ReceiverBody{
									B: []byte(
										"POST / HTTP/1.1\r\n" +
											"Host: localhost\r\n" +
											"User-Agent: curl/7.75.0\r\n" +
											"Accept: */*\r\n" +
											"Content-Length: 5250\r\n" +
											"Content-Type: multipart/form-data; boundary=------------------------c61fd8e07a9d3f9b\r\n" +
											"\r\n" +
											"--------------------------c61fd8e07a9d3f9b\r\n" +
											"Content-Disposition: form-data; name=\"alice\"\r\n" +
											"\r\n" +
											"azaza\r\n" +
											"--------------------------c61fd8e07a9d3f9b--"),
								},
							},
						},
					},
				},
			},
			wantRes: []byte("HTTP/1.1 200 OK\r\n" +
				"Content-Length: 6\r\n" +
				"Content-Type: text/html\r\n" +
				"\r\n" +
				"200 OK"),
		},

		{
			name: "len(req) > 512 && len(req) < 1024",
			R: &tpsReceiverStruct{
				A: &application.App{
					A: a,
					L: &SpyLogger{},
				},
			},
			req: "POST / HTTP/1.1\r\n" +
				"Host: localhost\r\n" +
				"User-Agent: curl/7.75.0\r\n" +
				"Accept: */*\r\n" +
				"Content-Length: 5250\r\n" +
				"Content-Type: multipart/form-data; boundary=------------------------c61fd8e07a9d3f9b\r\n" +
				"\r\n" +
				"--------------------------c61fd8e07a9d3f9b\r\n" +
				"Content-Disposition: form-data; name=\"alice\"; filename=\"long.txt\"\r\n" +
				"Content-Type: text/plain\r\n" +
				"\r\n" +
				"0\r\n" +
				"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\r\n" +
				"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111\r\n" +
				"222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222\r\n" +
				"333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333\r\n" +
				"444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444\r\n" +
				"555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555\r\n" +
				"666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666",
			TS: "qqq",
			wantR: &tpsReceiverStruct{
				A: &application.App{
					A: a,
					L: &SpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.ReceiverUnit{
								H: repo.ReceiverHeader{
									TS:   "qqq",
									Part: 0,
									Bou:  repo.Boundary{Prefix: []byte("--"), Root: []byte("------------------------c61fd8e07a9d3f9b")},
								},
								B: repo.ReceiverBody{
									B: []byte(
										"POST / HTTP/1.1\r\n" +
											"Host: localhost\r\n" +
											"User-Agent: curl/7.75.0\r\n" +
											"Accept: */*\r\n" +
											"Content-Length: 5250\r\n" +
											"Content-Type: multipart/form-data; boundary=------------------------c61fd8e07a9d3f9b\r\n" +
											"\r\n" +
											"--------------------------c61fd8e07a9d3f9b\r\n" +
											"Content-Disposition: form-data; name=\"alice\"; filename=\"long.txt\"\r\n" +
											"Content-Type: text/plain\r\n" +
											"\r\n" +
											"0\r\n" +
											"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\r\n" +
											"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111\r\n" +
											"222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222\r\n" +
											"333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333\r\n" +
											"444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444\r\n" +
											"555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555\r\n" +
											"666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666"),
								},
							},
						},
					},
				},
			},
			wantRes: []byte("HTTP/1.1 200 OK\r\n" +
				"Content-Length: 6\r\n" +
				"Content-Type: text/html\r\n" +
				"\r\n" +
				"200 OK"),
		},

		{
			name: "len(req) == 1024",
			R: &tpsReceiverStruct{
				A: &application.App{
					A: a,
					L: &SpyLogger{},
				},
			},
			req: "POST / HTTP/1.1\r\n" +
				"Host: localhost\r\n" +
				"User-Agent: curl/7.75.0\r\n" +
				"Accept: */*\r\n" +
				"Content-Length: 5250\r\n" +
				"Content-Type: multipart/form-data; boundary=------------------------c61fd8e07a9d3f9b\r\n" +
				"\r\n" +
				"--------------------------c61fd8e07a9d3f9b\r\n" +
				"Content-Disposition: form-data; name=\"alice\"; filename=\"long.txt\"\r\n" +
				"Content-Type: text/plain\r\n" +
				"\r\n" +
				"0\r\n" +
				"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\r\n" +
				"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111\r\n" +
				"222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222\r\n" +
				"333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333\r\n" +
				"444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444\r\n" +
				"555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555\r\n" +
				"6666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666",
			TS: "qqq",
			wantR: &tpsReceiverStruct{
				A: &application.App{
					A: a,
					L: &SpyLogger{
						calls: 1,
						params: []repo.AppUnit{
							repo.ReceiverUnit{
								H: repo.ReceiverHeader{
									TS:   "qqq",
									Part: 0,
									Bou:  repo.Boundary{Prefix: []byte("--"), Root: []byte("------------------------c61fd8e07a9d3f9b")},
								},
								B: repo.ReceiverBody{
									B: []byte(
										"POST / HTTP/1.1\r\n" +
											"Host: localhost\r\n" +
											"User-Agent: curl/7.75.0\r\n" +
											"Accept: */*\r\n" +
											"Content-Length: 5250\r\n" +
											"Content-Type: multipart/form-data; boundary=------------------------c61fd8e07a9d3f9b\r\n" +
											"\r\n" +
											"--------------------------c61fd8e07a9d3f9b\r\n" +
											"Content-Disposition: form-data; name=\"alice\"; filename=\"long.txt\"\r\n" +
											"Content-Type: text/plain\r\n" +
											"\r\n" +
											"0\r\n" +
											"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\r\n" +
											"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111\r\n" +
											"222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222\r\n" +
											"333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333\r\n" +
											"444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444\r\n" +
											"555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555\r\n" +
											"6666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666"),
								},
							},
						},
					},
				},
			},
			wantRes: []byte("HTTP/1.1 200 OK\r\n" +
				"Content-Length: 6\r\n" +
				"Content-Type: text/html\r\n" +
				"\r\n" +
				"200 OK"),
		},

		{
			name: "len(req) > 1024 && len(req) < 2048",
			R: &tpsReceiverStruct{
				A: &application.App{
					A: a,
					L: &SpyLogger{},
				},
			},
			req: "POST / HTTP/1.1\r\n" +
				"Host: localhost\r\n" +
				"User-Agent: curl/7.75.0\r\n" +
				"Accept: */*\r\n" +
				"Content-Length: 5250\r\n" +
				"Content-Type: multipart/form-data; boundary=------------------------c61fd8e07a9d3f9b\r\n" +
				"\r\n" +
				"--------------------------c61fd8e07a9d3f9b\r\n" +
				"Content-Disposition: form-data; name=\"alice\"; filename=\"long.txt\"\r\n" +
				"Content-Type: text/plain\r\n" +
				"\r\n" +
				"0\r\n" +
				"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\r\n" +
				"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111\r\n" +
				"222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222\r\n" +
				"333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333\r\n" +
				"444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444\r\n" +
				"555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555\r\n" +
				"66666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666",
			TS: "qqq",
			wantR: &tpsReceiverStruct{
				A: &application.App{
					A: a,
					L: &SpyLogger{
						calls: 2,
						params: []repo.AppUnit{
							repo.ReceiverUnit{
								H: repo.ReceiverHeader{
									TS:   "qqq",
									Part: 0,
									Bou:  repo.Boundary{Prefix: []byte("--"), Root: []byte("------------------------c61fd8e07a9d3f9b")},
								},
								B: repo.ReceiverBody{
									B: []byte(
										"POST / HTTP/1.1\r\n" +
											"Host: localhost\r\n" +
											"User-Agent: curl/7.75.0\r\n" +
											"Accept: */*\r\n" +
											"Content-Length: 5250\r\n" +
											"Content-Type: multipart/form-data; boundary=------------------------c61fd8e07a9d3f9b\r\n" +
											"\r\n" +
											"--------------------------c61fd8e07a9d3f9b\r\n" +
											"Content-Disposition: form-data; name=\"alice\"; filename=\"long.txt\"\r\n" +
											"Content-Type: text/plain\r\n" +
											"\r\n" +
											"0\r\n" +
											"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\r\n" +
											"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111\r\n" +
											"222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222\r\n" +
											"333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333\r\n" +
											"444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444\r\n" +
											"555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555\r\n" +
											"6666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666")}},
							repo.ReceiverUnit{
								H: repo.ReceiverHeader{
									TS:   "qqq",
									Part: 1,
									Bou:  repo.Boundary{Prefix: []byte("--"), Root: []byte("------------------------c61fd8e07a9d3f9b")},
								},
								B: repo.ReceiverBody{
									B: []byte("6")}},
						},
					},
				},
			},
			wantRes: []byte("HTTP/1.1 200 OK\r\n" +
				"Content-Length: 6\r\n" +
				"Content-Type: text/html\r\n" +
				"\r\n" +
				"200 OK"),
		},

		{
			name: "len(req) == 2048",
			R: &tpsReceiverStruct{
				A: &application.App{
					A: a,
					L: &SpyLogger{},
				},
			},
			req: "POST / HTTP/1.1\r\n" +
				"Host: localhost\r\n" +
				"User-Agent: curl/7.75.0\r\n" +
				"Accept: */*\r\n" +
				"Content-Length: 5250\r\n" +
				"Content-Type: multipart/form-data; boundary=------------------------c61fd8e07a9d3f9b\r\n" +
				"\r\n" +
				"--------------------------c61fd8e07a9d3f9b\r\n" +
				"Content-Disposition: form-data; name=\"alice\"; filename=\"long.txt\"\r\n" +
				"Content-Type: text/plain\r\n" +
				"\r\n" +
				"0\r\n" +
				"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\r\n" +
				"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111\r\n" +
				"222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222\r\n" +
				"333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333\r\n" +
				"444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444\r\n" +
				"555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555\r\n" +
				"666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666\r\n" +
				"777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777\r\n" +
				"888888888888888888888888888888888888888888888888888888888888888888888888888888888888888888888888888\r\n" +
				"999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999\r\n" +
				"1\r\n" +
				"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\r\n" +
				"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111\r\n" +
				"222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222\r\n",
			TS: "qqq",
			wantR: &tpsReceiverStruct{
				A: &application.App{
					A: a,
					L: &SpyLogger{
						calls: 2,
						params: []repo.AppUnit{
							repo.ReceiverUnit{
								H: repo.ReceiverHeader{
									TS:   "qqq",
									Part: 0,
									Bou:  repo.Boundary{Prefix: []byte("--"), Root: []byte("------------------------c61fd8e07a9d3f9b")},
								},
								B: repo.ReceiverBody{
									B: []byte(
										"POST / HTTP/1.1\r\n" +
											"Host: localhost\r\n" +
											"User-Agent: curl/7.75.0\r\n" +
											"Accept: */*\r\n" +
											"Content-Length: 5250\r\n" +
											"Content-Type: multipart/form-data; boundary=------------------------c61fd8e07a9d3f9b\r\n" +
											"\r\n" +
											"--------------------------c61fd8e07a9d3f9b\r\n" +
											"Content-Disposition: form-data; name=\"alice\"; filename=\"long.txt\"\r\n" +
											"Content-Type: text/plain\r\n" +
											"\r\n" +
											"0\r\n" +
											"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\r\n" +
											"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111\r\n" +
											"222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222\r\n" +
											"333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333\r\n" +
											"444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444\r\n" +
											"555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555555\r\n" +
											"6666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666")}},
							repo.ReceiverUnit{
								H: repo.ReceiverHeader{
									TS:   "qqq",
									Part: 1,
									Bou:  repo.Boundary{Prefix: []byte("--"), Root: []byte("------------------------c61fd8e07a9d3f9b")},
								},
								B: repo.ReceiverBody{
									B: []byte(
										"66666\r\n" +
											"777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777\r\n" +
											"888888888888888888888888888888888888888888888888888888888888888888888888888888888888888888888888888\r\n" +
											"999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999\r\n" +
											"1\r\n" +
											"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\r\n" +
											"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111\r\n" +
											"222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222\r\n")}},
						},
					},
				},
			},
			wantRes: []byte("HTTP/1.1 200 OK\r\n" +
				"Content-Length: 6\r\n" +
				"Content-Type: text/html\r\n" +
				"\r\n" +
				"200 OK"),
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {

			v.cl, v.sr = net.Pipe()

			v.wg.Add(1)

			go v.R.HandleRequestFull(v.sr, v.TS, &v.wg)

			fmt.Fprint(v.cl, v.req)
			time.Sleep(time.Millisecond * 50)
			s.Equal(v.wantRes, GetResponse(v.cl))

			v.wg.Wait()
			s.Equal(v.wantR, v.R)
			v.cl.Close()
			v.sr.Close()
		})
	}
}

type SpyLogger struct {
	calls  int
	params []repo.AppUnit
}

func (s *SpyLogger) LogStuff(ru repo.AppUnit) {
	s.calls++
	s.params = append(s.params, ru)
}

func GetResponse(conn net.Conn) []byte {
	r := make([]byte, 200)
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 25))
	n, _ := io.ReadFull(conn, r)
	if n < len(r) {
		r = r[:n]
	}
	return r
}
