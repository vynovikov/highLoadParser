package transmitters

import (
	"context"
	"fmt"
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

type transmittersStruct struct {
	saverKafkaWriter *kafka.Writer
	encoder          encoder.Encoder
}

func NewTransmitter(enc encoder.Encoder) *transmittersStruct {

	kafkaAddr := os.Getenv("KAFKA_ADDR")
	kafkaPort := os.Getenv("KAFKA_PORT")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	logger.L.Infof("in transmitters.NewTransmitter broker: %s:%s, topic: %s\n", kafkaAddr, kafkaPort, kafkaTopic)

	var (
		conn *kafka.Conn
		err  error
	)

	for i := 0; i < 5; i++ {

		logger.L.Infof("in transmitters.NewTransmitter %d attempt to connect to %s %s", i, "tcp", fmt.Sprintf("%s:%s", kafkaAddr, kafkaPort))

		conn, err = kafka.Dial("tcp", fmt.Sprintf("%s:%s", kafkaAddr, kafkaPort))
		if err != nil {
			logger.L.Errorf("in rpc.NewReceiver cannot dial: %v\n", err)
			time.Sleep(time.Second * 10)
			if conn != nil {
				conn.Close()
			}
		}
	}
	partitions, err := conn.ReadPartitions()
	if err != nil {

		logger.L.Errorln(err)
		os.Exit(1)
	}

	for _, p := range partitions {

		if p.Topic == kafkaTopic {

			logger.L.Infof("in rpc.NewReceiver topic %s is found\n", kafkaTopic)

			return &transmittersStruct{

				saverKafkaWriter: kafka.NewWriter(kafka.WriterConfig{
					Brokers:  []string{fmt.Sprintf("%s:%s", kafkaAddr, kafkaPort)},
					Topic:    kafkaTopic,
					Balancer: &kafka.RoundRobin{},
				}),
				encoder: enc,
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
