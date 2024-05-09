package controllers

import (
	"bytes"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/vynovikov/highLoadParser/internal/logger"
	"github.com/vynovikov/highLoadParser/internal/service"
	"github.com/vynovikov/highLoadParser/pkg/byteOps"
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

			c.ctrlMap[string(getBoundary(bou)[2:])] = struct{}{}

			respondContinue(conn)

			header = append(getBoundary(bou)[2:], []byte("\r\n")...)

		}
	}

	c.HandleRequestLast(conn, ts, bou, header, errFirst)
}

func (c *controllerStruct) HandleRequestLast(conn net.Conn, ts string, bou boundary, header []byte, errFirst error) {
	p := 0
	for {
		h := newParserControllerHeader(ts, p, bou)
		b, errSecond := analyzeBits(conn, 1024, p, header, errFirst)
		serviceDTO := newParserServiceDTO(h, b)

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

func (c *controllerStruct) CleanMap(s string) {
	if len(c.ctrlMap) > 1 {
		delete(c.ctrlMap, s)
	}
	c.ctrlMap = make(map[string]struct{})
}

func newParserServiceDTO(h parserControllerHeader, b parserControllerBody) service.ParserServiceDTO {

	dto := newParserServiceInitDTO(h, b)
	dto.Evolve(0)

	return service.ParserServiceDTO{
		U: dto.psus,
		S: dto.pssu,
	}
}

func (s *parserServiceInitDTO) Evolve(start int) {

	b := s.body[start:]
	boundaryCore := getBoundary(s.bou)[2:]
	ll := make([]byte, 0, maxLineLimit)

	if len(b) < len(boundaryCore)+2*len(sep) {

		s.last = true

		psu := service.NewParserServiceUnit(service.NewParserServiceHeader(s.ts, s.part, service.True, service.False), service.NewParserServiceBody(b))

		s.psus = append(s.psus, &psu)

		return
	}

	if start == 0 {
		be := make([]byte, 0)

		if len(b) > maxHeaderLimit {
			be = b[len(b)-maxLineLimit:]

		} else {
			be = b
		}

		lenbe := len(be)
		ll = GetLineWithCRLFLeft(be, lenbe-1, maxLineLimit, s.bou)

		if (len(ll) > 2 && byteOps.BeginningEqual(ll[2:], boundaryCore)) ||
			(len(ll) <= 2 && bytes.Contains([]byte(sep), ll)) { // last line equal to boundary begginning or vice versa

			if IsLastBoundary(ll, []byte(""), s.bou) { // last boundary in last line

				s.last = true

			} else {

				as := service.NewParserServiceSub(service.NewParserServiceSubHeader(s.ts, s.part), service.NewParserServiceSubBody(ll))
				s.pssu = &as
			}

			b = b[:len(b)-len(ll)]

			s.body = s.body[:len(s.body)-len(ll)]

			if bytes.Contains(b, []byte(boundaryField)) {

				bodyIdx := bytes.Index(b, boundaryCore)

				b = b[bodyIdx+len(sep)+len(boundaryCore):]

				start += bodyIdx + len(sep) + len(boundaryCore)
			}
		}
	}
	switch bytes.Count(b, boundaryCore) {
	case 0:
		psu := service.ParserServiceUnit{}

		if (s.pssu == nil && len(ll) > 2 && !byteOps.BeginningEqual(ll[2:], boundaryCore)) ||
			(s.pssu == nil && len(ll) <= 2 && !bytes.Contains([]byte(sep), ll)) ||
			(s.pssu == nil && s.last) {

			if s.last && start == 0 {
				psu = service.NewParserServiceUnit(service.NewParserServiceHeader(s.ts, s.part, service.True, service.False), service.NewParserServiceBody(b))

			} else if s.last && start != 0 {
				psu = service.NewParserServiceUnit(service.NewParserServiceHeader(s.ts, s.part, service.False, service.False), service.NewParserServiceBody(b))

			} else if !s.last && start == 0 {
				psu = service.NewParserServiceUnit(service.NewParserServiceHeader(s.ts, s.part, service.True, service.True), service.NewParserServiceBody(b))

			} else if !s.last && start != 0 {
				psu = service.NewParserServiceUnit(service.NewParserServiceHeader(s.ts, s.part, service.False, service.True), service.NewParserServiceBody(b))

			}

		} else if (s.pssu == nil && len(ll) > 2 && byteOps.BeginningEqual(ll[2:], boundaryCore)) ||
			(s.pssu == nil && len(ll) <= 2 && bytes.Contains([]byte(sep), ll)) {

			if start == 0 {
				psu = service.NewParserServiceUnit(service.NewParserServiceHeader(s.ts, s.part, service.True, service.False), service.NewParserServiceBody(b))

			} else {
				psu = service.NewParserServiceUnit(service.NewParserServiceHeader(s.ts, s.part, service.False, service.True), service.NewParserServiceBody(b))

			}

		} else {

			if start == 0 {
				psu = service.NewParserServiceUnit(service.NewParserServiceHeader(s.ts, s.part, service.True, service.Probably), service.NewParserServiceBody(b))

			} else {
				psu = service.NewParserServiceUnit(service.NewParserServiceHeader(s.ts, s.part, service.False, service.Probably), service.NewParserServiceBody(b))

			}
		}

		s.psus = append(s.psus, &psu)

	default:

		idx := bytes.Index(b, boundaryCore) + start

		psu := service.ParserServiceUnit{}

		if start == 0 {
			psu = service.NewParserServiceUnit(service.NewParserServiceHeader(s.ts, s.part, service.True, service.False), service.NewParserServiceBody(b[:idx-start-len(sep)]))

		} else {
			psu = service.NewParserServiceUnit(service.NewParserServiceHeader(s.ts, s.part, service.False, service.False), service.NewParserServiceBody(b[:idx-start-len(sep)]))
		}

		s.psus = append(s.psus, &psu)

		s.Evolve(idx + len(sep) + len(boundaryCore))

	}
}
