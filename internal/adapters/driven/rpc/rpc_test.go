package rpc

import (
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"
	"github.com/vynovikov/highLoadParser/internal/repo"
)

type rpcSuite struct {
	suite.Suite
}

func TestRpcSuite(t *testing.T) {
	suite.Run(t, new(rpcSuite))
}

func (s *rpcSuite) TestGenMessage() {
	tt := []struct {
		name  string
		adu   repo.AppDistributorUnit
		topic string
		wantM kafka.Message
	}{

		{
			name: "fileName absent, not last",
			adu: repo.AppDistributorUnit{
				H: repo.AppDistributorHeader{TS: "qqq", FormName: "alice"},
				B: repo.AppDistributorBody{B: []byte("azaza")},
			},
			wantM: kafka.Message{Key: []byte("qqq"), Value: []byte("\n\x03qqq\x12\x05alice\"\x05azaza")},
		},

		{
			name: "fileName absent, last",
			adu: repo.AppDistributorUnit{
				H: repo.AppDistributorHeader{TS: "qqq", FormName: "alice", IsLast: true},
				B: repo.AppDistributorBody{B: []byte("azaza")},
			},
			wantM: kafka.Message{Key: []byte("qqq"), Value: []byte("\n\x03qqq\x12\x05alice\"\x05azaza(\x01")},
		},

		{
			name: "fileName present, not last",
			adu: repo.AppDistributorUnit{
				H: repo.AppDistributorHeader{TS: "qqq", FormName: "alice", FileName: "short.txt"},
				B: repo.AppDistributorBody{B: []byte("azaza")},
			},
			wantM: kafka.Message{Key: []byte("qqq"), Value: []byte("\n\x03qqq\x12\x05alice\x1a\tshort.txt\"\x05azaza")},
		},

		{
			name: "fileName present, last",
			adu: repo.AppDistributorUnit{
				H: repo.AppDistributorHeader{TS: "qqq", FormName: "alice", FileName: "short.txt", IsLast: true},
				B: repo.AppDistributorBody{B: []byte("azaza")},
			},
			wantM: kafka.Message{Key: []byte("qqq"), Value: []byte("\n\x03qqq\x12\x05alice\x1a\tshort.txt\"\x05azaza(\x01")},
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			m, _ := GenMessage(v.adu, v.topic)
			//logger.L.Infof("m = %q\n", m)
			s.Equal(v.wantM, m)
		})
	}
}
