package repository

import (
	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/logger"
)

type ParserRepository interface {
	Register(RepositoryDTO) (TransferUnit, error)
	Check(RepositoryDTO) (dataHandler.Presence, error)
}

type repositoryStruct struct {
	dataHandler dataHandler.DataHandler
}

func NewParserRepository(dh dataHandler.DataHandler) *repositoryStruct {
	return &repositoryStruct{
		dataHandler: dh,
	}
}

func (r *repositoryStruct) Register(d RepositoryDTO) (TransferUnit, error) {

	logger.L.Infof("in repository.Register body: %s\n", d.Body())

	return TransferUnit{}, nil
}

func (r *repositoryStruct) Check(d RepositoryDTO) (dataHandler.Presence, error) {

	return r.dataHandler.Check(NewRepositoryDTOUnit(d))
}
