package infrastructure

type DataPiece interface {
	Part() int
	TS() string
	Body() []byte
	Header() string
}

type TransferUnit struct {
	TH TransferHeader
	TB TransferBody
}

type TransferHeader struct {
}

type TransferBody struct {
	B []byte
}
