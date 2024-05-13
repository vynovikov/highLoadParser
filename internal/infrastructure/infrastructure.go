package infrastructure

import (
	"github.com/vynovikov/highLoadParser/internal/repository"
	"github.com/vynovikov/highLoadParser/internal/transmitters"
)

type Infrastructure interface {
	Save() error
	Send() error
}

type infrastructureStruct struct {
	repo        repository.ParserRepository
	transmitter transmitters.ParserTransmitter
}

func NewInfraStructure(repo repository.ParserRepository, transmitter transmitters.ParserTransmitter) *infrastructureStruct {

	return &infrastructureStruct{

		repo:        repo,
		transmitter: transmitter,
	}
}

func (i *infrastructureStruct) Save() error {

	return nil
}

func (i *infrastructureStruct) Send() error {

	return nil
}
