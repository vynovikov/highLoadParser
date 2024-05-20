package service

import (
	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/entities"
)

type disposition int

const (
	False disposition = iota
	True
	Probably
	Sep            = "\r\n"
	MaxLineLimit   = 100
	MaxHeaderLimit = 210
)

type ParserServiceHeader struct {
	Part int
	TS   string
	B    disposition
	E    disposition
	last bool
}

func NewParserServiceHeader(ts string, p int, b, e disposition) ParserServiceHeader {
	return ParserServiceHeader{
		Part: p,
		TS:   ts,
		B:    b,
		E:    e,
	}
}

type ParserServiceBody struct {
	B []byte
}

func NewParserServiceSubBody(b []byte) ParserServiceSubBody {
	return ParserServiceSubBody{
		B: b,
	}
}

type ParserServiceUnit struct {
	PSH ParserServiceHeader
	PSB ParserServiceBody
}

func NewParserServiceUnit(psh ParserServiceHeader, psb ParserServiceBody) ParserServiceUnit {
	return ParserServiceUnit{
		PSH: psh,
		PSB: psb,
	}
}

type ParserServiceSubBody struct {
	B []byte
}

func NewParserServiceBody(b []byte) ParserServiceBody {
	return ParserServiceBody{
		B: b,
	}
}

type ParserServiceSubHeader struct {
	Part int
	TS   string
}

func NewParserServiceSubHeader(ts string, p int) ParserServiceSubHeader {
	return ParserServiceSubHeader{
		Part: p,
		TS:   ts,
	}
}

type ParserServiceSub struct {
	PSSH ParserServiceSubHeader
	PSSB ParserServiceSubBody
}

func NewParserServiceSub(pssh ParserServiceSubHeader, pssb ParserServiceSubBody) ParserServiceSub {
	return ParserServiceSub{
		PSSH: pssh,
		PSSB: pssb,
	}
}

type DataHandlerDTO interface {
	Part() int
	TS() string
	Body() []byte
	SetBody([]byte)
	B() dataHandler.Disposition
	E() dataHandler.Disposition
	Last() bool
	IsSub() bool
}

// ParserServiceUnit -> dataPiece

func (su *ParserServiceUnit) Part() int {
	return su.PSH.Part
}

func (su *ParserServiceUnit) TS() string {
	return su.PSH.TS
}

func (su *ParserServiceUnit) B() dataHandler.Disposition {
	return dataHandler.Disposition(su.PSH.B)
}

func (su *ParserServiceUnit) E() dataHandler.Disposition {
	return dataHandler.Disposition(su.PSH.E)
}

func (su *ParserServiceUnit) Body() []byte {
	return su.PSB.B
}

func (su *ParserServiceUnit) SetBody(b []byte) {
	su.PSB.B = b
}

func (su *ParserServiceUnit) Last() bool {
	return su.PSH.last
}

func (su *ParserServiceUnit) IsSub() bool {
	return false
}

// ParserServiceSub -> dataPiece

func (ss *ParserServiceSub) Part() int {
	return ss.PSSH.Part
}

func (ss *ParserServiceSub) TS() string {
	return ss.PSSH.TS
}

func (ss *ParserServiceSub) B() dataHandler.Disposition {
	return dataHandler.Disposition(False)
}

func (ss *ParserServiceSub) E() dataHandler.Disposition {
	return dataHandler.Disposition(Probably)
}

func (ss *ParserServiceSub) Body() []byte {
	return ss.PSSB.B
}

func (ss *ParserServiceSub) SetBody(b []byte) {
	ss.PSSB.B = b
}

func (ss *ParserServiceSub) Last() bool {
	return false
}

func (ss *ParserServiceSub) IsSub() bool {
	return true
}

type ParserServiceDTO struct {
	Part int
	TS   string
	Body []byte
	Bou  entities.Boundary
	last bool
	psus []*ParserServiceUnit
	pssu *ParserServiceSub
}
