// Transmitter.
// Uses kafka.
// Transmits data to saver service and logs to logger service.
package rpc

import (
	"context"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/vynovikov/highLoadParser/internal/adapters/driven/rpc/tosaver/pb"
	"github.com/vynovikov/highLoadParser/internal/logger"
	"google.golang.org/protobuf/proto"

	"github.com/vynovikov/highLoadParser/internal/repo"
)

type Transmitter interface {
	Transmit(repo.AppDistributorUnit)
	Log(string) error
	Stop() error
}

type TransmitAdapter struct {
	lock  sync.Mutex
	KC    *kafka.Conn
	Topic string
}

func NewTransmitter() *TransmitAdapter {
	var (
		conn *kafka.Conn
		err  error
	)
	kafkaAddr := os.Getenv("KAFKA_ADDR")
	kafkaTopicName := os.Getenv("KAFKA_TOPIC_NAME")
	logger.L.Infof("in rpc.NewTransmitter kafkaTopicName: %s\n", kafkaTopicName)
	hostName := os.Getenv("HOSTNAME")
	logger.L.Infof("in rpc.NewTransmitter hostName: %s\n", hostName)

	//kafkaPartitionString := os.Getenv("KAFKA_PARTITION")
	/*
		partition, err := strconv.Atoi(kafkaPartitionString)
		if err != nil {
			logger.L.Errorf("in rpc.GetKafkaProducer unble to convers %q into int %v\n", kafkaPartitionString, err)
		}
	*/
	//logger.L.Infof("addr = %s, kafkaTopicName = %s, partition = %d\n", kafkaAddr, kafkaTopicName, partition)
	index := hostName[strings.LastIndex(hostName, "-")+1:]
	logger.L.Infof("in rpc.NewTransmitter index: %s\n", index)
	//kafkaTopic:=kafkaTopicName+hostName[strings.Index(hostName,strings.LastIndex(hostName,"-"))]

	for {
		conn, err = kafka.DialLeader(context.Background(), "tcp", kafkaAddr, kafkaTopicName+"-"+index, 0)
		if err != nil {
			//logger.L.Errorf("in rpc.GetKafkaProducer error %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	return &TransmitAdapter{
		KC:    conn,
		Topic: kafkaTopicName,
	}
}

func CreateTopic(conn *kafka.Conn, t string) error {
	controller, err := conn.Controller()
	if err != nil {
		logger.L.Errorf("in rpc.CreateTopic error %v\n", err)
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		logger.L.Errorf("in rpc.CreateTopic error %v\n", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             t,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		logger.L.Errorf("in rpc.CreateTopic error %v\n", err)
	}
	//logger.L.Infof("in rpc.GetKafkaProducer kafkaTopicName %q created\n", t)
	return nil
}

func (t *TransmitAdapter) Transmit(adu repo.AppDistributorUnit) {
	//logger.L.Infof("in rpc.Transmit transmitting adu header %v body %q\n", adu.GetHeader(), adu.GetBody())
	var (
		m    kafka.Message
		err  error
		newT *TransmitAdapter
	)
	m, err = GenMessage(adu, t.Topic)
	if err != nil {
		logger.L.Errorf("in rpc.Transmit generating message error %v\n", err)
	}
	//logger.L.Infof("in rpc.Transmit for adu header %v body %q made m %q\n", adu.GetHeader(), adu.GetBody(), m)
	for {
		if newT != nil {
			_, err = newT.KC.WriteMessages(m)
		} else {
			_, err = t.KC.WriteMessages(m)
		}
		if err != nil {
			logger.L.Errorf("in rpc.Transmit writing message error %v\n", err)
			if strings.Contains(err.Error(), "broken pipe") {
				logger.L.Infof("transmitter restarting")
				newT = NewTransmitter()

			}
			time.Sleep(time.Second * 5)
			continue
		}
		logger.L.Infof("in rpc.Transmit message %q was sent\n", m)
		break
	}
}
func (t *TransmitAdapter) Stop() error {
	return t.KC.Close()
}
func (t *TransmitAdapter) Log(s string) error {

	return nil
}
func GenMessage(adu repo.AppDistributorUnit, topic string) (kafka.Message, error) {
	var m kafka.Message
	//logger.L.Infof("in rpc.GenMessage adu header %v body %q\n", adu.GetHeader(), adu.GetBody())
	serialized, err := serialize(adu)
	if err != nil {
		return m, err
	}
	m.Topic = topic // exceptional
	m.Key = []byte(adu.H.TS)
	m.Value = serialized
	//logger.L.Infof("in rpc.GenMessage m = %v\n", m)

	return m, nil
}

func serialize(adu repo.AppDistributorUnit) ([]byte, error) {

	pbMessage := &pb.Message{
		Ts:         adu.H.TS,
		FormName:   adu.H.FormName,
		FileName:   adu.H.FileName,
		FieldValue: adu.B.B,
		Last:       adu.H.IsLast,
	}
	return proto.Marshal(pbMessage)

}
