package infrastructure

import (
	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/repository"
	"github.com/vynovikov/highLoadParser/internal/transmitters"
)

type Infrastructure interface {
	Register(repository.RepositoryDTO) (dataHandler.ProducerUnit, error)
	Send(dataHandler.ProducerUnit) error
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
func (i *infrastructureStruct) Register(dto repository.RepositoryDTO) (dataHandler.ProducerUnit, error) {

	return i.repo.Register(dto)
}

func (i *infrastructureStruct) Send(unit dataHandler.ProducerUnit) error {

	return i.transmitter.TransmitToSaver(unit)
}
