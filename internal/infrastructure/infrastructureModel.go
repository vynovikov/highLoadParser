package infrastructure

type TransferUnit interface {
	Key() []byte
	Value() []byte
}

type TransferUnitStruct struct {
	I_key   []byte
	I_value []byte
}

func (t *TransferUnitStruct) Key() []byte {

	return t.I_key
}

func (t *TransferUnitStruct) Value() []byte {

	return t.I_value
}
