package encoder

import (
	"encoding/json"

	"github.com/vynovikov/highLoadParser/internal/logger"
)

type jsonEncoderStruct struct {
}

func NewJSONEncoder() jsonEncoderStruct {

	return jsonEncoderStruct{}
}

func (j jsonEncoderStruct) EncodeKey(u TransferUnit) []byte {

	keyMap := make(map[string]interface{}, 0)

	keyMap["ts"] = u.TS()
	keyMap["formName"] = u.FormName()
	keyMap["fileName"] = u.FileName()
	keyMap["first"] = u.Start()

	marshalledKey, err := json.Marshal(keyMap)
	if err != nil {

		logger.L.Println(err)
	}

	return marshalledKey
}

func (j jsonEncoderStruct) EncodeValue(u TransferUnit) []byte {

	valueMap := make(map[string]interface{}, 0)

	valueMap["body"] = string(u.Body())
	valueMap["last"] = u.Final()

	marshalledValue, err := json.Marshal(valueMap)
	if err != nil {

		logger.L.Println(err)
	}

	return marshalledValue
}
