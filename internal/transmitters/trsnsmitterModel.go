package transmitters

import (
	"github.com/segmentio/kafka-go"
	"github.com/vynovikov/highLoadParser/internal/encoder"
)

type TransferUnit interface {
	TS() string
	Part() int
	FormName() string
	FileName() string
	Body() []byte
	Start() bool
	IsSub() bool
	End() bool
	Final() bool
}

type transmittersStruct struct {
	saverKafkaWriter *kafka.Writer
	encoder          encoder.Encoder
}
type ProducerUnit interface {
}
