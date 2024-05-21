package repository

import (
	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/logger"
)

type ParserRepository interface {
	Register([]dataHandler.DataHandlerDTO, dataHandler.Boundary) error
}

type repositoryStruct struct {
	dataHandler dataHandler.DataHandler
}

func NewParserRepository(dh dataHandler.DataHandler) *repositoryStruct {
	return &repositoryStruct{
		dataHandler: dh,
	}
}

func (r *repositoryStruct) Register(dtos []dataHandler.DataHandlerDTO, bou dataHandler.Boundary) error {

	for _, d := range dtos {

		d = dataHandler.NewDataHandlerUnit(d)

		switch {

		case d.B() == dataHandler.False:

			err := r.dataHandler.Create(d, bou)
			if err != nil {

				logger.L.Infof("in repository.Register unable to create %s %d: %v\n", d.TS(), d.Part(), err)
			}

		case d.B() == dataHandler.True:

			err := r.dataHandler.Updade(d, bou)
			if err != nil {

				logger.L.Infof("in repository.Register unable to update %s %d: %v\n", d.TS(), d.Part(), err)
			}
		}

		if d.Last() {

			err := r.dataHandler.Delete(d.TS())

			if err != nil {

				logger.L.Infof("in repository.Register unable to delete %s %v\n", d.TS(), err)
			}
		}
	}

	return nil
}
