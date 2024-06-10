package transmitters

import (
	"context"
	"net"
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

	//kafkaAddr := os.Getenv("KAFKA_ADDR")
	kafkaAddr := "localhost"
	//kafkaPort := os.Getenv("KAFKA_PORT")
	kafkaPort := "29092"

	dialURI := net.JoinHostPort(kafkaAddr, kafkaPort)

	//kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaTopic := "topic1"

	logger.L.Infof("in transmitters.NewTransmitter dialURI: %s, topic: %s\n", dialURI, kafkaTopic)

	var (
		conn       *kafka.Conn
		err        error
		partitions []kafka.Partition
	)

	for i := 0; i < 5; i++ {

		logger.L.Infof("in transmitters.NewTransmitter %d attempt to connect to %s %s", i, "tcp", dialURI)

		conn, err = kafka.Dial("tcp", dialURI)
		if err != nil {

			logger.L.Errorf("in rpc.NewReceiver cannot dial: %v. Trying again\n", err)

			time.Sleep(time.Second * 5)

			continue
		}

		partitions, err = conn.ReadPartitions()

		if err != nil {

			logger.L.Errorf("in transmitters.NewTransmitter error reading partitions %v. Trying again", err)

			time.Sleep(time.Second * 5)

			continue

		}

		for _, p := range partitions {

			if p.Topic == kafkaTopic {

				logger.L.Infof("in rpc.NewReceiver topic %s is found\n", kafkaTopic)

				wc := kafka.WriterConfig{
					Brokers:  []string{dialURI},
					Topic:    kafkaTopic,
					Balancer: &kafka.RoundRobin{},
				}

				logger.L.Infof("in rpc.NewReceiver writerConfig %v\n", wc)

				kw := kafka.NewWriter(wc)

				logger.L.Infof("in rpc.NewReceiver kafkaWriter %v\n", kw)

				ts := &transmittersStruct{

					saverKafkaWriter: kw,
					encoder:          enc,
				}

				logger.L.Infof("in rpc.NewReceiver ts: %v:%v\n", ts.saverKafkaWriter, ts.encoder)

				return ts
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
