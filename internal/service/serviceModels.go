package service

import (
	"fmt"

	"github.com/vynovikov/highLoadParser/internal/entities"
	"github.com/vynovikov/highLoadParser/internal/repository"
)

type Disposition int

const (
	False Disposition = iota
	True
	Probably
	Sep            = "\r\n"
	MaxLineLimit   = 100
	MaxHeaderLimit = 210
)

type Boundary struct {
	Prefix []byte
	Root   []byte
	Suffix []byte
}

/*
type ParserServiceDTO struct {
	U []*ParserServiceUnit
	S *ParserServiceSub
}
*/

type ParserServiceHeader struct {
	Part int
	TS   string
	B    Disposition
	E    Disposition
}

func NewParserServiceHeader(ts string, p int, b, e Disposition) ParserServiceHeader {
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

type DataPiece interface {
	Part() int
	TS() string
	Body() []byte
	Header() string
}

// ParserServiceUnit -> dataPiece

func (su *ParserServiceUnit) Part() int {
	return su.PSH.Part
}

func (su *ParserServiceUnit) TS() string {
	return su.PSH.TS
}

func (su *ParserServiceUnit) B() repository.Disposition {
	return repository.Disposition(su.PSH.B)
}

func (su *ParserServiceUnit) E() repository.Disposition {
	return repository.Disposition(su.PSH.E)
}

func (su *ParserServiceUnit) Body() []byte {
	return su.PSB.B
}

func (su *ParserServiceUnit) Header() string {
	return fmt.Sprintf("TS = %s, Part = %d, B() = %d, E() = %d\n", su.PSH.TS, su.PSH.Part, su.PSH.B, su.PSH.E)
}

// ParserServiceSub -> dataPiece

func (ss *ParserServiceSub) Part() int {
	return ss.PSSH.Part
}

func (ss *ParserServiceSub) TS() string {
	return ss.PSSH.TS
}

func (ss *ParserServiceSub) B() repository.Disposition {
	return repository.Disposition(False)
}

func (ss *ParserServiceSub) E() repository.Disposition {
	return repository.Disposition(Probably)
}

func (ss *ParserServiceSub) Body() []byte {
	return ss.PSSB.B
}

func (ss *ParserServiceSub) Header() string {
	return fmt.Sprintf("TS = %s, Part = %d, B() = %d, E() = %d\n", ss.PSSH.TS, ss.PSSH.Part, False, Probably)
}

// --------------------------------

type ParserServiceDTO struct {
	Part       int
	TS         string
	Body       []byte
	startIndex int
	Bou        entities.Boundary
	last       bool
	psus       []*ParserServiceUnit
	pssu       *ParserServiceSub
}
