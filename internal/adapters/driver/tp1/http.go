// HTTP Reciever.
package tp1

import (
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/vynovikov/study/highLoadParser/internal/adapters/application"
	"github.com/vynovikov/study/highLoadParser/internal/logger"
	"github.com/vynovikov/study/highLoadParser/internal/repo"
)

type TpServer struct {
	l net.Listener
}

type TpReceiver interface {
	Run()
	HandleRequestFull(net.Conn, string, *sync.WaitGroup)
	Stop(*sync.WaitGroup)
}

type tpReceiverStruct struct {
	A     application.Application
	Saved map[string]struct{}
	srv   *TpServer
	wg    sync.WaitGroup
}

func NewTpReceiver(a application.Application) *tpReceiverStruct {

	li, err := net.Listen("tcp", ":3000")
	if err != nil {
		logger.L.Error(err)
	}
	logger.L.Info("listening localhost:3000")

	s := &TpServer{l: li}

	return &tpReceiverStruct{
		A:     a,
		Saved: make(map[string]struct{}),
		srv:   s,
	}
}

func (r *tpReceiverStruct) Run() {
	for {
		conn, err := r.srv.l.Accept()

		if err != nil && conn == nil {

			r.wg.Wait()
			r.A.ChanInClose()

			return

		}

		r.wg.Add(1)
		ts := repo.NewTS()
		// consider serial execution -> mutex
		go r.HandleRequestFull(conn, ts, &r.wg)

	}

}

// Tested in http_test.go
func (r *tpReceiverStruct) HandleRequestFull(conn net.Conn, ts string, wg *sync.WaitGroup) {
	defer wg.Done()
	bou, header, errFirst := repo.AnalyzeHeader(conn)
	if errFirst != nil && strings.Contains(errFirst.Error(), "100-continue") {
		r.Saved[string(repo.GenBoundary(bou)[2:])] = struct{}{}
		var wwg sync.WaitGroup
		wwg.Add(1)
		go repo.RespondContinue(conn, &wwg)
		wwg.Wait()

		header = append(repo.GenBoundary(bou)[2:], []byte("\r\n")...)
	}
	r.HandleRequestLast(conn, ts, bou, header, errFirst)

	if r.A.Stopping() {
		r.A.ChanInClose()
	}
	return

}
func (r *tpReceiverStruct) HandleRequestLast(conn net.Conn, ts string, bou repo.Boundary, header []byte, errFirst error) {
	p := 0
	for {
		h := repo.NewReceiverHeader(ts, p, bou)
		b, errSecond := repo.AnalyzeBits(conn, 1024, p, header, errFirst)
		//logger.L.Infof("in HandleReauestLast b = %q, errSecond = %v\n", b, errSecond)

		u := repo.NewReceiverUnit(h, b)
		if errFirst != nil {

			if errFirst == io.EOF || errFirst == io.ErrUnexpectedEOF || os.IsTimeout(errFirst) {
				r.A.AddToFeeder(u)
				break
			}
		}
		if errSecond != nil {
			if errSecond == io.EOF || errSecond == io.ErrUnexpectedEOF || os.IsTimeout(errSecond) {
				r.A.AddToFeeder(u)
				break
			}
			if strings.Contains(errSecond.Error(), "empty") {
				break
			}
		}
		//logger.L.Infof("in HandleReauest sending to feeder %v\n", u)
		r.A.AddToFeeder(u)

		p++
	}
	var wwg sync.WaitGroup
	wwg.Add(1)
	go repo.Respond(conn, &wwg)
	wwg.Wait()
}

func (r *tpReceiverStruct) CleanSaved(s string) {
	if len(r.Saved) > 1 {
		delete(r.Saved, s)
	}
	r.Saved = make(map[string]struct{})
}

func (r *tpReceiverStruct) Stop(wg *sync.WaitGroup) {

	r.srv.l.Close()

	r.wg.Wait()

	wg.Done()
}
