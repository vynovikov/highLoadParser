package transmitters

type TransferUnit interface {
	Key() []byte
	Value() []byte
}

type TransferUnitStruct struct {
	key   []byte
	value []byte
}

func NewTransferUnitStruct(t TransferUnit) *TransferUnitStruct {

	return &TransferUnitStruct{
		key:   t.Key(),
		value: t.Value(),
	}
}

func (t *TransferUnitStruct) Key() []byte {

	return t.key
}

func (t *TransferUnitStruct) Value() []byte {

	return t.value
}
