// HTTP Reciever.
package tp

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/vynovikov/highLoadParser/internal/adapters/application"
	"github.com/vynovikov/highLoadParser/internal/logger"
	"github.com/vynovikov/highLoadParser/internal/repo"
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
		logger.L.Infoln("got request")
		if err != nil && conn == nil {

			r.wg.Wait()
			r.A.ChanInClose()

			return

		}
		logger.L.Infoln("got request")
		r.wg.Add(1)
		ts := repo.NewTS()
		// consider serial execution -> mutex
		go r.HandleRequestFull(conn, ts, &r.wg)

	}

}

// Tested in http_test.go
func (r *tpReceiverStruct) HandleRequestFull(conn net.Conn, ts string, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(r.Saved) == 0 {
		bou, header, errFirst := repo.AnalyzeHeader(conn)
		logger.L.Infof("in HandleReauestFull 1 bou = %q, header = %q, errFirst = %v\n", bou, header, errFirst)
		if errFirst != nil && strings.Contains(errFirst.Error(), "100-continue") {
			r.Saved[string(repo.GenBoundary(bou)[2:])] = struct{}{}
			var wwg sync.WaitGroup
			wwg.Add(1)
			go repo.RespondContinue(conn, &wwg)
			wwg.Wait()

			header = append(repo.GenBoundary(bou)[2:], []byte("\r\n")...)
			//return
		}

		//go repo.Respond(conn)
		r.HandleRequestLast(conn, ts, bou, header, errFirst)

		if r.A.Stopping() {
			r.A.ChanInClose()
		}
		return
	}
	// len(r.Saved) > 0
	bouLen := 0
	for i, _ := range r.Saved {
		bouLen = len(i)
		break
	}
	firstN, err := repo.ReadFirst(conn, bouLen+4)
	if err != nil {
		logger.L.Errorf("in tp.HandleRequestFull error reading request %v\n", err)
	}
	c := string(firstN[2 : len(firstN)-2])
	if _, ok := r.Saved[c]; ok {
		r.HandleRequestLast(conn, ts, repo.NewBoundary(firstN[:2], firstN[2:bouLen+2]), []byte(""), fmt.Errorf("in tp.HandleRequestFull rest part after 100-continue"))
		r.CleanSaved(c)
	}

}
func (r *tpReceiverStruct) HandleRequestLast(conn net.Conn, ts string, bou repo.Boundary, header []byte, errFirst error) {
	p := 0
	for {
		h := repo.NewReceiverHeader(ts, p, bou)
		b, errSecond := repo.AnalyzeBits(conn, 1024, p, header, errFirst)
		logger.L.Infof("in HandleReauestLast b = %q, errSecond = %v\n", b, errSecond)

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
