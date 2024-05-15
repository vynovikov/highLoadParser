package dataHandler

type Presence struct {
}

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

type DataPiece interface {
	Part() int
	TS() string
	Body() []byte
	B() Disposition
	E() Disposition
}

type DataHandlerUnit struct {
	part int
	ts   string
	body []byte
	b    Disposition
	e    Disposition
}

func NewDataHandlerUnit(d DataPiece) *DataHandlerUnit {
	return &DataHandlerUnit{
		part: d.Part(),
		ts:   d.TS(),
		body: d.Body(),
		b:    Disposition(d.B()),
		e:    Disposition(d.E()),
	}
}

func (i *DataHandlerUnit) Part() int {
	return i.part
}

func (i *DataHandlerUnit) TS() string {
	return i.ts
}
func (i *DataHandlerUnit) Body() []byte {
	return i.body
}
func (i *DataHandlerUnit) B() Disposition {
	return i.b
}
func (i *DataHandlerUnit) E() Disposition {
	return i.e
}

type key struct {
	TS   string
	Part int
}

func newKey(d DataPiece) key {

	return key{
		TS:   d.TS(),
		Part: d.Part(),
	}
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
