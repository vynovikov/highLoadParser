package dataHandler

/*
func NewValue() value {

		return value{}
	}
*/
type Disposition int

const (
	False Disposition = iota
	True
	Probably
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

func newKey(d DataHandlerDTO) key {

	return key{
		TS:   d.TS(),
		Part: d.Part(),
	}
}

type value struct {
}

func newValue(d DataHandlerDTO) value {

	return value{}
}

/*
	type value struct {
		D headerData
		B BeginningData
		E probability
	}
*/
type headerData struct {
	FormName string
	FileName string
	H        []byte
}

type BeginningData struct {
	Part int
}

type Presence struct {
	value map[bool]value
}
