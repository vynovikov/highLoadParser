package controllers

import (
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/vynovikov/highLoadParser/internal/entities"
	"github.com/vynovikov/highLoadParser/internal/logger"
	"github.com/vynovikov/highLoadParser/internal/service"
)

type ParserController interface {
	HandleRequestFull(net.Conn, string, *sync.WaitGroup)
}

type controllerStruct struct {
	srv     service.ParcerService
	ctrlMap map[string]struct{}
}

func NewController(s service.ParcerService) *controllerStruct {
	cm := make(map[string]struct{})

	return &controllerStruct{
		srv:     s,
		ctrlMap: cm,
	}
}

// Tested in http_test.go
func (c *controllerStruct) HandleRequestFull(conn net.Conn, ts string, wg *sync.WaitGroup) {

	logger.L.Infof("in controller.HandleRequestFull new request\n")

	defer wg.Done()

	bou, header, errFirst := analyzeHeader(conn)

	if errFirst != nil {

		if strings.Contains(errFirst.Error(), "100-continue") {

			c.ctrlMap[string(entities.GetBoundary(bou)[2:])] = struct{}{}

			respondContinue(conn)

			header = append(entities.GetBoundary(bou)[2:], []byte("\r\n")...)

		}
	}

	c.HandleRequestLast(conn, ts, bou, header, errFirst)
}

func (c *controllerStruct) HandleRequestLast(conn net.Conn, ts string, bou entities.Boundary, header []byte, errFirst error) {
	p := 0
	for {
		h := newParserControllerHeader(ts, p, bou)
		b, errSecond := analyzeBits(conn, 1024, p, header, errFirst)
		serviceDTO := newParserServiceInitDTO(h, b)

		if errFirst != nil {

			if errFirst == io.EOF || errFirst == io.ErrUnexpectedEOF || os.IsTimeout(errFirst) {
				c.srv.Serve(serviceDTO)
				break
			}
		}
		if errSecond != nil {
			if errSecond == io.EOF || errSecond == io.ErrUnexpectedEOF || os.IsTimeout(errSecond) {
				c.srv.Serve(serviceDTO)
				break
			}
			if strings.Contains(errSecond.Error(), "empty") {
				break
			}
		}
		c.srv.Serve(serviceDTO)

		p++
	}
	respondOK(conn)

}
