package infrastructure

import (
	"github.com/vynovikov/highLoadParser/internal/repository"
	"github.com/vynovikov/highLoadParser/internal/transmitters"
)

type Infrastructure interface {
	Save([]DataPiece) ([]TransferUnit, []error)
	Check([]DataPiece) ([]Presence, []error)
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

func (i *infrastructureStruct) Check(datapieces []DataPiece) ([]Presence, []error) {

	res, errs := make([]Presence, 0, len(datapieces)), make([]error, 0, len(datapieces))

	for _, v := range datapieces {

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

func (i *infrastructureStruct) checkOne(d DataPiece) (Presence, error) {

	presense, err := i.repo.Check(d)
	if err != nil {

		return Presence{}, err
	}
	infPresence := newPresence(presense)

	return infPresence, nil

}

func newPresence(p repository.Presence) Presence {
	return Presence{}
}
