package repository

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type repositorySuite struct {
	suite.Suite
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestIsLastBoundaryPart() {
	tt := []struct {
		name string
	}{}
	for _, v := range tt {
		s.Run(v.name, func() {})
	}
}
