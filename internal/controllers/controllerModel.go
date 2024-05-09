package controllers

import "github.com/vynovikov/highLoadParser/internal/service"

const (
	boundaryField  = "boundary="
	sep            = "\r\n"
	maxLineLimit   = 100
	maxHeaderLimit = 210
)

type boundary struct {
	prefix []byte
	root   []byte
	suffix []byte
}

func getBoundary(bou boundary) []byte {

	boundary := make([]byte, 0)
	boundary = append(boundary, []byte("\r\n")...)
	boundary = append(boundary, bou.prefix...)
	boundary = append(boundary, bou.root...)

	return boundary
}

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
	bou  boundary
}

func newParserControllerHeader(ts string, p int, bou boundary) parserControllerHeader {

	return parserControllerHeader{
		part: p,
		ts:   ts,
		bou:  bou,
	}
}

type parserServiceInitDTO struct {
	part       int
	ts         string
	body       []byte
	startIndex int
	bou        boundary
	last       bool
	psus       []*service.ParserServiceUnit
	pssu       *service.ParserServiceSub
}

func newParserServiceInitDTO(h parserControllerHeader, b parserControllerBody) *parserServiceInitDTO {
	return &parserServiceInitDTO{
		body: b.B,
		bou:  h.bou,
		part: h.part,
		ts:   h.ts,
	}
}
