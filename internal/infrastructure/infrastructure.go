package infrastructure

import (
	"github.com/vynovikov/highLoadParser/internal/logger"
	"github.com/vynovikov/highLoadParser/internal/repository"
	"github.com/vynovikov/highLoadParser/internal/transmitters"
)

type Infrastructure interface {
	Save([]DataPiece) ([]TransferUnit, error)
	Send([]TransferUnit) error
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

func (i *infrastructureStruct) Save(datapieces []DataPiece) ([]TransferUnit, error) {

	res := make([]TransferUnit, len(datapieces))
	/*
		res, err := i.repo.Register(datapieces)
		if err != nil {

			return res, err
		}
	*/
	return res, nil
}

func (i *infrastructureStruct) Send(units []TransferUnit) error {

	for _, v := range units {

		logger.L.Infof("in infrastructure.Send trying to send %v\n", v)
	}

	return nil
}
