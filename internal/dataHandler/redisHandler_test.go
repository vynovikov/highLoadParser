package dataHandler

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type redisHandlerSuite struct {
	suite.Suite
}

func TestRedisHandlerSuite(t *testing.T) {
	suite.Run(t, new(redisHandlerSuite))
}

func (s *redisHandlerSuite) TestSet() {
	tt := []struct {
		name string
	}{}
	for _, v := range tt {
		s.Run(v.name, func() {})
	}
}