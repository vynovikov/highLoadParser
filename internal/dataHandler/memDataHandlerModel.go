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

type DataHandlerDTO interface {
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

func NewDataHandlerUnit(d DataHandlerDTO) *DataHandlerUnit {
	return &DataHandlerUnit{
		part: d.Part(),
		ts:   d.TS(),
		body: d.Body(),
		b:    d.B(),
		e:    d.E(),
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
