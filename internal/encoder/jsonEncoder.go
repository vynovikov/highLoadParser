package encoder

type jsonEncoderStruct struct {
}

func NewJSONEncoder() jsonEncoderStruct {

	return jsonEncoderStruct{}
}

func (j jsonEncoderStruct) Encode() []byte {

	return make([]byte, 0)
}
