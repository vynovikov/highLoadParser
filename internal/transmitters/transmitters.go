package transmitters

type ParserTransmitter interface {
	TransmitToSaver() error
	TransmitToLogger() error
}

type transmittersStruct struct {
}

func NewTransmitter() *transmittersStruct {

	return &transmittersStruct{}
}

func (t *transmittersStruct) TransmitToLogger() error {

	return nil
}

func (t *transmittersStruct) TransmitToSaver() error {

	return nil
}

func NewParserTransmitters() *transmittersStruct {

	return &transmittersStruct{}
}
