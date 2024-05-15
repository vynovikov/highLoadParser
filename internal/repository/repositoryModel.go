package repository

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
	Header() string
}

type TransferUnit struct {
}

type Presence struct {
}
