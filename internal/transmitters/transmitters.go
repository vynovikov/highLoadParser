package transmitters

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/vynovikov/highLoadParser/internal/logger"
)

type ParserTransmitter interface {
	TransmitToSaver([]TransferUnit) error
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

func (t *transmittersStruct) TransmitToSaver(units []TransferUnit) error {

	for _, v := range units {

		logger.L.Infof("in transmitter.TransmitToParser sending unit key: %s, value: %s\n", string(v.Key()), string(v.Value()))

		t.saverKafkaWriter.WriteMessages(context.Background(), kafka.Message{
			Key:   v.Key(),
			Value: v.Value(),
		})
	}

	return nil
}

func NewParserTransmitters() *transmittersStruct {

	return &transmittersStruct{}
}
