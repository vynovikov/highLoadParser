package transmitters

type TransferUnit interface {
	TS() string
	Part() int
	FormName() string
	FileName() string
	Body() []byte
	Start() bool
	IsSub() bool
	End() bool
	Final() bool
}

/*
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
*/
type ProducerUnit interface {
}
