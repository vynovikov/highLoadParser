package infrastructure

import (
	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/repository"
	"github.com/vynovikov/highLoadParser/internal/transmitters"
)

type Infrastructure interface {
	Check([]dataHandler.DataHandlerDTO) ([]dataHandler.Presence, []error)
	Send([]TransferUnit) []error
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

func (i *infrastructureStruct) Check(dtos []dataHandler.DataHandlerDTO) ([]dataHandler.Presence, []error) {

	res, errs := make([]dataHandler.Presence, 0, len(dtos)), make([]error, 0, len(dtos))

	for _, v := range dtos {

		presenceOne, err := i.checkOne(v)

		if err != nil {

			errs = append(errs, err)
		} else {

			res = append(res, presenceOne)
		}
	}

	return res, errs
}

func (i *infrastructureStruct) Send(units []TransferUnit) []error {

	errs := make([]error, 0, len(units))

	for _, v := range units {

		err := i.transmitter.TransmitToSaver(v)

		if err != nil {

			errs = append(errs, err)
		}
	}

	return errs
}

func (i *infrastructureStruct) checkOne(d dataHandler.DataHandlerDTO) (dataHandler.Presence, error) {

	return i.repo.Check(d)
}
