package encoder

import (
	"github.com/vynovikov/highLoadParser/internal/encoder/pb"
	"github.com/vynovikov/highLoadParser/internal/logger"
	"google.golang.org/protobuf/proto"
)

type protobufEncoderStruct struct {
}

func NewProtobufEncoder() protobufEncoderStruct {

	return protobufEncoderStruct{}
}

func (p protobufEncoderStruct) EncodeKey(u TransferUnit) []byte {

	protoKey := &pb.MessageHeader{
		Ts:       u.TS(),
		FormName: u.FormName(),
		FileName: u.FileName(),
		First:    u.Start(),
	}

	marshalledKey, err := proto.Marshal(protoKey)
	if err != nil {

		logger.L.Warn(err)
	}

	return marshalledKey
}

func (p protobufEncoderStruct) EncodeValue(u TransferUnit) []byte {

	protoValue := &pb.MessageBody{
		Body: u.Body(),
		Last: u.End(),
	}

	marshalledValue, err := proto.Marshal(protoValue)
	if err != nil {

		logger.L.Warn(err)
	}

	return marshalledValue
}
