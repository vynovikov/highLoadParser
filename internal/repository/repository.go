package repository

import (
	"github.com/vynovikov/study/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/study/highLoadParser/internal/logger"
)

type ParserRepository interface {
	Register(DataPiece) error
}

type repositoryStruct struct {
	dataHandler dataHandler.DataHandler
}

func NewParserRepository(dh dataHandler.DataHandler) *repositoryStruct {
	return &repositoryStruct{
		dataHandler: dh,
	}
}

func (r *repositoryStruct) Register(d DataPiece) error {
	logger.L.Infof("in repository.Register header: %s, body: %s\n", d.Header(), d.Body())

	return nil
}
