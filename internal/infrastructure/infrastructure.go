package infrastructure

import (
	"github.com/vynovikov/highLoadParser/internal/repository"
	"github.com/vynovikov/highLoadParser/internal/transmitters"
)

type Infrastructure interface {
	Save([]DataPiece) ([]TransferUnit, []error)
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

func (i *infrastructureStruct) Save(datapieces []DataPiece) ([]TransferUnit, []error) {

	res := make([]TransferUnit, 0, len(datapieces))
	errs := make([]error, 0, len(datapieces))

	for _, v := range datapieces {
		_, err := i.repo.Register(v)

		if err != nil {

			errs = append(errs, err)
		} else {

			res = append(res, &TransferUnitStruct{})
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
