package infrastructure

type TransferUnit interface {
	Key() []byte
	Value() []byte
}

type transferUnitStruct struct {
	key   []byte
	value []byte
}

/*
func NewTransferUnit(d InfrastructureDTO) *transferUnitStruct {

	return &transferUnitStruct{
		key:   []byte("alice"),
		value: []byte("azaza"),
	}
}
*/

func (t *transferUnitStruct) Key() []byte {

	return t.key
}

func (t *transferUnitStruct) Value() []byte {

	return t.value
}
