package infrastructure

import (
	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/repository"
	"github.com/vynovikov/highLoadParser/internal/transmitters"
)

type Infrastructure interface {
	Register([]dataHandler.DataHandlerDTO, dataHandler.Boundary) error
	Send(TransferUnit) error
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
func (i *infrastructureStruct) Register(dtos []dataHandler.DataHandlerDTO, bou dataHandler.Boundary) error {

	return i.repo.Register(dtos, bou)
}

func (i *infrastructureStruct) Send(unit TransferUnit) error {

	i.transmitter.TransmitToSaver(unit)

	return nil
}
