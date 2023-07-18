// HTTPS receiver.
//
// x509 pair should be in tls forlder inside root directory
package tps

import (
	"crypto/tls"
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

type TpsServer struct {
	l net.Listener
}

type TpsReceiver interface {
	Run()
	HandleRequestFull(conn net.Conn, ts string, wg *sync.WaitGroup)
	Stop(*sync.WaitGroup)
}

type tpsReceiverStruct struct {
	A     application.Application
	Saved map[string]struct{}
	srv   *TpsServer
	wg    sync.WaitGroup
}

func NewTpsReceiver(a application.Application) *tpsReceiverStruct {

	cer, err := tls.LoadX509KeyPair("tls/cert.pem", "tls/key.pem")
	if err != nil {
		logger.L.Errorf("in tps.NewTpsReceiver tls.LoadX509KeyPair returned err: %v\n", err)
		return nil
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	li, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		logger.L.Errorf("in driver.Run error: %v\n", err)
	}
	logger.L.Infoln("listening localhost:443")

	srv := &TpsServer{l: li}

	return &tpsReceiverStruct{
		A:   a,
		srv: srv,
	}
}

func (r *tpsReceiverStruct) Run() {

	for {
		conn, err := r.srv.l.Accept()
		if err != nil {

			if r.A.Stopping() {
				return
			}

			logger.L.Error(err)
		}
		r.wg.Add(1)

		ts := repo.NewTS()
		go r.HandleRequestFull(conn, ts, &r.wg)

	}

}

// Tested in https_test.go
func (r *tpsReceiverStruct) HandleRequestFull(conn net.Conn, ts string, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(r.Saved) == 0 {
		bou, header, errFirst := repo.AnalyzeHeader(conn)
		//logger.L.Infof("in HandleReauest bou = %q, header = %q, errFirst = %v\n", bou, header, errFirst)
		if errFirst != nil && strings.Contains(errFirst.Error(), "100-continue") {
			r.Saved[string(repo.GenBoundary(bou)[2:])] = struct{}{}
			go repo.RespondContinue(conn)
			return
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
func (r *tpsReceiverStruct) HandleRequestLast(conn net.Conn, ts string, bou repo.Boundary, header []byte, errFirst error) {
	p := 0
	for {
		h := repo.NewReceiverHeader(ts, p, bou)
		b, errSecond := repo.AnalyzeBits(conn, 1024, p, header, errFirst)
		//logger.L.Infof("in HandleReauest b = %q, errSecond = %v\n", b, errSecond)

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
	go repo.Respond(conn)
}

func (r *tpsReceiverStruct) CleanSaved(s string) {
	if len(r.Saved) > 1 {
		delete(r.Saved, s)
	}
	r.Saved = make(map[string]struct{})
}

func (r *tpsReceiverStruct) Stop(wg *sync.WaitGroup) {

	r.srv.l.Close()

	r.wg.Wait()

	wg.Done()
}
