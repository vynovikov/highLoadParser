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

func (s *byteOpsSuite) TestSameByteSlices() {

	tt := []struct {
		name string
		a    []byte
		b    []byte
	}{
		{
			name: "1. Short",
			a:    []byte("aaa"),
			b:    []byte("aaa"),
		},

		{
			name: "2. Composed",
			a:    append([]byte("aaa"), []byte("bbb")...),
			b:    []byte("aaabbb"),
		},

		{
			name: "3. Real boundary",
			a:    append([]byte("\r\n-----"), []byte("-----bRoot\r\n")[:10]...),
			b:    []byte("\r\n----------bRoot"),
		},
	}
	for _, v := range tt {

		s.True(SameByteSlices(v.a, v.b))
	}
}
