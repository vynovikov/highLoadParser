package dataHandler

type Presence struct {
}

func NewValue() value {

	return value{}
}

type DataPiece interface {
	Part() int
	TS() string
	Body() []byte
	Header() string
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

type value struct {
	D headerData
	B BeginningData
	E probability
}

type headerData struct {
	FormName string
	FileName string
	H        []byte
}

type BeginningData struct {
	Part int
}

type probability int

const (
	False probability = iota
	True
	Probably
)
