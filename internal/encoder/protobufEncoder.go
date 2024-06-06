package encoder

type protobufEncoderStruct struct {
}

func NewProtobufEncoder() protobufEncoderStruct {

	return protobufEncoderStruct{}
}

func (p protobufEncoderStruct) Encode() []byte {

	return make([]byte, 0)
}
