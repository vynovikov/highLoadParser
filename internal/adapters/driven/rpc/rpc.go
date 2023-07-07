// Transmitter.
// Uses kafka.
// Transmits data to saver service and logs to logger service.
package rpc

import (
	"context"
	"net"
	"strconv"
	"sync"

	"github.com/segmentio/kafka-go"
	"github.com/vynovikov/highLoadParser/internal/adapters/driven/rpc/tosaver/pb"
	"github.com/vynovikov/highLoadParser/internal/logger"
	"google.golang.org/protobuf/proto"

	"github.com/vynovikov/highLoadParser/internal/repo"
)

type Transmitter interface {
	Transmit(repo.AppDistributorUnit)
	Log(string) error
}

type TransmitAdapter struct {
	lock sync.Mutex
	KW   *kafka.Writer
	t    string
}

func NewTransmitter(t string) *TransmitAdapter {

	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		logger.L.Errorf("in rpc.GetKafkaProducer error %v", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		logger.L.Errorf("in rpc.GetKafkaProducer error %v", err)
	}
	if len(partitions) == 0 {
		logger.L.Infof("in rpc.GetKafkaProducer no topics found\n")
	}
	//logger.L.Infoln("azaza")
	for _, v := range partitions {
		if t == v.Topic {
			//logger.L.Infof("in rpc.GetKafkaProducer topic %q found\n", t)
			return &TransmitAdapter{
				KW: NewWriter(t),
				t:  t,
			}
		}
	}
	//logger.L.Infof("in rpc.GetKafkaProducer topic %q is not found should be created\n", t)
	err = CreateTopic(conn, t)
	if err != nil {
		logger.L.Errorf("in rpc.GetKafkaProducer error %v", err)
	}
	return &TransmitAdapter{
		KW: NewWriter(t),
		t:  t,
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
	//logger.L.Infof("in rpc.GetKafkaProducer topic %q created\n", t)
	return nil
}

func (t *TransmitAdapter) Transmit(adu repo.AppDistributorUnit) {
	var (
		m   kafka.Message
		err error
	)
	m, err = GenMessage(adu, t.t)
	if err != nil {
		logger.L.Errorf("in rpc.Transmit error %v\n", err)
	}
	//logger.L.Infof("in rpc.Transmit for adu header %v body %q made m value %q\n", adu.GetHeader(), adu.GetBody(), m.Value)
	err = t.KW.WriteMessages(context.Background(), m)
	if err != nil {
		logger.L.Errorf("in rpc.Transmit error %v\n", err)
	}
}
func (t *TransmitAdapter) Log(s string) error {

	return nil
}
func NewWriter(t string) *kafka.Writer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   t,
	})
	//logger.L.Infof("in rpc.NewWriter w = %v\n", w)
	return w
}
func GenMessage(adu repo.AppDistributorUnit, t string) (kafka.Message, error) {
	var m kafka.Message

	serialized, err := serialize(adu)
	if err != nil {
		return m, err
	}
	//m.Topic = t
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
