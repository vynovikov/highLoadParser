package infrastructure

type DataPiece interface {
	Part() int
	TS() string
	Body() []byte
	Header() string
}

type TransferUnitStruct struct {
	TH TransferHeader
	TB TransferBody
}

type TransferHeader struct {
}

type TransferBody struct {
	B []byte
}

type TransferUnit interface {
	Tx() error
}

func (t *TransferUnitStruct) Tx() error {

	return nil
}

type Presence struct {
}
