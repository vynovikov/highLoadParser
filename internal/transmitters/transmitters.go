package transmitters

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/vynovikov/highLoadParser/internal/logger"
	"github.com/vynovikov/highLoadParser/internal/service/pb"
	"google.golang.org/protobuf/proto"
)

type ParserTransmitter interface {
	TransmitToSaver(TransferUnit) error
	TransmitToLogger(TransferUnit) error
}

type transmittersStruct struct {
	saverKafkaWriter *kafka.Writer
}

func NewTransmitter() *transmittersStruct {

	return &transmittersStruct{

		saverKafkaWriter: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{"localhost:29092"},
			Topic:    "topic1",
			Balancer: &kafka.RoundRobin{},
		}),
	}
}

func (t *transmittersStruct) TransmitToLogger(TransferUnit) error {

	return nil
}

// TODO: Create encoding dependency
func (t *transmittersStruct) TransmitToSaver(unit TransferUnit) error {

	protoKey := &pb.MessageHeader{
		Ts:       unit.TS(),
		FormName: unit.FormName(),
		FileName: unit.FileName(),
		First:    unit.Start(),
	}

	protoValue := &pb.MessageBody{
		Body: unit.Body(),
		Last: unit.End(),
	}

	marshalledKey, err := proto.Marshal(protoKey)
	if err != nil {

		logger.L.Warn(err)
	}

	marshalledValue, err := proto.Marshal(protoValue)
	if err != nil {

		logger.L.Warn(err)
	}

	t.saverKafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   marshalledKey,
		Value: marshalledValue,
	})

	return nil
}

func NewParserTransmitters() *transmittersStruct {

	return &transmittersStruct{}
}
