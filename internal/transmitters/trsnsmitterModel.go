package transmitters

type TransferUnit interface {
	Key() []byte
	Value() []byte
}

type transferUnitStruct struct {
	key   []byte
	value []byte
}

func NewTransferUnitStruct(t TransferUnit) *transferUnitStruct {

	return &transferUnitStruct{
		key:   t.Key(),
		value: t.Value(),
	}
}

func (t *transferUnitStruct) Key() []byte {

	return t.key
}

func (t *transferUnitStruct) Value() []byte {

	return t.value
}
