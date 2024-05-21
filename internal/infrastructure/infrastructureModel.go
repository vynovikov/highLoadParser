package infrastructure

type TransferUnit interface {
	Key() []byte
	Value() []byte
}

type TransferUnitStruct struct {
}

func (t *TransferUnitStruct) Key() []byte {

	return nil
}

func (t *TransferUnitStruct) Value() []byte {

	return nil
}
