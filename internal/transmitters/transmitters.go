package transmitters

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/vynovikov/highLoadParser/internal/encoder"
)

type ParserTransmitter interface {
	TransmitToSaver(TransferUnit) error
	TransmitToLogger(TransferUnit) error
}

type transmittersStruct struct {
	saverKafkaWriter *kafka.Writer
	encoder          encoder.Encoder
}

func NewTransmitter(enc encoder.Encoder) *transmittersStruct {

	return &transmittersStruct{

		saverKafkaWriter: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{"localhost:29092"},
			Topic:    "topic1",
			Balancer: &kafka.RoundRobin{},
		}),
		encoder: enc,
	}
}

func (t *transmittersStruct) TransmitToLogger(TransferUnit) error {

	return nil
}

func (t *transmittersStruct) TransmitToSaver(unit TransferUnit) error {

	encodedKey := t.encoder.EncodeKey(unit)
	encodedValue := t.encoder.EncodeValue(unit)

	t.saverKafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   encodedKey,
		Value: encodedValue,
	})

	return nil
}

func NewParserTransmitters() *transmittersStruct {

	return &transmittersStruct{}
}
