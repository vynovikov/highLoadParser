package byteOps

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type byteOpsSuite struct {
	suite.Suite
}

type boundary struct {
	prefix []byte
	root   []byte
	suffix []byte
}

func TestByteOps(t *testing.T) {
	suite.Run(t, new(byteOpsSuite))
}
