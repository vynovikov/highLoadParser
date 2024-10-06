package transmitters

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/vynovikov/highLoadParser/internal/encoder"
	"github.com/vynovikov/highLoadParser/internal/logger"
)

type ParserTransmitter interface {
	TransmitToSaver(TransferUnit) error
	TransmitToLogger(TransferUnit) error
}

func NewTransmitter(enc encoder.Encoder) *transmittersStruct {

	var (
		conn       *kafka.Conn
		err        error
		partitions []kafka.Partition
	)

	kafkaAddr := os.Getenv("KAFKA_HOSTNAME")
	kafkaPort := os.Getenv("KAFKA_PORT")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	dialURI := net.JoinHostPort(kafkaAddr, kafkaPort)

	for i := 0; i < 5; i++ {

		conn, err = kafka.Dial("tcp", dialURI)
		if err != nil {

			logger.L.Errorf("in transmitters.NewTransmitter cannot dial: %v. Trying again\n", err)

			time.Sleep(time.Second * 10)

			continue
		}

		partitions, err = conn.ReadPartitions()

		if err != nil {

			logger.L.Errorf("in transmitters.NewTransmitter cannot read partitions: %v. Trying again\n", err)

			time.Sleep(time.Second * 10)

			continue

		}

		logger.L.Infof("in transmitters.NewTransmitter partitions: %v len %d\n", partitions, len(partitions))

		for _, p := range partitions {

			if p.Topic == kafkaTopic {

				ts := &transmittersStruct{

					saverKafkaWriter: kafka.NewWriter(kafka.WriterConfig{
						Brokers:  []string{dialURI},
						Topic:    kafkaTopic,
						Balancer: &kafka.RoundRobin{},
					}),
					encoder: enc,
				}

				return ts

			} else {

				time.Sleep(time.Second * 10)

				continue
			}
		}
	}

	return &transmittersStruct{}
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
