package transmitters

type ParserTransmitter interface {
	TransmitToSaver(TransferUnit) error
	TransmitToLogger(TransferUnit) error
}

type transmittersStruct struct {
}

func NewTransmitter() *transmittersStruct {

	return &transmittersStruct{}
}

func (t *transmittersStruct) TransmitToLogger(TransferUnit) error {

	return nil
}

func (t *transmittersStruct) TransmitToSaver(TransferUnit) error {

	return nil
}

func NewParserTransmitters() *transmittersStruct {

	return &transmittersStruct{}
}
