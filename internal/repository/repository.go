package repository

import (
	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/logger"
)

type ParserRepository interface {
	Register(DataPiece) (TransferUnit, error)
	Check(DataPiece) (Presence, error)
}

type repositoryStruct struct {
	dataHandler dataHandler.DataHandler
}

func NewParserRepository(dh dataHandler.DataHandler) *repositoryStruct {
	return &repositoryStruct{
		dataHandler: dh,
	}
}

func (r *repositoryStruct) Register(d DataPiece) (TransferUnit, error) {
	logger.L.Infof("in repository.Register header: %s, body: %s\n", d.Header(), d.Body())

	return TransferUnit{}, nil
}

func (r *repositoryStruct) Check(d DataPiece) (Presence, error) {

	return Presence{}, nil
}
