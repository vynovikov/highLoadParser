package controllers

import (
	"github.com/vynovikov/highLoadParser/internal/entities"
	"github.com/vynovikov/highLoadParser/internal/service"
)

const (
	boundaryField  = "boundary="
	sep            = "\r\n"
	maxLineLimit   = 100
	maxHeaderLimit = 210
)

type parserControllerBody struct {
	B []byte
}

func newParserControllerBody(n int) parserControllerBody {
	return parserControllerBody{
		B: make([]byte, n),
	}
}

type parserControllerHeader struct {
	part int
	ts   string
	bou  entities.Boundary
}

func newParserControllerHeader(ts string, p int, bou entities.Boundary) parserControllerHeader {

	return parserControllerHeader{
		part: p,
		ts:   ts,
		bou:  bou,
	}
}

func newParserServiceInitDTO(h parserControllerHeader, b parserControllerBody) service.ParserServiceDTO {
	return service.ParserServiceDTO{
		Body: b.B,
		Bou:  h.bou,
		Part: h.part,
		TS:   h.ts,
	}
}
