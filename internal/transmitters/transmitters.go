package transmitters

import (
	"context"

	"github.com/segmentio/kafka-go"
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

func (t *transmittersStruct) TransmitToSaver(unit TransferUnit) error {

	t.saverKafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   unit.Key(),
		Value: unit.Value(),
	})

	return nil
}

func NewParserTransmitters() *transmittersStruct {

	return &transmittersStruct{}
}
