package dataHandler

import "errors"

type (
	Disposition int
	sufficiency int
)

var (
	errHeaderEnding   error = errors.New("ending of header")
	errHeaderNotFull  error = errors.New("header is not full")
	errHeaderNotFound error = errors.New("header is not found")
)

const (
	False Disposition = iota
	True
	Probably
	incomplete sufficiency = iota
	sufficient
	insufficient
	sep            = "\r\n"
	maxLineLimit   = 100
	maxHeaderLimit = 210
)

type DataHandlerDTO interface {
	Part() int
	TS() string
	Body() []byte
	B() Disposition
	E() Disposition
	Last() bool
	IsSub() bool
}

type DataHandlerUnit struct {
	part  int
	ts    string
	body  []byte
	b     Disposition
	e     Disposition
	last  bool
	isSub bool
}

func NewDataHandlerUnit(d DataHandlerDTO) *DataHandlerUnit {
	return &DataHandlerUnit{
		part: d.Part(),
		ts:   d.TS(),
		body: d.Body(),
		b:    d.B(),
		e:    d.E(),
		last: d.Last(),
	}
}

func (d *DataHandlerUnit) Part() int {
	return d.part
}

func (d *DataHandlerUnit) TS() string {
	return d.ts
}

func (d *DataHandlerUnit) Body() []byte {
	return d.body
}

func (d *DataHandlerUnit) B() Disposition {
	return d.b
}

func (d *DataHandlerUnit) E() Disposition {
	return d.e
}

func (d *DataHandlerUnit) Last() bool {
	return d.last
}

func (d *DataHandlerUnit) IsSub() bool {
	return d.isSub
}

type key struct {
	TS   string
	Part int
}

type keyGeneral struct {
	ts string
}

func newKeyGeneral(d DataHandlerDTO) keyGeneral {

	return keyGeneral{
		ts: d.TS(),
	}
}

type keyDetailed struct {
	ts   string
	part int
}

func newKeyDetailed(d DataHandlerDTO) keyDetailed {

	return keyDetailed{
		ts:   d.TS(),
		part: d.Part(),
	}
}

func newKey(d DataHandlerDTO) key {

	return key{
		TS:   d.TS(),
		Part: d.Part(),
	}
}

type value struct {
	h headerData
	e Disposition
}

/*
	type value struct {
		D headerData
		B BeginningData
		E probability
	}
*/
type headerData struct {
	formName    string
	fileName    string
	headerBytes []byte
}

type BeginningData struct {
	Part int
}

type Presence struct {
	value map[bool]value
}

type Boundary struct {
	Prefix []byte
	Root   []byte
	Suffix []byte
}
