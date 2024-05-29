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
	CONTENT_DISPOSITION             = "Content-Disposition"
	False               Disposition = iota
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
	SetBody([]byte)
	B() int
	E() int
	Last() bool
	IsSub() bool
}

type DataHandlerUnit struct {
	Dh_part  int
	Dh_ts    string
	Dh_body  []byte
	Dh_b     int
	Dh_e     int
	Dh_last  bool
	Dh_isSub bool
}

func NewDataHandlerUnit(d DataHandlerDTO) *DataHandlerUnit {
	return &DataHandlerUnit{
		Dh_part: d.Part(),
		Dh_ts:   d.TS(),
		Dh_body: d.Body(),
		Dh_b:    d.B(),
		Dh_e:    d.E(),
		Dh_last: d.Last(),
	}
}

func (d *DataHandlerUnit) Part() int {
	return d.Dh_part
}

func (d *DataHandlerUnit) TS() string {
	return d.Dh_ts
}

func (d *DataHandlerUnit) Body() []byte {
	return d.Dh_body
}

func (d *DataHandlerUnit) SetBody(b []byte) {
	d.Dh_body = b
}

func (d *DataHandlerUnit) B() int {
	return d.Dh_b
}

func (d *DataHandlerUnit) E() int {
	return d.Dh_e
}

func (d *DataHandlerUnit) Last() bool {
	return d.Dh_last
}

func (d *DataHandlerUnit) IsSub() bool {
	return d.Dh_isSub
}

type keyGeneral struct {
	ts string
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

type value struct {
	h headerData
	e int
}

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

type ProducerUnit interface {
	TS() string
	Part() int
	FormName() string
	FileName() string
	Body() []byte
	Start() bool
	IsSub() bool
	End() bool
	Final() bool
}

type ProducerUnitStruct struct {
	Dh_TS       string
	Dh_Part     int
	Dh_FormName string
	Dh_FileName string
	Dh_Body     []byte
	Dh_Start    bool
	Dh_End      bool
	Dh_Final    bool
	Dh_IsSub    bool
}

func (t *ProducerUnitStruct) TS() string {

	return t.Dh_TS
}
func (t *ProducerUnitStruct) Part() int {

	return t.Dh_Part
}

func (t *ProducerUnitStruct) FormName() string {

	return t.Dh_FormName
}

func (t *ProducerUnitStruct) FileName() string {

	return t.Dh_FileName
}

func (t *ProducerUnitStruct) Body() []byte {

	return t.Dh_Body
}

func (t *ProducerUnitStruct) Start() bool {

	return t.Dh_Start
}

func (t *ProducerUnitStruct) End() bool {

	return t.Dh_End
}

func (t *ProducerUnitStruct) Final() bool {

	return t.Dh_Final
}

func (t *ProducerUnitStruct) IsSub() bool {

	return t.Dh_IsSub
}
